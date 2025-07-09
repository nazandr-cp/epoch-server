package storage

import (
	"context"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEpochStorageClient(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	
	// Create test logger
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	// Create client
	client := NewClientWithBaseDir(logger, tempDir)
	
	// Test data
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
		
		// Check file exists
		expectedPath := filepath.Join(tempDir, "epoch_16", "vault_0xf82b93f3d6a703b8b5949809771b1e725708590a.json")
		assert.FileExists(t, expectedPath)
		
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

func TestEpochStorageClientConcurrency(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	
	// Create test logger
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	// Create client
	client := NewClientWithBaseDir(logger, tempDir)
	
	ctx := context.Background()
	vaultID := "0xf82b93f3d6a703b8b5949809771b1e725708590a"
	
	// Test concurrent writes
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
					MerkleRoot:  "0x" + string(rune(epochNum+'0')),
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
	
	// Test concurrent reads
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

func TestEpochStorageClientFilesystem(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	
	// Create test logger
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	// Create client
	client := NewClientWithBaseDir(logger, tempDir)
	
	ctx := context.Background()
	vaultID := "0xf82b93f3d6a703b8b5949809771b1e725708590a"
	epochNumber := big.NewInt(20)
	
	testSnapshot := MerkleSnapshot{
		Entries: []MerkleEntry{
			{
				Address:     "0x3575b992c5337226aecf4e7f93dfbe80c576ce15",
				TotalEarned: big.NewInt(1000),
			},
		},
		MerkleRoot:  "0x1234567890abcdef",
		Timestamp:   time.Now().Unix(),
		VaultID:     vaultID,
		BlockNumber: 23534102,
	}
	
	t.Run("DirectoryCreation", func(t *testing.T) {
		err := client.SaveEpochSnapshot(ctx, epochNumber, testSnapshot)
		require.NoError(t, err)
		
		// Check that directory structure was created
		epochDir := filepath.Join(tempDir, "epoch_20")
		assert.DirExists(t, epochDir)
		
		vaultFile := filepath.Join(epochDir, "vault_0xf82b93f3d6a703b8b5949809771b1e725708590a.json")
		assert.FileExists(t, vaultFile)
	})
	
	t.Run("FileContent", func(t *testing.T) {
		vaultFile := filepath.Join(tempDir, "epoch_20", "vault_0xf82b93f3d6a703b8b5949809771b1e725708590a.json")
		
		data, err := os.ReadFile(vaultFile)
		require.NoError(t, err)
		
		// Basic JSON validation (epoch number is stored as number, not string)
		assert.Contains(t, string(data), `"epochNumber": 20`)
		assert.Contains(t, string(data), `"vaultId": "0xf82b93f3d6a703b8b5949809771b1e725708590a"`)
		assert.Contains(t, string(data), `"merkleRoot": "0x1234567890abcdef"`)
		assert.Contains(t, string(data), `"blockNumber": 23534102`)
	})
	
	t.Run("InvalidDirectoryHandling", func(t *testing.T) {
		// Create a file where we expect a directory
		badPath := filepath.Join(tempDir, "epoch_bad")
		err := os.WriteFile(badPath, []byte("not a directory"), 0644)
		require.NoError(t, err)
		
		badClient := NewClientWithBaseDir(logger, badPath)
		
		err = badClient.SaveEpochSnapshot(ctx, big.NewInt(100), testSnapshot)
		assert.Error(t, err)
	})
}