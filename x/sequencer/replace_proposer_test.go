package sequencer_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/openmetaearth/me-hub/testutil/keeper"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
	"github.com/stretchr/testify/require"
)

func TestReplaceProposerQueuesOldSequencerForUnbonding(t *testing.T) {
	keeper, ctx := keepertest.SequencerKeeper(t)
	now := time.Date(2026, 6, 6, 12, 0, 0, 0, time.UTC)
	ctx = ctx.WithBlockHeight(42).WithBlockTime(now)

	rollappID := "rollapp-1"
	oldProposer := "old-proposer"
	newProposer := "new-proposer"
	bond := sdk.NewCoins(sdk.NewCoin("dym", sdk.NewInt(100)))

	keeper.SetSequencer(ctx, types.Sequencer{
		SequencerAddress: oldProposer,
		RollappId:        rollappID,
		Proposer:         true,
		Status:           types.Bonded,
		Tokens:           bond,
	})
	keeper.SetSequencer(ctx, types.Sequencer{
		SequencerAddress: newProposer,
		RollappId:        rollappID,
		Status:           types.Bonded,
		Tokens:           bond,
	})
	require.NoError(t, keeper.SetReplaceProposer(ctx, &types.MsgRepalceProposer{
		RollappId:   rollappID,
		OldProposer: oldProposer,
		NewProposer: newProposer,
		BlockHeight: 10,
	}))

	err := keeper.ProcSequencerByPendingStates(ctx, rollappID, newProposer, &rollapptypes.StateInfo{
		StartHeight: 11,
		NumBlocks:   1,
	})
	require.NoError(t, err)

	oldSeq, found := keeper.GetSequencer(ctx, oldProposer)
	require.True(t, found)
	require.False(t, oldSeq.Proposer)
	require.Equal(t, types.Unbonding, oldSeq.Status)
	require.Equal(t, int64(42), oldSeq.UnbondingHeight)
	require.Equal(t, now.Add(keeper.UnbondingTime(ctx)), oldSeq.UnbondTime)

	require.Empty(t, keeper.GetMatureUnbondingSequencers(ctx, oldSeq.UnbondTime.Add(-time.Nanosecond)))
	matureSequencers := keeper.GetMatureUnbondingSequencers(ctx, oldSeq.UnbondTime)
	require.Len(t, matureSequencers, 1)
	require.Equal(t, oldProposer, matureSequencers[0].SequencerAddress)

	newSeq, found := keeper.GetSequencer(ctx, newProposer)
	require.True(t, found)
	require.True(t, newSeq.Proposer)
}
