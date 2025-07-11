package merkleimpl

import (
	"context"
	"math/big"
	"testing"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-pkgz/lgr"
)

func TestProofGenerator_OpenZeppelinCompatibility(t *testing.T) {
	// Create test service
	service := createTestService(t)

	// Test with sample data similar to what the contracts would use
	entries := []merkle.Entry{
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1000000000000000000)}, // 1 ETH equivalent
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(2000000000000000000)}, // 2 ETH equivalent
		{Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", TotalEarned: big.NewInt(500000000000000000)},  // 0.5 ETH equivalent
	}

	// Generate Merkle root
	root := service.BuildMerkleRootFromEntries(entries)

	// Test proof generation for each entry
	for i, entry := range entries {
		proof, calculatedRoot, err := service.GenerateProof(entries, entry.Address, entry.TotalEarned)
		if err != nil {
			t.Fatalf("Failed to generate proof for entry %d: %v", i, err)
		}

		// Verify the calculated root matches
		if calculatedRoot != root {
			t.Errorf("Root mismatch for entry %d: expected %s, got %s",
				i, common.Bytes2Hex(root[:]), common.Bytes2Hex(calculatedRoot[:]))
		}

		// Verify the proof manually (simulate OpenZeppelin's verify function)
		if !service.verifyProof(proof, root, entry.Address, entry.TotalEarned) {
			t.Errorf("Proof verification failed for entry %d", i)
		}
	}
}

func TestProofGenerator_LeafHashCompatibility(t *testing.T) {
	service := createTestService(t)

	// Test cases that match the Solidity abi.encodePacked format
	testCases := []struct {
		address string
		amount  *big.Int
		name    string
	}{
		{
			address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3",
			amount:  big.NewInt(1000000000000000000),
			name:    "1 ETH",
		},
		{
			address: "0x0000000000000000000000000000000000000000",
			amount:  big.NewInt(0),
			name:    "Zero address, zero amount",
		},
		{
			address: "0xffffffffffffffffffffffffffffffffffffffff",
			amount:  new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)), // 2^256 - 1
			name:    "Max address, max uint256",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			leaf := service.CreateLeafHash(tc.address, tc.amount)

			// Verify the leaf hash format
			if leaf == [32]byte{} {
				t.Error("Leaf hash should not be empty")
			}

			// Verify consistency - same inputs should produce same hash
			leaf2 := service.CreateLeafHash(tc.address, tc.amount)
			if leaf != leaf2 {
				t.Error("Leaf hash should be deterministic")
			}
		})
	}
}

func TestProofGenerator_EmptyAndSingleEntry(t *testing.T) {
	service := createTestService(t)

	// Test empty entries
	emptyRoot := service.BuildMerkleRootFromEntries([]merkle.Entry{})
	if emptyRoot != [32]byte{} {
		t.Error("Empty entries should produce zero root")
	}

	// Test single entry
	singleEntry := []merkle.Entry{
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1000000000000000000)},
	}

	singleRoot := service.BuildMerkleRootFromEntries(singleEntry)
	expectedLeaf := service.CreateLeafHash(singleEntry[0].Address, singleEntry[0].TotalEarned)

	if singleRoot != expectedLeaf {
		t.Errorf("Single entry root should equal the leaf hash: expected %s, got %s",
			common.Bytes2Hex(expectedLeaf[:]), common.Bytes2Hex(singleRoot[:]))
	}

	// Test proof for single entry
	proof, root, err := service.GenerateProof(singleEntry, singleEntry[0].Address, singleEntry[0].TotalEarned)
	if err != nil {
		t.Fatalf("Failed to generate proof for single entry: %v", err)
	}

	if len(proof) != 0 {
		t.Error("Single entry should have empty proof")
	}

	if root != singleRoot {
		t.Error("Single entry proof root should match BuildMerkleRoot result")
	}
}

func TestProofGenerator_DeterministicSorting(t *testing.T) {
	service := createTestService(t)

	// Test that different input orders produce the same root
	entries1 := []merkle.Entry{
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1000000000000000000)},
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(2000000000000000000)},
		{Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", TotalEarned: big.NewInt(500000000000000000)},
	}

	entries2 := []merkle.Entry{
		{Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", TotalEarned: big.NewInt(500000000000000000)},
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1000000000000000000)},
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(2000000000000000000)},
	}

	root1 := service.BuildMerkleRootFromEntries(entries1)
	root2 := service.BuildMerkleRootFromEntries(entries2)

	if root1 != root2 {
		t.Errorf("Different input orders should produce same root: %s vs %s",
			common.Bytes2Hex(root1[:]), common.Bytes2Hex(root2[:]))
	}
}

// verifyProof simulates OpenZeppelin's MerkleProof.verify function
func (s *Service) verifyProof(proof [][32]byte, root [32]byte, address string, amount *big.Int) bool {
	leaf := s.CreateLeafHash(address, amount)
	return s.processProof(proof, leaf) == root
}

// processProof simulates OpenZeppelin's MerkleProof.processProof function
func (s *Service) processProof(proof [][32]byte, leaf [32]byte) [32]byte {
	computedHash := leaf
	for _, proofElement := range proof {
		if s.IsLeftSmaller(computedHash, proofElement) {
			// computedHash goes on the left
			combined := append(computedHash[:], proofElement[:]...)
			computedHash = crypto.Keccak256Hash(combined)
		} else {
			// computedHash goes on the right
			combined := append(proofElement[:], computedHash[:]...)
			computedHash = crypto.Keccak256Hash(combined)
		}
	}
	return computedHash
}

func TestProofGenerator_LargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large dataset test in short mode")
	}

	service := createTestService(t)

	// Generate a larger dataset to test performance and correctness
	entries := make([]merkle.Entry, 100)
	for i := 0; i < 100; i++ {
		// Generate deterministic but varied addresses and amounts
		addr := common.BigToAddress(big.NewInt(int64(i * 12345)))
		amount := big.NewInt(int64((i + 1) * 1000000000000000000)) // i+1 ETH
		entries[i] = merkle.Entry{
			Address:     addr.Hex(),
			TotalEarned: amount,
		}
	}

	// Build root
	root := service.BuildMerkleRootFromEntries(entries)

	// Test a few random proofs
	testIndices := []int{0, 10, 50, 99}
	for _, idx := range testIndices {
		entry := entries[idx]
		proof, calculatedRoot, err := service.GenerateProof(entries, entry.Address, entry.TotalEarned)
		if err != nil {
			t.Fatalf("Failed to generate proof for entry %d: %v", idx, err)
		}

		if calculatedRoot != root {
			t.Errorf("Root mismatch for entry %d", idx)
		}

		if !service.verifyProof(proof, root, entry.Address, entry.TotalEarned) {
			t.Errorf("Proof verification failed for entry %d", idx)
		}
	}
}

// createTestService creates a service instance for testing
func createTestService(t *testing.T) *Service {
	tempDir := t.TempDir()
	logger := lgr.NoOp

	// Create badger database
	opts := badger.DefaultOptions(tempDir)
	opts.Logger = nil // Disable badger logging for tests
	db, err := badger.Open(opts)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create mock subgraph client
	mockClient := &testSubgraphClient{}

	return New(db, mockClient, logger)
}

// testSubgraphClient implements SubgraphClient for testing
type testSubgraphClient struct{}

func (m *testSubgraphClient) QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*subgraph.Epoch, error) {
	return &subgraph.Epoch{
		EpochNumber:                  epochNumber,
		StartTimestamp:               "1640000000",
		EndTimestamp:                 "1640086400",
		ProcessingCompletedTimestamp: "1640086400",
		CreatedAtBlock:               "12345678",
		UpdatedAtBlock:               "12345680",
	}, nil
}

func (m *testSubgraphClient) QueryCurrentActiveEpoch(ctx context.Context) (*subgraph.Epoch, error) {
	return &subgraph.Epoch{
		EpochNumber:                  "1",
		StartTimestamp:               "1640000000",
		EndTimestamp:                 "1640086400",
		ProcessingCompletedTimestamp: "",
		CreatedAtBlock:               "12345678",
		UpdatedAtBlock:               "12345680",
	}, nil
}

func (m *testSubgraphClient) QueryAccountSubsidiesForVault(ctx context.Context, vaultAddress string) ([]subgraph.AccountSubsidy, error) {
	return []subgraph.AccountSubsidy{}, nil
}

func (m *testSubgraphClient) ExecuteQuery(ctx context.Context, request subgraph.GraphQLRequest, response interface{}) error {
	return nil
}
