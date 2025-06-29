package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/andrey/epoch-server/internal/clients/epoch"
	"github.com/andrey/epoch-server/internal/clients/storage"
	"github.com/go-pkgz/lgr"
)

type AccountSubsidyPerCollection struct {
	Account            Account `json:"account"`
	SecondsAccumulated string  `json:"secondsAccumulated"`
	SecondsClaimed     string  `json:"secondsClaimed"`
	LastEffectiveValue string  `json:"lastEffectiveValue"`
	UpdatedAtTimestamp string  `json:"updatedAtTimestamp"`
}

type Account struct {
	ID string `json:"id"`
}

type EpochManagerClient interface {
	Current() epoch.EpochInfo
	FinalizeEpoch() error
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
	logger               lgr.L
}

func NewLazyDistributor(
	graphClient GraphClient,
	epochManagerClient EpochManagerClient,
	debtSubsidizerClient DebtSubsidizerClient,
	storageClient StorageClient,
	logger lgr.L,
) *LazyDistributor {
	return &LazyDistributor{
		graphClient:          graphClient,
		epochManagerClient:   epochManagerClient,
		debtSubsidizerClient: debtSubsidizerClient,
		storageClient:        storageClient,
		logger:               logger,
	}
}

func (ld *LazyDistributor) Run(ctx context.Context, vaultId string) error {
	ld.logger.Logf("INFO starting lazy distribution for vault %s", vaultId)

	subsidies, err := ld.queryLazySubsidies(ctx, vaultId)
	if err != nil {
		return fmt.Errorf("failed to query lazy subsidies: %w", err)
	}

	epochEnd := ld.epochManagerClient.Current().EndTime

	entries := make([]storage.MerkleEntry, 0, len(subsidies))
	for _, subsidy := range subsidies {
		totalEarned, err := ld.calculateTotalEarned(subsidy, epochEnd)
		if err != nil {
			return fmt.Errorf("failed to calculate total earned for account %s: %w", subsidy.Account.ID, err)
		}

		entries = append(entries, storage.MerkleEntry{
			Address:     subsidy.Account.ID,
			TotalEarned: totalEarned,
		})
	}

	merkleRoot := ld.buildMerkleRoot(entries)

	snapshot := storage.MerkleSnapshot{
		Entries:    entries,
		MerkleRoot: fmt.Sprintf("0x%x", merkleRoot),
		Timestamp:  time.Now().Unix(),
		VaultID:    vaultId,
	}

	if err := ld.storageClient.SaveSnapshot(ctx, snapshot); err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	var rootBytes [32]byte
	copy(rootBytes[:], merkleRoot[:])

	if err := ld.debtSubsidizerClient.UpdateMerkleRootAndWaitForConfirmation(ctx, vaultId, rootBytes); err != nil {
		return fmt.Errorf("failed to update merkle root and wait for confirmation: %w", err)
	}

	if err := ld.epochManagerClient.FinalizeEpoch(); err != nil {
		return fmt.Errorf("failed to finalize epoch: %w", err)
	}

	ld.logger.Logf("INFO completed lazy distribution for vault %s with %d entries", vaultId, len(entries))
	return nil
}

func (ld *LazyDistributor) queryLazySubsidies(ctx context.Context, vaultId string) ([]AccountSubsidyPerCollection, error) {
	query := `
		query GetLazySubsidies($vaultId: ID!, $first: Int!, $skip: Int!) {
			accountSubsidiesPerCollections(
				where: { vault: $vaultId, secondsAccumulated_gt: "0" }
				orderBy: account
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
		"vaultId": vaultId,
	}

	var response struct {
		Data struct {
			AccountSubsidiesPerCollections []AccountSubsidyPerCollection `json:"accountSubsidiesPerCollections"`
		} `json:"data"`
	}

	if err := ld.graphClient.ExecutePaginatedQuery(ctx, query, variables, "accountSubsidiesPerCollections", &response); err != nil {
		return nil, fmt.Errorf("failed to execute paginated GraphQL query: %w", err)
	}

	return response.Data.AccountSubsidiesPerCollections, nil
}

func (ld *LazyDistributor) calculateTotalEarned(subsidy AccountSubsidyPerCollection, epochEnd int64) (*big.Int, error) {
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

	deltaT := epochEnd - updatedAtTimestamp
	extraSeconds := new(big.Int).Mul(big.NewInt(deltaT), lastEffectiveValue)
	newTotalSeconds := new(big.Int).Add(secondsAccumulated, extraSeconds)

	totalEarned := ld.secondsToTokens(newTotalSeconds)
	return totalEarned, nil
}

func (ld *LazyDistributor) secondsToTokens(seconds *big.Int) *big.Int {
	conversionRate := big.NewInt(1000000000000000000)
	return new(big.Int).Div(seconds, conversionRate)
}

func (ld *LazyDistributor) buildMerkleRoot(entries []storage.MerkleEntry) [32]byte {
	if len(entries) == 0 {
		return [32]byte{}
	}

	hashes := make([][32]byte, len(entries))
	for i, entry := range entries {
		data := fmt.Sprintf("%s:%s", entry.Address, entry.TotalEarned.String())
		hashes[i] = sha256.Sum256([]byte(data))
	}

	for len(hashes) > 1 {
		var nextLevel [][32]byte
		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {
				combined := append(hashes[i][:], hashes[i+1][:]...)
				nextLevel = append(nextLevel, sha256.Sum256(combined))
			} else {
				nextLevel = append(nextLevel, hashes[i])
			}
		}
		hashes = nextLevel
	}

	return hashes[0]
}
