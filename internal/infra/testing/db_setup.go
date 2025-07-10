package testing

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-pkgz/lgr"
)

// SetupTestDB creates a BadgerDB instance running in a testcontainer and returns
// the database instance along with a cleanup function.
//
// This is a simplified interface for integration tests that handles all the
// container setup, configuration, and cleanup automatically.
//
// Usage:
//
//	db, cleanup, err := SetupTestDB(ctx)
//	require.NoError(t, err)
//	defer cleanup()
//
//	// Use db directly...
//	store := NewStore(db, logger)
func SetupTestDB(ctx context.Context) (*badger.DB, func(), error) {
	// Create logger for the test setup
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Use default test configuration optimized for integration tests
	config := DefaultTestConfig()
	config.BadgerDB.Debug = false // Reduce noise in tests unless debugging

	// Create unique directory for this test instance
	uniqueDir := filepath.Join(os.TempDir(), fmt.Sprintf("test-badger-%d", time.Now().UnixNano()))

	// Setup BadgerContainer with testcontainers
	containerConfig := BadgerContainerConfig{
		Options: config.BadgerDB.ToBadgerOptions(uniqueDir, logger),
		Logger:  logger,
		Debug:   false,
	}

	container, err := NewBadgerContainer(ctx, containerConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create BadgerDB container: %w", err)
	}

	// Get the database instance
	db := container.GetDB()
	if db == nil {
		if closeErr := container.Close(ctx); closeErr != nil {
			logger.Logf("WARN failed to close BadgerDB container during cleanup: %v", closeErr)
		}
		return nil, nil, fmt.Errorf("failed to get database instance from container")
	}

	// Create cleanup function that handles all teardown
	cleanup := func() {
		if err := container.Close(ctx); err != nil {
			logger.Logf("WARN failed to close BadgerDB container: %v", err)
		}
	}

	return db, cleanup, nil
}

// SetupTestDBWithConfig creates a BadgerDB instance with custom configuration.
// This provides more control over the database setup while still handling
// container management automatically.
func SetupTestDBWithConfig(ctx context.Context, config TestConfig) (*badger.DB, func(), error) {
	// Create logger for the test setup
	logger := lgr.New(lgr.Msec, lgr.Debug)
	if config.BadgerDB.Debug {
		logger = lgr.New(lgr.Msec, lgr.Debug)
	}

	// Create unique directory for this test instance
	uniqueDir := filepath.Join(os.TempDir(), fmt.Sprintf("test-badger-custom-%d", time.Now().UnixNano()))

	// Setup BadgerContainer with custom configuration
	containerConfig := BadgerContainerConfig{
		Options: config.BadgerDB.ToBadgerOptions(uniqueDir, logger),
		Logger:  logger,
		Debug:   config.BadgerDB.Debug,
	}

	container, err := NewBadgerContainer(ctx, containerConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create BadgerDB container with custom config: %w", err)
	}

	// Get the database instance
	db := container.GetDB()
	if db == nil {
		if closeErr := container.Close(ctx); closeErr != nil {
			logger.Logf("WARN failed to close BadgerDB container during cleanup: %v", closeErr)
		}
		return nil, nil, fmt.Errorf("failed to get database instance from container")
	}

	// Create cleanup function that handles all teardown
	cleanup := func() {
		if err := container.Close(ctx); err != nil {
			logger.Logf("WARN failed to close BadgerDB container: %v", err)
		}
	}

	return db, cleanup, nil
}

// SetupTestDBAndHelper creates both a database instance and a test helper.
// This is useful for tests that need the helper utilities.
func SetupTestDBAndHelper(ctx context.Context) (*badger.DB, *BadgerTestHelper, func(), error) {
	// Create logger for the test setup
	logger := lgr.New(lgr.Msec, lgr.Debug)

	// Use default test configuration optimized for integration tests
	config := DefaultTestConfig()
	config.BadgerDB.Debug = false

	// Create unique directory for this test instance
	uniqueDir := filepath.Join(os.TempDir(), fmt.Sprintf("test-badger-helper-%d", time.Now().UnixNano()))

	// Setup BadgerContainer with testcontainers
	containerConfig := BadgerContainerConfig{
		Options: config.BadgerDB.ToBadgerOptions(uniqueDir, logger),
		Logger:  logger,
		Debug:   false,
	}

	container, err := NewBadgerContainer(ctx, containerConfig)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create BadgerDB container: %w", err)
	}

	// Get the database instance
	db := container.GetDB()
	if db == nil {
		if closeErr := container.Close(ctx); closeErr != nil {
			logger.Logf("WARN failed to close BadgerDB container during cleanup: %v", closeErr)
		}
		return nil, nil, nil, fmt.Errorf("failed to get database instance from container")
	}

	// Create test helper
	helper := NewBadgerTestHelper(container, logger)

	// Create cleanup function that handles all teardown
	cleanup := func() {
		if err := container.Close(ctx); err != nil {
			logger.Logf("WARN failed to close BadgerDB container: %v", err)
		}
	}

	return db, helper, cleanup, nil
}
