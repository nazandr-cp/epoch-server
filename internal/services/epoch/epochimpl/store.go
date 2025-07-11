package epochimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/andrey/epoch-server/internal/infra/utils"
	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/dgraph-io/badger/v4"
	"github.com/go-pkgz/lgr"
)

// Store handles storage operations for epoch service
type Store struct {
	db     *badger.DB
	logger lgr.L
}

// NewStore creates a new store instance
func NewStore(db *badger.DB, logger lgr.L) *Store {
	return &Store{
		db:     db,
		logger: logger,
	}
}

// SaveEpoch saves epoch information
func (s *Store) SaveEpoch(ctx context.Context, epoch epoch.EpochInfo) error {
	epoch.UpdatedAt = time.Now()
	if epoch.CreatedAt.IsZero() {
		epoch.CreatedAt = time.Now()
	}

	key := s.buildEpochKey(epoch.Number, epoch.VaultID)
	data, err := json.Marshal(epoch)
	if err != nil {
		return fmt.Errorf("failed to marshal epoch: %w", err)
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
	if err != nil {
		return fmt.Errorf("failed to save epoch: %w", err)
	}

	// Update current epoch pointer if this is the latest
	currentKey := s.buildCurrentKey(epoch.VaultID)
	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(currentKey), []byte(epoch.Number.String()))
	})
	if err != nil {
		s.logger.Logf("WARN failed to update current epoch pointer: %v", err)
	}

	s.logger.Logf("INFO saved epoch %s for vault %s with status %s",
		epoch.Number.String(), epoch.VaultID, epoch.Status)
	return nil
}

// GetEpoch retrieves epoch information
func (s *Store) GetEpoch(ctx context.Context, epochNumber *big.Int, vaultID string) (*epoch.EpochInfo, error) {
	key := s.buildEpochKey(epochNumber, vaultID)

	var epochInfo epoch.EpochInfo
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &epochInfo)
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, fmt.Errorf("epoch not found for vault %s, epoch %s", vaultID, epochNumber.String())
		}
		return nil, fmt.Errorf("failed to get epoch: %w", err)
	}

	return &epochInfo, nil
}

// GetCurrentEpoch retrieves the current epoch for a vault
func (s *Store) GetCurrentEpoch(ctx context.Context, vaultID string) (*epoch.EpochInfo, error) {
	currentKey := s.buildCurrentKey(vaultID)

	var currentEpochStr string
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(currentKey))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			currentEpochStr = string(val)
			return nil
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, fmt.Errorf("no current epoch found for vault %s", vaultID)
		}
		return nil, fmt.Errorf("failed to get current epoch pointer: %w", err)
	}

	currentEpoch, ok := new(big.Int).SetString(currentEpochStr, 10)
	if !ok {
		return nil, fmt.Errorf("invalid current epoch number: %s", currentEpochStr)
	}

	return s.GetEpoch(ctx, currentEpoch, vaultID)
}

// ListEpochs retrieves multiple epochs for a vault
func (s *Store) ListEpochs(ctx context.Context, vaultID string, limit int) ([]epoch.EpochInfo, error) {
	prefix := s.buildVaultPrefix(vaultID)
	var epochs []epoch.EpochInfo

	err := s.db.View(func(txn *badger.Txn) error {
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

			// Skip non-epoch keys (like current pointer)
			if !strings.Contains(keyStr, ":epoch:") {
				continue
			}

			err := item.Value(func(val []byte) error {
				var epochInfo epoch.EpochInfo
				if err := json.Unmarshal(val, &epochInfo); err != nil {
					s.logger.Logf("WARN failed to unmarshal epoch: %v", err)
					return nil // Continue iteration
				}
				epochs = append(epochs, epochInfo)
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
		return nil, fmt.Errorf("failed to list epochs: %w", err)
	}

	return epochs, nil
}

// UpdateEpochStatus updates the status of an epoch
func (s *Store) UpdateEpochStatus(ctx context.Context, epochNumber *big.Int, vaultID string, status string) error {
	epoch, err := s.GetEpoch(ctx, epochNumber, vaultID)
	if err != nil {
		return fmt.Errorf("failed to get epoch for status update: %w", err)
	}

	epoch.Status = status
	epoch.UpdatedAt = time.Now()

	return s.SaveEpoch(ctx, *epoch)
}

// Key building functions
func (s *Store) buildEpochKey(epochNumber *big.Int, vaultID string) string {
	normalizedVaultID := utils.NormalizeAddress(vaultID)
	return fmt.Sprintf("epoch:vault:%s:epoch:%020s", normalizedVaultID, epochNumber.String())
}

func (s *Store) buildCurrentKey(vaultID string) string {
	normalizedVaultID := utils.NormalizeAddress(vaultID)
	return fmt.Sprintf("epoch:current:vault:%s", normalizedVaultID)
}

func (s *Store) buildVaultPrefix(vaultID string) string {
	normalizedVaultID := utils.NormalizeAddress(vaultID)
	return fmt.Sprintf("epoch:vault:%s:", normalizedVaultID)
}
