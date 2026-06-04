package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBridgeTokenClaimHashSeparatesNameAndSymbol(t *testing.T) {
	base := MsgBridgeTokenClaim{
		BlockHeight:   100,
		EventNonce:    7,
		TokenContract: "0x0000000000000000000000000000000000000001",
		Decimals:      18,
		ChainName:     "bsc",
	}

	first := base
	first.Name = "AAA"
	first.Symbol = "BBB/CCC"

	second := base
	second.Name = "AAA/BBB"
	second.Symbol = "CCC"

	require.False(t, bytes.Equal(first.ClaimHash(), second.ClaimHash()))
}

func TestBridgeTokenClaimHashIncludesChainName(t *testing.T) {
	first := MsgBridgeTokenClaim{
		BlockHeight:   100,
		EventNonce:    7,
		TokenContract: "0x0000000000000000000000000000000000000001",
		Name:          "AAA",
		Symbol:        "BBB",
		Decimals:      18,
		ChainName:     "bsc",
	}
	second := first
	second.ChainName = "tron"

	require.False(t, bytes.Equal(first.ClaimHash(), second.ClaimHash()))
}
