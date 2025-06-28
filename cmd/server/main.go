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

	graphClient := graph.NewClient("") // TODO: Get endpoint from config
	contractClient := contract.NewClient(logger)

	svc := service.NewService(graphClient, contractClient, logger)
	handler := handlers.NewHandler(svc, logger)

	ctx := context.Background()
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
