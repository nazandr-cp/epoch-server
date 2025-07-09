package merkleimpl

import (
	"context"
	"testing"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/go-pkgz/lgr"
)

// MockSubgraphClient for testing
type MockSubgraphClient struct{}

func (m *MockSubgraphClient) ExecuteQuery(ctx context.Context, request subgraph.GraphQLRequest, response interface{}) error {
	// Mock implementation for testing
	return nil
}

func (m *MockSubgraphClient) QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*subgraph.Epoch, error) {
	// Mock implementation for testing
	return &subgraph.Epoch{
		EpochNumber: epochNumber,
		StartTimestamp: "1640000000",
		EndTimestamp: "1640086400",
	}, nil
}

func (m *MockSubgraphClient) QueryCurrentActiveEpoch(ctx context.Context) (*subgraph.Epoch, error) {
	// Mock implementation for testing
	return &subgraph.Epoch{
		EpochNumber: "1",
		StartTimestamp: "1640000000",
		EndTimestamp: "1640086400",
	}, nil
}

func (m *MockSubgraphClient) ExecutePaginatedQuery(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, response interface{}) error {
	// Mock implementation for testing
	return nil
}

func (m *MockSubgraphClient) ExecuteQueryAtBlock(ctx context.Context, query string, variables map[string]interface{}, blockNumber int64, response interface{}) error {
	// Mock implementation for testing
	return nil
}

func (m *MockSubgraphClient) ExecutePaginatedQueryAtBlock(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, blockNumber int64, response interface{}) error {
	// Mock implementation for testing
	return nil
}

func (m *MockSubgraphClient) QueryEpochByNumber(ctx context.Context, epochNumber string) (*subgraph.Epoch, error) {
	// Mock implementation for testing
	return &subgraph.Epoch{}, nil
}

func (m *MockSubgraphClient) QueryAccountSubsidiesAtBlock(ctx context.Context, vaultAddress string, blockNumber int64) ([]subgraph.AccountSubsidy, error) {
	// Mock implementation for testing
	return []subgraph.AccountSubsidy{}, nil
}

func (m *MockSubgraphClient) QueryMerkleDistributionForEpoch(ctx context.Context, epochNumber string, vaultAddress string) (*subgraph.MerkleDistribution, error) {
	// Mock implementation for testing
	return &subgraph.MerkleDistribution{}, nil
}

func (m *MockSubgraphClient) QueryAccountSubsidiesForEpoch(ctx context.Context, vaultAddress string, epochEndTimestamp string) ([]subgraph.AccountSubsidy, error) {
	// Mock implementation for testing
	return []subgraph.AccountSubsidy{}, nil
}

func TestEpochBlockManager_NewEpochBlockManager(t *testing.T) {
	mockClient := &MockSubgraphClient{}
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
	mockClient := &MockSubgraphClient{}
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