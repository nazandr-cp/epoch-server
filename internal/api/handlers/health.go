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

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

// HandleHealth returns the health status of the service
// @Summary Health check
// @Description Returns the current health status of the epoch server
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse "Service is healthy"
// @Router /health [get]
func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}