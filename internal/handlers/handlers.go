package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/andrey/epoch-server/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-pkgz/lgr"
)

type Service interface {
	StartEpoch(ctx context.Context, epochID string) error
	DistributeSubsidies(ctx context.Context, vaultId string) error
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) DistributeSubsidies(w http.ResponseWriter, r *http.Request) {
	epochID := chi.URLParam(r, "id")
	
	// For now, use the vault address from the subgraph data
	// TODO: Make this configurable or accept as parameter
	vaultId := "0x4a4be724f522946296a51d8c82c7c2e8e5a62655"

	h.logger.Logf("INFO received distribute subsidies request for epoch %s, vault %s", epochID, vaultId)

	if err := h.service.DistributeSubsidies(r.Context(), vaultId); err != nil {
		h.logger.Logf("ERROR failed to distribute subsidies for epoch %s, vault %s: %v", epochID, vaultId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
