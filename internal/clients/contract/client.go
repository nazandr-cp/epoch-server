package contract

import (
	"context"

	"github.com/go-pkgz/lgr"
)

type EthereumConfig struct {
	RPCURL     string
	PrivateKey string
	GasLimit   uint64
	GasPrice   string
}

type ContractAddresses struct {
	Comptroller        string
	EpochManager       string
	DebtSubsidizer     string
	LendingManager     string
	CollectionRegistry string
}

type Client struct {
	logger    lgr.L
	ethConfig EthereumConfig
	contracts ContractAddresses
}

func NewClient(logger lgr.L) *Client {
	return &Client{
		logger: logger,
	}
}

func NewClientWithConfig(logger lgr.L, ethConfig EthereumConfig, contracts ContractAddresses) *Client {
	return &Client{
		logger:    logger,
		ethConfig: ethConfig,
		contracts: contracts,
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
