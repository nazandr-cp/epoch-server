package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/andrey/epoch-server/internal/clients/contract"
	"github.com/andrey/epoch-server/internal/clients/epoch"
	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/andrey/epoch-server/internal/clients/storage"
	"github.com/andrey/epoch-server/internal/clients/subsidizer"
	"github.com/andrey/epoch-server/internal/config"
	"github.com/go-pkgz/lgr"
)

// Predefined error types for different failure scenarios
var (
	ErrTransactionFailed = errors.New("blockchain transaction failed")
	ErrInvalidInput      = errors.New("invalid input parameters")
	ErrNotFound          = errors.New("resource not found")
	ErrTimeout           = errors.New("operation timed out")
)

type GraphClient interface {
	QueryAccounts(ctx context.Context) ([]graph.Account, error)
	ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error
	ExecutePaginatedQuery(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, response interface{}) error
}

type ContractClient interface {
	StartEpoch(ctx context.Context) error
	DistributeSubsidies(ctx context.Context, epochID string) error
	GetCurrentEpochId(ctx context.Context) (*big.Int, error)
	UpdateExchangeRate(ctx context.Context, lendingManagerAddress string) error
	AllocateYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string) error
	EndEpochWithSubsidies(ctx context.Context, epochId *big.Int, vaultAddress string, merkleRoot [32]byte, subsidiesDistributed *big.Int) error
}

type Service struct {
	graphClient    GraphClient
	contractClient ContractClient
	logger         lgr.L
	config         *config.Config
}

func NewService(graphClient *graph.Client, contractClient *contract.Client, logger lgr.L, cfg *config.Config) *Service {
	return &Service{
		graphClient:    graphClient,
		contractClient: contractClient,
		logger:         logger,
		config:         cfg,
	}
}

func (s *Service) StartEpoch(ctx context.Context) error {
	// Check if there's an active epoch that needs to be completed first
	currentEpochId, err := s.contractClient.GetCurrentEpochId(ctx)
	if err != nil {
		s.logger.Logf("ERROR failed to get current epoch ID: %v", err)
		return fmt.Errorf("failed to get current epoch ID: %w", err)
	}

	// If there's a current epoch (ID > 0), we need to validate it's completed
	// Note: The smart contract will reject the startEpoch call if epoch is still active
	// But we provide a more user-friendly error message here
	if currentEpochId.Cmp(big.NewInt(0)) > 0 {
		s.logger.Logf("INFO current epoch ID is %s, attempting to start new epoch", currentEpochId.String())
	}

	accounts, err := s.graphClient.QueryAccounts(ctx)
	if err != nil {
		s.logger.Logf("ERROR failed to query accounts: %v", err)
		return fmt.Errorf("failed to query accounts: %w", err)
	}

	s.logger.Logf("INFO found %d accounts for starting new epoch", len(accounts))

	if err := s.contractClient.StartEpoch(ctx); err != nil {
		s.logger.Logf("ERROR blockchain transaction failed for startEpoch: %v", err)
		// Check if the error is specifically about epoch still being active
		if isEpochStillActiveError(err) {
			return fmt.Errorf("%w: cannot start new epoch - current epoch %s is still active and must be completed first", ErrTransactionFailed, currentEpochId.String())
		}
		return fmt.Errorf("%w: failed to start epoch: %v", ErrTransactionFailed, err)
	}

	s.logger.Logf("INFO successfully initiated epoch start")
	return nil
}

func (s *Service) DistributeSubsidies(ctx context.Context, vaultId string) error {
	// Validate input
	if vaultId == "" {
		return fmt.Errorf("%w: vaultId cannot be empty", ErrInvalidInput)
	}

	s.logger.Logf("INFO starting subsidy distribution for vault %s", vaultId)

	epochClient := epoch.NewClientWithContract(s.logger, s.contractClient)
	subsidizerClient := subsidizer.NewClient(s.logger)
	storageClient := storage.NewClient(s.logger)

	lazyDistributor := NewLazyDistributor(
		s.graphClient,
		epochClient,
		subsidizerClient,
		storageClient,
		s.logger,
		s.config,
	)

	if err := lazyDistributor.Run(ctx, vaultId); err != nil {
		s.logger.Logf("ERROR subsidy distribution failed for vault %s: %v", vaultId, err)
		// Check if the error is transaction-related
		if isTransactionError(err) {
			return fmt.Errorf("%w: failed to run lazy distributor for vault %s: %v", ErrTransactionFailed, vaultId, err)
		}
		return fmt.Errorf("failed to run lazy distributor for vault %s: %w", vaultId, err)
	}

	s.logger.Logf("INFO successfully completed subsidy distribution for vault %s", vaultId)
	return nil
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

// isEpochStillActiveError checks if the error is specifically about an epoch still being active
func isEpochStillActiveError(err error) bool {
	errStr := err.Error()
	epochStillActiveErrors := []string{
		"EpochManager__EpochStillActive",
		"epoch still active",
		"EpochStillActive",
	}
	
	for _, epochErr := range epochStillActiveErrors {
		if contains(errStr, epochErr) {
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
