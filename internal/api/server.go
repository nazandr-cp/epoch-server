package api

import (
	"fmt"
	"net/http"

	"github.com/andrey/epoch-server/internal/api/handlers"
	"github.com/andrey/epoch-server/internal/api/middleware"
	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/go-pkgz/lgr"
	"github.com/go-pkgz/rest"
	"github.com/go-pkgz/routegroup"
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
	healthHandler := handlers.NewHealthHandler(s.logger)
	epochHandler := handlers.NewEpochHandler(s.epochService, s.logger, s.config)
	subsidyHandler := handlers.NewSubsidyHandler(s.subsidyService, s.logger, s.config)
	merkleHandler := handlers.NewMerkleHandler(s.merkleService, s.logger, s.config)

	// Create base router with routegroup
	router := routegroup.New(http.NewServeMux())

	// Apply global middlewares
	router.Use(rest.RealIP)
	router.Use(middleware.Auth(s.logger))
	router.Use(middleware.Logging(s.logger))
	router.Use(middleware.Recovery(s.logger))
	router.Use(rest.AppInfo("epoch-server", "andrey", "1.0.0"))
	router.Use(rest.Ping)

	// Health check route (no grouping needed)
	router.HandleFunc("GET /health", healthHandler.HandleHealth)

	// API routes group
	apiRouter := router.Mount("/api")
	
	// Epoch management routes
	epochRouter := apiRouter.Mount("/epochs")
	epochRouter.HandleFunc("POST /start", epochHandler.HandleStartEpoch)
	epochRouter.HandleFunc("POST /force-end", epochHandler.HandleForceEndEpoch)
	epochRouter.HandleFunc("POST /distribute", subsidyHandler.HandleDistributeSubsidies)

	// User-related routes
	userRouter := apiRouter.Mount("/users")
	userRouter.HandleFunc("GET /{address}/total-earned", epochHandler.HandleGetUserTotalEarned)
	userRouter.HandleFunc("GET /{address}/merkle-proof", merkleHandler.HandleGetUserMerkleProof)
	userRouter.HandleFunc("GET /{address}/merkle-proof/epoch/{epochNumber}", merkleHandler.HandleGetUserHistoricalMerkleProof)

	return router
}

// Start starts the HTTP server
func (s *Server) Start() error {
	handler := s.SetupRoutes()
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.logger.Logf("INFO starting server on %s", addr)

	return http.ListenAndServe(addr, handler)
}
