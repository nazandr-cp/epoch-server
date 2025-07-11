package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeAddress(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "uppercase address",
			input:    "0x742D35CC6BF8E65F8B95E6C5CB15F5C5D5B8DBC3",
			expected: "0x742d35cc6bf8e65f8b95e6c5cb15f5c5d5b8dbc3",
		},
		{
			name:     "mixed case address",
			input:    "0xAbCdEfAbCdEfAbCdEfAbCdEfAbCdEfAbCdEfAbCd",
			expected: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		},
		{
			name:     "already lowercase",
			input:    "0x1234567890123456789012345678901234567890",
			expected: "0x1234567890123456789012345678901234567890",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAddress(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidAddress(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid lowercase address",
			input:    "0x742d35cc6bf8e65f8b95e6c5cb15f5c5d5b8dbc3",
			expected: true,
		},
		{
			name:     "valid uppercase address",
			input:    "0x742D35CC6BF8E65F8B95E6C5CB15F5C5D5B8DBC3",
			expected: true,
		},
		{
			name:     "valid mixed case address",
			input:    "0xAbCdEfAbCdEfAbCdEfAbCdEfAbCdEfAbCdEfAbCd",
			expected: true,
		},
		{
			name:     "invalid - missing 0x prefix",
			input:    "742d35cc6bf8e65f8b95e6c5cb15f5c5d5b8dbc3",
			expected: false,
		},
		{
			name:     "invalid - too short",
			input:    "0x742d35cc6bf8e65f8b95e6c5cb15f5c5d5b8db",
			expected: false,
		},
		{
			name:     "invalid - too long",
			input:    "0x742d35cc6bf8e65f8b95e6c5cb15f5c5d5b8dbc333",
			expected: false,
		},
		{
			name:     "invalid - non-hex characters",
			input:    "0x742d35cc6bf8e65f8b95e6c5cb15f5c5d5b8dbgg",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "zero address",
			input:    "0x0000000000000000000000000000000000000000",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidAddress(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateAndNormalizeAddress(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		shouldError bool
	}{
		{
			name:        "valid uppercase address",
			input:       "0x742D35CC6BF8E65F8B95E6C5CB15F5C5D5B8DBC3",
			expected:    "0x742d35cc6bf8e65f8b95e6c5cb15f5c5d5b8dbc3",
			shouldError: false,
		},
		{
			name:        "valid mixed case address",
			input:       "0xAbCdEfAbCdEfAbCdEfAbCdEfAbCdEfAbCdEfAbCd",
			expected:    "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			shouldError: false,
		},
		{
			name:        "invalid address - missing 0x",
			input:       "742d35cc6bf8e65f8b95e6c5cb15f5c5d5b8dbc3",
			expected:    "",
			shouldError: true,
		},
		{
			name:        "invalid address - too short",
			input:       "0x742d35cc6bf8e65f8b95e6c5cb15f5c5d5b8db",
			expected:    "",
			shouldError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateAndNormalizeAddress(tt.input)

			if tt.shouldError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidAddress, err)
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
