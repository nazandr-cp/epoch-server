package contract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/andrey/epoch-server/pkg/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	bind_v2 "github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
	logger       lgr.L
	ethConfig    EthereumConfig
	contracts    ContractAddresses
	ethClient    *ethclient.Client
	privateKey   *ecdsa.PrivateKey
	epochManager *contracts.IEpochManager
}

func NewClient(logger lgr.L) *Client {
	return &Client{
		logger: logger,
	}
}

func NewClientWithConfig(logger lgr.L, ethConfig EthereumConfig, contracts ContractAddresses) (*Client, error) {
	client := &Client{
		logger:    logger,
		ethConfig: ethConfig,
		contracts: contracts,
	}
	
	// Initialize Ethereum client and contracts
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
	if c.contracts.EpochManager == "" {
		return fmt.Errorf("EpochManager contract address is required")
	}

	// Connect to Ethereum client
	ethClient, err := ethclient.Dial(c.ethConfig.RPCURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum RPC: %w", err)
	}
	c.ethClient = ethClient

	// Parse private key (strip 0x prefix if present)
	privateKeyHex := c.ethConfig.PrivateKey
	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	c.privateKey = privateKey

	// Initialize EpochManager contract
	c.epochManager = contracts.NewIEpochManager()

	return nil
}

func (c *Client) StartEpoch(ctx context.Context, epochID string) error {
	c.logger.Logf("INFO starting epoch %s", epochID)
	
	if c.ethClient == nil || c.privateKey == nil {
		c.logger.Logf("ERROR Ethereum client not initialized")
		return fmt.Errorf("ethereum client not initialized")
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

	// Create contract instance and call function
	contractAddr := common.HexToAddress(c.contracts.EpochManager)
	contractInstance := c.epochManager.Instance(c.ethClient, contractAddr)

	// Call startEpoch() function using simplified interface
	data := c.epochManager.PackStartEpoch()
	tx, err := contractInstance.RawTransact(opts, data)
	
	if err != nil {
		c.logger.Logf("ERROR failed to call startEpoch: %v", err)
		return fmt.Errorf("failed to call startEpoch: %w", err)
	}

	c.logger.Logf("INFO started epoch transaction sent: %s", tx.Hash().Hex())
	
	// Wait for transaction to be mined and check if it was successful
	c.logger.Logf("DEBUG about to call bind.WaitMined for transaction %s", tx.Hash().Hex())
	c.logger.Logf("INFO waiting for transaction %s to be mined...", tx.Hash().Hex())
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
	c.logger.Logf("INFO getting current epoch ID")
	
	if c.ethClient == nil || c.privateKey == nil {
		c.logger.Logf("WARN Ethereum client not initialized, returning epoch ID 1")
		return big.NewInt(1), nil
	}

	// Create contract instance for EpochManager
	contractAddr := common.HexToAddress(c.contracts.EpochManager)
	contractInstance := c.epochManager.Instance(c.ethClient, contractAddr)

	// Call getCurrentEpochId() function using abigen v2
	callOpts := &bind_v2.CallOpts{Context: ctx}
	var result []interface{}
	err := contractInstance.Call(callOpts, &result, "getCurrentEpochId")
	if err != nil {
		c.logger.Logf("ERROR failed to call getCurrentEpochId: %v", err)
		return nil, fmt.Errorf("failed to call getCurrentEpochId: %w", err)
	}

	// Extract epoch ID from result
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

func (c *Client) AllocateYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string) error {
	c.logger.Logf("INFO allocating yield to epoch %s for vault %s", epochId.String(), vaultAddress)
	
	if c.ethClient == nil || c.privateKey == nil {
		c.logger.Logf("WARN Ethereum client not initialized, skipping allocateYieldToEpoch call")
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

	// Create vault contract instance and call allocateYieldToEpoch function
	vaultAddr := common.HexToAddress(vaultAddress)
	
	// Create the function call data for allocateYieldToEpoch(uint256)
	methodID := crypto.Keccak256([]byte("allocateYieldToEpoch(uint256)"))[:4]
	epochIdPacked := common.LeftPadBytes(epochId.Bytes(), 32)
	data := append(methodID, epochIdPacked...)

	// Create vault contract instance (not epoch manager) since allocateYieldToEpoch is on the vault
	contractInstance := c.epochManager.Instance(c.ethClient, vaultAddr)
	tx, err := contractInstance.RawTransact(opts, data)
	
	if err != nil {
		c.logger.Logf("ERROR failed to call allocateYieldToEpoch: %v", err)
		return fmt.Errorf("failed to call allocateYieldToEpoch: %w", err)
	}

	c.logger.Logf("INFO allocateYieldToEpoch transaction sent: %s", tx.Hash().Hex())
	
	// Wait for transaction to be mined and check if it was successful
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

func (c *Client) EndEpochWithSubsidies(ctx context.Context, epochId *big.Int, vaultAddress string, merkleRoot [32]byte, subsidiesDistributed *big.Int) error {
	c.logger.Logf("INFO ending epoch %s with subsidies: vault=%s, merkleRoot=%x, subsidies=%s", 
		epochId.String(), vaultAddress, merkleRoot, subsidiesDistributed.String())
	
	if c.ethClient == nil || c.privateKey == nil {
		c.logger.Logf("ERROR Ethereum client not initialized")
		return fmt.Errorf("ethereum client not initialized")
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

	// Create contract instance and call function
	contractAddr := common.HexToAddress(c.contracts.EpochManager)
	contractInstance := c.epochManager.Instance(c.ethClient, contractAddr)

	// Convert vault address string to common.Address
	vaultAddr := common.HexToAddress(vaultAddress)

	// Call endEpochWithSubsidies function
	data := c.epochManager.PackEndEpochWithSubsidies(epochId, vaultAddr, merkleRoot, subsidiesDistributed)
	tx, err := contractInstance.RawTransact(opts, data)
	
	if err != nil {
		c.logger.Logf("ERROR failed to call endEpochWithSubsidies: %v", err)
		return fmt.Errorf("failed to call endEpochWithSubsidies: %w", err)
	}

	c.logger.Logf("INFO endEpochWithSubsidies transaction sent: %s", tx.Hash().Hex())
	
	// Wait for transaction to be mined and check if it was successful
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
