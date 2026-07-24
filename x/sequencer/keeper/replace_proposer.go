package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rollappTypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (k Keeper) SetReplaceProposer(ctx sdk.Context, data *types.MsgRepalceProposer) error {
	if data == nil {
		return fmt.Errorf("SetReplaceProposer data is nil")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	if store.Get(types.RepalceRollappProposerKey(data.RollappId)) != nil {
		return types.ErrExistingReplaceProposer
	}
	storeReplaceProposerInfo := &types.MsgStoreReplaceProposer{
		ReplaceProposer: *data,
		HubBlockHeight:  ctx.BlockHeight(),
	}
	bz, err := k.cdc.Marshal(storeReplaceProposerInfo)
	if err != nil {
		return err
	}
	store.Set(types.RepalceRollappProposerKey(data.RollappId), bz)
	return nil
}

func (k Keeper) GetReplaceProposer(ctx sdk.Context, rollappId string) (*types.MsgStoreReplaceProposer, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	bz := store.Get(types.RepalceRollappProposerKey(rollappId))
	if bz == nil {
		return nil, nil
	}
	var msg types.MsgStoreReplaceProposer
	if err := k.cdc.Unmarshal(bz, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (k Keeper) DeleteReplaceProposer(ctx sdk.Context, rollappId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	store.Delete(types.RepalceRollappProposerKey(rollappId))
}

func (k Keeper) IsHasReplaceProposer(ctx sdk.Context, rollappId string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	return store.Get(types.RepalceRollappProposerKey(rollappId)) != nil
}

// ProcSequencerByPendingStates executes a scheduled ReplaceProposer once the
// pending state crosses the authored block height. Adapted for v4 proposer store.
func (k Keeper) ProcSequencerByPendingStates(ctx sdk.Context, rollappId, creator string, rollappState *rollappTypes.StateInfo) error {
	val, err := k.GetReplaceProposer(ctx, rollappId)
	if err != nil {
		return err
	}
	if val == nil {
		return nil
	}
	if err = k.IsExceedAuthoredBlockHeight(ctx, rollappId, creator, rollappState.StartHeight, rollappState.NumBlocks); err != nil {
		return err
	}

	if (rollappState.StartHeight + rollappState.NumBlocks - 1) < uint64(val.ReplaceProposer.BlockHeight) {
		return nil
	}

	oldSequencer, err := k.RealSequencer(ctx, val.ReplaceProposer.OldProposer)
	if err != nil {
		return fmt.Errorf("can not found old sequencer: %s", val.ReplaceProposer.OldProposer)
	}
	if oldSequencer.RollappId != rollappId {
		return fmt.Errorf("old sequencer's rollapp(%s) dismatch to processing rollapp(%s)",
			oldSequencer.RollappId, rollappId)
	}

	proposer := k.GetProposer(ctx, rollappId)
	if proposer.Address != oldSequencer.Address && !oldSequencer.Bonded() {
		return nil
	}

	newSequencer, err := k.RealSequencer(ctx, val.ReplaceProposer.NewProposer)
	if err != nil {
		return fmt.Errorf("can not found new sequencer: %s", val.ReplaceProposer.NewProposer)
	}
	if newSequencer.RollappId != rollappId {
		return fmt.Errorf("new sequencer's rollapp(%s) dismatch to processing rollapp(%s)",
			newSequencer.RollappId, rollappId)
	}
	if !newSequencer.Bonded() {
		return fmt.Errorf("new sequencer %s is not bonded", val.ReplaceProposer.NewProposer)
	}

	// Swap proposer via v4 proposer store, then unbond the old proposer.
	k.SetProposer(ctx, rollappId, newSequencer.Address)
	k.unbond(ctx, &oldSequencer)
	k.SetSequencer(ctx, oldSequencer)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventProcReplaceProposer,
			sdk.NewAttribute(types.AttributeKeyRollappId, rollappId),
			sdk.NewAttribute(types.AttributeKeyOldProposer, val.ReplaceProposer.OldProposer),
			sdk.NewAttribute(types.AttributeKeyNewProposer, val.ReplaceProposer.NewProposer),
			sdk.NewAttribute(types.AttributeKeyPendingBlockHeight, fmt.Sprintf("%d-%d", rollappState.StartHeight,
				rollappState.StartHeight+rollappState.NumBlocks-1)),
		),
	)
	return nil
}

func (k Keeper) IsExceedAuthoredBlockHeight(ctx sdk.Context, rollappId, creator string, startHeight uint64, numBlocks uint64) error {
	val, err := k.GetReplaceProposer(ctx, rollappId)
	if err != nil {
		return err
	}
	if val == nil {
		return nil
	}
	endHeight := startHeight + numBlocks - 1
	if val.ReplaceProposer.OldProposer == creator {
		if endHeight > uint64(val.ReplaceProposer.BlockHeight) {
			k.Logger(ctx).Error("exceedAuthoredBlockHeight:", "old sequencer", creator,
				"authored_block_height", val.ReplaceProposer.BlockHeight, "request_block_height", fmt.Sprintf("%d-%d", startHeight, endHeight))
			return types.ErrorExceedAuthoredBlockHeight
		}
		return nil
	}
	if val.ReplaceProposer.NewProposer == creator {
		if startHeight <= uint64(val.ReplaceProposer.BlockHeight) {
			k.Logger(ctx).Error("exceedAuthoredBlockHeight:", "new sequencer", creator,
				"authored_block_height", val.ReplaceProposer.BlockHeight+1, "request_block_height", fmt.Sprintf("%d-%d", startHeight, endHeight))
			return types.ErrorExceedAuthoredBlockHeight
		}
		return nil
	}
	k.Logger(ctx).Error("exceedAuthoredBlockHeight:", "unknown creator", creator,
		"old sequencer", val.ReplaceProposer.OldProposer, "new sequencer", val.ReplaceProposer.NewProposer)
	return types.ErrorExceedAuthoredBlockHeight
}

// forceRemoveUnbondedSequencer refunds and clears an already-unbonded old proposer
// after the replace height is finalized.
func (k Keeper) forceRemoveUnbondedSequencer(ctx sdk.Context, seqAddr string) error {
	seq, err := k.RealSequencer(ctx, seqAddr)
	if err != nil {
		return err
	}
	if seq.Bonded() {
		// Still bonded somehow — unbond first.
		k.unbond(ctx, &seq)
	}
	if !seq.Tokens.IsZero() {
		if err := k.refund(ctx, &seq, seq.TokensCoin()); err != nil {
			return err
		}
	}
	k.SetSequencer(ctx, seq)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventDirectRemoveSequencer,
			sdk.NewAttribute(types.AttributeKeySequencer, seqAddr),
		),
	)
	return nil
}
