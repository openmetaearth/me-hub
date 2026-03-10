package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetExternalUnlockAmount(t *testing.T) {
	convert := sdk.NewDec(10).Power(6 - 6).TruncateInt()
	require.Equal(t, sdkmath.NewInt(1), convert, "expected convert to be 1 when decimals are equal")
	tests := []struct {
		name        string
		amount      sdk.Int
		chainName   string
		bridgeToken *BridgeToken
		expected    sdk.Int
	}{
		{
			name:      "BSC USDT with 18 decimals: converts 6-decimal amount to 18-decimal",
			amount:    sdkmath.NewInt(1_000_000), // 1 USDT in 6-decimal
			chainName: "bsc",
			bridgeToken: &BridgeToken{
				Symbol:  "USDT",
				Decimal: 18,
			},
			// 1_000_000 * 10^(18-6) = 1_000_000 * 10^12
			expected: sdkmath.NewInt(1_000_000).Mul(sdkmath.NewInt(1_000_000_000_000)),
		},
		{
			name:      "BSC USDC with 18 decimals: converts 6-decimal amount to 18-decimal",
			amount:    sdkmath.NewInt(5_000_000), // 5 USDC in 6-decimal
			chainName: "bsc",                     // case-insensitive
			bridgeToken: &BridgeToken{
				Symbol:  "USDC",
				Decimal: 18,
			},
			expected: sdkmath.NewInt(5_000_000).Mul(sdkmath.NewInt(1_000_000_000_000)),
		},
		{
			name:      "BSC USDT with decimal equal to 6: no conversion (multiply by 10^0 = 1)",
			amount:    sdkmath.NewInt(1_000_000),
			chainName: "bsc",
			bridgeToken: &BridgeToken{
				Symbol:  "USDT",
				Decimal: 6,
			},
			expected: sdkmath.NewInt(1_000_000),
		},
		{
			name:      "non-BSC chain (ETH) USDT: no conversion",
			amount:    sdkmath.NewInt(1_000_000),
			chainName: "bsc",
			bridgeToken: &BridgeToken{
				Symbol:  "USDT",
				Decimal: 6,
			},
			expected: sdkmath.NewInt(1_000_000),
		},
		{
			name:      "BSC non-USDT/USDC token: no conversion",
			amount:    sdkmath.NewInt(2_000_000),
			chainName: "bsc",
			bridgeToken: &BridgeToken{
				Symbol:  "BNB",
				Decimal: 18,
			},
			expected: sdkmath.NewInt(2_000_000),
		},
		{
			name:      "zero amount BSC USDT",
			amount:    sdkmath.ZeroInt(),
			chainName: "bsc",
			bridgeToken: &BridgeToken{
				Symbol:  "USDT",
				Decimal: 18,
			},
			expected: sdkmath.ZeroInt(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetExternalUnlockAmount(tc.amount, tc.chainName, tc.bridgeToken)
			require.True(t, tc.expected.Equal(result),
				"expected %s, got %s", tc.expected.String(), result.String())
		})
	}
}
