package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rollappTypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (k Keeper) SetReplaceProposer(ctx sdk.Context, data *types.MsgRepalceProposer) error {
	if nil == data {
		return fmt.Errorf("SetReplaceProposer data is nil")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	val := store.Get(types.RepalceRollappProposerKey(data.RollappId))
	if val != nil {
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
	err := k.cdc.Unmarshal(bz, &msg)
	if err != nil {
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
	bz := store.Get(types.RepalceRollappProposerKey(rollappId))
	if bz == nil {
		return false
	}
	return true
}

/*
func (k Keeper) SetReplacedSequencerAddress(ctx sdk.Context, rollappId, addr string, blockHeight int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	store.Set(types.ReplacedSequencerAddressKey(rollappId, addr), []byte(fmt.Sprintf("%d", blockHeight)))
}

func (k Keeper) GetReplacedSequencerAddress(ctx sdk.Context, rollappId, addr string) (int64, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	bz := store.Get(types.ReplacedSequencerAddressKey(rollappId, addr))
	if bz == nil {
		return 0, nil
	}

	val, err := strconv.ParseInt(string(bz), 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (k Keeper) DeleteReplacedSequencerAddress(ctx sdk.Context, rollappId, addr string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	store.Delete(types.ReplacedSequencerAddressKey(rollappId, addr))
}

func (k Keeper) IsReplacedSequencerAddress(ctx sdk.Context, rollappId, addr string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	bz := store.Get(types.ReplacedSequencerAddressKey(rollappId, addr))
	if bz == nil {
		return false
	}
	return true
}

*/

func (k Keeper) ProcSequencerByPendingStates(ctx sdk.Context, rollappId, creator string, rollappState *rollappTypes.StateInfo) error {
	val, err := k.GetReplaceProposer(ctx, rollappId)
	if err != nil {
		return err
	}
	if nil == val {
		return nil
	}
	if err = k.IsExceedAuthoredBlockHeight(ctx, rollappId, creator, rollappState.StartHeight, rollappState.NumBlocks); err != nil {
		return err
	}

	if (rollappState.StartHeight + rollappState.NumBlocks - 1) >= uint64(val.ReplaceProposer.BlockHeight) {
		//delete the replaced sequencer address record and set the new sequencer as proposer
		oldSequencer, found := k.GetSequencer(ctx, val.ReplaceProposer.OldProposer)
		if !found {
			return fmt.Errorf("can not found old sequencer: %s", val.ReplaceProposer.OldProposer)
		}
		if oldSequencer.RollappId != rollappId {
			return fmt.Errorf("old sequencer's rollapp(%s) dismatch to processing rollapp(%s)",
				oldSequencer.RollappId, rollappId)
		}
		if oldSequencer.IsProposer() || oldSequencer.Status == types.Bonded {
			oldSequencer.Proposer = false
			oldSequencer.Status = types.Unbonding
			oldSequencer.UnbondingHeight = ctx.BlockHeight()
			k.UpdateSequencer(ctx, oldSequencer, types.Bonded)
			newSequencer, found := k.GetSequencer(ctx, val.ReplaceProposer.NewProposer)
			if !found {
				return fmt.Errorf("can not found new sequencer: %s", val.ReplaceProposer.NewProposer)
			}
			if newSequencer.RollappId != rollappId {
				return fmt.Errorf("new sequencer's rollapp(%s) dismatch to processing rollapp(%s)",
					newSequencer.RollappId, rollappId)
			}
			if newSequencer.Status != types.Bonded {
				return fmt.Errorf("new sequencer %s status(%d) is not bonded", val.ReplaceProposer.NewProposer, newSequencer.Status)
			}
			newSequencer.Proposer = true
			k.UpdateSequencer(ctx, newSequencer, types.Bonded)
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
		} else {
			return nil
		}
	}
	return nil

}
func (k Keeper) IsExceedAuthoredBlockHeight(ctx sdk.Context, rollappId, creator string, startHeight uint64, numBlocks uint64) error {
	val, err := k.GetReplaceProposer(ctx, rollappId)
	if err != nil {
		return err
	}
	if nil == val {
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
	} else if val.ReplaceProposer.NewProposer == creator {
		if startHeight <= uint64(val.ReplaceProposer.BlockHeight) {
			k.Logger(ctx).Error("exceedAuthoredBlockHeight:", "new sequencer", creator,
				"authored_block_height", val.ReplaceProposer.BlockHeight+1, "request_block_height", fmt.Sprintf("%d-%d", startHeight, endHeight))
			return types.ErrorExceedAuthoredBlockHeight
		}
		return nil
	} else {
		k.Logger(ctx).Error("exceedAuthoredBlockHeight:", "unknown creator", creator,
			"old sequencer", val.ReplaceProposer.OldProposer, "new sequencer", val.ReplaceProposer.NewProposer)
		return types.ErrorExceedAuthoredBlockHeight
	}
}
