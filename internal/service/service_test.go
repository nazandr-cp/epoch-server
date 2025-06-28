package service

import (
	"context"
	"errors"
	"testing"

	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/go-pkgz/lgr"
)

type mockGraphClient struct {
	queryUsersFunc       func(ctx context.Context) ([]graph.User, error)
	queryEligibilityFunc func(ctx context.Context, epochID string) ([]graph.Eligibility, error)
}

func (m *mockGraphClient) QueryUsers(ctx context.Context) ([]graph.User, error) {
	if m.queryUsersFunc != nil {
		return m.queryUsersFunc(ctx)
	}
	return nil, nil
}

func (m *mockGraphClient) QueryEligibility(ctx context.Context, epochID string) ([]graph.Eligibility, error) {
	if m.queryEligibilityFunc != nil {
		return m.queryEligibilityFunc(ctx, epochID)
	}
	return nil, nil
}

func (m *mockGraphClient) ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error {
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
		name                    string
		epochID                 string
		mockGraphClient         *mockGraphClient
		mockContractClient      *mockContractClient
		wantErr                 bool
		expectedQueryUsersCalls int
		expectedQueryEligCalls  int
		expectedStartEpochCalls int
	}{
		{
			name:    "successful start epoch",
			epochID: "epoch1",
			mockGraphClient: &mockGraphClient{
				queryUsersFunc: func(ctx context.Context) ([]graph.User, error) {
					return []graph.User{
						{ID: "user1", TotalSecondsClaimed: "100"},
						{ID: "user2", TotalSecondsClaimed: "200"},
					}, nil
				},
				queryEligibilityFunc: func(ctx context.Context, epochID string) ([]graph.Eligibility, error) {
					return []graph.Eligibility{
						{ID: "eligibility1", IsEligible: true},
						{ID: "eligibility2", IsEligible: false},
					}, nil
				},
			},
			mockContractClient: &mockContractClient{
				startEpochFunc: func(ctx context.Context, epochID string) error {
					return nil
				},
			},
			wantErr:                 false,
			expectedQueryUsersCalls: 1,
			expectedQueryEligCalls:  1,
			expectedStartEpochCalls: 1,
		},
		{
			name:    "query users error",
			epochID: "epoch1",
			mockGraphClient: &mockGraphClient{
				queryUsersFunc: func(ctx context.Context) ([]graph.User, error) {
					return nil, errors.New("failed to query users")
				},
			},
			mockContractClient: &mockContractClient{
				startEpochFunc: func(ctx context.Context, epochID string) error {
					return nil
				},
			},
			wantErr:                 true,
			expectedQueryUsersCalls: 1,
			expectedQueryEligCalls:  0,
			expectedStartEpochCalls: 0,
		},
		{
			name:    "query eligibility error",
			epochID: "epoch1",
			mockGraphClient: &mockGraphClient{
				queryUsersFunc: func(ctx context.Context) ([]graph.User, error) {
					return []graph.User{{ID: "user1"}}, nil
				},
				queryEligibilityFunc: func(ctx context.Context, epochID string) ([]graph.Eligibility, error) {
					return nil, errors.New("failed to query eligibility")
				},
			},
			mockContractClient: &mockContractClient{
				startEpochFunc: func(ctx context.Context, epochID string) error {
					return nil
				},
			},
			wantErr:                 true,
			expectedQueryUsersCalls: 1,
			expectedQueryEligCalls:  1,
			expectedStartEpochCalls: 0,
		},
		{
			name:    "contract start epoch error",
			epochID: "epoch1",
			mockGraphClient: &mockGraphClient{
				queryUsersFunc: func(ctx context.Context) ([]graph.User, error) {
					return []graph.User{{ID: "user1"}}, nil
				},
				queryEligibilityFunc: func(ctx context.Context, epochID string) ([]graph.Eligibility, error) {
					return []graph.Eligibility{{ID: "eligibility1"}}, nil
				},
			},
			mockContractClient: &mockContractClient{
				startEpochFunc: func(ctx context.Context, epochID string) error {
					return errors.New("failed to start epoch")
				},
			},
			wantErr:                 true,
			expectedQueryUsersCalls: 1,
			expectedQueryEligCalls:  1,
			expectedStartEpochCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var queryUsersCalls, queryEligCalls, startEpochCalls int

			if tt.mockGraphClient.queryUsersFunc != nil {
				originalFunc := tt.mockGraphClient.queryUsersFunc
				tt.mockGraphClient.queryUsersFunc = func(ctx context.Context) ([]graph.User, error) {
					queryUsersCalls++
					return originalFunc(ctx)
				}
			}

			if tt.mockGraphClient.queryEligibilityFunc != nil {
				originalFunc := tt.mockGraphClient.queryEligibilityFunc
				tt.mockGraphClient.queryEligibilityFunc = func(ctx context.Context, epochID string) ([]graph.Eligibility, error) {
					queryEligCalls++
					return originalFunc(ctx, epochID)
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
			if queryUsersCalls != tt.expectedQueryUsersCalls {
				t.Errorf("Expected %d QueryUsers calls, got %d", tt.expectedQueryUsersCalls, queryUsersCalls)
			}
			if queryEligCalls != tt.expectedQueryEligCalls {
				t.Errorf("Expected %d QueryEligibility calls, got %d", tt.expectedQueryEligCalls, queryEligCalls)
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
