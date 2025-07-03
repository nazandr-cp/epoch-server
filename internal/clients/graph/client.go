package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
}

func NewClient(endpoint string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		endpoint: endpoint,
	}
}

// Account represents a user account in the new schema (previously User)
type Account struct {
	ID                           string `json:"id"`
	TotalSecondsClaimed          string `json:"totalSecondsClaimed"`
	TotalSubsidiesReceived       string `json:"totalSubsidiesReceived"`
	TotalYieldEarned             string `json:"totalYieldEarned"`
	TotalBorrowVolume            string `json:"totalBorrowVolume"`
	TotalNFTsOwned               string `json:"totalNFTsOwned"`
	TotalCollectionsParticipated string `json:"totalCollectionsParticipated"`
	CreatedAtBlock               string `json:"createdAtBlock"`
	CreatedAtTimestamp           string `json:"createdAtTimestamp"`
	UpdatedAtBlock               string `json:"updatedAtBlock"`
	UpdatedAtTimestamp           string `json:"updatedAtTimestamp"`
}

// AccountSubsidy represents the new consolidated account subsidy entity
type AccountSubsidy struct {
	ID                      string  `json:"id"`
	Account                 Account `json:"account"`
	AccountMarket           string  `json:"accountMarket"`
	CollectionParticipation string  `json:"collectionParticipation"`
	BalanceNFT              string  `json:"balanceNFT"`
	SecondsAccumulated      string  `json:"secondsAccumulated"`
	SecondsClaimed          string  `json:"secondsClaimed"`
	SubsidiesAccrued        string  `json:"subsidiesAccrued"`
	SubsidiesClaimed        string  `json:"subsidiesClaimed"`
	AverageHoldingPeriod    string  `json:"averageHoldingPeriod"`
	TotalRewardsEarned      string  `json:"totalRewardsEarned"`
	LastEffectiveValue      string  `json:"lastEffectiveValue"`
	UpdatedAtBlock          string  `json:"updatedAtBlock"`
	UpdatedAtTimestamp      string  `json:"updatedAtTimestamp"`
}

type Eligibility struct {
	ID                    string     `json:"id"`
	Account               Account    `json:"user"` // Changed to match subgraph schema
	Epoch                 Epoch      `json:"epoch"`
	Collection            Collection `json:"collection"`
	NFTBalance            string     `json:"nftBalance"`
	BorrowBalance         string     `json:"borrowBalance"`
	HoldingDuration       string     `json:"holdingDuration"`
	IsEligible            bool       `json:"isEligible"`
	SubsidyReceived       string     `json:"subsidyReceived"`
	YieldShare            string     `json:"yieldShare"`
	BonusMultiplier       string     `json:"bonusMultiplier"`
}

type Epoch struct {
	ID                            string `json:"id"`
	EpochNumber                   string `json:"epochNumber"`
	Status                        string `json:"status"`
	StartTimestamp                string `json:"startTimestamp"`
	EndTimestamp                  string `json:"endTimestamp"`
	ProcessingStartedTimestamp    string `json:"processingStartedTimestamp"`
	ProcessingCompletedTimestamp  string `json:"processingCompletedTimestamp"`
	TotalYieldAvailable           string `json:"totalYieldAvailable"`
	TotalYieldAllocated           string `json:"totalYieldAllocated"`
	TotalYieldDistributed         string `json:"totalYieldDistributed"`
	RemainingYield                string `json:"remainingYield"`
	TotalSubsidiesDistributed     string `json:"totalSubsidiesDistributed"`
	TotalEligibleUsers            string `json:"totalEligibleUsers"`
	TotalParticipatingCollections string `json:"totalParticipatingCollections"`
	ProcessingTimeMs              string `json:"processingTimeMs"`
	ProcessingGasUsed             string `json:"processingGasUsed"`
	ProcessingTransactionCount    string `json:"processingTransactionCount"`
	CreatedAtBlock                string `json:"createdAtBlock"`
	CreatedAtTimestamp            string `json:"createdAtTimestamp"`
	UpdatedAtBlock                string `json:"updatedAtBlock"`
	UpdatedAtTimestamp            string `json:"updatedAtTimestamp"`
}

type Collection struct {
	ID                    string `json:"id"`
	ContractAddress       string `json:"contractAddress"`
	Name                  string `json:"name"`
	Symbol                string `json:"symbol"`
	TotalSupply           string `json:"totalSupply"`
	CollectionType        string `json:"collectionType"`
	IsActive              bool   `json:"isActive"`
	YieldSharePercentage  string `json:"yieldSharePercentage"`
	WeightFunctionType    string `json:"weightFunctionType"`
	WeightFunctionP1      string `json:"weightFunctionP1"`
	WeightFunctionP2      string `json:"weightFunctionP2"`
	MinBorrowAmount       string `json:"minBorrowAmount"`
	MaxBorrowAmount       string `json:"maxBorrowAmount"`
	TotalNFTsDeposited    string `json:"totalNFTsDeposited"`
	UpdatedAtBlock        string `json:"updatedAtBlock"`
	UpdatedAtTimestamp    string `json:"updatedAtTimestamp"`
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// AccountsResponse represents the response for accounts query (new schema)
type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

// UsersResponse is kept for backward compatibility

type EligibilitiesResponse struct {
	UserEpochEligibilities []Eligibility `json:"userEpochEligibilities"`
}

// QueryAccounts queries accounts using the new schema
func (c *Client) QueryAccounts(ctx context.Context) ([]Account, error) {
	query := `
		query GetAccounts($first: Int!, $skip: Int!) {
			accounts(first: $first, skip: $skip) {
				id
				totalSecondsClaimed
				totalSubsidiesReceived
				totalYieldEarned
				totalBorrowVolume
				totalNFTsOwned
				totalCollectionsParticipated
				createdAtBlock
				createdAtTimestamp
				updatedAtBlock
				updatedAtTimestamp
			}
		}
	`

	var response AccountsResponse

	if err := c.ExecutePaginatedQuery(ctx, query, map[string]interface{}{}, "accounts", &response); err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}

	return response.Accounts, nil
}

func (c *Client) QueryEligibility(ctx context.Context, epochID string) ([]Eligibility, error) {
	query := `
		query GetEligibility($epochId: String!, $first: Int!, $skip: Int!) {
			userEpochEligibilities(where: { epoch: $epochId }, first: $first, skip: $skip) {
				id
				user {
					id
					totalSecondsClaimed
					totalSubsidiesReceived
					totalYieldEarned
					totalBorrowVolume
					totalNFTsOwned
					totalCollectionsParticipated
					createdAtBlock
					createdAtTimestamp
					updatedAtBlock
					updatedAtTimestamp
				}
				epoch {
					id
					epochNumber
					status
					startTimestamp
					endTimestamp
					processingStartedTimestamp
					processingCompletedTimestamp
					totalYieldAvailable
					totalYieldAllocated
					totalYieldDistributed
					remainingYield
					totalSubsidiesDistributed
					totalEligibleUsers
					totalParticipatingCollections
					processingTimeMs
					processingGasUsed
					processingTransactionCount
					createdAtBlock
					createdAtTimestamp
					updatedAtBlock
					updatedAtTimestamp
				}
				collection {
					id
					contractAddress
					name
					symbol
					totalSupply
					collectionType
					isActive
					yieldSharePercentage
					weightFunctionType
					weightFunctionP1
					weightFunctionP2
					minBorrowAmount
					maxBorrowAmount
					totalNFTsDeposited
					updatedAtBlock
					updatedAtTimestamp
				}
				nftBalance
				borrowBalance
				holdingDuration
				isEligible
				subsidyReceived
				yieldShare
				bonusMultiplier
			}
		}
	`

	var response EligibilitiesResponse

	variables := map[string]interface{}{
		"epochId": epochID,
	}

	if err := c.ExecutePaginatedQuery(ctx, query, variables, "userEpochEligibilities", &response); err != nil {
		return nil, fmt.Errorf("failed to query eligibility for epoch %s: %w", epochID, err)
	}

	return response.UserEpochEligibilities, nil
}

func (c *Client) executeQuery(ctx context.Context, request GraphQLRequest, response interface{}) error {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var graphQLResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphQLResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if len(graphQLResp.Errors) > 0 {
		return fmt.Errorf("GraphQL errors: %v", graphQLResp.Errors)
	}

	data, err := json.Marshal(graphQLResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal response data: %w", err)
	}

	if err := json.Unmarshal(data, response); err != nil {
		return fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	return nil
}

func (c *Client) ExecuteQuery(ctx context.Context, request GraphQLRequest, response interface{}) error {
	return c.executeQuery(ctx, request, response)
}

// ExecutePaginatedQuery executes a GraphQL query with pagination, automatically fetching all pages
func (c *Client) ExecutePaginatedQuery(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, response interface{}) error {
	const pageSize = 1000
	var allResults []json.RawMessage
	skip := 0

	for {
		// Add pagination variables
		paginatedVars := make(map[string]interface{})
		for k, v := range variables {
			paginatedVars[k] = v
		}
		paginatedVars["first"] = pageSize
		paginatedVars["skip"] = skip

		req := GraphQLRequest{
			Query:     queryTemplate,
			Variables: paginatedVars,
		}

		var data map[string]json.RawMessage

		if err := c.executeQuery(ctx, req, &data); err != nil {
			return fmt.Errorf("failed to execute paginated query at skip %d: %w", skip, err)
		}

		entitiesRaw, ok := data[entityField]
		if !ok {
			return fmt.Errorf("missing %s field in response", entityField)
		}

		// Parse entities as array
		var entities []json.RawMessage
		if err := json.Unmarshal(entitiesRaw, &entities); err != nil {
			return fmt.Errorf("failed to parse %s array: %w", entityField, err)
		}

		// If this page is empty, we've reached the end
		if len(entities) == 0 {
			break
		}

		allResults = append(allResults, entities...)

		// If this page has fewer results than pageSize, we've reached the end
		if len(entities) < pageSize {
			break
		}

		skip += pageSize
	}

	// Reconstruct the full response with nested data field to match GraphQL standard
	allEntitiesJson, err := json.Marshal(allResults)
	if err != nil {
		return fmt.Errorf("failed to marshal accumulated results: %w", err)
	}

	dataField := map[string]json.RawMessage{
		entityField: allEntitiesJson,
	}
	
	dataFieldJson, err := json.Marshal(dataField)
	if err != nil {
		return fmt.Errorf("failed to marshal data field: %w", err)
	}

	fullResponse := map[string]json.RawMessage{
		"data": dataFieldJson,
	}

	fullResponseJson, err := json.Marshal(fullResponse)
	if err != nil {
		return fmt.Errorf("failed to marshal full response: %w", err)
	}

	if err := json.Unmarshal(fullResponseJson, response); err != nil {
		return fmt.Errorf("failed to unmarshal into response type: %w", err)
	}

	return nil
}

// HealthCheck performs a basic connectivity check to verify the subgraph is accessible
// It executes a simple introspection query to validate the GraphQL endpoint is responding
func (c *Client) HealthCheck(ctx context.Context) error {
	// Use a simple introspection query to check if the subgraph is accessible
	query := `
		query HealthCheck {
			__schema {
				queryType {
					name
				}
			}
		}
	`

	req := GraphQLRequest{
		Query: query,
	}

	var response map[string]interface{}

	if err := c.executeQuery(ctx, req, &response); err != nil {
		return fmt.Errorf("subgraph health check failed: %w", err)
	}

	// Check if the response contains the expected schema information
	if schema, ok := response["__schema"]; ok {
		if schemaMap, ok := schema.(map[string]interface{}); ok {
			if queryType, ok := schemaMap["queryType"]; ok {
				if queryTypeMap, ok := queryType.(map[string]interface{}); ok {
					if name, ok := queryTypeMap["name"]; ok && name == "Query" {
						return nil
					}
				}
			}
		}
	}

	return fmt.Errorf("subgraph health check failed: unexpected response structure")
}
