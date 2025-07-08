package merkle

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-pkgz/lgr"
)

// ProofGenerator generates Merkle proofs compatible with OpenZeppelin's MerkleProof library
// This is the unified merkle tree implementation used across all services
type ProofGenerator struct {
	calculator       *Calculator
	timestampManager *TimestampManager
	logger           lgr.L
}

// Entry represents a leaf entry in the Merkle tree
type Entry struct {
	Address     string
	TotalEarned *big.Int
}

// NewProofGenerator creates a new proof generator
func NewProofGenerator() *ProofGenerator {
	return &ProofGenerator{
		calculator: NewCalculator(),
	}
}

// NewProofGeneratorWithDependencies creates a new proof generator with dependencies
func NewProofGeneratorWithDependencies(graphClient GraphClient, logger lgr.L) *ProofGenerator {
	return &ProofGenerator{
		calculator:       NewCalculator(),
		timestampManager: NewTimestampManager(graphClient, logger),
		logger:           logger,
	}
}

// GenerateProof generates a Merkle proof for a specific entry
func (pg *ProofGenerator) GenerateProof(entries []Entry, targetAddress string, targetAmount *big.Int) ([][32]byte, [32]byte, error) {
	if len(entries) == 0 {
		return nil, [32]byte{}, nil
	}

	// Sort entries deterministically by address
	sortedEntries := make([]Entry, len(entries))
	copy(sortedEntries, entries)
	pg.sortEntries(sortedEntries)

	// Find target index
	targetIndex := -1
	normalizedTargetAddress := strings.ToLower(targetAddress)
	for i, entry := range sortedEntries {
		if strings.ToLower(entry.Address) == normalizedTargetAddress && entry.TotalEarned.Cmp(targetAmount) == 0 {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return nil, [32]byte{}, nil
	}

	// Generate leaf hashes
	leafHashes := make([][32]byte, len(sortedEntries))
	for i, entry := range sortedEntries {
		leafHashes[i] = pg.createLeafHash(entry.Address, entry.TotalEarned)
	}

	// Generate proof and root
	proof := pg.generateMerkleProof(leafHashes, targetIndex)
	root := pg.buildMerkleRoot(leafHashes)

	return proof, root, nil
}

// BuildMerkleRoot builds the Merkle root from entries
func (pg *ProofGenerator) BuildMerkleRoot(entries []Entry) [32]byte {
	if len(entries) == 0 {
		return [32]byte{}
	}

	// Sort entries deterministically by address
	sortedEntries := make([]Entry, len(entries))
	copy(sortedEntries, entries)
	pg.sortEntries(sortedEntries)

	// Generate leaf hashes
	leafHashes := make([][32]byte, len(sortedEntries))
	for i, entry := range sortedEntries {
		leafHashes[i] = pg.createLeafHash(entry.Address, entry.TotalEarned)
	}

	return pg.buildMerkleRoot(leafHashes)
}

// sortEntries sorts entries by address to ensure deterministic ordering
// Uses case-sensitive comparison to match lazy_distributor sorting
func (pg *ProofGenerator) sortEntries(entries []Entry) {
	for i := 1; i < len(entries); i++ {
		key := entries[i]
		j := i - 1
		// Normalize addresses to lowercase for consistent comparison
		keyAddr := strings.ToLower(key.Address)
		for j >= 0 && strings.ToLower(entries[j].Address) > keyAddr {
			entries[j+1] = entries[j]
			j--
		}
		entries[j+1] = key
	}
}

// createLeafHash creates a leaf hash compatible with Solidity's abi.encodePacked(recipient, newTotal)
func (pg *ProofGenerator) createLeafHash(address string, amount *big.Int) [32]byte {
	// Convert address string to common.Address (normalize case first)
	// Use checksum format to ensure consistency with contract expectations
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

// buildMerkleRoot builds the Merkle root from leaf hashes
func (pg *ProofGenerator) buildMerkleRoot(leaves [][32]byte) [32]byte {
	if len(leaves) == 0 {
		return [32]byte{}
	}
	if len(leaves) == 1 {
		return leaves[0]
	}

	currentLevel := leaves
	for len(currentLevel) > 1 {
		var nextLevel [][32]byte
		for i := 0; i < len(currentLevel); i += 2 {
			if i+1 < len(currentLevel) {
				// Sort pair to match OpenZeppelin's ordering
				left, right := currentLevel[i], currentLevel[i+1]
				if !pg.isLeftSmaller(left, right) {
					left, right = right, left
				}
				// Hash the sorted pair using keccak256
				combined := append(left[:], right[:]...)
				nextLevel = append(nextLevel, crypto.Keccak256Hash(combined))
			} else {
				// Odd number of nodes, promote the last one
				nextLevel = append(nextLevel, currentLevel[i])
			}
		}
		currentLevel = nextLevel
	}

	return currentLevel[0]
}

// generateMerkleProof generates a Merkle proof for a leaf at the given index
func (pg *ProofGenerator) generateMerkleProof(leaves [][32]byte, leafIndex int) [][32]byte {
	if len(leaves) == 0 || leafIndex < 0 || leafIndex >= len(leaves) {
		return nil
	}

	var proof [][32]byte
	currentLevel := leaves
	currentIndex := leafIndex

	for len(currentLevel) > 1 {
		var nextLevel [][32]byte
		var nextIndex int

		for i := 0; i < len(currentLevel); i += 2 {
			if i+1 < len(currentLevel) {
				left, right := currentLevel[i], currentLevel[i+1]

				// Add sibling to proof if this pair contains our target
				if i == currentIndex || i+1 == currentIndex {
					if i == currentIndex {
						// Our node is on the left, add right sibling
						proof = append(proof, right)
					} else {
						// Our node is on the right, add left sibling
						proof = append(proof, left)
					}
					nextIndex = len(nextLevel) // Index in next level
				}

				// Sort pair to match OpenZeppelin's ordering
				if !pg.isLeftSmaller(left, right) {
					left, right = right, left
				}

				// Hash the sorted pair
				combined := append(left[:], right[:]...)
				nextLevel = append(nextLevel, crypto.Keccak256Hash(combined))
			} else {
				// Odd number of nodes, promote the last one
				if i == currentIndex {
					nextIndex = len(nextLevel)
				}
				nextLevel = append(nextLevel, currentLevel[i])
			}
		}

		currentLevel = nextLevel
		currentIndex = nextIndex
	}

	return proof
}

// isLeftSmaller determines if left hash should come before right hash in OpenZeppelin ordering
func (pg *ProofGenerator) isLeftSmaller(left, right [32]byte) bool {
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

// TreeResult contains the result of merkle tree generation
type TreeResult struct {
	Entries    []Entry
	MerkleRoot [32]byte
	Timestamp  int64
}

// GenerateTreeFromSubsidies generates a merkle tree from account subsidies using consistent timestamp
func (pg *ProofGenerator) GenerateTreeFromSubsidies(ctx context.Context, vaultAddress string, subsidies []graph.AccountSubsidy) (*TreeResult, error) {
	if pg.timestampManager == nil {
		return nil, fmt.Errorf("timestamp manager not initialized - use NewProofGeneratorWithDependencies")
	}

	// Get the consistent timestamp for this vault
	epochTimestamp, err := pg.timestampManager.GetLatestEpochTimestamp(ctx, vaultAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get epoch timestamp: %w", err)
	}

	if pg.logger != nil {
		pg.logger.Logf("INFO generating merkle tree for vault %s using epoch %s processingCompletedTimestamp %d", 
			vaultAddress, epochTimestamp.EpochNumber, epochTimestamp.ProcessingCompletedTimestamp)
	}

	// Process subsidies to entries with positive earnings
	entries, err := pg.calculator.ProcessAccountSubsidies(subsidies, epochTimestamp.ProcessingCompletedTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to process account subsidies: %w", err)
	}

	// Generate merkle root
	merkleRoot := pg.BuildMerkleRoot(entries)

	return &TreeResult{
		Entries:    entries,
		MerkleRoot: merkleRoot,
		Timestamp:  epochTimestamp.ProcessingCompletedTimestamp,
	}, nil
}

// GenerateHistoricalTreeFromSubsidies generates a merkle tree from historical account subsidies
func (pg *ProofGenerator) GenerateHistoricalTreeFromSubsidies(ctx context.Context, epochNumber string, subsidies []graph.AccountSubsidy) (*TreeResult, error) {
	if pg.timestampManager == nil {
		return nil, fmt.Errorf("timestamp manager not initialized - use NewProofGeneratorWithDependencies")
	}

	// Get the historical timestamp for this epoch
	epochTimestamp, err := pg.timestampManager.GetHistoricalEpochTimestamp(ctx, epochNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical epoch timestamp: %w", err)
	}

	if pg.logger != nil {
		pg.logger.Logf("INFO generating historical merkle tree for epoch %s using processingCompletedTimestamp %d", 
			epochNumber, epochTimestamp.ProcessingCompletedTimestamp)
	}

	// Process subsidies to entries with positive earnings
	entries, err := pg.calculator.ProcessAccountSubsidies(subsidies, epochTimestamp.ProcessingCompletedTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to process account subsidies: %w", err)
	}

	// Generate merkle root
	merkleRoot := pg.BuildMerkleRoot(entries)

	return &TreeResult{
		Entries:    entries,
		MerkleRoot: merkleRoot,
		Timestamp:  epochTimestamp.ProcessingCompletedTimestamp,
	}, nil
}

// BuildMerkleRootFromSubsidies builds a merkle root directly from subsidies (for LazyDistributor)
func (pg *ProofGenerator) BuildMerkleRootFromSubsidies(ctx context.Context, vaultAddress string, subsidies []graph.AccountSubsidy) ([32]byte, int64, error) {
	result, err := pg.GenerateTreeFromSubsidies(ctx, vaultAddress, subsidies)
	if err != nil {
		return [32]byte{}, 0, err
	}
	return result.MerkleRoot, result.Timestamp, nil
}
