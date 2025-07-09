package merkleimpl

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

// MerkleEntry represents a leaf entry in the Merkle tree
type MerkleEntry struct {
	Address     string   `json:"address"`
	TotalEarned *big.Int `json:"totalEarned"`
}

// MerkleSnapshot represents a complete snapshot of merkle tree data for an epoch
type MerkleSnapshot struct {
	EpochNumber *big.Int      `json:"epochNumber"`
	Entries     []MerkleEntry `json:"entries"`
	MerkleRoot  string        `json:"merkleRoot"`
	Timestamp   int64         `json:"timestamp"`
	VaultID     string        `json:"vaultId"`
	BlockNumber int64         `json:"blockNumber"`
	CreatedAt   time.Time     `json:"createdAt"`
}

// Store handles storage operations for merkle service
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

// SaveSnapshot saves a merkle snapshot for an epoch
func (s *Store) SaveSnapshot(ctx context.Context, epochNumber *big.Int, snapshot MerkleSnapshot) error {
	snapshot.EpochNumber = epochNumber
	snapshot.CreatedAt = time.Now()

	key := s.buildSnapshotKey(epochNumber, snapshot.VaultID)
	data, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
	if err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	// Update latest snapshot pointer
	latestKey := s.buildLatestKey(snapshot.VaultID)
	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(latestKey), []byte(epochNumber.String()))
	})
	if err != nil {
		s.logger.Logf("WARN failed to update latest snapshot pointer: %v", err)
	}

	s.logger.Logf("INFO saved merkle snapshot for vault %s, epoch %s with %d entries",
		snapshot.VaultID, epochNumber.String(), len(snapshot.Entries))
	return nil
}

// GetSnapshot retrieves a merkle snapshot for a specific epoch
func (s *Store) GetSnapshot(ctx context.Context, epochNumber *big.Int, vaultID string) (*MerkleSnapshot, error) {
	key := s.buildSnapshotKey(epochNumber, vaultID)
	
	var snapshot MerkleSnapshot
	err := s.db.View(func(txn *badger.Txn) error {
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
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	return &snapshot, nil
}

// GetLatestSnapshot retrieves the latest merkle snapshot for a vault
func (s *Store) GetLatestSnapshot(ctx context.Context, vaultID string) (*MerkleSnapshot, error) {
	latestKey := s.buildLatestKey(vaultID)
	
	var latestEpochStr string
	err := s.db.View(func(txn *badger.Txn) error {
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

	return s.GetSnapshot(ctx, latestEpoch, vaultID)
}

// ListSnapshots retrieves multiple snapshots for a vault
func (s *Store) ListSnapshots(ctx context.Context, vaultID string, limit int) ([]MerkleSnapshot, error) {
	prefix := s.buildVaultPrefix(vaultID)
	var snapshots []MerkleSnapshot
	
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
			
			// Skip non-snapshot keys (like latest pointer)
			if !strings.Contains(keyStr, ":epoch:") {
				continue
			}
			
			err := item.Value(func(val []byte) error {
				var snapshot MerkleSnapshot
				if err := json.Unmarshal(val, &snapshot); err != nil {
					s.logger.Logf("WARN failed to unmarshal snapshot: %v", err)
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

// Key building functions
func (s *Store) buildSnapshotKey(epochNumber *big.Int, vaultID string) string {
	normalizedVaultID := strings.ToLower(vaultID)
	return fmt.Sprintf("merkle:snapshot:vault:%s:epoch:%020s", normalizedVaultID, epochNumber.String())
}

func (s *Store) buildLatestKey(vaultID string) string {
	normalizedVaultID := strings.ToLower(vaultID)
	return fmt.Sprintf("merkle:latest:vault:%s", normalizedVaultID)
}

func (s *Store) buildVaultPrefix(vaultID string) string {
	normalizedVaultID := strings.ToLower(vaultID)
	return fmt.Sprintf("merkle:snapshot:vault:%s:", normalizedVaultID)
}