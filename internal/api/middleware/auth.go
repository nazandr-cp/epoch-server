package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/go-pkgz/lgr"
)

// Auth creates a middleware for authentication (placeholder implementation)
func Auth(logger lgr.L) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For now, this is a placeholder - you can implement actual auth logic here
			// For example, checking API keys, JWT tokens, etc.
			
			// Example: Check for API key in header
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				// For development, we'll allow requests without auth
				// In production, you might want to reject these
				logger.Logf("WARN request without API key from %s", r.RemoteAddr)
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAuth creates a middleware that requires authentication
func RequireAuth(logger lgr.L) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for authentication
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				logger.Logf("ERROR unauthorized request from %s", r.RemoteAddr)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": "Unauthorized",
					"code":  http.StatusUnauthorized,
				})
				return
			}
			
			// In a real implementation, you would validate the API key
			// For now, just log and proceed
			logger.Logf("DEBUG authenticated request with key: %s...", apiKey[:min(len(apiKey), 8)])
			
			next.ServeHTTP(w, r)
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}