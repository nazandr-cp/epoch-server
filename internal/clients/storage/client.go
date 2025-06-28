package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/go-pkgz/lgr"
)

type MerkleEntry struct {
	Address     string   `json:"address"`
	TotalEarned *big.Int `json:"totalEarned"`
}

type MerkleSnapshot struct {
	Entries    []MerkleEntry `json:"entries"`
	MerkleRoot string        `json:"merkleRoot"`
	Timestamp  int64         `json:"timestamp"`
	VaultID    string        `json:"vaultId"`
}

type Client struct {
	logger lgr.L
}

func NewClient(logger lgr.L) *Client {
	return &Client{
		logger: logger,
	}
}

func (c *Client) SaveSnapshot(ctx context.Context, snapshot MerkleSnapshot) error {
	data, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	c.logger.Logf("INFO saving snapshot for vault %s with %d entries: %s",
		snapshot.VaultID, len(snapshot.Entries), string(data))
	return nil
}
