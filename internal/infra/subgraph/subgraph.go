package subgraph

import "context"

//go:generate moq -out subgraph_mocks.go . SubgraphClient

// SubgraphClient defines the interface for subgraph operations
type SubgraphClient interface {
	// basic query operations
	ExecuteQuery(ctx context.Context, request GraphQLRequest, response interface{}) error
	HealthCheck(ctx context.Context) error

	// account queries
	QueryAccounts(ctx context.Context) ([]Account, error)
	QueryAccountSubsidiesForVault(ctx context.Context, vaultAddress string) ([]AccountSubsidy, error)
	QueryAccountSubsidiesAtBlock(
		ctx context.Context,
		vaultAddress string,
		blockNumber int64,
	) ([]AccountSubsidy, error)
	QueryAccountSubsidiesForEpoch(
		ctx context.Context,
		vaultAddress string,
		epochEndTimestamp string,
	) ([]AccountSubsidy, error)

	// epoch queries
	QueryCompletedEpochs(ctx context.Context) ([]Epoch, error)
	QueryEpochByNumber(ctx context.Context, epochNumber string) (*Epoch, error)
	QueryCurrentActiveEpoch(ctx context.Context) (*Epoch, error)
	QueryEpochWithBlockInfo(ctx context.Context, epochNumber string) (*Epoch, error)

	// merkle distribution queries
	QueryMerkleDistributionForEpoch(
		ctx context.Context,
		epochNumber string,
		vaultAddress string,
	) (*MerkleDistribution, error)

	// advanced query operations
	ExecutePaginatedQuery(
		ctx context.Context,
		queryTemplate string,
		variables map[string]interface{},
		entityField string,
		response interface{},
	) error
	ExecuteQueryAtBlock(
		ctx context.Context,
		query string,
		variables map[string]interface{},
		blockNumber int64,
		response interface{},
	) error
	ExecutePaginatedQueryAtBlock(
		ctx context.Context,
		queryTemplate string,
		variables map[string]interface{},
		entityField string,
		blockNumber int64,
		response interface{},
	) error
}

// Config represents the configuration for subgraph client
type Config struct {
	Endpoint string
	Timeout  string
}
