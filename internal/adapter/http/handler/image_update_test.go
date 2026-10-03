package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	"github.com/bishal05das/travelbuddy/internal/domain"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

type agencyImageFunc func(context.Context, domain.Actor, uuid.UUID, string) (string, error)

func (f agencyImageFunc) Execute(ctx context.Context, actor domain.Actor, agencyID uuid.UUID, path string) (string, error) {
	return f(ctx, actor, agencyID, path)
}

type tourImageFunc func(context.Context, domain.Actor, uuid.UUID, uuid.UUID, string) (string, error)

func (f tourImageFunc) Execute(ctx context.Context, actor domain.Actor, agencyID, tourID uuid.UUID, path string) (string, error) {
	return f(ctx, actor, agencyID, tourID, path)
}

func TestImageReplacementUploadsAndCleanup(t *testing.T) {
	t.Chdir(t.TempDir())
	agencyID, tourID, memberID := uuid.New(), uuid.New(), uuid.New()
	for _, target := range []string{"agencies", "tours"} {
		for _, tc := range []struct {
			name   string
			err    error
			status int
		}{
			{"success", nil, 200}, {"database failure", errors.New("private database detail"), 500},
			{"wrong agency", domain.ErrImageAccessDenied, 403}, {"missing target", domain.ErrImageTargetNotFound, 404},
			{"cancelled tour", domain.ErrTourCancelled, 409},
		} {
			t.Run(target+"/"+tc.name, func(t *testing.T) {
				dir := filepath.Join("images", target)
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatal(err)
				}
				oldPath := filepath.Join(dir, "old.png")
				if err := os.WriteFile(oldPath, pngBytes, 0600); err != nil {
					t.Fatal(err)
				}
				defer os.Remove(oldPath)
				var newPath string
				update := func(ctx context.Context, actor domain.Actor, id uuid.UUID, path string) (string, error) {
					if id != agencyID || actor.ID != memberID || actor.AgencyID == nil || *actor.AgencyID != agencyID {
						t.Fatalf("incorrect identity or scope: %+v, %s", actor, id)
					}
					newPath = path
					if got, err := os.ReadFile(path); err != nil || string(got) != string(pngBytes) {
						t.Fatalf("new image not saved before DB update: %v", err)
					}
					return oldPath, tc.err
				}
				var serve http.HandlerFunc
				if target == "agencies" {
					serve = handler.NewAgencyHandler(nil, nil, nil, nil, agencyImageFunc(update)).UpdateImage
				} else {
					serve = handler.NewTourHandler(nil, nil, nil, nil, nil, nil, tourImageFunc(func(ctx context.Context, actor domain.Actor, agency, tour uuid.UUID, path string) (string, error) {
						if tour != tourID {
							t.Fatal("incorrect tour scope")
						}
						return update(ctx, actor, agency, path)
					})).UpdateImage
				}
				req := newAgencyForm(t, nil, "new.png", pngBytes)
				req.SetPathValue("agency_id", agencyID.String())
				req.SetPathValue("tour_id", tourID.String())
				req = req.WithContext(util.WithPayload(req.Context(), &util.Payload{UserID: memberID, Role: domain.RoleMember, AgencyID: &agencyID}))
				rec := httptest.NewRecorder()
				serve(rec, req)
				if rec.Code != tc.status {
					t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
				}
				_, oldErr := os.Stat(oldPath)
				_, newErr := os.Stat(newPath)
				if tc.err == nil {
					if !errors.Is(oldErr, os.ErrNotExist) || newErr != nil {
						t.Fatalf("replacement cleanup: old=%v new=%v", oldErr, newErr)
					}
					var result struct {
						ImagePath string `json:"image_path"`
						ImageURL  string `json:"image_url"`
					}
					if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
						t.Fatal(err)
					}
					if result.ImagePath != newPath || result.ImageURL != "/"+newPath {
						t.Fatalf("wrong image URLs: %+v", result)
					}
					os.Remove(newPath)
				} else if oldErr != nil || !errors.Is(newErr, os.ErrNotExist) {
					t.Fatalf("failure cleanup: old=%v new=%v", oldErr, newErr)
				}
			})
		}
	}
}

func TestImageReplacementRejectsInvalidUploads(t *testing.T) {
	t.Chdir(t.TempDir())
	h := handler.NewAgencyHandler(nil, nil, nil, nil, agencyImageFunc(func(context.Context, domain.Actor, uuid.UUID, string) (string, error) {
		t.Fatal("invalid upload reached database")
		return "", nil
	}))
	for _, tc := range []struct {
		name, filename string
		content        []byte
		status         int
	}{
		{"missing", "", nil, 400}, {"empty", "empty.png", nil, 400},
		{"wrong extension", "image.svg", pngBytes, 400}, {"wrong content", "script.png", []byte("<script>alert(1)</script>"), 400},
		{"over size limit", "large.png", make([]byte, (10<<20)+1), 413},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := newAgencyForm(t, nil, tc.filename, tc.content)
			req.SetPathValue("agency_id", uuid.NewString())
			req = req.WithContext(util.WithPayload(req.Context(), &util.Payload{Role: domain.RoleSuper}))
			rec := httptest.NewRecorder()
			h.UpdateImage(rec, req)
			if rec.Code != tc.status {
				t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
	rec := httptest.NewRecorder()
	h.UpdateImage(rec, newAgencyForm(t, nil, "new.png", pngBytes))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("guest status: %d", rec.Code)
	}
}
