package merkle

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/andrey/epoch-server/internal/clients/storage"
	"github.com/andrey/epoch-server/internal/utils"
	"github.com/go-pkgz/lgr"
)

// GraphClient interface for merkle service operations
type GraphClient interface {
	QueryEpochByNumber(ctx context.Context, epochNumber string) (*graph.Epoch, error)
	QueryCurrentActiveEpoch(ctx context.Context) (*graph.Epoch, error)
	QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*graph.Epoch, error)
	QueryAccountSubsidiesAtBlock(ctx context.Context, vaultAddress string, blockNumber int64) ([]graph.AccountSubsidy, error)
	QueryMerkleDistributionForEpoch(ctx context.Context, epochNumber string, vaultAddress string) (*graph.MerkleDistribution, error)
	QueryAccountSubsidiesForEpoch(ctx context.Context, vaultAddress string, epochEndTimestamp string) ([]graph.AccountSubsidy, error)
	ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error
	ExecutePaginatedQuery(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, response interface{}) error
	ExecuteQueryAtBlock(ctx context.Context, query string, variables map[string]interface{}, blockNumber int64, response interface{}) error
	ExecutePaginatedQueryAtBlock(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, blockNumber int64, response interface{}) error
}

// EpochStorageClient interface for accessing cached snapshots
type EpochStorageClient interface {
	GetLatestEpochSnapshot(ctx context.Context, vaultID string) (*storage.MerkleSnapshot, error)
	GetEpochSnapshot(ctx context.Context, epochNumber *big.Int, vaultID string) (*storage.MerkleSnapshot, error)
}

// MerkleProofService handles on-demand merkle proof generation
type MerkleProofService struct {
	graphClient    GraphClient
	storageClient  EpochStorageClient
	proofGenerator *ProofGenerator
	logger         lgr.L
}

// NewMerkleProofService creates a new merkle proof service
func NewMerkleProofService(graphClient GraphClient, logger lgr.L) *MerkleProofService {
	return &MerkleProofService{
		graphClient:    graphClient,
		storageClient:  nil, // Will fallback to subgraph queries
		proofGenerator: NewProofGeneratorWithDependencies(graphClient, logger),
		logger:         logger,
	}
}

// NewMerkleProofServiceWithStorage creates a new merkle proof service with storage client
func NewMerkleProofServiceWithStorage(graphClient GraphClient, storageClient EpochStorageClient, logger lgr.L) *MerkleProofService {
	return &MerkleProofService{
		graphClient:    graphClient,
		storageClient:  storageClient,
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
	BlockNumber     int64    `json:"blockNumber"`   // Block number used for data consistency
	DataTimestamp   int64    `json:"dataTimestamp"` // Timestamp used for calculations
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

	// Only use cached snapshot - no subgraph fallback
	if mps.storageClient == nil {
		return nil, fmt.Errorf("storage client not available")
	}

	response, err := mps.generateProofFromCache(ctx, userAddress, vaultAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to generate proof from snapshot: %w", err)
	}

	mps.logger.Logf("INFO generated proof from cached snapshot for user %s", userAddress)
	return response, nil
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

	// Try to use cached snapshot first if storage client is available
	if mps.storageClient != nil {
		if response, err := mps.generateHistoricalProofFromCache(ctx, userAddress, vaultAddress, epochNumber); err == nil {
			mps.logger.Logf("INFO generated historical proof from cached snapshot for user %s epoch %s", userAddress, epochNumber)
			return response, nil
		} else {
			mps.logger.Logf("WARN failed to generate historical proof from cache, falling back to subgraph: %v", err)
		}
	}

	// Get the historical merkle distribution for validation
	merkleDistribution, err := mps.graphClient.QueryMerkleDistributionForEpoch(ctx, epochNumber, vaultAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to query merkle distribution for epoch %s vault %s: %w", epochNumber, vaultAddress, err)
	}

	// Get the historical epoch's block information for consistency
	epochInfo, err := mps.getHistoricalEpochBlock(ctx, epochNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical epoch block info: %w", err)
	}

	mps.logger.Logf("INFO using historical epoch %s created at block %d for proof generation",
		epochNumber, epochInfo.CreatedAtBlock)

	// Query account subsidies as they were at the specific block
	accountSubsidies, err := mps.queryAccountSubsidiesAtBlock(ctx, vaultAddress, epochInfo.CreatedAtBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to query historical account subsidies at block %d: %w", epochInfo.CreatedAtBlock, err)
	}

	mps.logger.Logf("INFO found %d account subsidies for vault %s at epoch %s block %d",
		len(accountSubsidies), vaultAddress, epochNumber, epochInfo.CreatedAtBlock)

	// Use the unified merkle package to generate the historical tree with block consistency
	result, err := mps.proofGenerator.GenerateTreeFromSubsidiesAtBlock(ctx, vaultAddress, accountSubsidies, epochInfo)
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
		BlockNumber:     epochInfo.CreatedAtBlock,
		DataTimestamp:   result.Timestamp,
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

// getLatestMerkleEpochBlock retrieves the latest epoch that has a merkle distribution for this vault
// This ensures proof generation uses the same block state as epoch processing
func (mps *MerkleProofService) getLatestMerkleEpochBlock(ctx context.Context, vaultAddress string) (*EpochTimestamp, error) {
	// Create an epoch block manager to get epoch block info
	epochBlockManager := NewEpochBlockManager(mps.graphClient, mps.logger)

	// Get the latest epoch that has a merkle distribution for this vault
	epochInfo, err := epochBlockManager.GetLatestEpochTimestamp(ctx, vaultAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest epoch with merkle distribution: %w", err)
	}

	return epochInfo, nil
}

// getHistoricalEpochBlock retrieves block information for a specific historical epoch
func (mps *MerkleProofService) getHistoricalEpochBlock(ctx context.Context, epochNumber string) (*EpochTimestamp, error) {
	// Create an epoch block manager to get epoch block info
	epochBlockManager := NewEpochBlockManager(mps.graphClient, mps.logger)

	// Get the block information for the specific epoch
	epochInfo, err := epochBlockManager.GetEpochBlockByNumber(ctx, epochNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get epoch %s block info: %w", epochNumber, err)
	}

	return epochInfo, nil
}

// queryAccountSubsidiesAtBlock queries account subsidies at a specific block number
func (mps *MerkleProofService) queryAccountSubsidiesAtBlock(ctx context.Context, vaultAddress string, blockNumber int64) ([]graph.AccountSubsidy, error) {
	// Use the new block-based query method from the graph client
	subsidies, err := mps.graphClient.QueryAccountSubsidiesAtBlock(ctx, strings.ToLower(vaultAddress), blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to query account subsidies at block %d: %w", blockNumber, err)
	}

	mps.logger.Logf("INFO queried %d account subsidies for vault %s at block %d",
		len(subsidies), vaultAddress, blockNumber)

	return subsidies, nil
}

// generateProofFromCache generates a merkle proof using cached snapshot data
func (mps *MerkleProofService) generateProofFromCache(ctx context.Context, userAddress, vaultAddress string) (*MerkleProofResponse, error) {
	// Get the latest cached snapshot
	snapshot, err := mps.storageClient.GetLatestEpochSnapshot(ctx, vaultAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest epoch snapshot: %w", err)
	}

	// Convert storage entries to merkle entries
	entries := make([]Entry, len(snapshot.Entries))
	for i, entry := range snapshot.Entries {
		entries[i] = Entry{
			Address:     entry.Address,
			TotalEarned: entry.TotalEarned,
		}
	}

	// Find the user's entry and amount
	var targetUserEarning *big.Int
	for _, entry := range entries {
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

	mps.logger.Logf("INFO user %s has claimable amount: %s from cached snapshot", userAddress, targetUserEarning.String())

	// Generate proof for the target user
	proof, _, err := mps.proofGenerator.GenerateProof(entries, userAddress, targetUserEarning)
	if err != nil {
		return nil, fmt.Errorf("failed to generate merkle proof from cache: %w", err)
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
		MerkleRoot:      snapshot.MerkleRoot,
		MerkleProof:     proofStrings,
		GeneratedAt:     time.Now().Unix(),
		EpochNumber:     snapshot.EpochNumber.String(),
		BlockNumber:     snapshot.BlockNumber,
		DataTimestamp:   snapshot.Timestamp,
	}, nil
}

// generateHistoricalProofFromCache generates a historical merkle proof using cached snapshot data
func (mps *MerkleProofService) generateHistoricalProofFromCache(ctx context.Context, userAddress, vaultAddress, epochNumber string) (*MerkleProofResponse, error) {
	// Parse epoch number
	epochNum, ok := new(big.Int).SetString(epochNumber, 10)
	if !ok {
		return nil, fmt.Errorf("invalid epoch number: %s", epochNumber)
	}

	// Get the cached snapshot for the specific epoch
	snapshot, err := mps.storageClient.GetEpochSnapshot(ctx, epochNum, vaultAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get epoch snapshot: %w", err)
	}

	// Convert storage entries to merkle entries
	entries := make([]Entry, len(snapshot.Entries))
	for i, entry := range snapshot.Entries {
		entries[i] = Entry{
			Address:     entry.Address,
			TotalEarned: entry.TotalEarned,
		}
	}

	// Find the user's entry and amount
	var targetUserEarning *big.Int
	for _, entry := range entries {
		normalizedAddr := utils.NormalizeAddress(entry.Address)
		if normalizedAddr == userAddress {
			targetUserEarning = entry.TotalEarned
			break
		}
	}

	// Check if target user has any earnings
	if targetUserEarning == nil {
		return nil, fmt.Errorf("user %s has no claimable subsidies in vault %s for epoch %s", userAddress, vaultAddress, epochNumber)
	}

	mps.logger.Logf("INFO user %s has claimable amount: %s from cached snapshot for epoch %s", userAddress, targetUserEarning.String(), epochNumber)

	// Generate proof for the target user
	proof, merkleRoot, err := mps.proofGenerator.GenerateProof(entries, userAddress, targetUserEarning)
	if err != nil {
		return nil, fmt.Errorf("failed to generate merkle proof from cache: %w", err)
	}

	// Validate that the generated merkle root matches the cached root
	generatedRootHex := fmt.Sprintf("0x%x", merkleRoot[:])
	if generatedRootHex != snapshot.MerkleRoot {
		return nil, fmt.Errorf("generated merkle root %s does not match cached root %s for epoch %s",
			generatedRootHex, snapshot.MerkleRoot, epochNumber)
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
		MerkleRoot:      snapshot.MerkleRoot,
		MerkleProof:     proofStrings,
		GeneratedAt:     time.Now().Unix(),
		EpochNumber:     epochNumber,
		BlockNumber:     snapshot.BlockNumber,
		DataTimestamp:   snapshot.Timestamp,
	}, nil
}
