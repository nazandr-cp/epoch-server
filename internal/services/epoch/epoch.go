package epoch

import (
	"context"
)

//go:generate moq -out epoch_mocks.go . Service

// Service defines the interface for epoch management operations
type Service interface {
	// StartEpoch initiates a new epoch
	StartEpoch(ctx context.Context) error
	
	// ForceEndEpoch forcibly ends an epoch with zero yield
	ForceEndEpoch(ctx context.Context, epochId uint64, vaultId string) error
	
	// GetUserTotalEarned calculates total earned subsidies for a user
	GetUserTotalEarned(ctx context.Context, userAddress, vaultId string) (*UserEarningsResponse, error)
}