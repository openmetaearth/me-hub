package keeper

import (
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	cometbfttypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app/apptesting"
)

func TestFreezeClientState(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	k := app.RollappKeeper
	clientID := "test-client-id"

	// Create and set a tendermint client state
	latestHeight := clienttypes.NewHeight(1, 100000) // revision 1, height 100000
	clientState := &cometbfttypes.ClientState{
		LatestHeight: latestHeight,
	}
	app.IBCKeeper.ClientKeeper.SetClientState(ctx, clientID, clientState)

	// Call freezeClientState (unexported helper)
	err := k.freezeClientState(ctx, clientID)
	require.NoError(t, err)

	// Retrieve the updated client state
	updatedState, found := app.IBCKeeper.ClientKeeper.GetClientState(ctx, clientID)
	require.True(t, found)

	tmClientState, ok := updatedState.(*cometbfttypes.ClientState)
	require.True(t, ok)

	// Verify FrozenHeight
	// revision number should be 1 (revisionNumber from latestHeight)
	// revision height should be 100000 (revisionHeight from latestHeight)
	require.Equal(t, uint64(1), tmClientState.FrozenHeight.GetRevisionNumber())
	require.Equal(t, uint64(100000), tmClientState.FrozenHeight.GetRevisionHeight())
}
