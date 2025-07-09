package blockchain

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
	AllocateCumulativeYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string, amount *big.Int) error
	EndEpochWithSubsidies(ctx context.Context, epochId *big.Int, vaultAddress string, merkleRoot [32]byte, subsidiesDistributed *big.Int) error
}

type EpochClient struct {
	logger         lgr.L
	contractClient ContractClient
}

func NewEpochClient(logger lgr.L) *EpochClient {
	return &EpochClient{
		logger: logger,
	}
}

func NewEpochClientWithContract(logger lgr.L, contractClient ContractClient) *EpochClient {
	return &EpochClient{
		logger:         logger,
		contractClient: contractClient,
	}
}

type EpochInfo struct {
	EndTime int64 `json:"endTime"`
}

func (c *EpochClient) Current() EpochInfo {
	return EpochInfo{
		EndTime: time.Now().Unix(),
	}
}

func (c *EpochClient) GetCurrentEpochId(ctx context.Context) (*big.Int, error) {
	c.logger.Logf("INFO getting current epoch ID")
	
	if c.contractClient == nil {
		c.logger.Logf("WARN contract client not initialized, returning epoch ID 1")
		return big.NewInt(1), nil
	}
	
	return c.contractClient.GetCurrentEpochId(ctx)
}

func (c *EpochClient) FinalizeEpoch() error {
	c.logger.Logf("INFO finalizing epoch")
	return nil
}

func (c *EpochClient) UpdateExchangeRate(ctx context.Context, lendingManagerAddress string) error {
	c.logger.Logf("INFO updating exchange rate for LendingManager %s", lendingManagerAddress)
	
	if c.contractClient == nil {
		c.logger.Logf("WARN contract client not initialized, skipping updateExchangeRate call")
		return nil
	}
	
	return c.contractClient.UpdateExchangeRate(ctx, lendingManagerAddress)
}

func (c *EpochClient) AllocateYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string) error {
	c.logger.Logf("INFO allocating yield to epoch %s for vault %s", epochId.String(), vaultAddress)
	
	if c.contractClient == nil {
		c.logger.Logf("WARN contract client not initialized, skipping allocateYieldToEpoch call")
		return nil
	}
	
	return c.contractClient.AllocateYieldToEpoch(ctx, epochId, vaultAddress)
}

func (c *EpochClient) AllocateCumulativeYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string, amount *big.Int) error {
	c.logger.Logf("INFO allocating cumulative yield %s to epoch %s for vault %s", amount.String(), epochId.String(), vaultAddress)
	
	if c.contractClient == nil {
		c.logger.Logf("WARN contract client not initialized, skipping allocateCumulativeYieldToEpoch call")
		return nil
	}
	
	return c.contractClient.AllocateCumulativeYieldToEpoch(ctx, epochId, vaultAddress, amount)
}

func (c *EpochClient) EndEpochWithSubsidies(ctx context.Context, epochId *big.Int, vaultAddress string, merkleRoot [32]byte, subsidiesDistributed *big.Int) error {
	
	if c.contractClient == nil {
		c.logger.Logf("WARN contract client not initialized, skipping endEpochWithSubsidies call")
		return nil
	}
	
	return c.contractClient.EndEpochWithSubsidies(ctx, epochId, vaultAddress, merkleRoot, subsidiesDistributed)
}
