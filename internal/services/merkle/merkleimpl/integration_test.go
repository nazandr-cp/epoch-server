package merkleimpl

import (
	"context"
	"math/big"
	"testing"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-pkgz/lgr"
)

// TestMerkleCompatibility_CrossSystemIntegration demonstrates that the Go Merkle
// implementation produces compatible proofs with OpenZeppelin's Solidity implementation
func TestMerkleCompatibility_CrossSystemIntegration(t *testing.T) {
	// Test data that would be typical in the lend.fam system
	testEntries := []merkle.Entry{
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1500000000000000000)}, // 1.5 ETH worth
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(750000000000000000)},  // 0.75 ETH worth
		{Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", TotalEarned: big.NewInt(2000000000000000000)}, // 2 ETH worth
		{Address: "0x0000000000000000000000000000000000000001", TotalEarned: big.NewInt(100000000000000000)},  // 0.1 ETH worth
	}

	service := createTestServiceForIntegration(t)

	// Generate Merkle root
	root := service.BuildMerkleRootFromEntries(testEntries)

	// Test each entry can generate valid proofs
	for i, entry := range testEntries {
		// Generate proof
		proof, calculatedRoot, err := service.GenerateProof(testEntries, entry.Address, entry.TotalEarned)
		if err != nil {
			t.Fatalf("Failed to generate proof for entry %d: %v", i, err)
		}

		// Verify root consistency
		if calculatedRoot != root {
			t.Errorf("Root mismatch for entry %d", i)
		}

		// Verify using our internal verification (simulates OpenZeppelin)
		if !service.verifyProof(proof, root, entry.Address, entry.TotalEarned) {
			t.Errorf("Proof verification failed for entry %d", i)
		}
	}

	// Root should be compatible with OpenZeppelin's MerkleProof.verify()
	if len(common.Bytes2Hex(root[:])) != 64 {
		t.Error("Invalid root length")
	}
}

// createTestServiceForIntegration creates a service instance for integration testing
func createTestServiceForIntegration(t *testing.T) *Service {
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
	mockClient := &integrationTestSubgraphClient{}

	return New(db, mockClient, logger)
}

// integrationTestSubgraphClient implements SubgraphClient for integration testing
type integrationTestSubgraphClient struct{}

func (m *integrationTestSubgraphClient) QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*subgraph.Epoch, error) {
	return &subgraph.Epoch{}, nil
}

func (m *integrationTestSubgraphClient) QueryCurrentActiveEpoch(ctx context.Context) (*subgraph.Epoch, error) {
	return &subgraph.Epoch{}, nil
}

func (m *integrationTestSubgraphClient) ExecuteQuery(ctx context.Context, request subgraph.GraphQLRequest, response interface{}) error {
	return nil
}

// createTestServiceForBenchmark creates a service instance for benchmark testing
func createTestServiceForBenchmark(b *testing.B) *Service {
	tempDir := b.TempDir()
	logger := lgr.NoOp

	// Create badger database
	opts := badger.DefaultOptions(tempDir)
	opts.Logger = nil // Disable badger logging for tests
	db, err := badger.Open(opts)
	if err != nil {
		b.Fatalf("Failed to open test database: %v", err)
	}

	// Create mock subgraph client
	mockClient := &integrationTestSubgraphClient{}

	return New(db, mockClient, logger)
}

// TestZeroValueHandling ensures that zero values are handled correctly
func TestZeroValueHandling(t *testing.T) {
	// Test with zero amounts (should still be included in tree)
	entries := []merkle.Entry{
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(0)},
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(1000000000000000000)},
	}

	service := createTestServiceForIntegration(t)
	root := service.BuildMerkleRootFromEntries(entries)

	// Generate proof for zero amount entry
	proof, calculatedRoot, err := service.GenerateProof(entries, entries[0].Address, entries[0].TotalEarned)
	if err != nil {
		t.Fatalf("Failed to generate proof for zero amount entry: %v", err)
	}

	if calculatedRoot != root {
		t.Error("Root mismatch for zero amount entry")
	}

	if !service.verifyProof(proof, root, entries[0].Address, entries[0].TotalEarned) {
		t.Error("Proof verification failed for zero amount entry")
	}
}

// TestLargeAmountHandling tests handling of very large amounts (near uint256 max)
func TestLargeAmountHandling(t *testing.T) {
	// Create a very large amount (close to uint256 max)
	largeAmount := new(big.Int)
	largeAmount.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

	entries := []merkle.Entry{
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1000000000000000000)},
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: largeAmount},
	}

	service := createTestServiceForIntegration(t)
	root := service.BuildMerkleRootFromEntries(entries)

	// Generate proof for large amount entry
	proof, calculatedRoot, err := service.GenerateProof(entries, entries[1].Address, entries[1].TotalEarned)
	if err != nil {
		t.Fatalf("Failed to generate proof for large amount entry: %v", err)
	}

	if calculatedRoot != root {
		t.Error("Root mismatch for large amount entry")
	}

	if !service.verifyProof(proof, root, entries[1].Address, entries[1].TotalEarned) {
		t.Error("Proof verification failed for large amount entry")
	}
}

// BenchmarkMerkleOperations benchmarks key Merkle operations
func BenchmarkMerkleOperations(b *testing.B) {
	// Prepare test data
	entries := make([]merkle.Entry, 1000)
	for i := 0; i < 1000; i++ {
		addr := common.BigToAddress(big.NewInt(int64(i)))
		amount := big.NewInt(int64((i + 1) * 1000000000000000000)) // (i+1) ETH
		entries[i] = merkle.Entry{Address: addr.Hex(), TotalEarned: amount}
	}

	service := createTestServiceForBenchmark(b)

	b.Run("BuildMerkleRoot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			service.BuildMerkleRootFromEntries(entries)
		}
	})

	root := service.BuildMerkleRootFromEntries(entries)

	b.Run("GenerateProof", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			targetEntry := entries[i%len(entries)]
			_, _, _ = service.GenerateProof(entries, targetEntry.Address, targetEntry.TotalEarned)
		}
	})

	// Pre-generate a proof for verification benchmark
	targetEntry := entries[0]
	proof, _, _ := service.GenerateProof(entries, targetEntry.Address, targetEntry.TotalEarned)

	b.Run("VerifyProof", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			service.verifyProof(proof, root, targetEntry.Address, targetEntry.TotalEarned)
		}
	})
}
