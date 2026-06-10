package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetMintCoin(t *testing.T) {
	tests := []struct {
		name        string
		amount      sdk.Int
		chainName   string
		bridgeToken *BridgeToken
		expected    sdk.Coin
	}{
		{
			name:      "BSC USDT Mint: divides 18-decimal BSC amount to 6-decimal Cosmos coin",
			amount:    sdkmath.NewInt(1_000_000).Mul(sdkmath.NewInt(1_000_000_000_000)), // 1 USDT in 18-decimal BSC
			chainName: "bsc",
			bridgeToken: &BridgeToken{
				Denom:  "gravity0xToken",
				Symbol: "USDT",
			},
			expected: sdk.NewCoin("gravity0xToken", sdkmath.NewInt(1_000_000)),
		},
		{
			name:      "ETH USDT Mint: no scaling, remains 6-decimal",
			amount:    sdkmath.NewInt(1_000_000), // 1 USDT in 6-decimal from ETH
			chainName: "eth",
			bridgeToken: &BridgeToken{
				Denom:  "gravity0xToken",
				Symbol: "USDT",
			},
			expected: sdk.NewCoin("gravity0xToken", sdkmath.NewInt(1_000_000)),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetMintCoin(tc.amount, tc.chainName, tc.bridgeToken)
			require.Equal(t, tc.expected, result)
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
			name:      "BSC USDT Unlock: multiplies 6-decimal Cosmos amount to 18-decimal BSC amount",
			amount:    sdkmath.NewInt(1_000_000), // 1 USDT in 6-decimal
			chainName: "bsc",
			bridgeToken: &BridgeToken{
				Symbol: "USDT",
			},
			expected: sdkmath.NewInt(1_000_000).Mul(sdkmath.NewInt(1_000_000_000_000)),
		},
		{
			name:      "ETH USDT Unlock: no scaling, remains 6-decimal",
			amount:    sdkmath.NewInt(1_000_000),
			chainName: "eth",
			bridgeToken: &BridgeToken{
				Symbol: "USDT",
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
