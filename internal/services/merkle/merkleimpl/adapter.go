package merkleimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/go-pkgz/lgr"
)

// MerkleProofService interface for generating merkle proofs  
type MerkleProofService interface {
	GenerateUserMerkleProof(ctx context.Context, userAddress, vaultAddress string) (*UserMerkleProofResponse, error)
	GenerateHistoricalMerkleProof(ctx context.Context, userAddress, vaultAddress, epochNumber string) (*UserMerkleProofResponse, error)
}

// UserMerkleProofResponse represents a merkle proof response (using existing structure)
type UserMerkleProofResponse struct {
	UserAddress   string   `json:"userAddress"`
	VaultAddress  string   `json:"vaultAddress"`
	EpochNumber   string   `json:"epochNumber,omitempty"`
	TotalEarned   string   `json:"totalEarned"`
	MerkleProof   []string `json:"merkleProof"`
	MerkleRoot    string   `json:"merkleRoot"`
	LeafIndex     int      `json:"leafIndex"`
	GeneratedAt   int64    `json:"generatedAt"`
}

// Service implements the merkle service interface
type Service struct {
	merkleProofService MerkleProofService
	logger             lgr.L
}

// New creates a new merkle service implementation
func New(merkleProofService MerkleProofService, logger lgr.L) *Service {
	return &Service{
		merkleProofService: merkleProofService,
		logger:             logger,
	}
}

// GenerateUserMerkleProof generates a merkle proof for a user's current earnings
func (s *Service) GenerateUserMerkleProof(ctx context.Context, userAddress, vaultAddress string) (*merkle.UserMerkleProofResponse, error) {
	if userAddress == "" {
		return nil, fmt.Errorf("%w: userAddress cannot be empty", merkle.ErrInvalidInput)
	}
	if vaultAddress == "" {
		return nil, fmt.Errorf("%w: vaultAddress cannot be empty", merkle.ErrInvalidInput)
	}

	s.logger.Logf("INFO generating merkle proof for user %s in vault %s", userAddress, vaultAddress)

	response, err := s.merkleProofService.GenerateUserMerkleProof(ctx, userAddress, vaultAddress)
	if err != nil {
		s.logger.Logf("ERROR failed to generate merkle proof: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	return &merkle.UserMerkleProofResponse{
		UserAddress:  response.UserAddress,
		VaultAddress: response.VaultAddress,
		EpochNumber:  response.EpochNumber,
		TotalEarned:  response.TotalEarned,
		MerkleProof:  response.MerkleProof,
		MerkleRoot:   response.MerkleRoot,
		LeafIndex:    response.LeafIndex,
		GeneratedAt:  time.Now().Unix(),
	}, nil
}

// GenerateHistoricalMerkleProof generates a merkle proof for a user's earnings at a specific epoch
func (s *Service) GenerateHistoricalMerkleProof(ctx context.Context, userAddress, vaultAddress, epochNumber string) (*merkle.UserMerkleProofResponse, error) {
	if userAddress == "" {
		return nil, fmt.Errorf("%w: userAddress cannot be empty", merkle.ErrInvalidInput)
	}
	if vaultAddress == "" {
		return nil, fmt.Errorf("%w: vaultAddress cannot be empty", merkle.ErrInvalidInput)
	}
	if epochNumber == "" {
		return nil, fmt.Errorf("%w: epochNumber cannot be empty", merkle.ErrInvalidInput)
	}

	s.logger.Logf("INFO generating historical merkle proof for user %s in vault %s for epoch %s", userAddress, vaultAddress, epochNumber)

	response, err := s.merkleProofService.GenerateHistoricalMerkleProof(ctx, userAddress, vaultAddress, epochNumber)
	if err != nil {
		s.logger.Logf("ERROR failed to generate historical merkle proof: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	return &merkle.UserMerkleProofResponse{
		UserAddress:  response.UserAddress,
		VaultAddress: response.VaultAddress,
		EpochNumber:  response.EpochNumber,
		TotalEarned:  response.TotalEarned,
		MerkleProof:  response.MerkleProof,
		MerkleRoot:   response.MerkleRoot,
		LeafIndex:    response.LeafIndex,
		GeneratedAt:  time.Now().Unix(),
	}, nil
}