package types

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/stretchr/testify/require"
)

const testRelayerSetClaimChain = "testchain942"

func init() {
	RegisterExternalAddress(testRelayerSetClaimChain, EthereumAddress{})
}

func TestMsgRelayerSetUpdateClaimValidateBasicRejectsZeroNonce(t *testing.T) {
	relayer := sdk.AccAddress(bytes.Repeat([]byte{0x01}, 20)).String()
	msg := &MsgRelayerSetUpdateClaim{
		EventNonce:      1,
		BlockHeight:     1,
		RelayerSetNonce: 0,
		Members: BridgeValidators{
			{
				Power:           1,
				ExternalAddress: "0x0000000000000000000000000000000000000001",
			},
		},
		RelayerAddress: relayer,
		ChainName:      testRelayerSetClaimChain,
	}

	require.ErrorContains(t, msg.ValidateBasic(), "zero relayer set nonce")

	msg.RelayerSetNonce = 1
	require.NoError(t, msg.ValidateBasic())
}
