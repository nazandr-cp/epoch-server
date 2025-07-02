package graph

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
			users, err := client.QueryUsers(context.Background())

			if tt.wantErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if len(users) != tt.wantUsersCount {
				t.Errorf("Expected %d users, got %d", tt.wantUsersCount, len(users))
			}

			if !tt.wantErr && len(users) > 0 {
				user := users[0]
				if user.ID != "user1" {
					t.Errorf("Expected first user ID to be 'user1', got %s", user.ID)
				}
				if user.TotalSecondsClaimed != "100" {
					t.Errorf("Expected TotalSecondsClaimed to be '100', got %s", user.TotalSecondsClaimed)
				}
			}
		})
	}
}

func TestClient_QueryEligibility(t *testing.T) {
	tests := []struct {
		name                 string
		epochID              string
		serverResponse       string
		statusCode           int
		wantErr              bool
		wantEligibilityCount int
	}{
		{
			name:    "successful query",
			epochID: "epoch1",
			serverResponse: `{
				"data": {
					"userEpochEligibilities": [
						{
							"id": "eligibility1",
							"user": {
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
							"epoch": {
								"id": "epoch1",
								"epochNumber": "1",
								"status": "ACTIVE",
								"startTimestamp": "1640995200",
								"endTimestamp": "1640995300",
								"processingStartedTimestamp": "",
								"processingCompletedTimestamp": "",
								"totalYieldAvailable": "1000",
								"totalYieldAllocated": "800",
								"totalYieldDistributed": "0",
								"remainingYield": "200",
								"totalSubsidiesDistributed": "0",
								"totalEligibleUsers": "10",
								"totalParticipatingCollections": "3",
								"participantCount": "15",
								"processingTimeMs": "0",
								"estimatedProcessingTime": "300",
								"processingGasUsed": "0",
								"processingTransactionCount": "0",
								"createdAtBlock": "1000",
								"createdAtTimestamp": "1640995000",
								"updatedAtBlock": "2000",
								"updatedAtTimestamp": "1640995300"
							},
							"collection": {
								"id": "collection1",
								"contractAddress": "0x1234567890abcdef",
								"name": "Test Collection",
								"symbol": "TEST",
								"totalSupply": "10000",
								"collectionType": "ERC721",
								"isActive": true,
								"yieldSharePercentage": "10",
								"weightFunctionType": "LINEAR",
								"weightFunctionP1": "1",
								"weightFunctionP2": "0",
								"minBorrowAmount": "100",
								"maxBorrowAmount": "10000",
								"totalNFTsDeposited": "500",
								"registeredAtBlock": "500",
								"registeredAtTimestamp": "1640994000",
								"updatedAtBlock": "2000",
								"updatedAtTimestamp": "1640995300"
							},
							"nftBalance": "3",
							"borrowBalance": "1500",
							"holdingDuration": "86400",
							"isEligible": true,
							"subsidyReceived": "25",
							"yieldShare": "15",
							"bonusMultiplier": "1.2",
							"calculatedAtBlock": "2000",
							"calculatedAtTimestamp": "1640995300"
						}
					]
				}
			}`,
			statusCode:           http.StatusOK,
			wantErr:              false,
			wantEligibilityCount: 1,
		},
		{
			name:                 "server error",
			epochID:              "epoch1",
			serverResponse:       `{"error": "internal server error"}`,
			statusCode:           http.StatusInternalServerError,
			wantErr:              true,
			wantEligibilityCount: 0,
		},
		{
			name:    "graphql errors",
			epochID: "epoch1",
			serverResponse: `{
				"data": null,
				"errors": [
					{"message": "Field 'userEpochEligibilities' doesn't exist on type 'Query'"}
				]
			}`,
			statusCode:           http.StatusOK,
			wantErr:              true,
			wantEligibilityCount: 0,
		},
		{
			name:                 "malformed json",
			epochID:              "epoch1",
			serverResponse:       `{"data": {"userEpochEligibilities": [invalid json`,
			statusCode:           http.StatusOK,
			wantErr:              true,
			wantEligibilityCount: 0,
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

				if !strings.Contains(req.Query, "userEpochEligibilities") {
					t.Errorf("Expected userEpochEligibilities query, got %s", req.Query)
				}

				if req.Variables["epochId"] != tt.epochID {
					t.Errorf("Expected epochId variable to be %s, got %v", tt.epochID, req.Variables["epochId"])
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
			eligibilities, err := client.QueryEligibility(context.Background(), tt.epochID)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if len(eligibilities) != tt.wantEligibilityCount {
				t.Errorf("Expected %d eligibilities, got %d", tt.wantEligibilityCount, len(eligibilities))
			}

			if !tt.wantErr && len(eligibilities) > 0 {
				eligibility := eligibilities[0]
				if eligibility.ID != "eligibility1" {
					t.Errorf("Expected first eligibility ID to be 'eligibility1', got %s", eligibility.ID)
				}
				if eligibility.User.ID != "user1" {
					t.Errorf("Expected user ID to be 'user1', got %s", eligibility.User.ID)
				}
				if !eligibility.IsEligible {
					t.Errorf("Expected isEligible to be true, got %v", eligibility.IsEligible)
				}
			}
		})
	}
}
