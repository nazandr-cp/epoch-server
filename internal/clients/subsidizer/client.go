package subsidizer

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pkgz/lgr"
)

type Client struct {
	logger lgr.L
}

func NewClient(logger lgr.L) *Client {
	return &Client{
		logger: logger,
	}
}

func (c *Client) UpdateMerkleRoot(ctx context.Context, vaultId string, root [32]byte) error {
	c.logger.Logf("INFO updating merkle root for vault %s: %x", vaultId, root)
	return nil
}

func (c *Client) UpdateMerkleRootAndWaitForConfirmation(ctx context.Context, vaultId string, root [32]byte) error {
	c.logger.Logf("INFO updating merkle root for vault %s: %x", vaultId, root)

	// Simulate transaction submission
	c.logger.Logf("INFO submitting UpdateMerkleRoot transaction for vault %s", vaultId)

	// Simulate waiting for mining confirmation
	c.logger.Logf("INFO waiting for transaction confirmation for vault %s", vaultId)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while waiting for confirmation: %w", ctx.Err())
	case <-time.After(100 * time.Millisecond): // Simulate mining time
		c.logger.Logf("INFO transaction confirmed for vault %s", vaultId)
		return nil
	}
}
