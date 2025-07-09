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

func (m *MockGraphClient) ExecuteQueryAtBlock(ctx context.Context, query string, variables map[string]interface{}, blockNumber int64, response interface{}) error {
	// Mock implementation for testing
	return nil
}

func (m *MockGraphClient) ExecutePaginatedQueryAtBlock(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, blockNumber int64, response interface{}) error {
	// Mock implementation for testing
	return nil
}

func (m *MockGraphClient) QueryEpochByNumber(ctx context.Context, epochNumber string) (*graph.Epoch, error) {
	// Mock implementation for testing
	return &graph.Epoch{}, nil
}

func (m *MockGraphClient) QueryCurrentActiveEpoch(ctx context.Context) (*graph.Epoch, error) {
	// Mock implementation for testing
	return &graph.Epoch{}, nil
}

func (m *MockGraphClient) QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*graph.Epoch, error) {
	// Mock implementation for testing
	return &graph.Epoch{
		EpochNumber:    epochNumber,
		CreatedAtBlock: "12345",
		UpdatedAtBlock: "12346",
		StartTimestamp: "1640995200",
		EndTimestamp:   "1641081600",
	}, nil
}

func (m *MockGraphClient) QueryAccountSubsidiesAtBlock(ctx context.Context, vaultAddress string, blockNumber int64) ([]graph.AccountSubsidy, error) {
	// Mock implementation for testing
	return []graph.AccountSubsidy{}, nil
}

func (m *MockGraphClient) QueryMerkleDistributionForEpoch(ctx context.Context, epochNumber string, vaultAddress string) (*graph.MerkleDistribution, error) {
	// Mock implementation for testing
	return &graph.MerkleDistribution{}, nil
}

func (m *MockGraphClient) QueryAccountSubsidiesForEpoch(ctx context.Context, vaultAddress string, epochEndTimestamp string) ([]graph.AccountSubsidy, error) {
	// Mock implementation for testing
	return []graph.AccountSubsidy{}, nil
}

func TestEpochBlockManager_NewEpochBlockManager(t *testing.T) {
	mockClient := &MockGraphClient{}
	logger := lgr.NoOp

	ebm := NewEpochBlockManager(mockClient, logger)
	
	if ebm == nil {
		t.Error("NewEpochBlockManager returned nil")
	}
	
	if ebm.graphClient == nil {
		t.Error("EpochBlockManager graphClient is nil")
	}
	
	if ebm.logger == nil {
		t.Error("EpochBlockManager logger is nil")
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
	
	if pg.epochBlockManager == nil {
		t.Error("ProofGenerator epochBlockManager is nil")
	}
}