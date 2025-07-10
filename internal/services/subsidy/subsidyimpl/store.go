package subsidyimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/dgraph-io/badger/v4"
	"github.com/go-pkgz/lgr"
)

// Store handles storage operations for subsidy service
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

// SaveDistribution saves a subsidy distribution record
func (s *Store) SaveDistribution(ctx context.Context, distribution subsidy.SubsidyDistribution) error {
	distribution.UpdatedAt = time.Now()
	if distribution.CreatedAt.IsZero() {
		distribution.CreatedAt = time.Now()
	}

	key := s.buildDistributionKey(distribution.ID)
	data, err := json.Marshal(distribution)
	if err != nil {
		return fmt.Errorf("failed to marshal distribution: %w", err)
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
	if err != nil {
		return fmt.Errorf("failed to save distribution: %w", err)
	}

	// Also save by epoch and vault for easier querying
	epochKey := s.buildEpochDistributionKey(distribution.EpochNumber, distribution.VaultID, distribution.ID)
	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(epochKey), []byte(distribution.ID))
	})
	if err != nil {
		s.logger.Logf("WARN failed to save epoch distribution index: %v", err)
	}

	s.logger.Logf("INFO saved subsidy distribution %s for epoch %s, vault %s, amount %s",
		distribution.ID, distribution.EpochNumber.String(), distribution.VaultID, distribution.Amount.String())
	return nil
}

// GetDistribution retrieves a subsidy distribution by ID
func (s *Store) GetDistribution(ctx context.Context, distributionID string) (*subsidy.SubsidyDistribution, error) {
	key := s.buildDistributionKey(distributionID)

	var distribution subsidy.SubsidyDistribution
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &distribution)
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, fmt.Errorf("distribution not found: %s", distributionID)
		}
		return nil, fmt.Errorf("failed to get distribution: %w", err)
	}

	return &distribution, nil
}

// ListDistributionsByEpoch retrieves all distributions for a specific epoch and vault
func (s *Store) ListDistributionsByEpoch(ctx context.Context, epochNumber *big.Int, vaultID string) ([]subsidy.SubsidyDistribution, error) {
	prefix := s.buildEpochVaultPrefix(epochNumber, vaultID)
	var distributions []subsidy.SubsidyDistribution

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(prefix)

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			err := item.Value(func(val []byte) error {
				distributionID := string(val)
				distribution, err := s.GetDistribution(ctx, distributionID)
				if err != nil {
					s.logger.Logf("WARN failed to get distribution %s: %v", distributionID, err)
					return nil // Continue iteration
				}
				distributions = append(distributions, *distribution)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list distributions: %w", err)
	}

	return distributions, nil
}

// ListDistributionsByStatus retrieves all distributions with a specific status
func (s *Store) ListDistributionsByStatus(ctx context.Context, status string, limit int) ([]subsidy.SubsidyDistribution, error) {
	var distributions []subsidy.SubsidyDistribution

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("subsidy:distribution:")

		it := txn.NewIterator(opts)
		defer it.Close()

		count := 0
		for it.Rewind(); it.Valid() && (limit == 0 || count < limit); it.Next() {
			item := it.Item()

			err := item.Value(func(val []byte) error {
				var distribution subsidy.SubsidyDistribution
				if err := json.Unmarshal(val, &distribution); err != nil {
					s.logger.Logf("WARN failed to unmarshal distribution: %v", err)
					return nil // Continue iteration
				}

				if distribution.Status == status {
					distributions = append(distributions, distribution)
					count++
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list distributions by status: %w", err)
	}

	return distributions, nil
}

// UpdateDistributionStatus updates the status of a distribution
func (s *Store) UpdateDistributionStatus(ctx context.Context, distributionID string, status string, txHash string, blockNumber int64) error {
	distribution, err := s.GetDistribution(ctx, distributionID)
	if err != nil {
		return fmt.Errorf("failed to get distribution for status update: %w", err)
	}

	distribution.Status = status
	// Only update TxHash and BlockNumber if they are provided (not empty/zero)
	if txHash != "" {
		distribution.TxHash = txHash
	}
	if blockNumber > 0 {
		distribution.BlockNumber = blockNumber
	}
	distribution.UpdatedAt = time.Now()

	return s.SaveDistribution(ctx, *distribution)
}

// Key building functions
func (s *Store) buildDistributionKey(distributionID string) string {
	return fmt.Sprintf("subsidy:distribution:%s", distributionID)
}

func (s *Store) buildEpochDistributionKey(epochNumber *big.Int, vaultID string, distributionID string) string {
	normalizedVaultID := strings.ToLower(vaultID)
	return fmt.Sprintf("subsidy:epoch:%020s:vault:%s:distribution:%s", epochNumber.String(), normalizedVaultID, distributionID)
}

func (s *Store) buildEpochVaultPrefix(epochNumber *big.Int, vaultID string) string {
	normalizedVaultID := strings.ToLower(vaultID)
	return fmt.Sprintf("subsidy:epoch:%020s:vault:%s:", epochNumber.String(), normalizedVaultID)
}
