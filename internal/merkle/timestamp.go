package merkle

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/go-pkgz/lgr"
)

// TimestampManager handles consistent timestamp resolution for merkle tree operations
type TimestampManager struct {
	graphClient GraphClient
	logger      lgr.L
}


// NewTimestampManager creates a new timestamp manager
func NewTimestampManager(graphClient GraphClient, logger lgr.L) *TimestampManager {
	return &TimestampManager{
		graphClient: graphClient,
		logger:      logger,
	}
}

// EpochTimestamp represents epoch timing information
type EpochTimestamp struct {
	EpochNumber                   string
	ProcessingCompletedTimestamp  int64
	StartTimestamp                int64
	EndTimestamp                  int64
}

// GetLatestEpochTimestamp retrieves the latest epoch timestamp that has a merkle distribution
// This ensures both root generation and proof generation use the same timestamp source
func (tm *TimestampManager) GetLatestEpochTimestamp(ctx context.Context, vaultAddress string) (*EpochTimestamp, error) {
	// Get the latest processed epoch for this vault (one that has merkle distribution)
	latestEpoch, err := tm.getLatestProcessedEpochForVault(ctx, vaultAddress)
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
		tm.logger.Logf("WARN using startTimestamp as fallback for epoch %s", latestEpoch.EpochNumber)
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

	tm.logger.Logf("INFO resolved epoch %s timestamp: processingCompleted=%d, start=%d, end=%d", 
		latestEpoch.EpochNumber, processingTime, startTime, endTime)

	return &EpochTimestamp{
		EpochNumber:                   latestEpoch.EpochNumber,
		ProcessingCompletedTimestamp:  processingTime,
		StartTimestamp:                startTime,
		EndTimestamp:                  endTime,
	}, nil
}

// GetHistoricalEpochTimestamp retrieves timestamp information for a specific epoch
func (tm *TimestampManager) GetHistoricalEpochTimestamp(ctx context.Context, epochNumber string) (*EpochTimestamp, error) {
	epoch, err := tm.getEpochByNumber(ctx, epochNumber)
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

	tm.logger.Logf("INFO resolved historical epoch %s timestamp: processingCompleted=%d, start=%d, end=%d", 
		epochNumber, processingTime, startTime, endTime)

	return &EpochTimestamp{
		EpochNumber:                   epochNumber,
		ProcessingCompletedTimestamp:  processingTime,
		StartTimestamp:                startTime,
		EndTimestamp:                  endTime,
	}, nil
}

// getLatestProcessedEpochForVault retrieves the most recent epoch that has a merkle distribution for the vault
func (tm *TimestampManager) getLatestProcessedEpochForVault(ctx context.Context, vaultAddress string) (*graph.Epoch, error) {
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

	if err := tm.graphClient.ExecuteQuery(ctx, graph.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}, &response); err != nil {
		return nil, fmt.Errorf("failed to query latest processed epoch: %w", err)
	}

	if len(response.MerkleDistributions) == 0 {
		return nil, fmt.Errorf("no processed epochs found for vault %s", vaultAddress)
	}

	tm.logger.Logf("INFO found merkle distribution for epoch %s with root %s", 
		response.MerkleDistributions[0].Epoch.EpochNumber, 
		response.MerkleDistributions[0].MerkleRoot)

	return &response.MerkleDistributions[0].Epoch, nil
}

// getEpochByNumber retrieves epoch information by epoch number
func (tm *TimestampManager) getEpochByNumber(ctx context.Context, epochNumber string) (*graph.Epoch, error) {
	query := `
		query GetEpochByNumber($epochNumber: String!) {
			epoches(where: { epochNumber: $epochNumber }) {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
				processingCompletedTimestamp
			}
		}
	`

	variables := map[string]interface{}{
		"epochNumber": epochNumber,
	}

	var response struct {
		Epoches []graph.Epoch `json:"epoches"`
	}

	if err := tm.graphClient.ExecuteQuery(ctx, graph.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}, &response); err != nil {
		return nil, fmt.Errorf("failed to query epoch %s: %w", epochNumber, err)
	}

	if len(response.Epoches) == 0 {
		return nil, fmt.Errorf("epoch %s not found", epochNumber)
	}

	return &response.Epoches[0], nil
}