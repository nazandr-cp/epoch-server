package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andrey/epoch-server/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-pkgz/lgr"
)

type Service interface {
	StartEpoch(ctx context.Context, epochID string) error
	DistributeSubsidies(ctx context.Context, vaultId string) error
}

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

type Handler struct {
	service Service
	logger  lgr.L
}

func NewHandler(service *service.Service, logger lgr.L) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) StartEpoch(w http.ResponseWriter, r *http.Request) {
	epochID := chi.URLParam(r, "id")

	h.logger.Logf("INFO received start epoch request for epoch %s", epochID)

	if err := h.service.StartEpoch(r.Context(), epochID); err != nil {
		h.logger.Logf("ERROR failed to start epoch %s: %v", epochID, err)
		h.writeErrorResponse(w, err, "Failed to start epoch")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"epochID": epochID,
		"message": "Epoch start initiated successfully",
	})
}

func (h *Handler) DistributeSubsidies(w http.ResponseWriter, r *http.Request) {
	epochID := chi.URLParam(r, "id")
	
	// Use the correct vault address from our deployment
	vaultId := "0xf82C7D08E65B74bf926552726305ff9ff0b0f700"

	h.logger.Logf("INFO received distribute subsidies request for epoch %s, vault %s", epochID, vaultId)

	if err := h.service.DistributeSubsidies(r.Context(), vaultId); err != nil {
		h.logger.Logf("ERROR failed to distribute subsidies for epoch %s, vault %s: %v", epochID, vaultId, err)
		h.writeErrorResponse(w, err, "Failed to distribute subsidies")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"epochID": epochID,
		"vaultID": vaultId,
		"message": "Subsidy distribution initiated successfully",
	})
}

// writeErrorResponse writes a structured error response based on the error type
func (h *Handler) writeErrorResponse(w http.ResponseWriter, err error, message string) {
	w.Header().Set("Content-Type", "application/json")
	
	var errResponse ErrorResponse
	errResponse.Error = message
	errResponse.Details = err.Error()

	// Determine appropriate HTTP status code based on error type
	if errors.Is(err, service.ErrTransactionFailed) {
		errResponse.Code = http.StatusBadGateway
		w.WriteHeader(http.StatusBadGateway)
	} else if errors.Is(err, service.ErrInvalidInput) {
		errResponse.Code = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
	} else if errors.Is(err, service.ErrNotFound) {
		errResponse.Code = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
	} else if errors.Is(err, service.ErrTimeout) {
		errResponse.Code = http.StatusRequestTimeout
		w.WriteHeader(http.StatusRequestTimeout)
	} else {
		// Default to internal server error
		errResponse.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(errResponse)
}
