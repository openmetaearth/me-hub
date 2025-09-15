package helpers_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func TestPrivKeyFromMnemonic(t *testing.T) {
	mnemonic := "afraid bottom consider camera dragon seek text other addict shift tent zebra feel swim approve prize promote secret process method quick engine return pitch"
	privKey, err := helpers.PrivKeyFromMnemonic(mnemonic, hd.Secp256k1Type, 0, 0)
	require.NoError(t, err)

	addr, err := sdk.AccAddressFromHexUnsafe("3E557FB31C0C94D05C2B214AFAD6EA97A41A368E")
	require.NoError(t, err)
	require.Equal(t, addr, sdk.AccAddress(privKey.PubKey().Address().Bytes()))
}
