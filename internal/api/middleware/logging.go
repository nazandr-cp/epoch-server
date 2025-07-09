package middleware

import (
	"net/http"
	"time"

	"github.com/go-pkgz/lgr"
)

// Logging creates a middleware for request logging
func Logging(logger lgr.L) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture status code
			wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Process the request
			next.ServeHTTP(wrapper, r)

			// Log the request
			duration := time.Since(start)
			logger.Logf("INFO %s %s %d %v %s",
				r.Method,
				r.URL.Path,
				wrapper.statusCode,
				duration,
				r.RemoteAddr,
			)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
