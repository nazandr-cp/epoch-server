package subsidizer

import (
	"context"

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
