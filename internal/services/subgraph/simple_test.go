package subgraph

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/go-pkgz/lgr"
)

func TestClient_DirectQuery(t *testing.T) {
	// Test direct query without pagination
	serverResponse := `{
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
				}
			]
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(serverResponse))
	}))
	defer server.Close()

	client := ProvideClient(server.URL, lgr.NoOp)

	// Test ExecuteQuery directly
	request := subgraph.GraphQLRequest{
		Query: `query { accounts { id totalSecondsClaimed } }`,
	}

	var response subgraph.AccountsResponse

	err := client.ExecuteQuery(context.Background(), request, &response)
	if err != nil {
		t.Fatalf("ExecuteQuery failed: %v", err)
	}

	if len(response.Accounts) != 1 {
		t.Errorf("Expected 1 account, got %d", len(response.Accounts))
	}

	if len(response.Accounts) > 0 && response.Accounts[0].ID != "user1" {
		t.Errorf("Expected user1, got %s", response.Accounts[0].ID)
	}
}
