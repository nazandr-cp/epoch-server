package subsidyimpl

import (
	"context"
	"fmt"
	"math/big"

	"github.com/andrey/epoch-server/internal/infra/blockchain"
	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/andrey/epoch-server/internal/services/merkle/merkleimpl"
	"github.com/go-pkgz/lgr"
)

// LazyDistributor implements the subsidy.LazyDistributor interface
// It handles the actual distribution of subsidies by generating merkle roots
// and updating the blockchain contract
type LazyDistributor struct {
	blockchainClient *blockchain.Client
	merkleService    merkle.Service
	subgraphClient   *subgraph.Client
	logger           lgr.L
}

// NewLazyDistributor creates a new LazyDistributor instance
func NewLazyDistributor(
	blockchainClient *blockchain.Client,
	merkleService merkle.Service,
	subgraphClient *subgraph.Client,
	logger lgr.L,
) *LazyDistributor {
	return &LazyDistributor{
		blockchainClient: blockchainClient,
		merkleService:    merkleService,
		subgraphClient:   subgraphClient,
		logger:           logger,
	}
}

// Run executes the subsidy distribution for a given vault
func (d *LazyDistributor) Run(ctx context.Context, vaultId string) error {
	if vaultId == "" {
		return fmt.Errorf("vaultId cannot be empty")
	}

	d.logger.Logf("INFO starting lazy distributor for vault %s", vaultId)

	// Get account subsidies for the vault from subgraph
	subsidies, err := d.subgraphClient.QueryAccountSubsidiesForVault(ctx, vaultId)
	if err != nil {
		d.logger.Logf("ERROR failed to get account subsidies for vault %s: %v", vaultId, err)
		return fmt.Errorf("failed to get account subsidies: %w", err)
	}

	if len(subsidies) == 0 {
		d.logger.Logf("INFO no subsidies found for vault %s, skipping distribution", vaultId)
		return nil
	}

	// Convert subsidies to merkle entries
	entries, totalSubsidies, err := d.convertSubsidiesToEntries(subsidies)
	if err != nil {
		d.logger.Logf("ERROR failed to convert subsidies to entries: %v", err)
		return fmt.Errorf("failed to convert subsidies to entries: %w", err)
	}

	if len(entries) == 0 {
		d.logger.Logf("INFO no valid entries found for vault %s, skipping distribution", vaultId)
		return nil
	}

	// Generate merkle root from entries
	merkleRoot, err := d.generateMerkleRoot(entries)
	if err != nil {
		d.logger.Logf("ERROR failed to generate merkle root: %v", err)
		return fmt.Errorf("failed to generate merkle root: %w", err)
	}

	d.logger.Logf("INFO generated merkle root for vault %s: %x", vaultId, merkleRoot)
	d.logger.Logf("INFO total subsidies for vault %s: %s", vaultId, totalSubsidies.String())

	// Update merkle root on blockchain via DebtSubsidizer contract
	if err := d.updateMerkleRoot(ctx, vaultId, merkleRoot, totalSubsidies); err != nil {
		d.logger.Logf("ERROR failed to update merkle root on blockchain: %v", err)
		return fmt.Errorf("failed to update merkle root on blockchain: %w", err)
	}

	d.logger.Logf("INFO successfully completed lazy distributor for vault %s", vaultId)
	return nil
}

// convertSubsidiesToEntries converts subgraph subsidies to merkle entries
func (d *LazyDistributor) convertSubsidiesToEntries(subsidies []subgraph.AccountSubsidy) ([]merkle.Entry, *big.Int, error) {
	entries := make([]merkle.Entry, 0, len(subsidies))
	totalSubsidies := big.NewInt(0)

	for _, subsidy := range subsidies {
		// Parse the total rewards earned amount
		amount, ok := new(big.Int).SetString(subsidy.TotalRewardsEarned, 10)
		if !ok {
			d.logger.Logf("WARN invalid total rewards earned amount for account %s: %s", subsidy.Account.ID, subsidy.TotalRewardsEarned)
			continue
		}

		// Skip zero amounts
		if amount.Sign() <= 0 {
			continue
		}

		// Create merkle entry
		entry := merkle.Entry{
			Address:     subsidy.Account.ID,
			TotalEarned: amount,
		}

		entries = append(entries, entry)
		totalSubsidies.Add(totalSubsidies, amount)
	}

	return entries, totalSubsidies, nil
}

// generateMerkleRoot generates the merkle root from entries using the merkle service
func (d *LazyDistributor) generateMerkleRoot(entries []merkle.Entry) ([32]byte, error) {
	// Use the merkle service to build the merkle root
	// We need to access the internal implementation to get the BuildMerkleRootFromEntries method
	merkleImpl, ok := d.merkleService.(*merkleimpl.Service)
	if !ok {
		return [32]byte{}, fmt.Errorf("merkle service is not the expected implementation type")
	}

	root := merkleImpl.BuildMerkleRootFromEntries(entries)
	return root, nil
}

// updateMerkleRoot updates the merkle root on the blockchain via DebtSubsidizer contract
func (d *LazyDistributor) updateMerkleRoot(ctx context.Context, vaultId string, merkleRoot [32]byte, totalSubsidies *big.Int) error {
	// For now, we'll use the existing subsidizer client pattern
	// In a real implementation, we would need to add a method to the blockchain client
	// to handle DebtSubsidizer contract calls

	// Create a subsidizer client using the same config as the main blockchain client
	// This is a temporary approach until we integrate subsidizer into the main client
	subsidizer := blockchain.NewSubsidizerClient(d.logger)

	// Call the subsidizer to update merkle root
	return subsidizer.UpdateMerkleRootAndWaitForConfirmation(ctx, vaultId, merkleRoot, totalSubsidies)
}
