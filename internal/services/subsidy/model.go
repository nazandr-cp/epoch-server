package subsidy

import (
	"context"
	"math/big"
	"time"
)

// SubsidyDistributionRequest represents a request to distribute subsidies
type SubsidyDistributionRequest struct {
	VaultID   string `json:"vaultId"`
	EpochID   string `json:"epochId,omitempty"`
	ForceMode bool   `json:"forceMode,omitempty"`
}

// SubsidyDistributionResponse represents the response from subsidy distribution
type SubsidyDistributionResponse struct {
	VaultID           string `json:"vaultId"`
	EpochID           string `json:"epochId"`
	TotalSubsidies    string `json:"totalSubsidies"`
	AccountsProcessed int    `json:"accountsProcessed"`
	MerkleRoot        string `json:"merkleRoot"`
	TransactionHash   string `json:"transactionHash,omitempty"`
	Status            string `json:"status"`
}

// DistributionResult represents the result of a subsidy distribution
type DistributionResult struct {
	TotalSubsidies    *big.Int `json:"totalSubsidies"`
	AccountsProcessed int      `json:"accountsProcessed"`
	MerkleRoot        string   `json:"merkleRoot"`
}

// LazyDistributor interface for subsidy distribution
type LazyDistributor interface {
	Run(ctx context.Context, vaultId string) (*DistributionResult, error)
	RunWithEpoch(ctx context.Context, vaultId string, epochNumber *big.Int) (*DistributionResult, error)
}

// SubsidyDistribution represents a subsidy distribution record
type SubsidyDistribution struct {
	ID                string    `json:"id"`
	EpochNumber       *big.Int  `json:"epochNumber"`
	VaultID           string    `json:"vaultId"`
	CollectionAddress string    `json:"collectionAddress"`
	Amount            *big.Int  `json:"amount"`
	Status            string    `json:"status"` // "pending", "distributed", "failed"
	TxHash            string    `json:"txHash,omitempty"`
	BlockNumber       int64     `json:"blockNumber,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
