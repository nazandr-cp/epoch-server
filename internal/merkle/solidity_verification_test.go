package merkle

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestExactSolidityCompatibility performs ultra-precise validation that our Go
// implementation matches the exact behavior expected by the Solidity contracts
func TestExactSolidityCompatibility(t *testing.T) {
	pg := NewProofGenerator()
	
	// Test case that exactly matches DebtSubsidizer.sol expected format
	// This mimics the exact data that would come from the lazy distributor
	entries := []Entry{
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1500000000000000000)}, // 1.5 ETH
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(750000000000000000)},  // 0.75 ETH  
		{Address: "0xAbCdEfAbCdEfAbCdEfAbCdEfAbCdEfAbCdEfAbCd", TotalEarned: big.NewInt(2000000000000000000)}, // 2 ETH
	}
	
	// Generate merkle root using our implementation
	root := pg.BuildMerkleRoot(entries)
	
	t.Logf("Generated Merkle Root: 0x%x", root)
	
	// For each entry, generate proof and verify it matches expected Solidity behavior
	for i, entry := range entries {
		t.Run(entry.Address, func(t *testing.T) {
			// Generate proof
			proof, calculatedRoot, err := pg.GenerateProof(entries, entry.Address, entry.TotalEarned)
			if err != nil {
				t.Fatalf("Failed to generate proof: %v", err)
			}
			
			// Verify root consistency
			if calculatedRoot != root {
				t.Errorf("Root mismatch: expected %x, got %x", root, calculatedRoot)
			}
			
			// Create the exact leaf that the Solidity contract would create
			expectedLeaf := createSolidityCompatibleLeaf(entry.Address, entry.TotalEarned)
			actualLeaf := pg.createLeafHash(entry.Address, entry.TotalEarned)
			
			if expectedLeaf != actualLeaf {
				t.Errorf("Leaf hash mismatch:\nExpected: %x\nActual:   %x", expectedLeaf, actualLeaf)
			}
			
			// Verify proof using OpenZeppelin-compatible verification
			if !simulateOpenZeppelinVerify(proof, root, expectedLeaf) {
				t.Errorf("OpenZeppelin verification failed")
			}
			
			// Log details for manual verification
			t.Logf("Entry %d:", i)
			t.Logf("  Address: %s", entry.Address)
			t.Logf("  Amount: %s", entry.TotalEarned.String())
			t.Logf("  Leaf: 0x%x", expectedLeaf)
			t.Logf("  Proof length: %d", len(proof))
			for j, p := range proof {
				t.Logf("  Proof[%d]: 0x%x", j, p)
			}
		})
	}
}

// createSolidityCompatibleLeaf creates a leaf hash that exactly matches
// the Solidity contract's: keccak256(abi.encodePacked(recipient, newTotal))
func createSolidityCompatibleLeaf(address string, amount *big.Int) [32]byte {
	// Convert address to bytes exactly as Solidity would
	addr := common.HexToAddress(address)
	
	// Create abi.encodePacked equivalent
	// In Solidity: abi.encodePacked(address recipient, uint256 newTotal)
	// - address is 20 bytes
	// - uint256 is 32 bytes in big-endian format
	packed := make([]byte, 52) // 20 + 32 = 52 bytes
	
	// Copy address bytes (20 bytes)
	copy(packed[0:20], addr.Bytes())
	
	// Copy amount as 32-byte big-endian
	amountBytes := make([]byte, 32)
	amount.FillBytes(amountBytes) // FillBytes fills with big-endian representation
	copy(packed[20:52], amountBytes)
	
	return crypto.Keccak256Hash(packed)
}

// simulateOpenZeppelinVerify simulates OpenZeppelin's MerkleProof.verify function
func simulateOpenZeppelinVerify(proof [][32]byte, root [32]byte, leaf [32]byte) bool {
	computedHash := leaf
	
	for _, proofElement := range proof {
		// OpenZeppelin sorts the pair before hashing: keccak256(abi.encodePacked(a, b))
		// where a and b are sorted so that a <= b
		if isHashSmaller(computedHash, proofElement) {
			// computedHash <= proofElement, so computedHash goes first
			combined := append(computedHash[:], proofElement[:]...)
			computedHash = crypto.Keccak256Hash(combined)
		} else {
			// proofElement < computedHash, so proofElement goes first
			combined := append(proofElement[:], computedHash[:]...)
			computedHash = crypto.Keccak256Hash(combined)
		}
	}
	
	return computedHash == root
}

// isHashSmaller compares two hashes using the same logic as OpenZeppelin
func isHashSmaller(a, b [32]byte) bool {
	for i := 0; i < 32; i++ {
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return false // Equal hashes
}

// TestKnownVectorCompatibility tests against known vectors that could be
// pre-computed and verified against actual Solidity contract calls
func TestKnownVectorCompatibility(t *testing.T) {
	pg := NewProofGenerator()
	
	// Known test vector - these values would be verified against actual contract
	testVector := struct {
		address      string
		amount       *big.Int
		expectedLeaf string // This would come from actual Solidity execution
	}{
		address:      "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3",
		amount:       big.NewInt(1000000000000000000), // 1 ETH
		expectedLeaf: "", // Would be filled with actual contract result
	}
	
	// Generate leaf using our implementation
	actualLeaf := pg.createLeafHash(testVector.address, testVector.amount)
	
	// For now, just verify it's deterministic and not zero
	if actualLeaf == [32]byte{} {
		t.Error("Leaf hash should not be zero")
	}
	
	// Verify the leaf matches our Solidity-compatible function
	expectedLeaf := createSolidityCompatibleLeaf(testVector.address, testVector.amount)
	if actualLeaf != expectedLeaf {
		t.Errorf("Leaf mismatch:\nActual:   %x\nExpected: %x", actualLeaf, expectedLeaf)
	}
	
	t.Logf("Address: %s", testVector.address)
	t.Logf("Amount: %s", testVector.amount.String())
	t.Logf("Leaf: 0x%x", actualLeaf)
}

// TestCaseNormalization verifies that addresses are handled correctly
// regardless of case, matching Ethereum's case-insensitive address handling
func TestCaseNormalization(t *testing.T) {
	pg := NewProofGenerator()
	amount := big.NewInt(1000000000000000000)
	
	// Same address in different cases
	addresses := []string{
		"0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", // Mixed case
		"0x742D35CC6BF8E65F8B95E6C5CB15F5C5D5B8DBC3", // Upper case
		"0x742d35cc6bf8e65f8b95e6c5cb15f5c5d5b8dbc3", // Lower case
	}
	
	var expectedLeaf [32]byte
	for i, addr := range addresses {
		leaf := pg.createLeafHash(addr, amount)
		if i == 0 {
			expectedLeaf = leaf
			t.Logf("Reference leaf for %s: 0x%x", addr, leaf)
		} else {
			if leaf != expectedLeaf {
				t.Errorf("Address %s produced different leaf hash: 0x%x (expected 0x%x)", 
					addr, leaf, expectedLeaf)
			} else {
				t.Logf("Address %s correctly produced same leaf: 0x%x", addr, leaf)
			}
		}
	}
}

// BenchmarkSolidityCompatibleOperations benchmarks the Solidity-compatible functions
func BenchmarkSolidityCompatibleOperations(b *testing.B) {
	pg := NewProofGenerator()
	address := "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3"
	amount := big.NewInt(1000000000000000000)
	
	b.Run("CreateLeafHash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pg.createLeafHash(address, amount)
		}
	})
	
	b.Run("SolidityCompatibleLeaf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			createSolidityCompatibleLeaf(address, amount)
		}
	})
	
	// Test with multiple entries
	entries := []Entry{
		{Address: "0x742d35Cc6bF8E65f8b95E6c5CB15F5C5D5b8DbC3", TotalEarned: big.NewInt(1000000000000000000)},
		{Address: "0x1234567890123456789012345678901234567890", TotalEarned: big.NewInt(2000000000000000000)},
		{Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", TotalEarned: big.NewInt(500000000000000000)},
	}
	
	b.Run("BuildMerkleRoot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pg.BuildMerkleRoot(entries)
		}
	})
}