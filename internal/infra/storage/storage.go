package storage

import "github.com/dgraph-io/badger/v4"

//go:generate moq -out storage_mocks.go . StorageClient

// StorageClient defines the interface for storage operations
type StorageClient interface {
	GetDB() *badger.DB
	Close() error
}

// Config contains configuration for storage
type Config struct {
	Type string `yaml:"type"` // "badger" or "memory"
	Path string `yaml:"path"` // path for badger database
}
