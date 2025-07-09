package subgraph

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_QueryUsers(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		statusCode     int
		wantErr        bool
		wantUsersCount int
	}{
		{
			name: "successful query",
			serverResponse: `{
				"data": {
					"accounts": [
						{
							"id": "user1",
							"totalSecondsClaimed": "100",
							"totalSubsidiesReceived": "50",
							"totalYieldEarned": "25",
							"totalBorrowVolume": "1000",
							"totalNFTsOwned": "5",
							"totalCollectionsParticipated": "2",
							"createdAtBlock": "1000",
							"createdAtTimestamp": "1640995200",
							"updatedAtBlock": "2000",
							"updatedAtTimestamp": "1640995300"
						},
						{
							"id": "user2",
							"totalSecondsClaimed": "200",
							"totalSubsidiesReceived": "75",
							"totalYieldEarned": "35",
							"totalBorrowVolume": "2000",
							"totalNFTsOwned": "10",
							"totalCollectionsParticipated": "3",
							"createdAtBlock": "1100",
							"createdAtTimestamp": "1640995400",
							"updatedAtBlock": "2100",
							"updatedAtTimestamp": "1640995500"
						}
					]
				}
			}`,
			statusCode:     http.StatusOK,
			wantErr:        false,
			wantUsersCount: 2,
		},
		{
			name:           "server error",
			serverResponse: `{"error": "internal server error"}`,
			statusCode:     http.StatusInternalServerError,
			wantErr:        true,
			wantUsersCount: 0,
		},
		{
			name: "graphql errors",
			serverResponse: `{
				"data": null,
				"errors": [
					{"message": "Field 'accounts' doesn't exist on type 'Query'"}
				]
			}`,
			statusCode:     http.StatusOK,
			wantErr:        true,
			wantUsersCount: 0,
		},
		{
			name:           "malformed json",
			serverResponse: `{"data": {"accounts": [invalid json`,
			statusCode:     http.StatusOK,
			wantErr:        true,
			wantUsersCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
				}

				var req GraphQLRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("Failed to decode request: %v", err)
				}

				// QueryUsers now delegates to QueryAccounts, so we expect accounts query
				if !strings.Contains(req.Query, "accounts(first:") {
					t.Errorf("Expected accounts query (via QueryUsers delegation), got %s", req.Query)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.serverResponse))
				if err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			}))
			defer server.Close()

			client := NewClient(server.URL)
			accounts, err := client.QueryAccounts(context.Background())

			if tt.wantErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if len(accounts) != tt.wantUsersCount {
				t.Errorf("Expected %d accounts, got %d (error: %v)", tt.wantUsersCount, len(accounts), err)
			}

			if !tt.wantErr && len(accounts) > 0 {
				account := accounts[0]
				if account.ID != "user1" {
					t.Errorf("Expected first account ID to be 'user1', got %s", account.ID)
				}
				if account.TotalSecondsClaimed != "100" {
					t.Errorf("Expected TotalSecondsClaimed to be '100', got %s", account.TotalSecondsClaimed)
				}
			}
		})
	}
}

func TestClient_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		statusCode     int
		wantErr        bool
	}{
		{
			name: "successful health check",
			serverResponse: `{
				"data": {
					"__schema": {
						"queryType": {
							"name": "Query"
						}
					}
				}
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:           "server error",
			serverResponse: `{"error": "internal server error"}`,
			statusCode:     http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name: "graphql errors",
			serverResponse: `{
				"data": null,
				"errors": [
					{"message": "Schema introspection not allowed"}
				]
			}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:           "malformed json",
			serverResponse: `{"data": {"__schema": invalid json`,
			statusCode:     http.StatusOK,
			wantErr:        true,
		},
		{
			name: "unexpected response structure",
			serverResponse: `{
				"data": {
					"__schema": {
						"queryType": {
							"name": "Mutation"
						}
					}
				}
			}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
				}

				var req GraphQLRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("Failed to decode request: %v", err)
				}

				if !strings.Contains(req.Query, "__schema") {
					t.Errorf("Expected __schema introspection query, got %s", req.Query)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.serverResponse))
				if err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			}))
			defer server.Close()

			client := NewClient(server.URL)
			err := client.HealthCheck(context.Background())

			if tt.wantErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

