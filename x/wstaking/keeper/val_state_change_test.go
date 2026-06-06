package keeper

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	ed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/stretchr/testify/require"
)

func TestDeduplicateValidatorUpdatesKeepsFinalPubKeyPower(t *testing.T) {
	oldPubKey, err := cryptocodec.ToTmProtoPublicKey(ed25519.GenPrivKey().PubKey())
	require.NoError(t, err)
	newPubKey, err := cryptocodec.ToTmProtoPublicKey(ed25519.GenPrivKey().PubKey())
	require.NoError(t, err)

	updates := []abci.ValidatorUpdate{
		{
			PubKey: oldPubKey,
			Power:  200,
		},
		{
			PubKey: oldPubKey,
			Power:  0,
		},
		{
			PubKey: newPubKey,
			Power:  200,
		},
	}

	deduplicated := deduplicateValidatorUpdates(updates)

	require.Len(t, deduplicated, 2)
	require.Equal(t, oldPubKey, deduplicated[0].PubKey)
	require.EqualValues(t, 0, deduplicated[0].Power)
	require.Equal(t, newPubKey, deduplicated[1].PubKey)
	require.EqualValues(t, 200, deduplicated[1].Power)
}
