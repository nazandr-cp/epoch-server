package epoch

import (
	"context"
	"math/big"
	"time"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
)

// UserEarningsResponse represents the response for user total earned query
type UserEarningsResponse struct {
	UserAddress   string `json:"userAddress"`
	VaultAddress  string `json:"vaultAddress"`
	TotalEarned   string `json:"totalEarned"`
	CalculatedAt  int64  `json:"calculatedAt"`
	DataTimestamp int64  `json:"dataTimestamp"` // Timestamp used for calculations
}

// ContractClient interface for blockchain operations
type ContractClient interface {
	StartEpoch(ctx context.Context) error
	GetCurrentEpochId(ctx context.Context) (*big.Int, error)
	ForceEndEpochWithZeroYield(ctx context.Context, epochId *big.Int, vaultAddress string) error
}

// SubgraphClient interface for querying subgraph data
type SubgraphClient interface {
	QueryAccounts(ctx context.Context) ([]subgraph.Account, error)
	ExecuteQuery(ctx context.Context, request subgraph.GraphQLRequest, response interface{}) error
}

// Calculator interface for earnings calculations
type Calculator interface {
	CalculateTotalEarned(subsidy subgraph.AccountSubsidy, epochEndTime int64) (*big.Int, error)
}

// EpochInfo represents information about an epoch
type EpochInfo struct {
	Number      *big.Int  `json:"number"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	BlockNumber int64     `json:"blockNumber"`
	Status      string    `json:"status"` // "pending", "active", "completed"
	VaultID     string    `json:"vaultId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
