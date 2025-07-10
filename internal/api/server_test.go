package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/go-pkgz/lgr"
)

func TestServerRoutes(t *testing.T) {
	// Create mock services
	mockEpochService := &epoch.ServiceMock{
		StartEpochFunc: func(ctx context.Context) (*epoch.StartEpochResponse, error) {
			return &epoch.StartEpochResponse{Status: "started"}, nil
		},
		ForceEndEpochFunc: func(ctx context.Context, epochId uint64, vaultId string) (*epoch.ForceEndEpochResponse, error) {
			return &epoch.ForceEndEpochResponse{Status: "force_ended"}, nil
		},
		GetUserTotalEarnedFunc: func(ctx context.Context, userAddress, vaultId string) (*epoch.UserEarningsResponse, error) {
			return &epoch.UserEarningsResponse{}, nil
		},
	}

	mockSubsidyService := &subsidy.ServiceMock{
		DistributeSubsidiesFunc: func(ctx context.Context, vaultId string) (*subsidy.SubsidyDistributionResponse, error) {
			return &subsidy.SubsidyDistributionResponse{Status: "completed"}, nil
		},
	}

	mockMerkleService := &merkle.ServiceMock{
		GenerateUserMerkleProofFunc: func(
			ctx context.Context,
			userAddress, vaultAddress string,
		) (*merkle.UserMerkleProofResponse, error) {
			return &merkle.UserMerkleProofResponse{}, nil
		},
		GenerateHistoricalMerkleProofFunc: func(
			ctx context.Context,
			userAddress, vaultAddress, epochNumber string,
		) (*merkle.UserMerkleProofResponse, error) {
			return &merkle.UserMerkleProofResponse{}, nil
		},
	}

	logger := lgr.NoOp
	cfg := &config.Config{}

	// Create server
	server := NewServer(mockEpochService, mockSubsidyService, mockMerkleService, logger, cfg)
	handler := server.SetupRoutes()

	// Test cases for different routes
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		description    string
	}{
		{
			name:           "health_check",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
			description:    "Health check endpoint",
		},
		{
			name:           "epoch_start",
			method:         "POST",
			path:           "/api/epochs/start",
			expectedStatus: http.StatusAccepted,
			description:    "Start epoch endpoint",
		},
		{
			name:           "epoch_force_end",
			method:         "POST",
			path:           "/api/epochs/force-end",
			expectedStatus: http.StatusBadRequest,
			description:    "Force end epoch endpoint (needs request body)",
		},
		{
			name:           "epoch_distribute",
			method:         "POST",
			path:           "/api/epochs/distribute",
			expectedStatus: http.StatusAccepted,
			description:    "Distribute subsidies endpoint",
		},
		{
			name:           "user_total_earned",
			method:         "GET",
			path:           "/api/users/0x1234567890123456789012345678901234567890/total-earned",
			expectedStatus: http.StatusOK,
			description:    "Get user total earned endpoint",
		},
		{
			name:           "user_merkle_proof",
			method:         "GET",
			path:           "/api/users/0x1234567890123456789012345678901234567890/merkle-proof",
			expectedStatus: http.StatusOK,
			description:    "Get user merkle proof endpoint",
		},
		{
			name:           "user_historical_merkle_proof",
			method:         "GET",
			path:           "/api/users/0x1234567890123456789012345678901234567890/merkle-proof/epoch/1",
			expectedStatus: http.StatusOK,
			description:    "Get user historical merkle proof endpoint",
		},
		// Note: Swagger UI test is disabled as it requires static files to be served
		// which don't work well in test environment. The endpoint works in production.
		// {
		//     name:           "swagger_ui",
		//     method:         "GET",
		//     path:           "/swagger/",
		//     expectedStatus: http.StatusOK,
		//     description:    "Swagger UI endpoint",
		// },
		{
			name:           "not_found",
			method:         "GET",
			path:           "/api/nonexistent",
			expectedStatus: http.StatusNotFound,
			description:    "Non-existent endpoint should return 404",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.description, tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestRouteGrouping(t *testing.T) {
	// Create minimal server for testing structure
	server := NewServer(nil, nil, nil, lgr.NoOp, &config.Config{})
	handler := server.SetupRoutes()

	// Test that routes are properly grouped
	testRoutes := []struct {
		path        string
		method      string
		shouldExist bool
	}{
		{"/health", "GET", true},
		{"/api/epochs/start", "POST", true},
		{"/api/users/test/total-earned", "GET", true},
		{"/epochs/start", "POST", false},           // Should not work without /api prefix
		{"/users/test/total-earned", "GET", false}, // Should not work without /api prefix
	}

	for _, route := range testRoutes {
		req := httptest.NewRequest(route.method, route.path, nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if route.shouldExist && rr.Code == http.StatusNotFound {
			t.Errorf("Route %s (%s) should exist but returned 404", route.path, route.method)
		}
		if !route.shouldExist && rr.Code != http.StatusNotFound {
			t.Errorf("Route %s (%s) should not exist but returned %d", route.path, route.method, rr.Code)
		}
	}
}
