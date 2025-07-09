package merkle

import (
	"context"
)

// Service defines the interface for merkle proof operations
type Service interface {
	// GenerateUserMerkleProof generates a merkle proof for a user's current earnings
	GenerateUserMerkleProof(ctx context.Context, userAddress, vaultAddress string) (*UserMerkleProofResponse, error)
	
	// GenerateHistoricalMerkleProof generates a merkle proof for a user's earnings at a specific epoch
	GenerateHistoricalMerkleProof(ctx context.Context, userAddress, vaultAddress, epochNumber string) (*UserMerkleProofResponse, error)
}