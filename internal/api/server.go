package api

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/andrey/epoch-server/docs"
	"github.com/andrey/epoch-server/internal/api/handlers"
	"github.com/andrey/epoch-server/internal/api/middleware"
	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/go-pkgz/lgr"
	"github.com/go-pkgz/rest"
	"github.com/go-pkgz/routegroup"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Server represents the HTTP server
type Server struct {
	epochService   epoch.Service
	subsidyService subsidy.Service
	merkleService  merkle.Service
	logger         lgr.L
	config         *config.Config
}

// NewServer creates a new HTTP server
func NewServer(
	epochService epoch.Service,
	subsidyService subsidy.Service,
	merkleService merkle.Service,
	logger lgr.L,
	cfg *config.Config,
) *Server {
	return &Server{
		epochService:   epochService,
		subsidyService: subsidyService,
		merkleService:  merkleService,
		logger:         logger,
		config:         cfg,
	}
}

// SetupRoutes configures all HTTP routes and middleware
func (s *Server) SetupRoutes() http.Handler {
	// Create handlers
	healthHandler := handlers.NewHealthHandler(s.logger, s.checkEpochService, s.checkSubsidyService, s.checkMerkleService)
	epochHandler := handlers.NewEpochHandler(s.epochService, s.logger, s.config)
	subsidyHandler := handlers.NewSubsidyHandler(s.subsidyService, s.logger, s.config)
	merkleHandler := handlers.NewMerkleHandler(s.merkleService, s.logger, s.config)

	// Create base router with routegroup
	router := routegroup.New(http.NewServeMux())

	// Apply global middlewares
	router.Use(rest.RealIP)
	router.Use(rest.Trace)                  // Add request tracing
	router.Use(rest.SizeLimit(1024 * 1024)) // 1MB request size limit
	// router.Use(middleware.Auth(s.logger))
	router.Use(middleware.Logging(s.logger)) // Keep custom logging middleware
	router.Use(middleware.Recovery(s.logger))
	router.Use(rest.AppInfo("epoch-server", "andrey", "1.0.0"))
	router.Use(rest.Ping)

	// Health check route (no grouping needed)
	router.HandleFunc("GET /health", healthHandler.HandleHealth)

	// Swagger documentation route
	router.HandleFunc("GET /swagger/*", httpSwagger.Handler())

	// API routes group
	router.Group().Mount("/api").Route(func(apiRouter *routegroup.Bundle) {
		// Epoch management routes
		apiRouter.Group().Mount("/epochs").Route(func(epochRouter *routegroup.Bundle) {
			epochRouter.HandleFunc("POST /start", epochHandler.HandleStartEpoch)
			epochRouter.HandleFunc("POST /force-end", epochHandler.HandleForceEndEpoch)
			epochRouter.HandleFunc("POST /distribute", subsidyHandler.HandleDistributeSubsidies)
		})

		// User-related routes
		apiRouter.Group().Mount("/users").Route(func(userRouter *routegroup.Bundle) {
			userRouter.HandleFunc("GET /{address}/total-earned", epochHandler.HandleGetUserTotalEarned)
			userRouter.HandleFunc("GET /{address}/merkle-proof", merkleHandler.HandleGetUserMerkleProof)
			userRouter.HandleFunc(
				"GET /{address}/merkle-proof/epoch/{epochNumber}",
				merkleHandler.HandleGetUserHistoricalMerkleProof,
			)
		})
	})

	return router
}

// Start starts the HTTP server with proper timeouts
func (s *Server) Start() error {
	handler := s.SetupRoutes()
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.logger.Logf("INFO starting server on %s", addr)

	// Create server with security timeouts
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server.ListenAndServe()
}

// Health check functions for services
func (s *Server) checkEpochService() error {
	// Basic health check - could be enhanced with actual service checks
	if s.epochService == nil {
		return fmt.Errorf("epoch service not initialized")
	}
	return nil
}

func (s *Server) checkSubsidyService() error {
	// Basic health check - could be enhanced with actual service checks
	if s.subsidyService == nil {
		return fmt.Errorf("subsidy service not initialized")
	}
	return nil
}

func (s *Server) checkMerkleService() error {
	// Basic health check - could be enhanced with actual service checks
	if s.merkleService == nil {
		return fmt.Errorf("merkle service not initialized")
	}
	return nil
}
