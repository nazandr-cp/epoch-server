package utils

import "strings"

// NormalizeAddress converts an Ethereum address to lowercase for consistent comparison
// This ensures that address comparisons are case-insensitive, which is the standard
// for Ethereum addresses.
func NormalizeAddress(address string) string {
	return strings.ToLower(address)
}
