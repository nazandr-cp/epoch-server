package testing

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/andrey/epoch-server/internal/infra/utils"
	"github.com/dgraph-io/badger/v4"
	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BadgerTestHelper provides utilities for testing BadgerDB operations
type BadgerTestHelper struct {
	container *BadgerContainer
	logger    lgr.L
}

// NewBadgerTestHelper creates a new BadgerDB test helper
func NewBadgerTestHelper(container *BadgerContainer, logger lgr.L) *BadgerTestHelper {
	return &BadgerTestHelper{
		container: container,
		logger:    logger,
	}
}

// AssertKeyExists checks if a key exists in the database
func (h *BadgerTestHelper) AssertKeyExists(t require.TestingT, key string) {
	err := h.container.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})
	require.NoError(t, err, "Key %s should exist", key)
}

// AssertKeyNotExists checks if a key does not exist in the database
func (h *BadgerTestHelper) AssertKeyNotExists(t require.TestingT, key string) {
	err := h.container.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})
	require.Error(t, err, "Key %s should not exist", key)
	require.Equal(t, badger.ErrKeyNotFound, err)
}

// AssertKeyValue checks if a key has the expected value
func (h *BadgerTestHelper) AssertKeyValue(t require.TestingT, key string, expectedValue []byte) {
	err := h.container.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			require.Equal(t, expectedValue, val, "Value mismatch for key %s", key)
			return nil
		})
	})
	require.NoError(t, err)
}

// AssertKeyCount checks if the total number of keys matches expected count
func (h *BadgerTestHelper) AssertKeyCount(t require.TestingT, expectedCount int) {
	count, err := h.container.GetKeyCount()
	require.NoError(t, err)
	require.Equal(t, expectedCount, count, "Key count mismatch")
}

// AssertKeysWithPrefix checks if keys with prefix match expected count
func (h *BadgerTestHelper) AssertKeysWithPrefix(t require.TestingT, prefix string, expectedCount int) {
	keys, err := h.container.GetKeysWithPrefix(prefix)
	require.NoError(t, err)
	require.Equal(t, expectedCount, len(keys), "Key count with prefix %s mismatch", prefix)
}

// AssertDBEmpty checks if the database is empty
func (h *BadgerTestHelper) AssertDBEmpty(t require.TestingT) {
	h.AssertKeyCount(t, 0)
}

// AssertTransactionIsolation tests transaction isolation by running concurrent transactions
func (h *BadgerTestHelper) AssertTransactionIsolation(t require.TestingT, key string, initialValue []byte) {
	// Set initial value
	err := h.container.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), initialValue)
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	errors := make(chan error, 2)

	// Transaction 1: Read and update
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.container.db.Update(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(key))
			if err != nil {
				return err
			}

			var val []byte
			err = item.Value(func(v []byte) error {
				val = append(val, v...)
				return nil
			})
			if err != nil {
				return err
			}

			// Simulate work
			time.Sleep(50 * time.Millisecond)

			// Update value
			newVal := append(val, []byte("-tx1")...)
			return txn.Set([]byte(key), newVal)
		})
		if err != nil {
			errors <- err
		}
	}()

	// Transaction 2: Read and update
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.container.db.Update(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(key))
			if err != nil {
				return err
			}

			var val []byte
			err = item.Value(func(v []byte) error {
				val = append(val, v...)
				return nil
			})
			if err != nil {
				return err
			}

			// Simulate work
			time.Sleep(50 * time.Millisecond)

			// Update value
			newVal := append(val, []byte("-tx2")...)
			return txn.Set([]byte(key), newVal)
		})
		if err != nil {
			errors <- err
		}
	}()

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		require.NoError(t, err, "Transaction isolation test failed")
	}

	// Verify final value
	err = h.container.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			h.logger.Logf("DEBUG final value: %s", string(val))
			// One of the transactions should have won
			assert.True(t, len(val) > len(initialValue), "Final value should be modified")
			return nil
		})
	})
	require.NoError(t, err)
}

// BenchmarkOperations runs benchmark operations on BadgerDB
func (h *BadgerTestHelper) BenchmarkOperations(b *testing.B, operations func(db *badger.DB, b *testing.B)) {
	b.ResetTimer()
	operations(h.container.db, b)
}

// GenerateRandomKey generates a random key for testing
func (h *BadgerTestHelper) GenerateRandomKey(prefix string) string {
	// use crypto/rand for better security
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		// fallback to timestamp-based key if crypto/rand fails
		return fmt.Sprintf("%s:%d", prefix, time.Now().UnixNano())
	}
	// convert bytes to uint32 for consistent formatting
	randomValue := uint32(randomBytes[0])<<24 | uint32(randomBytes[1])<<16 | uint32(randomBytes[2])<<8 | uint32(randomBytes[3])
	return fmt.Sprintf("%s:%d:%d", prefix, time.Now().UnixNano(), randomValue)
}

// GenerateRandomValue generates random value data for testing
func (h *BadgerTestHelper) GenerateRandomValue(size int) []byte {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		// fallback to zero bytes if crypto/rand fails
		return make([]byte, size)
	}
	return data
}

// PopulateWithTestData populates the database with test data
func (h *BadgerTestHelper) PopulateWithTestData(count int, keyPrefix string, valueSize int) error {
	return h.container.db.Update(func(txn *badger.Txn) error {
		for i := 0; i < count; i++ {
			key := fmt.Sprintf("%s:%06d", keyPrefix, i)
			value := h.GenerateRandomValue(valueSize)
			if err := txn.Set([]byte(key), value); err != nil {
				return err
			}
		}
		return nil
	})
}

// MeasureOperationTime measures the time taken for a database operation
func (h *BadgerTestHelper) MeasureOperationTime(operation func() error) (time.Duration, error) {
	start := time.Now()
	err := operation()
	duration := time.Since(start)
	return duration, err
}

// RunConcurrentOperations runs multiple operations concurrently
func (h *BadgerTestHelper) RunConcurrentOperations(operations []func() error) []error {
	var wg sync.WaitGroup
	errors := make([]error, len(operations))

	for i, op := range operations {
		wg.Add(1)
		go func(index int, operation func() error) {
			defer wg.Done()
			errors[index] = operation()
		}(i, op)
	}

	wg.Wait()
	return errors
}

// CollectMetrics collects BadgerDB metrics
func (h *BadgerTestHelper) CollectMetrics() BadgerMetrics {
	lsm, vlog := h.container.GetStats()
	keyCount, _ := h.container.GetKeyCount()

	return BadgerMetrics{
		KeyCount:  keyCount,
		LSMSize:   lsm,
		VLogSize:  vlog,
		TotalSize: lsm + vlog,
		Timestamp: time.Now(),
	}
}

// BadgerMetrics holds BadgerDB metrics
type BadgerMetrics struct {
	KeyCount  int
	LSMSize   int64
	VLogSize  int64
	TotalSize int64
	Timestamp time.Time
}

// TestScenario represents a test scenario
type TestScenario struct {
	Name        string
	Description string
	Setup       func(helper *BadgerTestHelper) error
	Execute     func(helper *BadgerTestHelper) error
	Verify      func(helper *BadgerTestHelper, t require.TestingT) error
	Cleanup     func(helper *BadgerTestHelper) error
}

// RunScenario executes a test scenario
func (h *BadgerTestHelper) RunScenario(t require.TestingT, scenario TestScenario) {
	h.logger.Logf("INFO Running scenario: %s", scenario.Name)

	// Setup
	if scenario.Setup != nil {
		err := scenario.Setup(h)
		require.NoError(t, err, "Setup failed for scenario: %s", scenario.Name)
	}

	// Execute
	if scenario.Execute != nil {
		err := scenario.Execute(h)
		require.NoError(t, err, "Execute failed for scenario: %s", scenario.Name)
	}

	// Verify
	if scenario.Verify != nil {
		err := scenario.Verify(h, t)
		require.NoError(t, err, "Verify failed for scenario: %s", scenario.Name)
	}

	// Cleanup
	if scenario.Cleanup != nil {
		err := scenario.Cleanup(h)
		require.NoError(t, err, "Cleanup failed for scenario: %s", scenario.Name)
	}

	h.logger.Logf("INFO Scenario completed: %s", scenario.Name)
}

// CreateEpochKey creates a test epoch key
func (h *BadgerTestHelper) CreateEpochKey(vaultID string, epochNumber *big.Int) string {
	normalizedVaultID := utils.NormalizeAddress(vaultID)
	return fmt.Sprintf("epoch:vault:%s:epoch:%020s", normalizedVaultID, epochNumber.String())
}

// CreateMerkleKey creates a test merkle key
func (h *BadgerTestHelper) CreateMerkleKey(vaultID string, epochNumber *big.Int) string {
	normalizedVaultID := utils.NormalizeAddress(vaultID)
	return fmt.Sprintf("merkle:snapshot:vault:%s:epoch:%020s", normalizedVaultID, epochNumber.String())
}

// CreateSubsidyKey creates a test subsidy key
func (h *BadgerTestHelper) CreateSubsidyKey(distributionID string) string {
	return fmt.Sprintf("subsidy:distribution:%s", distributionID)
}

// CreateSubsidyEpochKey creates a test subsidy epoch index key
func (h *BadgerTestHelper) CreateSubsidyEpochKey(epochNumber *big.Int, vaultID string, distributionID string) string {
	normalizedVaultID := utils.NormalizeAddress(vaultID)
	return fmt.Sprintf("subsidy:epoch:%020s:vault:%s:distribution:%s", epochNumber.String(), normalizedVaultID, distributionID)
}

// WaitForCondition waits for a condition to be true with timeout
func (h *BadgerTestHelper) WaitForCondition(condition func() bool, timeout time.Duration, interval time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(interval)
	}
	return false
}

// GetHost returns the container host
func (h *BadgerTestHelper) GetHost(ctx context.Context) (string, error) {
	return h.container.GetHost(ctx)
}

// Sync forces a sync of the BadgerDB
func (h *BadgerTestHelper) Sync() error {
	return h.container.Sync()
}

// RunGC runs garbage collection on the BadgerDB
func (h *BadgerTestHelper) RunGC(ctx context.Context) error {
	return h.container.RunGC(ctx)
}

// Clear removes all data from the BadgerDB
func (h *BadgerTestHelper) Clear() error {
	return h.container.Clear()
}
