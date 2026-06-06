package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSignersPanicsOnInvalidAddresses(t *testing.T) {
	require.Panics(t, func() {
		MsgNewClass{Sender: "not-a-bech32-address"}.GetSigners()
	})

	require.Panics(t, func() {
		MsgMintNFT{Creator: "not-a-bech32-address"}.GetSigners()
	})

	require.Panics(t, func() {
		MsgSend{Sender: "not-a-bech32-address"}.GetSigners()
	})
}
