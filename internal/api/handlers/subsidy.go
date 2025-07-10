package handlers

import (
	"net/http"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/go-pkgz/lgr"
	"github.com/go-pkgz/rest"
)

// SubsidyHandler handles subsidy-related HTTP requests
type SubsidyHandler struct {
	subsidyService subsidy.Service
	logger         lgr.L
	config         *config.Config
}

// NewSubsidyHandler creates a new subsidy handler
func NewSubsidyHandler(subsidyService subsidy.Service, logger lgr.L, cfg *config.Config) *SubsidyHandler {
	return &SubsidyHandler{
		subsidyService: subsidyService,
		logger:         logger,
		config:         cfg,
	}
}

// HandleDistributeSubsidies handles subsidy distribution requests
// @Summary Distribute subsidies
// @Description Initiates the distribution of subsidies for the current epoch
// @Tags epochs
// @Accept json
// @Produce json
// @Success 202 {object} subsidy.SubsidyDistributionResponse "Subsidy distribution accepted"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/epochs/distribute [post]
func (h *SubsidyHandler) HandleDistributeSubsidies(w http.ResponseWriter, r *http.Request) {
	// Use the vault address from configuration
	vaultId := h.config.Contracts.CollectionsVault

	h.logger.Logf("INFO received distribute subsidies request for vault %s", vaultId)

	response, err := h.subsidyService.DistributeSubsidies(r.Context(), vaultId)
	if err != nil {
		h.logger.Logf("ERROR failed to distribute subsidies for vault %s: %v", vaultId, err)
		writeErrorResponse(w, r, h.logger, err, "Failed to distribute subsidies")
		return
	}

	if err := rest.EncodeJSON(w, http.StatusAccepted, response); err != nil {
		h.logger.Logf("ERROR failed to encode JSON response: %v", err)
	}
}
