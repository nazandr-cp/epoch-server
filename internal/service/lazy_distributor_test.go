package service

import (
	"context"
	"math/big"
	"testing"

	"github.com/andrey/epoch-server/internal/clients/epoch"
	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/andrey/epoch-server/internal/clients/storage"
	"github.com/andrey/epoch-server/internal/clients/subsidizer"
	"github.com/go-pkgz/lgr"
)

func TestLazyDistributor_CalculateTotalEarned(t *testing.T) {
	logger := lgr.Default()
	epochClient := epoch.NewClient(logger)
	subsidizerClient := subsidizer.NewClient(logger)
	storageClient := storage.NewClient(logger)

	graphClient := &mockLazyGraphClient{}
	ld := NewLazyDistributor(graphClient, epochClient, subsidizerClient, storageClient, logger)

	subsidy := AccountSubsidyPerCollection{
		Account:            Account{ID: "0x123"},
		SecondsAccumulated: "1000000000000000000",
		SecondsClaimed:     "0",
		LastEffectiveValue: "1000000000000000000",
		UpdatedAtTimestamp: "1000",
	}

	epochEnd := int64(2000)
	totalEarned, err := ld.calculateTotalEarned(subsidy, epochEnd)
	if err != nil {
		t.Fatalf("calculateTotalEarned failed: %v", err)
	}

	expected := big.NewInt(1001)
	if totalEarned.Cmp(expected) != 0 {
		t.Errorf("expected %s, got %s", expected.String(), totalEarned.String())
	}
}

func TestLazyDistributor_BuildMerkleRoot(t *testing.T) {
	logger := lgr.Default()
	epochClient := epoch.NewClient(logger)
	subsidizerClient := subsidizer.NewClient(logger)
	storageClient := storage.NewClient(logger)

	graphClient := &mockLazyGraphClient{}
	ld := NewLazyDistributor(graphClient, epochClient, subsidizerClient, storageClient, logger)

	entries := []storage.MerkleEntry{
		{Address: "0x123", TotalEarned: big.NewInt(100)},
		{Address: "0x456", TotalEarned: big.NewInt(200)},
	}

	root := ld.buildMerkleRoot(entries)
	if root == [32]byte{} {
		t.Error("merkle root should not be empty")
	}
}

type mockLazyGraphClient struct{}

func (m *mockLazyGraphClient) QueryUsers(ctx context.Context) ([]graph.User, error) {
	return nil, nil
}

func (m *mockLazyGraphClient) QueryEligibility(ctx context.Context, epochID string) ([]graph.Eligibility, error) {
	return nil, nil
}

func (m *mockLazyGraphClient) ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error {
	mockResponse := struct {
		Data struct {
			AccountSubsidiesPerCollections []AccountSubsidyPerCollection `json:"accountSubsidiesPerCollections"`
		} `json:"data"`
	}{
		Data: struct {
			AccountSubsidiesPerCollections []AccountSubsidyPerCollection `json:"accountSubsidiesPerCollections"`
		}{
			AccountSubsidiesPerCollections: []AccountSubsidyPerCollection{
				{
					Account:            Account{ID: "0x123"},
					SecondsAccumulated: "1000000000000000000",
					SecondsClaimed:     "0",
					LastEffectiveValue: "1000000000000000000",
					UpdatedAtTimestamp: "1000",
				},
			},
		},
	}

	*response.(*struct {
		Data struct {
			AccountSubsidiesPerCollections []AccountSubsidyPerCollection `json:"accountSubsidiesPerCollections"`
		} `json:"data"`
	}) = mockResponse

	return nil
}
