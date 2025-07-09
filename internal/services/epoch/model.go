package epoch

// UserEarningsResponse represents the response for user total earned query
type UserEarningsResponse struct {
	UserAddress   string `json:"userAddress"`
	VaultAddress  string `json:"vaultAddress"`
	TotalEarned   string `json:"totalEarned"`
	CalculatedAt  int64  `json:"calculatedAt"`
	DataTimestamp int64  `json:"dataTimestamp"` // Timestamp used for calculations
}
