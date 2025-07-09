package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/go-pkgz/lgr"
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
func (h *SubsidyHandler) HandleDistributeSubsidies(w http.ResponseWriter, r *http.Request) {
	// Use the vault address from configuration
	vaultId := h.config.Contracts.CollectionsVault

	h.logger.Logf("INFO received distribute subsidies request for vault %s", vaultId)

	if err := h.subsidyService.DistributeSubsidies(r.Context(), vaultId); err != nil {
		h.logger.Logf("ERROR failed to distribute subsidies for vault %s: %v", vaultId, err)
		writeErrorResponse(w, err, "Failed to distribute subsidies")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"vaultID": vaultId,
		"message": "Subsidy distribution initiated successfully",
	})
}