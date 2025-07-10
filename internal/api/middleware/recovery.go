package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/go-pkgz/lgr"
)

// Recovery creates a middleware for panic recovery
func Recovery(logger lgr.L) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic with stack trace
					logger.Logf("ERROR panic recovered: %v\nStack trace:\n%s", err, debug.Stack())

					// Return a 500 error response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					if err := json.NewEncoder(w).Encode(map[string]interface{}{
						"error": "Internal server error",
						"code":  http.StatusInternalServerError,
					}); err != nil {
						logger.Logf("ERROR failed to encode recovery error response: %v", err)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
