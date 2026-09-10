package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/utils/gerrc"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (k msgServer) ReplaceProposer(goCtx context.Context, msg *types.MsgReplaceProposerRequest) (*types.MsgReplaceProposerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.ReplaceProposer == nil {
		return nil, errorsmod.Wrap(gerrc.ErrInvalidArgument, "replace proposer is nil")
	}

	rollapp, found := k.rollappKeeper.GetRollapp(ctx, msg.ReplaceProposer.RollappId)
	if !found {
		return nil, types.ErrUnknownRollappID
	}
	if msg.Creator != rollapp.Owner {
		return nil, errorsmod.Wrapf(gerrc.ErrPermissionDenied, "only rollapp owner %s can replace proposer, but got %s", rollapp.Owner, msg.Creator)
	}

	if k.IsHasReplaceProposer(ctx, msg.ReplaceProposer.RollappId) {
		return nil, errorsmod.Wrap(gerrc.ErrInvalidArgument, "there is already a pending replace proposer request")
	}

	oldSequencer, err := k.RealSequencer(ctx, msg.ReplaceProposer.OldProposer)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrSequencerNotFound, "old proposer %s not found", msg.ReplaceProposer.OldProposer)
	}
	if oldSequencer.RollappId != msg.ReplaceProposer.RollappId {
		return nil, errorsmod.Wrapf(gerrc.ErrInvalidArgument, "old proposer %s belongs to different rollapp", msg.ReplaceProposer.OldProposer)
	}
	proposer := k.GetProposer(ctx, msg.ReplaceProposer.RollappId)
	if proposer.Address != msg.ReplaceProposer.OldProposer {
		return nil, errorsmod.Wrapf(types.ErrNotProposer, "old proposer %s is not a proposer", msg.ReplaceProposer.OldProposer)
	}

	newSequencer, err := k.RealSequencer(ctx, msg.ReplaceProposer.NewProposer)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrSequencerNotFound, "new proposer %s not found", msg.ReplaceProposer.NewProposer)
	}
	if newSequencer.RollappId != msg.ReplaceProposer.RollappId {
		return nil, errorsmod.Wrapf(gerrc.ErrInvalidArgument, "new proposer %s belongs to different rollapp", msg.ReplaceProposer.NewProposer)
	}
	if !newSequencer.Bonded() {
		return nil, errorsmod.Wrapf(gerrc.ErrFailedPrecondition, "new proposer %s is not bonded", msg.ReplaceProposer.NewProposer)
	}

	stateInfoIndex, found := k.rollappKeeper.GetLatestStateInfoIndex(ctx, msg.ReplaceProposer.RollappId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrUnknownRequest, "no state info index found for rollapp %s", msg.ReplaceProposer.RollappId)
	}
	stateInfo, found := k.rollappKeeper.GetStateInfo(ctx, msg.ReplaceProposer.RollappId, stateInfoIndex.Index)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrUnknownRequest, "no state info found for rollapp %s at index %d", msg.ReplaceProposer.RollappId, stateInfoIndex.Index)
	}
	if msg.ReplaceProposer.BlockHeight <= int64(stateInfo.StartHeight+stateInfo.NumBlocks) {
		return nil, errorsmod.Wrapf(gerrc.ErrInvalidArgument, "replace proposer block height %d must be greater than last state info end height %d", msg.ReplaceProposer.BlockHeight, stateInfo.StartHeight+stateInfo.NumBlocks)
	}

	if err := k.SetReplaceProposer(ctx, msg.ReplaceProposer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventReplaceProposer,
			sdk.NewAttribute(types.AttributeKeyRollappId, msg.ReplaceProposer.RollappId),
			sdk.NewAttribute(types.AttributeKeyOldProposer, msg.ReplaceProposer.OldProposer),
			sdk.NewAttribute(types.AttributeKeyNewProposer, msg.ReplaceProposer.NewProposer),
		),
	)
	return &types.MsgReplaceProposerResponse{}, nil
}
