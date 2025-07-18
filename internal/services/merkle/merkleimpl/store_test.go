package merkleimpl

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/dgraph-io/badger/v4"
	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerkleStore(t *testing.T) {
	tempDir := t.TempDir()
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Create badger database
	opts := badger.DefaultOptions(tempDir)
	opts.Logger = &testLogger{logger}
	db, err := badger.Open(opts)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, db.Close())
	}()

	// Create mock subgraph client
	mockClient := &mockSubgraphClient{}

	// Create unified service
	service := New(db, mockClient, logger)

	epochNumber := big.NewInt(16)
	vaultID := "0xf82b93f3d6a703b8b5949809771b1e725708590a"

	testSnapshot := merkle.MerkleSnapshot{
		Entries: []merkle.MerkleEntry{
			{
				Address:     "0x3575b992c5337226aecf4e7f93dfbe80c576ce15",
				TotalEarned: big.NewInt(1000),
			},
			{
				Address:     "0x8f37c5c4fa708e06a656d858003ef7dc5f60a29b",
				TotalEarned: big.NewInt(500),
			},
		},
		MerkleRoot:  "0x1234567890abcdef",
		Timestamp:   time.Now().Unix(),
		VaultID:     vaultID,
		BlockNumber: 23534102,
	}

	ctx := context.Background()

	t.Run("SaveSnapshot", func(t *testing.T) {
		err := service.store.SaveSnapshot(ctx, epochNumber, testSnapshot)
		require.NoError(t, err)

		// Verify by reading the saved snapshot
		saved, err := service.store.GetSnapshot(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.Equal(t, epochNumber, saved.EpochNumber)
		assert.False(t, saved.CreatedAt.IsZero())
	})

	t.Run("GetSnapshot", func(t *testing.T) {
		retrieved, err := service.store.GetSnapshot(ctx, epochNumber, vaultID)
		require.NoError(t, err)

		assert.Equal(t, epochNumber, retrieved.EpochNumber)
		assert.Equal(t, vaultID, retrieved.VaultID)
		assert.Equal(t, testSnapshot.MerkleRoot, retrieved.MerkleRoot)
		assert.Equal(t, testSnapshot.Timestamp, retrieved.Timestamp)
		assert.Equal(t, testSnapshot.BlockNumber, retrieved.BlockNumber)
		assert.Len(t, retrieved.Entries, 2)
	})

	t.Run("GetLatestSnapshot", func(t *testing.T) {
		latest, err := service.store.GetLatestSnapshot(ctx, vaultID)
		require.NoError(t, err)

		assert.Equal(t, epochNumber, latest.EpochNumber)
		assert.Equal(t, vaultID, latest.VaultID)
		assert.Equal(t, testSnapshot.MerkleRoot, latest.MerkleRoot)
	})

	t.Run("ListSnapshots", func(t *testing.T) {
		snapshots, err := service.store.ListSnapshots(ctx, vaultID, 10)
		require.NoError(t, err)

		assert.Len(t, snapshots, 1)
		assert.Equal(t, epochNumber, snapshots[0].EpochNumber)
		assert.Equal(t, vaultID, snapshots[0].VaultID)
	})

	t.Run("SaveMultipleSnapshots", func(t *testing.T) {
		// Save another snapshot
		epoch17 := big.NewInt(17)
		testSnapshot2 := testSnapshot
		testSnapshot2.MerkleRoot = "0xabcdef1234567890"
		testSnapshot2.Timestamp = time.Now().Unix() + 1000

		err := service.store.SaveSnapshot(ctx, epoch17, testSnapshot2)
		require.NoError(t, err)

		// List snapshots should return both, sorted by epoch number (descending)
		snapshots, err := service.store.ListSnapshots(ctx, vaultID, 10)
		require.NoError(t, err)

		assert.Len(t, snapshots, 2)
		assert.Equal(t, epoch17, snapshots[0].EpochNumber) // Latest first
		assert.Equal(t, epochNumber, snapshots[1].EpochNumber)

		// Latest should now be epoch 17
		latest, err := service.store.GetLatestSnapshot(ctx, vaultID)
		require.NoError(t, err)
		assert.Equal(t, epoch17, latest.EpochNumber)
	})

	t.Run("GetNonExistentSnapshot", func(t *testing.T) {
		_, err := service.store.GetSnapshot(ctx, big.NewInt(999), vaultID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "snapshot not found")
	})

	t.Run("GetLatestFromEmptyVault", func(t *testing.T) {
		_, err := service.store.GetLatestSnapshot(ctx, "0xnonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no snapshots found")
	})
}

// testLogger implements badger.Logger for testing
type testLogger struct {
	lgr lgr.L
}

func (l *testLogger) Errorf(format string, args ...interface{}) {
	l.lgr.Logf("ERROR "+format, args...)
}

func (l *testLogger) Warningf(format string, args ...interface{}) {
	l.lgr.Logf("WARN "+format, args...)
}

func (l *testLogger) Infof(format string, args ...interface{}) {
	l.lgr.Logf("INFO "+format, args...)
}

func (l *testLogger) Debugf(format string, args ...interface{}) {
	l.lgr.Logf("DEBUG "+format, args...)
}

// mockSubgraphClient implements merkle.SubgraphClient for testing
type mockSubgraphClient struct{}

func (m *mockSubgraphClient) QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*subgraph.Epoch, error) {
	return &subgraph.Epoch{
		EpochNumber:                  epochNumber,
		StartTimestamp:               "1640000000",
		EndTimestamp:                 "1640086400",
		ProcessingCompletedTimestamp: "1640086400",
		CreatedAtBlock:               "12345678",
		UpdatedAtBlock:               "12345680",
	}, nil
}

func (m *mockSubgraphClient) QueryCurrentActiveEpoch(ctx context.Context) (*subgraph.Epoch, error) {
	return &subgraph.Epoch{
		EpochNumber:                  "1",
		StartTimestamp:               "1640000000",
		EndTimestamp:                 "1640086400",
		ProcessingCompletedTimestamp: "",
		CreatedAtBlock:               "12345678",
		UpdatedAtBlock:               "12345680",
	}, nil
}

func (m *mockSubgraphClient) QueryAccountSubsidiesForVault(ctx context.Context, vaultAddress string) ([]subgraph.AccountSubsidy, error) {
	return []subgraph.AccountSubsidy{}, nil
}

func (m *mockSubgraphClient) ExecuteQuery(ctx context.Context, request subgraph.GraphQLRequest, response interface{}) error {
	return nil
}
