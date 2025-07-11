package merkleimpl

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
	"github.com/andrey/epoch-server/internal/infra/utils"
	"github.com/andrey/epoch-server/internal/services/merkle"
	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-pkgz/lgr"
)

type Service struct {
	store       *Store
	graphClient merkle.SubgraphClient
	logger      lgr.L
}

func New(db *badger.DB, graphClient merkle.SubgraphClient, logger lgr.L) *Service {
	return &Service{
		store:       NewStore(db, logger),
		graphClient: graphClient,
		logger:      logger,
	}
}

func (s *Service) GenerateUserMerkleProof(ctx context.Context, userAddress, vaultAddress string) (*merkle.UserMerkleProofResponse, error) {
	if userAddress == "" {
		return nil, fmt.Errorf("%w: userAddress cannot be empty", merkle.ErrInvalidInput)
	}
	if vaultAddress == "" {
		return nil, fmt.Errorf("%w: vaultAddress cannot be empty", merkle.ErrInvalidInput)
	}

	s.logger.Logf("INFO generating merkle proof for user %s in vault %s", userAddress, vaultAddress)

	// First try to get from stored snapshot (prioritize snapshot over subgraph)
	latestSnapshot, err := s.store.GetLatestSnapshot(ctx, vaultAddress)
	if err == nil && latestSnapshot != nil {
		s.logger.Logf("INFO found latest snapshot for vault %s, epoch %s with %d entries, root: %s",
			vaultAddress, latestSnapshot.EpochNumber.String(), len(latestSnapshot.Entries), latestSnapshot.MerkleRoot)
		return s.generateProofFromSnapshot(latestSnapshot, userAddress)
	}

	s.logger.Logf("WARN no snapshot found for vault %s, falling back to subgraph: %v", vaultAddress, err)

	// IMPORTANT: When using subgraph fallback, we must use the exact same epoch and data
	// that was used during the last distribution to ensure merkle root consistency

	// Fallback: Get the latest processed epoch for this vault from subgraph
	latestEpoch, err := s.getLatestProcessedEpochForVault(ctx, vaultAddress)
	if err != nil {
		s.logger.Logf("ERROR failed to get latest processed epoch: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	// Get epoch timestamp information
	epochTimestamp, err := s.parseEpochTimestamp(latestEpoch)
	if err != nil {
		s.logger.Logf("ERROR failed to parse epoch timestamp: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	// Get account subsidies for the vault from subgraph
	subsidies, err := s.getAccountSubsidiesForVault(ctx, vaultAddress)
	if err != nil {
		s.logger.Logf("ERROR failed to get account subsidies: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	// Process subsidies to entries with positive earnings
	entries, err := s.processAccountSubsidies(subsidies, epochTimestamp.ProcessingCompletedTimestamp)
	if err != nil {
		s.logger.Logf("ERROR failed to process account subsidies: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	// Find the user's entry
	normalizedUserAddress := utils.NormalizeAddress(userAddress)
	var userEntry *merkle.Entry
	for _, entry := range entries {
		if utils.NormalizeAddress(entry.Address) == normalizedUserAddress {
			userEntry = &entry
			break
		}
	}

	if userEntry == nil {
		s.logger.Logf("WARN user %s not found in vault %s entries", userAddress, vaultAddress)
		return nil, fmt.Errorf("%w: user not found in vault entries", merkle.ErrNotFound)
	}

	// Generate merkle proof
	proof, root, err := s.GenerateProof(entries, userEntry.Address, userEntry.TotalEarned)
	if err != nil {
		s.logger.Logf("ERROR failed to generate merkle proof: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	// Find leaf index
	leafIndex := s.findLeafIndex(entries, userEntry.Address, userEntry.TotalEarned)

	// Convert proof to string array
	proofStrings := make([]string, len(proof))
	for i, p := range proof {
		proofStrings[i] = common.Bytes2Hex(p[:])
	}

	return &merkle.UserMerkleProofResponse{
		UserAddress:  userAddress,
		VaultAddress: vaultAddress,
		EpochNumber:  latestEpoch.EpochNumber,
		TotalEarned:  userEntry.TotalEarned.String(),
		MerkleProof:  proofStrings,
		MerkleRoot:   common.Bytes2Hex(root[:]),
		LeafIndex:    leafIndex,
		GeneratedAt:  time.Now().Unix(),
	}, nil
}

func (s *Service) GenerateHistoricalMerkleProof(ctx context.Context, userAddress, vaultAddress, epochNumber string) (*merkle.UserMerkleProofResponse, error) {
	if userAddress == "" {
		return nil, fmt.Errorf("%w: userAddress cannot be empty", merkle.ErrInvalidInput)
	}
	if vaultAddress == "" {
		return nil, fmt.Errorf("%w: vaultAddress cannot be empty", merkle.ErrInvalidInput)
	}
	if epochNumber == "" {
		return nil, fmt.Errorf("%w: epochNumber cannot be empty", merkle.ErrInvalidInput)
	}

	s.logger.Logf("INFO generating historical merkle proof for user %s in vault %s for epoch %s", userAddress, vaultAddress, epochNumber)

	// First try to get from stored snapshot
	epochNum, ok := new(big.Int).SetString(epochNumber, 10)
	if !ok {
		return nil, fmt.Errorf("%w: invalid epoch number format", merkle.ErrInvalidInput)
	}

	snapshot, err := s.store.GetSnapshot(ctx, epochNum, vaultAddress)
	if err == nil {
		// Found stored snapshot, generate proof from it
		return s.generateProofFromSnapshot(snapshot, userAddress)
	}

	// If snapshot not found, generate from subgraph data
	s.logger.Logf("INFO snapshot not found for epoch %s, generating from subgraph data", epochNumber)

	// Get historical epoch information
	epoch, err := s.getEpochByNumber(ctx, epochNumber)
	if err != nil {
		s.logger.Logf("ERROR failed to get historical epoch: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	epochTimestamp, err := s.parseEpochTimestamp(epoch)
	if err != nil {
		s.logger.Logf("ERROR failed to parse historical epoch timestamp: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	// Get historical account subsidies for the vault
	subsidies, err := s.getHistoricalAccountSubsidiesForVault(ctx, vaultAddress, epochNumber)
	if err != nil {
		s.logger.Logf("ERROR failed to get historical account subsidies: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	// Process subsidies to entries with positive earnings
	entries, err := s.processAccountSubsidies(subsidies, epochTimestamp.ProcessingCompletedTimestamp)
	if err != nil {
		s.logger.Logf("ERROR failed to process historical account subsidies: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	// Find the user's entry
	normalizedUserAddress := utils.NormalizeAddress(userAddress)
	var userEntry *merkle.Entry
	for _, entry := range entries {
		if utils.NormalizeAddress(entry.Address) == normalizedUserAddress {
			userEntry = &entry
			break
		}
	}

	if userEntry == nil {
		s.logger.Logf("WARN user %s not found in vault %s entries for epoch %s", userAddress, vaultAddress, epochNumber)
		return nil, fmt.Errorf("%w: user not found in vault entries for epoch", merkle.ErrNotFound)
	}

	// Generate merkle proof
	proof, root, err := s.GenerateProof(entries, userEntry.Address, userEntry.TotalEarned)
	if err != nil {
		s.logger.Logf("ERROR failed to generate historical merkle proof: %v", err)
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	// Find leaf index
	leafIndex := s.findLeafIndex(entries, userEntry.Address, userEntry.TotalEarned)

	// Convert proof to string array
	proofStrings := make([]string, len(proof))
	for i, p := range proof {
		proofStrings[i] = common.Bytes2Hex(p[:])
	}

	return &merkle.UserMerkleProofResponse{
		UserAddress:  userAddress,
		VaultAddress: vaultAddress,
		EpochNumber:  epochNumber,
		TotalEarned:  userEntry.TotalEarned.String(),
		MerkleProof:  proofStrings,
		MerkleRoot:   common.Bytes2Hex(root[:]),
		LeafIndex:    leafIndex,
		GeneratedAt:  time.Now().Unix(),
	}, nil
}

func (s *Service) CalculateTotalEarned(subsidy subgraph.AccountSubsidy, endTimestamp int64) (*big.Int, error) {
	secondsAccumulated, ok := new(big.Int).SetString(subsidy.SecondsAccumulated, 10)
	if !ok {
		return nil, fmt.Errorf("invalid secondsAccumulated: %s", subsidy.SecondsAccumulated)
	}

	lastEffectiveValue, ok := new(big.Int).SetString(subsidy.LastEffectiveValue, 10)
	if !ok {
		return nil, fmt.Errorf("invalid lastEffectiveValue: %s", subsidy.LastEffectiveValue)
	}

	updatedAtTimestamp, err := strconv.ParseInt(subsidy.UpdatedAtTimestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid updatedAtTimestamp: %s", subsidy.UpdatedAtTimestamp)
	}

	// Calculate additional seconds from updatedAtTimestamp to end timestamp
	deltaT := endTimestamp - updatedAtTimestamp
	extraSeconds := new(big.Int).Mul(big.NewInt(deltaT), lastEffectiveValue)
	newTotalSeconds := new(big.Int).Add(secondsAccumulated, extraSeconds)

	// Convert seconds to tokens
	totalEarned := s.secondsToTokens(newTotalSeconds)
	return totalEarned, nil
}

func (s *Service) secondsToTokens(seconds *big.Int) *big.Int {
	conversionRate := big.NewInt(1000000000000000000) // 1e18
	return new(big.Int).Div(seconds, conversionRate)
}

func (s *Service) processAccountSubsidies(subsidies []subgraph.AccountSubsidy, endTimestamp int64) ([]merkle.Entry, error) {
	var entries []merkle.Entry

	for _, subsidy := range subsidies {
		totalEarned, err := s.CalculateTotalEarned(subsidy, endTimestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate total earned for account %s: %w", subsidy.Account.ID, err)
		}

		// Only include accounts with positive earnings
		if totalEarned.Cmp(big.NewInt(0)) > 0 {
			entries = append(entries, merkle.Entry{
				Address:     subsidy.Account.ID,
				TotalEarned: totalEarned,
			})
		}
	}

	return entries, nil
}

func (s *Service) GenerateProof(entries []merkle.Entry, targetAddress string, targetAmount *big.Int) ([][32]byte, [32]byte, error) {
	if len(entries) == 0 {
		return nil, [32]byte{}, nil
	}

	// Sort entries deterministically by address
	sortedEntries := make([]merkle.Entry, len(entries))
	copy(sortedEntries, entries)
	s.sortEntries(sortedEntries)

	// Find target index
	targetIndex := -1
	normalizedTargetAddress := utils.NormalizeAddress(targetAddress)
	for i, entry := range sortedEntries {
		if utils.NormalizeAddress(entry.Address) == normalizedTargetAddress && entry.TotalEarned.Cmp(targetAmount) == 0 {
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
		leafHashes[i] = s.CreateLeafHash(entry.Address, entry.TotalEarned)
	}

	// Generate proof and root
	proof := s.generateMerkleProof(leafHashes, targetIndex)
	root := s.buildMerkleRoot(leafHashes)

	return proof, root, nil
}

func (s *Service) BuildMerkleRootFromEntries(entries []merkle.Entry) [32]byte {
	if len(entries) == 0 {
		return [32]byte{}
	}

	// Sort entries deterministically by address
	sortedEntries := make([]merkle.Entry, len(entries))
	copy(sortedEntries, entries)
	s.sortEntries(sortedEntries)

	// Generate leaf hashes
	leafHashes := make([][32]byte, len(sortedEntries))
	for i, entry := range sortedEntries {
		leafHashes[i] = s.CreateLeafHash(entry.Address, entry.TotalEarned)
	}

	return s.buildMerkleRoot(leafHashes)
}

func (s *Service) sortEntries(entries []merkle.Entry) {
	for i := 1; i < len(entries); i++ {
		key := entries[i]
		j := i - 1
		// Normalize addresses to lowercase for consistent comparison
		keyAddr := utils.NormalizeAddress(key.Address)
		for j >= 0 && utils.NormalizeAddress(entries[j].Address) > keyAddr {
			entries[j+1] = entries[j]
			j--
		}
		entries[j+1] = key
	}
}

func (s *Service) CreateLeafHash(address string, amount *big.Int) [32]byte {
	// Convert address string to common.Address (normalize case first)
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

func (s *Service) buildMerkleRoot(leaves [][32]byte) [32]byte {
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
				if !s.IsLeftSmaller(left, right) {
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

func (s *Service) generateMerkleProof(leaves [][32]byte, leafIndex int) [][32]byte {
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
				if !s.IsLeftSmaller(left, right) {
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

func (s *Service) IsLeftSmaller(left, right [32]byte) bool {
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

func (s *Service) findLeafIndex(entries []merkle.Entry, targetAddress string, targetAmount *big.Int) int {
	// Sort entries deterministically by address
	sortedEntries := make([]merkle.Entry, len(entries))
	copy(sortedEntries, entries)
	s.sortEntries(sortedEntries)

	normalizedTargetAddress := utils.NormalizeAddress(targetAddress)
	for i, entry := range sortedEntries {
		if utils.NormalizeAddress(entry.Address) == normalizedTargetAddress && entry.TotalEarned.Cmp(targetAmount) == 0 {
			return i
		}
	}
	return -1
}

func (s *Service) parseEpochTimestamp(epoch *subgraph.Epoch) (*merkle.EpochTimestamp, error) {
	// Parse processingCompletedTimestamp (preferred) or fallback to startTimestamp
	var processingTime int64
	var err error

	if epoch.ProcessingCompletedTimestamp != "" {
		processingTime, err = strconv.ParseInt(epoch.ProcessingCompletedTimestamp, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid processing completed timestamp: %s", epoch.ProcessingCompletedTimestamp)
		}
	} else {
		// Fallback to startTimestamp if processingCompletedTimestamp is not available
		processingTime, err = strconv.ParseInt(epoch.StartTimestamp, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid start timestamp: %s", epoch.StartTimestamp)
		}
		s.logger.Logf("WARN using startTimestamp as fallback for epoch %s", epoch.EpochNumber)
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
		return nil, fmt.Errorf("invalid created at block: %s", epoch.CreatedAtBlock)
	}

	updatedAtBlock, err := strconv.ParseInt(epoch.UpdatedAtBlock, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid updated at block: %s", epoch.UpdatedAtBlock)
	}

	s.logger.Logf("INFO resolved epoch %s timestamp: processingCompleted=%d, start=%d, end=%d",
		epoch.EpochNumber, processingTime, startTime, endTime)

	return &merkle.EpochTimestamp{
		EpochNumber:                  epoch.EpochNumber,
		ProcessingCompletedTimestamp: processingTime,
		StartTimestamp:               startTime,
		EndTimestamp:                 endTime,
		CreatedAtBlock:               createdAtBlock,
		UpdatedAtBlock:               updatedAtBlock,
	}, nil
}

func (s *Service) getLatestProcessedEpochForVault(ctx context.Context, vaultAddress string) (*subgraph.Epoch, error) {
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
					createdAtBlock
					updatedAtBlock
				}
				merkleRoot
				timestamp
			}
		}
	`

	variables := map[string]interface{}{
		"vaultAddress": utils.NormalizeAddress(vaultAddress),
	}

	var response struct {
		MerkleDistributions []struct {
			Epoch      subgraph.Epoch `json:"epoch"`
			MerkleRoot string         `json:"merkleRoot"`
			Timestamp  string         `json:"timestamp"`
		} `json:"merkleDistributions"`
	}

	if err := s.graphClient.ExecuteQuery(ctx, subgraph.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}, &response); err != nil {
		return nil, fmt.Errorf("failed to query latest processed epoch: %w", err)
	}

	if len(response.MerkleDistributions) == 0 {
		return nil, fmt.Errorf("no processed epochs found for vault %s", vaultAddress)
	}

	s.logger.Logf("INFO found merkle distribution for epoch %s with root %s",
		response.MerkleDistributions[0].Epoch.EpochNumber,
		response.MerkleDistributions[0].MerkleRoot)

	return &response.MerkleDistributions[0].Epoch, nil
}

func (s *Service) getEpochByNumber(ctx context.Context, epochNumber string) (*subgraph.Epoch, error) {
	return s.graphClient.QueryEpochWithBlockInfo(ctx, epochNumber)
}

func (s *Service) generateProofFromSnapshot(snapshot *merkle.MerkleSnapshot, userAddress string) (*merkle.UserMerkleProofResponse, error) {
	// Convert MerkleEntry to Entry
	entries := make([]merkle.Entry, len(snapshot.Entries))
	for i, entry := range snapshot.Entries {
		entries[i] = merkle.Entry(entry)
	}

	// Find the user's entry
	normalizedUserAddress := utils.NormalizeAddress(userAddress)
	var userEntry *merkle.Entry
	for _, entry := range entries {
		if utils.NormalizeAddress(entry.Address) == normalizedUserAddress {
			userEntry = &entry
			break
		}
	}

	if userEntry == nil {
		return nil, fmt.Errorf("%w: user not found in snapshot", merkle.ErrNotFound)
	}

	// Generate merkle proof
	proof, root, err := s.GenerateProof(entries, userEntry.Address, userEntry.TotalEarned)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", merkle.ErrProofGeneration, err)
	}

	// Find leaf index
	leafIndex := s.findLeafIndex(entries, userEntry.Address, userEntry.TotalEarned)

	// Convert proof to string array
	proofStrings := make([]string, len(proof))
	for i, p := range proof {
		proofStrings[i] = common.Bytes2Hex(p[:])
	}

	return &merkle.UserMerkleProofResponse{
		UserAddress:  userAddress,
		VaultAddress: snapshot.VaultID,
		EpochNumber:  snapshot.EpochNumber.String(),
		TotalEarned:  userEntry.TotalEarned.String(),
		MerkleProof:  proofStrings,
		MerkleRoot:   common.Bytes2Hex(root[:]),
		LeafIndex:    leafIndex,
		GeneratedAt:  time.Now().Unix(),
	}, nil
}

func (s *Service) SaveSnapshot(ctx context.Context, epochNumber *big.Int, snapshot merkle.MerkleSnapshot) error {
	return s.store.SaveSnapshot(ctx, epochNumber, snapshot)
}

func (s *Service) getAccountSubsidiesForVault(ctx context.Context, vaultAddress string) ([]subgraph.AccountSubsidy, error) {
	return s.graphClient.QueryAccountSubsidiesForVault(ctx, vaultAddress)
}

func (s *Service) getHistoricalAccountSubsidiesForVault(ctx context.Context, vaultAddress, epochNumber string) ([]subgraph.AccountSubsidy, error) {
	return []subgraph.AccountSubsidy{}, nil
}
