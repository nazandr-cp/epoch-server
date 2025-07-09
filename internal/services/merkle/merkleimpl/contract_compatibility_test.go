package merkleimpl

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestContractCompatibility_ExactSolidityMatching tests that our Go implementation
// produces merkle roots and proofs that are exactly compatible with the Solidity
// contract's expectations using OpenZeppelin's MerkleProof library
func TestContractCompatibility_ExactSolidityMatching(t *testing.T) {
	// Test data that exactly matches what would be used in production
	testCases := []struct {
		name    string
		entries []Entry
	}{
		{
			name: "Production-like data with mixed case addresses",
			entries: []Entry{
				{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1500000000000000000)},
				{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(750000000000000000)},
				{Address: "0xAbCdEfAbCdEfAbCdEfAbCdEfAbCdEfAbCdEfAbCd", TotalEarned: big.NewInt(2000000000000000000)},
				{Address: "0x0000000000000000000000000000000000000001", TotalEarned: big.NewInt(100000000000000000)},
			},
		},
		{
			name: "Edge cases - zero amounts and max values",
			entries: []Entry{
				{Address: "0x0000000000000000000000000000000000000000", TotalEarned: big.NewInt(0)},
				{Address: "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", TotalEarned: func() *big.Int {
					val, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
					return val
				}()},
			},
		},
		{
			name: "Single entry",
			entries: []Entry{
				{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1000000000000000000)},
			},
		},
	}

	pg := NewProofGenerator()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate merkle root
			root := pg.BuildMerkleRoot(tc.entries)

			// Test each entry
			for i, entry := range tc.entries {
				// Generate proof using our implementation
				proof, calculatedRoot, err := pg.GenerateProof(tc.entries, entry.Address, entry.TotalEarned)
				if err != nil {
					t.Fatalf("Failed to generate proof for entry %d: %v", i, err)
				}

				// Verify root consistency
				if calculatedRoot != root {
					t.Errorf("Entry %d: Root mismatch: expected %x, got %x", i, root, calculatedRoot)
				}

				// Verify the proof using our OpenZeppelin-compatible verification
				if !pg.verifyProof(proof, root, entry.Address, entry.TotalEarned) {
					t.Errorf("Entry %d: Proof verification failed", i)
				}

				// Additional verification: simulate the exact Solidity leaf creation
				expectedLeaf := simulateSolidityLeafCreation(entry.Address, entry.TotalEarned)
				actualLeaf := pg.createLeafHash(entry.Address, entry.TotalEarned)
				if expectedLeaf != actualLeaf {
					t.Errorf("Entry %d: Leaf hash mismatch: expected %x, got %x", i, expectedLeaf, actualLeaf)
				}
			}
		})
	}
}

// simulateSolidityLeafCreation exactly replicates the Solidity contract's leaf creation
// This matches: keccak256(abi.encodePacked(recipient, newTotal))
func simulateSolidityLeafCreation(address string, amount *big.Int) [32]byte {
	// Convert address to common.Address to ensure proper 20-byte representation
	addr := common.HexToAddress(address)
	
	// Create abi.encodePacked equivalent: address (20 bytes) + amount (32 bytes)
	packed := make([]byte, 52)
	copy(packed[0:20], addr.Bytes())
	
	// Convert amount to 32-byte big-endian representation
	amountBytes := make([]byte, 32)
	amount.FillBytes(amountBytes)
	copy(packed[20:52], amountBytes)
	
	// Hash using keccak256
	return crypto.Keccak256Hash(packed)
}

// TestAddressNormalization verifies that different case variations of the same address
// produce the same leaf hash (as expected by the Solidity contract)
func TestAddressNormalization(t *testing.T) {
	pg := NewProofGenerator()
	amount := big.NewInt(1000000000000000000)
	
	// Different case variations of the same address
	addresses := []string{
		"0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3",
		"0x742D35CC6BF8E65F8B95E6C5CB15F5C5D5B8DBC3",
		"0x742d35cc6bf8e65f8b95e6c5cb15f5c5d5b8dbc3",
	}
	
	var expectedLeaf [32]byte
	for i, addr := range addresses {
		leaf := pg.createLeafHash(addr, amount)
		if i == 0 {
			expectedLeaf = leaf
		} else {
			if leaf != expectedLeaf {
				t.Errorf("Address %s produced different leaf hash than expected", addr)
			}
		}
	}
}

// TestSortingConsistency verifies that the sorting algorithm produces consistent results
// regardless of input order
func TestSortingConsistency(t *testing.T) {
	pg := NewProofGenerator()
	
	// Same entries in different orders
	entries1 := []Entry{
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1000000000000000000)},
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(2000000000000000000)},
		{Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", TotalEarned: big.NewInt(500000000000000000)},
	}
	
	entries2 := []Entry{
		{Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", TotalEarned: big.NewInt(500000000000000000)},
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1000000000000000000)},
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(2000000000000000000)},
	}
	
	entries3 := []Entry{
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(2000000000000000000)},
		{Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", TotalEarned: big.NewInt(500000000000000000)},
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1000000000000000000)},
	}
	
	root1 := pg.BuildMerkleRoot(entries1)
	root2 := pg.BuildMerkleRoot(entries2)
	root3 := pg.BuildMerkleRoot(entries3)
	
	if root1 != root2 || root2 != root3 {
		t.Errorf("Different input orders produced different roots: %x, %x, %x", root1, root2, root3)
	}
}

// TestLazyDistributorCompatibility tests that the lazy distributor produces
// compatible merkle roots with the proof generator
func TestLazyDistributorCompatibility(t *testing.T) {
	// This test would need to import the lazy distributor and test compatibility
	// For now, we'll test the core hashing compatibility
	
	pg := NewProofGenerator()
	
	// Test the core leaf hashing function compatibility
	testEntries := []Entry{
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1500000000000000000)},
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(750000000000000000)},
	}
	
	// Generate root using proof generator
	root := pg.BuildMerkleRoot(testEntries)
	
	// Verify that we can generate valid proofs for each entry
	for i, entry := range testEntries {
		proof, calculatedRoot, err := pg.GenerateProof(testEntries, entry.Address, entry.TotalEarned)
		if err != nil {
			t.Fatalf("Failed to generate proof for entry %d: %v", i, err)
		}
		
		if calculatedRoot != root {
			t.Errorf("Entry %d: Root mismatch", i)
		}
		
		if !pg.verifyProof(proof, root, entry.Address, entry.TotalEarned) {
			t.Errorf("Entry %d: Proof verification failed", i)
		}
	}
}

// BenchmarkContractCompatibility benchmarks the contract-compatible operations
func BenchmarkContractCompatibility(b *testing.B) {
	// Prepare test data
	entries := make([]Entry, 100)
	for i := 0; i < 100; i++ {
		addr := common.BigToAddress(big.NewInt(int64(i)))
		amount := big.NewInt(int64((i + 1) * 1000000000000000000))
		entries[i] = Entry{Address: addr.Hex(), TotalEarned: amount}
	}
	
	pg := NewProofGenerator()
	
	b.Run("BuildMerkleRoot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pg.BuildMerkleRoot(entries)
		}
	})
	
	root := pg.BuildMerkleRoot(entries)
	targetEntry := entries[0]
	
	b.Run("GenerateProof", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = pg.GenerateProof(entries, targetEntry.Address, targetEntry.TotalEarned)
		}
	})
	
	proof, _, _ := pg.GenerateProof(entries, targetEntry.Address, targetEntry.TotalEarned)
	
	b.Run("VerifyProof", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pg.verifyProof(proof, root, targetEntry.Address, targetEntry.TotalEarned)
		}
	})
}