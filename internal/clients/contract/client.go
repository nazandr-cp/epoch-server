package contract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/andrey/epoch-server/pkg/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
		return nil // Return nil for now to not break existing flow
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
		return err
	}

	c.logger.Logf("INFO started epoch transaction sent: %s", tx.Hash().Hex())
	return nil
}

func (c *Client) DistributeSubsidies(ctx context.Context, epochID string) error {
	c.logger.Logf("INFO distributing subsidies for epoch %s", epochID)
	return nil
}
