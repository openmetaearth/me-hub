package types

import (
	"encoding/hex"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestClaimHashBindsChainName(t *testing.T) {
	t.Run("send to me", func(t *testing.T) {
		claim := &MsgSendToMeClaim{
			EventNonce:    2,
			BlockHeight:   100,
			TokenContract: "0x0000000000000000000000000000000000000001",
			Sender:        "0x0000000000000000000000000000000000000002",
			Amount:        sdkmath.NewInt(100),
			Receiver:      "me1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a",
			ChainName:     "bsc",
		}
		other := *claim
		other.ChainName = "gravity"

		require.NotEqual(t, hex.EncodeToString(claim.ClaimHash()), hex.EncodeToString(other.ClaimHash()))
	})

	t.Run("send to external", func(t *testing.T) {
		claim := &MsgSendToExternalClaim{
			EventNonce:    3,
			BlockHeight:   101,
			BatchNonce:    5,
			TokenContract: "0x0000000000000000000000000000000000000001",
			ChainName:     "bsc",
		}
		other := *claim
		other.ChainName = "gravity"

		require.NotEqual(t, hex.EncodeToString(claim.ClaimHash()), hex.EncodeToString(other.ClaimHash()))
	})

	t.Run("bridge token", func(t *testing.T) {
		claim := &MsgBridgeTokenClaim{
			EventNonce:    4,
			BlockHeight:   102,
			TokenContract: "0x0000000000000000000000000000000000000001",
			Name:          "Tether USD",
			Symbol:        "USDT",
			Decimals:      18,
			ChainName:     "bsc",
		}
		other := *claim
		other.ChainName = "gravity"

		require.NotEqual(t, hex.EncodeToString(claim.ClaimHash()), hex.EncodeToString(other.ClaimHash()))
	})

	t.Run("relayer set update", func(t *testing.T) {
		claim := &MsgRelayerSetUpdateClaim{
			EventNonce:      5,
			BlockHeight:     103,
			RelayerSetNonce: 7,
			Members: []BridgeValidator{
				{Power: 100, ExternalAddress: "0x0000000000000000000000000000000000000001"},
			},
			ChainName: "bsc",
		}
		other := *claim
		other.ChainName = "gravity"

		require.NotEqual(t, hex.EncodeToString(claim.ClaimHash()), hex.EncodeToString(other.ClaimHash()))
	})
}
