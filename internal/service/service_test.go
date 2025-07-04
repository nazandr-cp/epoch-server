package service

import (
	"context"
	"errors"
	"testing"

	"github.com/andrey/epoch-server/internal/clients/graph"
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
			},
			wantErr:                    true,
			expectedQueryAccountsCalls: 1,
			expectedStartEpochCalls:    0,
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

			service := &Service{
				graphClient:    tt.mockGraphClient,
				contractClient: tt.mockContractClient,
				logger:         lgr.NoOp,
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
			service := &Service{
				graphClient: tt.mockGraphClient,
				logger:      lgr.NoOp,
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
