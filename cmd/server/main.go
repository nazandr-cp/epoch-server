package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/andrey/epoch-server/internal/clients/contract"
	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/andrey/epoch-server/internal/config"
	"github.com/andrey/epoch-server/internal/handlers"
	internalLog "github.com/andrey/epoch-server/internal/log"
	"github.com/andrey/epoch-server/internal/scheduler"
	"github.com/andrey/epoch-server/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger := internalLog.New(cfg.Logging.Level)

	graphClient := graph.NewClient(cfg.Subgraph.Endpoint)

	// Perform subgraph health check during startup
	logger.Logf("INFO checking subgraph connectivity at %s", cfg.Subgraph.Endpoint)
	ctx := context.Background()
	if err := graphClient.HealthCheck(ctx); err != nil {
		log.Fatalf("Failed to connect to subgraph: %v", err)
	}
	logger.Logf("INFO subgraph health check passed")

	ethConfig := contract.EthereumConfig{
		RPCURL:     cfg.Ethereum.RPCURL,
		PrivateKey: cfg.Ethereum.PrivateKey,
		GasLimit:   cfg.Ethereum.GasLimit,
		GasPrice:   cfg.Ethereum.GasPrice,
	}

	contractAddresses := contract.ContractAddresses{
		Comptroller:        cfg.Contracts.Comptroller,
		EpochManager:       cfg.Contracts.EpochManager,
		DebtSubsidizer:     cfg.Contracts.DebtSubsidizer,
		LendingManager:     cfg.Contracts.LendingManager,
		CollectionRegistry: cfg.Contracts.CollectionRegistry,
	}

	contractClient, err := contract.NewClientWithConfig(logger, ethConfig, contractAddresses)
	if err != nil {
		log.Fatalf("Failed to initialize contract client: %v", err)
	}

	svc := service.NewService(graphClient, contractClient, logger)
	handler := handlers.NewHandler(svc, logger)
	schedulerInstance := scheduler.NewScheduler(cfg.Scheduler.Interval, svc, logger)
	go schedulerInstance.Start(ctx)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Get("/health", handler.Health)
	r.Post("/epochs/{id}/start", handler.StartEpoch)
	r.Post("/epochs/{id}/distribute", handler.DistributeSubsidies)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Logf("INFO starting server on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Logf("ERROR server failed to start: %v", err)
	}
}
