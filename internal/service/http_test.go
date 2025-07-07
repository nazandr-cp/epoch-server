package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrey/epoch-server/internal/config"
	"github.com/go-pkgz/lgr"
)

func TestService_NewHTTPHandler(t *testing.T) {
	// Create a mock service
	mockGraphClient := &mockGraphClient{}
	mockContractClient := &mockContractClient{}
	logger := lgr.NoOp
	cfg := &config.Config{}
	cfg.Contracts.CollectionsVault = "0x1234567890123456789012345678901234567890"

	service := &Service{
		graphClient:    mockGraphClient,
		contractClient: mockContractClient,
		logger:         logger,
		config:         cfg,
	}

	// Get the HTTP handler
	handler := service.NewHTTPHandler()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "health check",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ping endpoint",
			method:         "GET",
			path:           "/ping",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "start epoch",
			method:         "POST",
			path:           "/epochs/start",
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "distribute subsidies",
			method:         "POST",
			path:           "/epochs/distribute",
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "method not allowed",
			method:         "DELETE",
			path:           "/epochs/start",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "not found",
			method:         "GET",
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}