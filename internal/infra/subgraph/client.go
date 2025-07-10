package subgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-pkgz/lgr"
)

//go:generate moq -out subgraph_mocks.go . SubgraphClient

// SubgraphClient defines the interface for subgraph operations
type SubgraphClient interface {
	ExecuteQuery(ctx context.Context, request GraphQLRequest, response interface{}) error
	QueryAccounts(ctx context.Context) ([]Account, error)
	QueryAccountSubsidiesForVault(ctx context.Context, vaultAddress string) ([]AccountSubsidy, error)
	QueryCompletedEpochs(ctx context.Context) ([]Epoch, error)
	QueryEpochByNumber(ctx context.Context, epochNumber string) (*Epoch, error)
	QueryCurrentActiveEpoch(ctx context.Context) (*Epoch, error)
	QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*Epoch, error)
	QueryMerkleDistributionForEpoch(
		ctx context.Context,
		epochNumber string,
		vaultAddress string,
	) (*MerkleDistribution, error)
	QueryAccountSubsidiesAtBlock(
		ctx context.Context,
		vaultAddress string,
		blockNumber int64,
	) ([]AccountSubsidy, error)
	QueryAccountSubsidiesForEpoch(
		ctx context.Context,
		vaultAddress string,
		epochEndTimestamp string,
	) ([]AccountSubsidy, error)
	HealthCheck(ctx context.Context) error
	ExecutePaginatedQuery(
		ctx context.Context,
		queryTemplate string,
		variables map[string]interface{},
		entityField string,
		response interface{},
	) error
	ExecuteQueryAtBlock(
		ctx context.Context,
		query string,
		variables map[string]interface{},
		blockNumber int64,
		response interface{},
	) error
	ExecutePaginatedQueryAtBlock(
		ctx context.Context,
		queryTemplate string,
		variables map[string]interface{},
		entityField string,
		blockNumber int64,
		response interface{},
	) error
}

type Client struct {
	httpClient *http.Client
	endpoint   string
	logger     lgr.L
}

// Ensure Client implements SubgraphClient
var _ SubgraphClient = (*Client)(nil)

func NewClient(endpoint string, logger lgr.L) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		endpoint: endpoint,
		logger:   logger,
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
	ID                   string `json:"id"`
	ContractAddress      string `json:"contractAddress"`
	Name                 string `json:"name"`
	Symbol               string `json:"symbol"`
	TotalSupply          string `json:"totalSupply"`
	CollectionType       string `json:"collectionType"`
	IsActive             bool   `json:"isActive"`
	YieldSharePercentage string `json:"yieldSharePercentage"`
	WeightFunctionType   string `json:"weightFunctionType"`
	WeightFunctionP1     string `json:"weightFunctionP1"`
	WeightFunctionP2     string `json:"weightFunctionP2"`
	MinBorrowAmount      string `json:"minBorrowAmount"`
	MaxBorrowAmount      string `json:"maxBorrowAmount"`
	TotalNFTsDeposited   string `json:"totalNFTsDeposited"`
	UpdatedAtBlock       string `json:"updatedAtBlock"`
	UpdatedAtTimestamp   string `json:"updatedAtTimestamp"`
}

type MerkleDistribution struct {
	ID                 string `json:"id"`
	Epoch              Epoch  `json:"epoch"`
	Vault              string `json:"vault"`
	MerkleRoot         string `json:"merkleRoot"`
	TotalAmount        string `json:"totalAmount"`
	Timestamp          string `json:"timestamp"`
	UpdatedAtBlock     string `json:"updatedAtBlock"`
	UpdatedAtTimestamp string `json:"updatedAtTimestamp"`
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
	Block     *BlockParameter        `json:"block,omitempty"`
}

// BlockParameter represents a block constraint for GraphQL queries
type BlockParameter struct {
	Number *int64  `json:"number,omitempty"`
	Hash   *string `json:"hash,omitempty"`
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

// AccountSubsidiesResponse represents the response for account subsidies query
type AccountSubsidiesResponse struct {
	AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
}

// EpochsResponse represents the response for epochs query
type EpochsResponse struct {
	Epochs []Epoch `json:"epoches"`
}

// MerkleDistributionsResponse represents the response for merkle distributions query
type MerkleDistributionsResponse struct {
	MerkleDistributions []MerkleDistribution `json:"merkleDistributions"`
}

// UsersResponse is kept for backward compatibility

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

// QueryAccountSubsidiesForVault queries all account subsidies for a specific vault to generate merkle proofs
func (c *Client) QueryAccountSubsidiesForVault(ctx context.Context, vaultAddress string) ([]AccountSubsidy, error) {
	query := `
		query GetAccountSubsidies($vaultId: String!, $first: Int!, $skip: Int!) {
			accountSubsidies(
				where: { 
					collectionParticipation_: { vault: $vaultId }
					secondsAccumulated_gt: "0" 
				}
				orderBy: id
				orderDirection: asc
				first: $first
				skip: $skip
			) {
				id
				account { id }
				secondsAccumulated
				secondsClaimed
				lastEffectiveValue
				updatedAtTimestamp
				collectionParticipation
			}
		}
	`

	variables := map[string]interface{}{
		"vaultId": vaultAddress,
	}

	var response AccountSubsidiesResponse

	if err := c.ExecutePaginatedQuery(ctx, query, variables, "accountSubsidies", &response); err != nil {
		return nil, fmt.Errorf("failed to query account subsidies for vault %s: %w", vaultAddress, err)
	}

	return response.AccountSubsidies, nil
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
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Logf("WARN failed to close response body: %v", closeErr)
		}
	}()

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
func (c *Client) ExecutePaginatedQuery(
	ctx context.Context,
	queryTemplate string,
	variables map[string]interface{},
	entityField string,
	response interface{},
) error {
	allResults, err := c.fetchAllPages(ctx, queryTemplate, variables, entityField, nil)
	if err != nil {
		return err
	}

	return c.reconstructResponse(allResults, entityField, response)
}

// ExecutePaginatedQueryAtBlock executes a GraphQL query with pagination at a specific block
func (c *Client) ExecutePaginatedQueryAtBlock(
	ctx context.Context,
	queryTemplate string,
	variables map[string]interface{},
	entityField string,
	blockNumber int64,
	response interface{},
) error {
	blockParam := &BlockParameter{
		Number: &blockNumber,
	}

	allResults, err := c.fetchAllPages(ctx, queryTemplate, variables, entityField, blockParam)
	if err != nil {
		return err
	}

	return c.reconstructResponse(allResults, entityField, response)
}

// ExecuteQueryAtBlock executes a GraphQL query at a specific block
func (c *Client) ExecuteQueryAtBlock(ctx context.Context, query string, variables map[string]interface{}, blockNumber int64, response interface{}) error {
	req := GraphQLRequest{
		Query:     query,
		Variables: variables,
		Block: &BlockParameter{
			Number: &blockNumber,
		},
	}
	return c.executeQuery(ctx, req, response)
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

// QueryCompletedEpochs queries all completed epochs, ordered by epoch number descending
func (c *Client) QueryCompletedEpochs(ctx context.Context) ([]Epoch, error) {
	query := `
		query GetCompletedEpochs($first: Int!, $skip: Int!) {
			epoches(
				where: { 
					status: "COMPLETED"
					processingCompletedTimestamp_not: null
				}
				orderBy: epochNumber
				orderDirection: desc
				first: $first
				skip: $skip
			) {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
				processingCompletedTimestamp
				totalSubsidiesDistributed
				totalYieldDistributed
				updatedAtTimestamp
			}
		}
	`

	var response EpochsResponse

	if err := c.ExecutePaginatedQuery(ctx, query, map[string]interface{}{}, "epoches", &response); err != nil {
		return nil, fmt.Errorf("failed to query completed epochs: %w", err)
	}

	return response.Epochs, nil
}

// QueryEpochByNumber queries a specific epoch by its number
func (c *Client) QueryEpochByNumber(ctx context.Context, epochNumber string) (*Epoch, error) {
	query := `
		query GetEpochByNumber($epochNumber: String!) {
			epoches(
				where: { 
					epochNumber: $epochNumber
					status: "COMPLETED"
					processingCompletedTimestamp_not: null
				}
				first: 1
			) {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
				processingCompletedTimestamp
				totalSubsidiesDistributed
				totalYieldDistributed
				updatedAtTimestamp
			}
		}
	`

	variables := map[string]interface{}{
		"epochNumber": epochNumber,
	}

	req := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	var response EpochsResponse

	if err := c.executeQuery(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to query epoch by number %s: %w", epochNumber, err)
	}

	if len(response.Epochs) == 0 {
		return nil, fmt.Errorf("epoch %s not found or not completed", epochNumber)
	}

	return &response.Epochs[0], nil
}

// QueryCurrentActiveEpoch queries the current active epoch to get its creation block
func (c *Client) QueryCurrentActiveEpoch(ctx context.Context) (*Epoch, error) {
	query := `
		query GetCurrentActiveEpoch {
			epoches(
				where: { 
					status: "ACTIVE"
				}
				orderBy: epochNumber
				orderDirection: desc
				first: 1
			) {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
				processingCompletedTimestamp
				createdAtBlock
				createdAtTimestamp
				updatedAtBlock
				updatedAtTimestamp
			}
		}
	`

	req := GraphQLRequest{
		Query: query,
	}

	var response EpochsResponse

	if err := c.executeQuery(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to query current active epoch: %w", err)
	}

	if len(response.Epochs) == 0 {
		return nil, fmt.Errorf("no active epoch found")
	}

	return &response.Epochs[0], nil
}

// QueryEpochWithBlockInfo queries a specific epoch by its number including block information
func (c *Client) QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*Epoch, error) {
	query := `
		query GetEpochWithBlockInfo($epochNumber: String!) {
			epoches(
				where: { 
					epochNumber: $epochNumber
				}
				first: 1
			) {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
				processingCompletedTimestamp
				totalSubsidiesDistributed
				totalYieldDistributed
				createdAtBlock
				createdAtTimestamp
				updatedAtBlock
				updatedAtTimestamp
			}
		}
	`

	variables := map[string]interface{}{
		"epochNumber": epochNumber,
	}

	req := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	var response EpochsResponse

	if err := c.executeQuery(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to query epoch with block info %s: %w", epochNumber, err)
	}

	if len(response.Epochs) == 0 {
		return nil, fmt.Errorf("epoch %s not found", epochNumber)
	}

	return &response.Epochs[0], nil
}

// QueryMerkleDistributionForEpoch queries the merkle distribution for a specific epoch and vault
func (c *Client) QueryMerkleDistributionForEpoch(ctx context.Context, epochNumber string, vaultAddress string) (*MerkleDistribution, error) {
	query := `
		query GetMerkleDistribution($epochNumber: String!, $vaultAddress: String!) {
			merkleDistributions(
				where: { 
					epoch_: { epochNumber: $epochNumber }
					vault: $vaultAddress
				}
				first: 1
			) {
				id
				vault
				merkleRoot
				totalAmount
				timestamp
				epoch {
					id
					epochNumber
					processingCompletedTimestamp
					endTimestamp
				}
			}
		}
	`

	variables := map[string]interface{}{
		"epochNumber":  epochNumber,
		"vaultAddress": vaultAddress,
	}

	req := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	var response MerkleDistributionsResponse

	if err := c.executeQuery(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to query merkle distribution for epoch %s vault %s: %w", epochNumber, vaultAddress, err)
	}

	if len(response.MerkleDistributions) == 0 {
		return nil, fmt.Errorf("merkle distribution not found for epoch %s vault %s", epochNumber, vaultAddress)
	}

	return &response.MerkleDistributions[0], nil
}

// QueryAccountSubsidiesAtBlock queries all account subsidies for a specific vault at a specific block
// This ensures block-consistent data for merkle tree generation
func (c *Client) QueryAccountSubsidiesAtBlock(ctx context.Context, vaultAddress string, blockNumber int64) ([]AccountSubsidy, error) {
	query := `
		query GetAccountSubsidiesAtBlock($vaultId: String!, $first: Int!, $skip: Int!) {
			accountSubsidies(
				where: { 
					collectionParticipation_: { vault: $vaultId }
					secondsAccumulated_gt: "0" 
				}
				orderBy: id
				orderDirection: asc
				first: $first
				skip: $skip
			) {
				id
				account { id }
				secondsAccumulated
				secondsClaimed
				lastEffectiveValue
				updatedAtTimestamp
				updatedAtBlock
				collectionParticipation
			}
		}
	`

	variables := map[string]interface{}{
		"vaultId": vaultAddress,
	}

	var response AccountSubsidiesResponse

	if err := c.ExecutePaginatedQueryAtBlock(ctx, query, variables, "accountSubsidies", blockNumber, &response); err != nil {
		return nil, fmt.Errorf("failed to query account subsidies at block %d for vault %s: %w", blockNumber, vaultAddress, err)
	}

	return response.AccountSubsidies, nil
}

// QueryAccountSubsidiesForEpoch queries account subsidies using the epoch completion timestamp
// This ensures we get the same data that was used during epoch processing
func (c *Client) QueryAccountSubsidiesForEpoch(ctx context.Context, vaultAddress string, epochEndTimestamp string) ([]AccountSubsidy, error) {
	query := `
		query GetAccountSubsidiesForEpoch($vaultId: String!, $epochEndTimestamp: String!, $first: Int!, $skip: Int!) {
			accountSubsidies(
				where: { 
					collectionParticipation_: { vault: $vaultId }
					secondsAccumulated_gt: "0"
					updatedAtTimestamp_lte: $epochEndTimestamp
				}
				orderBy: id
				orderDirection: asc
				first: $first
				skip: $skip
			) {
				id
				account { id }
				secondsAccumulated
				secondsClaimed
				lastEffectiveValue
				updatedAtTimestamp
				collectionParticipation
			}
		}
	`

	variables := map[string]interface{}{
		"vaultId":           vaultAddress,
		"epochEndTimestamp": epochEndTimestamp,
	}

	var response AccountSubsidiesResponse

	if err := c.ExecutePaginatedQuery(ctx, query, variables, "accountSubsidies", &response); err != nil {
		return nil, fmt.Errorf("failed to query account subsidies for epoch timestamp %s vault %s: %w", epochEndTimestamp, vaultAddress, err)
	}

	return response.AccountSubsidies, nil
}

// fetchAllPages fetches all pages of a paginated GraphQL query
func (c *Client) fetchAllPages(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, blockParam *BlockParameter) ([]json.RawMessage, error) {
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
			Block:     blockParam,
		}

		var data map[string]json.RawMessage

		if err := c.executeQuery(ctx, req, &data); err != nil {
			if blockParam != nil {
				return nil, fmt.Errorf("failed to execute paginated query at block %d skip %d: %w", *blockParam.Number, skip, err)
			}
			return nil, fmt.Errorf("failed to execute paginated query at skip %d: %w", skip, err)
		}

		entitiesRaw, ok := data[entityField]
		if !ok {
			return nil, fmt.Errorf("missing %s field in response", entityField)
		}

		// Parse entities as array
		var entities []json.RawMessage
		if err := json.Unmarshal(entitiesRaw, &entities); err != nil {
			return nil, fmt.Errorf("failed to parse %s array: %w", entityField, err)
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

	return allResults, nil
}

// reconstructResponse reconstructs the full GraphQL response from accumulated results
func (c *Client) reconstructResponse(allResults []json.RawMessage, entityField string, response interface{}) error {
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
