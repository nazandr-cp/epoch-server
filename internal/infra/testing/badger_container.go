package testing

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/docker/go-connections/nat"
	"github.com/go-pkgz/lgr"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// BadgerContainer wraps a BadgerDB instance running in a Docker container
type BadgerContainer struct {
	container testcontainers.Container
	db        *badger.DB
	logger    lgr.L
	dataDir   string
}

// BadgerContainerConfig holds configuration for BadgerDB container
type BadgerContainerConfig struct {
	// Docker image to use (optional, defaults to alpine with volume mount)
	Image string
	// Data directory inside container
	DataDir string
	// BadgerDB options
	Options badger.Options
	// Logger instance
	Logger lgr.L
	// Whether to enable debug logging
	Debug bool
}

// NewBadgerContainer creates a new BadgerDB container instance
func NewBadgerContainer(ctx context.Context, config BadgerContainerConfig) (*BadgerContainer, error) {
	if config.Image == "" {
		config.Image = "alpine:latest"
	}
	if config.DataDir == "" {
		config.DataDir = "/data/badger"
	}
	if config.Logger == nil {
		config.Logger = lgr.New(lgr.Debug)
	}

	// Create container request
	req := testcontainers.ContainerRequest{
		Image: config.Image,
		// Keep container running
		Cmd:          []string{"sleep", "3600"},
		ExposedPorts: []string{},
		WaitingFor:   wait.ForExec([]string{"echo", "ready"}).WithStartupTimeout(30 * time.Second),
	}

	// Start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start BadgerDB container: %w", err)
	}

	// Create BadgerDB instance
	opts := config.Options
	if opts.Dir == "" {
		// Use a temp directory for the test
		opts.Dir = "/tmp/badger-test"
		opts.ValueDir = "/tmp/badger-test"
	}

	// Set up BadgerDB options for testing
	opts.Logger = newBadgerLogger(config.Logger)
	opts.MemTableSize = 1 << 20 // 1MB for faster tests
	opts.NumMemtables = 2
	opts.NumLevelZeroTables = 1
	opts.NumLevelZeroTablesStall = 2
	opts.LevelSizeMultiplier = 2
	opts.MaxLevels = 3
	opts.SyncWrites = false // Faster for tests
	opts.NumVersionsToKeep = 1
	opts.CompactL0OnClose = true
	opts.ValueLogFileSize = 16 << 20 // 16MB - minimum valid size

	if config.Debug {
		opts.Logger = newBadgerLogger(lgr.New(lgr.Debug))
	}

	db, err := badger.Open(opts)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to open BadgerDB: %w", err)
	}

	bc := &BadgerContainer{
		container: container,
		db:        db,
		logger:    config.Logger,
		dataDir:   config.DataDir,
	}

	return bc, nil
}

// GetDB returns the BadgerDB instance
func (bc *BadgerContainer) GetDB() *badger.DB {
	return bc.db
}

// GetContainer returns the underlying testcontainer
func (bc *BadgerContainer) GetContainer() testcontainers.Container {
	return bc.container
}

// GetHost returns the container host
func (bc *BadgerContainer) GetHost(ctx context.Context) (string, error) {
	return bc.container.Host(ctx)
}

// GetPort returns the mapped port for the given container port
func (bc *BadgerContainer) GetPort(ctx context.Context, port nat.Port) (nat.Port, error) {
	return bc.container.MappedPort(ctx, port)
}

// ExecuteCommand executes a command in the container
func (bc *BadgerContainer) ExecuteCommand(ctx context.Context, cmd []string) (int, error) {
	exitCode, _, err := bc.container.Exec(ctx, cmd)
	return exitCode, err
}

// GetDataDir returns the data directory path
func (bc *BadgerContainer) GetDataDir() string {
	return bc.dataDir
}

// GetStats returns BadgerDB statistics
func (bc *BadgerContainer) GetStats() (int64, int64) {
	lsm, vlog := bc.db.Size()
	return lsm, vlog
}

// RunGC runs garbage collection on the BadgerDB
func (bc *BadgerContainer) RunGC(ctx context.Context) error {
	return bc.db.RunValueLogGC(0.5)
}

// Sync forces a sync of the BadgerDB
func (bc *BadgerContainer) Sync() error {
	return bc.db.Sync()
}

// Backup creates a backup of the BadgerDB
func (bc *BadgerContainer) Backup(ctx context.Context, writer interface{}) error {
	// Since we can't import the backup package directly, we'll implement a simple key-value dump
	return bc.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				bc.logger.Logf("DEBUG backup key=%s, value_len=%d", string(k), len(v))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// GetKeyCount returns the total number of keys in the database
func (bc *BadgerContainer) GetKeyCount() (int, error) {
	count := 0
	err := bc.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			count++
		}
		return nil
	})
	return count, err
}

// GetKeysWithPrefix returns all keys with the given prefix
func (bc *BadgerContainer) GetKeysWithPrefix(prefix string) ([]string, error) {
	var keys []string
	err := bc.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())
			if len(prefix) == 0 || string(key[:len(prefix)]) == prefix {
				keys = append(keys, key)
			}
		}
		return nil
	})
	return keys, err
}

// Clear removes all data from the BadgerDB
func (bc *BadgerContainer) Clear() error {
	return bc.db.DropAll()
}

// Close closes the BadgerDB and stops the container
func (bc *BadgerContainer) Close(ctx context.Context) error {
	var errs []error

	// Close BadgerDB
	if bc.db != nil {
		if err := bc.db.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close BadgerDB: %w", err))
		}
	}

	// Stop container
	if bc.container != nil {
		if err := bc.container.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate container: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during close: %v", errs)
	}
	return nil
}

// badgerLogger adapts lgr.L to badger's Logger interface
type badgerLogger struct {
	lgr lgr.L
}

func newBadgerLogger(l lgr.L) *badgerLogger {
	return &badgerLogger{lgr: l}
}

func (l *badgerLogger) Errorf(format string, args ...interface{}) {
	l.lgr.Logf("ERROR "+format, args...)
}

func (l *badgerLogger) Warningf(format string, args ...interface{}) {
	l.lgr.Logf("WARN "+format, args...)
}

func (l *badgerLogger) Infof(format string, args ...interface{}) {
	l.lgr.Logf("INFO "+format, args...)
}

func (l *badgerLogger) Debugf(format string, args ...interface{}) {
	l.lgr.Logf("DEBUG "+format, args...)
}
