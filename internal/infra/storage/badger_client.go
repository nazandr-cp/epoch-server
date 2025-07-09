package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-pkgz/lgr"
)

type BadgerClient struct {
	db     *badger.DB
	logger lgr.L
}

func NewBadgerClient(logger lgr.L, dbPath string) (*BadgerClient, error) {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = newBadgerLogger(logger)
	
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger database: %w", err)
	}

	return &BadgerClient{
		db:     db,
		logger: logger,
	}, nil
}

func (c *BadgerClient) Close() error {
	return c.db.Close()
}

func (c *BadgerClient) SaveEpochSnapshot(ctx context.Context, epochNumber *big.Int, snapshot MerkleSnapshot) error {
	snapshot.EpochNumber = epochNumber
	snapshot.CreatedAt = time.Now()

	key := c.buildSnapshotKey(epochNumber, snapshot.VaultID)
	data, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	err = c.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
	if err != nil {
		return fmt.Errorf("failed to save snapshot to badger: %w", err)
	}

	// Update latest snapshot pointer
	latestKey := c.buildLatestKey(snapshot.VaultID)
	err = c.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(latestKey), []byte(epochNumber.String()))
	})
	if err != nil {
		c.logger.Logf("WARN failed to update latest snapshot pointer: %v", err)
	}

	c.logger.Logf("INFO saved epoch snapshot for vault %s, epoch %s with %d entries",
		snapshot.VaultID, epochNumber.String(), len(snapshot.Entries))
	return nil
}

func (c *BadgerClient) GetEpochSnapshot(ctx context.Context, epochNumber *big.Int, vaultID string) (*MerkleSnapshot, error) {
	key := c.buildSnapshotKey(epochNumber, vaultID)
	
	var snapshot MerkleSnapshot
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &snapshot)
		})
	})
	
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, fmt.Errorf("snapshot not found for vault %s, epoch %s", vaultID, epochNumber.String())
		}
		return nil, fmt.Errorf("failed to get snapshot from badger: %w", err)
	}

	return &snapshot, nil
}

func (c *BadgerClient) GetLatestEpochSnapshot(ctx context.Context, vaultID string) (*MerkleSnapshot, error) {
	latestKey := c.buildLatestKey(vaultID)
	
	var latestEpochStr string
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(latestKey))
		if err != nil {
			return err
		}
		
		return item.Value(func(val []byte) error {
			latestEpochStr = string(val)
			return nil
		})
	})
	
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, fmt.Errorf("no snapshots found for vault %s", vaultID)
		}
		return nil, fmt.Errorf("failed to get latest snapshot pointer: %w", err)
	}

	latestEpoch, ok := new(big.Int).SetString(latestEpochStr, 10)
	if !ok {
		return nil, fmt.Errorf("invalid latest epoch number: %s", latestEpochStr)
	}

	return c.GetEpochSnapshot(ctx, latestEpoch, vaultID)
}

func (c *BadgerClient) ListEpochSnapshots(ctx context.Context, vaultID string, limit int) ([]MerkleSnapshot, error) {
	prefix := c.buildVaultPrefix(vaultID)
	var snapshots []MerkleSnapshot
	
	err := c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Reverse = true // Get latest first
		
		it := txn.NewIterator(opts)
		defer it.Close()
		
		count := 0
		for it.Rewind(); it.Valid() && (limit == 0 || count < limit); it.Next() {
			item := it.Item()
			key := item.Key()
			keyStr := string(key)
			
			// Skip keys that don't match our vault prefix
			if !strings.HasPrefix(keyStr, prefix) {
				continue
			}
			
			// Skip non-snapshot keys (like latest pointer)
			if !strings.Contains(keyStr, ":epoch:") {
				continue
			}
			err := item.Value(func(val []byte) error {
				var snapshot MerkleSnapshot
				if err := json.Unmarshal(val, &snapshot); err != nil {
					c.logger.Logf("WARN failed to unmarshal snapshot: %v", err)
					return nil // Continue iteration
				}
				snapshots = append(snapshots, snapshot)
				count++
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %w", err)
	}

	return snapshots, nil
}

func (c *BadgerClient) SaveSnapshot(ctx context.Context, snapshot MerkleSnapshot) error {
	if snapshot.EpochNumber != nil {
		return c.SaveEpochSnapshot(ctx, snapshot.EpochNumber, snapshot)
	}

	data, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	c.logger.Logf("INFO saving snapshot for vault %s with %d entries: %s",
		snapshot.VaultID, len(snapshot.Entries), string(data))
	return nil
}

// Key building functions
func (c *BadgerClient) buildSnapshotKey(epochNumber *big.Int, vaultID string) string {
	normalizedVaultID := strings.ToLower(vaultID)
	// Use zero-padded epoch number for proper sorting
	return fmt.Sprintf("snapshot:vault:%s:epoch:%020s", normalizedVaultID, epochNumber.String())
}

func (c *BadgerClient) buildLatestKey(vaultID string) string {
	normalizedVaultID := strings.ToLower(vaultID)
	return fmt.Sprintf("latest:vault:%s", normalizedVaultID)
}

func (c *BadgerClient) buildVaultPrefix(vaultID string) string {
	normalizedVaultID := strings.ToLower(vaultID)
	return fmt.Sprintf("snapshot:vault:%s:", normalizedVaultID)
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