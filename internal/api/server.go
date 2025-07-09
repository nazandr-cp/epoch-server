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

	// Create new ServeMux with Go 1.22+ routing patterns
	mux := http.NewServeMux()

	// Register health check route
	mux.HandleFunc("GET /health", healthHandler.HandleHealth)

	// Register epoch routes
	mux.HandleFunc("POST /epochs/start", epochHandler.HandleStartEpoch)
	mux.HandleFunc("POST /epochs/force-end", epochHandler.HandleForceEndEpoch)
	mux.HandleFunc("GET /users/{address}/total-earned", epochHandler.HandleGetUserTotalEarned)

	// Register subsidy routes
	mux.HandleFunc("POST /epochs/distribute", subsidyHandler.HandleDistributeSubsidies)

	// Register merkle proof routes
	mux.HandleFunc("GET /users/{address}/merkle-proof", merkleHandler.HandleGetUserMerkleProof)
	mux.HandleFunc("GET /users/{address}/merkle-proof/epoch/{epochNumber}", merkleHandler.HandleGetUserHistoricalMerkleProof)

	// Apply middlewares using go-pkgz/rest and custom middleware
	var handler http.Handler = mux

	// Apply middlewares in reverse order (last applied = outermost)
	handler = rest.Ping(handler)
	handler = rest.AppInfo("epoch-server", "andrey", "1.0.0")(handler)
	handler = middleware.Recovery(s.logger)(handler)
	handler = middleware.Logging(s.logger)(handler)
	handler = middleware.Auth(s.logger)(handler)
	handler = rest.RealIP(handler)

	return handler
}

// Start starts the HTTP server
func (s *Server) Start() error {
	handler := s.SetupRoutes()
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.logger.Logf("INFO starting server on %s", addr)

	return http.ListenAndServe(addr, handler)
}
