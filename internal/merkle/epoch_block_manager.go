package merkle

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/go-pkgz/lgr"
)

// EpochBlockManager handles consistent block and timestamp resolution for merkle tree operations
// This ensures both epoch processing and proof generation use the same blockchain state
type EpochBlockManager struct {
	graphClient GraphClient
	logger      lgr.L
}

// NewEpochBlockManager creates a new epoch block manager
func NewEpochBlockManager(graphClient GraphClient, logger lgr.L) *EpochBlockManager {
	return &EpochBlockManager{
		graphClient: graphClient,
		logger:      logger,
	}
}

// NewTimestampManager creates a new epoch block manager (deprecated name for backwards compatibility)
func NewTimestampManager(graphClient GraphClient, logger lgr.L) *EpochBlockManager {
	return NewEpochBlockManager(graphClient, logger)
}

// EpochTimestamp represents epoch timing and block information
type EpochTimestamp struct {
	EpochNumber                   string
	ProcessingCompletedTimestamp  int64
	StartTimestamp                int64
	EndTimestamp                  int64
	CreatedAtBlock                int64  // Block number where epoch was created
	UpdatedAtBlock                int64  // Block number where epoch was last updated
}

// GetLatestEpochTimestamp retrieves the latest epoch timestamp that has a merkle distribution
// This ensures both root generation and proof generation use the same timestamp source
func (ebm *EpochBlockManager) GetLatestEpochTimestamp(ctx context.Context, vaultAddress string) (*EpochTimestamp, error) {
	// Get the latest processed epoch for this vault (one that has merkle distribution)
	latestEpoch, err := ebm.getLatestProcessedEpochForVault(ctx, vaultAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest processed epoch: %w", err)
	}

	// Parse processingCompletedTimestamp (preferred) or fallback to startTimestamp
	var processingTime int64
	if latestEpoch.ProcessingCompletedTimestamp != "" {
		processingTime, err = strconv.ParseInt(latestEpoch.ProcessingCompletedTimestamp, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid processing completed timestamp: %s", latestEpoch.ProcessingCompletedTimestamp)
		}
	} else {
		// Fallback to startTimestamp if processingCompletedTimestamp is not available
		processingTime, err = strconv.ParseInt(latestEpoch.StartTimestamp, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid start timestamp: %s", latestEpoch.StartTimestamp)
		}
		ebm.logger.Logf("WARN using startTimestamp as fallback for epoch %s", latestEpoch.EpochNumber)
	}

	// Parse other timestamps
	startTime, err := strconv.ParseInt(latestEpoch.StartTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid start timestamp: %s", latestEpoch.StartTimestamp)
	}

	endTime, err := strconv.ParseInt(latestEpoch.EndTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid end timestamp: %s", latestEpoch.EndTimestamp)
	}

	ebm.logger.Logf("INFO resolved epoch %s timestamp: processingCompleted=%d, start=%d, end=%d", 
		latestEpoch.EpochNumber, processingTime, startTime, endTime)

	// Parse block numbers
	createdAtBlock, err := strconv.ParseInt(latestEpoch.CreatedAtBlock, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid created at block: %s", latestEpoch.CreatedAtBlock)
	}

	updatedAtBlock, err := strconv.ParseInt(latestEpoch.UpdatedAtBlock, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid updated at block: %s", latestEpoch.UpdatedAtBlock)
	}

	return &EpochTimestamp{
		EpochNumber:                   latestEpoch.EpochNumber,
		ProcessingCompletedTimestamp:  processingTime,
		StartTimestamp:                startTime,
		EndTimestamp:                  endTime,
		CreatedAtBlock:                createdAtBlock,
		UpdatedAtBlock:                updatedAtBlock,
	}, nil
}

// GetHistoricalEpochTimestamp retrieves timestamp information for a specific epoch
func (ebm *EpochBlockManager) GetHistoricalEpochTimestamp(ctx context.Context, epochNumber string) (*EpochTimestamp, error) {
	epoch, err := ebm.getEpochByNumber(ctx, epochNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get epoch %s: %w", epochNumber, err)
	}

	// Parse processingCompletedTimestamp (required for historical epochs)
	processingTime, err := strconv.ParseInt(epoch.ProcessingCompletedTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid processing completed timestamp for epoch %s: %s", epochNumber, epoch.ProcessingCompletedTimestamp)
	}

	// Parse other timestamps
	startTime, err := strconv.ParseInt(epoch.StartTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid start timestamp: %s", epoch.StartTimestamp)
	}

	endTime, err := strconv.ParseInt(epoch.EndTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid end timestamp: %s", epoch.EndTimestamp)
	}

	// Parse block numbers
	createdAtBlock, err := strconv.ParseInt(epoch.CreatedAtBlock, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid created at block for epoch %s: %s", epochNumber, epoch.CreatedAtBlock)
	}

	updatedAtBlock, err := strconv.ParseInt(epoch.UpdatedAtBlock, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid updated at block for epoch %s: %s", epochNumber, epoch.UpdatedAtBlock)
	}

	ebm.logger.Logf("INFO resolved historical epoch %s timestamp: processingCompleted=%d, start=%d, end=%d, createdBlock=%d", 
		epochNumber, processingTime, startTime, endTime, createdAtBlock)

	return &EpochTimestamp{
		EpochNumber:                   epochNumber,
		ProcessingCompletedTimestamp:  processingTime,
		StartTimestamp:                startTime,
		EndTimestamp:                  endTime,
		CreatedAtBlock:                createdAtBlock,
		UpdatedAtBlock:                updatedAtBlock,
	}, nil
}

// GetCurrentActiveEpochBlock retrieves the current active epoch's creation block number
// This is the critical method for ensuring block consistency between epoch processing and proof generation
func (ebm *EpochBlockManager) GetCurrentActiveEpochBlock(ctx context.Context) (*EpochTimestamp, error) {
	// Use the new graph client method to get current active epoch
	activeEpoch, err := ebm.graphClient.QueryCurrentActiveEpoch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current active epoch: %w", err)
	}

	// Parse timestamps
	startTime, err := strconv.ParseInt(activeEpoch.StartTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid start timestamp: %s", activeEpoch.StartTimestamp)
	}

	endTime, err := strconv.ParseInt(activeEpoch.EndTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid end timestamp: %s", activeEpoch.EndTimestamp)
	}

	// Parse processing completed timestamp (may be null for active epochs)
	var processingTime int64
	if activeEpoch.ProcessingCompletedTimestamp != "" {
		processingTime, err = strconv.ParseInt(activeEpoch.ProcessingCompletedTimestamp, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid processing completed timestamp: %s", activeEpoch.ProcessingCompletedTimestamp)
		}
	}

	// Parse block numbers
	createdAtBlock, err := strconv.ParseInt(activeEpoch.CreatedAtBlock, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid created at block: %s", activeEpoch.CreatedAtBlock)
	}

	updatedAtBlock, err := strconv.ParseInt(activeEpoch.UpdatedAtBlock, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid updated at block: %s", activeEpoch.UpdatedAtBlock)
	}

	ebm.logger.Logf("INFO found active epoch %s: createdAtBlock=%d, start=%d, end=%d", 
		activeEpoch.EpochNumber, createdAtBlock, startTime, endTime)

	return &EpochTimestamp{
		EpochNumber:                   activeEpoch.EpochNumber,
		ProcessingCompletedTimestamp:  processingTime,
		StartTimestamp:                startTime,
		EndTimestamp:                  endTime,
		CreatedAtBlock:                createdAtBlock,
		UpdatedAtBlock:                updatedAtBlock,
	}, nil
}

// GetEpochBlockByNumber retrieves block information for a specific epoch
func (ebm *EpochBlockManager) GetEpochBlockByNumber(ctx context.Context, epochNumber string) (*EpochTimestamp, error) {
	// Use the new graph client method to get epoch with block info
	epoch, err := ebm.graphClient.QueryEpochWithBlockInfo(ctx, epochNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get epoch %s with block info: %w", epochNumber, err)
	}

	// Parse timestamps
	startTime, err := strconv.ParseInt(epoch.StartTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid start timestamp: %s", epoch.StartTimestamp)
	}

	endTime, err := strconv.ParseInt(epoch.EndTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid end timestamp: %s", epoch.EndTimestamp)
	}

	// Parse processing completed timestamp
	var processingTime int64
	if epoch.ProcessingCompletedTimestamp != "" {
		processingTime, err = strconv.ParseInt(epoch.ProcessingCompletedTimestamp, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid processing completed timestamp: %s", epoch.ProcessingCompletedTimestamp)
		}
	}

	// Parse block numbers
	createdAtBlock, err := strconv.ParseInt(epoch.CreatedAtBlock, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid created at block: %s", epoch.CreatedAtBlock)
	}

	updatedAtBlock, err := strconv.ParseInt(epoch.UpdatedAtBlock, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid updated at block: %s", epoch.UpdatedAtBlock)
	}

	ebm.logger.Logf("INFO found epoch %s: createdAtBlock=%d, start=%d, end=%d", 
		epochNumber, createdAtBlock, startTime, endTime)

	return &EpochTimestamp{
		EpochNumber:                   epochNumber,
		ProcessingCompletedTimestamp:  processingTime,
		StartTimestamp:                startTime,
		EndTimestamp:                  endTime,
		CreatedAtBlock:                createdAtBlock,
		UpdatedAtBlock:                updatedAtBlock,
	}, nil
}

// getLatestProcessedEpochForVault retrieves the most recent epoch that has a merkle distribution for the vault
func (ebm *EpochBlockManager) getLatestProcessedEpochForVault(ctx context.Context, vaultAddress string) (*graph.Epoch, error) {
	query := `
		query GetLatestProcessedEpoch($vaultAddress: String!) {
			merkleDistributions(
				where: { vault: $vaultAddress }
				orderBy: timestamp
				orderDirection: desc
				first: 1
			) {
				epoch {
					id
					epochNumber
					status
					startTimestamp
					endTimestamp
					processingCompletedTimestamp
				}
				merkleRoot
				timestamp
			}
		}
	`

	variables := map[string]interface{}{
		"vaultAddress": strings.ToLower(vaultAddress),
	}

	var response struct {
		MerkleDistributions []struct {
			Epoch      graph.Epoch `json:"epoch"`
			MerkleRoot string      `json:"merkleRoot"`
			Timestamp  string      `json:"timestamp"`
		} `json:"merkleDistributions"`
	}

	if err := ebm.graphClient.ExecuteQuery(ctx, graph.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}, &response); err != nil {
		return nil, fmt.Errorf("failed to query latest processed epoch: %w", err)
	}

	if len(response.MerkleDistributions) == 0 {
		return nil, fmt.Errorf("no processed epochs found for vault %s", vaultAddress)
	}

	ebm.logger.Logf("INFO found merkle distribution for epoch %s with root %s", 
		response.MerkleDistributions[0].Epoch.EpochNumber, 
		response.MerkleDistributions[0].MerkleRoot)

	return &response.MerkleDistributions[0].Epoch, nil
}

// getEpochByNumber retrieves epoch information by epoch number including block info
func (ebm *EpochBlockManager) getEpochByNumber(ctx context.Context, epochNumber string) (*graph.Epoch, error) {
	// Use the new QueryEpochWithBlockInfo method to get complete epoch information
	return ebm.graphClient.QueryEpochWithBlockInfo(ctx, epochNumber)
}