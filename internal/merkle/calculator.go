package merkle

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/andrey/epoch-server/internal/clients/graph"
)

// Calculator provides shared calculation utilities for merkle tree operations
type Calculator struct{}

// NewCalculator creates a new calculator instance
func NewCalculator() *Calculator {
	return &Calculator{}
}

// CalculateTotalEarned calculates the total earned for an account using the provided timestamp
// This is the unified calculation logic used across all services
func (c *Calculator) CalculateTotalEarned(subsidy graph.AccountSubsidy, endTimestamp int64) (*big.Int, error) {
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
	totalEarned := c.SecondsToTokens(newTotalSeconds)
	return totalEarned, nil
}

// SecondsToTokens converts seconds to token amounts using the standard conversion rate
func (c *Calculator) SecondsToTokens(seconds *big.Int) *big.Int {
	conversionRate := big.NewInt(1000000000000000000) // 1e18
	return new(big.Int).Div(seconds, conversionRate)
}

// CalculateTotalSubsidies calculates the total subsidies from a list of entries
func (c *Calculator) CalculateTotalSubsidies(entries []Entry) *big.Int {
	totalSubsidies := big.NewInt(0)
	for _, entry := range entries {
		totalSubsidies.Add(totalSubsidies, entry.TotalEarned)
	}
	return totalSubsidies
}

// ProcessAccountSubsidies processes account subsidies and returns entries with positive earnings
func (c *Calculator) ProcessAccountSubsidies(subsidies []graph.AccountSubsidy, endTimestamp int64) ([]Entry, error) {
	var entries []Entry
	
	for _, subsidy := range subsidies {
		totalEarned, err := c.CalculateTotalEarned(subsidy, endTimestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate total earned for account %s: %w", subsidy.Account.ID, err)
		}

		// Only include accounts with positive earnings
		if totalEarned.Cmp(big.NewInt(0)) > 0 {
			entries = append(entries, Entry{
				Address:     subsidy.Account.ID,
				TotalEarned: totalEarned,
			})
		}
	}
	
	return entries, nil
}