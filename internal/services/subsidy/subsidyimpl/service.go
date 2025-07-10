package subsidyimpl

import (
	"context"
	"fmt"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/go-pkgz/lgr"
)

// Service implements the subsidy service interface
type Service struct {
	lazyDistributor subsidy.LazyDistributor
	logger          lgr.L
	config          *config.Config
}

// New creates a new subsidy service implementation
func New(lazyDistributor subsidy.LazyDistributor, logger lgr.L, cfg *config.Config) *Service {
	return &Service{
		lazyDistributor: lazyDistributor,
		logger:          logger,
		config:          cfg,
	}
}

// DistributeSubsidies manages the distribution of subsidies for a vault
func (s *Service) DistributeSubsidies(ctx context.Context, vaultId string) (*subsidy.SubsidyDistributionResponse, error) {
	// Validate input
	if vaultId == "" {
		return nil, fmt.Errorf("%w: vaultId cannot be empty", subsidy.ErrInvalidInput)
	}

	s.logger.Logf("INFO starting subsidy distribution for vault %s", vaultId)

	if err := s.lazyDistributor.Run(ctx, vaultId); err != nil {
		s.logger.Logf("ERROR subsidy distribution failed for vault %s: %v", vaultId, err)
		// Check if the error is transaction-related
		if isTransactionError(err) {
			return nil, fmt.Errorf("%w: failed to run lazy distributor for vault %s: %v", subsidy.ErrTransactionFailed, vaultId, err)
		}
		return nil, fmt.Errorf("failed to run lazy distributor for vault %s: %w", vaultId, err)
	}

	s.logger.Logf("INFO successfully completed subsidy distribution for vault %s", vaultId)

	return &subsidy.SubsidyDistributionResponse{
		VaultID:           vaultId,
		EpochID:           "current",
		TotalSubsidies:    "0",
		AccountsProcessed: 0,
		MerkleRoot:        "",
		Status:            "completed",
	}, nil
}

// isTransactionError determines if an error is related to blockchain transaction failures
func isTransactionError(err error) bool {
	errStr := err.Error()
	// Check for common blockchain transaction error patterns
	transactionErrors := []string{
		"failed to call",
		"transaction failed",
		"gas",
		"revert",
		"nonce",
		"insufficient funds",
		"execution reverted",
		"failed to send transaction",
		"transaction timeout",
	}

	for _, txErr := range transactionErrors {
		if contains(errStr, txErr) {
			return true
		}
	}
	return false
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
