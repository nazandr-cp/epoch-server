package storage

import (
	"fmt"

	"github.com/andrey/epoch-server/internal/infra/storage"
	"github.com/dgraph-io/badger/v4"
	"github.com/go-pkgz/lgr"
)

type Client struct {
	db     *badger.DB
	logger lgr.L
}

// ProvideClient creates a storage client implementation
func ProvideClient(config storage.Config, logger lgr.L) (storage.StorageClient, error) {
	switch config.Type {
	case "badger":
		opts := badger.DefaultOptions(config.Path)
		opts.Logger = newBadgerLogger(logger)

		db, err := badger.Open(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to open badger database: %w", err)
		}

		return &Client{
			db:     db,
			logger: logger,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}
}

func (c *Client) GetDB() *badger.DB {
	return c.db
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

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

