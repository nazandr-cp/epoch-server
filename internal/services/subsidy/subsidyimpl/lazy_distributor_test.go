package subsidyimpl

import (
	"math/big"
	"testing"

	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andrey/epoch-server/internal/infra/subgraph"
)

func TestLazyDistributor_CalculateTotalEarned(t *testing.T) {
	logger := lgr.Default()
	distributor := &LazyDistributor{
		logger: logger,
	}

	tests := []struct {
		name             string
		subsidy          subgraph.AccountSubsidy
		endTimestamp     int64
		expectedEarnings string
		expectError      bool
	}{
		{
			name: "real_user_1_significant_seconds",
			subsidy: subgraph.AccountSubsidy{
				Account: subgraph.Account{
					ID: "0x8f37c5c4fa708e06a656d858003ef7dc5f60a29b",
				},
				SecondsAccumulated: "439236",
				LastEffectiveValue: "9000000000000000000",
				UpdatedAtTimestamp: "1752211061",
				TotalRewardsEarned: "0",
			},
			endTimestamp:     1752211061 + 3600,
			expectedEarnings: "32400",
			expectError:      false,
		},
		{
			name: "real_user_2_higher_seconds",
			subsidy: subgraph.AccountSubsidy{
				Account: subgraph.Account{
					ID: "0x3575b992c5337226aecf4e7f93dfbe80c576ce15",
				},
				SecondsAccumulated: "1024884",
				LastEffectiveValue: "21000000000000000000",
				UpdatedAtTimestamp: "1752211061",
				TotalRewardsEarned: "0",
			},
			endTimestamp:     1752211061 + 3600,
			expectedEarnings: "75600",
			expectError:      false,
		},
		{
			name: "zero_effective_value",
			subsidy: subgraph.AccountSubsidy{
				Account: subgraph.Account{
					ID: "0x1234567890123456789012345678901234567890",
				},
				SecondsAccumulated: "1000000",
				LastEffectiveValue: "0",
				UpdatedAtTimestamp: "1752211061",
				TotalRewardsEarned: "0",
			},
			endTimestamp:     1752211061 + 3600,
			expectedEarnings: "0",
			expectError:      false,
		},
		{
			name: "invalid_seconds_accumulated",
			subsidy: subgraph.AccountSubsidy{
				Account: subgraph.Account{
					ID: "0x1234567890123456789012345678901234567890",
				},
				SecondsAccumulated: "invalid",
				LastEffectiveValue: "9000000000000000000",
				UpdatedAtTimestamp: "1752211061",
				TotalRewardsEarned: "0",
			},
			endTimestamp: 1752211061 + 3600,
			expectError:  true,
		},
		{
			name: "invalid_effective_value",
			subsidy: subgraph.AccountSubsidy{
				Account: subgraph.Account{
					ID: "0x1234567890123456789012345678901234567890",
				},
				SecondsAccumulated: "1000000",
				LastEffectiveValue: "invalid",
				UpdatedAtTimestamp: "1752211061",
				TotalRewardsEarned: "0",
			},
			endTimestamp: 1752211061 + 3600,
			expectError:  true,
		},
		{
			name: "same_timestamp_no_additional_time",
			subsidy: subgraph.AccountSubsidy{
				Account: subgraph.Account{
					ID: "0x1234567890123456789012345678901234567890",
				},
				SecondsAccumulated: "1000000000000000000",
				LastEffectiveValue: "9000000000000000000",
				UpdatedAtTimestamp: "1752211061",
				TotalRewardsEarned: "0",
			},
			endTimestamp:     1752211061,
			expectedEarnings: "1",
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := distributor.calculateTotalEarned(tt.subsidy, tt.endTimestamp)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedEarnings, result.String())
		})
	}
}

func TestLazyDistributor_ConvertSubsidiesToEntries(t *testing.T) {
	logger := lgr.Default()
	distributor := &LazyDistributor{
		logger: logger,
	}

	t.Run("mixed_totalRewardsEarned_scenarios", func(t *testing.T) {
		subsidies := []subgraph.AccountSubsidy{
			{
				Account: subgraph.Account{
					ID: "0x8f37c5c4fa708e06a656d858003ef7dc5f60a29b",
				},
				SecondsAccumulated: "439236",
				LastEffectiveValue: "9000000000000000000",
				UpdatedAtTimestamp: "1752211061",
				TotalRewardsEarned: "0",
			},
			{
				Account: subgraph.Account{
					ID: "0x3575b992c5337226aecf4e7f93dfbe80c576ce15",
				},
				SecondsAccumulated: "1024884",
				LastEffectiveValue: "21000000000000000000",
				UpdatedAtTimestamp: "1752211061",
				TotalRewardsEarned: "50000000000000000000",
			},
			{
				Account: subgraph.Account{
					ID: "0x1111111111111111111111111111111111111111",
				},
				SecondsAccumulated: "0",
				LastEffectiveValue: "0",
				UpdatedAtTimestamp: "1752211061",
				TotalRewardsEarned: "0",
			},
		}

		entries, totalSubsidies, err := distributor.convertSubsidiesToEntries(subsidies)

		require.NoError(t, err)
		assert.Len(t, entries, 2, "Should have 2 valid entries (excluding zero earnings)")

		assert.Equal(t, "0x8f37c5c4fa708e06a656d858003ef7dc5f60a29b", entries[0].Address)
		assert.True(t, entries[0].TotalEarned.Sign() > 0)

		assert.Equal(t, "0x3575b992c5337226aecf4e7f93dfbe80c576ce15", entries[1].Address)
		expectedPreCalculated := new(big.Int)
		expectedPreCalculated.SetString("50000000000000000000", 10)
		assert.Equal(t, expectedPreCalculated, entries[1].TotalEarned)

		expectedTotal := new(big.Int).Add(entries[0].TotalEarned, entries[1].TotalEarned)
		assert.Equal(t, expectedTotal, totalSubsidies)
	})

	t.Run("all_zero_totalRewardsEarned", func(t *testing.T) {
		subsidies := []subgraph.AccountSubsidy{
			{
				Account: subgraph.Account{
					ID: "0x8f37c5c4fa708e06a656d858003ef7dc5f60a29b",
				},
				SecondsAccumulated: "439236",
				LastEffectiveValue: "9000000000000000000",
				UpdatedAtTimestamp: "1752211061",
				TotalRewardsEarned: "",
			},
			{
				Account: subgraph.Account{
					ID: "0x3575b992c5337226aecf4e7f93dfbe80c576ce15",
				},
				SecondsAccumulated: "1024884",
				LastEffectiveValue: "21000000000000000000",
				UpdatedAtTimestamp: "1752211061",
				TotalRewardsEarned: "0",
			},
		}

		entries, totalSubsidies, err := distributor.convertSubsidiesToEntries(subsidies)

		require.NoError(t, err)
		assert.Len(t, entries, 2, "Should have 2 valid entries from fallback calculations")
		assert.True(t, totalSubsidies.Sign() > 0, "Total subsidies should be positive from fallback calculations")

		for i, entry := range entries {
			assert.True(t, entry.TotalEarned.Sign() > 0, "Entry %d should have positive earnings", i)
		}
	})
}

func TestLazyDistributor_ConvertSubsidiesToEntries_RealData(t *testing.T) {
	logger := lgr.Default()
	distributor := &LazyDistributor{
		logger: logger,
	}

	subsidies := []subgraph.AccountSubsidy{
		{
			Account: subgraph.Account{
				ID: "0x8f37c5c4fa708e06a656d858003ef7dc5f60a29b",
			},
			SecondsAccumulated: "439236",
			LastEffectiveValue: "9000000000000000000",
			UpdatedAtTimestamp: "1752211061",
			TotalRewardsEarned: "0",
		},
		{
			Account: subgraph.Account{
				ID: "0x3575b992c5337226aecf4e7f93dfbe80c576ce15",
			},
			SecondsAccumulated: "1024884",
			LastEffectiveValue: "21000000000000000000",
			UpdatedAtTimestamp: "1752211061",
			TotalRewardsEarned: "0",
		},
	}

	entries, totalSubsidies, err := distributor.convertSubsidiesToEntries(subsidies)

	require.NoError(t, err)
	assert.Len(t, entries, 2, "Should convert both real users to valid entries")
	assert.True(t, totalSubsidies.Sign() > 0, "Total subsidies should be positive")

	assert.Equal(t, "0x8f37c5c4fa708e06a656d858003ef7dc5f60a29b", entries[0].Address)
	assert.Equal(t, "0x3575b992c5337226aecf4e7f93dfbe80c576ce15", entries[1].Address)

	assert.True(t, entries[0].TotalEarned.Sign() > 0, "User 1 should have positive earnings")
	assert.True(t, entries[1].TotalEarned.Sign() > 0, "User 2 should have positive earnings")

	assert.True(t, entries[1].TotalEarned.Cmp(entries[0].TotalEarned) > 0,
		"User 2 should have higher earnings than User 1")

	t.Logf("User 1 earnings: %s", entries[0].TotalEarned.String())
	t.Logf("User 2 earnings: %s", entries[1].TotalEarned.String())
	t.Logf("Total subsidies: %s", totalSubsidies.String())
}
