package middleware

import (
	"log"
	"net/http"
	"slices"

	"github.com/bishal05das/travelbuddy/internal/domain"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

// These middlewares run after Authentication, which puts the verified
// claims on the request context.

// RequireRole allows only callers whose token carries one of roles.
func (m *MiddlewareManager) RequireRole(roles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p, ok := util.PayloadFromContext(r.Context())
			if !ok || !slices.Contains(roles, p.Role) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireSelfOrSuper allows a caller to act on their own account (the
// path value named param) or a super user to act on any account.
func (m *MiddlewareManager) RequireSelfOrSuper(param string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p, ok := util.PayloadFromContext(r.Context())
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			if p.Role == domain.RoleSuper {
				next.ServeHTTP(w, r)
				return
			}
			id, err := uuid.Parse(r.PathValue(param))
			if err != nil || p.Role != domain.RoleUser || id != p.UserID {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission allows super users, and agency members whose role
// grants resource:action within their own agency. When the route has an
// {agency_id} path value it must be the member's agency; routes without one
// rely on the use case to scope the change to the member's agency.
func (m *MiddlewareManager) RequirePermission(resource, action string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p, ok := util.PayloadFromContext(r.Context())
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			if p.Role == domain.RoleSuper {
				next.ServeHTTP(w, r)
				return
			}
			if p.Role != domain.RoleMember || p.AgencyID == nil {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			if raw := r.PathValue("agency_id"); raw != "" {
				pathAgency, err := uuid.Parse(raw)
				if err != nil {
					http.Error(w, "invalid agency id", http.StatusBadRequest)
					return
				}
				if pathAgency != *p.AgencyID {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
			}

			allowed, err := m.authorizer.MemberHasPermission(r.Context(), p.UserID, *p.AgencyID, resource, action)
			if err != nil {
				log.Println("permission check failed:", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if !allowed {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
