package handlers

import (
	"net/http"
	"strconv"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/infra/utils"
	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/go-pkgz/lgr"
	"github.com/go-pkgz/rest"
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
// @Summary Start epoch
// @Description Initiates the start of a new epoch for yield distribution
// @Tags epochs
// @Accept json
// @Produce json
// @Success 202 {object} epoch.StartEpochResponse "Epoch start accepted"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/epochs/start [post]
func (h *EpochHandler) HandleStartEpoch(w http.ResponseWriter, r *http.Request) {
	h.logger.Logf("INFO received start epoch request")

	response, err := h.epochService.StartEpoch(r.Context())
	if err != nil {
		h.logger.Logf("ERROR failed to start epoch: %v", err)
		writeErrorResponse(w, r, h.logger, err, "Failed to start epoch")
		return
	}

	if err := rest.EncodeJSON(w, http.StatusAccepted, response); err != nil {
		h.logger.Logf("ERROR failed to encode JSON response: %v", err)
	}
}

// HandleForceEndEpoch handles force end epoch requests
// @Summary Force end epoch
// @Description Forcibly ends an epoch with zero yield distribution
// @Tags epochs
// @Accept json
// @Produce json
// @Param epochId query uint64 true "Epoch ID to force end"
// @Success 202 {object} epoch.ForceEndEpochResponse "Epoch force end accepted"
// @Failure 400 {object} ErrorResponse "Bad request - missing or invalid epochId"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/epochs/force-end [post]
func (h *EpochHandler) HandleForceEndEpoch(w http.ResponseWriter, r *http.Request) {
	// Parse epoch ID from query parameter
	epochIdStr := r.URL.Query().Get("epochId")
	if epochIdStr == "" {
		h.logger.Logf("ERROR missing epochId parameter")
		writeErrorResponse(w, r, h.logger, epoch.ErrInvalidInput, "epochId parameter is required")
		return
	}

	epochId, err := strconv.ParseUint(epochIdStr, 10, 64)
	if err != nil {
		h.logger.Logf("ERROR invalid epochId parameter: %v", err)
		writeErrorResponse(w, r, h.logger, epoch.ErrInvalidInput, "invalid epochId parameter")
		return
	}

	// Use the vault address from configuration
	vaultId := h.config.Contracts.CollectionsVault

	h.logger.Logf("INFO received force end epoch request for epoch %d, vault %s", epochId, vaultId)

	response, err := h.epochService.ForceEndEpoch(r.Context(), epochId, vaultId)
	if err != nil {
		h.logger.Logf("ERROR failed to force end epoch %d for vault %s: %v", epochId, vaultId, err)
		writeErrorResponse(w, r, h.logger, err, "Failed to force end epoch")
		return
	}

	if err := rest.EncodeJSON(w, http.StatusAccepted, response); err != nil {
		h.logger.Logf("ERROR failed to encode JSON response: %v", err)
	}
}

// HandleGetUserTotalEarned handles user total earned requests
// @Summary Get user total earned
// @Description Retrieves the total amount earned by a user across all epochs
// @Tags users
// @Accept json
// @Produce json
// @Param address path string true "User wallet address" example:"0x1234567890123456789012345678901234567890"
// @Success 200 {object} epoch.UserEarningsResponse "User earnings information"
// @Failure 400 {object} ErrorResponse "Bad request - invalid address"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/users/{address}/total-earned [get]
func (h *EpochHandler) HandleGetUserTotalEarned(w http.ResponseWriter, r *http.Request) {
	// Extract user address from URL path
	userAddress := r.PathValue("address")
	if userAddress == "" {
		writeErrorResponse(w, r, h.logger, epoch.ErrInvalidInput, "Missing user address")
		return
	}

	// Use the vault address from configuration (normalize to lowercase)
	vaultId := utils.NormalizeAddress(h.config.Contracts.CollectionsVault)

	h.logger.Logf("INFO received get total earned request for user %s in vault %s", userAddress, vaultId)

	response, err := h.epochService.GetUserTotalEarned(r.Context(), userAddress, vaultId)
	if err != nil {
		h.logger.Logf("ERROR failed to get total earned for user %s: %v", userAddress, err)
		writeErrorResponse(w, r, h.logger, err, "Failed to get user total earned")
		return
	}

	rest.RenderJSON(w, response)
}
