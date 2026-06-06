package types

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/openmetaearth/me-hub/app/params"
	"github.com/stretchr/testify/require"
)

const testBridgeTokenClaimChain = "testchain941"

func init() {
	RegisterExternalAddress(testBridgeTokenClaimChain, EthereumAddress{})
}

func TestMsgBridgeTokenClaimValidateBasicRejectsOversizedDecimals(t *testing.T) {
	relayer := sdk.AccAddress(bytes.Repeat([]byte{0x01}, 20)).String()
	msg := &MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1,
		TokenContract:  "0x0000000000000000000000000000000000000001",
		Name:           "Token",
		Symbol:         "TOK",
		Decimals:       MaxBridgeTokenDecimals + 1,
		RelayerAddress: relayer,
		ChainName:      testBridgeTokenClaimChain,
	}

	require.ErrorContains(t, msg.ValidateBasic(), "bridge token decimals")

	msg.Decimals = MaxBridgeTokenDecimals
	require.NoError(t, msg.ValidateBasic())
}
