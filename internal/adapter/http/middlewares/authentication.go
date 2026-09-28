package middleware

import (
	"net/http"
	"strings"

	util "github.com/bishal05das/travelbuddy/utils"
)

// Authentication verifies the Bearer token and stores its claims on the
// request context for handlers (util.GetPayload) and later middleware.
func (m *MiddlewareManager) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, err := util.ParseJWT(m.cfg.JWTSecretkey, strings.TrimSpace(token))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(util.WithPayload(r.Context(), payload)))
	})
}
