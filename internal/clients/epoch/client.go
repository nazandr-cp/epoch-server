package epoch

import (
	"context"
	"math/big"
	"time"

	"github.com/go-pkgz/lgr"
)

type ContractClient interface {
	GetCurrentEpochId(ctx context.Context) (*big.Int, error)
	UpdateExchangeRate(ctx context.Context, lendingManagerAddress string) error
	AllocateYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string) error
	EndEpochWithSubsidies(ctx context.Context, epochId *big.Int, vaultAddress string, merkleRoot [32]byte, subsidiesDistributed *big.Int) error
}

type Client struct {
	logger         lgr.L
	contractClient ContractClient
}

func NewClient(logger lgr.L) *Client {
	return &Client{
		logger: logger,
	}
}

func NewClientWithContract(logger lgr.L, contractClient ContractClient) *Client {
	return &Client{
		logger:         logger,
		contractClient: contractClient,
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

func (c *Client) GetCurrentEpochId(ctx context.Context) (*big.Int, error) {
	c.logger.Logf("INFO getting current epoch ID")
	
	if c.contractClient == nil {
		c.logger.Logf("WARN contract client not initialized, returning epoch ID 1")
		return big.NewInt(1), nil
	}
	
	return c.contractClient.GetCurrentEpochId(ctx)
}

func (c *Client) FinalizeEpoch() error {
	c.logger.Logf("INFO finalizing epoch")
	return nil
}

func (c *Client) UpdateExchangeRate(ctx context.Context, lendingManagerAddress string) error {
	c.logger.Logf("INFO updating exchange rate for LendingManager %s", lendingManagerAddress)
	
	if c.contractClient == nil {
		c.logger.Logf("WARN contract client not initialized, skipping updateExchangeRate call")
		return nil
	}
	
	return c.contractClient.UpdateExchangeRate(ctx, lendingManagerAddress)
}

func (c *Client) AllocateYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string) error {
	c.logger.Logf("INFO allocating yield to epoch %s for vault %s", epochId.String(), vaultAddress)
	
	if c.contractClient == nil {
		c.logger.Logf("WARN contract client not initialized, skipping allocateYieldToEpoch call")
		return nil
	}
	
	return c.contractClient.AllocateYieldToEpoch(ctx, epochId, vaultAddress)
}

func (c *Client) EndEpochWithSubsidies(ctx context.Context, epochId *big.Int, vaultAddress string, merkleRoot [32]byte, subsidiesDistributed *big.Int) error {
	
	if c.contractClient == nil {
		c.logger.Logf("WARN contract client not initialized, skipping endEpochWithSubsidies call")
		return nil
	}
	
	return c.contractClient.EndEpochWithSubsidies(ctx, epochId, vaultAddress, merkleRoot, subsidiesDistributed)
}
