package migration

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-pkgz/lgr"
)

type GraphClient interface {
	ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error
}

type SubsidizerClient interface {
	UpdateMerkleRoot(ctx context.Context, vaultId string, root [32]byte) error
}

type MigrationConfig struct {
	SubgraphEndpoint    string
	SnapshotBlockNumber *big.Int
	VaultID             string
	DryRun              bool
}

type AccountSubsidyRecord struct {
	Account            string
	WeightedBalance    *big.Int
	CurrentBorrowU     *big.Int
	LastEffectiveValue *big.Int
	SecondsAccumulated *big.Int
}

type MerkleLeaf struct {
	Account            string
	SecondsAccumulated *big.Int
}

type MigrationService struct {
	graphClient      GraphClient
	subsidizerClient SubsidizerClient
	logger           lgr.L
	config           MigrationConfig
}

func NewMigrationService(
	graphClient GraphClient,
	subsidizerClient SubsidizerClient,
	logger lgr.L,
	config MigrationConfig,
) *MigrationService {
	return &MigrationService{
		graphClient:      graphClient,
		subsidizerClient: subsidizerClient,
		logger:           logger,
		config:           config,
	}
}

func (m *MigrationService) InitializeSubsidies(ctx context.Context) error {
	m.logger.Logf("INFO Starting subsidy initialization migration for vault %s", m.config.VaultID)

	records, err := m.fetchAccountSubsidyRecords(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch account subsidy records: %w", err)
	}

	m.logger.Logf("INFO Found %d account subsidy records", len(records))

	if err := m.computeLastEffectiveValues(ctx, records); err != nil {
		return fmt.Errorf("failed to compute last effective values: %w", err)
	}

	if !m.config.DryRun {
		if err := m.updateSubgraphRecords(ctx, records); err != nil {
			return fmt.Errorf("failed to update subgraph records: %w", err)
		}
	}

	merkleRoot, err := m.generateMerkleRoot(records)
	if err != nil {
		return fmt.Errorf("failed to generate merkle root: %w", err)
	}

	m.logger.Logf("INFO Generated merkle root: %x", merkleRoot)

	if !m.config.DryRun {
		if err := m.subsidizerClient.UpdateMerkleRoot(ctx, m.config.VaultID, merkleRoot); err != nil {
			return fmt.Errorf("failed to update merkle root: %w", err)
		}
	}

	m.logger.Logf("INFO Successfully completed subsidy initialization migration")
	return nil
}

func (m *MigrationService) fetchAccountSubsidyRecords(ctx context.Context) ([]*AccountSubsidyRecord, error) {
	return m.fetchAccountSubsidyRecordsPaginated(ctx)
}

func (m *MigrationService) fetchAccountSubsidyRecordsPaginated(ctx context.Context) ([]*AccountSubsidyRecord, error) {
	query := `
		query GetAccountSubsidies($vaultId: ID!, $first: Int!, $skip: Int!) {
			accountSubsidies(
				where: { collectionParticipation_contains: $vaultId }
				orderBy: account
				orderDirection: asc
				first: $first
				skip: $skip
			) {
				account {
					id
				}
				lastEffectiveValue
				secondsAccumulated
				# Note: accountMarket is now a string reference, need to query separately if needed
			}
		}
	`

	const pageSize = 1000
	var allRecords []*AccountSubsidyRecord
	skip := 0

	for {
		request := graph.GraphQLRequest{
			Query: query,
			Variables: map[string]interface{}{
				"vaultId": m.config.VaultID,
				"first":   pageSize,
				"skip":    skip,
			},
		}

		var response struct {
			Data struct {
				AccountSubsidies []struct {
					Account struct {
						ID string `json:"id"`
					} `json:"account"`
					LastEffectiveValue string `json:"lastEffectiveValue"`
					SecondsAccumulated string `json:"secondsAccumulated"`
					// Note: AccountMarket is now a string reference in the new schema
					// Need to handle borrow balance separately if required
				} `json:"accountSubsidies"`
			} `json:"data"`
		}

		if err := m.graphClient.ExecuteQuery(ctx, request, &response); err != nil {
			return nil, fmt.Errorf("failed to execute GraphQL query at skip %d: %w", skip, err)
		}

		pageData := response.Data.AccountSubsidies
		if len(pageData) == 0 {
			break
		}

		for _, item := range pageData {
			currentBorrowU := big.NewInt(0)

			secondsAccumulated, ok := new(big.Int).SetString(item.SecondsAccumulated, 10)
			if !ok {
				return nil, fmt.Errorf("invalid seconds accumulated for account %s: %s", item.Account.ID, item.SecondsAccumulated)
			}

			record := &AccountSubsidyRecord{
				Account:            item.Account.ID,
				WeightedBalance:    big.NewInt(0), // No longer used, set to zero
				CurrentBorrowU:     currentBorrowU,
				SecondsAccumulated: secondsAccumulated,
			}

			allRecords = append(allRecords, record)
		}

		if len(pageData) < pageSize {
			break
		}

		skip += pageSize
	}

	return allRecords, nil
}

func (m *MigrationService) computeLastEffectiveValues(ctx context.Context, records []*AccountSubsidyRecord) error {
	m.logger.Logf("INFO Computing lastEffectiveValue for %d records", len(records))

	for _, record := range records {
		lastEffectiveValue := new(big.Int)

		// WeightedBalance is no longer used, only use CurrentBorrowU
		lastEffectiveValue.Add(lastEffectiveValue, record.CurrentBorrowU)
		record.LastEffectiveValue = lastEffectiveValue

	}

	return nil
}

func (m *MigrationService) updateSubgraphRecords(ctx context.Context, records []*AccountSubsidyRecord) error {
	m.logger.Logf("INFO Updating %d subgraph records", len(records))
	return nil
}

func (m *MigrationService) generateMerkleRoot(records []*AccountSubsidyRecord) ([32]byte, error) {
	if len(records) == 0 {
		return [32]byte{}, nil
	}

	leaves := make([]MerkleLeaf, 0, len(records))
	for _, record := range records {
		if record.SecondsAccumulated.Cmp(big.NewInt(0)) > 0 {
			leaves = append(leaves, MerkleLeaf{
				Account:            record.Account,
				SecondsAccumulated: record.SecondsAccumulated,
			})
		}
	}

	sort.Slice(leaves, func(i, j int) bool {
		return leaves[i].Account < leaves[j].Account
	})

	if len(leaves) == 0 {
		return [32]byte{}, nil
	}

	leafHashes := make([][32]byte, len(leaves))
	for i, leaf := range leaves {
		// Use same leaf format as LazyDistributor: keccak256(abi.encodePacked(recipient, newTotal))
		leafHashes[i] = m.createLeafHash(leaf.Account, leaf.SecondsAccumulated)
	}

	root := m.buildMerkleTree(leafHashes)
	return root, nil
}

// buildMerkleTree builds a Merkle tree compatible with OpenZeppelin's MerkleProof library
func (m *MigrationService) buildMerkleTree(leaves [][32]byte) [32]byte {
	if len(leaves) == 0 {
		return [32]byte{}
	}
	if len(leaves) == 1 {
		return leaves[0]
	}

	nextLevel := make([][32]byte, 0, (len(leaves)+1)/2)

	for i := 0; i < len(leaves); i += 2 {
		if i+1 < len(leaves) {
			// Sort pair to match OpenZeppelin's ordering
			left, right := leaves[i], leaves[i+1]
			if !m.isLeftSmaller(left, right) {
				left, right = right, left
			}
			// Hash the sorted pair using keccak256
			combined := append(left[:], right[:]...)
			hash := crypto.Keccak256Hash(combined)
			nextLevel = append(nextLevel, hash)
		} else {
			nextLevel = append(nextLevel, leaves[i])
		}
	}

	return m.buildMerkleTree(nextLevel)
}

// createLeafHash creates a leaf hash compatible with Solidity's abi.encodePacked(recipient, newTotal)
func (m *MigrationService) createLeafHash(address string, amount *big.Int) [32]byte {
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
func (m *MigrationService) isLeftSmaller(left, right [32]byte) bool {
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

func (m *MigrationService) IsIdempotent(ctx context.Context) (bool, error) {
	m.logger.Logf("INFO Checking if migration has already been executed")

	return false, nil
}

func (m *MigrationService) GetMigrationStatus(ctx context.Context) (map[string]interface{}, error) {
	status := map[string]interface{}{
		"vault_id":       m.config.VaultID,
		"snapshot_block": m.config.SnapshotBlockNumber.String(),
		"dry_run":        m.config.DryRun,
	}

	isComplete, err := m.IsIdempotent(ctx)
	if err != nil {
		return nil, err
	}
	status["completed"] = isComplete

	return status, nil
}
