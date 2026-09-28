package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bishal05das/travelbuddy/config"
	"github.com/bishal05das/travelbuddy/internal/domain"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

// fakeAuthorizer grants exactly the listed "resource:action" pairs to any member.
type fakeAuthorizer map[string]bool

func (f fakeAuthorizer) MemberHasPermission(_ context.Context, _, _ uuid.UUID, resource, action string) (bool, error) {
	return f[resource+":"+action], nil
}

var ok200 = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

func request(p *util.Payload, pathValues map[string]string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	for k, v := range pathValues {
		req.SetPathValue(k, v)
	}
	if p != nil {
		req = req.WithContext(util.WithPayload(req.Context(), p))
	}
	return req
}

func status(h http.Handler, req *http.Request) int {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code
}

func TestRequirePermission(t *testing.T) {
	m := NewMiddlewareManager(&config.Config{}, fakeAuthorizer{"tour:create": true})
	h := m.RequirePermission("tour", "create")(ok200)

	agency := uuid.New()
	other := uuid.New()
	member := &util.Payload{UserID: uuid.New(), Role: domain.RoleMember, AgencyID: &agency}

	tests := map[string]struct {
		payload *util.Payload
		path    map[string]string
		want    int
	}{
		"member of the agency with permission": {member, map[string]string{"agency_id": agency.String()}, http.StatusOK},
		"member acting on another agency":      {member, map[string]string{"agency_id": other.String()}, http.StatusForbidden},
		"member without agency path value":     {member, nil, http.StatusOK},
		"super user":                           {&util.Payload{Role: domain.RoleSuper}, map[string]string{"agency_id": other.String()}, http.StatusOK},
		"plain user":                           {&util.Payload{UserID: uuid.New(), Role: domain.RoleUser}, nil, http.StatusForbidden},
		"member token without agency":          {&util.Payload{UserID: uuid.New(), Role: domain.RoleMember}, nil, http.StatusForbidden},
		"no claims":                            {nil, nil, http.StatusForbidden},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := status(h, request(tt.payload, tt.path)); got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}

	t.Run("member lacking the permission", func(t *testing.T) {
		h := m.RequirePermission("tour", "delete")(ok200)
		if got := status(h, request(member, nil)); got != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", got)
		}
	})
}

func TestRequireSelfOrSuper(t *testing.T) {
	m := NewMiddlewareManager(&config.Config{}, nil)
	h := m.RequireSelfOrSuper("user_id")(ok200)
	me := uuid.New()

	tests := map[string]struct {
		payload *util.Payload
		target  string
		want    int
	}{
		"own account":         {&util.Payload{UserID: me, Role: domain.RoleUser}, me.String(), http.StatusOK},
		"someone else":        {&util.Payload{UserID: me, Role: domain.RoleUser}, uuid.NewString(), http.StatusForbidden},
		"member with same id": {&util.Payload{UserID: me, Role: domain.RoleMember}, me.String(), http.StatusForbidden},
		"super user":          {&util.Payload{UserID: me, Role: domain.RoleSuper}, uuid.NewString(), http.StatusOK},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := status(h, request(tt.payload, map[string]string{"user_id": tt.target}))
			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	m := NewMiddlewareManager(&config.Config{}, nil)
	h := m.RequireRole(domain.RoleSuper)(ok200)

	if got := status(h, request(&util.Payload{Role: domain.RoleSuper}, nil)); got != http.StatusOK {
		t.Fatalf("super: expected 200, got %d", got)
	}
	if got := status(h, request(&util.Payload{Role: domain.RoleUser}, nil)); got != http.StatusForbidden {
		t.Fatalf("user: expected 403, got %d", got)
	}
}
