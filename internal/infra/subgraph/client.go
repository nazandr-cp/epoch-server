package subgraph

// GraphQLRequest represents a GraphQL query request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
	Block     *BlockParameter        `json:"block,omitempty"`
}

// BlockParameter represents a block constraint for GraphQL queries
type BlockParameter struct {
	Number *int64  `json:"number,omitempty"`
	Hash   *string `json:"hash,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// Account represents a user account
type Account struct {
	ID                           string `json:"id"`
	TotalSecondsClaimed          string `json:"totalSecondsClaimed"`
	TotalSubsidiesReceived       string `json:"totalSubsidiesReceived"`
	TotalYieldEarned             string `json:"totalYieldEarned"`
	TotalBorrowVolume            string `json:"totalBorrowVolume"`
	TotalNFTsOwned               string `json:"totalNFTsOwned"`
	TotalCollectionsParticipated string `json:"totalCollectionsParticipated"`
	CreatedAtBlock               string `json:"createdAtBlock"`
	CreatedAtTimestamp           string `json:"createdAtTimestamp"`
	UpdatedAtBlock               string `json:"updatedAtBlock"`
	UpdatedAtTimestamp           string `json:"updatedAtTimestamp"`
}

// AccountSubsidy represents account subsidy data
type AccountSubsidy struct {
	ID                      string  `json:"id"`
	Account                 Account `json:"account"`
	AccountMarket           string  `json:"accountMarket"`
	CollectionParticipation string  `json:"collectionParticipation"`
	BalanceNFT              string  `json:"balanceNFT"`
	SecondsAccumulated      string  `json:"secondsAccumulated"`
	SecondsClaimed          string  `json:"secondsClaimed"`
	SubsidiesAccrued        string  `json:"subsidiesAccrued"`
	SubsidiesClaimed        string  `json:"subsidiesClaimed"`
	AverageHoldingPeriod    string  `json:"averageHoldingPeriod"`
	TotalRewardsEarned      string  `json:"totalRewardsEarned"`
	LastEffectiveValue      string  `json:"lastEffectiveValue"`
	UpdatedAtBlock          string  `json:"updatedAtBlock"`
	UpdatedAtTimestamp      string  `json:"updatedAtTimestamp"`
}

type Epoch struct {
	ID                            string `json:"id"`
	EpochNumber                   string `json:"epochNumber"`
	Status                        string `json:"status"`
	StartTimestamp                string `json:"startTimestamp"`
	EndTimestamp                  string `json:"endTimestamp"`
	ProcessingStartedTimestamp    string `json:"processingStartedTimestamp"`
	ProcessingCompletedTimestamp  string `json:"processingCompletedTimestamp"`
	TotalYieldAvailable           string `json:"totalYieldAvailable"`
	TotalYieldAllocated           string `json:"totalYieldAllocated"`
	TotalYieldDistributed         string `json:"totalYieldDistributed"`
	RemainingYield                string `json:"remainingYield"`
	TotalSubsidiesDistributed     string `json:"totalSubsidiesDistributed"`
	TotalEligibleUsers            string `json:"totalEligibleUsers"`
	TotalParticipatingCollections string `json:"totalParticipatingCollections"`
	ProcessingTimeMs              string `json:"processingTimeMs"`
	ProcessingGasUsed             string `json:"processingGasUsed"`
	ProcessingTransactionCount    string `json:"processingTransactionCount"`
	CreatedAtBlock                string `json:"createdAtBlock"`
	CreatedAtTimestamp            string `json:"createdAtTimestamp"`
	UpdatedAtBlock                string `json:"updatedAtBlock"`
	UpdatedAtTimestamp            string `json:"updatedAtTimestamp"`
}

type Collection struct {
	ID                   string `json:"id"`
	ContractAddress      string `json:"contractAddress"`
	Name                 string `json:"name"`
	Symbol               string `json:"symbol"`
	TotalSupply          string `json:"totalSupply"`
	CollectionType       string `json:"collectionType"`
	IsActive             bool   `json:"isActive"`
	YieldSharePercentage string `json:"yieldSharePercentage"`
	WeightFunctionType   string `json:"weightFunctionType"`
	WeightFunctionP1     string `json:"weightFunctionP1"`
	WeightFunctionP2     string `json:"weightFunctionP2"`
	MinBorrowAmount      string `json:"minBorrowAmount"`
	MaxBorrowAmount      string `json:"maxBorrowAmount"`
	TotalNFTsDeposited   string `json:"totalNFTsDeposited"`
	UpdatedAtBlock       string `json:"updatedAtBlock"`
	UpdatedAtTimestamp   string `json:"updatedAtTimestamp"`
}

type MerkleDistribution struct {
	ID                 string `json:"id"`
	Epoch              Epoch  `json:"epoch"`
	Vault              string `json:"vault"`
	MerkleRoot         string `json:"merkleRoot"`
	TotalAmount        string `json:"totalAmount"`
	Timestamp          string `json:"timestamp"`
	UpdatedAtBlock     string `json:"updatedAtBlock"`
	UpdatedAtTimestamp string `json:"updatedAtTimestamp"`
}

// AccountsResponse represents the response for accounts query
type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

// AccountSubsidiesResponse represents account subsidies query response
type AccountSubsidiesResponse struct {
	AccountSubsidies []AccountSubsidy `json:"accountSubsidies"`
}

// EpochsResponse represents epochs query response
type EpochsResponse struct {
	Epochs []Epoch `json:"epoches"`
}

// MerkleDistributionsResponse represents merkle distributions query response
type MerkleDistributionsResponse struct {
	MerkleDistributions []MerkleDistribution `json:"merkleDistributions"`
}

// UsersResponse for backward compatibility
type UsersResponse struct {
	Users []Account `json:"users"`
}

