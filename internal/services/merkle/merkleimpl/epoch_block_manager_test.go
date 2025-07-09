package merkleimpl

import (
	"context"
	"testing"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/go-pkgz/lgr"
)

func TestEpochBlockManager_NewEpochBlockManager(t *testing.T) {
	mockClient := &subgraph.SubgraphClientMock{
		QueryEpochWithBlockInfoFunc: func(ctx context.Context, epochNumber string) (*subgraph.Epoch, error) {
			return &subgraph.Epoch{
				EpochNumber: epochNumber,
				StartTimestamp: "1640000000",
				EndTimestamp: "1640086400",
			}, nil
		},
	}
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
	mockClient := &subgraph.SubgraphClientMock{
		QueryEpochWithBlockInfoFunc: func(ctx context.Context, epochNumber string) (*subgraph.Epoch, error) {
			return &subgraph.Epoch{
				EpochNumber: epochNumber,
				StartTimestamp: "1640000000",
				EndTimestamp: "1640086400",
			}, nil
		},
	}
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