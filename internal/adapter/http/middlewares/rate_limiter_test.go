package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bishal05das/travelbuddy/config"
)

func TestClientIP(t *testing.T) {
	direct := NewMiddlewareManager(&config.Config{}, nil)
	proxied := NewMiddlewareManager(&config.Config{TrustProxyHeaders: true}, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.7:51234"
	req.Header.Set("X-Forwarded-For", "198.51.100.1, 10.0.0.1")

	if got := direct.clientIP(req); got != "203.0.113.7" {
		t.Fatalf("direct: expected remote host without port, got %q", got)
	}
	if got := proxied.clientIP(req); got != "198.51.100.1" {
		t.Fatalf("proxied: expected first forwarded address, got %q", got)
	}
}

func TestRateLimiterSharesBudgetAcrossConnections(t *testing.T) {
	m := NewMiddlewareManager(&config.Config{}, nil)
	h := m.RateLimiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	limited := false
	for i := 0; i < requestsPerMinute+1; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// A new source port per request, like new connections, and a spoofed
		// header that must be ignored.
		req.RemoteAddr = "203.0.113.7:" + string(rune('0'+i%10)) + "000"
		req.Header.Set("X-Forwarded-For", "10.0.0."+string(rune('0'+i%10)))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code == http.StatusTooManyRequests {
			limited = true
		}
	}
	if !limited {
		t.Fatal("expected the same IP to be limited regardless of port or forwarded headers")
	}
}
