package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetMintAmount(t *testing.T) {
	tests := []struct {
		name        string
		amount      sdk.Int
		chainName   string
		bridgeToken *BridgeToken
		expected    sdk.Int
	}{
		{
			name:      "BSC USDT Mint (6 to 18): multiply by 10^12",
			amount:    sdkmath.NewInt(1_000_000), // 1 USDT in 6-decimal from BSC
			chainName: "bsc",
			bridgeToken: &BridgeToken{
				Symbol:  "USDT",
				Decimal: 18,
			},
			// 1_000_000 * 10^(18-6)
			expected: sdkmath.NewInt(1_000_000).Mul(sdkmath.NewInt(1_000_000_000_000)),
		},
		{
			name:      "ETH USDT Mint (no special scaling): no change",
			amount:    sdkmath.NewInt(1_000_000),
			chainName: "eth",
			bridgeToken: &BridgeToken{
				Symbol:  "USDT",
				Decimal: 6,
			},
			expected: sdkmath.NewInt(1_000_000),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetMintAmount(tc.amount, tc.chainName, tc.bridgeToken)
			require.True(t, tc.expected.Equal(result), "expected %s, got %s", tc.expected.String(), result.String())
		})
	}
}

func TestGetExternalUnlockAmount(t *testing.T) {
	tests := []struct {
		name        string
		amount      sdk.Int
		chainName   string
		bridgeToken *BridgeToken
		expected    sdk.Int
	}{
		{
			name:      "BSC USDT Unlock (18 to 6): divide by 10^12",
			amount:    sdkmath.NewInt(1_000_000).Mul(sdkmath.NewInt(1_000_000_000_000)), // 1 USDT in 18-decimal ME
			chainName: "bsc",
			bridgeToken: &BridgeToken{
				Symbol:  "USDT",
				Decimal: 18,
			},
			// 1_000_000_000_000_000_000 / 10^(18-6) = 1_000_000
			expected: sdkmath.NewInt(1_000_000),
		},
		{
			name:      "ETH USDT Unlock (no special scaling): no change",
			amount:    sdkmath.NewInt(1_000_000),
			chainName: "eth",
			bridgeToken: &BridgeToken{
				Symbol:  "USDT",
				Decimal: 6,
			},
			expected: sdkmath.NewInt(1_000_000),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetExternalUnlockAmount(tc.amount, tc.chainName, tc.bridgeToken)
			require.True(t, tc.expected.Equal(result), "expected %s, got %s", tc.expected.String(), result.String())
		})
	}
}
