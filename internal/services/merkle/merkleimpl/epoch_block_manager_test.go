package merkleimpl

import (
	"context"
	"math/big"
	"testing"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/dgraph-io/badger/v4"
	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_New(t *testing.T) {
	tempDir := t.TempDir()
	logger := lgr.NoOp

	// Create badger database
	opts := badger.DefaultOptions(tempDir)
	opts.Logger = nil
	db, err := badger.Open(opts)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, db.Close())
	}()

	mockClient := &testServiceSubgraphClient{}

	service := New(db, mockClient, logger)

	assert.NotNil(t, service)
	assert.NotNil(t, service.store)
	assert.NotNil(t, service.graphClient)
	assert.NotNil(t, service.logger)
}

func TestService_CalculateTotalEarned(t *testing.T) {
	service := createTestServiceForEpochTests(t)

	subsidy := subgraph.AccountSubsidy{
		SecondsAccumulated: "1000000000000000000000", // 1000 seconds worth
		LastEffectiveValue: "1000000000000000000",    // 1 token per second
		UpdatedAtTimestamp: "1640000000",             // Some timestamp
	}

	endTimestamp := int64(1640001000) // 1000 seconds later

	totalEarned, err := service.CalculateTotalEarned(subsidy, endTimestamp)
	require.NoError(t, err)

	// Should be approximately 2000 tokens (1000 accumulated + 1000 new)
	expected := big.NewInt(2000)
	assert.Equal(t, expected, totalEarned)
}

func TestService_SecondsToTokens(t *testing.T) {
	service := createTestServiceForEpochTests(t)

	seconds := new(big.Int)
	seconds.SetString("1000000000000000000000", 10) // 1000 * 1e18
	tokens := service.secondsToTokens(seconds)

	expected := big.NewInt(1000) // Should be 1000 tokens
	assert.Equal(t, expected, tokens)
}

func TestService_ParseEpochTimestamp(t *testing.T) {
	service := createTestServiceForEpochTests(t)

	epoch := &subgraph.Epoch{
		EpochNumber:                  "1",
		StartTimestamp:               "1640000000",
		EndTimestamp:                 "1640086400",
		ProcessingCompletedTimestamp: "1640086400",
		CreatedAtBlock:               "12345678",
		UpdatedAtBlock:               "12345680",
	}

	timestamp, err := service.parseEpochTimestamp(epoch)
	require.NoError(t, err)

	assert.Equal(t, "1", timestamp.EpochNumber)
	assert.Equal(t, int64(1640000000), timestamp.StartTimestamp)
	assert.Equal(t, int64(1640086400), timestamp.EndTimestamp)
	assert.Equal(t, int64(1640086400), timestamp.ProcessingCompletedTimestamp)
	assert.Equal(t, int64(12345678), timestamp.CreatedAtBlock)
	assert.Equal(t, int64(12345680), timestamp.UpdatedAtBlock)
}

// createTestServiceForEpochTests creates a service instance for epoch-related tests
func createTestServiceForEpochTests(t *testing.T) *Service {
	tempDir := t.TempDir()
	logger := lgr.NoOp

	// Create badger database
	opts := badger.DefaultOptions(tempDir)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create mock subgraph client
	mockClient := &testServiceSubgraphClient{}

	return New(db, mockClient, logger)
}

// testServiceSubgraphClient implements SubgraphClient for testing
type testServiceSubgraphClient struct{}

func (m *testServiceSubgraphClient) QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*subgraph.Epoch, error) {
	return &subgraph.Epoch{
		EpochNumber:                  epochNumber,
		StartTimestamp:               "1640000000",
		EndTimestamp:                 "1640086400",
		ProcessingCompletedTimestamp: "1640086400",
		CreatedAtBlock:               "12345678",
		UpdatedAtBlock:               "12345680",
	}, nil
}

func (m *testServiceSubgraphClient) QueryCurrentActiveEpoch(ctx context.Context) (*subgraph.Epoch, error) {
	return &subgraph.Epoch{
		EpochNumber:                  "1",
		StartTimestamp:               "1640000000",
		EndTimestamp:                 "1640086400",
		ProcessingCompletedTimestamp: "",
		CreatedAtBlock:               "12345678",
		UpdatedAtBlock:               "12345680",
	}, nil
}

func (m *testServiceSubgraphClient) ExecuteQuery(ctx context.Context, request subgraph.GraphQLRequest, response interface{}) error {
	return nil
}
