package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	"github.com/bishal05das/travelbuddy/internal/domain"
	mocks "github.com/bishal05das/travelbuddy/internal/mocks/usecase"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

func TestOwnerProtectionHandlerReturnsForbiddenForEveryCaller(t *testing.T) {
	agencyID, ownerID := uuid.New(), uuid.New()
	for _, caller := range []struct {
		name    string
		payload *util.Payload
	}{
		{"staff", &util.Payload{UserID: uuid.New(), Role: domain.RoleMember, AgencyID: &agencyID}},
		{"owner self", &util.Payload{UserID: ownerID, Role: domain.RoleMember, AgencyID: &agencyID}},
		{"platform admin", &util.Payload{UserID: uuid.New(), Role: domain.RoleSuper}},
	} {
		for _, action := range []string{"delete", "change permissions"} {
			t.Run(caller.name+"/"+action, func(t *testing.T) {
				called := false
				checkTarget := func(actor domain.Actor, target uuid.UUID) error {
					called = true
					if actor.ID != caller.payload.UserID || actor.Role != caller.payload.Role || target != ownerID {
						t.Fatalf("unexpected caller or target: actor=%+v target=%s", actor, target)
					}
					return fmt.Errorf("protected owner: %w", domain.ErrOwnerProtected)
				}
				deleteUC := &mocks.MockDeleteMember{ExecuteFunc: func(_ context.Context, actor domain.Actor, target uuid.UUID) error {
					return checkTarget(actor, target)
				}}
				updateUC := &mocks.MockUpdatePermission{ExecuteFunc: func(_ context.Context, actor domain.Actor, target uuid.UUID, _ *domain.UpdatePermissionRequest) error {
					return checkTarget(actor, target)
				}}
				h := handler.NewMemberHandler(nil, deleteUC, nil, updateUC, nil, nil)
				req := httptest.NewRequest(http.MethodDelete, "/members/"+ownerID.String(), nil)
				serve := h.DeleteMember
				if action == "change permissions" {
					req = httptest.NewRequest(http.MethodPut, "/members/"+ownerID.String()+"/permissions", strings.NewReader(`{"permissions":[]}`))
					serve = h.UpdateMemberPermissions
				}
				req.SetPathValue("member_id", ownerID.String())
				req = req.WithContext(util.WithPayload(req.Context(), caller.payload))
				rec := httptest.NewRecorder()
				serve(rec, req)
				if !called {
					t.Fatal("request did not reach the owner protection check")
				}
				if rec.Code != http.StatusForbidden {
					t.Fatalf("expected 403, got %d: %s", rec.Code, rec.Body.String())
				}
			})
		}
	}
}

type memberProfileUseCaseFunc func(context.Context, domain.Actor) (*domain.MemberProfile, error)

func (f memberProfileUseCaseFunc) Execute(ctx context.Context, actor domain.Actor) (*domain.MemberProfile, error) {
	return f(ctx, actor)
}

func TestMyProfileUsesAuthenticatedIdentity(t *testing.T) {
	for _, isOwner := range []bool{false, true} {
		t.Run(fmt.Sprintf("owner=%t", isOwner), func(t *testing.T) {
			memberID, agencyID := uuid.New(), uuid.New()
			want := domain.MemberProfile{
				MemberID: memberID, AgencyID: agencyID, AgencyName: "Test Agency",
				Name: "Signed In Member", Email: "member@example.com", Phone: "+8801700000000",
				RoleName: "Manager", IsOwner: isOwner,
			}
			called := false
			profileUC := memberProfileUseCaseFunc(func(_ context.Context, actor domain.Actor) (*domain.MemberProfile, error) {
				called = true
				if actor.ID != memberID || actor.Role != domain.RoleMember || actor.AgencyID == nil || *actor.AgencyID != agencyID {
					t.Fatalf("profile lookup must use verified claims: %+v", actor)
				}
				return &want, nil
			})
			h := handler.NewMemberHandler(nil, nil, nil, nil, nil, profileUC)
			req := httptest.NewRequest(http.MethodGet, "/members/me?member_id="+uuid.NewString()+"&agency_id="+uuid.NewString(), nil)
			req.SetPathValue("member_id", uuid.NewString())
			req.SetPathValue("agency_id", uuid.NewString())
			req = req.WithContext(util.WithPayload(req.Context(), &util.Payload{UserID: memberID, Role: domain.RoleMember, AgencyID: &agencyID}))
			rec := httptest.NewRecorder()
			h.GetMyProfile(rec, req)
			if !called || rec.Code != http.StatusOK {
				t.Fatalf("expected profile lookup and 200, called=%t status=%d body=%s", called, rec.Code, rec.Body.String())
			}
			var got domain.MemberProfile
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Fatalf("profile mismatch: got %+v, want %+v", got, want)
			}
		})
	}
}

func TestMyProfileRejectsUnauthenticatedAndNonMemberCallers(t *testing.T) {
	for _, tc := range []struct {
		name    string
		payload *util.Payload
		status  int
	}{
		{"guest", nil, http.StatusUnauthorized},
		{"customer", &util.Payload{UserID: uuid.New(), Role: domain.RoleUser}, http.StatusForbidden},
		{"platform admin", &util.Payload{UserID: uuid.New(), Role: domain.RoleSuper}, http.StatusForbidden},
		{"member without agency", &util.Payload{UserID: uuid.New(), Role: domain.RoleMember}, http.StatusForbidden},
	} {
		t.Run(tc.name, func(t *testing.T) {
			profileUC := memberProfileUseCaseFunc(func(context.Context, domain.Actor) (*domain.MemberProfile, error) {
				t.Fatal("unauthorized request reached profile lookup")
				return nil, nil
			})
			h := handler.NewMemberHandler(nil, nil, nil, nil, nil, profileUC)
			req := httptest.NewRequest(http.MethodGet, "/members/me", nil)
			if tc.payload != nil {
				req = req.WithContext(util.WithPayload(req.Context(), tc.payload))
			}
			rec := httptest.NewRecorder()
			h.GetMyProfile(rec, req)
			if rec.Code != tc.status {
				t.Fatalf("expected %d, got %d: %s", tc.status, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestMyProfileRejectsDeletedMemberSession(t *testing.T) {
	agencyID := uuid.New()
	profileUC := memberProfileUseCaseFunc(func(context.Context, domain.Actor) (*domain.MemberProfile, error) {
		return nil, fmt.Errorf("profile lookup: %w", domain.ErrMemberNotFound)
	})
	h := handler.NewMemberHandler(nil, nil, nil, nil, nil, profileUC)
	req := httptest.NewRequest(http.MethodGet, "/members/me", nil)
	req = req.WithContext(util.WithPayload(req.Context(), &util.Payload{UserID: uuid.New(), Role: domain.RoleMember, AgencyID: &agencyID}))
	rec := httptest.NewRecorder()
	h.GetMyProfile(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for deleted member, got %d: %s", rec.Code, rec.Body.String())
	}
}
