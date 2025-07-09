package testing

import (
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-pkgz/lgr"
)

// TestConfig holds configuration for integration tests
type TestConfig struct {
	// BadgerDB configuration
	BadgerDB BadgerDBTestConfig `yaml:"badger_db"`
	
	// Test execution configuration
	Execution ExecutionConfig `yaml:"execution"`
	
	// Performance test configuration
	Performance PerformanceConfig `yaml:"performance"`
	
	// Concurrency test configuration
	Concurrency ConcurrencyConfig `yaml:"concurrency"`
}

// BadgerDBTestConfig contains BadgerDB-specific test configuration
type BadgerDBTestConfig struct {
	// Base directory for test databases
	BaseDir string `yaml:"base_dir"`
	
	// Whether to enable debug logging
	Debug bool `yaml:"debug"`
	
	// Memory table size for tests (in bytes)
	MemTableSize int64 `yaml:"mem_table_size"`
	
	// Number of memory tables
	NumMemtables int `yaml:"num_memtables"`
	
	
	// Whether to sync writes (slower but safer)
	SyncWrites bool `yaml:"sync_writes"`
	
	// Number of versions to keep
	NumVersionsToKeep int `yaml:"num_versions_to_keep"`
	
	// Compact L0 on close
	CompactL0OnClose bool `yaml:"compact_l0_on_close"`
	
	// Maximum levels in LSM tree
	MaxLevels int `yaml:"max_levels"`
	
	// Level size multiplier
	LevelSizeMultiplier int `yaml:"level_size_multiplier"`
}

// ExecutionConfig contains test execution parameters
type ExecutionConfig struct {
	// Timeout for individual tests
	TestTimeout time.Duration `yaml:"test_timeout"`
	
	// Timeout for container startup
	StartupTimeout time.Duration `yaml:"startup_timeout"`
	
	// Whether to run tests in parallel
	ParallelExecution bool `yaml:"parallel_execution"`
	
	// Whether to clean up containers after tests
	CleanupContainers bool `yaml:"cleanup_containers"`
	
	// Whether to preserve test data for debugging
	PreserveData bool `yaml:"preserve_data"`
}

// PerformanceConfig contains performance test parameters
type PerformanceConfig struct {
	// Number of operations for performance tests
	OperationCount int `yaml:"operation_count"`
	
	// Size of test data (in bytes)
	DataSize int `yaml:"data_size"`
	
	// Number of concurrent operations
	ConcurrentOperations int `yaml:"concurrent_operations"`
	
	// Duration for stress tests
	StressDuration time.Duration `yaml:"stress_duration"`
	
	// Memory limit for tests (in bytes)
	MemoryLimit int64 `yaml:"memory_limit"`
	
	// Whether to enable memory profiling
	MemoryProfiling bool `yaml:"memory_profiling"`
	
	// Whether to enable CPU profiling
	CPUProfiling bool `yaml:"cpu_profiling"`
}

// ConcurrencyConfig contains concurrency test parameters
type ConcurrencyConfig struct {
	// Number of concurrent goroutines
	NumGoroutines int `yaml:"num_goroutines"`
	
	// Number of concurrent readers
	NumReaders int `yaml:"num_readers"`
	
	// Number of concurrent writers
	NumWriters int `yaml:"num_writers"`
	
	// Duration for concurrency tests
	Duration time.Duration `yaml:"duration"`
	
	// Whether to test transaction conflicts
	TestTransactionConflicts bool `yaml:"test_transaction_conflicts"`
	
	// Whether to test deadlock scenarios
	TestDeadlocks bool `yaml:"test_deadlocks"`
}

// DefaultTestConfig returns a default test configuration
func DefaultTestConfig() TestConfig {
	return TestConfig{
		BadgerDB: BadgerDBTestConfig{
			BaseDir:               "/tmp/badger-test",
			Debug:                 false,
			MemTableSize:          1 << 20, // 1MB
			NumMemtables:          2,
			SyncWrites:            false,
			NumVersionsToKeep:     1,
			CompactL0OnClose:      true,
			MaxLevels:             3,
			LevelSizeMultiplier:   2,
		},
		Execution: ExecutionConfig{
			TestTimeout:       5 * time.Minute,
			StartupTimeout:    30 * time.Second,
			ParallelExecution: true,
			CleanupContainers: true,
			PreserveData:      false,
		},
		Performance: PerformanceConfig{
			OperationCount:       10000,
			DataSize:             1024,
			ConcurrentOperations: 10,
			StressDuration:       30 * time.Second,
			MemoryLimit:          100 << 20, // 100MB
			MemoryProfiling:      false,
			CPUProfiling:         false,
		},
		Concurrency: ConcurrencyConfig{
			NumGoroutines:            10,
			NumReaders:               5,
			NumWriters:               3,
			Duration:                 10 * time.Second,
			TestTransactionConflicts: true,
			TestDeadlocks:            true,
		},
	}
}

// PerformanceTestConfig returns a configuration optimized for performance tests
func PerformanceTestConfig() TestConfig {
	config := DefaultTestConfig()
	config.Performance.OperationCount = 100000
	config.Performance.ConcurrentOperations = 50
	config.Performance.StressDuration = 2 * time.Minute
	config.Performance.MemoryProfiling = true
	config.Performance.CPUProfiling = true
	config.Execution.TestTimeout = 10 * time.Minute
	return config
}

// ConcurrencyTestConfig returns a configuration optimized for concurrency tests
func ConcurrencyTestConfig() TestConfig {
	config := DefaultTestConfig()
	config.Concurrency.NumGoroutines = 50
	config.Concurrency.NumReaders = 20
	config.Concurrency.NumWriters = 10
	config.Concurrency.Duration = 30 * time.Second
	config.Execution.TestTimeout = 5 * time.Minute
	return config
}

// ToBadgerOptions converts BadgerDBTestConfig to badger.Options
func (c BadgerDBTestConfig) ToBadgerOptions(dir string, logger lgr.L) badger.Options {
	opts := badger.DefaultOptions(dir)
	opts.Logger = newBadgerLogger(logger)
	opts.MemTableSize = c.MemTableSize
	opts.NumMemtables = c.NumMemtables
	opts.SyncWrites = c.SyncWrites
	opts.NumVersionsToKeep = c.NumVersionsToKeep
	opts.CompactL0OnClose = c.CompactL0OnClose
	opts.MaxLevels = c.MaxLevels
	opts.LevelSizeMultiplier = c.LevelSizeMultiplier
	
	// Test-specific optimizations
	opts.NumLevelZeroTables = 1
	opts.NumLevelZeroTablesStall = 2
	opts.ValueThreshold = 32 // Store smaller values in LSM tree
	opts.ValueLogFileSize = 16 << 20 // 16MB - minimum valid size
	
	return opts
}