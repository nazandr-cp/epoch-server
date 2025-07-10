package testing

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// TestDataGenerator generates realistic test data for BadgerDB integration tests
type TestDataGenerator struct {
	// removed rand field as we use crypto/rand directly
}

// secureRandomInt generates a secure random integer in range [0, max)
func (g *TestDataGenerator) secureRandomInt(max int) int {
	if max <= 0 {
		return 0
	}

	// use crypto/rand for secure random generation
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		// fallback to simple hash-based approach if crypto/rand fails
		return int(time.Now().UnixNano()) % max
	}

	// convert bytes to uint32 and mod by max
	randomValue := uint32(randomBytes[0])<<24 | uint32(randomBytes[1])<<16 | uint32(randomBytes[2])<<8 | uint32(randomBytes[3])
	return int(randomValue) % max
}

// secureRandomBigInt generates a secure random big.Int in range [0, max)
func (g *TestDataGenerator) secureRandomBigInt(max *big.Int) *big.Int {
	if max.Sign() <= 0 {
		return big.NewInt(0)
	}

	// use crypto/rand for secure random generation
	random, err := rand.Int(rand.Reader, max)
	if err != nil {
		// fallback to zero if crypto/rand fails
		return big.NewInt(0)
	}

	return random
}

// NewTestDataGenerator creates a new test data generator
func NewTestDataGenerator(seed int64) *TestDataGenerator {
	return &TestDataGenerator{}
}

// EpochData represents epoch test data
type EpochData struct {
	Number      *big.Int  `json:"number"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	BlockNumber int64     `json:"blockNumber"`
	Status      string    `json:"status"`
	VaultID     string    `json:"vaultId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// MerkleData represents merkle snapshot test data
type MerkleData struct {
	Entries     []MerkleEntry `json:"entries"`
	MerkleRoot  string        `json:"merkleRoot"`
	Timestamp   int64         `json:"timestamp"`
	VaultID     string        `json:"vaultId"`
	BlockNumber int64         `json:"blockNumber"`
	EpochNumber *big.Int      `json:"epochNumber"`
	CreatedAt   time.Time     `json:"createdAt"`
}

// MerkleEntry represents a merkle tree entry
type MerkleEntry struct {
	Address     string   `json:"address"`
	TotalEarned *big.Int `json:"totalEarned"`
}

// SubsidyData represents subsidy distribution test data
type SubsidyData struct {
	ID                string    `json:"id"`
	EpochNumber       *big.Int  `json:"epochNumber"`
	VaultID           string    `json:"vaultId"`
	CollectionAddress string    `json:"collectionAddress"`
	Amount            *big.Int  `json:"amount"`
	Status            string    `json:"status"`
	TxHash            string    `json:"txHash,omitempty"`
	BlockNumber       int64     `json:"blockNumber,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// GenerateEpochData generates realistic epoch test data
func (g *TestDataGenerator) GenerateEpochData(vaultID string, epochNumber *big.Int) EpochData {
	now := time.Now()
	startTime := now.Add(-time.Duration(g.secureRandomInt(3600)) * time.Second)
	endTime := startTime.Add(time.Duration(g.secureRandomInt(3600)+1800) * time.Second)

	statuses := []string{"pending", "active", "completed"}
	status := statuses[g.secureRandomInt(len(statuses))]

	return EpochData{
		Number:      epochNumber,
		StartTime:   startTime,
		EndTime:     endTime,
		BlockNumber: int64(g.secureRandomInt(1000000)) + 18000000, // Realistic block number
		Status:      status,
		VaultID:     vaultID,
		CreatedAt:   now.Add(-time.Duration(g.secureRandomInt(86400)) * time.Second),
		UpdatedAt:   now,
	}
}

// GenerateMultipleEpochData generates multiple epoch data entries
func (g *TestDataGenerator) GenerateMultipleEpochData(vaultID string, count int) []EpochData {
	epochs := make([]EpochData, count)
	for i := 0; i < count; i++ {
		epochNumber := big.NewInt(int64(i + 1))
		epochs[i] = g.GenerateEpochData(vaultID, epochNumber)
	}
	return epochs
}

// GenerateMerkleData generates realistic merkle snapshot test data
func (g *TestDataGenerator) GenerateMerkleData(vaultID string, epochNumber *big.Int, entryCount int) MerkleData {
	entries := make([]MerkleEntry, entryCount)
	for i := 0; i < entryCount; i++ {
		address := g.GenerateRandomAddress()
		amount := g.GenerateRandomAmount()
		entries[i] = MerkleEntry{
			Address:     address,
			TotalEarned: amount,
		}
	}

	return MerkleData{
		Entries:     entries,
		MerkleRoot:  g.GenerateRandomHash(),
		Timestamp:   time.Now().Unix(),
		VaultID:     vaultID,
		BlockNumber: int64(g.secureRandomInt(1000000)) + 18000000,
		EpochNumber: epochNumber,
		CreatedAt:   time.Now(),
	}
}

// GenerateSubsidyData generates realistic subsidy distribution test data
func (g *TestDataGenerator) GenerateSubsidyData(vaultID string, epochNumber *big.Int) SubsidyData {
	statuses := []string{"pending", "distributed", "failed"}
	status := statuses[g.secureRandomInt(len(statuses))]

	now := time.Now()

	data := SubsidyData{
		ID:                g.GenerateRandomID(),
		EpochNumber:       epochNumber,
		VaultID:           vaultID,
		CollectionAddress: g.GenerateRandomAddress(),
		Amount:            g.GenerateRandomAmount(),
		Status:            status,
		CreatedAt:         now.Add(-time.Duration(g.secureRandomInt(86400)) * time.Second),
		UpdatedAt:         now,
	}

	// Add tx hash and block number for distributed status
	if status == "distributed" {
		data.TxHash = g.GenerateRandomHash()
		data.BlockNumber = int64(g.secureRandomInt(1000000)) + 18000000
	}

	return data
}

// GenerateRandomAddress generates a random Ethereum address
func (g *TestDataGenerator) GenerateRandomAddress() string {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		// fallback to zero address if crypto/rand fails
		return common.Address{}.Hex()
	}
	return common.BytesToAddress(bytes).Hex()
}

// GenerateRandomHash generates a random hash
func (g *TestDataGenerator) GenerateRandomHash() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// fallback to zero hash if crypto/rand fails
		return hexutil.Encode(make([]byte, 32))
	}
	return hexutil.Encode(bytes)
}

// GenerateRandomAmount generates a random amount (up to 1000 ETH)
func (g *TestDataGenerator) GenerateRandomAmount() *big.Int {
	// Generate amount between 0.01 ETH and 1000 ETH
	min := big.NewInt(10000000000000000) // 0.01 ETH
	max := new(big.Int)
	max.SetString("1000000000000000000000", 10) // 1000 ETH

	diff := new(big.Int).Sub(max, min)
	random := g.secureRandomBigInt(diff)
	return random.Add(random, min)
}

// GenerateRandomID generates a random ID
func (g *TestDataGenerator) GenerateRandomID() string {
	return fmt.Sprintf("test-%d-%d", time.Now().UnixNano(), g.secureRandomInt(1000000))
}

// GenerateVaultID generates a random vault ID
func (g *TestDataGenerator) GenerateVaultID() string {
	return g.GenerateRandomAddress()
}

// TestDataSet represents a complete set of test data
type TestDataSet struct {
	VaultID   string
	Epochs    []EpochData
	Merkles   []MerkleData
	Subsidies []SubsidyData
}

// GenerateTestDataSet generates a complete test data set
func (g *TestDataGenerator) GenerateTestDataSet(vaultID string, epochCount int, merkleEntriesPerEpoch int, subsidiesPerEpoch int) TestDataSet {
	epochs := g.GenerateMultipleEpochData(vaultID, epochCount)

	var merkles []MerkleData
	var subsidies []SubsidyData

	for i := 0; i < epochCount; i++ {
		epochNumber := big.NewInt(int64(i + 1))

		// Generate merkle data for this epoch
		merkleData := g.GenerateMerkleData(vaultID, epochNumber, merkleEntriesPerEpoch)
		merkles = append(merkles, merkleData)

		// Generate subsidy data for this epoch
		for j := 0; j < subsidiesPerEpoch; j++ {
			subsidyData := g.GenerateSubsidyData(vaultID, epochNumber)
			subsidies = append(subsidies, subsidyData)
		}
	}

	return TestDataSet{
		VaultID:   vaultID,
		Epochs:    epochs,
		Merkles:   merkles,
		Subsidies: subsidies,
	}
}

// ToJSON converts data to JSON bytes
func (e EpochData) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (m MerkleData) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (s SubsidyData) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// TestScenarios provides common test scenarios
type TestScenarios struct {
	generator *TestDataGenerator
}

// NewTestScenarios creates a new test scenarios generator
func NewTestScenarios(seed int64) *TestScenarios {
	return &TestScenarios{
		generator: NewTestDataGenerator(seed),
	}
}

// BasicDataFlow creates a basic data flow scenario
func (ts *TestScenarios) BasicDataFlow(vaultID string) TestDataSet {
	return ts.generator.GenerateTestDataSet(vaultID, 5, 10, 2)
}

// HighVolumeScenario creates a high volume data scenario
func (ts *TestScenarios) HighVolumeScenario(vaultID string) TestDataSet {
	return ts.generator.GenerateTestDataSet(vaultID, 100, 100, 10)
}

// ConcurrencyScenario creates data for concurrency testing
func (ts *TestScenarios) ConcurrencyScenario(vaultCount int) []TestDataSet {
	var datasets []TestDataSet
	for i := 0; i < vaultCount; i++ {
		vaultID := ts.generator.GenerateVaultID()
		dataset := ts.generator.GenerateTestDataSet(vaultID, 10, 20, 5)
		datasets = append(datasets, dataset)
	}
	return datasets
}

// StressTestScenario creates data for stress testing
func (ts *TestScenarios) StressTestScenario(vaultCount int) []TestDataSet {
	var datasets []TestDataSet
	for i := 0; i < vaultCount; i++ {
		vaultID := ts.generator.GenerateVaultID()
		dataset := ts.generator.GenerateTestDataSet(vaultID, 1000, 1000, 50)
		datasets = append(datasets, dataset)
	}
	return datasets
}

// EdgeCaseScenario creates edge case test data
func (ts *TestScenarios) EdgeCaseScenario(vaultID string) TestDataSet {
	// Generate data with edge cases
	epochs := []EpochData{
		// Zero epoch
		ts.generator.GenerateEpochData(vaultID, big.NewInt(0)),
		// Large epoch number
		ts.generator.GenerateEpochData(vaultID, big.NewInt(999999999)),
		// Current epoch
		ts.generator.GenerateEpochData(vaultID, big.NewInt(time.Now().Unix())),
	}

	// Merkle with zero entries
	merkleEmpty := ts.generator.GenerateMerkleData(vaultID, big.NewInt(0), 0)

	// Merkle with single entry
	merkleSingle := ts.generator.GenerateMerkleData(vaultID, big.NewInt(1), 1)

	// Merkle with max entries
	merkleMax := ts.generator.GenerateMerkleData(vaultID, big.NewInt(999999999), 10000)

	merkles := []MerkleData{merkleEmpty, merkleSingle, merkleMax}

	// Subsidy with zero amount
	subsidyZero := ts.generator.GenerateSubsidyData(vaultID, big.NewInt(0))
	subsidyZero.Amount = big.NewInt(0)

	// Subsidy with max amount
	subsidyMax := ts.generator.GenerateSubsidyData(vaultID, big.NewInt(999999999))
	subsidyMax.Amount, _ = new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

	subsidies := []SubsidyData{subsidyZero, subsidyMax}

	return TestDataSet{
		VaultID:   vaultID,
		Epochs:    epochs,
		Merkles:   merkles,
		Subsidies: subsidies,
	}
}
