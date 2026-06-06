package keeper_test

import (
	"testing"

	keepertest "github.com/openmetaearth/me-hub/testutil/keeper"
	common "github.com/openmetaearth/me-hub/x/common/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/stretchr/testify/require"
)

func TestRevertPendingStatesSkipsMissingStateInfo(t *testing.T) {
	k, ctx := keepertest.RollappKeeper(t)

	const (
		fraudRollappID = "fraud-rollapp"
		otherRollappID = "other-rollapp"
		creationHeight = uint64(10)
	)

	otherState := types.StateInfo{
		StateInfoIndex: types.StateInfoIndex{RollappId: otherRollappID, Index: 1},
		Status:         common.Status_PENDING,
	}
	k.SetStateInfo(ctx, otherState)
	k.SetBlockHeightToFinalizationQueue(ctx, types.BlockHeightToFinalizationQueue{
		CreationHeight: creationHeight,
		FinalizationQueue: []types.StateInfoIndex{
			{RollappId: fraudRollappID, Index: 7},
			otherState.StateInfoIndex,
		},
	})

	k.RevertPendingStates(ctx, fraudRollappID)

	_, found := k.GetStateInfo(ctx, fraudRollappID, 7)
	require.False(t, found)
	_, found = k.GetStateInfo(ctx, "", 0)
	require.False(t, found)

	storedOtherState, found := k.GetStateInfo(ctx, otherRollappID, 1)
	require.True(t, found)
	require.Equal(t, common.Status_PENDING, storedOtherState.Status)

	queue, found := k.GetBlockHeightToFinalizationQueue(ctx, creationHeight)
	require.True(t, found)
	require.Equal(t, []types.StateInfoIndex{otherState.StateInfoIndex}, queue.FinalizationQueue)
}
