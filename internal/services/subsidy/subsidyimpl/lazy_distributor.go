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

// LazyDistributor handles subsidy distribution by generating merkle roots
type LazyDistributor struct {
	blockchainClient blockchain.BlockchainClient
	merkleService    merkle.Service
	subgraphClient   subgraph.SubgraphClient
	logger           lgr.L
}

// NewLazyDistributor creates a new LazyDistributor
func NewLazyDistributor(
	blockchainClient blockchain.BlockchainClient,
	merkleService merkle.Service,
	subgraphClient subgraph.SubgraphClient,
	logger lgr.L,
) *LazyDistributor {
	return &LazyDistributor{
		blockchainClient: blockchainClient,
		merkleService:    merkleService,
		subgraphClient:   subgraphClient,
		logger:           logger,
	}
}

// Run executes subsidy distribution for a vault
func (d *LazyDistributor) Run(ctx context.Context, vaultId string) error {
	if vaultId == "" {
		return fmt.Errorf("vaultId cannot be empty")
	}

	d.logger.Logf("INFO starting lazy distributor for vault %s", vaultId)

	subsidies, err := d.subgraphClient.QueryAccountSubsidiesForVault(ctx, vaultId)
	if err != nil {
		d.logger.Logf("ERROR failed to get account subsidies for vault %s: %v", vaultId, err)
		return fmt.Errorf("failed to get account subsidies: %w", err)
	}

	if len(subsidies) == 0 {
		d.logger.Logf("INFO no subsidies found for vault %s, skipping distribution", vaultId)
		return nil
	}

	entries, totalSubsidies, err := d.convertSubsidiesToEntries(subsidies)
	if err != nil {
		d.logger.Logf("ERROR failed to convert subsidies to entries: %v", err)
		return fmt.Errorf("failed to convert subsidies to entries: %w", err)
	}

	if len(entries) == 0 {
		d.logger.Logf("INFO no valid entries found for vault %s, skipping distribution", vaultId)
		return nil
	}

	merkleRoot, err := d.generateMerkleRoot(entries)
	if err != nil {
		d.logger.Logf("ERROR failed to generate merkle root: %v", err)
		return fmt.Errorf("failed to generate merkle root: %w", err)
	}

	d.logger.Logf("INFO generated merkle root for vault %s: %x", vaultId, merkleRoot)
	d.logger.Logf("INFO total subsidies for vault %s: %s", vaultId, totalSubsidies.String())

	if err := d.updateMerkleRoot(ctx, vaultId, merkleRoot, totalSubsidies); err != nil {
		d.logger.Logf("ERROR failed to update merkle root on blockchain: %v", err)
		return fmt.Errorf("failed to update merkle root on blockchain: %w", err)
	}

	d.logger.Logf("INFO successfully completed lazy distributor for vault %s", vaultId)
	return nil
}

// convertSubsidiesToEntries converts subsidies to merkle entries
func (d *LazyDistributor) convertSubsidiesToEntries(
	subsidies []subgraph.AccountSubsidy,
) ([]merkle.Entry, *big.Int, error) {
	entries := make([]merkle.Entry, 0, len(subsidies))
	totalSubsidies := big.NewInt(0)

	for _, subsidy := range subsidies {
		amount, ok := new(big.Int).SetString(subsidy.TotalRewardsEarned, 10)
		if !ok {
			d.logger.Logf(
				"WARN invalid total rewards earned amount for account %s: %s",
				subsidy.Account.ID,
				subsidy.TotalRewardsEarned,
			)
			continue
		}

		if amount.Sign() <= 0 {
			continue
		}

		entry := merkle.Entry{
			Address:     subsidy.Account.ID,
			TotalEarned: amount,
		}

		entries = append(entries, entry)
		totalSubsidies.Add(totalSubsidies, amount)
	}

	return entries, totalSubsidies, nil
}

// generateMerkleRoot generates merkle root from entries
func (d *LazyDistributor) generateMerkleRoot(entries []merkle.Entry) ([32]byte, error) {
	merkleImpl, ok := d.merkleService.(*merkleimpl.Service)
	if !ok {
		return [32]byte{}, fmt.Errorf("merkle service is not the expected implementation type")
	}

	root := merkleImpl.BuildMerkleRootFromEntries(entries)
	return root, nil
}

// updateMerkleRoot updates merkle root on blockchain
func (d *LazyDistributor) updateMerkleRoot(
	ctx context.Context,
	vaultId string,
	merkleRoot [32]byte,
	totalSubsidies *big.Int,
) error {
	return d.blockchainClient.UpdateMerkleRootAndWaitForConfirmation(ctx, vaultId, merkleRoot, totalSubsidies)
}
