package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/andrey/epoch-server/internal/clients/contract"
	"github.com/andrey/epoch-server/internal/clients/epoch"
	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/andrey/epoch-server/internal/clients/storage"
	"github.com/andrey/epoch-server/internal/clients/subsidizer"
	"github.com/andrey/epoch-server/internal/config"
	"github.com/andrey/epoch-server/internal/merkle"
	"github.com/andrey/epoch-server/internal/utils"
	"github.com/go-pkgz/lgr"
	"github.com/go-pkgz/rest"
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
	QueryAccountSubsidiesForVault(ctx context.Context, vaultAddress string) ([]graph.AccountSubsidy, error)
	QueryCompletedEpochs(ctx context.Context) ([]graph.Epoch, error)
	QueryEpochByNumber(ctx context.Context, epochNumber string) (*graph.Epoch, error)
	QueryMerkleDistributionForEpoch(ctx context.Context, epochNumber string, vaultAddress string) (*graph.MerkleDistribution, error)
	QueryAccountSubsidiesForEpoch(ctx context.Context, vaultAddress string, epochEndTimestamp string) ([]graph.AccountSubsidy, error)
	ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error
	ExecutePaginatedQuery(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, response interface{}) error
}

type ContractClient interface {
	StartEpoch(ctx context.Context) error
	DistributeSubsidies(ctx context.Context, epochID string) error
	GetCurrentEpochId(ctx context.Context) (*big.Int, error)
	UpdateExchangeRate(ctx context.Context, lendingManagerAddress string) error
	AllocateYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string) error
	AllocateCumulativeYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string, amount *big.Int) error
	EndEpochWithSubsidies(ctx context.Context, epochId *big.Int, vaultAddress string, merkleRoot [32]byte, subsidiesDistributed *big.Int) error
}

type Service struct {
	graphClient       GraphClient
	contractClient    ContractClient
	merkleService     *merkle.MerkleProofService
	calculator        *merkle.Calculator
	logger            lgr.L
	config            *config.Config
}

func NewService(graphClient *graph.Client, contractClient *contract.Client, logger lgr.L, cfg *config.Config) *Service {
	merkleService := merkle.NewMerkleProofService(graphClient, logger)
	
	return &Service{
		graphClient:    graphClient,
		contractClient: contractClient,
		merkleService:  merkleService,
		calculator:     merkle.NewCalculator(),
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
	
	// Create real subsidizer client with blockchain configuration
	ethConfig := subsidizer.EthereumConfig{
		RPCURL:     s.config.Ethereum.RPCURL,
		PrivateKey: s.config.Ethereum.PrivateKey,
		GasLimit:   s.config.Ethereum.GasLimit,
		GasPrice:   s.config.Ethereum.GasPrice,
	}
	
	subsidizerClient, err := subsidizer.NewClientWithConfig(s.logger, ethConfig, s.config.Contracts.DebtSubsidizer)
	if err != nil {
		s.logger.Logf("WARN failed to create real subsidizer client, falling back to mock: %v", err)
		subsidizerClient = subsidizer.NewClient(s.logger)
	}
	
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

func (s *Service) GetUserTotalEarned(ctx context.Context, userAddress, vaultId string) (*UserEarningsResponse, error) {
	// Validate input
	if userAddress == "" {
		return nil, fmt.Errorf("%w: userAddress cannot be empty", ErrInvalidInput)
	}
	if vaultId == "" {
		return nil, fmt.Errorf("%w: vaultId cannot be empty", ErrInvalidInput)
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

	s.logger.Logf("DEBUG executing query: %s", query)
	
	// Use a generic response to capture any structure
	var genericResponse map[string]interface{}
	if err := s.graphClient.ExecuteQuery(ctx, graph.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}, &genericResponse); err != nil {
		s.logger.Logf("ERROR GraphQL query failed: %v", err)
		return nil, fmt.Errorf("failed to query user subsidy data: %w", err)
	}

	s.logger.Logf("DEBUG received raw response: %+v", genericResponse)

	// Now try to parse it properly
	if err := s.graphClient.ExecuteQuery(ctx, graph.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}, &response); err != nil {
		s.logger.Logf("ERROR GraphQL query failed on structured parse: %v", err)
		return nil, fmt.Errorf("failed to query user subsidy data: %w", err)
	}

	s.logger.Logf("DEBUG received %d account subsidies from subgraph", len(response.AccountSubsidies))
	for i, subsidy := range response.AccountSubsidies {
		s.logger.Logf("DEBUG subsidy %d: account=%s, vault=%s", i, subsidy.Account.ID, subsidy.CollectionParticipation.Vault.ID)
	}

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
		return nil, fmt.Errorf("%w: no subsidy data found for user %s in vault %s", ErrNotFound, userAddress, vaultId)
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

	// Convert to graph.AccountSubsidy format for calculation
	subsidyForCalc := graph.AccountSubsidy{
		Account: graph.Account{
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

	response_data := &UserEarningsResponse{
		UserAddress:  userAddress,
		VaultAddress: vaultId,
		TotalEarned:  totalEarned.String(),
		CalculatedAt: time.Now().Unix(),
	}

	s.logger.Logf("INFO calculated total earned for user %s: %s (using epoch end: %d)", userAddress, totalEarned.String(), epochEndTime)
	return response_data, nil
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

// NewHTTPHandler creates and configures the HTTP handler with routes and middlewares
func (s *Service) NewHTTPHandler() http.Handler {
	// Create new ServeMux with Go 1.22+ routing patterns
	mux := http.NewServeMux()
	
	// Register routes using modern Go routing patterns
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /epochs/start", s.handleStartEpoch)
	mux.HandleFunc("POST /epochs/distribute", s.handleDistributeSubsidies)
	mux.HandleFunc("GET /users/{address}/total-earned", s.handleGetUserTotalEarned)
	mux.HandleFunc("GET /users/{address}/merkle-proof", s.handleGetUserMerkleProof)
	mux.HandleFunc("GET /users/{address}/merkle-proof/epoch/{epochNumber}", s.handleGetUserHistoricalMerkleProof)
	
	// Apply middlewares using go-pkgz/rest
	// Chain middlewares manually: outermost -> innermost
	var handler http.Handler = mux
	
	// Apply middlewares in reverse order (last applied = outermost)
	handler = rest.Ping(handler)
	handler = rest.AppInfo("epoch-server", "andrey", "1.0.0")(handler)
	handler = rest.Recoverer(s.logger)(handler)
	handler = rest.RealIP(handler)
	
	return handler
}

// HTTP Handler methods

func (s *Service) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Service) handleStartEpoch(w http.ResponseWriter, r *http.Request) {
	s.logger.Logf("INFO received start epoch request")

	if err := s.StartEpoch(r.Context()); err != nil {
		s.logger.Logf("ERROR failed to start epoch: %v", err)
		s.writeErrorResponse(w, err, "Failed to start epoch")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"message": "Epoch start initiated successfully",
	})
}

func (s *Service) handleDistributeSubsidies(w http.ResponseWriter, r *http.Request) {
	// Use the vault address from configuration
	vaultId := s.config.Contracts.CollectionsVault

	s.logger.Logf("INFO received distribute subsidies request for vault %s", vaultId)

	if err := s.DistributeSubsidies(r.Context(), vaultId); err != nil {
		s.logger.Logf("ERROR failed to distribute subsidies for vault %s: %v", vaultId, err)
		s.writeErrorResponse(w, err, "Failed to distribute subsidies")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"vaultID": vaultId,
		"message": "Subsidy distribution initiated successfully",
	})
}

func (s *Service) handleGetUserTotalEarned(w http.ResponseWriter, r *http.Request) {
	// Extract user address from URL path
	userAddress := r.PathValue("address")
	if userAddress == "" {
		s.writeErrorResponse(w, fmt.Errorf("%w: missing user address in path", ErrInvalidInput), "Missing user address")
		return
	}

	// Use the vault address from configuration (normalize to lowercase)
	vaultId := utils.NormalizeAddress(s.config.Contracts.CollectionsVault)

	s.logger.Logf("INFO received get total earned request for user %s in vault %s", userAddress, vaultId)

	response, err := s.GetUserTotalEarned(r.Context(), userAddress, vaultId)
	if err != nil {
		s.logger.Logf("ERROR failed to get total earned for user %s: %v", userAddress, err)
		s.writeErrorResponse(w, err, "Failed to get user total earned")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Service) handleGetUserMerkleProof(w http.ResponseWriter, r *http.Request) {
	// Extract user address from URL path
	userAddress := r.PathValue("address")
	if userAddress == "" {
		s.writeErrorResponse(w, fmt.Errorf("%w: missing user address in path", ErrInvalidInput), "Missing user address")
		return
	}

	// Get vault address from query parameter or use default from config
	vaultAddress := r.URL.Query().Get("vault")
	if vaultAddress == "" {
		vaultAddress = s.config.Contracts.CollectionsVault
	}
	vaultAddress = utils.NormalizeAddress(vaultAddress)

	s.logger.Logf("INFO received merkle proof request for user %s in vault %s", userAddress, vaultAddress)

	response, err := s.merkleService.GenerateUserMerkleProof(r.Context(), userAddress, vaultAddress)
	if err != nil {
		s.logger.Logf("ERROR failed to generate merkle proof for user %s: %v", userAddress, err)
		s.writeErrorResponse(w, err, "Failed to generate merkle proof")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// UserEarningsResponse represents the response for user total earned query
type UserEarningsResponse struct {
	UserAddress   string `json:"userAddress"`
	VaultAddress  string `json:"vaultAddress"`
	TotalEarned   string `json:"totalEarned"`
	CalculatedAt  int64  `json:"calculatedAt"`
}

// writeErrorResponse writes a structured error response based on the error type
func (s *Service) writeErrorResponse(w http.ResponseWriter, err error, message string) {
	w.Header().Set("Content-Type", "application/json")
	
	var errResponse ErrorResponse
	errResponse.Error = message
	errResponse.Details = err.Error()

	// Determine appropriate HTTP status code based on error type
	if errors.Is(err, ErrTransactionFailed) {
		errResponse.Code = http.StatusBadGateway
		w.WriteHeader(http.StatusBadGateway)
	} else if errors.Is(err, ErrInvalidInput) {
		errResponse.Code = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
	} else if errors.Is(err, ErrNotFound) {
		errResponse.Code = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
	} else if errors.Is(err, ErrTimeout) {
		errResponse.Code = http.StatusRequestTimeout
		w.WriteHeader(http.StatusRequestTimeout)
	} else {
		// Default to internal server error
		errResponse.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(errResponse)
}

func (s *Service) handleGetUserHistoricalMerkleProof(w http.ResponseWriter, r *http.Request) {
	// Extract user address and epoch number from URL path
	userAddress := r.PathValue("address")
	epochNumber := r.PathValue("epochNumber")
	
	if userAddress == "" {
		s.writeErrorResponse(w, fmt.Errorf("%w: missing user address in path", ErrInvalidInput), "Missing user address")
		return
	}
	
	if epochNumber == "" {
		s.writeErrorResponse(w, fmt.Errorf("%w: missing epoch number in path", ErrInvalidInput), "Missing epoch number")
		return
	}

	// Get vault address from query parameter or use default from config
	vaultAddress := r.URL.Query().Get("vault")
	if vaultAddress == "" {
		vaultAddress = s.config.Contracts.CollectionsVault
	}
	vaultAddress = utils.NormalizeAddress(vaultAddress)

	s.logger.Logf("INFO received historical merkle proof request for user %s in vault %s for epoch %s", userAddress, vaultAddress, epochNumber)

	response, err := s.merkleService.GenerateHistoricalMerkleProof(r.Context(), userAddress, vaultAddress, epochNumber)
	if err != nil {
		s.logger.Logf("ERROR failed to generate historical merkle proof for user %s epoch %s: %v", userAddress, epochNumber, err)
		s.writeErrorResponse(w, err, "Failed to generate historical merkle proof")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
