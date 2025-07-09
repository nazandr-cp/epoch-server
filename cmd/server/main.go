package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andrey/epoch-server/internal/api"
	"github.com/andrey/epoch-server/internal/infra/blockchain"
	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/infra/logging"
	"github.com/andrey/epoch-server/internal/infra/storage"
	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/andrey/epoch-server/internal/services/epoch/epochimpl"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/andrey/epoch-server/internal/services/merkle/merkleimpl"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logging
	logConfig := logging.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
		Output: cfg.Logging.Output,
	}
	logger := logging.NewWithConfig(logConfig)

	// Setup infrastructure components
	subgraphClient := subgraph.NewClient(cfg.Subgraph.Endpoint)

	// Perform subgraph health check during startup
	logger.Logf("INFO checking subgraph connectivity at %s", cfg.Subgraph.Endpoint)
	ctx := context.Background()
	if err := subgraphClient.HealthCheck(ctx); err != nil {
		log.Fatalf("Failed to connect to subgraph: %v", err)
	}
	logger.Logf("INFO subgraph health check passed")

	// Setup blockchain clients
	ethConfig := blockchain.EthereumConfig{
		RPCURL:     cfg.Ethereum.RPCURL,
		PrivateKey: cfg.Ethereum.PrivateKey,
		GasLimit:   cfg.Ethereum.GasLimit,
		GasPrice:   cfg.Ethereum.GasPrice,
	}

	contractAddresses := blockchain.ContractAddresses{
		Comptroller:        cfg.Contracts.Comptroller,
		EpochManager:       cfg.Contracts.EpochManager,
		DebtSubsidizer:     cfg.Contracts.DebtSubsidizer,
		LendingManager:     cfg.Contracts.LendingManager,
		CollectionRegistry: cfg.Contracts.CollectionRegistry,
	}

	contractClient, err := blockchain.NewClientWithConfig(logger, ethConfig, contractAddresses)
	if err != nil {
		log.Fatalf("Failed to initialize contract client: %v", err)
	}

	// Setup storage
	_ = storage.NewClient(logger)

	// Setup merkle service dependencies
	calculator := merkleimpl.NewCalculator()

	// Setup services
	epochService := epochimpl.New(contractClient, subgraphClient, calculator, logger, cfg)
	
	// Create a mock merkle service for now
	merkleService := &mockMerkleService{logger: logger}
	
	// Create a mock subsidy service for now  
	subsidyService := &mockSubsidyService{logger: logger}

	// Setup scheduler - TODO: This needs to be updated to work with the new service interfaces
	// schedulerInstance := scheduler.NewScheduler(cfg.Scheduler.Interval, epochService, logger, cfg)
	// go schedulerInstance.Start(ctx)

	// Setup and start HTTP server
	server := api.NewServer(epochService, subsidyService, merkleService, logger, cfg)
	
	if err := server.Start(); err != nil {
		logger.Logf("ERROR server failed to start: %v", err)
	}
}

// Mock services for now - these will be replaced with proper implementations

type mockMerkleService struct {
	logger interface{}
}

func (m *mockMerkleService) GenerateUserMerkleProof(ctx context.Context, userAddress, vaultAddress string) (*merkle.UserMerkleProofResponse, error) {
	return nil, fmt.Errorf("mock service not implemented")
}

func (m *mockMerkleService) GenerateHistoricalMerkleProof(ctx context.Context, userAddress, vaultAddress, epochNumber string) (*merkle.UserMerkleProofResponse, error) {
	return nil, fmt.Errorf("mock service not implemented")
}

type mockSubsidyService struct {
	logger interface{}
}

func (m *mockSubsidyService) DistributeSubsidies(ctx context.Context, vaultId string) error {
	return fmt.Errorf("mock service not implemented")
}