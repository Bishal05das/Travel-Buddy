package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	defaultRequestsPerMinute = 30
	visitorIdleTTL           = 3 * time.Minute
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type ipStore struct {
	mu        sync.Mutex
	visitors  map[string]*visitor
	lastSweep time.Time
	perMinute int
}

func newIPStore(perMinute int) *ipStore {
	if perMinute <= 0 {
		perMinute = defaultRequestsPerMinute
	}
	return &ipStore{visitors: make(map[string]*visitor), lastSweep: time.Now(), perMinute: perMinute}
}

func (s *ipStore) get(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	// Evict idle visitors so the map cannot grow without bound.
	if now.Sub(s.lastSweep) > visitorIdleTTL {
		for key, v := range s.visitors {
			if now.Sub(v.lastSeen) > visitorIdleTTL {
				delete(s.visitors, key)
			}
		}
		s.lastSweep = now
	}

	v, ok := s.visitors[ip]
	if !ok {
		v = &visitor{limiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(s.perMinute)), s.perMinute)}
		s.visitors[ip] = v
	}
	v.lastSeen = now
	return v.limiter
}

func (m *MiddlewareManager) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := m.clientIP(r)
		if !m.limiter.get(ip).Allow() {
			http.Error(w, `{"error":"too many requests"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// clientIP identifies the caller by IP address (without the port, which
// changes per connection). Forwarding headers can be set by any client,
// so they are only honoured when the app runs behind a trusted proxy.
func (m *MiddlewareManager) clientIP(r *http.Request) string {
	if m.cfg.TrustProxyHeaders {
		if ip := strings.TrimSpace(r.Header.Get("X-Real-IP")); ip != "" {
			return ip
		}
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			return strings.TrimSpace(strings.Split(fwd, ",")[0])
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
