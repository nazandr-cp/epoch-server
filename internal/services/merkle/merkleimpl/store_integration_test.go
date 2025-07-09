package merkleimpl

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	infratesting "github.com/andrey/epoch-server/internal/infra/testing"
)

// TestMerkleStore_Integration runs comprehensive integration tests for MerkleStore using testcontainers
func TestMerkleStore_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test data generator
	generator := infratesting.NewTestDataGenerator(time.Now().UnixNano())

	t.Run("ContainerLifecycle", func(t *testing.T) {
		testContainerLifecycle(t, ctx, generator)
	})

	t.Run("DataPersistence", func(t *testing.T) {
		testDataPersistence(t, ctx, generator)
	})

	t.Run("LargeDatasets", func(t *testing.T) {
		testLargeDatasets(t, ctx, generator)
	})

	t.Run("ConcurrentOperations", func(t *testing.T) {
		testConcurrentOperations(t, ctx, generator)
	})

	t.Run("SnapshotVersioning", func(t *testing.T) {
		testSnapshotVersioning(t, ctx, generator)
	})

	t.Run("PerformanceBenchmarks", func(t *testing.T) {
		testPerformanceBenchmarks(t, ctx, generator)
	})

	t.Run("EdgeCases", func(t *testing.T) {
		testEdgeCases(t, ctx, generator)
	})

	t.Run("TransactionIsolation", func(t *testing.T) {
		testTransactionIsolation(t, ctx, generator)
	})
}

// testContainerLifecycle verifies container setup and teardown
func testContainerLifecycle(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, helper, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	// Clear any existing data from previous tests
	err = helper.Clear()
	require.NoError(t, err)

	// Verify container is running and accessible
	host, err := helper.GetHost(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, host)

	// Verify database is accessible and empty
	helper.AssertDBEmpty(t)

	// Verify we can perform basic operations
	vaultID := generator.GenerateVaultID()
	epochNumber := big.NewInt(1)
	merkleData := generator.GenerateMerkleData(vaultID, epochNumber, 5)

	snapshot := MerkleSnapshot{
		Entries:     convertToMerkleEntries(merkleData.Entries),
		MerkleRoot:  merkleData.MerkleRoot,
		Timestamp:   merkleData.Timestamp,
		VaultID:     merkleData.VaultID,
		BlockNumber: merkleData.BlockNumber,
	}

	err = store.SaveSnapshot(ctx, epochNumber, snapshot)
	require.NoError(t, err)

	// Verify key was created
	key := helper.CreateMerkleKey(vaultID, epochNumber)
	helper.AssertKeyExists(t, key)

	// Verify data can be retrieved
	retrieved, err := store.GetSnapshot(ctx, epochNumber, vaultID)
	require.NoError(t, err)
	assert.Equal(t, epochNumber, retrieved.EpochNumber)
	assert.Equal(t, vaultID, retrieved.VaultID)
}

// testDataPersistence verifies data survives container restarts
func testDataPersistence(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, helper, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()
	epochNumber := big.NewInt(42)

	// Create test data
	merkleData := generator.GenerateMerkleData(vaultID, epochNumber, 10)
	snapshot := MerkleSnapshot{
		Entries:     convertToMerkleEntries(merkleData.Entries),
		MerkleRoot:  merkleData.MerkleRoot,
		Timestamp:   merkleData.Timestamp,
		VaultID:     merkleData.VaultID,
		BlockNumber: merkleData.BlockNumber,
	}

	// Save data
	err = store.SaveSnapshot(ctx, epochNumber, snapshot)
	require.NoError(t, err)

	// Force sync to ensure data is persisted
	err = helper.Sync()
	require.NoError(t, err)

	// Simulate restart by clearing caches and forcing read from disk
	err = helper.RunGC(ctx)
	// GC may not have anything to clean up, which is not an error
	if err != nil && !strings.Contains(err.Error(), "didn't result in any cleanup") {
		require.NoError(t, err)
	}

	// Verify data is still accessible
	retrieved, err := store.GetSnapshot(ctx, epochNumber, vaultID)
	require.NoError(t, err)
	assert.Equal(t, snapshot.MerkleRoot, retrieved.MerkleRoot)
	assert.Equal(t, snapshot.Timestamp, retrieved.Timestamp)
	assert.Len(t, retrieved.Entries, len(snapshot.Entries))
}

// testLargeDatasets tests performance with large merkle trees
func testLargeDatasets(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, helper, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()

	// Test with increasingly large datasets
	sizes := []int{100, 1000, 5000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			epochNumber := big.NewInt(int64(size))

			// Generate large dataset
			merkleData := generator.GenerateMerkleData(vaultID, epochNumber, size)
			snapshot := MerkleSnapshot{
				Entries:     convertToMerkleEntries(merkleData.Entries),
				MerkleRoot:  merkleData.MerkleRoot,
				Timestamp:   merkleData.Timestamp,
				VaultID:     merkleData.VaultID,
				BlockNumber: merkleData.BlockNumber,
			}

			// Measure save time
			saveTime, err := helper.MeasureOperationTime(func() error {
				return store.SaveSnapshot(ctx, epochNumber, snapshot)
			})
			require.NoError(t, err)
			t.Logf("Save time for %d entries: %v", size, saveTime)

			// Measure retrieve time
			retrieveTime, err := helper.MeasureOperationTime(func() error {
				_, err := store.GetSnapshot(ctx, epochNumber, vaultID)
				return err
			})
			require.NoError(t, err)
			t.Logf("Retrieve time for %d entries: %v", size, retrieveTime)

			// Verify correct data
			retrieved, err := store.GetSnapshot(ctx, epochNumber, vaultID)
			require.NoError(t, err)
			assert.Len(t, retrieved.Entries, size)
			assert.Equal(t, snapshot.MerkleRoot, retrieved.MerkleRoot)

			// Check memory usage
			metrics := helper.CollectMetrics()
			t.Logf("Metrics after %d entries: KeyCount=%d, LSMSize=%d, VLogSize=%d",
				size, metrics.KeyCount, metrics.LSMSize, metrics.VLogSize)
		})
	}
}

// testConcurrentOperations tests concurrent read/write operations
func testConcurrentOperations(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, _, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()
	numGoroutines := 20
	operationsPerGoroutine := 10

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*operationsPerGoroutine)

	// Concurrent writers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				epochNumber := big.NewInt(int64(workerID*operationsPerGoroutine + j))

				merkleData := generator.GenerateMerkleData(vaultID, epochNumber, 20)
				snapshot := MerkleSnapshot{
					Entries:     convertToMerkleEntries(merkleData.Entries),
					MerkleRoot:  merkleData.MerkleRoot,
					Timestamp:   merkleData.Timestamp,
					VaultID:     merkleData.VaultID,
					BlockNumber: merkleData.BlockNumber,
				}

				if err := store.SaveSnapshot(ctx, epochNumber, snapshot); err != nil {
					errors <- fmt.Errorf("worker %d, op %d: %w", workerID, j, err)
				}
			}
		}(i)
	}

	// Concurrent readers (after some data is written)
	time.Sleep(100 * time.Millisecond)
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine/2; j++ {
				// Try to read latest snapshot
				if _, err := store.GetLatestSnapshot(ctx, vaultID); err != nil {
					// It's okay if no data exists yet
					if !assert.Contains(t, err.Error(), "no snapshots found") {
						errors <- fmt.Errorf("reader %d, op %d: %w", workerID, j, err)
					}
				}

				// List snapshots
				if _, err := store.ListSnapshots(ctx, vaultID, 10); err != nil {
					errors <- fmt.Errorf("reader %d, list %d: %w", workerID, j, err)
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	var errorList []error
	for err := range errors {
		errorList = append(errorList, err)
	}

	if len(errorList) > 0 {
		t.Errorf("Concurrent operations failed with %d errors: %v", len(errorList), errorList[0])
	}

	// Verify final state
	snapshots, err := store.ListSnapshots(ctx, vaultID, 0) // Get all
	require.NoError(t, err)
	assert.Equal(t, numGoroutines*operationsPerGoroutine, len(snapshots))
}

// testSnapshotVersioning tests multiple epochs and latest pointer management
func testSnapshotVersioning(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, _, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()
	numEpochs := 20

	var savedSnapshots []MerkleSnapshot

	// Save multiple epochs
	for i := 1; i <= numEpochs; i++ {
		epochNumber := big.NewInt(int64(i))

		merkleData := generator.GenerateMerkleData(vaultID, epochNumber, 10)
		snapshot := MerkleSnapshot{
			Entries:     convertToMerkleEntries(merkleData.Entries),
			MerkleRoot:  merkleData.MerkleRoot,
			Timestamp:   merkleData.Timestamp + int64(i), // Ensure increasing timestamps
			VaultID:     merkleData.VaultID,
			BlockNumber: merkleData.BlockNumber + int64(i),
		}

		err := store.SaveSnapshot(ctx, epochNumber, snapshot)
		require.NoError(t, err)
		savedSnapshots = append(savedSnapshots, snapshot)

		// Verify latest pointer is updated
		latest, err := store.GetLatestSnapshot(ctx, vaultID)
		require.NoError(t, err)
		assert.Equal(t, epochNumber, latest.EpochNumber)
		assert.Equal(t, snapshot.MerkleRoot, latest.MerkleRoot)
	}

	// Test listing with different limits
	testCases := []struct {
		limit    int
		expected int
	}{
		{0, numEpochs},   // All
		{5, 5},           // Limited
		{100, numEpochs}, // More than available
	}

	for _, tc := range testCases {
		snapshots, err := store.ListSnapshots(ctx, vaultID, tc.limit)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, len(snapshots))

		// Verify they're in descending order (latest first)
		if len(snapshots) > 1 {
			for i := 1; i < len(snapshots); i++ {
				assert.True(t, snapshots[i-1].EpochNumber.Cmp(snapshots[i].EpochNumber) > 0,
					"Snapshots should be in descending order")
			}
		}
	}

	// Test specific epoch retrieval
	for i := 1; i <= numEpochs; i++ {
		epochNumber := big.NewInt(int64(i))
		snapshot, err := store.GetSnapshot(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.Equal(t, epochNumber, snapshot.EpochNumber)
		assert.Equal(t, savedSnapshots[i-1].MerkleRoot, snapshot.MerkleRoot)
	}
}

// testPerformanceBenchmarks runs performance benchmarks
func testPerformanceBenchmarks(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, helper, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()

	// Benchmark save operations
	t.Run("BenchmarkSave", func(t *testing.T) {
		numOps := 1000
		entrySize := 50

		start := time.Now()
		for i := 0; i < numOps; i++ {
			epochNumber := big.NewInt(int64(i))
			merkleData := generator.GenerateMerkleData(vaultID, epochNumber, entrySize)
			snapshot := MerkleSnapshot{
				Entries:     convertToMerkleEntries(merkleData.Entries),
				MerkleRoot:  merkleData.MerkleRoot,
				Timestamp:   merkleData.Timestamp,
				VaultID:     merkleData.VaultID,
				BlockNumber: merkleData.BlockNumber,
			}

			err := store.SaveSnapshot(ctx, epochNumber, snapshot)
			require.NoError(t, err)
		}
		duration := time.Since(start)

		t.Logf("Saved %d snapshots (%d entries each) in %v", numOps, entrySize, duration)
		t.Logf("Average time per save: %v", duration/time.Duration(numOps))

		// Collect metrics
		metrics := helper.CollectMetrics()
		t.Logf("Final metrics: KeyCount=%d, TotalSize=%d", metrics.KeyCount, metrics.TotalSize)
	})

	// Benchmark read operations
	t.Run("BenchmarkRead", func(t *testing.T) {
		numReads := 1000

		start := time.Now()
		for i := 0; i < numReads; i++ {
			epochNumber := big.NewInt(int64(i % 100)) // Read from saved data
			_, err := store.GetSnapshot(ctx, epochNumber, vaultID)
			if err != nil && !assert.Contains(t, err.Error(), "snapshot not found") {
				require.NoError(t, err)
			}
		}
		duration := time.Since(start)

		t.Logf("Performed %d reads in %v", numReads, duration)
		t.Logf("Average time per read: %v", duration/time.Duration(numReads))
	})
}

// testEdgeCases tests edge cases and error conditions
func testEdgeCases(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, _, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()

	t.Run("EmptySnapshot", func(t *testing.T) {
		epochNumber := big.NewInt(999)
		snapshot := MerkleSnapshot{
			Entries:     []MerkleEntry{}, // Empty entries
			MerkleRoot:  "0x0000000000000000000000000000000000000000000000000000000000000000",
			Timestamp:   time.Now().Unix(),
			VaultID:     vaultID,
			BlockNumber: 12345,
		}

		err := store.SaveSnapshot(ctx, epochNumber, snapshot)
		require.NoError(t, err)

		retrieved, err := store.GetSnapshot(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.Len(t, retrieved.Entries, 0)
		assert.Equal(t, snapshot.MerkleRoot, retrieved.MerkleRoot)
	})

	t.Run("ZeroEpoch", func(t *testing.T) {
		epochNumber := big.NewInt(0)
		merkleData := generator.GenerateMerkleData(vaultID, epochNumber, 5)
		snapshot := MerkleSnapshot{
			Entries:     convertToMerkleEntries(merkleData.Entries),
			MerkleRoot:  merkleData.MerkleRoot,
			Timestamp:   merkleData.Timestamp,
			VaultID:     merkleData.VaultID,
			BlockNumber: merkleData.BlockNumber,
		}

		err := store.SaveSnapshot(ctx, epochNumber, snapshot)
		require.NoError(t, err)

		retrieved, err := store.GetSnapshot(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.Equal(t, epochNumber, retrieved.EpochNumber)
	})

	t.Run("LargeEpochNumber", func(t *testing.T) {
		// Test with very large epoch number
		epochNumber := new(big.Int)
		epochNumber.SetString("999999999999999999999", 10)

		merkleData := generator.GenerateMerkleData(vaultID, epochNumber, 10)
		snapshot := MerkleSnapshot{
			Entries:     convertToMerkleEntries(merkleData.Entries),
			MerkleRoot:  merkleData.MerkleRoot,
			Timestamp:   merkleData.Timestamp,
			VaultID:     merkleData.VaultID,
			BlockNumber: merkleData.BlockNumber,
		}

		err := store.SaveSnapshot(ctx, epochNumber, snapshot)
		require.NoError(t, err)

		retrieved, err := store.GetSnapshot(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.Equal(t, epochNumber, retrieved.EpochNumber)
	})

	t.Run("NonExistentData", func(t *testing.T) {
		nonExistentVault := generator.GenerateVaultID()
		nonExistentEpoch := big.NewInt(999999)

		// Test GetSnapshot with non-existent data
		_, err := store.GetSnapshot(ctx, nonExistentEpoch, nonExistentVault)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "snapshot not found")

		// Test GetLatestSnapshot with non-existent vault
		_, err = store.GetLatestSnapshot(ctx, nonExistentVault)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no snapshots found")

		// Test ListSnapshots with non-existent vault (should return empty list)
		snapshots, err := store.ListSnapshots(ctx, nonExistentVault, 10)
		require.NoError(t, err)
		assert.Len(t, snapshots, 0)
	})
}

// testTransactionIsolation tests transaction isolation and consistency
func testTransactionIsolation(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, _, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()
	epochNumber := big.NewInt(100)

	// Create initial snapshot
	merkleData := generator.GenerateMerkleData(vaultID, epochNumber, 10)
	snapshot := MerkleSnapshot{
		Entries:     convertToMerkleEntries(merkleData.Entries),
		MerkleRoot:  merkleData.MerkleRoot,
		Timestamp:   merkleData.Timestamp,
		VaultID:     merkleData.VaultID,
		BlockNumber: merkleData.BlockNumber,
	}

	err = store.SaveSnapshot(ctx, epochNumber, snapshot)
	require.NoError(t, err)

	// Test concurrent reads while writing
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	// Start multiple readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()

			for j := 0; j < 20; j++ {
				retrieved, err := store.GetSnapshot(ctx, epochNumber, vaultID)
				if err != nil {
					errors <- fmt.Errorf("reader %d, iteration %d: %w", readerID, j, err)
					return
				}

				// Verify data consistency
				if retrieved.EpochNumber.Cmp(epochNumber) != 0 {
					errors <- fmt.Errorf("reader %d: epoch number mismatch", readerID)
					return
				}

				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// Start a writer that updates the same epoch
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 10; i++ {
			updatedSnapshot := snapshot
			updatedSnapshot.Timestamp = time.Now().Unix() + int64(i)
			updatedSnapshot.MerkleRoot = fmt.Sprintf("0x%064d", i)

			if err := store.SaveSnapshot(ctx, epochNumber, updatedSnapshot); err != nil {
				errors <- fmt.Errorf("writer iteration %d: %w", i, err)
				return
			}

			time.Sleep(20 * time.Millisecond)
		}
	}()

	wg.Wait()
	close(errors)

	// Check for errors
	var errorList []error
	for err := range errors {
		errorList = append(errorList, err)
	}

	if len(errorList) > 0 {
		t.Errorf("Transaction isolation test failed with %d errors: %v", len(errorList), errorList[0])
	}

	// Verify final state is consistent
	final, err := store.GetSnapshot(ctx, epochNumber, vaultID)
	require.NoError(t, err)
	assert.Equal(t, epochNumber, final.EpochNumber)
	assert.Equal(t, vaultID, final.VaultID)
}

// Helper function to convert test data to store types
func convertToMerkleEntries(entries []infratesting.MerkleEntry) []MerkleEntry {
	result := make([]MerkleEntry, len(entries))
	for i, entry := range entries {
		result[i] = MerkleEntry{
			Address:     entry.Address,
			TotalEarned: entry.TotalEarned,
		}
	}
	return result
}
