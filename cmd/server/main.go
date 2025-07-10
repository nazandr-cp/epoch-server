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

	"github.com/andrey/epoch-server/internal/api"
	"github.com/andrey/epoch-server/internal/infra/blockchain"
	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/infra/logging"
	"github.com/andrey/epoch-server/internal/infra/storage"
	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/andrey/epoch-server/internal/services/epoch/epochimpl"
	"github.com/andrey/epoch-server/internal/services/merkle/merkleimpl"
	"github.com/andrey/epoch-server/internal/services/scheduler"
	"github.com/go-pkgz/lgr"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logging
	logger := setupLogging(cfg)

	// Setup infrastructure components
	ctx := context.Background()
	subgraphClient := setupSubgraphClient(cfg, logger, ctx)

	// Setup blockchain clients
	contractClient := setupBlockchainClient(cfg, logger)

	// Setup database
	dbWrapper := setupDatabase(cfg, logger)
	defer func() {
		if closeErr := dbWrapper.Close(); closeErr != nil {
			logger.Logf("WARN failed to close database: %v", closeErr)
		}
	}()

	// Setup services
	epochService, subsidyService, merkleService := setupServices(cfg, logger, contractClient, subgraphClient, dbWrapper)

	// Setup scheduler
	setupScheduler(cfg, logger, ctx, epochService, subsidyService)

	// Setup and start HTTP server
	startServer(cfg, logger, epochService, subsidyService, merkleService)
}

// setupLogging configures the logger with the provided configuration
func setupLogging(cfg *config.Config) lgr.L {
	logConfig := logging.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
		Output: cfg.Logging.Output,
	}
	return logging.NewWithConfig(logConfig)
}

// setupSubgraphClient creates and health-checks the subgraph client
func setupSubgraphClient(cfg *config.Config, logger lgr.L, ctx context.Context) *subgraph.Client {
	subgraphClient := subgraph.NewClient(cfg.Subgraph.Endpoint, logger)

	// Perform subgraph health check during startup
	logger.Logf("INFO checking subgraph connectivity at %s", cfg.Subgraph.Endpoint)
	if err := subgraphClient.HealthCheck(ctx); err != nil {
		log.Fatalf("Failed to connect to subgraph: %v", err)
	}
	logger.Logf("INFO subgraph health check passed")

	return subgraphClient
}

// setupBlockchainClient creates and configures the blockchain client
func setupBlockchainClient(cfg *config.Config, logger lgr.L) *blockchain.Client {
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

	return contractClient
}

// setupDatabase creates and configures the database wrapper
func setupDatabase(cfg *config.Config, logger lgr.L) *storage.DatabaseWrapper {
	storageConfig := storage.StorageConfig{
		Type: cfg.Database.Type,
		Path: cfg.Database.ConnectionString,
	}
	dbWrapper, err := storage.NewDatabaseWrapper(storageConfig, logger)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	return dbWrapper
}

// setupServices creates and configures all business services
func setupServices(
	cfg *config.Config,
	logger lgr.L,
	contractClient *blockchain.Client,
	subgraphClient *subgraph.Client,
	dbWrapper *storage.DatabaseWrapper,
) (*epochimpl.Service, *mockSubsidyService, *merkleimpl.Service) {
	// Setup merkle service with unified implementation
	merkleService := merkleimpl.New(dbWrapper.GetDB(), subgraphClient, logger)

	// Setup services
	epochService := epochimpl.New(contractClient, subgraphClient, merkleService, logger, cfg)

	// Create a mock subsidy service for now
	subsidyService := &mockSubsidyService{logger: logger}

	return epochService, subsidyService, merkleService
}

// setupScheduler creates and starts the scheduler
func setupScheduler(
	cfg *config.Config,
	logger lgr.L,
	ctx context.Context,
	epochService *epochimpl.Service,
	subsidyService *mockSubsidyService,
) {
	schedulerInterval := cfg.Scheduler.Interval
	schedulerInstance := scheduler.NewScheduler(epochService, subsidyService, schedulerInterval, logger, cfg)
	go schedulerInstance.Start(ctx)
}

// startServer creates and starts the HTTP server
func startServer(
	cfg *config.Config,
	logger lgr.L,
	epochService *epochimpl.Service,
	subsidyService *mockSubsidyService,
	merkleService *merkleimpl.Service,
) {
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
