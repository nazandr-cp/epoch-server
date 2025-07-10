package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/andrey/epoch-server/pkg/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-pkgz/lgr"
)

type SubsidizerEthereumConfig struct {
	RPCURL     string
	PrivateKey string
	GasLimit   uint64
	GasPrice   string
}

type SubsidizerClient struct {
	logger       lgr.L
	ethConfig    SubsidizerEthereumConfig
	contractAddr string
	ethClient    *ethclient.Client
	privateKey   *ecdsa.PrivateKey
	subsidizer   *contracts.IDebtSubsidizer
}

// NewClient creates a mock client for backward compatibility
func NewSubsidizerClient(logger lgr.L) *SubsidizerClient {
	return &SubsidizerClient{
		logger: logger,
	}
}

// NewClientWithConfig creates a real blockchain client
func NewSubsidizerClientWithConfig(
	logger lgr.L,
	ethConfig SubsidizerEthereumConfig,
	contractAddr string,
) (*SubsidizerClient, error) {
	client := &SubsidizerClient{
		logger:       logger,
		ethConfig:    ethConfig,
		contractAddr: contractAddr,
	}

	// Initialize Ethereum client and contract
	if err := client.initialize(); err != nil {
		logger.Logf("ERROR failed to initialize subsidizer client: %v", err)
		return nil, err
	}

	return client, nil
}

func (c *SubsidizerClient) initialize() error {
	// Validate required configuration
	if c.ethConfig.RPCURL == "" {
		return fmt.Errorf("RPC URL is required")
	}
	if c.ethConfig.PrivateKey == "" {
		return fmt.Errorf("private key is required")
	}
	if c.contractAddr == "" {
		return fmt.Errorf("DebtSubsidizer contract address is required")
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

	// Initialize DebtSubsidizer contract
	c.subsidizer = contracts.NewIDebtSubsidizer()

	c.logger.Logf("INFO DebtSubsidizer client initialized successfully")
	return nil
}

func (c *SubsidizerClient) UpdateMerkleRoot(
	ctx context.Context,
	vaultId string,
	root [32]byte,
	totalSubsidies *big.Int,
) error {
	// If not initialized (mock mode), just log
	if c.ethClient == nil {
		c.logger.Logf("INFO [MOCK] updating merkle root for vault %s: %x", vaultId, root)
		return nil
	}

	// Real implementation
	c.logger.Logf("INFO updating merkle root for vault %s: %x", vaultId, root)

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

	// Build transaction data
	vaultAddress := common.HexToAddress(vaultId)
	data := c.subsidizer.PackUpdateMerkleRoot(vaultAddress, root, totalSubsidies)

	// Create contract instance and submit transaction
	contractAddr := common.HexToAddress(c.contractAddr)
	contractInstance := c.subsidizer.Instance(c.ethClient, contractAddr)
	tx, err := contractInstance.RawTransact(opts, data)

	if err != nil {
		c.logger.Logf("ERROR failed to call updateMerkleRoot: %v", err)
		return fmt.Errorf("failed to call updateMerkleRoot: %w", err)
	}

	c.logger.Logf("INFO updateMerkleRoot transaction sent: %s", tx.Hash().Hex())
	return nil
}

func (c *SubsidizerClient) UpdateMerkleRootAndWaitForConfirmation(
	ctx context.Context,
	vaultId string,
	root [32]byte,
	totalSubsidies *big.Int,
) error {
	// If not initialized (mock mode), simulate the old behavior
	if c.ethClient == nil {
		c.logger.Logf("INFO [MOCK] updating merkle root for vault %s: %x", vaultId, root)
		c.logger.Logf("INFO [MOCK] submitting UpdateMerkleRoot transaction for vault %s", vaultId)
		c.logger.Logf("INFO [MOCK] waiting for transaction confirmation for vault %s", vaultId)
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for confirmation: %w", ctx.Err())
		case <-time.After(100 * time.Millisecond): // Simulate mining time
			c.logger.Logf("INFO [MOCK] transaction confirmed for vault %s", vaultId)
			return nil
		}
	}

	// Real implementation
	c.logger.Logf("INFO updating merkle root for vault %s: %x", vaultId, root)

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

	// Build transaction data
	vaultAddress := common.HexToAddress(vaultId)
	data := c.subsidizer.PackUpdateMerkleRoot(vaultAddress, root, totalSubsidies)

	// Create contract instance and submit transaction
	contractAddr := common.HexToAddress(c.contractAddr)
	contractInstance := c.subsidizer.Instance(c.ethClient, contractAddr)
	tx, err := contractInstance.RawTransact(opts, data)

	if err != nil {
		c.logger.Logf("ERROR failed to call updateMerkleRoot: %v", err)
		return fmt.Errorf("failed to call updateMerkleRoot: %w", err)
	}

	c.logger.Logf("INFO submitting UpdateMerkleRoot transaction for vault %s", vaultId)
	c.logger.Logf("INFO updateMerkleRoot transaction sent: %s", tx.Hash().Hex())

	// Wait for transaction to be mined and check if it was successful
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
