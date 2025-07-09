package subsidy

import (
	"context"
)

//go:generate moq -out subsidy_mocks.go . Service

// Service defines the interface for subsidy distribution operations
type Service interface {
	// DistributeSubsidies manages the distribution of subsidies for a vault
	DistributeSubsidies(ctx context.Context, vaultId string) error
}