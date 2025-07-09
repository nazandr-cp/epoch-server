package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-pkgz/lgr"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	logger lgr.L
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(logger lgr.L) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

// HandleHealth returns the health status of the service
func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}