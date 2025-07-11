package merkle

import (
	"context"
	"math/big"
	"time"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
)

// UserMerkleProofResponse represents a merkle proof response for a user
type UserMerkleProofResponse struct {
	UserAddress  string   `json:"userAddress"`
	VaultAddress string   `json:"vaultAddress"`
	EpochNumber  string   `json:"epochNumber,omitempty"`
	TotalEarned  string   `json:"totalEarned"`
	MerkleProof  []string `json:"merkleProof"`
	MerkleRoot   string   `json:"merkleRoot"`
	LeafIndex    int      `json:"leafIndex"`
	GeneratedAt  int64    `json:"generatedAt"`
}

// MerkleDistribution represents merkle distribution data for an epoch
type MerkleDistribution struct {
	EpochNumber       string   `json:"epochNumber"`
	VaultAddress      string   `json:"vaultAddress"`
	MerkleRoot        string   `json:"merkleRoot"`
	TotalSubsidies    string   `json:"totalSubsidies"`
	AccountsProcessed int      `json:"accountsProcessed"`
	Proofs            []string `json:"proofs"`
	CreatedAt         int64    `json:"createdAt"`
}

// SubgraphClient interface for subgraph operations
type SubgraphClient interface {
	QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*subgraph.Epoch, error)
	QueryCurrentActiveEpoch(ctx context.Context) (*subgraph.Epoch, error)
	QueryAccountSubsidiesForVault(ctx context.Context, vaultAddress string) ([]subgraph.AccountSubsidy, error)
	ExecuteQuery(ctx context.Context, request subgraph.GraphQLRequest, response interface{}) error
}

// Entry represents a leaf entry in the Merkle tree
type Entry struct {
	Address     string
	TotalEarned *big.Int
}

// EpochTimestamp represents epoch timing and block information
type EpochTimestamp struct {
	EpochNumber                  string
	ProcessingCompletedTimestamp int64
	StartTimestamp               int64
	EndTimestamp                 int64
	CreatedAtBlock               int64 // Block number where epoch was created
	UpdatedAtBlock               int64 // Block number where epoch was last updated
}

// TreeResult contains the result of merkle tree generation
type TreeResult struct {
	Entries     []Entry
	MerkleRoot  [32]byte
	Timestamp   int64
	BlockNumber int64 // Block number used for data consistency
}

// MerkleEntry represents a leaf entry in the Merkle tree
type MerkleEntry struct {
	Address     string   `json:"address"`
	TotalEarned *big.Int `json:"totalEarned"`
}

// MerkleSnapshot represents a complete snapshot of merkle tree data for an epoch
type MerkleSnapshot struct {
	EpochNumber *big.Int      `json:"epochNumber"`
	Entries     []MerkleEntry `json:"entries"`
	MerkleRoot  string        `json:"merkleRoot"`
	Timestamp   int64         `json:"timestamp"`
	VaultID     string        `json:"vaultId"`
	BlockNumber int64         `json:"blockNumber"`
	CreatedAt   time.Time     `json:"createdAt"`
}
