package service

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/andrey/epoch-server/internal/config"
	"github.com/go-pkgz/lgr"
)

type mockGraphClient struct {
	queryAccountsFunc func(ctx context.Context) ([]graph.Account, error)
}

func (m *mockGraphClient) QueryAccounts(ctx context.Context) ([]graph.Account, error) {
	if m.queryAccountsFunc != nil {
		return m.queryAccountsFunc(ctx)
	}
	return nil, nil
}


func (m *mockGraphClient) ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error {
	return nil
}

func (m *mockGraphClient) ExecutePaginatedQuery(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, response interface{}) error {
	return nil
}

type mockContractClient struct {
	startEpochFunc          func(ctx context.Context, epochID string) error
	distributeSubsidiesFunc func(ctx context.Context, epochID string) error
	getCurrentEpochIdFunc   func(ctx context.Context) (*big.Int, error)
	updateExchangeRateFunc  func(ctx context.Context, lendingManagerAddress string) error
	allocateYieldToEpochFunc func(ctx context.Context, epochId *big.Int, vaultAddress string) error
	endEpochWithSubsidiesFunc func(ctx context.Context, epochId *big.Int, vaultAddress string, merkleRoot [32]byte, subsidiesDistributed *big.Int) error
}

func (m *mockContractClient) StartEpoch(ctx context.Context, epochID string) error {
	if m.startEpochFunc != nil {
		return m.startEpochFunc(ctx, epochID)
	}
	return nil
}

func (m *mockContractClient) DistributeSubsidies(ctx context.Context, epochID string) error {
	if m.distributeSubsidiesFunc != nil {
		return m.distributeSubsidiesFunc(ctx, epochID)
	}
	return nil
}

func (m *mockContractClient) GetCurrentEpochId(ctx context.Context) (*big.Int, error) {
	if m.getCurrentEpochIdFunc != nil {
		return m.getCurrentEpochIdFunc(ctx)
	}
	// Default to returning epoch ID 1 for tests
	return big.NewInt(1), nil
}

func (m *mockContractClient) AllocateYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string) error {
	if m.allocateYieldToEpochFunc != nil {
		return m.allocateYieldToEpochFunc(ctx, epochId, vaultAddress)
	}
	return nil
}

func (m *mockContractClient) UpdateExchangeRate(ctx context.Context, lendingManagerAddress string) error {
	if m.updateExchangeRateFunc != nil {
		return m.updateExchangeRateFunc(ctx, lendingManagerAddress)
	}
	return nil
}

func (m *mockContractClient) EndEpochWithSubsidies(ctx context.Context, epochId *big.Int, vaultAddress string, merkleRoot [32]byte, subsidiesDistributed *big.Int) error {
	if m.endEpochWithSubsidiesFunc != nil {
		return m.endEpochWithSubsidiesFunc(ctx, epochId, vaultAddress, merkleRoot, subsidiesDistributed)
	}
	return nil
}

func TestService_StartEpoch(t *testing.T) {
	tests := []struct {
		name                       string
		epochID                    string
		mockGraphClient            *mockGraphClient
		mockContractClient         *mockContractClient
		wantErr                    bool
		expectedQueryAccountsCalls int
		expectedStartEpochCalls    int
	}{
		{
			name:    "successful start epoch",
			epochID: "epoch1",
			mockGraphClient: &mockGraphClient{
				queryAccountsFunc: func(ctx context.Context) ([]graph.Account, error) {
					return []graph.Account{
						{ID: "user1", TotalSecondsClaimed: "100"},
						{ID: "user2", TotalSecondsClaimed: "200"},
					}, nil
				},
			},
			mockContractClient: &mockContractClient{
				startEpochFunc: func(ctx context.Context, epochID string) error {
					return nil
				},
				getCurrentEpochIdFunc: func(ctx context.Context) (*big.Int, error) {
					return big.NewInt(0), nil // No current epoch
				},
			},
			wantErr:                    false,
			expectedQueryAccountsCalls: 1,
			expectedStartEpochCalls:    1,
		},
		{
			name:    "query accounts error",
			epochID: "epoch1",
			mockGraphClient: &mockGraphClient{
				queryAccountsFunc: func(ctx context.Context) ([]graph.Account, error) {
					return nil, errors.New("failed to query accounts")
				},
			},
			mockContractClient: &mockContractClient{
				startEpochFunc: func(ctx context.Context, epochID string) error {
					return nil
				},
				getCurrentEpochIdFunc: func(ctx context.Context) (*big.Int, error) {
					return big.NewInt(0), nil // No current epoch
				},
			},
			wantErr:                    true,
			expectedQueryAccountsCalls: 1,
			expectedStartEpochCalls:    0,
		},
		{
			name:    "epoch still active error",
			epochID: "epoch2",
			mockGraphClient: &mockGraphClient{
				queryAccountsFunc: func(ctx context.Context) ([]graph.Account, error) {
					return []graph.Account{
						{ID: "user1", TotalSecondsClaimed: "100"},
					}, nil
				},
			},
			mockContractClient: &mockContractClient{
				startEpochFunc: func(ctx context.Context, epochID string) error {
					return errors.New("execution reverted: EpochManager__EpochStillActive")
				},
				getCurrentEpochIdFunc: func(ctx context.Context) (*big.Int, error) {
					return big.NewInt(1), nil // Current epoch 1 is active
				},
			},
			wantErr:                    true,
			expectedQueryAccountsCalls: 1,
			expectedStartEpochCalls:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var queryAccountsCalls, startEpochCalls int

			if tt.mockGraphClient.queryAccountsFunc != nil {
				originalFunc := tt.mockGraphClient.queryAccountsFunc
				tt.mockGraphClient.queryAccountsFunc = func(ctx context.Context) ([]graph.Account, error) {
					queryAccountsCalls++
					return originalFunc(ctx)
				}
			}


			if tt.mockContractClient.startEpochFunc != nil {
				originalFunc := tt.mockContractClient.startEpochFunc
				tt.mockContractClient.startEpochFunc = func(ctx context.Context, epochID string) error {
					startEpochCalls++
					return originalFunc(ctx, epochID)
				}
			}

			cfg := &config.Config{}
			service := &Service{
				graphClient:    tt.mockGraphClient,
				contractClient: tt.mockContractClient,
				logger:         lgr.NoOp,
				config:         cfg,
			}

			err := service.StartEpoch(context.Background(), tt.epochID)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if queryAccountsCalls != tt.expectedQueryAccountsCalls {
				t.Errorf("Expected %d QueryAccounts calls, got %d", tt.expectedQueryAccountsCalls, queryAccountsCalls)
			}
			if startEpochCalls != tt.expectedStartEpochCalls {
				t.Errorf("Expected %d StartEpoch calls, got %d", tt.expectedStartEpochCalls, startEpochCalls)
			}
		})
	}
}

func TestService_DistributeSubsidies(t *testing.T) {
	tests := []struct {
		name            string
		vaultId         string
		mockGraphClient *mockGraphClient
		wantErr         bool
	}{
		{
			name:            "successful distribute subsidies",
			vaultId:         "vault1",
			mockGraphClient: &mockGraphClient{},
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			cfg.Contracts.LendingManager = "0x64Bd8C3294956E039EDf1a4058b6588de3731248"
			service := &Service{
				graphClient:    tt.mockGraphClient,
				contractClient: &mockContractClient{},
				logger:         lgr.NoOp,
				config:         cfg,
			}

			err := service.DistributeSubsidies(context.Background(), tt.vaultId)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
