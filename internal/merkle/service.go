package merkle

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/andrey/epoch-server/internal/utils"
	"github.com/go-pkgz/lgr"
)

// GraphClient interface for merkle service operations
type GraphClient interface {
	QueryEpochByNumber(ctx context.Context, epochNumber string) (*graph.Epoch, error)
	QueryMerkleDistributionForEpoch(ctx context.Context, epochNumber string, vaultAddress string) (*graph.MerkleDistribution, error)
	QueryAccountSubsidiesForEpoch(ctx context.Context, vaultAddress string, epochEndTimestamp string) ([]graph.AccountSubsidy, error)
	ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error
	ExecutePaginatedQuery(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, response interface{}) error
}

// MerkleProofService handles on-demand merkle proof generation
type MerkleProofService struct {
	graphClient    GraphClient
	proofGenerator *ProofGenerator
	logger         lgr.L
}

// NewMerkleProofService creates a new merkle proof service
func NewMerkleProofService(graphClient GraphClient, logger lgr.L) *MerkleProofService {
	return &MerkleProofService{
		graphClient:    graphClient,
		proofGenerator: NewProofGeneratorWithDependencies(graphClient, logger),
		logger:         logger,
	}
}

// MerkleProofResponse represents the response structure for merkle proof requests
type MerkleProofResponse struct {
	UserAddress     string   `json:"userAddress"`
	VaultAddress    string   `json:"vaultAddress"`
	ClaimableAmount string   `json:"claimableAmount"`
	MerkleRoot      string   `json:"merkleRoot"`
	MerkleProof     []string `json:"merkleProof"`
	GeneratedAt     int64    `json:"generatedAt"`
	EpochNumber     string   `json:"epochNumber,omitempty"`
}

// GenerateUserMerkleProof generates a merkle proof for a specific user's subsidy claim
func (mps *MerkleProofService) GenerateUserMerkleProof(ctx context.Context, userAddress, vaultAddress string) (*MerkleProofResponse, error) {
	// Validate inputs
	if userAddress == "" {
		return nil, fmt.Errorf("userAddress cannot be empty")
	}
	if vaultAddress == "" {
		return nil, fmt.Errorf("vaultAddress cannot be empty")
	}

	// Normalize addresses
	userAddress = utils.NormalizeAddress(userAddress)
	vaultAddress = utils.NormalizeAddress(vaultAddress)

	mps.logger.Logf("INFO generating merkle proof for user %s in vault %s", userAddress, vaultAddress)

	// Query account subsidies for the vault
	subsidies, err := mps.queryAccountSubsidiesForVault(ctx, vaultAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to query account subsidies: %w", err)
	}

	// Use the unified merkle package to generate the tree
	result, err := mps.proofGenerator.GenerateTreeFromSubsidies(ctx, vaultAddress, subsidies)
	if err != nil {
		return nil, fmt.Errorf("failed to generate merkle tree: %w", err)
	}

	// Find the user's entry and amount
	var targetUserEarning *big.Int
	
	for _, entry := range result.Entries {
		normalizedAddr := utils.NormalizeAddress(entry.Address)
		
		if normalizedAddr == userAddress {
			targetUserEarning = entry.TotalEarned
			break
		}
	}

	// Check if target user has any earnings
	if targetUserEarning == nil {
		return nil, fmt.Errorf("user %s has no claimable subsidies in vault %s", userAddress, vaultAddress)
	}

	mps.logger.Logf("INFO user %s has claimable amount: %s", userAddress, targetUserEarning.String())

	// Generate proof for the target user
	proof, merkleRoot, err := mps.proofGenerator.GenerateProof(result.Entries, userAddress, targetUserEarning)
	if err != nil {
		return nil, fmt.Errorf("failed to generate merkle proof: %w", err)
	}

	// Convert proof to string array
	proofStrings := make([]string, len(proof))
	for i, p := range proof {
		proofStrings[i] = fmt.Sprintf("0x%x", p[:])
	}

	return &MerkleProofResponse{
		UserAddress:     userAddress,
		VaultAddress:    vaultAddress,
		ClaimableAmount: targetUserEarning.String(),
		MerkleRoot:      fmt.Sprintf("0x%x", merkleRoot[:]),
		MerkleProof:     proofStrings,
		GeneratedAt:     result.Timestamp,
	}, nil
}


// GenerateHistoricalMerkleProof generates a merkle proof for a specific user's subsidy claim from a historical epoch
// This ensures the proof is generated using the exact same data that was used during epoch processing
func (mps *MerkleProofService) GenerateHistoricalMerkleProof(ctx context.Context, userAddress, vaultAddress, epochNumber string) (*MerkleProofResponse, error) {
	// Validate inputs
	if userAddress == "" {
		return nil, fmt.Errorf("userAddress cannot be empty")
	}
	if vaultAddress == "" {
		return nil, fmt.Errorf("vaultAddress cannot be empty")
	}
	if epochNumber == "" {
		return nil, fmt.Errorf("epochNumber cannot be empty")
	}

	// Normalize addresses
	userAddress = utils.NormalizeAddress(userAddress)
	vaultAddress = utils.NormalizeAddress(vaultAddress)

	mps.logger.Logf("INFO generating historical merkle proof for user %s in vault %s for epoch %s", userAddress, vaultAddress, epochNumber)

	// Get the historical merkle distribution for validation
	merkleDistribution, err := mps.graphClient.QueryMerkleDistributionForEpoch(ctx, epochNumber, vaultAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to query merkle distribution for epoch %s vault %s: %w", epochNumber, vaultAddress, err)
	}

	// Query account subsidies as they were at the time of epoch processing
	accountSubsidies, err := mps.graphClient.QueryAccountSubsidiesForEpoch(ctx, vaultAddress, epochNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to query historical account subsidies: %w", err)
	}

	mps.logger.Logf("INFO found %d account subsidies for vault %s at epoch %s", len(accountSubsidies), vaultAddress, epochNumber)

	// Use the unified merkle package to generate the historical tree
	result, err := mps.proofGenerator.GenerateHistoricalTreeFromSubsidies(ctx, epochNumber, accountSubsidies)
	if err != nil {
		return nil, fmt.Errorf("failed to generate historical merkle tree: %w", err)
	}

	// Find the target user's earning
	var targetUserEarning *big.Int
	for _, entry := range result.Entries {
		if utils.NormalizeAddress(entry.Address) == userAddress {
			targetUserEarning = entry.TotalEarned
			break
		}
	}

	// Check if target user has any earnings
	if targetUserEarning == nil {
		return nil, fmt.Errorf("user %s has no claimable subsidies in vault %s for epoch %s", userAddress, vaultAddress, epochNumber)
	}

	mps.logger.Logf("INFO user %s has claimable amount: %s for epoch %s", userAddress, targetUserEarning.String(), epochNumber)
	mps.logger.Logf("INFO total users with earnings: %d", len(result.Entries))

	// Generate proof for the target user
	proof, merkleRoot, err := mps.proofGenerator.GenerateProof(result.Entries, userAddress, targetUserEarning)
	if err != nil {
		return nil, fmt.Errorf("failed to generate merkle proof: %w", err)
	}

	// Validate that the generated merkle root matches the stored historical root
	generatedRootHex := fmt.Sprintf("0x%x", merkleRoot[:])
	if generatedRootHex != merkleDistribution.MerkleRoot {
		return nil, fmt.Errorf("generated merkle root %s does not match stored historical root %s for epoch %s", 
			generatedRootHex, merkleDistribution.MerkleRoot, epochNumber)
	}

	// Convert proof to string array
	proofStrings := make([]string, len(proof))
	for i, p := range proof {
		proofStrings[i] = fmt.Sprintf("0x%x", p[:])
	}

	mps.logger.Logf("INFO generated historical merkle proof with %d elements for user %s epoch %s", len(proofStrings), userAddress, epochNumber)
	mps.logger.Logf("INFO merkle root validation passed: %s", generatedRootHex)

	return &MerkleProofResponse{
		UserAddress:     userAddress,
		VaultAddress:    vaultAddress,
		ClaimableAmount: targetUserEarning.String(),
		MerkleRoot:      generatedRootHex,
		MerkleProof:     proofStrings,
		GeneratedAt:     time.Now().Unix(),
		EpochNumber:     epochNumber,
	}, nil
}

// queryAccountSubsidiesForVault queries all account subsidies for a specific vault
func (mps *MerkleProofService) queryAccountSubsidiesForVault(ctx context.Context, vaultId string) ([]graph.AccountSubsidy, error) {
	query := `
		query GetAccountSubsidies($vaultId: ID!, $first: Int!, $skip: Int!) {
			accountSubsidies(
				where: { collectionParticipation_: { vault: $vaultId }, secondsAccumulated_gt: "0" }
				orderBy: id
				orderDirection: asc
				first: $first
				skip: $skip
			) {
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

	variables := map[string]interface{}{
		"vaultId": strings.ToLower(vaultId),
	}

	var response struct {
		Data struct {
			AccountSubsidies []graph.AccountSubsidy `json:"accountSubsidies"`
		} `json:"data"`
	}

	if err := mps.graphClient.ExecutePaginatedQuery(ctx, query, variables, "accountSubsidies", &response); err != nil {
		return nil, fmt.Errorf("failed to execute paginated GraphQL query: %w", err)
	}

	return response.Data.AccountSubsidies, nil
}

