package service

import (
	"context"
	"fmt"
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

	subsidy := AccountSubsidy{
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
			AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
		} `json:"data"`
	}{
		Data: struct {
			AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
		}{
			AccountSubsidies: []AccountSubsidy{
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
			AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
		} `json:"data"`
	}) = mockResponse

	return nil
}

func (m *mockLazyGraphClient) ExecutePaginatedQuery(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, response interface{}) error {
	mockResponse := struct {
		Data struct {
			AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
		} `json:"data"`
	}{
		Data: struct {
			AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
		}{
			AccountSubsidies: []AccountSubsidy{
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
			AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
		} `json:"data"`
	}) = mockResponse

	return nil
}

func TestLazyDistributor_Run_PaginationSuccess(t *testing.T) {
	logger := lgr.Default()
	epochClient := epoch.NewClient(logger)
	subsidizerClient := subsidizer.NewClient(logger)
	storageClient := storage.NewClient(logger)

	graphClient := &mockPaginatedGraphClient{
		totalRecords: 2500,
	}
	ld := NewLazyDistributor(graphClient, epochClient, subsidizerClient, storageClient, logger)

	if err := ld.Run(context.Background(), "vault1"); err != nil {
		t.Fatalf("expected no error with pagination, got: %v", err)
	}

	if graphClient.callCount < 1 {
		t.Errorf("expected at least 1 call to ExecutePaginatedQuery, got %d", graphClient.callCount)
	}
}

type mockPaginatedGraphClient struct {
	totalRecords int
	callCount    int
}

func (m *mockPaginatedGraphClient) QueryUsers(ctx context.Context) ([]graph.User, error) {
	return nil, nil
}

func (m *mockPaginatedGraphClient) QueryEligibility(ctx context.Context, epochID string) ([]graph.Eligibility, error) {
	return nil, nil
}

func (m *mockPaginatedGraphClient) ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error {
	return nil
}

func (m *mockPaginatedGraphClient) ExecutePaginatedQuery(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, response interface{}) error {
	m.callCount++

	// Simulate pagination behavior
	const pageSize = 1000

	var skip int
	if skipVal, ok := variables["skip"]; ok && skipVal != nil {
		skip = skipVal.(int)
	}

	remaining := m.totalRecords - skip

	var recordsInThisPage int
	if remaining > pageSize {
		recordsInThisPage = pageSize
	} else if remaining > 0 {
		recordsInThisPage = remaining
	} else {
		recordsInThisPage = 0
	}

	// Generate mock records for this page
	records := make([]AccountSubsidyPerCollection, recordsInThisPage)
	for i := 0; i < recordsInThisPage; i++ {
		records[i] = AccountSubsidyPerCollection{
			Account:            Account{ID: fmt.Sprintf("0x%d", skip+i)},
			SecondsAccumulated: "1000000000000000000",
			SecondsClaimed:     "0",
			LastEffectiveValue: "1000000000000000000",
			UpdatedAtTimestamp: "1000",
		}
	}

	mockResponse := struct {
		Data struct {
			AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
		} `json:"data"`
	}{
		Data: struct {
			AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
		}{
			AccountSubsidies: records,
		},
	}

	*response.(*struct {
		Data struct {
			AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
		} `json:"data"`
	}) = mockResponse

	return nil
}

func TestLazyDistributor_TwoPhaseCommit_FailedUpdateMerkleRoot(t *testing.T) {
	logger := lgr.Default()
	epochClient := &mockEpochClient{}
	subsidizerClient := &mockSubsidizerClient{shouldFailUpdate: true}
	storageClient := storage.NewClient(logger)

	graphClient := &mockLazyGraphClient{}
	ld := NewLazyDistributor(graphClient, epochClient, subsidizerClient, storageClient, logger)

	err := ld.Run(context.Background(), "vault1")
	if err == nil {
		t.Fatal("expected error when UpdateMerkleRootAndWaitForConfirmation fails")
	}

	if epochClient.finalizeEpochCalled {
		t.Error("FinalizeEpoch should not be called when UpdateMerkleRootAndWaitForConfirmation fails")
	}
}

func TestLazyDistributor_TwoPhaseCommit_SuccessfulFlow(t *testing.T) {
	logger := lgr.Default()
	epochClient := &mockEpochClient{}
	subsidizerClient := &mockSubsidizerClient{shouldFailUpdate: false}
	storageClient := storage.NewClient(logger)

	graphClient := &mockLazyGraphClient{}
	ld := NewLazyDistributor(graphClient, epochClient, subsidizerClient, storageClient, logger)

	err := ld.Run(context.Background(), "vault1")
	if err != nil {
		t.Fatalf("expected no error with successful flow, got: %v", err)
	}

	if !subsidizerClient.updateMerkleRootAndWaitCalled {
		t.Error("UpdateMerkleRootAndWaitForConfirmation should be called")
	}

	if !epochClient.finalizeEpochCalled {
		t.Error("FinalizeEpoch should be called after successful UpdateMerkleRootAndWaitForConfirmation")
	}
}

type mockEpochClient struct {
	finalizeEpochCalled bool
}

func (m *mockEpochClient) Current() epoch.EpochInfo {
	return epoch.EpochInfo{
		EndTime: 2000,
	}
}

func (m *mockEpochClient) FinalizeEpoch() error {
	m.finalizeEpochCalled = true
	return nil
}

type mockSubsidizerClient struct {
	shouldFailUpdate              bool
	updateMerkleRootCalled        bool
	updateMerkleRootAndWaitCalled bool
}

func (m *mockSubsidizerClient) UpdateMerkleRoot(ctx context.Context, vaultId string, root [32]byte) error {
	m.updateMerkleRootCalled = true
	if m.shouldFailUpdate {
		return fmt.Errorf("simulated UpdateMerkleRoot failure")
	}
	return nil
}

func (m *mockSubsidizerClient) UpdateMerkleRootAndWaitForConfirmation(ctx context.Context, vaultId string, root [32]byte) error {
	m.updateMerkleRootAndWaitCalled = true
	if m.shouldFailUpdate {
		return fmt.Errorf("simulated UpdateMerkleRootAndWaitForConfirmation failure")
	}
	return nil
}
