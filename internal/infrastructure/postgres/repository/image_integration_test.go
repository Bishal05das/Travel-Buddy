package repository_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/infrastructure/postgres/repository"
	agencyusecase "github.com/bishal05das/travelbuddy/internal/usecase/agency"
	homeusecase "github.com/bishal05das/travelbuddy/internal/usecase/home"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Requires a disposable database URL; does not load the application's .env.
func TestImageReplacementIntegration(t *testing.T) {
	dsn := os.Getenv("IMAGE_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set IMAGE_TEST_DATABASE_URL to a disposable PostgreSQL database URL")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	setup, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer setup.Close()
	setup.SetMaxOpenConns(1)
	schema := "image_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := setup.ExecContext(ctx, "CREATE SCHEMA "+pq.QuoteIdentifier(schema)); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := setup.Exec("DROP SCHEMA " + pq.QuoteIdentifier(schema) + " CASCADE"); err != nil {
			t.Error(err)
		}
	}()
	if _, err := setup.ExecContext(ctx, "SET search_path TO "+pq.QuoteIdentifier(schema)+", public"); err != nil {
		t.Fatal(err)
	}
	paths, err := filepath.Glob("../../../../migrations/*.up.sql")
	if err != nil || len(paths) == 0 {
		t.Fatalf("migrations: %v", err)
	}
	sort.Strings(paths)
	for _, path := range paths {
		sql, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := setup.ExecContext(ctx, string(sql)); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
	}
	// Every connection uses the isolated schema, allowing real concurrent writes.
	scopedURL, err := url.Parse(dsn)
	if err != nil || (scopedURL.Scheme != "postgres" && scopedURL.Scheme != "postgresql") {
		t.Fatal("IMAGE_TEST_DATABASE_URL must be a PostgreSQL URL")
	}
	query := scopedURL.Query()
	query.Set("search_path", schema+",public")
	scopedURL.RawQuery = query.Encode()
	db, err := sqlx.ConnectContext(ctx, "postgres", scopedURL.String())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(4)
	agencies, tours := repository.NewAgencyRepositoryDB(db), repository.NewTourRepositoryDB(db)
	agencyID, otherAgencyID := uuid.New(), uuid.New()
	if _, err := db.ExecContext(ctx, `INSERT INTO agency (agency_id,name) VALUES ($1,'Image Agency'),($2,'Other Agency')`, agencyID, otherAgencyID); err != nil {
		t.Fatal(err)
	}
	tour := &domain.Tour{AgencyID: agencyID, Name: "Image Tour", StartDate: time.Now().Add(48 * time.Hour), EndDate: time.Now().Add(72 * time.Hour), LastEnrollmentDate: time.Now().Add(24 * time.Hour), TotalSeat: 20, Price: 1000, Description: "An image replacement test", ImagePath: "tours/initial.png"}
	if err := tours.CreateTour(ctx, tour); err != nil {
		t.Fatal(err)
	}

	for _, target := range []struct {
		name    string
		update  func(string) (string, error)
		current func() (string, error)
		count   func() (int, error)
	}{
		{"agency", func(path string) (string, error) { return agencies.UpdateAgencyImage(ctx, agencyID, path) },
			func() (string, error) {
				a, err := agencies.GetAgency(ctx, agencyID)
				if err != nil {
					return "", err
				}
				return a.ImagePath, nil
			},
			func() (int, error) {
				var n int
				err := db.GetContext(ctx, &n, `SELECT COUNT(*) FROM agency_images WHERE agency_id=$1 AND is_active AND deleted_at IS NULL`, agencyID)
				return n, err
			}},
		{"tour", func(path string) (string, error) { return tours.UpdateTourImage(ctx, agencyID, tour.TourID, path) },
			func() (string, error) {
				t, err := tours.GetByID(ctx, tour.TourID)
				if err != nil {
					return "", err
				}
				return t.ImagePath, nil
			},
			func() (int, error) {
				var n int
				err := db.GetContext(ctx, &n, `SELECT COUNT(*) FROM tour_images WHERE tour_id=$1 AND is_active AND deleted_at IS NULL`, tour.TourID)
				return n, err
			}},
	} {
		t.Run(target.name, func(t *testing.T) {
			first := "images/" + target.name + "/first.png"
			second := "images/" + target.name + "/second.png"
			if _, err := target.update(first); err != nil {
				t.Fatal(err)
			}
			if old, err := target.update(second); err != nil || old != first {
				t.Fatalf("replace old=%q err=%v", old, err)
			}
			if got, err := target.current(); err != nil || got != second {
				t.Fatalf("public image=%q err=%v", got, err)
			}
			// Failure after deactivating the old record must roll the transaction back.
			if _, err := target.update(strings.Repeat("x", 301)); err == nil {
				t.Fatal("expected invalid database path to fail")
			}
			if got, err := target.current(); err != nil || got != second {
				t.Fatalf("rollback image=%q err=%v", got, err)
			}
			var workers sync.WaitGroup
			failures := make(chan error, 8)
			for i := 0; i < 8; i++ {
				workers.Add(1)
				go func() {
					defer workers.Done()
					_, err := target.update("images/" + target.name + "/" + uuid.NewString() + ".png")
					if err != nil {
						failures <- err
					}
				}()
			}
			workers.Wait()
			close(failures)
			for err := range failures {
				t.Error(err)
			}
			if count, err := target.count(); err != nil || count != 1 {
				t.Fatalf("active images=%d err=%v", count, err)
			}
		})
	}
	if _, err := tours.UpdateTourImage(ctx, otherAgencyID, tour.TourID, "foreign.png"); !errors.Is(err, domain.ErrImageTargetNotFound) {
		t.Fatalf("cross agency update: %v", err)
	}
	if _, err := agencies.UpdateAgencyImage(ctx, uuid.New(), "missing.png"); !errors.Is(err, domain.ErrImageTargetNotFound) {
		t.Fatalf("missing agency: %v", err)
	}
	imageUC := agencyusecase.NewUpdateAgencyImageUseCase(agencies)
	for _, actor := range []domain.Actor{{Role: domain.RoleUser, AgencyID: &agencyID}, {Role: domain.RoleMember}, {Role: domain.RoleMember, AgencyID: &otherAgencyID}} {
		if _, err := imageUC.Execute(ctx, actor, agencyID, "forbidden.png"); !errors.Is(err, domain.ErrImageAccessDenied) {
			t.Fatalf("agency scope check: %v", err)
		}
	}
	t.Run("home includes current agency images without inflating tour counts", func(t *testing.T) {
		current, err := agencies.GetAgency(ctx, agencyID)
		if err != nil {
			t.Fatal(err)
		}
		// A deleted image can still have is_active=true; it must be ignored.
		if _, err := db.ExecContext(ctx, `INSERT INTO agency_images (agency_id, image_path, is_active, deleted_at) VALUES ($1, 'images/agencies/deleted.png', TRUE, CURRENT_TIMESTAMP)`, agencyID); err != nil {
			t.Fatal(err)
		}
		inactiveID := uuid.New()
		if _, err := db.ExecContext(ctx, `INSERT INTO agency (agency_id, name, is_active, rating) VALUES ($1, 'Inactive Agency', FALSE, 5)`, inactiveID); err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/home", nil).WithContext(ctx)
		handler.NewHomeHandler(homeusecase.NewHomeUseCase(repository.NewHomeRepositoryDB(db))).GetHome(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("home response: %d %s", rec.Code, rec.Body.String())
		}
		var response struct {
			Success bool                `json:"success"`
			Data    domain.HomeResponse `json:"data"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if !response.Success || len(response.Data.TopAgencies) != 2 {
			t.Fatalf("expected two active agencies: %+v", response)
		}
		for _, agency := range response.Data.TopAgencies {
			switch agency.AgencyID {
			case agencyID:
				if agency.ImagePath != current.ImagePath || agency.TotalTours != 1 {
					t.Fatalf("current image or tour count mismatch: %+v", agency)
				}
			case otherAgencyID:
				if agency.ImagePath != "" || agency.TotalTours != 0 {
					t.Fatalf("agency without an image or tours: %+v", agency)
				}
			default:
				t.Fatalf("unexpected agency: %+v", agency)
			}
		}
	})
	if err := tours.UpdateTourStatus(ctx, tour.TourID, "cancelled", &agencyID); err != nil {
		t.Fatal(err)
	}
	if _, err := tours.UpdateTourImage(ctx, agencyID, tour.TourID, "cancelled.png"); !errors.Is(err, domain.ErrTourCancelled) {
		t.Fatalf("cancelled update: %v", err)
	}
}
