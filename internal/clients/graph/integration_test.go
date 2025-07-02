package graph

import (
	"context"
	"net/http"
	"testing"
	"time"
)

const realSubgraphEndpoint = "https://subgraph.satsuma-prod.com/63265bbf8342/analog-renaissances-team--450535/subsidiz/version/v0.0.4/api"

func TestClient_RealSubgraph_QueryAccounts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create client with longer timeout for integration tests
	client := &Client{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		endpoint: realSubgraphEndpoint,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	accounts, err := client.QueryAccounts(ctx)
	if err != nil {
		t.Fatalf("QueryAccounts failed: %v", err)
	}

	t.Logf("Retrieved %d accounts from real subgraph", len(accounts))

	if len(accounts) > 0 {
		account := accounts[0]
		t.Logf("First account: ID=%s, TotalSecondsClaimed=%s, TotalSubsidiesReceived=%s",
			account.ID, account.TotalSecondsClaimed, account.TotalSubsidiesReceived)

		// Validate account structure
		if account.ID == "" {
			t.Error("Account ID should not be empty")
		}
		if account.CreatedAtTimestamp == "" {
			t.Error("CreatedAtTimestamp should not be empty")
		}
		if account.UpdatedAtTimestamp == "" {
			t.Error("UpdatedAtTimestamp should not be empty")
		}
	}
}

func TestClient_RealSubgraph_QueryAccounts_Updated(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create client with longer timeout for integration tests
	client := &Client{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		endpoint: realSubgraphEndpoint,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	accounts, err := client.QueryAccounts(ctx)
	if err != nil {
		t.Fatalf("QueryAccounts failed: %v", err)
	}

	t.Logf("Retrieved %d accounts", len(accounts))

	if len(accounts) > 0 {
		account := accounts[0]
		t.Logf("First account: ID=%s, TotalSecondsClaimed=%s", account.ID, account.TotalSecondsClaimed)
	}
}

func TestClient_RealSubgraph_DirectQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create client with longer timeout for integration tests
	client := &Client{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		endpoint: realSubgraphEndpoint,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test direct query execution with minimal fields first
	request := GraphQLRequest{
		Query: `query { 
			accounts(first: 2) { 
				id 
				totalSecondsClaimed
			} 
		}`,
	}

	var response AccountsResponse

	err := client.ExecuteQuery(ctx, request, &response)
	if err != nil {
		t.Fatalf("ExecuteQuery failed: %v", err)
	}

	t.Logf("Direct query retrieved %d accounts", len(response.Accounts))

	if len(response.Accounts) > 0 {
		account := response.Accounts[0]
		t.Logf("Account data: %+v", account)
	}
}

func TestClient_RealSubgraph_AccountSubsidies(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create client with longer timeout for integration tests
	client := &Client{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		endpoint: realSubgraphEndpoint,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test accountSubsidies query for epoch server compatibility
	request := GraphQLRequest{
		Query: `query { 
			accountSubsidies(first: 5) { 
				id
				account {
					id
				}
				secondsAccumulated
				secondsClaimed
				lastEffectiveValue
				updatedAtTimestamp
			} 
		}`,
	}

	var response struct {
		AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
	}

	err := client.ExecuteQuery(ctx, request, &response)
	if err != nil {
		t.Fatalf("AccountSubsidies query failed: %v", err)
	}

	t.Logf("Retrieved %d account subsidies", len(response.AccountSubsidies))

	if len(response.AccountSubsidies) > 0 {
		subsidy := response.AccountSubsidies[0]
		t.Logf("First subsidy: ID=%s, SecondsAccumulated=%s, SecondsClaimed=%s",
			subsidy.ID, subsidy.SecondsAccumulated, subsidy.SecondsClaimed)

		// Validate subsidy structure
		if subsidy.Account.ID == "" {
			t.Error("Account ID in subsidy should not be empty")
		}
		if subsidy.SecondsAccumulated == "" {
			t.Error("SecondsAccumulated should not be empty")
		}
	}
}

func TestClient_RealSubgraph_PaginationStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create client with longer timeout for integration tests
	client := &Client{
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		endpoint: realSubgraphEndpoint,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Test pagination with accountSubsidies
	query := `
		query GetAccountSubsidies($first: Int!, $skip: Int!) {
			accountSubsidies(first: $first, skip: $skip) {
				id
				account {
					id
				}
				secondsAccumulated
				secondsClaimed
				lastEffectiveValue
				updatedAtTimestamp
			}
		}
	`

	var response struct {
		AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
	}

	err := client.ExecutePaginatedQuery(ctx, query, map[string]interface{}{}, "accountSubsidies", &response)
	if err != nil {
		t.Fatalf("Paginated query failed: %v", err)
	}

	t.Logf("Paginated query retrieved %d total account subsidies", len(response.AccountSubsidies))

	// Verify we got some data
	if len(response.AccountSubsidies) > 0 {
		t.Logf("Successfully retrieved dataset via pagination: %d records", len(response.AccountSubsidies))

		// Log first few records
		for i, subsidy := range response.AccountSubsidies[:min(3, len(response.AccountSubsidies))] {
			t.Logf("Record %d: ID=%s, Account=%s, SecondsAccumulated=%s",
				i+1, subsidy.ID, subsidy.Account.ID, subsidy.SecondsAccumulated)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
