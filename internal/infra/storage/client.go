package storage

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-pkgz/lgr"
)

// StorageConfig contains configuration for storage
type StorageConfig struct {
	Type string `yaml:"type"` // "badger" or "memory"
	Path string `yaml:"path"` // Path for badger database
}

// DatabaseWrapper provides a database connection
type DatabaseWrapper struct {
	db     *badger.DB
	logger lgr.L
}

// NewDatabaseWrapper creates a new database wrapper
func NewDatabaseWrapper(config StorageConfig, logger lgr.L) (*DatabaseWrapper, error) {
	switch config.Type {
	case "badger":
		opts := badger.DefaultOptions(config.Path)
		opts.Logger = newBadgerLogger(logger)
		
		db, err := badger.Open(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to open badger database: %w", err)
		}
		
		return &DatabaseWrapper{
			db:     db,
			logger: logger,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}
}

// GetDB returns the database instance
func (w *DatabaseWrapper) GetDB() *badger.DB {
	return w.db
}

// Close closes the database connection
func (w *DatabaseWrapper) Close() error {
	if w.db != nil {
		return w.db.Close()
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
