package epoch

import (
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

type EpochInfo struct {
	EndTime int64 `json:"endTime"`
}

func (c *Client) Current() EpochInfo {
	return EpochInfo{
		EndTime: time.Now().Unix(),
	}
}

func (c *Client) FinalizeEpoch() error {
	c.logger.Logf("INFO finalizing epoch")
	return nil
}
