package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bishal05das/travelbuddy/config"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

func TestAuthentication(t *testing.T) {
	m := NewMiddlewareManager(&config.Config{JWTSecretkey: "secret"}, nil)
	userID := uuid.New()
	token, _ := util.CreateJWT("secret", util.Payload{UserID: userID, Role: "user"}, time.Hour)

	var seen *util.Payload
	h := m.Authentication(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen, _ = util.PayloadFromContext(r.Context())
	}))

	tests := map[string]struct {
		header string
		want   int
	}{
		"valid bearer token": {"Bearer " + token, http.StatusOK},
		"missing header":     {"", http.StatusUnauthorized},
		"wrong scheme":       {"Basic " + token, http.StatusUnauthorized},
		"no scheme":          {token, http.StatusUnauthorized},
		"bad token":          {"Bearer a.b.c", http.StatusUnauthorized},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			seen = nil
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, rec.Code)
			}
			if tt.want == http.StatusOK && (seen == nil || seen.UserID != userID) {
				t.Fatalf("claims not stored on context: %+v", seen)
			}
		})
	}
}
