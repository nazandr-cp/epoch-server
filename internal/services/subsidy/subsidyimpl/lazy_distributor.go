package subsidyimpl

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/andrey/epoch-server/internal/infra/blockchain"
	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/andrey/epoch-server/internal/services/merkle/merkleimpl"
	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/go-pkgz/lgr"
)

type LazyDistributor struct {
	blockchainClient blockchain.BlockchainClient
	merkleService    merkle.Service
	subgraphClient   subgraph.SubgraphClient
	logger           lgr.L
}

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

func (d *LazyDistributor) Run(ctx context.Context, vaultId string) (*subsidy.DistributionResult, error) {
	return d.RunWithEpoch(ctx, vaultId, nil)
}

func (d *LazyDistributor) RunWithEpoch(ctx context.Context, vaultId string, epochNumber *big.Int) (*subsidy.DistributionResult, error) {
	if vaultId == "" {
		return nil, fmt.Errorf("vaultId cannot be empty")
	}

	d.logger.Logf("INFO starting lazy distributor for vault %s", vaultId)

	d.logger.Logf("DEBUG querying account subsidies for vault %s", vaultId)
	subsidies, err := d.subgraphClient.QueryAccountSubsidiesForVault(ctx, vaultId)
	if err != nil {
		d.logger.Logf("ERROR failed to get account subsidies for vault %s: %v", vaultId, err)
		return nil, fmt.Errorf("failed to get account subsidies: %w", err)
	}
	d.logger.Logf("DEBUG query completed successfully, returned %d subsidies", len(subsidies))

	d.logger.Logf("DEBUG found %d potential subsidies for vault %s", len(subsidies), vaultId)
	for i, subsidy := range subsidies {
		d.logger.Logf(
			"DEBUG subsidy[%d]: account=%s, secondsAccumulated=%s, lastEffectiveValue=%s, totalRewardsEarned=%s, updatedAt=%s",
			i,
			subsidy.Account.ID,
			subsidy.SecondsAccumulated,
			subsidy.LastEffectiveValue,
			subsidy.TotalRewardsEarned,
			subsidy.UpdatedAtTimestamp,
		)
	}

	if len(subsidies) == 0 {
		d.logger.Logf("INFO no subsidies found for vault %s, skipping distribution", vaultId)
		return &subsidy.DistributionResult{
			TotalSubsidies:    big.NewInt(0),
			AccountsProcessed: 0,
			MerkleRoot:        "",
		}, nil
	}

	entries, totalSubsidies, err := d.convertSubsidiesToEntries(subsidies)
	if err != nil {
		d.logger.Logf("ERROR failed to convert subsidies to entries: %v", err)
		return nil, fmt.Errorf("failed to convert subsidies to entries: %w", err)
	}

	if len(entries) == 0 {
		d.logger.Logf("INFO no valid entries found for vault %s, skipping distribution", vaultId)
		return &subsidy.DistributionResult{
			TotalSubsidies:    big.NewInt(0),
			AccountsProcessed: 0,
			MerkleRoot:        "",
		}, nil
	}

	merkleRoot, err := d.generateMerkleRoot(entries)
	if err != nil {
		d.logger.Logf("ERROR failed to generate merkle root: %v", err)
		return nil, fmt.Errorf("failed to generate merkle root: %w", err)
	}

	d.logger.Logf("INFO generated merkle root for vault %s: %x", vaultId, merkleRoot)
	d.logger.Logf("INFO total subsidies for vault %s: %s", vaultId, totalSubsidies.String())

	if epochNumber != nil {
		if err := d.saveSnapshot(ctx, vaultId, entries, merkleRoot, epochNumber); err != nil {
			d.logger.Logf("WARN failed to save merkle snapshot: %v", err)
		}
	}

	if err := d.updateMerkleRoot(ctx, vaultId, merkleRoot, totalSubsidies); err != nil {
		d.logger.Logf("ERROR failed to update merkle root on blockchain: %v", err)
		return nil, fmt.Errorf("failed to update merkle root on blockchain: %w", err)
	}

	d.logger.Logf("INFO successfully completed lazy distributor for vault %s", vaultId)
	return &subsidy.DistributionResult{
		TotalSubsidies:    totalSubsidies,
		AccountsProcessed: len(entries),
		MerkleRoot:        fmt.Sprintf("%x", merkleRoot),
	}, nil
}

func (d *LazyDistributor) convertSubsidiesToEntries(
	subsidies []subgraph.AccountSubsidy,
) ([]merkle.Entry, *big.Int, error) {
	entries := make([]merkle.Entry, 0, len(subsidies))
	totalSubsidies := big.NewInt(0)
	currentTimestamp := time.Now().Unix()

	for _, subsidy := range subsidies {
		amount, ok := new(big.Int).SetString(subsidy.TotalRewardsEarned, 10)
		if !ok || amount.Sign() <= 0 {
			calculatedAmount, err := d.calculateTotalEarned(subsidy, currentTimestamp)
			if err != nil {
				d.logger.Logf(
					"WARN failed to calculate total earned for account %s: %v, using totalRewardsEarned=%s",
					subsidy.Account.ID,
					err,
					subsidy.TotalRewardsEarned,
				)
				continue
			}
			amount = calculatedAmount
			d.logger.Logf(
				"DEBUG calculated earnings for account %s: %s (secondsAccumulated=%s, lastEffectiveValue=%s)",
				subsidy.Account.ID,
				amount.String(),
				subsidy.SecondsAccumulated,
				subsidy.LastEffectiveValue,
			)
		}

		if amount.Sign() <= 0 {
			d.logger.Logf(
				"DEBUG skipping account %s with zero earnings (amount=%s)",
				subsidy.Account.ID,
				amount.String(),
			)
			continue
		}

		entry := merkle.Entry{
			Address:     subsidy.Account.ID,
			TotalEarned: amount,
		}

		entries = append(entries, entry)
		totalSubsidies.Add(totalSubsidies, amount)
	}

	d.logger.Logf("INFO processed %d subsidies, generated %d valid entries", len(subsidies), len(entries))
	return entries, totalSubsidies, nil
}

func (d *LazyDistributor) generateMerkleRoot(entries []merkle.Entry) ([32]byte, error) {
	merkleImpl, ok := d.merkleService.(*merkleimpl.Service)
	if !ok {
		return [32]byte{}, fmt.Errorf("merkle service is not the expected implementation type")
	}

	root := merkleImpl.BuildMerkleRootFromEntries(entries)
	return root, nil
}

func (d *LazyDistributor) calculateTotalEarned(subsidy subgraph.AccountSubsidy, endTimestamp int64) (*big.Int, error) {
	secondsAccumulated, ok := new(big.Int).SetString(subsidy.SecondsAccumulated, 10)
	if !ok {
		return nil, fmt.Errorf("invalid secondsAccumulated: %s", subsidy.SecondsAccumulated)
	}

	lastEffectiveValue, ok := new(big.Int).SetString(subsidy.LastEffectiveValue, 10)
	if !ok {
		return nil, fmt.Errorf("invalid lastEffectiveValue: %s", subsidy.LastEffectiveValue)
	}

	updatedAtTimestamp, err := strconv.ParseInt(subsidy.UpdatedAtTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid updatedAtTimestamp: %s", subsidy.UpdatedAtTimestamp)
	}

	deltaT := endTimestamp - updatedAtTimestamp
	extraSeconds := new(big.Int).Mul(big.NewInt(deltaT), lastEffectiveValue)
	newTotalSeconds := new(big.Int).Add(secondsAccumulated, extraSeconds)

	conversionRate := big.NewInt(1000000000000000000)
	totalEarned := new(big.Int).Div(newTotalSeconds, conversionRate)
	return totalEarned, nil
}

func (d *LazyDistributor) updateMerkleRoot(
	ctx context.Context,
	vaultId string,
	merkleRoot [32]byte,
	totalSubsidies *big.Int,
) error {
	return d.blockchainClient.UpdateMerkleRootAndWaitForConfirmation(ctx, vaultId, merkleRoot, totalSubsidies)
}

func (d *LazyDistributor) saveSnapshot(
	ctx context.Context,
	vaultId string,
	entries []merkle.Entry,
	merkleRoot [32]byte,
	epochNumber *big.Int,
) error {
	merkleEntries := make([]merkle.MerkleEntry, len(entries))
	for i, entry := range entries {
		merkleEntries[i] = merkle.MerkleEntry{
			Address:     entry.Address,
			TotalEarned: entry.TotalEarned,
		}
	}

	snapshot := merkle.MerkleSnapshot{
		VaultID:     vaultId,
		MerkleRoot:  fmt.Sprintf("%x", merkleRoot),
		Entries:     merkleEntries,
		EpochNumber: epochNumber,
	}

	merkleImpl, ok := d.merkleService.(*merkleimpl.Service)
	if !ok {
		return fmt.Errorf("merkle service is not the expected implementation type")
	}

	if err := merkleImpl.SaveSnapshot(ctx, epochNumber, snapshot); err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	d.logger.Logf("INFO saved merkle snapshot for vault %s, epoch %s with %d entries",
		vaultId, epochNumber.String(), len(entries))
	return nil
}
