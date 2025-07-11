package subsidyimpl

import (
	"context"
	"fmt"
	"math/big"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/go-pkgz/lgr"
)

type Service struct {
	lazyDistributor subsidy.LazyDistributor
	epochService    epoch.Service
	logger          lgr.L
	config          *config.Config
}

func New(lazyDistributor subsidy.LazyDistributor, epochService epoch.Service, logger lgr.L, cfg *config.Config) *Service {
	return &Service{
		lazyDistributor: lazyDistributor,
		epochService:    epochService,
		logger:          logger,
		config:          cfg,
	}
}

func (s *Service) DistributeSubsidies(ctx context.Context, vaultId string) (*subsidy.SubsidyDistributionResponse, error) {
	if vaultId == "" {
		return nil, fmt.Errorf("%w: vaultId cannot be empty", subsidy.ErrInvalidInput)
	}

	s.logger.Logf("INFO starting subsidy distribution for vault %s", vaultId)

	currentEpochId, err := s.epochService.GetCurrentEpochId(ctx)
	if err != nil {
		s.logger.Logf("ERROR failed to get current epoch ID: %v", err)
		return nil, fmt.Errorf("failed to get current epoch ID: %w", err)
	}

	if currentEpochId == 0 {
		s.logger.Logf("ERROR no active epoch found for distribution")
		return nil, fmt.Errorf("%w: no active epoch found (epoch ID is 0)", subsidy.ErrInvalidEpochState)
	}

	s.logger.Logf("INFO distributing subsidies for epoch %d in vault %s", currentEpochId, vaultId)

	distributionResult, err := s.lazyDistributor.RunWithEpoch(ctx, vaultId, big.NewInt(int64(currentEpochId)))
	if err != nil {
		s.logger.Logf("ERROR subsidy distribution failed for vault %s: %v", vaultId, err)
		if isTransactionError(err) {
			return nil, fmt.Errorf("%w: failed to run lazy distributor for vault %s: %v", subsidy.ErrTransactionFailed, vaultId, err)
		}
		return nil, fmt.Errorf("failed to run lazy distributor for vault %s: %w", vaultId, err)
	}

	s.logger.Logf("INFO successfully completed subsidy distribution for vault %s", vaultId)

	epochResponse, err := s.epochService.CompleteEpochAfterDistribution(ctx, currentEpochId, vaultId)
	if err != nil {
		s.logger.Logf("ERROR failed to complete epoch %d after distribution for vault %s: %v", currentEpochId, vaultId, err)
		return nil, fmt.Errorf("failed to complete epoch %d after subsidy distribution for vault %s: %w", currentEpochId, vaultId, err)
	}

	s.logger.Logf("INFO successfully completed epoch %s after distribution for vault %s", epochResponse.EpochID, vaultId)

	return &subsidy.SubsidyDistributionResponse{
		VaultID:           vaultId,
		EpochID:           epochResponse.EpochID,
		TotalSubsidies:    distributionResult.TotalSubsidies.String(),
		AccountsProcessed: distributionResult.AccountsProcessed,
		MerkleRoot:        distributionResult.MerkleRoot,
		Status:            "completed",
	}, nil
}

func isTransactionError(err error) bool {
	errStr := err.Error()
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
