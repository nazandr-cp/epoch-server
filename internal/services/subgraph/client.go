package subgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/go-pkgz/lgr"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
	logger     lgr.L
}

var _ subgraph.SubgraphClient = (*Client)(nil)

// ProvideClient creates a subgraph client
func ProvideClient(endpoint string, logger lgr.L) subgraph.SubgraphClient {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		endpoint: endpoint,
		logger:   logger,
	}
}

func (c *Client) QueryAccounts(ctx context.Context) ([]subgraph.Account, error) {
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

	var response subgraph.AccountsResponse

	if err := c.ExecutePaginatedQuery(ctx, query, map[string]interface{}{}, "accounts", &response); err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}

	return response.Accounts, nil
}

func (c *Client) QueryAccountSubsidiesForVault(
	ctx context.Context,
	vaultAddress string,
) ([]subgraph.AccountSubsidy, error) {
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

	var response subgraph.AccountSubsidiesResponse

	if err := c.ExecutePaginatedQuery(ctx, query, variables, "accountSubsidies", &response); err != nil {
		return nil, fmt.Errorf("failed to query account subsidies for vault %s: %w", vaultAddress, err)
	}

	return response.AccountSubsidies, nil
}

func (c *Client) executeQuery(ctx context.Context, request subgraph.GraphQLRequest, response interface{}) error {
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

	var graphQLResp subgraph.GraphQLResponse
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

func (c *Client) ExecuteQuery(ctx context.Context, request subgraph.GraphQLRequest, response interface{}) error {
	return c.executeQuery(ctx, request, response)
}

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

func (c *Client) ExecutePaginatedQueryAtBlock(
	ctx context.Context,
	queryTemplate string,
	variables map[string]interface{},
	entityField string,
	blockNumber int64,
	response interface{},
) error {
	blockParam := &subgraph.BlockParameter{
		Number: &blockNumber,
	}

	allResults, err := c.fetchAllPages(ctx, queryTemplate, variables, entityField, blockParam)
	if err != nil {
		return err
	}

	return c.reconstructResponse(allResults, entityField, response)
}

func (c *Client) ExecuteQueryAtBlock(
	ctx context.Context,
	query string,
	variables map[string]interface{},
	blockNumber int64,
	response interface{},
) error {
	req := subgraph.GraphQLRequest{
		Query:     query,
		Variables: variables,
		Block: &subgraph.BlockParameter{
			Number: &blockNumber,
		},
	}
	return c.executeQuery(ctx, req, response)
}

func (c *Client) HealthCheck(ctx context.Context) error {
	query := `
		query HealthCheck {
			__schema {
				queryType {
					name
				}
			}
		}
	`

	req := subgraph.GraphQLRequest{
		Query: query,
	}

	var response map[string]interface{}

	if err := c.executeQuery(ctx, req, &response); err != nil {
		return fmt.Errorf("subgraph health check failed: %w", err)
	}

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

func (c *Client) QueryCompletedEpochs(ctx context.Context) ([]subgraph.Epoch, error) {
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

	var response subgraph.EpochsResponse

	if err := c.ExecutePaginatedQuery(ctx, query, map[string]interface{}{}, "epoches", &response); err != nil {
		return nil, fmt.Errorf("failed to query completed epochs: %w", err)
	}

	return response.Epochs, nil
}

func (c *Client) QueryEpochByNumber(ctx context.Context, epochNumber string) (*subgraph.Epoch, error) {
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

	req := subgraph.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	var response subgraph.EpochsResponse

	if err := c.executeQuery(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to query epoch by number %s: %w", epochNumber, err)
	}

	if len(response.Epochs) == 0 {
		return nil, fmt.Errorf("epoch %s not found or not completed", epochNumber)
	}

	return &response.Epochs[0], nil
}

func (c *Client) QueryCurrentActiveEpoch(ctx context.Context) (*subgraph.Epoch, error) {
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

	req := subgraph.GraphQLRequest{
		Query: query,
	}

	var response subgraph.EpochsResponse

	if err := c.executeQuery(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to query current active epoch: %w", err)
	}

	if len(response.Epochs) == 0 {
		return nil, fmt.Errorf("no active epoch found")
	}

	return &response.Epochs[0], nil
}

func (c *Client) QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*subgraph.Epoch, error) {
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

	req := subgraph.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	var response subgraph.EpochsResponse

	if err := c.executeQuery(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to query epoch with block info %s: %w", epochNumber, err)
	}

	if len(response.Epochs) == 0 {
		return nil, fmt.Errorf("epoch %s not found", epochNumber)
	}

	return &response.Epochs[0], nil
}

func (c *Client) QueryMerkleDistributionForEpoch(
	ctx context.Context,
	epochNumber string,
	vaultAddress string,
) (*subgraph.MerkleDistribution, error) {
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

	req := subgraph.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	var response subgraph.MerkleDistributionsResponse

	if err := c.executeQuery(ctx, req, &response); err != nil {
		return nil, fmt.Errorf(
			"failed to query merkle distribution for epoch %s vault %s: %w",
			epochNumber,
			vaultAddress,
			err,
		)
	}

	if len(response.MerkleDistributions) == 0 {
		return nil, fmt.Errorf("merkle distribution not found for epoch %s vault %s", epochNumber, vaultAddress)
	}

	return &response.MerkleDistributions[0], nil
}

func (c *Client) QueryAccountSubsidiesAtBlock(
	ctx context.Context,
	vaultAddress string,
	blockNumber int64,
) ([]subgraph.AccountSubsidy, error) {
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

	var response subgraph.AccountSubsidiesResponse

	if err := c.ExecutePaginatedQueryAtBlock(
		ctx, query, variables, "accountSubsidies", blockNumber, &response,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to query account subsidies at block %d for vault %s: %w",
			blockNumber,
			vaultAddress,
			err,
		)
	}

	return response.AccountSubsidies, nil
}

func (c *Client) QueryAccountSubsidiesForEpoch(
	ctx context.Context,
	vaultAddress string,
	epochEndTimestamp string,
) ([]subgraph.AccountSubsidy, error) {
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

	var response subgraph.AccountSubsidiesResponse

	if err := c.ExecutePaginatedQuery(ctx, query, variables, "accountSubsidies", &response); err != nil {
		return nil, fmt.Errorf(
			"failed to query account subsidies for epoch timestamp %s vault %s: %w",
			epochEndTimestamp,
			vaultAddress,
			err,
		)
	}

	return response.AccountSubsidies, nil
}

func (c *Client) fetchAllPages(
	ctx context.Context,
	queryTemplate string,
	variables map[string]interface{},
	entityField string,
	blockParam *subgraph.BlockParameter,
) ([]json.RawMessage, error) {
	const pageSize = 1000
	var allResults []json.RawMessage
	skip := 0

	for {
		paginatedVars := make(map[string]interface{})
		for k, v := range variables {
			paginatedVars[k] = v
		}
		paginatedVars["first"] = pageSize
		paginatedVars["skip"] = skip

		req := subgraph.GraphQLRequest{
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

		var entities []json.RawMessage
		if err := json.Unmarshal(entitiesRaw, &entities); err != nil {
			return nil, fmt.Errorf("failed to parse %s array: %w", entityField, err)
		}

		if len(entities) == 0 {
			break
		}

		allResults = append(allResults, entities...)

		if len(entities) < pageSize {
			break
		}

		skip += pageSize
	}

	return allResults, nil
}

func (c *Client) reconstructResponse(allResults []json.RawMessage, entityField string, response interface{}) error {
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
