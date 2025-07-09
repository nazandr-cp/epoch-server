package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-pkgz/lgr"
)

type MerkleEntry struct {
	Address     string   `json:"address"`
	TotalEarned *big.Int `json:"totalEarned"`
}

type MerkleSnapshot struct {
	EpochNumber *big.Int      `json:"epochNumber"`
	Entries     []MerkleEntry `json:"entries"`
	MerkleRoot  string        `json:"merkleRoot"`
	Timestamp   int64         `json:"timestamp"`
	VaultID     string        `json:"vaultId"`
	BlockNumber int64         `json:"blockNumber"` // Block number used for data consistency
	CreatedAt   time.Time     `json:"createdAt"`
}

// EpochStorageClient defines the interface for epoch-based snapshot storage
type EpochStorageClient interface {
	SaveEpochSnapshot(ctx context.Context, epochNumber *big.Int, snapshot MerkleSnapshot) error
	GetEpochSnapshot(ctx context.Context, epochNumber *big.Int, vaultID string) (*MerkleSnapshot, error)
	GetLatestEpochSnapshot(ctx context.Context, vaultID string) (*MerkleSnapshot, error)
	ListEpochSnapshots(ctx context.Context, vaultID string, limit int) ([]MerkleSnapshot, error)
	// Legacy method for backward compatibility
	SaveSnapshot(ctx context.Context, snapshot MerkleSnapshot) error
}

type Client struct {
	logger    lgr.L
	mu        sync.RWMutex
	baseDir   string
	latestMap map[string]*MerkleSnapshot // vaultID -> latest snapshot
}

func NewClient(logger lgr.L) *Client {
	return NewClientWithBaseDir(logger, "snapshots")
}

func NewClientWithBaseDir(logger lgr.L, baseDir string) *Client {
	return &Client{
		logger:    logger,
		baseDir:   baseDir,
		latestMap: make(map[string]*MerkleSnapshot),
	}
}

func (c *Client) SaveEpochSnapshot(ctx context.Context, epochNumber *big.Int, snapshot MerkleSnapshot) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ensure epoch number is set
	snapshot.EpochNumber = epochNumber
	snapshot.CreatedAt = time.Now()

	// Create directory structure
	epochDir := filepath.Join(c.baseDir, fmt.Sprintf("epoch_%s", epochNumber.String()))
	if err := os.MkdirAll(epochDir, 0755); err != nil {
		return fmt.Errorf("failed to create epoch directory: %w", err)
	}

	// Save to file (normalize vault ID for consistent filename)
	normalizedVaultID := strings.ToLower(snapshot.VaultID)
	filePath := filepath.Join(epochDir, fmt.Sprintf("vault_%s.json", normalizedVaultID))
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	// Update latest snapshot cache
	c.latestMap[snapshot.VaultID] = &snapshot

	c.logger.Logf("INFO saved epoch snapshot for vault %s, epoch %s with %d entries to %s",
		snapshot.VaultID, epochNumber.String(), len(snapshot.Entries), filePath)
	return nil
}

func (c *Client) GetEpochSnapshot(ctx context.Context, epochNumber *big.Int, vaultID string) (*MerkleSnapshot, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	filePath := filepath.Join(c.baseDir, fmt.Sprintf("epoch_%s", epochNumber.String()), fmt.Sprintf("vault_%s.json", strings.ToLower(vaultID)))
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot not found for vault %s, epoch %s", vaultID, epochNumber.String())
		}
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}

	var snapshot MerkleSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	return &snapshot, nil
}

func (c *Client) GetLatestEpochSnapshot(ctx context.Context, vaultID string) (*MerkleSnapshot, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Check cache first
	if latest, exists := c.latestMap[vaultID]; exists {
		return latest, nil
	}

	// Fallback to file system scan
	snapshots, err := c.ListEpochSnapshots(ctx, vaultID, 1)
	if err != nil {
		return nil, err
	}

	if len(snapshots) == 0 {
		return nil, fmt.Errorf("no snapshots found for vault %s", vaultID)
	}

	return &snapshots[0], nil
}

func (c *Client) ListEpochSnapshots(ctx context.Context, vaultID string, limit int) ([]MerkleSnapshot, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Read all epoch directories
	entries, err := os.ReadDir(c.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []MerkleSnapshot{}, nil
		}
		return nil, fmt.Errorf("failed to read snapshots directory: %w", err)
	}

	var snapshots []MerkleSnapshot
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Parse epoch number from directory name
		if len(entry.Name()) < 7 || entry.Name()[:6] != "epoch_" {
			continue
		}

		epochStr := entry.Name()[6:]
		if _, ok := new(big.Int).SetString(epochStr, 10); !ok {
			continue
		}

		// Check if vault file exists in this epoch
		vaultFile := filepath.Join(c.baseDir, entry.Name(), fmt.Sprintf("vault_%s.json", strings.ToLower(vaultID)))
		if _, err := os.Stat(vaultFile); os.IsNotExist(err) {
			continue
		}

		// Read snapshot
		data, err := os.ReadFile(vaultFile)
		if err != nil {
			c.logger.Logf("WARN failed to read snapshot file %s: %v", vaultFile, err)
			continue
		}

		var snapshot MerkleSnapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			c.logger.Logf("WARN failed to unmarshal snapshot from %s: %v", vaultFile, err)
			continue
		}

		snapshots = append(snapshots, snapshot)
	}

	// Sort by epoch number (descending)
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].EpochNumber.Cmp(snapshots[j].EpochNumber) > 0
	})

	// Apply limit
	if limit > 0 && len(snapshots) > limit {
		snapshots = snapshots[:limit]
	}

	return snapshots, nil
}

// Legacy method for backward compatibility
func (c *Client) SaveSnapshot(ctx context.Context, snapshot MerkleSnapshot) error {
	// For backward compatibility, we'll use epoch number from the snapshot if available
	if snapshot.EpochNumber != nil {
		return c.SaveEpochSnapshot(ctx, snapshot.EpochNumber, snapshot)
	}

	// If no epoch number, just log (original behavior)
	data, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	c.logger.Logf("INFO saving snapshot for vault %s with %d entries: %s",
		snapshot.VaultID, len(snapshot.Entries), string(data))
	return nil
}
