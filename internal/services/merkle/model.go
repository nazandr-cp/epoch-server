package merkle

// UserMerkleProofResponse represents a merkle proof response for a user
type UserMerkleProofResponse struct {
	UserAddress   string   `json:"userAddress"`
	VaultAddress  string   `json:"vaultAddress"`
	EpochNumber   string   `json:"epochNumber,omitempty"`
	TotalEarned   string   `json:"totalEarned"`
	MerkleProof   []string `json:"merkleProof"`
	MerkleRoot    string   `json:"merkleRoot"`
	LeafIndex     int      `json:"leafIndex"`
	GeneratedAt   int64    `json:"generatedAt"`
}

// MerkleDistribution represents merkle distribution data for an epoch
type MerkleDistribution struct {
	EpochNumber        string   `json:"epochNumber"`
	VaultAddress       string   `json:"vaultAddress"`
	MerkleRoot         string   `json:"merkleRoot"`
	TotalSubsidies     string   `json:"totalSubsidies"`
	AccountsProcessed  int      `json:"accountsProcessed"`
	Proofs             []string `json:"proofs"`
	CreatedAt          int64    `json:"createdAt"`
}