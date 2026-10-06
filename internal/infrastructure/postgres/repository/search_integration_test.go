package repository_test

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/infrastructure/postgres/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Uses only an explicitly supplied disposable database, never the app's .env.
// All migrations, fixtures and queries run in a fresh, isolated schema.
func TestSearchIntegration(t *testing.T) {
	dsn := os.Getenv("SEARCH_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set SEARCH_TEST_DATABASE_URL to a disposable PostgreSQL database URL")
	}
	databaseURL, err := url.Parse(dsn)
	if err != nil || (databaseURL.Scheme != "postgres" && databaseURL.Scheme != "postgresql") {
		t.Fatal("SEARCH_TEST_DATABASE_URL must be a PostgreSQL URL")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		t.Fatalf("connect to search test database: %v", err)
	}
	defer db.Close()
	// Keep SET search_path valid for every query and repository transaction.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	schema := "search_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	exec := func(t *testing.T, query string, args ...any) {
		t.Helper()
		if _, err := db.ExecContext(ctx, query, args...); err != nil {
			t.Fatalf("search fixture SQL: %v", err)
		}
	}
	exec(t, "CREATE SCHEMA "+pq.QuoteIdentifier(schema))
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		if _, err := db.ExecContext(cleanupCtx, "DROP SCHEMA "+pq.QuoteIdentifier(schema)+" CASCADE"); err != nil {
			t.Errorf("remove isolated search schema: %v", err)
		}
	}()
	exec(t, "SET search_path TO "+pq.QuoteIdentifier(schema)+", public")
	paths, err := filepath.Glob("../../../../migrations/*.up.sql")
	if err != nil || len(paths) == 0 {
		t.Fatalf("find search-test migrations: count=%d err=%v", len(paths), err)
	}
	sort.Strings(paths)
	for _, path := range paths {
		migration, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read migration %s: %v", path, err)
		}
		if _, err := db.ExecContext(ctx, string(migration)); err != nil {
			t.Fatalf("apply migration %s: %v", path, err)
		}
	}

	repo := repository.NewSearchRepository(db)
	baseDate := time.Date(2040, time.January, 20, 0, 0, 0, 0, time.UTC)
	created := time.Date(2030, time.January, 1, 0, 0, 0, 0, time.UTC)
	reset := func(t *testing.T) {
		t.Helper()
		exec(t, `DELETE FROM tour_images; DELETE FROM tours; DELETE FROM agency_images; DELETE FROM agency;`)
		created = time.Date(2030, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	insertAgency := func(t *testing.T, name string, active bool) uuid.UUID {
		t.Helper()
		id := uuid.New()
		created = created.Add(time.Hour)
		exec(t, `INSERT INTO agency (agency_id, name, address, is_active, created_at)
			VALUES ($1, $2, 'Dhaka, Bangladesh', $3, $4)`, id, name, active, created)
		return id
	}
	insertTour := func(t *testing.T, agency uuid.UUID, name, description string, price, discount int, start time.Time, status string) uuid.UUID {
		t.Helper()
		id := uuid.New()
		created = created.Add(time.Hour)
		exec(t, `INSERT INTO tours
			(tour_id, agency_id, name, description, price, discount, start_date, end_date,
			last_enrollment_date, total_seat, available_seat, status, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 20, 20, $10, $11)`,
			id, agency, name, description, price, discount, start, start.AddDate(0, 0, 2), start.AddDate(0, 0, -1), status, created)
		return id
	}
	searchTours := func(t *testing.T, filter domain.TourSearchFilter) []domain.TourSearchResponse {
		t.Helper()
		if filter.Limit == 0 {
			filter.Limit = 100
		}
		results, err := repo.SearchTours(ctx, filter)
		if err != nil {
			t.Fatalf("search tours with query %q: %v", filter.Query, err)
		}
		return results
	}
	assertTours := func(t *testing.T, results []domain.TourSearchResponse, want ...uuid.UUID) {
		t.Helper()
		got := make([]uuid.UUID, 0, len(results))
		for _, result := range results {
			got = append(got, result.TourID)
		}
		if len(got) != len(want) {
			t.Fatalf("tour IDs=%v, want %v", got, want)
		}
		wantSet := make(map[uuid.UUID]bool, len(want))
		for _, id := range want {
			wantSet[id] = true
		}
		for _, id := range got {
			if !wantSet[id] {
				t.Fatalf("unexpected or duplicate tour %s in %v, want %v", id, got, want)
			}
			delete(wantSet, id)
		}
	}
	assertCount := func(t *testing.T, filter domain.TourSearchFilter, want int) {
		t.Helper()
		got, err := repo.CountTours(ctx, filter)
		if err != nil || got != want {
			t.Fatalf("count tours with query %q=%d err=%v, want %d", filter.Query, got, err, want)
		}
	}

	t.Run("typos case and extra whitespace", func(t *testing.T) {
		reset(t)
		agency := insertAgency(t, "Compass Adventures", true)
		sajek := insertTour(t, agency, "Sajek Valley Tour", "Mountain scenery and local food", 1000, 0, baseDate, domain.TourOpen)
		insertTour(t, agency, "Coastal Retreat", "Ocean scenery and sailing", 1000, 0, baseDate, domain.TourOpen)
		for _, query := range []string{"sajek", "sajeek", "sajke", "  SAJEK \t  VALLEY \n TOUR  "} {
			filter := domain.TourSearchFilter{Query: query}
			assertTours(t, searchTours(t, filter), sajek)
			assertCount(t, filter, 1)
		}
		assertTours(t, searchTours(t, domain.TourSearchFilter{Query: "qzxvbnmpl"}))
		assertCount(t, domain.TourSearchFilter{Query: "qzxvbnmpl"}, 0)
	})

	t.Run("all query words can match different fields in any order", func(t *testing.T) {
		reset(t)
		agency := insertAgency(t, "Compass Adventures", true)
		weekend := insertTour(t, agency, "Sajek Valley Tour", "Weekend trekking with a local guide", 1000, 0, baseDate, domain.TourOpen)
		insertTour(t, agency, "Sajek Holiday", "A relaxing mountain holiday", 1000, 0, baseDate, domain.TourOpen)
		for _, query := range []string{"weekend sajek", "sajek weekend", "weekend sajeek", "trekking compass"} {
			filter := domain.TourSearchFilter{Query: query}
			assertTours(t, searchTours(t, filter), weekend)
			assertCount(t, filter, 1)
		}
		assertTours(t, searchTours(t, domain.TourSearchFilter{Query: "sajek submarine"}))
	})

	t.Run("title relevance outweighs recency description and typo matches", func(t *testing.T) {
		reset(t)
		agency := insertAgency(t, "Compass Adventures", true)
		exact := insertTour(t, agency, "Sajek", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		prefix := insertTour(t, agency, "Sajek Valley Tour", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		substring := insertTour(t, agency, "Weekend in Sajek Valley", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		fuzzy := insertTour(t, agency, "Sajke Valley Tour", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		description := insertTour(t, agency, "Riverside Picnic", "Discover Sajek with a local guide", 1000, 0, baseDate, domain.TourOpen)
		results := searchTours(t, domain.TourSearchFilter{Query: "sajek"})
		assertTours(t, results, exact, prefix, substring, fuzzy, description)
		positions := make(map[uuid.UUID]int, len(results))
		for i, result := range results {
			positions[result.TourID] = i
		}
		if positions[exact] != 0 || positions[prefix] >= positions[substring] || positions[substring] >= positions[fuzzy] || positions[substring] >= positions[description] {
			t.Fatalf("expected exact, prefix and substring titles ahead of fuzzy/description matches: %+v", results)
		}
	})

	t.Run("wildcards and escape characters are literal and short queries stay conservative", func(t *testing.T) {
		reset(t)
		agency := insertAgency(t, "Compass Adventures", true)
		percent := insertTour(t, agency, "100% Coastal Tour", "Ocean scenery", 1000, 0, baseDate, domain.TourOpen)
		underscore := insertTour(t, agency, "Photo_Adventure", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		backslash := insertTour(t, agency, `Trail\Compass`, "Forest scenery", 1000, 0, baseDate, domain.TourOpen)
		sajek := insertTour(t, agency, "Sajek Valley Tour", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		insertTour(t, agency, "Sunrise Walk", "Ocean scenery", 1000, 0, baseDate, domain.TourOpen)
		for _, example := range []struct {
			query string
			want  uuid.UUID
		}{{"%", percent}, {"_", underscore}, {`\`, backslash}, {"sa", sajek}} {
			filter := domain.TourSearchFilter{Query: example.query}
			assertTours(t, searchTours(t, filter), example.want)
			assertCount(t, filter, 1)
		}
	})

	t.Run("fuzzy search respects dates and integer discounted prices", func(t *testing.T) {
		reset(t)
		agency := insertAgency(t, "Compass Adventures", true)
		// Matches Tour.UnitPrice: 1001 - floor(1001*33/100) = 671.
		discounted := insertTour(t, agency, "Sajek Valley Tour", "Mountain scenery", 1001, 33, baseDate, domain.TourOpen)
		insertTour(t, agency, "Sajek Premium Tour", "Mountain scenery", 2200, 0, baseDate, domain.TourOpen)
		insertTour(t, agency, "Sajek Budget Tour", "Mountain scenery", 600, 0, baseDate, domain.TourOpen)
		insertTour(t, agency, "Sajek February Tour", "Mountain scenery", 1001, 33, baseDate.AddDate(0, 1, 0), domain.TourOpen)
		insertTour(t, agency, "Sajek Earlier Tour", "Mountain scenery", 1001, 33, baseDate.AddDate(0, 0, -1), domain.TourOpen)
		min, max := float64(671), float64(671)
		start, end := baseDate, baseDate.AddDate(0, 0, 2)
		filter := domain.TourSearchFilter{Query: "sajeek", MinPrice: &min, MaxPrice: &max, StartDate: &start, EndDate: &end}
		assertTours(t, searchTours(t, filter), discounted)
		assertCount(t, filter, 1)
		min, max = 670, 670
		assertTours(t, searchTours(t, filter))
		assertCount(t, filter, 0)
		// HTTP filters accept decimal values even though payable tour prices are integers.
		min, max = 670.5, 671.5
		assertTours(t, searchTours(t, filter), discounted)
		assertCount(t, filter, 1)
		min, max = 671.5, 672
		assertTours(t, searchTours(t, filter))
		assertCount(t, filter, 0)
		min, max = 600, 700
		filter.MinPrice = nil
		// The 1001 list price must not exclude its discounted tour.
		results := searchTours(t, filter)
		if len(results) != 2 {
			t.Fatalf("expected discounted and budget tours below 700: %+v", results)
		}
		assertCount(t, filter, 2)
		end = baseDate.AddDate(0, 0, 1)
		assertTours(t, searchTours(t, filter))
		assertCount(t, filter, 0)
	})

	t.Run("discount calculation does not overflow PostgreSQL integer multiplication", func(t *testing.T) {
		reset(t)
		agency := insertAgency(t, "Compass Adventures", true)
		tour := insertTour(t, agency, "Sajek Premium Tour", "Mountain scenery", 2_000_000_000, 99, baseDate, domain.TourOpen)
		price := float64(20_000_000)
		filter := domain.TourSearchFilter{Query: "sajke", MinPrice: &price, MaxPrice: &price}
		assertTours(t, searchTours(t, filter), tour)
		assertCount(t, filter, 1)
	})

	t.Run("statuses remain visible unless explicitly filtered and agency filters apply", func(t *testing.T) {
		reset(t)
		agency := insertAgency(t, "Compass Adventures", true)
		other := insertAgency(t, "Harbor Trips", true)
		open := insertTour(t, agency, "Sajek Valley Tour", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		closed := insertTour(t, agency, "Sajek Closed Tour", "Mountain scenery", 1000, 0, baseDate, domain.TourClosed)
		cancelled := insertTour(t, agency, "Sajek Cancelled Tour", "Mountain scenery", 1000, 0, baseDate, domain.TourCancelled)
		otherTour := insertTour(t, other, "Sajek Other Tour", "Mountain scenery", 1000, 0, baseDate, domain.TourClosed)
		assertTours(t, searchTours(t, domain.TourSearchFilter{Query: "sajek"}), open, closed, cancelled, otherTour)
		status := domain.TourClosed
		filter := domain.TourSearchFilter{Query: "sajke", AgencyID: &agency, Status: &status}
		assertTours(t, searchTours(t, filter), closed)
		assertCount(t, filter, 1)
	})

	t.Run("current images and pagination preserve the matching total", func(t *testing.T) {
		reset(t)
		agency := insertAgency(t, "Compass Adventures", true)
		first := insertTour(t, agency, "Sajek", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		second := insertTour(t, agency, "Sajek Valley Tour", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		third := insertTour(t, agency, "Weekend in Sajek Valley", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		exec(t, `INSERT INTO tour_images (tour_id, image_path, is_active, deleted_at) VALUES
			($1, 'images/tours/old.png', FALSE, NULL),
			($1, 'images/tours/deleted.png', TRUE, CURRENT_TIMESTAMP),
			($1, 'images/tours/current.png', TRUE, NULL)`, first)
		all := searchTours(t, domain.TourSearchFilter{Query: "sajek"})
		assertTours(t, all, first, second, third)
		for _, result := range all {
			wantImage := ""
			if result.TourID == first {
				wantImage = "images/tours/current.png"
			}
			if result.ImagePath != wantImage {
				t.Errorf("tour %s image=%q, want %q", result.TourID, result.ImagePath, wantImage)
			}
		}
		for offset, expected := range all {
			filter := domain.TourSearchFilter{Query: "sajek", Limit: 1, Offset: offset}
			page := searchTours(t, filter)
			assertTours(t, page, expected.TourID)
			assertCount(t, filter, 3)
		}
		filter := domain.TourSearchFilter{Query: "sajek", Limit: 1, Offset: 3}
		assertTours(t, searchTours(t, filter))
		assertCount(t, filter, 3)
	})

	t.Run("inactive agencies and their tours are excluded", func(t *testing.T) {
		reset(t)
		active := insertAgency(t, "Sajek Travels", true)
		inactive := insertAgency(t, "Sajek Hidden Agency", false)
		visible := insertTour(t, active, "Sajek Valley Tour", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		insertTour(t, inactive, "Sajek Secret Tour", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		for _, query := range []string{"sajek", "sajeek", ""} {
			filter := domain.TourSearchFilter{Query: query}
			assertTours(t, searchTours(t, filter), visible)
			assertCount(t, filter, 1)
			agencies, err := repo.SearchAgencies(ctx, query, 100)
			if query == "" {
				if err != nil || len(agencies) != 0 {
					t.Fatalf("blank agency search results=%+v err=%v", agencies, err)
				}
				continue
			}
			if err != nil || len(agencies) != 1 || agencies[0].AgencyID != active || !agencies[0].IsActive {
				t.Fatalf("active agency search query=%q results=%+v err=%v", query, agencies, err)
			}
		}
	})

	t.Run("agency typo ranking normalization literal queries and current image", func(t *testing.T) {
		reset(t)
		exact := insertAgency(t, "Sajek", true)
		prefix := insertAgency(t, "Sajek Adventures", true)
		fuzzy := insertAgency(t, "Sajke Travels", true)
		insertAgency(t, "Harbor Excursions", true)
		exec(t, `INSERT INTO agency_images (agency_id, image_path, is_active, deleted_at) VALUES
			($1, 'images/agencies/old.png', FALSE, NULL),
			($1, 'images/agencies/deleted.png', TRUE, CURRENT_TIMESTAMP),
			($1, 'images/agencies/current.png', TRUE, NULL)`, exact)
		agencies, err := repo.SearchAgencies(ctx, "sajek", 100)
		if err != nil {
			t.Fatal(err)
		}
		if len(agencies) != 3 || agencies[0].AgencyID != exact || agencies[1].AgencyID != prefix || agencies[2].AgencyID != fuzzy {
			t.Fatalf("expected exact, prefix, fuzzy agency ranking: %+v", agencies)
		}
		if agencies[0].ImagePath != "images/agencies/current.png" || agencies[0].Address != "Dhaka, Bangladesh" || agencies[1].ImagePath != "" {
			t.Fatalf("agency image or public fields incorrect: %+v", agencies)
		}
		for _, query := range []string{"sajeek", "sajke"} {
			results, err := repo.SearchAgencies(ctx, query, 100)
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for _, agency := range results {
				found = found || agency.AgencyID == prefix
			}
			if !found {
				t.Fatalf("agency typo query %q misses Sajek Adventures: %+v", query, results)
			}
		}
		results, err := repo.SearchAgencies(ctx, " \t SAJEK   ADVENTURES \n", 100)
		if err != nil || len(results) == 0 || results[0].AgencyID != prefix {
			t.Fatalf("normalized agency search: %+v err=%v", results, err)
		}
		results, err = repo.SearchAgencies(ctx, "qzxvbnmpl", 100)
		if err != nil || len(results) != 0 {
			t.Fatalf("unrelated agency query: %+v err=%v", results, err)
		}
		percent := insertAgency(t, "100% Travel", true)
		underscore := insertAgency(t, "Photo_Travel", true)
		for _, example := range []struct {
			query string
			want  uuid.UUID
		}{{"%", percent}, {"_", underscore}} {
			results, err := repo.SearchAgencies(ctx, example.query, 100)
			if err != nil || len(results) != 1 || results[0].AgencyID != example.want {
				t.Fatalf("literal agency query %q: %+v err=%v", example.query, results, err)
			}
		}
	})

	t.Run("fuzzy thresholds do not leak to the pooled connection", func(t *testing.T) {
		reset(t)
		agency := insertAgency(t, "Compass Adventures", true)
		insertTour(t, agency, "Sajek Valley Tour", "Mountain scenery", 1000, 0, baseDate, domain.TourOpen)
		thresholds := func() []string {
			t.Helper()
			var values []string
			for _, setting := range []string{"pg_trgm.word_similarity_threshold", "pg_trgm.strict_word_similarity_threshold"} {
				var value string
				if err := db.GetContext(ctx, &value, "SELECT current_setting($1)", setting); err != nil {
					t.Fatal(err)
				}
				values = append(values, value)
			}
			return values
		}
		before := thresholds()
		searchTours(t, domain.TourSearchFilter{Query: "sajke"})
		assertCount(t, domain.TourSearchFilter{Query: "sajke"}, 1)
		if _, err := repo.SearchAgencies(ctx, "sajeek", 100); err != nil {
			t.Fatal(err)
		}
		if after := thresholds(); !reflect.DeepEqual(before, after) {
			t.Fatalf("fuzzy threshold leaked across pooled queries: before=%v after=%v", before, after)
		}
	})
}
