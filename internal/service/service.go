package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/andrey/epoch-server/internal/clients/contract"
	"github.com/andrey/epoch-server/internal/clients/epoch"
	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/andrey/epoch-server/internal/clients/storage"
	"github.com/andrey/epoch-server/internal/clients/subsidizer"
	"github.com/andrey/epoch-server/internal/config"
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

// NewHTTPHandler creates and configures the HTTP handler with routes and middlewares
func (s *Service) NewHTTPHandler() http.Handler {
	// Create new ServeMux with Go 1.22+ routing patterns
	mux := http.NewServeMux()
	
	// Register routes using modern Go routing patterns
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /epochs/start", s.handleStartEpoch)
	mux.HandleFunc("POST /epochs/distribute", s.handleDistributeSubsidies)
	
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

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
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
