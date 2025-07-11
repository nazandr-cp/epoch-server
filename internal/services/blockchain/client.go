package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/andrey/epoch-server/internal/infra/blockchain"
	"github.com/andrey/epoch-server/pkg/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	bind_v2 "github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-pkgz/lgr"
)

type Client struct {
	logger       lgr.L
	ethConfig    blockchain.Config
	ethClient    *ethclient.Client
	privateKey   *ecdsa.PrivateKey
	epochManager *contracts.IEpochManager
	subsidizer   *contracts.IDebtSubsidizer
}

// ProvideClient creates a new blockchain client implementation
func ProvideClient(logger lgr.L) blockchain.BlockchainClient {
	return &Client{
		logger: logger,
	}
}

// ProvideClientWithConfig creates a blockchain client with configuration
func ProvideClientWithConfig(logger lgr.L, config blockchain.Config) (blockchain.BlockchainClient, error) {
	client := &Client{
		logger:    logger,
		ethConfig: config,
	}

	if err := client.initialize(); err != nil {
		logger.Logf("ERROR failed to initialize contract client: %v", err)
		return nil, err
	}

	return client, nil
}

func (c *Client) initialize() error {
	// Validate required configuration
	if c.ethConfig.RPCURL == "" {
		return fmt.Errorf("RPC URL is required")
	}
	if c.ethConfig.PrivateKey == "" {
		return fmt.Errorf("private key is required")
	}
	if c.ethConfig.EpochManager == "" {
		return fmt.Errorf("EpochManager contract address is required")
	}

	ethClient, err := ethclient.Dial(c.ethConfig.RPCURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum RPC: %w", err)
	}
	c.ethClient = ethClient

	privateKeyHex := c.ethConfig.PrivateKey
	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	c.privateKey = privateKey
	c.epochManager = contracts.NewIEpochManager()
	c.subsidizer = contracts.NewIDebtSubsidizer()

	return nil
}

func (c *Client) StartEpoch(ctx context.Context) error {
	c.logger.Logf("INFO starting epoch")

	if c.ethClient == nil || c.privateKey == nil {
		c.logger.Logf("ERROR Ethereum client not initialized")
		return fmt.Errorf("ethereum client not initialized")
	}

	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		c.logger.Logf("ERROR failed to get chain ID: %v", err)
		return err
	}

	gasPrice, _ := new(big.Int).SetString(c.ethConfig.GasPrice, 10)
	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		c.logger.Logf("ERROR failed to create transactor: %v", err)
		return err
	}
	opts.GasLimit = c.ethConfig.GasLimit
	opts.GasPrice = gasPrice
	opts.Context = ctx

	contractAddr := common.HexToAddress(c.ethConfig.EpochManager)
	contractInstance := c.epochManager.Instance(c.ethClient, contractAddr)

	data := c.epochManager.PackStartEpoch()
	tx, err := contractInstance.RawTransact(opts, data)

	if err != nil {
		c.logger.Logf("ERROR failed to call startEpoch: %v", err)
		return fmt.Errorf("failed to call startEpoch: %w", err)
	}

	c.logger.Logf("INFO started epoch transaction sent: %s", tx.Hash().Hex())

	receipt, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		c.logger.Logf("ERROR failed to wait for startEpoch transaction %s: %v", tx.Hash().Hex(), err)
		return fmt.Errorf("failed to wait for startEpoch transaction: %w", err)
	}

	c.logger.Logf("INFO transaction %s mined in block %d", tx.Hash().Hex(), receipt.BlockNumber.Uint64())

	if receipt.Status == 0 {
		c.logger.Logf("ERROR startEpoch transaction failed: %s", tx.Hash().Hex())
		return fmt.Errorf("startEpoch transaction failed with hash %s", tx.Hash().Hex())
	}

	c.logger.Logf("INFO startEpoch transaction successful: %s", tx.Hash().Hex())
	return nil
}

func (c *Client) DistributeSubsidies(ctx context.Context, epochID string) error {
	c.logger.Logf("INFO distributing subsidies for epoch %s", epochID)
	return nil
}

func (c *Client) GetCurrentEpochId(ctx context.Context) (*big.Int, error) {
	if c.ethClient == nil || c.privateKey == nil {
		c.logger.Logf("WARN Ethereum client not initialized, returning epoch ID 1")
		return big.NewInt(1), nil
	}

	contractAddr := common.HexToAddress(c.ethConfig.EpochManager)
	contractInstance := c.epochManager.Instance(c.ethClient, contractAddr)

	callOpts := &bind_v2.CallOpts{Context: ctx}
	var result []interface{}
	err := contractInstance.Call(callOpts, &result, "getCurrentEpochId")
	if err != nil {
		c.logger.Logf("ERROR failed to call getCurrentEpochId: %v", err)
		return nil, fmt.Errorf("failed to call getCurrentEpochId: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no result returned from getCurrentEpochId")
	}
	epochId, ok := result[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from getCurrentEpochId")
	}

	c.logger.Logf("INFO current epoch ID: %s", epochId.String())
	return epochId, nil
}

func (c *Client) UpdateExchangeRate(ctx context.Context, lendingManagerAddress string) error {
	c.logger.Logf("INFO updating exchange rate for LendingManager %s", lendingManagerAddress)

	if c.ethClient == nil || c.privateKey == nil {
		c.logger.Logf("ERROR Ethereum client not initialized")
		return fmt.Errorf("ethereum client not initialized")
	}

	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		c.logger.Logf("ERROR failed to get chain ID: %v", err)
		return err
	}

	gasPrice, _ := new(big.Int).SetString(c.ethConfig.GasPrice, 10)
	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		c.logger.Logf("ERROR failed to create transactor: %v", err)
		return err
	}
	opts.GasLimit = c.ethConfig.GasLimit
	opts.GasPrice = gasPrice
	opts.Context = ctx

	lendingManagerAddr := common.HexToAddress(lendingManagerAddress)
	methodID := crypto.Keccak256([]byte("updateExchangeRate()"))[:4]
	data := methodID

	contractInstance := c.epochManager.Instance(c.ethClient, lendingManagerAddr)
	tx, err := contractInstance.RawTransact(opts, data)

	if err != nil {
		c.logger.Logf("ERROR failed to call updateExchangeRate: %v", err)
		return fmt.Errorf("failed to call updateExchangeRate: %w", err)
	}

	c.logger.Logf("INFO updateExchangeRate transaction sent: %s", tx.Hash().Hex())

	receipt, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		c.logger.Logf("ERROR failed to wait for updateExchangeRate transaction: %v", err)
		return fmt.Errorf("failed to wait for updateExchangeRate transaction: %w", err)
	}

	if receipt.Status == 0 {
		c.logger.Logf("ERROR updateExchangeRate transaction failed: %s", tx.Hash().Hex())
		return fmt.Errorf("updateExchangeRate transaction failed with hash %s", tx.Hash().Hex())
	}

	c.logger.Logf("INFO updateExchangeRate transaction successful: %s", tx.Hash().Hex())
	return nil
}

func (c *Client) AllocateYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string) error {
	c.logger.Logf("INFO allocating yield to epoch %s for vault %s", epochId.String(), vaultAddress)

	if c.ethClient == nil || c.privateKey == nil {
		c.logger.Logf("WARN Ethereum client not initialized, skipping allocateYieldToEpoch call")
		return nil
	}

	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		c.logger.Logf("ERROR failed to get chain ID: %v", err)
		return err
	}

	gasPrice, _ := new(big.Int).SetString(c.ethConfig.GasPrice, 10)
	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		c.logger.Logf("ERROR failed to create transactor: %v", err)
		return err
	}
	opts.GasLimit = c.ethConfig.GasLimit
	opts.GasPrice = gasPrice
	opts.Context = ctx

	vaultAddr := common.HexToAddress(vaultAddress)
	methodID := crypto.Keccak256([]byte("allocateYieldToEpoch(uint256)"))[:4]
	epochIdPacked := common.LeftPadBytes(epochId.Bytes(), 32)
	data := append(methodID, epochIdPacked...)

	contractInstance := c.epochManager.Instance(c.ethClient, vaultAddr)
	tx, err := contractInstance.RawTransact(opts, data)

	if err != nil {
		c.logger.Logf("ERROR failed to call allocateYieldToEpoch: %v", err)
		return fmt.Errorf("failed to call allocateYieldToEpoch: %w", err)
	}

	c.logger.Logf("INFO allocateYieldToEpoch transaction sent: %s", tx.Hash().Hex())

	receipt, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		c.logger.Logf("ERROR failed to wait for allocateYieldToEpoch transaction: %v", err)
		return fmt.Errorf("failed to wait for allocateYieldToEpoch transaction: %w", err)
	}

	if receipt.Status == 0 {
		c.logger.Logf("ERROR allocateYieldToEpoch transaction failed: %s", tx.Hash().Hex())
		return fmt.Errorf("allocateYieldToEpoch transaction failed with hash %s", tx.Hash().Hex())
	}

	c.logger.Logf("INFO allocateYieldToEpoch transaction successful: %s", tx.Hash().Hex())
	return nil
}

func (c *Client) AllocateCumulativeYieldToEpoch(
	ctx context.Context,
	epochId *big.Int,
	vaultAddress string,
	amount *big.Int,
) error {
	c.logger.Logf(
		"INFO allocating cumulative yield %s to epoch %s for vault %s",
		amount.String(),
		epochId.String(),
		vaultAddress,
	)

	if c.ethClient == nil || c.privateKey == nil {
		c.logger.Logf("WARN Ethereum client not initialized, skipping allocateCumulativeYieldToEpoch call")
		return nil
	}

	// Get chain ID for signing
	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		c.logger.Logf("ERROR failed to get chain ID: %v", err)
		return err
	}

	// Create transaction options with signer
	gasPrice, _ := new(big.Int).SetString(c.ethConfig.GasPrice, 10)
	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		c.logger.Logf("ERROR failed to create transactor: %v", err)
		return err
	}
	opts.GasLimit = c.ethConfig.GasLimit
	opts.GasPrice = gasPrice
	opts.Context = ctx

	vaultAddr := common.HexToAddress(vaultAddress)
	methodID := crypto.Keccak256([]byte("allocateCumulativeYieldToEpoch(uint256,uint256)"))[:4]
	epochIdPacked := common.LeftPadBytes(epochId.Bytes(), 32)
	amountPacked := common.LeftPadBytes(amount.Bytes(), 32)
	data := append(methodID, epochIdPacked...)
	data = append(data, amountPacked...)

	contractInstance := c.epochManager.Instance(c.ethClient, vaultAddr)
	tx, err := contractInstance.RawTransact(opts, data)

	if err != nil {
		c.logger.Logf("ERROR failed to call allocateCumulativeYieldToEpoch: %v", err)
		return fmt.Errorf("failed to call allocateCumulativeYieldToEpoch: %w", err)
	}

	c.logger.Logf("INFO allocateCumulativeYieldToEpoch transaction sent: %s", tx.Hash().Hex())

	receipt, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		c.logger.Logf("ERROR failed to wait for allocateCumulativeYieldToEpoch transaction: %v", err)
		return fmt.Errorf("failed to wait for allocateCumulativeYieldToEpoch transaction: %w", err)
	}

	if receipt.Status == 0 {
		c.logger.Logf("ERROR allocateCumulativeYieldToEpoch transaction failed: %s", tx.Hash().Hex())
		return fmt.Errorf("allocateCumulativeYieldToEpoch transaction failed with hash %s", tx.Hash().Hex())
	}

	c.logger.Logf("INFO allocateCumulativeYieldToEpoch transaction successful: %s", tx.Hash().Hex())
	return nil
}

func (c *Client) EndEpochWithSubsidies(
	ctx context.Context,
	epochId *big.Int,
	vaultAddress string,
	merkleRoot [32]byte,
	subsidiesDistributed *big.Int,
) error {
	c.logger.Logf("INFO ending epoch %s with subsidies: vault=%s, merkleRoot=%x, subsidies=%s",
		epochId.String(), vaultAddress, merkleRoot, subsidiesDistributed.String())

	if c.ethClient == nil || c.privateKey == nil {
		c.logger.Logf("ERROR Ethereum client not initialized")
		return fmt.Errorf("ethereum client not initialized")
	}

	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		c.logger.Logf("ERROR failed to get chain ID: %v", err)
		return err
	}

	gasPrice, _ := new(big.Int).SetString(c.ethConfig.GasPrice, 10)
	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		c.logger.Logf("ERROR failed to create transactor: %v", err)
		return err
	}
	opts.GasLimit = c.ethConfig.GasLimit
	opts.GasPrice = gasPrice
	opts.Context = ctx

	contractAddr := common.HexToAddress(c.ethConfig.EpochManager)
	contractInstance := c.epochManager.Instance(c.ethClient, contractAddr)
	vaultAddr := common.HexToAddress(vaultAddress)
	data := c.epochManager.PackEndEpochWithSubsidies(epochId, vaultAddr, merkleRoot, subsidiesDistributed)
	tx, err := contractInstance.RawTransact(opts, data)

	if err != nil {
		c.logger.Logf("ERROR failed to call endEpochWithSubsidies: %v", err)
		return fmt.Errorf("failed to call endEpochWithSubsidies: %w", err)
	}

	c.logger.Logf("INFO endEpochWithSubsidies transaction sent: %s", tx.Hash().Hex())

	receipt, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		c.logger.Logf("ERROR failed to wait for endEpochWithSubsidies transaction: %v", err)
		return fmt.Errorf("failed to wait for endEpochWithSubsidies transaction: %w", err)
	}

	if receipt.Status == 0 {
		c.logger.Logf("ERROR endEpochWithSubsidies transaction failed: %s", tx.Hash().Hex())
		return fmt.Errorf("endEpochWithSubsidies transaction failed with hash %s", tx.Hash().Hex())
	}

	c.logger.Logf("INFO endEpochWithSubsidies transaction successful: %s", tx.Hash().Hex())
	return nil
}

func (c *Client) ForceEndEpochWithZeroYield(ctx context.Context, epochId *big.Int, vaultAddress string) error {
	c.logger.Logf("INFO force ending epoch %s with zero yield: vault=%s", epochId.String(), vaultAddress)

	if c.ethClient == nil || c.privateKey == nil {
		c.logger.Logf("ERROR Ethereum client not initialized")
		return fmt.Errorf("ethereum client not initialized")
	}

	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		c.logger.Logf("ERROR failed to get chain ID: %v", err)
		return err
	}

	gasPrice, _ := new(big.Int).SetString(c.ethConfig.GasPrice, 10)
	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		c.logger.Logf("ERROR failed to create transactor: %v", err)
		return err
	}
	opts.GasLimit = c.ethConfig.GasLimit
	opts.GasPrice = gasPrice
	opts.Context = ctx

	vaultAddr := common.HexToAddress(vaultAddress)
	data := c.epochManager.PackForceEndEpochWithZeroYield(epochId, vaultAddr)

	contractAddr := common.HexToAddress(c.ethConfig.EpochManager)
	contractInstance := c.epochManager.Instance(c.ethClient, contractAddr)
	tx, err := contractInstance.RawTransact(opts, data)
	if err != nil {
		c.logger.Logf("ERROR failed to call forceEndEpochWithZeroYield: %v", err)
		return fmt.Errorf("failed to call forceEndEpochWithZeroYield: %w", err)
	}

	c.logger.Logf("INFO forceEndEpochWithZeroYield transaction sent: %s", tx.Hash().Hex())

	c.logger.Logf("INFO forceEndEpochWithZeroYield transaction successful: %s", tx.Hash().Hex())
	return nil
}

func (c *Client) UpdateMerkleRoot(
	ctx context.Context,
	vaultId string,
	root [32]byte,
	totalSubsidies *big.Int,
) error {
	if c.ethClient == nil {
		c.logger.Logf("INFO [MOCK] updating merkle root for vault %s: %x", vaultId, root)
		return nil
	}

	c.logger.Logf("INFO updating merkle root for vault %s: %x", vaultId, root)

	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		c.logger.Logf("ERROR failed to get chain ID: %v", err)
		return err
	}

	gasPrice, _ := new(big.Int).SetString(c.ethConfig.GasPrice, 10)
	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		c.logger.Logf("ERROR failed to create transactor: %v", err)
		return err
	}
	opts.GasLimit = c.ethConfig.GasLimit
	opts.GasPrice = gasPrice
	opts.Context = ctx

	vaultAddress := common.HexToAddress(vaultId)
	data := c.subsidizer.PackUpdateMerkleRoot(vaultAddress, root, totalSubsidies)

	contractAddr := common.HexToAddress(c.ethConfig.DebtSubsidizer)
	contractInstance := c.subsidizer.Instance(c.ethClient, contractAddr)
	tx, err := contractInstance.RawTransact(opts, data)

	if err != nil {
		c.logger.Logf("ERROR failed to call updateMerkleRoot: %v", err)
		return fmt.Errorf("failed to call updateMerkleRoot: %w", err)
	}

	c.logger.Logf("INFO updateMerkleRoot transaction sent: %s", tx.Hash().Hex())
	return nil
}

func (c *Client) UpdateMerkleRootAndWaitForConfirmation(
	ctx context.Context,
	vaultId string,
	root [32]byte,
	totalSubsidies *big.Int,
) error {
	if c.ethClient == nil {
		c.logger.Logf("INFO [MOCK] updating merkle root for vault %s: %x", vaultId, root)
		c.logger.Logf("INFO [MOCK] submitting UpdateMerkleRoot transaction for vault %s", vaultId)
		c.logger.Logf("INFO [MOCK] waiting for transaction confirmation for vault %s", vaultId)
		return nil
	}

	c.logger.Logf("INFO updating merkle root for vault %s: %x", vaultId, root)

	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		c.logger.Logf("ERROR failed to get chain ID: %v", err)
		return err
	}

	gasPrice, _ := new(big.Int).SetString(c.ethConfig.GasPrice, 10)
	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		c.logger.Logf("ERROR failed to create transactor: %v", err)
		return err
	}
	opts.GasLimit = c.ethConfig.GasLimit
	opts.GasPrice = gasPrice
	opts.Context = ctx

	vaultAddress := common.HexToAddress(vaultId)
	data := c.subsidizer.PackUpdateMerkleRoot(vaultAddress, root, totalSubsidies)

	contractAddr := common.HexToAddress(c.ethConfig.DebtSubsidizer)
	contractInstance := c.subsidizer.Instance(c.ethClient, contractAddr)
	tx, err := contractInstance.RawTransact(opts, data)

	if err != nil {
		c.logger.Logf("ERROR failed to call updateMerkleRoot: %v", err)
		return fmt.Errorf("failed to call updateMerkleRoot: %w", err)
	}

	c.logger.Logf("INFO submitting UpdateMerkleRoot transaction for vault %s", vaultId)
	c.logger.Logf("INFO updateMerkleRoot transaction sent: %s", tx.Hash().Hex())

	c.logger.Logf("INFO waiting for transaction confirmation for vault %s", vaultId)
	receipt, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		c.logger.Logf("ERROR failed to wait for updateMerkleRoot transaction: %v", err)
		return fmt.Errorf("failed to wait for updateMerkleRoot transaction: %w", err)
	}

	if receipt.Status == 0 {
		c.logger.Logf("ERROR updateMerkleRoot transaction failed: %s", tx.Hash().Hex())
		return fmt.Errorf("updateMerkleRoot transaction failed with hash %s", tx.Hash().Hex())
	}

	c.logger.Logf(
		"INFO transaction confirmed for vault %s (block: %d, gas used: %d)",
		vaultId,
		receipt.BlockNumber.Uint64(),
		receipt.GasUsed,
	)
	return nil
}
