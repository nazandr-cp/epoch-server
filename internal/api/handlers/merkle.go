package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/infra/utils"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/go-pkgz/lgr"
)

// MerkleHandler handles merkle proof-related HTTP requests
type MerkleHandler struct {
	merkleService merkle.Service
	logger        lgr.L
	config        *config.Config
}

// NewMerkleHandler creates a new merkle handler
func NewMerkleHandler(merkleService merkle.Service, logger lgr.L, cfg *config.Config) *MerkleHandler {
	return &MerkleHandler{
		merkleService: merkleService,
		logger:        logger,
		config:        cfg,
	}
}

// HandleGetUserMerkleProof handles user merkle proof requests
// @Summary Get user merkle proof
// @Description Generates a merkle proof for a user's current earnings
// @Tags users
// @Accept json
// @Produce json
// @Param address path string true "User wallet address" example:"0x1234567890123456789012345678901234567890"
// @Param vault query string false "Vault address (optional, uses default if not provided)" example:"0x1234567890123456789012345678901234567890"
// @Success 200 {object} merkle.UserMerkleProofResponse "Merkle proof generated successfully"
// @Failure 400 {object} ErrorResponse "Bad request - invalid address"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/users/{address}/merkle-proof [get]
func (h *MerkleHandler) HandleGetUserMerkleProof(w http.ResponseWriter, r *http.Request) {
	// Extract user address from URL path
	userAddress := r.PathValue("address")
	if userAddress == "" {
		writeErrorResponse(w, merkle.ErrInvalidInput, "Missing user address")
		return
	}

	// Get vault address from query parameter or use default from config
	vaultAddress := r.URL.Query().Get("vault")
	if vaultAddress == "" {
		vaultAddress = h.config.Contracts.CollectionsVault
	}
	vaultAddress = utils.NormalizeAddress(vaultAddress)

	h.logger.Logf("INFO received merkle proof request for user %s in vault %s", userAddress, vaultAddress)

	response, err := h.merkleService.GenerateUserMerkleProof(r.Context(), userAddress, vaultAddress)
	if err != nil {
		h.logger.Logf("ERROR failed to generate merkle proof for user %s: %v", userAddress, err)
		writeErrorResponse(w, err, "Failed to generate merkle proof")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleGetUserHistoricalMerkleProof handles historical merkle proof requests
// @Summary Get historical merkle proof
// @Description Generates a merkle proof for a user's earnings at a specific epoch
// @Tags users
// @Accept json
// @Produce json
// @Param address path string true "User wallet address" example:"0x1234567890123456789012345678901234567890"
// @Param epochNumber path string true "Epoch number" example:"1"
// @Param vault query string false "Vault address (optional, uses default if not provided)" example:"0x1234567890123456789012345678901234567890"
// @Success 200 {object} merkle.UserMerkleProofResponse "Historical merkle proof generated successfully"
// @Failure 400 {object} ErrorResponse "Bad request - invalid address or epoch"
// @Failure 404 {object} ErrorResponse "User or epoch not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/users/{address}/merkle-proof/epoch/{epochNumber} [get]
func (h *MerkleHandler) HandleGetUserHistoricalMerkleProof(w http.ResponseWriter, r *http.Request) {
	// Extract user address and epoch number from URL path
	userAddress := r.PathValue("address")
	epochNumber := r.PathValue("epochNumber")
	
	if userAddress == "" {
		writeErrorResponse(w, merkle.ErrInvalidInput, "Missing user address")
		return
	}
	
	if epochNumber == "" {
		writeErrorResponse(w, merkle.ErrInvalidInput, "Missing epoch number")
		return
	}

	// Get vault address from query parameter or use default from config
	vaultAddress := r.URL.Query().Get("vault")
	if vaultAddress == "" {
		vaultAddress = h.config.Contracts.CollectionsVault
	}
	vaultAddress = utils.NormalizeAddress(vaultAddress)

	h.logger.Logf("INFO received historical merkle proof request for user %s in vault %s for epoch %s", userAddress, vaultAddress, epochNumber)

	response, err := h.merkleService.GenerateHistoricalMerkleProof(r.Context(), userAddress, vaultAddress, epochNumber)
	if err != nil {
		h.logger.Logf("ERROR failed to generate historical merkle proof for user %s epoch %s: %v", userAddress, epochNumber, err)
		writeErrorResponse(w, err, "Failed to generate historical merkle proof")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}