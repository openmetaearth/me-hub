package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

var _ rollapptypes.RollappHooks = rollappHook{}

// Hooks wrapper struct for rollapp keeper.
type rollappHook struct {
	k Keeper
}

// RollappHooks returns the wrapper struct.
func (k Keeper) RollappHooks() rollapptypes.RollappHooks {
	return rollappHook{
		k,
	}
}

func (hook rollappHook) BeforeUpdateState(ctx sdk.Context, seqAddr string, rollappId string) error {
	// check to see if the sequencer has been registered before
	sequencer, found := hook.k.GetSequencer(ctx, seqAddr)
	if !found {
		return types.ErrUnknownSequencer
	}

	// check to see if the rollappId matches the one of the sequencer
	if sequencer.RollappId != rollappId {
		return types.ErrSequencerRollappMismatch
	}

	// check to see if the sequencer is active and can make the update
	if sequencer.Status != types.Bonded {
		return types.ErrInvalidSequencerStatus
	}

	if !sequencer.Proposer {
		return types.ErrNotActiveSequencer
	}
	return nil
}

func (hook rollappHook) AfterStateFinalized(ctx sdk.Context, rollappID string, stateInfo *rollapptypes.StateInfo) error {
	val, err := hook.k.GetReplaceProposer(ctx, rollappID)
	if err != nil {
		return err
	}
	if val != nil {
		if (stateInfo.StartHeight + stateInfo.NumBlocks - 1) >= uint64(val.ReplaceProposer.BlockHeight) {
			err = hook.k.forceRemoveUnbondingSequencer(ctx, val.ReplaceProposer.OldProposer, stateInfo.StartHeight, stateInfo.NumBlocks)
			if err != nil {
				hook.k.Logger(ctx).Error("forceRemoveUnbondingSequencer error.", "sequencer", val.ReplaceProposer.OldProposer,
					"rollapp", rollappID, "state_block_info", fmt.Sprintf("%d-%d", stateInfo.StartHeight,
						stateInfo.StartHeight+stateInfo.NumBlocks-1), "error", err.Error())
				return fmt.Errorf("forceRemoveUnbondingSequencer error in AfterStateFinalized.sequencer=%s,"+
					" rollapp = %s, err = %s", val.ReplaceProposer.OldProposer, rollappID, err.Error())
			}
			hook.k.DeleteReplaceProposer(ctx, rollappID)
			hook.k.Logger(ctx).Info("AfterStateFinalized processed ReplaceProposer.", "rollapp", rollappID,
				"old_sequencer", val.ReplaceProposer.OldProposer, "block_height", val.ReplaceProposer.BlockHeight,
				"state_block_info", fmt.Sprintf("%d-%d", stateInfo.StartHeight, stateInfo.StartHeight+stateInfo.NumBlocks-1))
		}
	}
	return nil
}

// FraudSubmitted implements the RollappHooks interface
// It slashes the sequencer and unbonds all other bonded sequencers
func (hook rollappHook) FraudSubmitted(ctx sdk.Context, rollappID string, height uint64, seqAddr string) error {
	err := hook.k.Slashing(ctx, seqAddr)
	if err != nil {
		return err
	}

	// unbond all other bonded sequencers
	sequencers := hook.k.GetSequencersByRollappByStatus(ctx, rollappID, types.Bonded)
	for _, sequencer := range sequencers {
		err := hook.k.forceUnbondSequencer(ctx, sequencer.SequencerAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

// RollappCreated implements types.RollappHooks.
func (hook rollappHook) RollappCreated(ctx sdk.Context, rollappID string) error {
	return nil
}

func (hook rollappHook) ProcPendingStates(ctx sdk.Context, rollappID, creator string, stateInfo *rollapptypes.StateInfo) error {
	return hook.k.ProcSequencerByPendingStates(ctx, rollappID, creator, stateInfo)
}
