// @title Epoch Server API
// @version 1.0
// @description Epoch Server for managing NFT collection-backed lending epochs, subsidies, and merkle proofs
// @termsOfService http://lend.fam/terms/
// @contact.name API Support
// @contact.url http://lend.fam/support
// @contact.email support@lend.fam
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8088
// @BasePath /
// @schemes http https
// @accept json
// @produce json
package main

import (
	"context"
	"log"

	"github.com/andrey/epoch-server/internal/api"
	"github.com/andrey/epoch-server/internal/infra/blockchain"
	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/infra/logging"
	"github.com/andrey/epoch-server/internal/infra/storage"
	"github.com/andrey/epoch-server/internal/infra/subgraph"
	blockchainService "github.com/andrey/epoch-server/internal/services/blockchain"
	"github.com/andrey/epoch-server/internal/services/epoch/epochimpl"
	"github.com/andrey/epoch-server/internal/services/merkle/merkleimpl"
	"github.com/andrey/epoch-server/internal/services/scheduler"
	storageService "github.com/andrey/epoch-server/internal/services/storage"
	subgraphService "github.com/andrey/epoch-server/internal/services/subgraph"
	"github.com/andrey/epoch-server/internal/services/subsidy/subsidyimpl"
	"github.com/go-pkgz/lgr"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger := setupLogging(cfg)

	ctx := context.Background()
	subgraphClient := setupSubgraphClient(cfg, logger, ctx)
	contractClient := setupBlockchainClient(cfg, logger)

	storageClient := setupDatabase(cfg, logger)
	defer func() {
		if closeErr := storageClient.Close(); closeErr != nil {
			logger.Logf("WARN failed to close database: %v", closeErr)
		}
	}()

	epochService, subsidyService, merkleService := setupServices(cfg, logger, contractClient, subgraphClient, storageClient)

	setupScheduler(cfg, logger, ctx, epochService, subsidyService)
	startServer(cfg, logger, epochService, subsidyService, merkleService)
}

func setupLogging(cfg *config.Config) lgr.L {
	return logging.NewWithConfig(logging.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
		Output: cfg.Logging.Output,
	})
}

func setupSubgraphClient(cfg *config.Config, logger lgr.L, ctx context.Context) subgraph.SubgraphClient {
	subgraphClient := subgraphService.ProvideClient(cfg.Subgraph.Endpoint, logger)

	logger.Logf("INFO checking subgraph connectivity at %s", cfg.Subgraph.Endpoint)
	if err := subgraphClient.HealthCheck(ctx); err != nil {
		log.Fatalf("Failed to connect to subgraph: %v", err)
	}
	logger.Logf("INFO subgraph health check passed")

	return subgraphClient
}

func setupBlockchainClient(cfg *config.Config, logger lgr.L) blockchain.BlockchainClient {
	contractClient, err := blockchainService.ProvideClientWithConfig(logger, blockchain.Config{
		RPCURL:             cfg.Ethereum.RPCURL,
		PrivateKey:         cfg.Ethereum.PrivateKey,
		GasLimit:           cfg.Ethereum.GasLimit,
		GasPrice:           cfg.Ethereum.GasPrice,
		Comptroller:        cfg.Contracts.Comptroller,
		EpochManager:       cfg.Contracts.EpochManager,
		DebtSubsidizer:     cfg.Contracts.DebtSubsidizer,
		LendingManager:     cfg.Contracts.LendingManager,
		CollectionRegistry: cfg.Contracts.CollectionRegistry,
	})
	if err != nil {
		log.Fatalf("Failed to initialize contract client: %v", err)
	}

	return contractClient
}

func setupDatabase(cfg *config.Config, logger lgr.L) storage.StorageClient {
	storageClient, err := storageService.ProvideClient(storage.Config{
		Type: cfg.Database.Type,
		Path: cfg.Database.ConnectionString,
	}, logger)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	return storageClient
}

func setupServices(
	cfg *config.Config,
	logger lgr.L,
	contractClient blockchain.BlockchainClient,
	subgraphClient subgraph.SubgraphClient,
	storageClient storage.StorageClient,
) (*epochimpl.Service, *subsidyimpl.Service, *merkleimpl.Service) {
	merkleService := merkleimpl.New(storageClient.GetDB(), subgraphClient, logger)
	epochService := epochimpl.New(contractClient, subgraphClient, merkleService, logger, cfg)
	lazyDistributor := subsidyimpl.NewLazyDistributor(contractClient, merkleService, subgraphClient, logger)
	subsidyService := subsidyimpl.New(lazyDistributor, epochService, logger, cfg)

	return epochService, subsidyService, merkleService
}

func setupScheduler(
	cfg *config.Config,
	logger lgr.L,
	ctx context.Context,
	epochService *epochimpl.Service,
	subsidyService *subsidyimpl.Service,
) {
	schedulerInstance := scheduler.NewScheduler(epochService, subsidyService, cfg.Scheduler.Interval, logger, cfg)
	go schedulerInstance.Start(ctx)
}

func startServer(
	cfg *config.Config,
	logger lgr.L,
	epochService *epochimpl.Service,
	subsidyService *subsidyimpl.Service,
	merkleService *merkleimpl.Service,
) {
	server := api.NewServer(epochService, subsidyService, merkleService, logger, cfg)

	if err := server.Start(); err != nil {
		logger.Logf("ERROR server failed to start: %v", err)
	}
}
