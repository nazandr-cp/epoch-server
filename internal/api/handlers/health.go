package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-pkgz/lgr"
	"github.com/go-pkgz/rest"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	logger       lgr.L
	healthChecks []func() error
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(logger lgr.L, healthChecks ...func() error) *HealthHandler {
	return &HealthHandler{
		logger:       logger,
		healthChecks: healthChecks,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string            `json:"status" example:"ok"`
	Checks map[string]string `json:"checks,omitempty"`
}

// HandleHealth returns the health status of the service
// @Summary Health check
// @Description Returns the current health status of the epoch server
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse "Service is healthy"
// @Failure 503 {object} HealthResponse "Service is unhealthy"
// @Router /health [get]
func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "ok",
		Checks: make(map[string]string),
	}

	// Run all health checks
	healthy := true
	for i, check := range h.healthChecks {
		checkName := fmt.Sprintf("check_%d", i)
		if err := check(); err != nil {
			healthy = false
			response.Checks[checkName] = fmt.Sprintf("FAIL: %s", err.Error())
		} else {
			response.Checks[checkName] = "OK"
		}
	}

	if !healthy {
		response.Status = "unhealthy"
		rest.EncodeJSON(w, http.StatusServiceUnavailable, response)
		return
	}

	rest.RenderJSON(w, response)
}
