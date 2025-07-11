package blockchain

import (
	"context"
	"math/big"
)

//go:generate moq -out blockchain_mocks.go . BlockchainClient

// BlockchainClient defines the interface for all blockchain operations
type BlockchainClient interface {
	// epoch management
	StartEpoch(ctx context.Context) error
	GetCurrentEpochId(ctx context.Context) (*big.Int, error)
	EndEpochWithSubsidies(
		ctx context.Context,
		epochId *big.Int,
		vaultAddress string,
		merkleRoot [32]byte,
		subsidiesDistributed *big.Int,
	) error
	ForceEndEpochWithZeroYield(ctx context.Context, epochId *big.Int, vaultAddress string) error

	// lending operations
	UpdateExchangeRate(ctx context.Context, lendingManagerAddress string) error

	// vault operations
	AllocateYieldToEpoch(ctx context.Context, epochId *big.Int, vaultAddress string) error
	AllocateCumulativeYieldToEpoch(
		ctx context.Context,
		epochId *big.Int,
		vaultAddress string,
		amount *big.Int,
	) error

	// subsidy distribution
	UpdateMerkleRoot(
		ctx context.Context,
		vaultId string,
		root [32]byte,
		totalSubsidies *big.Int,
	) error
	UpdateMerkleRootAndWaitForConfirmation(
		ctx context.Context,
		vaultId string,
		root [32]byte,
		totalSubsidies *big.Int,
	) error
	DistributeSubsidies(ctx context.Context, epochID string) error
}

// Config represents the configuration needed for blockchain clients
type Config struct {
	RPCURL             string
	PrivateKey         string
	GasLimit           uint64
	GasPrice           string
	Comptroller        string
	EpochManager       string
	DebtSubsidizer     string
	LendingManager     string
	CollectionRegistry string
}
