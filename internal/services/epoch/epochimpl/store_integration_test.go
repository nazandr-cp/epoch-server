package epochimpl

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

// TestEpochStore_Integration runs comprehensive integration tests for EpochStore using testcontainers
func TestEpochStore_Integration(t *testing.T) {
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
	
	t.Run("EpochLifecycleManagement", func(t *testing.T) {
		testEpochLifecycleManagement(t, ctx, generator)
	})
	
	t.Run("CurrentEpochPointer", func(t *testing.T) {
		testCurrentEpochPointer(t, ctx, generator)
	})
	
	t.Run("MultiVaultCoordination", func(t *testing.T) {
		testMultiVaultCoordination(t, ctx, generator)
	})
	
	t.Run("ConcurrentOperations", func(t *testing.T) {
		testConcurrentOperations(t, ctx, generator)
	})
	
	t.Run("StatusUpdateAtomicity", func(t *testing.T) {
		testStatusUpdateAtomicity(t, ctx, generator)
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
	epochData := generator.GenerateEpochData(vaultID, big.NewInt(1))
	
	epoch := EpochInfo{
		Number:      epochData.Number,
		StartTime:   epochData.StartTime,
		EndTime:     epochData.EndTime,
		BlockNumber: epochData.BlockNumber,
		Status:      epochData.Status,
		VaultID:     epochData.VaultID,
	}
	
	err = store.SaveEpoch(ctx, epoch)
	require.NoError(t, err)
	
	// Verify key was created
	key := helper.CreateEpochKey(vaultID, epochData.Number)
	helper.AssertKeyExists(t, key)
	
	// Verify data can be retrieved
	retrieved, err := store.GetEpoch(ctx, epochData.Number, vaultID)
	require.NoError(t, err)
	assert.Equal(t, epochData.Number, retrieved.Number)
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
	epochData := generator.GenerateEpochData(vaultID, epochNumber)
	epoch := EpochInfo{
		Number:      epochData.Number,
		StartTime:   epochData.StartTime,
		EndTime:     epochData.EndTime,
		BlockNumber: epochData.BlockNumber,
		Status:      epochData.Status,
		VaultID:     epochData.VaultID,
	}
	
	// Save data
	err = store.SaveEpoch(ctx, epoch)
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
	retrieved, err := store.GetEpoch(ctx, epochNumber, vaultID)
	require.NoError(t, err)
	assert.Equal(t, epoch.Status, retrieved.Status)
	assert.Equal(t, epoch.StartTime.Unix(), retrieved.StartTime.Unix())
	assert.Equal(t, epoch.EndTime.Unix(), retrieved.EndTime.Unix())
	assert.Equal(t, epoch.BlockNumber, retrieved.BlockNumber)
}

// testEpochLifecycleManagement tests epoch status transitions
func testEpochLifecycleManagement(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	// Setup test database with container
	db, _, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()
	epochNumber := big.NewInt(1)
	
	// Create initial epoch with "pending" status
	epochData := generator.GenerateEpochData(vaultID, epochNumber)
	epoch := EpochInfo{
		Number:      epochData.Number,
		StartTime:   epochData.StartTime,
		EndTime:     epochData.EndTime,
		BlockNumber: epochData.BlockNumber,
		Status:      "pending",
		VaultID:     epochData.VaultID,
	}
	
	// Save initial epoch
	err = store.SaveEpoch(ctx, epoch)
	require.NoError(t, err)
	
	// Verify initial status
	retrieved, err := store.GetEpoch(ctx, epochNumber, vaultID)
	require.NoError(t, err)
	assert.Equal(t, "pending", retrieved.Status)
	
	// Test status progression: pending -> active -> completed
	statuses := []string{"active", "completed"}
	
	for _, status := range statuses {
		originalUpdatedAt := retrieved.UpdatedAt
		time.Sleep(10 * time.Millisecond) // Ensure timestamp difference
		
		err = store.UpdateEpochStatus(ctx, epochNumber, vaultID, status)
		require.NoError(t, err)
		
		// Verify status update
		retrieved, err = store.GetEpoch(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.Equal(t, status, retrieved.Status)
		assert.True(t, retrieved.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should be updated")
		
		// Verify other fields remain unchanged
		assert.Equal(t, epoch.Number, retrieved.Number)
		assert.Equal(t, epoch.VaultID, retrieved.VaultID)
		assert.Equal(t, epoch.StartTime.Unix(), retrieved.StartTime.Unix())
		assert.Equal(t, epoch.EndTime.Unix(), retrieved.EndTime.Unix())
		assert.Equal(t, epoch.BlockNumber, retrieved.BlockNumber)
	}
}

// testCurrentEpochPointer tests current epoch pointer management
func testCurrentEpochPointer(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	// Setup test database with container
	db, _, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()
	numEpochs := 10
	
	var savedEpochs []EpochInfo
	
	// Save multiple epochs
	for i := 1; i <= numEpochs; i++ {
		epochNumber := big.NewInt(int64(i))
		epochData := generator.GenerateEpochData(vaultID, epochNumber)
		
		epoch := EpochInfo{
			Number:      epochData.Number,
			StartTime:   epochData.StartTime.Add(time.Duration(i) * time.Hour), // Ensure different times
			EndTime:     epochData.EndTime.Add(time.Duration(i) * time.Hour),
			BlockNumber: epochData.BlockNumber + int64(i),
			Status:      epochData.Status,
			VaultID:     epochData.VaultID,
		}
		
		err = store.SaveEpoch(ctx, epoch)
		require.NoError(t, err)
		savedEpochs = append(savedEpochs, epoch)
		
		// Verify current epoch pointer is updated to latest
		current, err := store.GetCurrentEpoch(ctx, vaultID)
		require.NoError(t, err)
		assert.Equal(t, epochNumber, current.Number)
		assert.Equal(t, epoch.Status, current.Status)
	}
	
	// Test that current epoch points to the highest epoch number
	current, err := store.GetCurrentEpoch(ctx, vaultID)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(int64(numEpochs)), current.Number)
	
	// Test listing epochs with different limits
	testCases := []struct {
		limit    int
		expected int
	}{
		{0, numEpochs},      // All
		{5, 5},              // Limited
		{100, numEpochs},    // More than available
	}
	
	for _, tc := range testCases {
		epochs, err := store.ListEpochs(ctx, vaultID, tc.limit)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, len(epochs))
		
		// Verify they're in descending order (latest first)
		if len(epochs) > 1 {
			for i := 1; i < len(epochs); i++ {
				assert.True(t, epochs[i-1].Number.Cmp(epochs[i].Number) > 0,
					"Epochs should be in descending order")
			}
		}
	}
}

// testMultiVaultCoordination tests operations across multiple vaults
func testMultiVaultCoordination(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	// Setup test database with container
	db, _, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	numVaults := 5
	epochsPerVault := 10
	
	vaultIDs := make([]string, numVaults)
	for i := 0; i < numVaults; i++ {
		vaultIDs[i] = generator.GenerateVaultID()
	}
	
	// Create epochs for each vault
	for _, vaultID := range vaultIDs {
		for j := 1; j <= epochsPerVault; j++ {
			epochNumber := big.NewInt(int64(j))
			epochData := generator.GenerateEpochData(vaultID, epochNumber)
			
			epoch := EpochInfo{
				Number:      epochData.Number,
				StartTime:   epochData.StartTime,
				EndTime:     epochData.EndTime,
				BlockNumber: epochData.BlockNumber,
				Status:      epochData.Status,
				VaultID:     epochData.VaultID,
			}
			
			err = store.SaveEpoch(ctx, epoch)
			require.NoError(t, err)
		}
	}
	
	// Verify each vault has its own current epoch
	for _, vaultID := range vaultIDs {
		current, err := store.GetCurrentEpoch(ctx, vaultID)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(int64(epochsPerVault)), current.Number)
		assert.Equal(t, vaultID, current.VaultID)
		
		// Verify correct number of epochs per vault
		epochs, err := store.ListEpochs(ctx, vaultID, 0)
		require.NoError(t, err)
		assert.Equal(t, epochsPerVault, len(epochs))
	}
	
	// Verify vaults are isolated (changing one doesn't affect others)
	testVaultID := vaultIDs[0]
	testEpochNumber := big.NewInt(5)
	
	err = store.UpdateEpochStatus(ctx, testEpochNumber, testVaultID, "completed")
	require.NoError(t, err)
	
	// Verify only the specific vault's epoch was updated
	updatedEpoch, err := store.GetEpoch(ctx, testEpochNumber, testVaultID)
	require.NoError(t, err)
	assert.Equal(t, "completed", updatedEpoch.Status)
	
	// Verify other vaults' epochs remain unchanged
	for _, otherVaultID := range vaultIDs[1:] {
		otherEpoch, err := store.GetEpoch(ctx, testEpochNumber, otherVaultID)
		require.NoError(t, err)
		assert.NotEqual(t, "completed", otherEpoch.Status)
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
				epochData := generator.GenerateEpochData(vaultID, epochNumber)
				
				epoch := EpochInfo{
					Number:      epochData.Number,
					StartTime:   epochData.StartTime,
					EndTime:     epochData.EndTime,
					BlockNumber: epochData.BlockNumber,
					Status:      epochData.Status,
					VaultID:     epochData.VaultID,
				}
				
				if err = store.SaveEpoch(ctx, epoch); err != nil {
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
				// Try to read current epoch
				if _, err := store.GetCurrentEpoch(ctx, vaultID); err != nil {
					// It's okay if no data exists yet
					if !assert.Contains(t, err.Error(), "no current epoch found") {
						errors <- fmt.Errorf("reader %d, op %d: %w", workerID, j, err)
					}
				}
				
				// List epochs
				if _, err := store.ListEpochs(ctx, vaultID, 10); err != nil {
					errors <- fmt.Errorf("reader %d, list %d: %w", workerID, j, err)
				}
			}
		}(i)
	}
	
	// Concurrent status updaters
	time.Sleep(200 * time.Millisecond)
	for i := 0; i < numGoroutines/4; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			for j := 0; j < operationsPerGoroutine/4; j++ {
				epochNumber := big.NewInt(int64(j))
				statuses := []string{"pending", "active", "completed"}
				status := statuses[j%len(statuses)]
				
				if err := store.UpdateEpochStatus(ctx, epochNumber, vaultID, status); err != nil {
					// It's okay if epoch doesn't exist yet
					if !assert.Contains(t, err.Error(), "epoch not found") {
						errors <- fmt.Errorf("updater %d, op %d: %w", workerID, j, err)
					}
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
	epochs, err := store.ListEpochs(ctx, vaultID, 0) // Get all
	require.NoError(t, err)
	assert.Equal(t, numGoroutines*operationsPerGoroutine, len(epochs))
}

// testStatusUpdateAtomicity tests atomic status updates
func testStatusUpdateAtomicity(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)
	
	// Setup test database with container
	db, _, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()
	epochNumber := big.NewInt(1)
	
	// Create initial epoch
	epochData := generator.GenerateEpochData(vaultID, epochNumber)
	epoch := EpochInfo{
		Number:      epochData.Number,
		StartTime:   epochData.StartTime,
		EndTime:     epochData.EndTime,
		BlockNumber: epochData.BlockNumber,
		Status:      "pending",
		VaultID:     epochData.VaultID,
	}
	
	err = store.SaveEpoch(ctx, epoch)
	require.NoError(t, err)
	
	// Test concurrent status updates
	numUpdaters := 10
	statuses := []string{"active", "completed", "failed", "pending"}
	
	var wg sync.WaitGroup
	errors := make(chan error, numUpdaters)
	
	for i := 0; i < numUpdaters; i++ {
		wg.Add(1)
		go func(updaterID int) {
			defer wg.Done()
			
			status := statuses[updaterID%len(statuses)]
			if err := store.UpdateEpochStatus(ctx, epochNumber, vaultID, status); err != nil {
				errors <- fmt.Errorf("updater %d: %w", updaterID, err)
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
		t.Errorf("Status update atomicity test failed with %d errors: %v", len(errorList), errorList[0])
	}
	
	// Verify final state is consistent
	final, err := store.GetEpoch(ctx, epochNumber, vaultID)
	require.NoError(t, err)
	assert.Contains(t, statuses, final.Status) // Should be one of the valid statuses
	assert.Equal(t, epochNumber, final.Number)
	assert.Equal(t, vaultID, final.VaultID)
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
		
		start := time.Now()
		for i := 0; i < numOps; i++ {
			epochNumber := big.NewInt(int64(i))
			epochData := generator.GenerateEpochData(vaultID, epochNumber)
			
			epoch := EpochInfo{
				Number:      epochData.Number,
				StartTime:   epochData.StartTime,
				EndTime:     epochData.EndTime,
				BlockNumber: epochData.BlockNumber,
				Status:      epochData.Status,
				VaultID:     epochData.VaultID,
			}
			
			err = store.SaveEpoch(ctx, epoch)
			require.NoError(t, err)
		}
		duration := time.Since(start)
		
		t.Logf("Saved %d epochs in %v", numOps, duration)
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
			_, err := store.GetEpoch(ctx, epochNumber, vaultID)
			if err != nil && !assert.Contains(t, err.Error(), "epoch not found") {
				require.NoError(t, err)
			}
		}
		duration := time.Since(start)
		
		t.Logf("Performed %d reads in %v", numReads, duration)
		t.Logf("Average time per read: %v", duration/time.Duration(numReads))
	})
	
	// Benchmark status updates
	t.Run("BenchmarkStatusUpdates", func(t *testing.T) {
		numUpdates := 500
		statuses := []string{"pending", "active", "completed"}
		
		start := time.Now()
		for i := 0; i < numUpdates; i++ {
			epochNumber := big.NewInt(int64(i % 100))
			status := statuses[i%len(statuses)]
			
			err := store.UpdateEpochStatus(ctx, epochNumber, vaultID, status)
			if err != nil && !assert.Contains(t, err.Error(), "epoch not found") {
				require.NoError(t, err)
			}
		}
		duration := time.Since(start)
		
		t.Logf("Performed %d status updates in %v", numUpdates, duration)
		t.Logf("Average time per update: %v", duration/time.Duration(numUpdates))
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
	
	t.Run("ZeroEpoch", func(t *testing.T) {
		epochNumber := big.NewInt(0)
		epochData := generator.GenerateEpochData(vaultID, epochNumber)
		
		epoch := EpochInfo{
			Number:      epochData.Number,
			StartTime:   epochData.StartTime,
			EndTime:     epochData.EndTime,
			BlockNumber: epochData.BlockNumber,
			Status:      epochData.Status,
			VaultID:     epochData.VaultID,
		}
		
		err = store.SaveEpoch(ctx, epoch)
		require.NoError(t, err)
		
		retrieved, err := store.GetEpoch(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.Equal(t, epochNumber, retrieved.Number)
	})
	
	t.Run("LargeEpochNumber", func(t *testing.T) {
		// Test with very large epoch number
		epochNumber := new(big.Int)
		epochNumber.SetString("999999999999999999999", 10)
		
		epochData := generator.GenerateEpochData(vaultID, epochNumber)
		epoch := EpochInfo{
			Number:      epochData.Number,
			StartTime:   epochData.StartTime,
			EndTime:     epochData.EndTime,
			BlockNumber: epochData.BlockNumber,
			Status:      epochData.Status,
			VaultID:     epochData.VaultID,
		}
		
		err = store.SaveEpoch(ctx, epoch)
		require.NoError(t, err)
		
		retrieved, err := store.GetEpoch(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.Equal(t, epochNumber, retrieved.Number)
	})
	
	t.Run("InvalidStatus", func(t *testing.T) {
		epochNumber := big.NewInt(99)
		epochData := generator.GenerateEpochData(vaultID, epochNumber)
		
		epoch := EpochInfo{
			Number:      epochData.Number,
			StartTime:   epochData.StartTime,
			EndTime:     epochData.EndTime,
			BlockNumber: epochData.BlockNumber,
			Status:      "invalid_status", // Invalid status
			VaultID:     epochData.VaultID,
		}
		
		// Should still save successfully (validation is business logic, not storage logic)
		err = store.SaveEpoch(ctx, epoch)
		require.NoError(t, err)
		
		retrieved, err := store.GetEpoch(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.Equal(t, "invalid_status", retrieved.Status)
	})
	
	t.Run("NonExistentData", func(t *testing.T) {
		nonExistentVault := generator.GenerateVaultID()
		nonExistentEpoch := big.NewInt(999999)
		
		// Test GetEpoch with non-existent data
		_, err := store.GetEpoch(ctx, nonExistentEpoch, nonExistentVault)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "epoch not found")
		
		// Test GetCurrentEpoch with non-existent vault
		_, err = store.GetCurrentEpoch(ctx, nonExistentVault)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no current epoch found")
		
		// Test UpdateEpochStatus with non-existent epoch
		err = store.UpdateEpochStatus(ctx, nonExistentEpoch, nonExistentVault, "active")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "epoch not found")
		
		// Test ListEpochs with non-existent vault (should return empty list)
		epochs, err := store.ListEpochs(ctx, nonExistentVault, 10)
		require.NoError(t, err)
		assert.Len(t, epochs, 0)
	})
	
	t.Run("TimeZeroValues", func(t *testing.T) {
		epochNumber := big.NewInt(100)
		epoch := EpochInfo{
			Number:      epochNumber,
			StartTime:   time.Time{}, // Zero time
			EndTime:     time.Time{}, // Zero time
			BlockNumber: 0,
			Status:      "pending",
			VaultID:     vaultID,
		}
		
		err = store.SaveEpoch(ctx, epoch)
		require.NoError(t, err)
		
		retrieved, err := store.GetEpoch(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.True(t, retrieved.StartTime.IsZero())
		assert.True(t, retrieved.EndTime.IsZero())
		assert.False(t, retrieved.CreatedAt.IsZero()) // Should be set automatically
		assert.False(t, retrieved.UpdatedAt.IsZero()) // Should be set automatically
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
	
	// Create initial epoch
	epochData := generator.GenerateEpochData(vaultID, epochNumber)
	epoch := EpochInfo{
		Number:      epochData.Number,
		StartTime:   epochData.StartTime,
		EndTime:     epochData.EndTime,
		BlockNumber: epochData.BlockNumber,
		Status:      "pending",
		VaultID:     epochData.VaultID,
	}
	
	err = store.SaveEpoch(ctx, epoch)
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
				retrieved, err := store.GetEpoch(ctx, epochNumber, vaultID)
				if err != nil {
					errors <- fmt.Errorf("reader %d, iteration %d: %w", readerID, j, err)
					return
				}
				
				// Verify data consistency
				if retrieved.Number.Cmp(epochNumber) != 0 {
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
		
		statuses := []string{"active", "completed", "failed"}
		for i, status := range statuses {
			if err := store.UpdateEpochStatus(ctx, epochNumber, vaultID, status); err != nil {
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
	final, err := store.GetEpoch(ctx, epochNumber, vaultID)
	require.NoError(t, err)
	assert.Equal(t, epochNumber, final.Number)
	assert.Equal(t, vaultID, final.VaultID)
	assert.Contains(t, []string{"active", "completed", "failed"}, final.Status)
}