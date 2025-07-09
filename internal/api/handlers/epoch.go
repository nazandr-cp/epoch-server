package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/infra/utils"
	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/go-pkgz/lgr"
)

// EpochHandler handles epoch-related HTTP requests
type EpochHandler struct {
	epochService epoch.Service
	logger       lgr.L
	config       *config.Config
}

// NewEpochHandler creates a new epoch handler
func NewEpochHandler(epochService epoch.Service, logger lgr.L, cfg *config.Config) *EpochHandler {
	return &EpochHandler{
		epochService: epochService,
		logger:       logger,
		config:       cfg,
	}
}

// HandleStartEpoch handles epoch start requests
func (h *EpochHandler) HandleStartEpoch(w http.ResponseWriter, r *http.Request) {
	h.logger.Logf("INFO received start epoch request")

	if err := h.epochService.StartEpoch(r.Context()); err != nil {
		h.logger.Logf("ERROR failed to start epoch: %v", err)
		writeErrorResponse(w, err, "Failed to start epoch")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"message": "Epoch start initiated successfully",
	})
}

// HandleForceEndEpoch handles force end epoch requests
func (h *EpochHandler) HandleForceEndEpoch(w http.ResponseWriter, r *http.Request) {
	// Parse epoch ID from query parameter
	epochIdStr := r.URL.Query().Get("epochId")
	if epochIdStr == "" {
		h.logger.Logf("ERROR missing epochId parameter")
		writeErrorResponse(w, epoch.ErrInvalidInput, "epochId parameter is required")
		return
	}

	epochId, err := strconv.ParseUint(epochIdStr, 10, 64)
	if err != nil {
		h.logger.Logf("ERROR invalid epochId parameter: %v", err)
		writeErrorResponse(w, epoch.ErrInvalidInput, "invalid epochId parameter")
		return
	}

	// Use the vault address from configuration
	vaultId := h.config.Contracts.CollectionsVault

	h.logger.Logf("INFO received force end epoch request for epoch %d, vault %s", epochId, vaultId)

	if err := h.epochService.ForceEndEpoch(r.Context(), epochId, vaultId); err != nil {
		h.logger.Logf("ERROR failed to force end epoch %d for vault %s: %v", epochId, vaultId, err)
		writeErrorResponse(w, err, "Failed to force end epoch")
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

// HandleGetUserTotalEarned handles user total earned requests
func (h *EpochHandler) HandleGetUserTotalEarned(w http.ResponseWriter, r *http.Request) {
	// Extract user address from URL path
	userAddress := r.PathValue("address")
	if userAddress == "" {
		writeErrorResponse(w, epoch.ErrInvalidInput, "Missing user address")
		return
	}

	// Use the vault address from configuration (normalize to lowercase)
	vaultId := utils.NormalizeAddress(h.config.Contracts.CollectionsVault)

	h.logger.Logf("INFO received get total earned request for user %s in vault %s", userAddress, vaultId)

	response, err := h.epochService.GetUserTotalEarned(r.Context(), userAddress, vaultId)
	if err != nil {
		h.logger.Logf("ERROR failed to get total earned for user %s: %v", userAddress, err)
		writeErrorResponse(w, err, "Failed to get user total earned")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}