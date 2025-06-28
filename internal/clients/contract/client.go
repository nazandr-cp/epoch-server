package contract

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

func (c *Client) StartEpoch(ctx context.Context, epochID string) error {
	c.logger.Logf("INFO starting epoch %s", epochID)
	return nil
}

func (c *Client) DistributeSubsidies(ctx context.Context, epochID string) error {
	c.logger.Logf("INFO distributing subsidies for epoch %s", epochID)
	return nil
}
