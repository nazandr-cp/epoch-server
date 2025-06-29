package migration

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"sort"

	"github.com/andrey/epoch-server/internal/clients/graph"
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
			accountSubsidiesPerCollections(
				where: { vault: $vaultId }
				orderBy: account
				orderDirection: asc
				first: $first
				skip: $skip
			) {
				account {
					id
				}
				weightedBalance
				lastEffectiveValue
				secondsAccumulated
				accountMarket {
					borrowBalance
				}
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
				AccountSubsidiesPerCollections []struct {
					Account struct {
						ID string `json:"id"`
					} `json:"account"`
					WeightedBalance    string `json:"weightedBalance"`
					LastEffectiveValue string `json:"lastEffectiveValue"`
					SecondsAccumulated string `json:"secondsAccumulated"`
					AccountMarket      struct {
						BorrowBalance string `json:"borrowBalance"`
					} `json:"accountMarket"`
				} `json:"accountSubsidiesPerCollections"`
			} `json:"data"`
		}

		if err := m.graphClient.ExecuteQuery(ctx, request, &response); err != nil {
			return nil, fmt.Errorf("failed to execute GraphQL query at skip %d: %w", skip, err)
		}

		pageData := response.Data.AccountSubsidiesPerCollections
		if len(pageData) == 0 {
			break
		}

		for _, item := range pageData {
			weightedBalance, ok := new(big.Int).SetString(item.WeightedBalance, 10)
			if !ok {
				return nil, fmt.Errorf("invalid weighted balance for account %s: %s", item.Account.ID, item.WeightedBalance)
			}

			currentBorrowU, ok := new(big.Int).SetString(item.AccountMarket.BorrowBalance, 10)
			if !ok {
				return nil, fmt.Errorf("invalid borrow balance for account %s: %s", item.Account.ID, item.AccountMarket.BorrowBalance)
			}

			secondsAccumulated, ok := new(big.Int).SetString(item.SecondsAccumulated, 10)
			if !ok {
				return nil, fmt.Errorf("invalid seconds accumulated for account %s: %s", item.Account.ID, item.SecondsAccumulated)
			}

			record := &AccountSubsidyRecord{
				Account:            item.Account.ID,
				WeightedBalance:    weightedBalance,
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

		if record.WeightedBalance.Cmp(big.NewInt(0)) > 0 {
			lastEffectiveValue.Add(lastEffectiveValue, record.WeightedBalance)
		}

		lastEffectiveValue.Add(lastEffectiveValue, record.CurrentBorrowU)
		record.LastEffectiveValue = lastEffectiveValue

		m.logger.Logf("DEBUG Account %s: weightedBalance=%s, currentBorrowU=%s, lastEffectiveValue=%s",
			record.Account,
			record.WeightedBalance.String(),
			record.CurrentBorrowU.String(),
			record.LastEffectiveValue.String())
	}

	return nil
}

func (m *MigrationService) updateSubgraphRecords(ctx context.Context, records []*AccountSubsidyRecord) error {
	m.logger.Logf("INFO Updating subgraph records (placeholder - requires actual subgraph mutation endpoint)")

	for _, record := range records {
		m.logger.Logf("INFO Would update account %s with lastEffectiveValue=%s",
			record.Account,
			record.LastEffectiveValue.String())
	}

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
		leafData := fmt.Sprintf("%s:%s", leaf.Account, leaf.SecondsAccumulated.String())
		leafHashes[i] = sha256.Sum256([]byte(leafData))
	}

	root := m.buildMerkleTree(leafHashes)
	return root, nil
}

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
			combined := append(leaves[i][:], leaves[i+1][:]...)
			hash := sha256.Sum256(combined)
			nextLevel = append(nextLevel, hash)
		} else {
			nextLevel = append(nextLevel, leaves[i])
		}
	}

	return m.buildMerkleTree(nextLevel)
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
