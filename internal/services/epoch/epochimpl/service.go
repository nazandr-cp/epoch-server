package epochimpl

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/andrey/epoch-server/internal/infra/utils"
	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/go-pkgz/lgr"
)

// ContractClient interface for blockchain operations
type ContractClient interface {
	StartEpoch(ctx context.Context) error
	GetCurrentEpochId(ctx context.Context) (*big.Int, error)
	ForceEndEpochWithZeroYield(ctx context.Context, epochId *big.Int, vaultAddress string) error
}

// SubgraphClient interface for querying subgraph data
type SubgraphClient interface {
	QueryAccounts(ctx context.Context) ([]subgraph.Account, error)
	ExecuteQuery(ctx context.Context, request subgraph.GraphQLRequest, response interface{}) error
}

// Calculator interface for earnings calculations
type Calculator interface {
	CalculateTotalEarned(subsidy subgraph.AccountSubsidy, epochEndTime int64) (*big.Int, error)
}

// Service implements the epoch service interface
type Service struct {
	contractClient ContractClient
	subgraphClient SubgraphClient
	calculator     Calculator
	logger         lgr.L
	config         *config.Config
}

// New creates a new epoch service implementation
func New(contractClient ContractClient, subgraphClient SubgraphClient, calculator Calculator, logger lgr.L, cfg *config.Config) *Service {
	return &Service{
		contractClient: contractClient,
		subgraphClient: subgraphClient,
		calculator:     calculator,
		logger:         logger,
		config:         cfg,
	}
}

// StartEpoch initiates a new epoch
func (s *Service) StartEpoch(ctx context.Context) error {
	// Check if there's an active epoch that needs to be completed first
	currentEpochId, err := s.contractClient.GetCurrentEpochId(ctx)
	if err != nil {
		s.logger.Logf("ERROR failed to get current epoch ID: %v", err)
		return fmt.Errorf("failed to get current epoch ID: %w", err)
	}

	// If there's a current epoch (ID > 0), we need to validate it's completed
	if currentEpochId.Cmp(big.NewInt(0)) > 0 {
		s.logger.Logf("INFO current epoch ID is %s, attempting to start new epoch", currentEpochId.String())
	}

	accounts, err := s.subgraphClient.QueryAccounts(ctx)
	if err != nil {
		s.logger.Logf("ERROR failed to query accounts: %v", err)
		return fmt.Errorf("failed to query accounts: %w", err)
	}

	s.logger.Logf("INFO found %d accounts for starting new epoch", len(accounts))

	if err := s.contractClient.StartEpoch(ctx); err != nil {
		s.logger.Logf("ERROR blockchain transaction failed for startEpoch: %v", err)
		// Check if the error is specifically about epoch still being active
		if isEpochStillActiveError(err) {
			return fmt.Errorf("%w: cannot start new epoch - current epoch %s is still active and must be completed first", epoch.ErrTransactionFailed, currentEpochId.String())
		}
		return fmt.Errorf("%w: failed to start epoch: %v", epoch.ErrTransactionFailed, err)
	}

	s.logger.Logf("INFO successfully initiated epoch start")
	return nil
}

// ForceEndEpoch forcibly ends an epoch with zero yield
func (s *Service) ForceEndEpoch(ctx context.Context, epochId uint64, vaultId string) error {
	// Validate input
	if vaultId == "" {
		return fmt.Errorf("%w: vaultId cannot be empty", epoch.ErrInvalidInput)
	}

	s.logger.Logf("INFO force ending epoch %d for vault %s", epochId, vaultId)

	// First, check current epoch to avoid trying to end an outdated epoch
	currentEpochId, err := s.contractClient.GetCurrentEpochId(ctx)
	if err != nil {
		s.logger.Logf("WARN failed to get current epoch ID, proceeding anyway: %v", err)
	} else {
		currentEpochInt := currentEpochId.Uint64()
		if epochId < currentEpochInt {
			s.logger.Logf("INFO epoch %d is already past (current: %d), considering it completed", epochId, currentEpochInt)
			return nil
		}
		if epochId > currentEpochInt {
			return fmt.Errorf("%w: cannot force end future epoch %d (current: %d)", epoch.ErrInvalidInput, epochId, currentEpochInt)
		}
	}

	// Convert epochId to big.Int
	epochIdBig := big.NewInt(int64(epochId))

	s.logger.Logf("INFO calling ForceEndEpochWithZeroYield for epoch %d", epochId)

	if err := s.contractClient.ForceEndEpochWithZeroYield(ctx, epochIdBig, vaultId); err != nil {
		s.logger.Logf("ERROR ForceEndEpochWithZeroYield failed for epoch %d: %v", epochId, err)
		if isTransactionError(err) {
			return fmt.Errorf("%w: failed to force end epoch %d for vault %s: %v", epoch.ErrTransactionFailed, epochId, vaultId, err)
		}
		return fmt.Errorf("failed to force end epoch %d for vault %s: %w", epochId, vaultId, err)
	}

	s.logger.Logf("INFO successfully force ended epoch %d for vault %s with zero yield", epochId, vaultId)
	return nil
}

// GetUserTotalEarned calculates total earned subsidies for a user
func (s *Service) GetUserTotalEarned(ctx context.Context, userAddress, vaultId string) (*epoch.UserEarningsResponse, error) {
	// Validate input
	if userAddress == "" {
		return nil, fmt.Errorf("%w: userAddress cannot be empty", epoch.ErrInvalidInput)
	}
	if vaultId == "" {
		return nil, fmt.Errorf("%w: vaultId cannot be empty", epoch.ErrInvalidInput)
	}

	// Normalize user address to lowercase
	userAddress = utils.NormalizeAddress(userAddress)

	s.logger.Logf("INFO getting total earned for user %s in vault %s", userAddress, vaultId)

	// Query user's account subsidy data and latest epoch end timestamp from subgraph
	query := fmt.Sprintf(`
		query {
			accountSubsidies(
				where: {
					account: "%s"
				}
			) {
				account {
					id
				}
				secondsAccumulated
				lastEffectiveValue
				updatedAtTimestamp
				collectionParticipation {
					vault {
						id
					}
				}
			}
			epoches(
				orderBy: epochNumber
				orderDirection: desc
				first: 1
			) {
				endTimestamp
			}
		}
	`, userAddress)

	variables := map[string]interface{}{}

	var response struct {
		AccountSubsidies []struct {
			Account struct {
				ID string `json:"id"`
			} `json:"account"`
			SecondsAccumulated      string `json:"secondsAccumulated"`
			LastEffectiveValue      string `json:"lastEffectiveValue"`
			UpdatedAtTimestamp      string `json:"updatedAtTimestamp"`
			CollectionParticipation struct {
				Vault struct {
					ID string `json:"id"`
				} `json:"vault"`
			} `json:"collectionParticipation"`
		} `json:"accountSubsidies"`
		Epoches []struct {
			EndTimestamp string `json:"endTimestamp"`
		} `json:"epoches"`
	}

	if err := s.subgraphClient.ExecuteQuery(ctx, subgraph.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}, &response); err != nil {
		s.logger.Logf("ERROR GraphQL query failed: %v", err)
		return nil, fmt.Errorf("failed to query user subsidy data: %w", err)
	}

	s.logger.Logf("DEBUG received %d account subsidies from subgraph", len(response.AccountSubsidies))

	// Filter account subsidies by vault
	var matchingSubsidy *struct {
		Account struct {
			ID string `json:"id"`
		} `json:"account"`
		SecondsAccumulated      string `json:"secondsAccumulated"`
		LastEffectiveValue      string `json:"lastEffectiveValue"`
		UpdatedAtTimestamp      string `json:"updatedAtTimestamp"`
		CollectionParticipation struct {
			Vault struct {
				ID string `json:"id"`
			} `json:"vault"`
		} `json:"collectionParticipation"`
	}

	for _, subsidy := range response.AccountSubsidies {
		if subsidy.CollectionParticipation.Vault.ID == vaultId {
			matchingSubsidy = &subsidy
			break
		}
	}

	if matchingSubsidy == nil {
		return nil, fmt.Errorf("%w: no subsidy data found for user %s in vault %s", epoch.ErrNotFound, userAddress, vaultId)
	}

	// Get epoch end timestamp from the latest epoch
	var epochEndTime int64
	if len(response.Epoches) > 0 {
		epochEndStr := response.Epoches[0].EndTimestamp
		var err error
		epochEndTime, err = strconv.ParseInt(epochEndStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid epoch end timestamp: %s", epochEndStr)
		}
		s.logger.Logf("INFO using epoch end timestamp: %d", epochEndTime)
	} else {
		// Fallback to current time if no epoch found
		epochEndTime = time.Now().Unix()
		s.logger.Logf("WARN no epoch found, using current time: %d", epochEndTime)
	}

	// Convert to subgraph.AccountSubsidy format for calculation
	subsidyForCalc := subgraph.AccountSubsidy{
		Account: subgraph.Account{
			ID: matchingSubsidy.Account.ID,
		},
		SecondsAccumulated: matchingSubsidy.SecondsAccumulated,
		LastEffectiveValue: matchingSubsidy.LastEffectiveValue,
		UpdatedAtTimestamp: matchingSubsidy.UpdatedAtTimestamp,
	}

	// Calculate total earned using the unified calculator
	totalEarned, err := s.calculator.CalculateTotalEarned(subsidyForCalc, epochEndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total earned: %w", err)
	}

	response_data := &epoch.UserEarningsResponse{
		UserAddress:   userAddress,
		VaultAddress:  vaultId,
		TotalEarned:   totalEarned.String(),
		CalculatedAt:  time.Now().Unix(),
		DataTimestamp: epochEndTime,
	}

	s.logger.Logf("INFO calculated total earned for user %s: %s (using epoch end: %d)", userAddress, totalEarned.String(), epochEndTime)
	return response_data, nil
}

// Helper functions

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
