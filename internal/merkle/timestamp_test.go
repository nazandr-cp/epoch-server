package merkle

import (
	"context"
	"testing"

	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/go-pkgz/lgr"
)

// MockGraphClient for testing
type MockGraphClient struct{}

func (m *MockGraphClient) ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error {
	// Mock implementation for testing
	return nil
}

func (m *MockGraphClient) ExecutePaginatedQuery(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, response interface{}) error {
	// Mock implementation for testing
	return nil
}

func (m *MockGraphClient) QueryEpochByNumber(ctx context.Context, epochNumber string) (*graph.Epoch, error) {
	// Mock implementation for testing
	return &graph.Epoch{}, nil
}

func (m *MockGraphClient) QueryMerkleDistributionForEpoch(ctx context.Context, epochNumber string, vaultAddress string) (*graph.MerkleDistribution, error) {
	// Mock implementation for testing
	return &graph.MerkleDistribution{}, nil
}

func (m *MockGraphClient) QueryAccountSubsidiesForEpoch(ctx context.Context, vaultAddress string, epochEndTimestamp string) ([]graph.AccountSubsidy, error) {
	// Mock implementation for testing
	return []graph.AccountSubsidy{}, nil
}

func TestTimestampManager_NewTimestampManager(t *testing.T) {
	mockClient := &MockGraphClient{}
	logger := lgr.NoOp

	tm := NewTimestampManager(mockClient, logger)
	
	if tm == nil {
		t.Error("NewTimestampManager returned nil")
	}
	
	if tm.graphClient == nil {
		t.Error("TimestampManager graphClient is nil")
	}
	
	if tm.logger == nil {
		t.Error("TimestampManager logger is nil")
	}
}

func TestCalculator_NewCalculator(t *testing.T) {
	calc := NewCalculator()
	
	if calc == nil {
		t.Error("NewCalculator returned nil")
	}
}

func TestProofGenerator_NewProofGeneratorWithDependencies(t *testing.T) {
	mockClient := &MockGraphClient{}
	logger := lgr.NoOp

	pg := NewProofGeneratorWithDependencies(mockClient, logger)
	
	if pg == nil {
		t.Error("NewProofGeneratorWithDependencies returned nil")
	}
	
	if pg.calculator == nil {
		t.Error("ProofGenerator calculator is nil")
	}
	
	if pg.timestampManager == nil {
		t.Error("ProofGenerator timestampManager is nil")
	}
}