package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	cometbfttypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/stretchr/testify/require"
)

type testIBCClientKeeper struct {
	states map[string]exported.ClientState
}

func (k *testIBCClientKeeper) GetClientState(_ sdk.Context, clientID string) (exported.ClientState, bool) {
	state, found := k.states[clientID]
	return state, found
}

func (k *testIBCClientKeeper) SetClientState(_ sdk.Context, clientID string, clientState exported.ClientState) {
	k.states[clientID] = clientState
}

func TestFreezeClientStateUsesLatestHeightRevisionOrder(t *testing.T) {
	clientID := "test-client-id"
	clientKeeper := &testIBCClientKeeper{
		states: map[string]exported.ClientState{
			clientID: &cometbfttypes.ClientState{
				LatestHeight: clienttypes.NewHeight(1, 100000),
			},
		},
	}
	k := Keeper{
		ibcClientKeeper: clientKeeper,
	}

	err := k.freezeClientState(sdk.Context{}, clientID)
	require.NoError(t, err)

	updatedState, found := clientKeeper.GetClientState(sdk.Context{}, clientID)
	require.True(t, found)

	tmClientState, ok := updatedState.(*cometbfttypes.ClientState)
	require.True(t, ok)

	require.Equal(t, uint64(1), tmClientState.FrozenHeight.GetRevisionNumber())
	require.Equal(t, uint64(100000), tmClientState.FrozenHeight.GetRevisionHeight())
}
