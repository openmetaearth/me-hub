package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
	"github.com/openmetaearth/me-hub/utils/gerrc"
)

var _ rollapptypes.RollappHooks = rollappHook{}

type rollappHook struct {
	rollapptypes.StubRollappCreatedHooks
	k Keeper
}

func (k Keeper) RollappHooks() rollapptypes.RollappHooks {
	return rollappHook{k: k}
}

// BeforeUpdateState will reject if the caller is not proposer
// if lastStateUpdateBySequencer is true, validate that the sequencer is in the middle of a rotation
func (hook rollappHook) BeforeUpdateState(ctx sdk.Context, seqAddr, rollappId string, lastStateUpdateBySequencer bool) error {
	proposer := hook.k.GetProposer(ctx, rollappId)
	if seqAddr != proposer.Address {
		return types.ErrNotProposer
	}

	// if lastStateUpdateBySequencer is true, validate that the sequencer is in the middle of a rotation
	if lastStateUpdateBySequencer && !hook.k.AwaitingLastProposerBlock(ctx, rollappId) {
		return errorsmod.Wrap(gerrc.ErrInvalidArgument, "sequencer is not in the middle of a rotation")
	}

	return nil
}

// AfterUpdateState checks if rotation is completed and the nextProposer is changed.
// Also processes ME-specific scheduled ReplaceProposer swaps.
func (hook rollappHook) AfterUpdateState(ctx sdk.Context, stateInfo *rollapptypes.StateInfoMeta) error {
	if err := hook.k.ProcSequencerByPendingStates(ctx, stateInfo.Rollapp, stateInfo.Sequencer, &stateInfo.StateInfo); err != nil {
		return err
	}
	proposer := hook.k.GetProposer(ctx, stateInfo.Rollapp)
	return hook.k.afterStateUpdate(ctx, proposer, stateInfo.Sequencer != stateInfo.NextProposer)
}

// AfterStateFinalized cleans up a completed ReplaceProposer once the authored height is finalized.
func (hook rollappHook) AfterStateFinalized(ctx sdk.Context, rollappID string, stateInfo *rollapptypes.StateInfo) error {
	val, err := hook.k.GetReplaceProposer(ctx, rollappID)
	if err != nil {
		return err
	}
	if val == nil {
		return nil
	}
	if (stateInfo.StartHeight + stateInfo.NumBlocks - 1) < uint64(val.ReplaceProposer.BlockHeight) {
		return nil
	}
	if err = hook.k.forceRemoveUnbondedSequencer(ctx, val.ReplaceProposer.OldProposer); err != nil {
		hook.k.Logger(ctx).Error("forceRemoveUnbondedSequencer error.", "sequencer", val.ReplaceProposer.OldProposer,
			"rollapp", rollappID, "state_block_info", fmt.Sprintf("%d-%d", stateInfo.StartHeight,
				stateInfo.StartHeight+stateInfo.NumBlocks-1), "error", err.Error())
		return fmt.Errorf("forceRemoveUnbondedSequencer error in AfterStateFinalized.sequencer=%s, rollapp=%s, err=%s",
			val.ReplaceProposer.OldProposer, rollappID, err.Error())
	}
	hook.k.DeleteReplaceProposer(ctx, rollappID)
	hook.k.Logger(ctx).Info("AfterStateFinalized processed ReplaceProposer.", "rollapp", rollappID,
		"old_sequencer", val.ReplaceProposer.OldProposer, "block_height", val.ReplaceProposer.BlockHeight,
		"state_block_info", fmt.Sprintf("%d-%d", stateInfo.StartHeight, stateInfo.StartHeight+stateInfo.NumBlocks-1))
	return nil
}

// OnHardFork implements the RollappHooks interface
// unbonds all rollapp sequencers
// slashing / jailing is handled by the caller, outside of this function
func (hook rollappHook) OnHardFork(ctx sdk.Context, rollappID string, _ uint64) error {
	err := hook.k.optOutAllSequencers(ctx, rollappID)
	if err != nil {
		return errorsmod.Wrap(err, "opt out all sequencers")
	}

	// clear current proposer and successor
	hook.k.abruptRemoveProposer(ctx, rollappID)
	hook.k.SetSuccessor(ctx, rollappID, types.SentinelSeqAddr)

	return nil
}
