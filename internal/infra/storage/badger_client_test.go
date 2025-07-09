package storage

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBadgerClient(t *testing.T) {
	tempDir := t.TempDir()
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	client, err := NewBadgerClient(logger, tempDir)
	require.NoError(t, err)
	defer client.Close()
	
	epochNumber := big.NewInt(16)
	vaultID := "0xf82b93f3d6a703b8b5949809771b1e725708590a"
	
	testSnapshot := MerkleSnapshot{
		Entries: []MerkleEntry{
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
	
	t.Run("SaveEpochSnapshot", func(t *testing.T) {
		err := client.SaveEpochSnapshot(ctx, epochNumber, testSnapshot)
		require.NoError(t, err)
		
		// Verify by reading the saved snapshot
		saved, err := client.GetEpochSnapshot(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.Equal(t, epochNumber, saved.EpochNumber)
		assert.False(t, saved.CreatedAt.IsZero())
	})
	
	t.Run("GetEpochSnapshot", func(t *testing.T) {
		retrieved, err := client.GetEpochSnapshot(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		
		assert.Equal(t, epochNumber, retrieved.EpochNumber)
		assert.Equal(t, vaultID, retrieved.VaultID)
		assert.Equal(t, testSnapshot.MerkleRoot, retrieved.MerkleRoot)
		assert.Equal(t, testSnapshot.Timestamp, retrieved.Timestamp)
		assert.Equal(t, testSnapshot.BlockNumber, retrieved.BlockNumber)
		assert.Len(t, retrieved.Entries, 2)
	})
	
	t.Run("GetLatestEpochSnapshot", func(t *testing.T) {
		latest, err := client.GetLatestEpochSnapshot(ctx, vaultID)
		require.NoError(t, err)
		
		assert.Equal(t, epochNumber, latest.EpochNumber)
		assert.Equal(t, vaultID, latest.VaultID)
		assert.Equal(t, testSnapshot.MerkleRoot, latest.MerkleRoot)
	})
	
	t.Run("ListEpochSnapshots", func(t *testing.T) {
		// First save the snapshot for this test
		err := client.SaveEpochSnapshot(ctx, epochNumber, testSnapshot)
		require.NoError(t, err)
		
		snapshots, err := client.ListEpochSnapshots(ctx, vaultID, 10)
		require.NoError(t, err)
		
		assert.Len(t, snapshots, 1)
		assert.Equal(t, epochNumber, snapshots[0].EpochNumber)
		assert.Equal(t, vaultID, snapshots[0].VaultID)
	})
	
	t.Run("SaveMultipleEpochs", func(t *testing.T) {
		// Save another epoch
		epoch17 := big.NewInt(17)
		testSnapshot2 := testSnapshot
		testSnapshot2.MerkleRoot = "0xabcdef1234567890"
		testSnapshot2.Timestamp = time.Now().Unix() + 1000
		
		err := client.SaveEpochSnapshot(ctx, epoch17, testSnapshot2)
		require.NoError(t, err)
		
		// List snapshots should return both, sorted by epoch number (descending)
		snapshots, err := client.ListEpochSnapshots(ctx, vaultID, 10)
		require.NoError(t, err)
		
		assert.Len(t, snapshots, 2)
		assert.Equal(t, epoch17, snapshots[0].EpochNumber) // Latest first
		assert.Equal(t, epochNumber, snapshots[1].EpochNumber)
		
		// Latest should now be epoch 17
		latest, err := client.GetLatestEpochSnapshot(ctx, vaultID)
		require.NoError(t, err)
		assert.Equal(t, epoch17, latest.EpochNumber)
	})
	
	t.Run("GetNonExistentSnapshot", func(t *testing.T) {
		_, err := client.GetEpochSnapshot(ctx, big.NewInt(999), vaultID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "snapshot not found")
	})
	
	t.Run("GetLatestFromEmptyVault", func(t *testing.T) {
		_, err := client.GetLatestEpochSnapshot(ctx, "0xnonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no snapshots found")
	})
	
	t.Run("LegacySaveSnapshot", func(t *testing.T) {
		// Test backward compatibility
		legacySnapshot := MerkleSnapshot{
			EpochNumber: big.NewInt(18),
			Entries: []MerkleEntry{
				{
					Address:     "0x1234567890abcdef",
					TotalEarned: big.NewInt(750),
				},
			},
			MerkleRoot:  "0xlegacyroot",
			Timestamp:   time.Now().Unix(),
			VaultID:     vaultID,
			BlockNumber: 23534200,
		}
		
		err := client.SaveSnapshot(ctx, legacySnapshot)
		require.NoError(t, err)
		
		// Should be able to retrieve using epoch-based methods
		retrieved, err := client.GetEpochSnapshot(ctx, big.NewInt(18), vaultID)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(18), retrieved.EpochNumber)
		assert.Equal(t, "0xlegacyroot", retrieved.MerkleRoot)
	})
}

func TestBadgerClientKeyEncoding(t *testing.T) {
	tempDir := t.TempDir()
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	client, err := NewBadgerClient(logger, tempDir)
	require.NoError(t, err)
	defer client.Close()
	
	t.Run("KeyFormatting", func(t *testing.T) {
		vaultID := "0xABC123def456"
		epochNumber := big.NewInt(42)
		
		key := client.buildSnapshotKey(epochNumber, vaultID)
		expectedKey := "snapshot:vault:0xabc123def456:epoch:00000000000000000042"
		assert.Equal(t, expectedKey, key)
		
		latestKey := client.buildLatestKey(vaultID)
		expectedLatestKey := "latest:vault:0xabc123def456"
		assert.Equal(t, expectedLatestKey, latestKey)
		
		prefix := client.buildVaultPrefix(vaultID)
		expectedPrefix := "snapshot:vault:0xabc123def456:"
		assert.Equal(t, expectedPrefix, prefix)
	})
	
	t.Run("EpochSorting", func(t *testing.T) {
		vaultID := "0xtest"
		ctx := context.Background()
		
		// Save epochs out of order
		epochs := []*big.Int{
			big.NewInt(100),
			big.NewInt(5),
			big.NewInt(50),
			big.NewInt(1),
		}
		
		for i, epoch := range epochs {
			snapshot := MerkleSnapshot{
				Entries: []MerkleEntry{
					{
						Address:     "0x1234567890abcdef",
						TotalEarned: big.NewInt(int64(i * 100)),
					},
				},
				MerkleRoot:  fmt.Sprintf("0x%d", i),
				Timestamp:   time.Now().Unix(),
				VaultID:     vaultID,
				BlockNumber: int64(23534102 + i),
			}
			
			err := client.SaveEpochSnapshot(ctx, epoch, snapshot)
			require.NoError(t, err)
		}
		
		// List should return in descending order (latest first)
		snapshots, err := client.ListEpochSnapshots(ctx, vaultID, 10)
		require.NoError(t, err)
		
		assert.Len(t, snapshots, 4)
		assert.Equal(t, big.NewInt(100), snapshots[0].EpochNumber)
		assert.Equal(t, big.NewInt(50), snapshots[1].EpochNumber)
		assert.Equal(t, big.NewInt(5), snapshots[2].EpochNumber)
		assert.Equal(t, big.NewInt(1), snapshots[3].EpochNumber)
	})
}

func TestBadgerClientConcurrency(t *testing.T) {
	tempDir := t.TempDir()
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	client, err := NewBadgerClient(logger, tempDir)
	require.NoError(t, err)
	defer client.Close()
	
	ctx := context.Background()
	vaultID := "0xf82b93f3d6a703b8b5949809771b1e725708590a"
	
	t.Run("ConcurrentWrites", func(t *testing.T) {
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(epochNum int) {
				defer func() { done <- true }()
				
				snapshot := MerkleSnapshot{
					Entries: []MerkleEntry{
						{
							Address:     "0x3575b992c5337226aecf4e7f93dfbe80c576ce15",
							TotalEarned: big.NewInt(int64(epochNum * 100)),
						},
					},
					MerkleRoot:  fmt.Sprintf("0x%d", epochNum),
					Timestamp:   time.Now().Unix(),
					VaultID:     vaultID,
					BlockNumber: int64(23534102 + epochNum),
				}
				
				err := client.SaveEpochSnapshot(ctx, big.NewInt(int64(epochNum)), snapshot)
				assert.NoError(t, err)
			}(i)
		}
		
		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
		
		// Verify all snapshots were saved
		snapshots, err := client.ListEpochSnapshots(ctx, vaultID, 20)
		require.NoError(t, err)
		assert.Len(t, snapshots, 10)
	})
	
	t.Run("ConcurrentReads", func(t *testing.T) {
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(epochNum int) {
				defer func() { done <- true }()
				
				snapshot, err := client.GetEpochSnapshot(ctx, big.NewInt(int64(epochNum)), vaultID)
				assert.NoError(t, err)
				assert.Equal(t, big.NewInt(int64(epochNum)), snapshot.EpochNumber)
			}(i)
		}
		
		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestBadgerClientMultipleVaults(t *testing.T) {
	tempDir := t.TempDir()
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	client, err := NewBadgerClient(logger, tempDir)
	require.NoError(t, err)
	defer client.Close()
	
	ctx := context.Background()
	vaultID1 := "0xvault1"
	vaultID2 := "0xvault2"
	
	t.Run("IsolatedVaultData", func(t *testing.T) {
		// Save snapshots for both vaults
		snapshot1 := MerkleSnapshot{
			Entries: []MerkleEntry{
				{
					Address:     "0x1111111111111111",
					TotalEarned: big.NewInt(1000),
				},
			},
			MerkleRoot:  "0xvault1root",
			Timestamp:   time.Now().Unix(),
			VaultID:     vaultID1,
			BlockNumber: 100,
		}
		
		snapshot2 := MerkleSnapshot{
			Entries: []MerkleEntry{
				{
					Address:     "0x2222222222222222",
					TotalEarned: big.NewInt(2000),
				},
			},
			MerkleRoot:  "0xvault2root",
			Timestamp:   time.Now().Unix(),
			VaultID:     vaultID2,
			BlockNumber: 200,
		}
		
		err := client.SaveEpochSnapshot(ctx, big.NewInt(1), snapshot1)
		require.NoError(t, err)
		
		err = client.SaveEpochSnapshot(ctx, big.NewInt(1), snapshot2)
		require.NoError(t, err)
		
		// Verify each vault has its own data
		vault1Snapshots, err := client.ListEpochSnapshots(ctx, vaultID1, 10)
		require.NoError(t, err)
		assert.Len(t, vault1Snapshots, 1)
		assert.Equal(t, "0xvault1root", vault1Snapshots[0].MerkleRoot)
		
		vault2Snapshots, err := client.ListEpochSnapshots(ctx, vaultID2, 10)
		require.NoError(t, err)
		assert.Len(t, vault2Snapshots, 1)
		assert.Equal(t, "0xvault2root", vault2Snapshots[0].MerkleRoot)
		
		// Verify latest snapshots are isolated
		latest1, err := client.GetLatestEpochSnapshot(ctx, vaultID1)
		require.NoError(t, err)
		assert.Equal(t, "0xvault1root", latest1.MerkleRoot)
		
		latest2, err := client.GetLatestEpochSnapshot(ctx, vaultID2)
		require.NoError(t, err)
		assert.Equal(t, "0xvault2root", latest2.MerkleRoot)
	})
}

func TestBadgerClientErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	client, err := NewBadgerClient(logger, tempDir)
	require.NoError(t, err)
	defer client.Close()
	
	ctx := context.Background()
	
	t.Run("InvalidEpochNumber", func(t *testing.T) {
		vaultID := "0xtest"
		
		// Try to get snapshot with non-existent epoch
		_, err := client.GetEpochSnapshot(ctx, big.NewInt(999999), vaultID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "snapshot not found")
	})
	
	t.Run("EmptyVaultID", func(t *testing.T) {
		// Save snapshot with empty vault ID should work but create odd keys
		snapshot := MerkleSnapshot{
			Entries: []MerkleEntry{
				{
					Address:     "0x1234567890abcdef",
					TotalEarned: big.NewInt(1000),
				},
			},
			MerkleRoot:  "0xemptytest",
			Timestamp:   time.Now().Unix(),
			VaultID:     "",
			BlockNumber: 100,
		}
		
		err := client.SaveEpochSnapshot(ctx, big.NewInt(1), snapshot)
		require.NoError(t, err)
		
		// Should be able to retrieve it
		retrieved, err := client.GetEpochSnapshot(ctx, big.NewInt(1), "")
		require.NoError(t, err)
		assert.Equal(t, "0xemptytest", retrieved.MerkleRoot)
	})
}