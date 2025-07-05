package service

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"time"

	"github.com/andrey/epoch-server/internal/clients/epoch"
	"github.com/andrey/epoch-server/internal/clients/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-pkgz/lgr"
)

// AccountSubsidy represents the new consolidated account subsidy entity
type AccountSubsidy struct {
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
	GetCurrentEpochId(ctx context.Context) (*big.Int, error)
	FinalizeEpoch() error
	UpdateExchangeRate(ctx context.Context, lendingManagerAddress string) error
	AllocateYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string) error
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

	ld.logger.Logf("INFO found %d subsidies from query", len(subsidies))

	epochEnd := ld.epochManagerClient.Current().EndTime
	ld.logger.Logf("INFO epoch end time: %d", epochEnd)

	entries := make([]storage.MerkleEntry, 0, len(subsidies))
	for _, subsidy := range subsidies {
		ld.logger.Logf("INFO processing account %s: secondsAccumulated=%s, lastEffectiveValue=%s, updatedAtTimestamp=%s", 
			subsidy.Account.ID, subsidy.SecondsAccumulated, subsidy.LastEffectiveValue, subsidy.UpdatedAtTimestamp)
		
		totalEarned, err := ld.calculateTotalEarned(subsidy, epochEnd)
		if err != nil {
			return fmt.Errorf("failed to calculate total earned for account %s: %w", subsidy.Account.ID, err)
		}

		ld.logger.Logf("INFO calculated totalEarned for account %s: %s", subsidy.Account.ID, totalEarned.String())

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

	ld.logger.Logf("INFO updating merkle root for vault %s: %x", vaultId, rootBytes)
	if err := ld.debtSubsidizerClient.UpdateMerkleRootAndWaitForConfirmation(ctx, vaultId, rootBytes); err != nil {
		ld.logger.Logf("ERROR failed to update merkle root for vault %s: %v", vaultId, err)
		return fmt.Errorf("failed to call updateMerkleRoot: %w", err)
	}

	// Calculate total subsidies distributed
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
	lendingManagerAddress := "0x64Bd8C3294956E039EDf1a4058b6588de3731248"
	ld.logger.Logf("INFO updating exchange rate for LendingManager %s", lendingManagerAddress)
	if err := ld.epochManagerClient.UpdateExchangeRate(ctx, lendingManagerAddress); err != nil {
		ld.logger.Logf("ERROR failed to update exchange rate for LendingManager %s: %v", lendingManagerAddress, err)
		return fmt.Errorf("failed to call updateExchangeRate: %w", err)
	}

	// Allocate yield to epoch before ending it with subsidies
	ld.logger.Logf("INFO allocating yield to epoch %s for vault %s", epochId.String(), vaultId)
	if err := ld.epochManagerClient.AllocateYieldToEpoch(ctx, epochId, vaultId); err != nil {
		ld.logger.Logf("ERROR failed to allocate yield to epoch %s for vault %s: %v", epochId.String(), vaultId, err)
		return fmt.Errorf("failed to call allocateYieldToEpoch: %w", err)
	}

	ld.logger.Logf("INFO ending epoch %s with subsidies for vault %s (total subsidies: %s)", epochId.String(), vaultId, totalSubsidies.String())
	if err := ld.epochManagerClient.EndEpochWithSubsidies(ctx, epochId, vaultId, rootBytes, totalSubsidies); err != nil {
		ld.logger.Logf("ERROR failed to end epoch %s with subsidies for vault %s: %v", epochId.String(), vaultId, err)
		return fmt.Errorf("failed to call endEpochWithSubsidies: %w", err)
	}

	ld.logger.Logf("INFO completed lazy distribution for vault %s with %d entries", vaultId, len(entries))
	return nil
}

func (ld *LazyDistributor) queryLazySubsidies(ctx context.Context, vaultId string) ([]AccountSubsidy, error) {
	query := `
		query GetLazySubsidies($vaultId: ID!, $first: Int!, $skip: Int!) {
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
		"vaultId": vaultId,
	}

	var response struct {
		Data struct {
			AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
		} `json:"data"`
	}

	if err := ld.graphClient.ExecutePaginatedQuery(ctx, query, variables, "accountSubsidies", &response); err != nil {
		return nil, fmt.Errorf("failed to execute paginated GraphQL query: %w", err)
	}

	return response.Data.AccountSubsidies, nil
}

func (ld *LazyDistributor) calculateTotalEarned(subsidy AccountSubsidy, epochEnd int64) (*big.Int, error) {
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

// buildMerkleRoot creates a Merkle tree compatible with OpenZeppelin's MerkleProof library
// Uses the same keccak256 hashing and abi.encodePacked format as the Solidity contracts
func (ld *LazyDistributor) buildMerkleRoot(entries []storage.MerkleEntry) [32]byte {
	if len(entries) == 0 {
		return [32]byte{}
	}

	// Sort entries deterministically by address to ensure consistent Merkle roots
	sortedEntries := make([]storage.MerkleEntry, len(entries))
	copy(sortedEntries, entries)
	sort.Slice(sortedEntries, func(i, j int) bool {
		return sortedEntries[i].Address < sortedEntries[j].Address
	})

	// Generate leaf hashes using keccak256 and abi.encodePacked format
	hashes := make([][32]byte, len(sortedEntries))
	for i, entry := range sortedEntries {
		// Create leaf using same format as Solidity: keccak256(abi.encodePacked(recipient, newTotal))
		leaf := ld.createLeafHash(entry.Address, entry.TotalEarned)
		hashes[i] = leaf
	}

	// Build tree bottom-up using OpenZeppelin-compatible logic
	for len(hashes) > 1 {
		var nextLevel [][32]byte
		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {
				// Sort pair to match OpenZeppelin's ordering
				left, right := hashes[i], hashes[i+1]
				if !ld.isLeftSmaller(left, right) {
					left, right = right, left
				}
				// Hash the sorted pair
				combined := append(left[:], right[:]...)
				nextLevel = append(nextLevel, crypto.Keccak256Hash(combined))
			} else {
				// Odd number of nodes, promote the last one
				nextLevel = append(nextLevel, hashes[i])
			}
		}
		hashes = nextLevel
	}

	return hashes[0]
}

// createLeafHash creates a leaf hash compatible with Solidity's abi.encodePacked(recipient, newTotal)
func (ld *LazyDistributor) createLeafHash(address string, amount *big.Int) [32]byte {
	// Convert address string to common.Address
	addr := common.HexToAddress(address)

	// Create packed encoding: address (20 bytes) + amount (32 bytes)
	packed := make([]byte, 0, 52)
	packed = append(packed, addr.Bytes()...)

	// Convert amount to 32-byte representation (big-endian)
	amountBytes := make([]byte, 32)
	amount.FillBytes(amountBytes)
	packed = append(packed, amountBytes...)

	// Hash using keccak256
	return crypto.Keccak256Hash(packed)
}

// isLeftSmaller determines if left hash should come before right hash in OpenZeppelin ordering
func (ld *LazyDistributor) isLeftSmaller(left, right [32]byte) bool {
	for i := 0; i < 32; i++ {
		if left[i] < right[i] {
			return true
		}
		if left[i] > right[i] {
			return false
		}
	}
	return false // Equal hashes, doesn't matter which comes first
}
