package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// ErrInvalidAddress is returned when an address is not a valid Ethereum address
	ErrInvalidAddress = errors.New("invalid Ethereum address format")
)

// NormalizeAddress converts an Ethereum address to lowercase for consistent comparison
// This ensures that address comparisons are case-insensitive, which is the standard
// for Ethereum addresses.
func NormalizeAddress(address string) string {
	return strings.ToLower(address)
}

// IsValidAddress checks if a string is a valid Ethereum address format
func IsValidAddress(address string) bool {
	if address == "" {
		return false
	}

	// Check if it's a valid hex string with 0x prefix and 40 hex characters
	matched, err := regexp.MatchString("^0x[0-9a-fA-F]{40}$", address)
	if err != nil {
		return false
	}

	if !matched {
		return false
	}

	// Use go-ethereum's common.IsHexAddress for additional validation
	return common.IsHexAddress(address)
}

// ValidateAndNormalizeAddress validates an address and returns it normalized to lowercase
// Returns an error if the address is invalid
func ValidateAndNormalizeAddress(address string) (string, error) {
	if !IsValidAddress(address) {
		return "", ErrInvalidAddress
	}
	return NormalizeAddress(address), nil
}
