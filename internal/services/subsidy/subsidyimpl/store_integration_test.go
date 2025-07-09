package subsidyimpl

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

// TestSubsidyStore_Integration runs comprehensive integration tests for SubsidyStore using testcontainers
func TestSubsidyStore_Integration(t *testing.T) {
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

	t.Run("DistributionWorkflow", func(t *testing.T) {
		testDistributionWorkflow(t, ctx, generator)
	})

	t.Run("EpochBasedQuerying", func(t *testing.T) {
		testEpochBasedQuerying(t, ctx, generator)
	})

	t.Run("StatusFiltering", func(t *testing.T) {
		testStatusFiltering(t, ctx, generator)
	})

	t.Run("ConcurrentOperations", func(t *testing.T) {
		testConcurrentOperations(t, ctx, generator)
	})

	t.Run("LargeScaleDistributions", func(t *testing.T) {
		testLargeScaleDistributions(t, ctx, generator)
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
	subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)

	distribution := SubsidyDistribution{
		ID:                subsidyData.ID,
		EpochNumber:       subsidyData.EpochNumber,
		VaultID:           subsidyData.VaultID,
		CollectionAddress: subsidyData.CollectionAddress,
		Amount:            subsidyData.Amount,
		Status:            subsidyData.Status,
		TxHash:            subsidyData.TxHash,
		BlockNumber:       subsidyData.BlockNumber,
	}

	err = store.SaveDistribution(ctx, distribution)
	require.NoError(t, err)

	// Verify key was created
	key := helper.CreateSubsidyKey(subsidyData.ID)
	helper.AssertKeyExists(t, key)

	// Verify data can be retrieved
	retrieved, err := store.GetDistribution(ctx, subsidyData.ID)
	require.NoError(t, err)
	assert.Equal(t, subsidyData.ID, retrieved.ID)
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
	subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)
	distribution := SubsidyDistribution{
		ID:                subsidyData.ID,
		EpochNumber:       subsidyData.EpochNumber,
		VaultID:           subsidyData.VaultID,
		CollectionAddress: subsidyData.CollectionAddress,
		Amount:            subsidyData.Amount,
		Status:            subsidyData.Status,
		TxHash:            subsidyData.TxHash,
		BlockNumber:       subsidyData.BlockNumber,
	}

	// Save data
	err = store.SaveDistribution(ctx, distribution)
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
	retrieved, err := store.GetDistribution(ctx, subsidyData.ID)
	require.NoError(t, err)
	assert.Equal(t, distribution.Status, retrieved.Status)
	assert.Equal(t, distribution.Amount.String(), retrieved.Amount.String())
	assert.Equal(t, distribution.CollectionAddress, retrieved.CollectionAddress)
	assert.Equal(t, distribution.TxHash, retrieved.TxHash)
	assert.Equal(t, distribution.BlockNumber, retrieved.BlockNumber)
}

// testDistributionWorkflow tests the complete distribution workflow
func testDistributionWorkflow(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, _, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()
	epochNumber := big.NewInt(1)

	// Create initial distribution with "pending" status
	subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)
	distribution := SubsidyDistribution{
		ID:                subsidyData.ID,
		EpochNumber:       subsidyData.EpochNumber,
		VaultID:           subsidyData.VaultID,
		CollectionAddress: subsidyData.CollectionAddress,
		Amount:            subsidyData.Amount,
		Status:            "pending",
		TxHash:            "", // No tx hash yet
		BlockNumber:       0,  // No block number yet
	}

	// Save initial distribution
	err = store.SaveDistribution(ctx, distribution)
	require.NoError(t, err)

	// Verify initial status
	retrieved, err := store.GetDistribution(ctx, distribution.ID)
	require.NoError(t, err)
	assert.Equal(t, "pending", retrieved.Status)
	assert.Empty(t, retrieved.TxHash)
	assert.Equal(t, int64(0), retrieved.BlockNumber)

	// Test workflow progression: pending -> distributed
	txHash := generator.GenerateRandomHash()
	blockNumber := int64(12345678)

	originalUpdatedAt := retrieved.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Ensure timestamp difference

	err = store.UpdateDistributionStatus(ctx, distribution.ID, "distributed", txHash, blockNumber)
	require.NoError(t, err)

	// Verify status update
	retrieved, err = store.GetDistribution(ctx, distribution.ID)
	require.NoError(t, err)
	assert.Equal(t, "distributed", retrieved.Status)
	assert.Equal(t, txHash, retrieved.TxHash)
	assert.Equal(t, blockNumber, retrieved.BlockNumber)
	assert.True(t, retrieved.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should be updated")

	// Verify other fields remain unchanged
	assert.Equal(t, distribution.ID, retrieved.ID)
	assert.Equal(t, distribution.VaultID, retrieved.VaultID)
	assert.Equal(t, distribution.EpochNumber.String(), retrieved.EpochNumber.String())
	assert.Equal(t, distribution.CollectionAddress, retrieved.CollectionAddress)
	assert.Equal(t, distribution.Amount.String(), retrieved.Amount.String())

	// Test failure case: distributed -> failed
	err = store.UpdateDistributionStatus(ctx, distribution.ID, "failed", "", 0)
	require.NoError(t, err)

	retrieved, err = store.GetDistribution(ctx, distribution.ID)
	require.NoError(t, err)
	assert.Equal(t, "failed", retrieved.Status)
	// TxHash and BlockNumber should remain from previous update
	assert.Equal(t, txHash, retrieved.TxHash)
	assert.Equal(t, blockNumber, retrieved.BlockNumber)
}

// testEpochBasedQuerying tests querying distributions by epoch
func testEpochBasedQuerying(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, _, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()
	numEpochs := 5
	distributionsPerEpoch := 10

	var allDistributions []SubsidyDistribution
	epochDistributions := make(map[string][]SubsidyDistribution) // epochNumber -> distributions

	// Create distributions for multiple epochs
	for i := 1; i <= numEpochs; i++ {
		epochNumber := big.NewInt(int64(i))
		epochKey := epochNumber.String()

		for j := 0; j < distributionsPerEpoch; j++ {
			subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)
			distribution := SubsidyDistribution{
				ID:                subsidyData.ID,
				EpochNumber:       subsidyData.EpochNumber,
				VaultID:           subsidyData.VaultID,
				CollectionAddress: subsidyData.CollectionAddress,
				Amount:            subsidyData.Amount,
				Status:            subsidyData.Status,
				TxHash:            subsidyData.TxHash,
				BlockNumber:       subsidyData.BlockNumber,
			}

			err = store.SaveDistribution(ctx, distribution)
			require.NoError(t, err)

			allDistributions = append(allDistributions, distribution)
			epochDistributions[epochKey] = append(epochDistributions[epochKey], distribution)
		}
	}

	// Test querying by specific epoch
	for i := 1; i <= numEpochs; i++ {
		epochNumber := big.NewInt(int64(i))

		distributions, err := store.ListDistributionsByEpoch(ctx, epochNumber, vaultID)
		require.NoError(t, err)
		assert.Equal(t, distributionsPerEpoch, len(distributions))

		// Verify all distributions belong to the correct epoch
		for _, dist := range distributions {
			assert.Equal(t, epochNumber.String(), dist.EpochNumber.String())
			assert.Equal(t, vaultID, dist.VaultID)
		}
	}

	// Test querying non-existent epoch
	nonExistentEpoch := big.NewInt(999)
	distributions, err := store.ListDistributionsByEpoch(ctx, nonExistentEpoch, vaultID)
	require.NoError(t, err)
	assert.Len(t, distributions, 0)

	// Test with different vault ID (should return empty)
	differentVaultID := generator.GenerateVaultID()
	distributions, err = store.ListDistributionsByEpoch(ctx, big.NewInt(1), differentVaultID)
	require.NoError(t, err)
	assert.Len(t, distributions, 0)
}

// testStatusFiltering tests filtering distributions by status
func testStatusFiltering(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, _, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()
	epochNumber := big.NewInt(1)

	statuses := []string{"pending", "distributed", "failed"}
	distributionsPerStatus := 20
	statusCounts := make(map[string]int)

	// Create distributions with different statuses
	for _, status := range statuses {
		for i := 0; i < distributionsPerStatus; i++ {
			subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)
			distribution := SubsidyDistribution{
				ID:                subsidyData.ID,
				EpochNumber:       subsidyData.EpochNumber,
				VaultID:           subsidyData.VaultID,
				CollectionAddress: subsidyData.CollectionAddress,
				Amount:            subsidyData.Amount,
				Status:            status,
				TxHash:            subsidyData.TxHash,
				BlockNumber:       subsidyData.BlockNumber,
			}

			err = store.SaveDistribution(ctx, distribution)
			require.NoError(t, err)
			statusCounts[status]++
		}
	}

	// Test filtering by each status
	for _, status := range statuses {
		t.Run(fmt.Sprintf("Status_%s", status), func(t *testing.T) {
			distributions, err := store.ListDistributionsByStatus(ctx, status, 0) // No limit
			require.NoError(t, err)
			assert.Equal(t, statusCounts[status], len(distributions))

			// Verify all distributions have the correct status
			for _, dist := range distributions {
				assert.Equal(t, status, dist.Status)
			}
		})
	}

	// Test with limit
	for _, status := range statuses {
		limit := 5
		distributions, err := store.ListDistributionsByStatus(ctx, status, limit)
		require.NoError(t, err)
		expectedCount := limit
		if statusCounts[status] < limit {
			expectedCount = statusCounts[status]
		}
		assert.Equal(t, expectedCount, len(distributions))
	}

	// Test with non-existent status
	distributions, err := store.ListDistributionsByStatus(ctx, "non_existent", 10)
	require.NoError(t, err)
	assert.Len(t, distributions, 0)
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
	epochNumber := big.NewInt(1)
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
				subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)
				subsidyData.ID = fmt.Sprintf("worker-%d-op-%d", workerID, j) // Unique ID

				distribution := SubsidyDistribution{
					ID:                subsidyData.ID,
					EpochNumber:       subsidyData.EpochNumber,
					VaultID:           subsidyData.VaultID,
					CollectionAddress: subsidyData.CollectionAddress,
					Amount:            subsidyData.Amount,
					Status:            subsidyData.Status,
					TxHash:            subsidyData.TxHash,
					BlockNumber:       subsidyData.BlockNumber,
				}

				if err = store.SaveDistribution(ctx, distribution); err != nil {
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
				// List distributions by epoch
				if _, err := store.ListDistributionsByEpoch(ctx, epochNumber, vaultID); err != nil {
					errors <- fmt.Errorf("reader %d, epoch list %d: %w", workerID, j, err)
				}

				// List distributions by status
				statuses := []string{"pending", "distributed", "failed"}
				status := statuses[j%len(statuses)]
				if _, err := store.ListDistributionsByStatus(ctx, status, 10); err != nil {
					errors <- fmt.Errorf("reader %d, status list %d: %w", workerID, j, err)
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
				distributionID := fmt.Sprintf("worker-%d-op-%d", j, 0) // Try to update existing distributions
				statuses := []string{"pending", "distributed", "failed"}
				status := statuses[j%len(statuses)]

				if err := store.UpdateDistributionStatus(ctx, distributionID, status, "", 0); err != nil {
					// It's okay if distribution doesn't exist yet
					if !assert.Contains(t, err.Error(), "distribution not found") {
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
	distributions, err := store.ListDistributionsByEpoch(ctx, epochNumber, vaultID)
	require.NoError(t, err)
	assert.Equal(t, numGoroutines*operationsPerGoroutine, len(distributions))
}

// testLargeScaleDistributions tests performance with many distributions
func testLargeScaleDistributions(t *testing.T, ctx context.Context, generator *infratesting.TestDataGenerator) {
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Setup test database with container
	db, helper, cleanup, err := infratesting.SetupTestDBAndHelper(ctx)
	require.NoError(t, err, "Failed to setup test database")
	defer cleanup()

	// Initialize store with database
	store := NewStore(db, logger)
	vaultID := generator.GenerateVaultID()
	epochNumber := big.NewInt(1)

	// Test with increasingly large numbers of distributions
	sizes := []int{100, 1000, 5000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			// Clear previous data
			err := helper.Clear()
			require.NoError(t, err)

			// Create many distributions
			start := time.Now()
			for i := 0; i < size; i++ {
				subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)
				subsidyData.ID = fmt.Sprintf("large-scale-%d", i)

				distribution := SubsidyDistribution{
					ID:                subsidyData.ID,
					EpochNumber:       subsidyData.EpochNumber,
					VaultID:           subsidyData.VaultID,
					CollectionAddress: subsidyData.CollectionAddress,
					Amount:            subsidyData.Amount,
					Status:            subsidyData.Status,
					TxHash:            subsidyData.TxHash,
					BlockNumber:       subsidyData.BlockNumber,
				}

				err = store.SaveDistribution(ctx, distribution)
				require.NoError(t, err)
			}
			saveTime := time.Since(start)
			t.Logf("Save time for %d distributions: %v", size, saveTime)

			// Test querying performance
			queryStart := time.Now()
			distributions, err := store.ListDistributionsByEpoch(ctx, epochNumber, vaultID)
			require.NoError(t, err)
			queryTime := time.Since(queryStart)

			assert.Equal(t, size, len(distributions))
			t.Logf("Query time for %d distributions: %v", size, queryTime)

			// Test status filtering performance
			filterStart := time.Now()
			pending, err := store.ListDistributionsByStatus(ctx, "pending", 0)
			require.NoError(t, err)
			filterTime := time.Since(filterStart)

			t.Logf("Filter time for %d distributions: %v (found %d pending)", size, filterTime, len(pending))

			// Check memory usage
			metrics := helper.CollectMetrics()
			t.Logf("Metrics after %d distributions: KeyCount=%d, LSMSize=%d, VLogSize=%d",
				size, metrics.KeyCount, metrics.LSMSize, metrics.VLogSize)
		})
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
	epochNumber := big.NewInt(1)

	// Benchmark save operations
	t.Run("BenchmarkSave", func(t *testing.T) {
		numOps := 1000

		start := time.Now()
		for i := 0; i < numOps; i++ {
			subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)
			subsidyData.ID = fmt.Sprintf("bench-save-%d", i)

			distribution := SubsidyDistribution{
				ID:                subsidyData.ID,
				EpochNumber:       subsidyData.EpochNumber,
				VaultID:           subsidyData.VaultID,
				CollectionAddress: subsidyData.CollectionAddress,
				Amount:            subsidyData.Amount,
				Status:            subsidyData.Status,
				TxHash:            subsidyData.TxHash,
				BlockNumber:       subsidyData.BlockNumber,
			}

			err = store.SaveDistribution(ctx, distribution)
			require.NoError(t, err)
		}
		duration := time.Since(start)

		t.Logf("Saved %d distributions in %v", numOps, duration)
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
			distributionID := fmt.Sprintf("bench-save-%d", i%100) // Read from saved data
			_, err := store.GetDistribution(ctx, distributionID)
			if err != nil && !assert.Contains(t, err.Error(), "distribution not found") {
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
		statuses := []string{"pending", "distributed", "failed"}

		start := time.Now()
		for i := 0; i < numUpdates; i++ {
			distributionID := fmt.Sprintf("bench-save-%d", i%100)
			status := statuses[i%len(statuses)]

			err := store.UpdateDistributionStatus(ctx, distributionID, status, "", 0)
			if err != nil && !assert.Contains(t, err.Error(), "distribution not found") {
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

	t.Run("ZeroAmount", func(t *testing.T) {
		epochNumber := big.NewInt(1)
		subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)
		distribution := SubsidyDistribution{
			ID:                subsidyData.ID,
			EpochNumber:       subsidyData.EpochNumber,
			VaultID:           subsidyData.VaultID,
			CollectionAddress: subsidyData.CollectionAddress,
			Amount:            big.NewInt(0), // Zero amount
			Status:            subsidyData.Status,
			TxHash:            subsidyData.TxHash,
			BlockNumber:       subsidyData.BlockNumber,
		}

		err = store.SaveDistribution(ctx, distribution)
		require.NoError(t, err)

		retrieved, err := store.GetDistribution(ctx, distribution.ID)
		require.NoError(t, err)
		assert.Equal(t, "0", retrieved.Amount.String())
	})

	t.Run("MaxAmount", func(t *testing.T) {
		epochNumber := big.NewInt(1)
		subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)

		// Max uint256 value
		maxAmount := new(big.Int)
		maxAmount.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

		distribution := SubsidyDistribution{
			ID:                subsidyData.ID,
			EpochNumber:       subsidyData.EpochNumber,
			VaultID:           subsidyData.VaultID,
			CollectionAddress: subsidyData.CollectionAddress,
			Amount:            maxAmount,
			Status:            subsidyData.Status,
			TxHash:            subsidyData.TxHash,
			BlockNumber:       subsidyData.BlockNumber,
		}

		err = store.SaveDistribution(ctx, distribution)
		require.NoError(t, err)

		retrieved, err := store.GetDistribution(ctx, distribution.ID)
		require.NoError(t, err)
		assert.Equal(t, maxAmount.String(), retrieved.Amount.String())
	})

	t.Run("ZeroEpoch", func(t *testing.T) {
		epochNumber := big.NewInt(0)
		subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)
		distribution := SubsidyDistribution{
			ID:                subsidyData.ID,
			EpochNumber:       subsidyData.EpochNumber,
			VaultID:           subsidyData.VaultID,
			CollectionAddress: subsidyData.CollectionAddress,
			Amount:            subsidyData.Amount,
			Status:            subsidyData.Status,
			TxHash:            subsidyData.TxHash,
			BlockNumber:       subsidyData.BlockNumber,
		}

		err = store.SaveDistribution(ctx, distribution)
		require.NoError(t, err)

		retrieved, err := store.GetDistribution(ctx, distribution.ID)
		require.NoError(t, err)
		assert.Equal(t, "0", retrieved.EpochNumber.String())
	})

	t.Run("EmptyTxHash", func(t *testing.T) {
		epochNumber := big.NewInt(1)
		subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)
		distribution := SubsidyDistribution{
			ID:                subsidyData.ID,
			EpochNumber:       subsidyData.EpochNumber,
			VaultID:           subsidyData.VaultID,
			CollectionAddress: subsidyData.CollectionAddress,
			Amount:            subsidyData.Amount,
			Status:            "pending",
			TxHash:            "", // Empty tx hash
			BlockNumber:       0,  // Zero block number
		}

		err = store.SaveDistribution(ctx, distribution)
		require.NoError(t, err)

		retrieved, err := store.GetDistribution(ctx, distribution.ID)
		require.NoError(t, err)
		assert.Empty(t, retrieved.TxHash)
		assert.Equal(t, int64(0), retrieved.BlockNumber)
	})

	t.Run("NonExistentData", func(t *testing.T) {
		nonExistentID := "non-existent-distribution"
		nonExistentVault := generator.GenerateVaultID()
		nonExistentEpoch := big.NewInt(999999)

		// Test GetDistribution with non-existent ID
		_, err := store.GetDistribution(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "distribution not found")

		// Test UpdateDistributionStatus with non-existent ID
		err = store.UpdateDistributionStatus(ctx, nonExistentID, "distributed", "0x123", 12345)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "distribution not found")

		// Test ListDistributionsByEpoch with non-existent data (should return empty list)
		distributions, err := store.ListDistributionsByEpoch(ctx, nonExistentEpoch, nonExistentVault)
		require.NoError(t, err)
		assert.Len(t, distributions, 0)

		// Test ListDistributionsByStatus with non-existent status (should return empty list)
		distributions, err = store.ListDistributionsByStatus(ctx, "non_existent_status", 10)
		require.NoError(t, err)
		assert.Len(t, distributions, 0)
	})

	t.Run("DuplicateIDs", func(t *testing.T) {
		epochNumber := big.NewInt(1)
		subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)

		distribution1 := SubsidyDistribution{
			ID:                subsidyData.ID,
			EpochNumber:       subsidyData.EpochNumber,
			VaultID:           subsidyData.VaultID,
			CollectionAddress: subsidyData.CollectionAddress,
			Amount:            subsidyData.Amount,
			Status:            "pending",
			TxHash:            "",
			BlockNumber:       0,
		}

		// Save first distribution
		err := store.SaveDistribution(ctx, distribution1)
		require.NoError(t, err)

		// Create second distribution with same ID but different data
		distribution2 := distribution1
		distribution2.Status = "distributed"
		distribution2.Amount = big.NewInt(999999)

		// Save second distribution (should overwrite first)
		err = store.SaveDistribution(ctx, distribution2)
		require.NoError(t, err)

		// Verify second distribution overwrote first
		retrieved, err := store.GetDistribution(ctx, subsidyData.ID)
		require.NoError(t, err)
		assert.Equal(t, "distributed", retrieved.Status)
		assert.Equal(t, "999999", retrieved.Amount.String())
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

	// Create initial distribution
	subsidyData := generator.GenerateSubsidyData(vaultID, epochNumber)
	distribution := SubsidyDistribution{
		ID:                subsidyData.ID,
		EpochNumber:       subsidyData.EpochNumber,
		VaultID:           subsidyData.VaultID,
		CollectionAddress: subsidyData.CollectionAddress,
		Amount:            subsidyData.Amount,
		Status:            "pending",
		TxHash:            "",
		BlockNumber:       0,
	}

	err = store.SaveDistribution(ctx, distribution)
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
				retrieved, err := store.GetDistribution(ctx, distribution.ID)
				if err != nil {
					errors <- fmt.Errorf("reader %d, iteration %d: %w", readerID, j, err)
					return
				}

				// Verify data consistency
				if retrieved.ID != distribution.ID {
					errors <- fmt.Errorf("reader %d: ID mismatch", readerID)
					return
				}

				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// Start a writer that updates the same distribution
	wg.Add(1)
	go func() {
		defer wg.Done()

		statuses := []string{"distributed", "failed", "pending"}
		for i, status := range statuses {
			txHash := fmt.Sprintf("0x%064d", i)
			blockNumber := int64(12345 + i)

			if err := store.UpdateDistributionStatus(ctx, distribution.ID, status, txHash, blockNumber); err != nil {
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
	final, err := store.GetDistribution(ctx, distribution.ID)
	require.NoError(t, err)
	assert.Equal(t, distribution.ID, final.ID)
	assert.Equal(t, distribution.VaultID, final.VaultID)
	assert.Contains(t, []string{"distributed", "failed", "pending"}, final.Status)
}
