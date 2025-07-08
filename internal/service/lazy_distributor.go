package service

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/andrey/epoch-server/internal/clients/epoch"
	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/andrey/epoch-server/internal/clients/storage"
	"github.com/andrey/epoch-server/internal/config"
	"github.com/andrey/epoch-server/internal/merkle"
	"github.com/go-pkgz/lgr"
)


type EpochManagerClient interface {
	Current() epoch.EpochInfo
	GetCurrentEpochId(ctx context.Context) (*big.Int, error)
	FinalizeEpoch() error
	UpdateExchangeRate(ctx context.Context, lendingManagerAddress string) error
	AllocateYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string) error
	AllocateCumulativeYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string, amount *big.Int) error
	EndEpochWithSubsidies(ctx context.Context, epochId *big.Int, vaultAddress string, merkleRoot [32]byte, subsidiesDistributed *big.Int) error
}

type DebtSubsidizerClient interface {
	UpdateMerkleRoot(ctx context.Context, vaultId string, root [32]byte) error
	UpdateMerkleRootAndWaitForConfirmation(ctx context.Context, vaultId string, root [32]byte) error
}

type StorageClient interface {
	SaveSnapshot(ctx context.Context, snapshot storage.MerkleSnapshot) error
}

type LazyDistributor struct {
	graphClient          GraphClient
	epochManagerClient   EpochManagerClient
	debtSubsidizerClient DebtSubsidizerClient
	storageClient        StorageClient
	proofGenerator       *merkle.ProofGenerator
	calculator           *merkle.Calculator
	logger               lgr.L
	config               *config.Config
}

func NewLazyDistributor(
	graphClient GraphClient,
	epochManagerClient EpochManagerClient,
	debtSubsidizerClient DebtSubsidizerClient,
	storageClient StorageClient,
	logger lgr.L,
	cfg *config.Config,
) *LazyDistributor {
	return &LazyDistributor{
		graphClient:          graphClient,
		epochManagerClient:   epochManagerClient,
		debtSubsidizerClient: debtSubsidizerClient,
		storageClient:        storageClient,
		proofGenerator:       merkle.NewProofGeneratorWithDependencies(graphClient, logger),
		calculator:           merkle.NewCalculator(),
		logger:               logger,
		config:               cfg,
	}
}

func (ld *LazyDistributor) Run(ctx context.Context, vaultId string) error {
	ld.logger.Logf("INFO starting lazy distribution for vault %s", vaultId)

	// Query account subsidies for the vault
	subsidies, err := ld.queryAccountSubsidiesForVault(ctx, vaultId)
	if err != nil {
		return fmt.Errorf("failed to query account subsidies: %w", err)
	}

	ld.logger.Logf("INFO found %d account subsidies for vault %s", len(subsidies), vaultId)

	// Use the unified merkle package directly
	result, err := ld.proofGenerator.GenerateTreeFromSubsidies(ctx, vaultId, subsidies)
	if err != nil {
		return fmt.Errorf("failed to generate merkle tree: %w", err)
	}

	// Convert to storage.MerkleEntry format
	entries := make([]storage.MerkleEntry, len(result.Entries))
	for i, entry := range result.Entries {
		entries[i] = storage.MerkleEntry{
			Address:     entry.Address,
			TotalEarned: entry.TotalEarned,
		}
	}

	ld.logger.Logf("INFO generated merkle tree with %d entries using processing time %d", len(entries), result.Timestamp)

	merkleRoot := result.MerkleRoot

	snapshot := storage.MerkleSnapshot{
		Entries:    entries,
		MerkleRoot: fmt.Sprintf("0x%x", merkleRoot),
		Timestamp:  result.Timestamp, // Use the consistent timestamp from the merkle tree generation
		VaultID:    vaultId,
	}

	if err := ld.storageClient.SaveSnapshot(ctx, snapshot); err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	var rootBytes [32]byte
	copy(rootBytes[:], merkleRoot[:])

	ld.logger.Logf("INFO updating merkle root for vault %s: %x", vaultId, rootBytes)
	if err := ld.debtSubsidizerClient.UpdateMerkleRootAndWaitForConfirmation(ctx, vaultId, rootBytes); err != nil {
		ld.logger.Logf("ERROR failed to update merkle root for vault %s: %v", vaultId, err)
		return fmt.Errorf("failed to call updateMerkleRoot: %w", err)
	}

	// Calculate total subsidies distributed using the unified calculator
	totalSubsidies := big.NewInt(0)
	for _, entry := range entries {
		totalSubsidies.Add(totalSubsidies, entry.TotalEarned)
	}

	// Get current epoch ID from the epoch manager
	epochId, err := ld.epochManagerClient.GetCurrentEpochId(ctx)
	if err != nil {
		ld.logger.Logf("ERROR failed to get current epoch ID: %v", err)
		return fmt.Errorf("failed to get current epoch ID: %w", err)
	}
	ld.logger.Logf("INFO using epoch ID %s for subsidy distribution", epochId.String())

	// Update exchange rate to ensure we have the latest yield calculations
	lendingManagerAddress := ld.config.Contracts.LendingManager
	ld.logger.Logf("INFO updating exchange rate for LendingManager %s", lendingManagerAddress)
	if err := ld.epochManagerClient.UpdateExchangeRate(ctx, lendingManagerAddress); err != nil {
		ld.logger.Logf("ERROR failed to update exchange rate for LendingManager %s: %v", lendingManagerAddress, err)
		return fmt.Errorf("failed to call updateExchangeRate: %w", err)
	}

	// Validate merkle tree total matches calculated subsidies
	ld.logger.Logf("INFO merkle tree contains %d entries with total subsidies: %s", len(entries), totalSubsidies.String())
	
	// Log breakdown for validation
	for i, entry := range entries {
		ld.logger.Logf("INFO entry %d: %s -> %s tokens", i, entry.Address, entry.TotalEarned.String())
	}
	
	// Pre-validate that the vault has sufficient yield to cover the merkle tree total
	// Note: This would require adding a contract call to check validateMerkleTreeAllocation
	// For now, we rely on the allocation function's built-in validation
	ld.logger.Logf("INFO proceeding with allocation of %s tokens for %d users", totalSubsidies.String(), len(entries))
	
	// Allocate exact cumulative yield amount needed for all subsidies
	ld.logger.Logf("INFO allocating cumulative yield %s to epoch %s for vault %s", totalSubsidies.String(), epochId.String(), vaultId)
	if err := ld.epochManagerClient.AllocateCumulativeYieldToEpoch(ctx, epochId, vaultId, totalSubsidies); err != nil {
		ld.logger.Logf("ERROR failed to allocate cumulative yield %s to epoch %s for vault %s: %v", totalSubsidies.String(), epochId.String(), vaultId, err)
		return fmt.Errorf("failed to call allocateCumulativeYieldToEpoch: %w", err)
	}

	if err := ld.epochManagerClient.EndEpochWithSubsidies(ctx, epochId, vaultId, rootBytes, totalSubsidies); err != nil {
		ld.logger.Logf("ERROR failed to end epoch %s with subsidies for vault %s: %v", epochId.String(), vaultId, err)
		return fmt.Errorf("failed to call endEpochWithSubsidies: %w", err)
	}

	ld.logger.Logf("INFO completed lazy distribution for vault %s with %d entries", vaultId, len(entries))
	return nil
}

// queryAccountSubsidiesForVault queries all account subsidies for a specific vault
func (ld *LazyDistributor) queryAccountSubsidiesForVault(ctx context.Context, vaultId string) ([]graph.AccountSubsidy, error) {
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

	if err := ld.graphClient.ExecutePaginatedQuery(ctx, query, variables, "accountSubsidies", &response); err != nil {
		return nil, fmt.Errorf("failed to execute paginated GraphQL query: %w", err)
	}

	return response.Data.AccountSubsidies, nil
}

