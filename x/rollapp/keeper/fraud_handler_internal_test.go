package keeper

import (
	"testing"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	cometbfttypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/stretchr/testify/require"
)

func TestFrozenHeightFromLatestPreservesRevisionOrder(t *testing.T) {
	latest := clienttypes.NewHeight(2, 1234)
	clientState := &cometbfttypes.ClientState{LatestHeight: latest}

	frozen := frozenHeightFromLatest(clientState)

	require.Equal(t, latest, frozen)
	require.Equal(t, uint64(2), frozen.RevisionNumber)
	require.Equal(t, uint64(1234), frozen.RevisionHeight)
}
