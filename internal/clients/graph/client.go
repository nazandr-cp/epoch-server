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

type User struct {
	ID                           string `json:"id"`
	TotalSecondsClaimed          string `json:"totalSecondsClaimed"`
	TotalSubsidiesReceived       string `json:"totalSubsidiesReceived"`
	TotalYieldEarned             string `json:"totalYieldEarned"`
	TotalBorrowVolume            string `json:"totalBorrowVolume"`
	TotalNFTsOwned               string `json:"totalNFTsOwned"`
	TotalCollectionsParticipated string `json:"totalCollectionsParticipated"`
	FirstInteractionBlock        string `json:"firstInteractionBlock"`
	FirstInteractionTimestamp    string `json:"firstInteractionTimestamp"`
	UpdatedAtBlock               string `json:"updatedAtBlock"`
	UpdatedAtTimestamp           string `json:"updatedAtTimestamp"`
}

type Eligibility struct {
	ID                    string     `json:"id"`
	User                  User       `json:"user"`
	Epoch                 Epoch      `json:"epoch"`
	Collection            Collection `json:"collection"`
	NFTBalance            string     `json:"nftBalance"`
	BorrowBalance         string     `json:"borrowBalance"`
	HoldingDuration       string     `json:"holdingDuration"`
	IsEligible            bool       `json:"isEligible"`
	SubsidyReceived       string     `json:"subsidyReceived"`
	YieldShare            string     `json:"yieldShare"`
	BonusMultiplier       string     `json:"bonusMultiplier"`
	CalculatedAtBlock     string     `json:"calculatedAtBlock"`
	CalculatedAtTimestamp string     `json:"calculatedAtTimestamp"`
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
	ParticipantCount              string `json:"participantCount"`
	ProcessingTimeMs              string `json:"processingTimeMs"`
	EstimatedProcessingTime       string `json:"estimatedProcessingTime"`
	ProcessingGasUsed             string `json:"processingGasUsed"`
	ProcessingTransactionCount    string `json:"processingTransactionCount"`
	CreatedAtBlock                string `json:"createdAtBlock"`
	CreatedAtTimestamp            string `json:"createdAtTimestamp"`
	UpdatedAtBlock                string `json:"updatedAtBlock"`
	UpdatedAtTimestamp            string `json:"updatedAtTimestamp"`
}

type Collection struct {
	ID                     string `json:"id"`
	ContractAddress        string `json:"contractAddress"`
	Name                   string `json:"name"`
	Symbol                 string `json:"symbol"`
	TotalSupply            string `json:"totalSupply"`
	CollectionType         string `json:"collectionType"`
	IsActive               bool   `json:"isActive"`
	YieldSharePercentage   string `json:"yieldSharePercentage"`
	WeightFunctionType     string `json:"weightFunctionType"`
	WeightFunctionP1       string `json:"weightFunctionP1"`
	WeightFunctionP2       string `json:"weightFunctionP2"`
	MinBorrowAmount        string `json:"minBorrowAmount"`
	MaxBorrowAmount        string `json:"maxBorrowAmount"`
	TotalNFTsDeposited     string `json:"totalNFTsDeposited"`
	TotalBorrowVolume      string `json:"totalBorrowVolume"`
	TotalYieldGenerated    string `json:"totalYieldGenerated"`
	TotalSubsidiesReceived string `json:"totalSubsidiesReceived"`
	RegisteredAtBlock      string `json:"registeredAtBlock"`
	RegisteredAtTimestamp  string `json:"registeredAtTimestamp"`
	UpdatedAtBlock         string `json:"updatedAtBlock"`
	UpdatedAtTimestamp     string `json:"updatedAtTimestamp"`
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

type UsersResponse struct {
	Users []User `json:"users"`
}

type EligibilitiesResponse struct {
	UserEpochEligibilities []Eligibility `json:"userEpochEligibilities"`
}

func (c *Client) QueryUsers(ctx context.Context) ([]User, error) {
	query := `
		query GetUsers($first: Int!, $skip: Int!) {
			users(first: $first, skip: $skip) {
				id
				totalSecondsClaimed
				totalSubsidiesReceived
				totalYieldEarned
				totalBorrowVolume
				totalNFTsOwned
				totalCollectionsParticipated
				firstInteractionBlock
				firstInteractionTimestamp
				updatedAtBlock
				updatedAtTimestamp
			}
		}
	`

	var response struct {
		Data UsersResponse `json:"data"`
	}

	if err := c.ExecutePaginatedQuery(ctx, query, map[string]interface{}{}, "users", &response); err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}

	return response.Data.Users, nil
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
					firstInteractionBlock
					firstInteractionTimestamp
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
					participantCount
					processingTimeMs
					estimatedProcessingTime
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
					totalBorrowVolume
					totalYieldGenerated
					totalSubsidiesReceived
					registeredAtBlock
					registeredAtTimestamp
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
				calculatedAtBlock
				calculatedAtTimestamp
			}
		}
	`

	var response struct {
		Data EligibilitiesResponse `json:"data"`
	}

	variables := map[string]interface{}{
		"epochId": epochID,
	}

	if err := c.ExecutePaginatedQuery(ctx, query, variables, "userEpochEligibilities", &response); err != nil {
		return nil, fmt.Errorf("failed to query eligibility for epoch %s: %w", epochID, err)
	}

	return response.Data.UserEpochEligibilities, nil
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

		var pageResponse struct {
			Data json.RawMessage `json:"data"`
		}

		if err := c.executeQuery(ctx, req, &pageResponse); err != nil {
			return fmt.Errorf("failed to execute paginated query at skip %d: %w", skip, err)
		}

		// Parse the data to extract the entity field
		var data map[string]json.RawMessage
		if err := json.Unmarshal(pageResponse.Data, &data); err != nil {
			return fmt.Errorf("failed to parse response data: %w", err)
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

	// Reconstruct the full response
	allEntitiesJson, err := json.Marshal(allResults)
	if err != nil {
		return fmt.Errorf("failed to marshal accumulated results: %w", err)
	}

	fullResponseData := map[string]json.RawMessage{
		entityField: allEntitiesJson,
	}

	fullResponseDataJson, err := json.Marshal(fullResponseData)
	if err != nil {
		return fmt.Errorf("failed to marshal response data: %w", err)
	}

	fullResponse := map[string]json.RawMessage{
		"data": fullResponseDataJson,
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
