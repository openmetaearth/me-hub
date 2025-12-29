package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/sequencer/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) RepalceProposer(goCtx context.Context, msg *types.MsgRepalceProposerRequest) (*types.MsgRepalceProposerResponse, error){
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
	oldSequencer,err := k.GetSequencer(ctx,msg.ReplaceProposer.OldProposer)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrUnknownSequencer, "old proposer %s not found", msg.ReplaceProposer.OldProposer)
	}
	if !oldSequencer.Proposer {
		return nil, errorsmod.Wrapf(types.ErrInvalidSequencerStatus, "old proposer %s is not a proposer", msg.ReplaceProposer.OldProposer)
	}
	newSequencer,err := k.GetSequencer(ctx,msg.ReplaceProposer.NewProposer)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrUnknownSequencer, "new proposer %s not found", msg.ReplaceProposer.NewProposer)
	}	
	if !newSequencer.IsBonded() {
		return nil, errorsmod.Wrapf(types.ErrInvalidSequencerStatus, "new proposer %s is not bonded", msg.ReplaceProposer.NewProposer)
	}
	stateInfoIndex,found := k.rollappKeeper.GetLatestStateInfoIndex(ctx,msg.ReplaceProposer.RollappId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrUnknownRequest, "no state info index found for rollapp %s", msg.ReplaceProposer.RollappId)
	}
	stateInfo,found := k.rollappKeeper.GetStateInfo(ctx,msg.ReplaceProposer.RollappId,stateInfoIndex.Index)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrUnknownRequest, "no state info found for rollapp %s at index %d", msg.ReplaceProposer.RollappId,stateInfoIndex.Index)
	}
	if msg.ReplaceProposer.BlockHeight <= (stateInfo.StartHeight + stateInfo.NumBlocks){
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "replace proposer block height %d must be greater than last state info end height %d", msg.ReplaceProposer.BlockHeight, stateInfo.StartHeight + stateInfo.NumBlocks)
	}
	

	

	// check if the msg.Creator is the current proposer
	if rollapp. != msg.OldProposer {
		return nil, errorsmod.Wrapf(types.ErrUnauthorized, "only current proposer %s can replace proposer, but got %s", rollapp.Proposer, msg.Creator)
	}

}