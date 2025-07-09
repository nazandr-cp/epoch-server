// @title Epoch Server API
// @version 1.0
// @description Epoch Server for managing NFT collection-backed lending epochs, subsidies, and merkle proofs
// @termsOfService http://lend.fam/terms/
// @contact.name API Support
// @contact.url http://lend.fam/support
// @contact.email support@lend.fam
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
// @accept json
// @produce json
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/andrey/epoch-server/internal/api"
	"github.com/andrey/epoch-server/internal/infra/blockchain"
	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/infra/logging"
	"github.com/andrey/epoch-server/internal/infra/storage"
	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/andrey/epoch-server/internal/services/epoch/epochimpl"
	"github.com/andrey/epoch-server/internal/services/merkle/merkleimpl"
	"github.com/andrey/epoch-server/internal/services/scheduler"
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

	// Setup database
	storageConfig := storage.StorageConfig{
		Type: cfg.Database.Type,
		Path: cfg.Database.ConnectionString,
	}
	dbWrapper, err := storage.NewDatabaseWrapper(storageConfig, logger)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbWrapper.Close()

	// Setup merkle service with unified implementation
	merkleService := merkleimpl.New(dbWrapper.GetDB(), subgraphClient, logger)

	// Setup services
	epochService := epochimpl.New(contractClient, subgraphClient, merkleService, logger, cfg)

	// Create a mock subsidy service for now
	subsidyService := &mockSubsidyService{logger: logger}

	// Setup scheduler with proper service interfaces
	schedulerInterval := time.Duration(cfg.Scheduler.Interval) * time.Second
	schedulerInstance := scheduler.NewScheduler(epochService, subsidyService, schedulerInterval, logger, cfg)
	go schedulerInstance.Start(ctx)

	// Setup and start HTTP server
	server := api.NewServer(epochService, subsidyService, merkleService, logger, cfg)

	if err := server.Start(); err != nil {
		logger.Logf("ERROR server failed to start: %v", err)
	}
}

// Mock subsidy service for now - this will be replaced with proper implementation

type mockSubsidyService struct {
	logger interface{}
}

func (m *mockSubsidyService) DistributeSubsidies(ctx context.Context, vaultId string) error {
	return fmt.Errorf("mock service not implemented")
}
