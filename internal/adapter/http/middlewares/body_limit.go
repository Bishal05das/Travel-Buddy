package middleware

import "net/http"

// LimitBody caps request bodies at maxBytes; reading past it fails and the
// multipart and JSON decoders return an error instead of filling memory
// or disk.
func LimitBody(maxBytes int64) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
