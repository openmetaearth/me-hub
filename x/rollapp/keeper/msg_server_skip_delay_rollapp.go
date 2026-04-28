package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func (k msgServer) SkipDelayRollapp(goCtx context.Context, msg *types.MsgSkipDelayRollapp) (*types.MsgSkipDelayRollappResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsDao(ctx, msg.Creator) {
		return nil, types.ErrCheckGlobalDao
	}

	_, isFound := k.GetRollapp(ctx, msg.RollappId)
	if !isFound {
		return nil, types.ErrUnknownRollappID
	}

	k.SetSkipDelayRollapp(ctx, msg.RollappId, msg.Skip)
	return &types.MsgSkipDelayRollappResponse{}, nil
}

func (k Keeper) SetSkipDelayRollapp(ctx sdk.Context, rollappId string, skip bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SkipDelayRollappKeyPrefix))
	store.Set([]byte(rollappId), []byte(strconv.FormatBool(skip)))
}

func (k Keeper) IsSkipDelayRollapp(ctx sdk.Context, rollappId string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SkipDelayRollappKeyPrefix))
	bz := store.Get([]byte(rollappId))
	if bz == nil {
		return false
	}
	skip, err := strconv.ParseBool(string(bz))
	if err != nil {
		return false
	}
	return skip
}

func (k Keeper) GetSkipDelayRollapps(ctx sdk.Context) (rollapps []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SkipDelayRollappKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		skip, err := strconv.ParseBool(string(iterator.Value()))
		if err == nil && skip {
			rollapps = append(rollapps, string(iterator.Key()))
		}
	}
	return rollapps
}
