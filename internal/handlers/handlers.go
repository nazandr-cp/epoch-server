package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/andrey/epoch-server/internal/config"
	"github.com/andrey/epoch-server/internal/service"
	"github.com/go-pkgz/lgr"
)

type Service interface {
	StartEpoch(ctx context.Context) error
	DistributeSubsidies(ctx context.Context, vaultId string) error
	ForceEndEpoch(ctx context.Context, epochId uint64, vaultId string) error
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
	config  *config.Config
}

func NewHandler(service *service.Service, logger lgr.L, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
		config:  cfg,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) StartEpoch(w http.ResponseWriter, r *http.Request) {
	h.logger.Logf("INFO received start epoch request")

	if err := h.service.StartEpoch(r.Context()); err != nil {
		h.logger.Logf("ERROR failed to start epoch: %v", err)
		h.writeErrorResponse(w, err, "Failed to start epoch")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"message": "Epoch start initiated successfully",
	})
}

func (h *Handler) DistributeSubsidies(w http.ResponseWriter, r *http.Request) {
	// Use the vault address from configuration
	vaultId := h.config.Contracts.CollectionsVault

	h.logger.Logf("INFO received distribute subsidies request for vault %s", vaultId)

	if err := h.service.DistributeSubsidies(r.Context(), vaultId); err != nil {
		h.logger.Logf("ERROR failed to distribute subsidies for vault %s: %v", vaultId, err)
		h.writeErrorResponse(w, err, "Failed to distribute subsidies")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"vaultID": vaultId,
		"message": "Subsidy distribution initiated successfully",
	})
}

func (h *Handler) ForceEndEpoch(w http.ResponseWriter, r *http.Request) {
	// Parse epoch ID from query parameter
	epochIdStr := r.URL.Query().Get("epochId")
	if epochIdStr == "" {
		h.logger.Logf("ERROR missing epochId parameter")
		h.writeErrorResponse(w, service.ErrInvalidInput, "epochId parameter is required")
		return
	}

	epochId, err := strconv.ParseUint(epochIdStr, 10, 64)
	if err != nil {
		h.logger.Logf("ERROR invalid epochId parameter: %v", err)
		h.writeErrorResponse(w, service.ErrInvalidInput, "invalid epochId parameter")
		return
	}

	// Use the vault address from configuration
	vaultId := h.config.Contracts.CollectionsVault

	h.logger.Logf("INFO received force end epoch request for epoch %d, vault %s", epochId, vaultId)

	if err := h.service.ForceEndEpoch(r.Context(), epochId, vaultId); err != nil {
		h.logger.Logf("ERROR failed to force end epoch %d for vault %s: %v", epochId, vaultId, err)
		h.writeErrorResponse(w, err, "Failed to force end epoch")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "accepted",
		"epochId": epochId,
		"vaultID": vaultId,
		"message": "Force end epoch initiated successfully",
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
