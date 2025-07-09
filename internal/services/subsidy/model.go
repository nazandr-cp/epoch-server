package subsidy

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
