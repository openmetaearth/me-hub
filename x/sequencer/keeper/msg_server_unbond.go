package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

// Unbond defines a method for removing coins from sequencer's bond
func (k msgServer) Unbond(goCtx context.Context, msg *types.MsgUnbond) (*types.MsgUnbondResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	completionTime, err := k.setSequencerToUnbonding(ctx, msg.Creator)
	if err != nil {
		return nil, err
	}

	return &types.MsgUnbondResponse{
		CompletionTime: completionTime,
	}, nil
}

func (k Keeper) setSequencerToUnbonding(ctx sdk.Context, seqAddr string) (time.Time, error) {
	seq, found := k.GetSequencer(ctx, seqAddr)
	if !found {
		return time.Time{}, types.ErrUnknownSequencer
	}

	if !seq.IsBonded() {
		return time.Time{}, errorsmod.Wrapf(
			types.ErrInvalidSequencerStatus,
			"sequencer status is not bonded: got %s",
			seq.Status.String(),
		)
	}

	if err := k.validateNoPendingReplaceProposerUnbond(ctx, seq); err != nil {
		return time.Time{}, err
	}

	oldStatus := seq.Status
	wasProposer := seq.Proposer

	// set the status to unbonding
	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	seq.Status = types.Unbonding
	seq.Proposer = false
	seq.UnbondingHeight = ctx.BlockHeight()
	seq.UnbondTime = completionTime

	k.UpdateSequencer(ctx, seq, oldStatus)
	k.setUnbondingSequencerQueue(ctx, seq)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUnbonding,
			sdk.NewAttribute(types.AttributeKeySequencer, seqAddr),
			sdk.NewAttribute(types.AttributeKeyBond, seq.Tokens.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.String()),
		),
	)

	if wasProposer {
		k.RotateProposer(ctx, seq.RollappId)
	}

	return completionTime, nil
}

func (k Keeper) validateNoPendingReplaceProposerUnbond(ctx sdk.Context, seq types.Sequencer) error {
	pending, err := k.GetReplaceProposer(ctx, seq.RollappId)
	if err != nil {
		return err
	}
	if pending == nil {
		return nil
	}

	replaceProposer := pending.ReplaceProposer
	if seq.SequencerAddress != replaceProposer.OldProposer && seq.SequencerAddress != replaceProposer.NewProposer {
		return nil
	}

	return errorsmod.Wrapf(
		types.ErrInvalidRequest,
		"sequencer %s cannot unbond while replace proposer is pending for rollapp %s",
		seq.SequencerAddress,
		seq.RollappId,
	)
}
