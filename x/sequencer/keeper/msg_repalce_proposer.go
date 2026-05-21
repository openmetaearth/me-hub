package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) ReplaceProposer(goCtx context.Context, msg *types.MsgReplaceProposerRequest) (*types.MsgReplaceProposerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check to see if the rollapp has been registered before
	rollapp, found := k.rollappKeeper.GetRollapp(ctx, msg.ReplaceProposer.RollappId)
	if !found {
		return nil, types.ErrUnknownRollappID
	}
	if rollapp.Frozen {
		return nil, types.ErrRollappJailed
	}
	if msg.Creator != rollapp.Creator {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "only rollapp creator %s can replace proposer, but got %s", rollapp.Creator, msg.Creator)

	}

	if found := k.IsHasReplaceProposer(ctx, msg.ReplaceProposer.RollappId); found {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "there is already a pending replace proposer request")
	}

	oldSequencer, found := k.GetSequencer(ctx, msg.ReplaceProposer.OldProposer)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrUnknownSequencer, "old proposer %s not found", msg.ReplaceProposer.OldProposer)
	}
	if !oldSequencer.Proposer {
		return nil, errorsmod.Wrapf(types.ErrInvalidSequencerStatus, "old proposer %s is not a proposer", msg.ReplaceProposer.OldProposer)
	}
	newSequencer, found := k.GetSequencer(ctx, msg.ReplaceProposer.NewProposer)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrUnknownSequencer, "new proposer %s not found", msg.ReplaceProposer.NewProposer)
	}
	if !newSequencer.IsBonded() {
		return nil, errorsmod.Wrapf(types.ErrInvalidSequencerStatus, "new proposer %s is not bonded", msg.ReplaceProposer.NewProposer)
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
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "replace proposer block height %d must be greater than last state info end height %d", msg.ReplaceProposer.BlockHeight, stateInfo.StartHeight+stateInfo.NumBlocks)
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
