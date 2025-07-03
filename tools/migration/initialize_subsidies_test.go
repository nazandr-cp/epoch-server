package migration

import (
	"context"
	"math/big"
	"testing"

	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/go-pkgz/lgr"
)

type mockGraphClient struct {
	records map[string]*AccountSubsidyRecord
}

func (m *mockGraphClient) ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error {
	resp := response.(*struct {
		Data struct {
			AccountSubsidies []struct {
				Account struct {
					ID string `json:"id"`
				} `json:"account"`
				LastEffectiveValue string `json:"lastEffectiveValue"`
				SecondsAccumulated string `json:"secondsAccumulated"`
				// Note: AccountMarket removed from new schema structure
			} `json:"accountSubsidies"`
		} `json:"data"`
	})

	for _, record := range m.records {
		item := struct {
			Account struct {
				ID string `json:"id"`
			} `json:"account"`
			LastEffectiveValue string `json:"lastEffectiveValue"`
			SecondsAccumulated string `json:"secondsAccumulated"`
			// Note: AccountMarket removed from new schema structure
		}{}

		item.Account.ID = record.Account
		item.SecondsAccumulated = record.SecondsAccumulated.String()
		// Note: BorrowBalance no longer directly available in AccountSubsidy

		resp.Data.AccountSubsidies = append(resp.Data.AccountSubsidies, item)
	}

	return nil
}

type mockSubsidizerClient struct {
	updatedRoots map[string][32]byte
}

func (m *mockSubsidizerClient) UpdateMerkleRoot(ctx context.Context, vaultId string, root [32]byte) error {
	if m.updatedRoots == nil {
		m.updatedRoots = make(map[string][32]byte)
	}
	m.updatedRoots[vaultId] = root
	return nil
}

func TestMigrationService_ComputeLastEffectiveValues(t *testing.T) {
	logger := lgr.NoOp
	config := MigrationConfig{
		VaultID: "test-vault",
		DryRun:  true,
	}

	service := NewMigrationService(nil, nil, logger, config)

	tests := []struct {
		name          string
		record        *AccountSubsidyRecord
		expectedValue *big.Int
	}{
		{
			name: "zero borrow balance",
			record: &AccountSubsidyRecord{
				Account:            "0x123",
				CurrentBorrowU:     big.NewInt(1000),
				SecondsAccumulated: big.NewInt(3600),
			},
			expectedValue: big.NewInt(1000),
		},
		{
			name: "positive borrow balance",
			record: &AccountSubsidyRecord{
				Account:            "0x456",
				CurrentBorrowU:     big.NewInt(1000),
				SecondsAccumulated: big.NewInt(7200),
			},
			expectedValue: big.NewInt(1000),
		},
		{
			name: "zero borrow balance case 2",
			record: &AccountSubsidyRecord{
				Account:            "0x789",
				CurrentBorrowU:     big.NewInt(0),
				SecondsAccumulated: big.NewInt(1800),
			},
			expectedValue: big.NewInt(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records := []*AccountSubsidyRecord{tt.record}

			err := service.computeLastEffectiveValues(context.Background(), records)
			if err != nil {
				t.Fatalf("computeLastEffectiveValues failed: %v", err)
			}

			if tt.record.LastEffectiveValue.Cmp(tt.expectedValue) != 0 {
				t.Errorf("expected lastEffectiveValue %s, got %s",
					tt.expectedValue.String(),
					tt.record.LastEffectiveValue.String())
			}
		})
	}
}

func TestMigrationService_GenerateMerkleRoot(t *testing.T) {
	logger := lgr.NoOp
	config := MigrationConfig{
		VaultID: "test-vault",
		DryRun:  true,
	}

	service := NewMigrationService(nil, nil, logger, config)

	tests := []struct {
		name     string
		records  []*AccountSubsidyRecord
		expected bool // true if root should be non-zero
	}{
		{
			name:     "empty records",
			records:  []*AccountSubsidyRecord{},
			expected: false,
		},
		{
			name: "records with zero seconds accumulated",
			records: []*AccountSubsidyRecord{
				{
					Account:            "0x123",
					SecondsAccumulated: big.NewInt(0),
				},
			},
			expected: false,
		},
		{
			name: "single record with positive seconds",
			records: []*AccountSubsidyRecord{
				{
					Account:            "0x123",
					SecondsAccumulated: big.NewInt(3600),
				},
			},
			expected: true,
		},
		{
			name: "multiple records with positive seconds",
			records: []*AccountSubsidyRecord{
				{
					Account:            "0x123",
					SecondsAccumulated: big.NewInt(3600),
				},
				{
					Account:            "0x456",
					SecondsAccumulated: big.NewInt(7200),
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := service.generateMerkleRoot(tt.records)
			if err != nil {
				t.Fatalf("generateMerkleRoot failed: %v", err)
			}

			isEmpty := true
			for _, b := range root {
				if b != 0 {
					isEmpty = false
					break
				}
			}

			if tt.expected && isEmpty {
				t.Error("expected non-empty merkle root, got empty")
			}
			if !tt.expected && !isEmpty {
				t.Error("expected empty merkle root, got non-empty")
			}
		})
	}
}

func TestMigrationService_InitializeSubsidies(t *testing.T) {
	logger := lgr.NoOp
	config := MigrationConfig{
		VaultID: "test-vault",
		DryRun:  false,
	}

	mockGraph := &mockGraphClient{
		records: map[string]*AccountSubsidyRecord{
			"0x123": {
				Account:            "0x123",
				CurrentBorrowU:     big.NewInt(500),
				SecondsAccumulated: big.NewInt(3600),
			},
			"0x456": {
				Account:            "0x456",
				CurrentBorrowU:     big.NewInt(2000),
				SecondsAccumulated: big.NewInt(7200),
			},
		},
	}

	mockSubsidizer := &mockSubsidizerClient{}

	service := &MigrationService{
		graphClient:      mockGraph,
		subsidizerClient: mockSubsidizer,
		logger:           logger,
		config:           config,
	}

	err := service.InitializeSubsidies(context.Background())
	if err != nil {
		t.Fatalf("InitializeSubsidies failed: %v", err)
	}

	if _, exists := mockSubsidizer.updatedRoots["test-vault"]; !exists {
		t.Error("expected merkle root to be updated for test-vault")
	}
}

func TestMigrationService_DeterministicMerkleRoot(t *testing.T) {
	logger := lgr.NoOp
	config := MigrationConfig{
		VaultID: "test-vault",
		DryRun:  true,
	}

	service := NewMigrationService(nil, nil, logger, config)

	records1 := []*AccountSubsidyRecord{
		{
			Account:            "0x123",
			SecondsAccumulated: big.NewInt(3600),
		},
		{
			Account:            "0x456",
			SecondsAccumulated: big.NewInt(7200),
		},
	}

	records2 := []*AccountSubsidyRecord{
		{
			Account:            "0x456",
			SecondsAccumulated: big.NewInt(7200),
		},
		{
			Account:            "0x123",
			SecondsAccumulated: big.NewInt(3600),
		},
	}

	root1, err := service.generateMerkleRoot(records1)
	if err != nil {
		t.Fatalf("generateMerkleRoot failed for records1: %v", err)
	}

	root2, err := service.generateMerkleRoot(records2)
	if err != nil {
		t.Fatalf("generateMerkleRoot failed for records2: %v", err)
	}

	if root1 != root2 {
		t.Error("merkle roots should be deterministic regardless of input order")
	}
}

func TestMigrationService_BuildMerkleTree(t *testing.T) {
	logger := lgr.NoOp
	config := MigrationConfig{}

	service := NewMigrationService(nil, nil, logger, config)

	tests := []struct {
		name        string
		leavesCount int
	}{
		{"empty", 0},
		{"single leaf", 1},
		{"two leaves", 2},
		{"three leaves", 3},
		{"four leaves", 4},
		{"five leaves", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaves := make([][32]byte, tt.leavesCount)
			for i := 0; i < tt.leavesCount; i++ {
				leaves[i] = [32]byte{byte(i + 1)}
			}

			root := service.buildMerkleTree(leaves)

			if tt.leavesCount == 0 {
				expectedEmpty := [32]byte{}
				if root != expectedEmpty {
					t.Error("expected empty root for empty leaves")
				}
			} else {
				emptyRoot := [32]byte{}
				if root == emptyRoot {
					t.Error("expected non-empty root for non-empty leaves")
				}
			}
		})
	}
}
