package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-pkgz/lgr"
)

type mockService struct {
	startEpochFunc          func(ctx context.Context, epochID string) error
	distributeSubsidiesFunc func(ctx context.Context, vaultID string) error
}

func (m *mockService) StartEpoch(ctx context.Context, epochID string) error {
	if m.startEpochFunc != nil {
		return m.startEpochFunc(ctx, epochID)
	}
	return nil
}

func (m *mockService) DistributeSubsidies(ctx context.Context, vaultID string) error {
	if m.distributeSubsidiesFunc != nil {
		return m.distributeSubsidiesFunc(ctx, vaultID)
	}
	return nil
}

func TestHandler_Health(t *testing.T) {
	handler := NewHandler(nil, lgr.NoOp)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	handler.Health(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	expectedContentType := "application/json"
	if contentType := recorder.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
	}

	var response map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	expectedStatus := "ok"
	if status := response["status"]; status != expectedStatus {
		t.Errorf("Expected status %s, got %s", expectedStatus, status)
	}
}

func TestHandler_StartEpoch(t *testing.T) {
	tests := []struct {
		name               string
		epochID            string
		mockService        *mockService
		expectedStatusCode int
		wantServiceCall    bool
	}{
		{
			name:    "successful start epoch",
			epochID: "epoch1",
			mockService: &mockService{
				startEpochFunc: func(ctx context.Context, epochID string) error {
					return nil
				},
			},
			expectedStatusCode: http.StatusAccepted,
			wantServiceCall:    true,
		},
		{
			name:    "service error",
			epochID: "epoch1",
			mockService: &mockService{
				startEpochFunc: func(ctx context.Context, epochID string) error {
					return errors.New("service error")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			wantServiceCall:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var serviceCalled bool
			if tt.mockService.startEpochFunc != nil {
				originalFunc := tt.mockService.startEpochFunc
				tt.mockService.startEpochFunc = func(ctx context.Context, epochID string) error {
					serviceCalled = true
					if epochID != tt.epochID {
						t.Errorf("Expected epochID %s, got %s", tt.epochID, epochID)
					}
					return originalFunc(ctx, epochID)
				}
			}

			handler := &Handler{
				service: tt.mockService,
				logger:  lgr.NoOp,
			}

			router := chi.NewRouter()
			router.Post("/epochs/{id}/start", handler.StartEpoch)

			req := httptest.NewRequest(http.MethodPost, "/epochs/"+tt.epochID+"/start", nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, recorder.Code)
			}

			if tt.wantServiceCall && !serviceCalled {
				t.Errorf("Expected service to be called, but it wasn't")
			}
			if !tt.wantServiceCall && serviceCalled {
				t.Errorf("Expected service not to be called, but it was")
			}
		})
	}
}

func TestHandler_DistributeSubsidies(t *testing.T) {
	tests := []struct {
		name               string
		epochID            string
		mockService        *mockService
		expectedStatusCode int
		wantServiceCall    bool
	}{
		{
			name:    "successful distribute subsidies",
			epochID: "epoch1",
			mockService: &mockService{
				distributeSubsidiesFunc: func(ctx context.Context, vaultID string) error {
					return nil
				},
			},
			expectedStatusCode: http.StatusAccepted,
			wantServiceCall:    true,
		},
		{
			name:    "service error",
			epochID: "epoch1",
			mockService: &mockService{
				distributeSubsidiesFunc: func(ctx context.Context, vaultID string) error {
					return errors.New("service error")
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			wantServiceCall:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var serviceCalled bool
			if tt.mockService.distributeSubsidiesFunc != nil {
				originalFunc := tt.mockService.distributeSubsidiesFunc
				tt.mockService.distributeSubsidiesFunc = func(ctx context.Context, vaultID string) error {
					serviceCalled = true
					expectedVaultID := "0x4a4be724f522946296a51d8c82c7c2e8e5a62655"
					if vaultID != expectedVaultID {
						t.Errorf("Expected vaultID %s, got %s", expectedVaultID, vaultID)
					}
					return originalFunc(ctx, vaultID)
				}
			}

			handler := &Handler{
				service: tt.mockService,
				logger:  lgr.NoOp,
			}

			router := chi.NewRouter()
			router.Post("/epochs/{id}/distribute", handler.DistributeSubsidies)

			req := httptest.NewRequest(http.MethodPost, "/epochs/"+tt.epochID+"/distribute", nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, recorder.Code)
			}

			if tt.wantServiceCall && !serviceCalled {
				t.Errorf("Expected service to be called, but it wasn't")
			}
			if !tt.wantServiceCall && serviceCalled {
				t.Errorf("Expected service not to be called, but it was")
			}
		})
	}
}
