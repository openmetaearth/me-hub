package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// SetMeid set a specific meid in the store from its index
func (k Keeper) SetMeid(ctx sdk.Context, meid types.Meid) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidKeyPrefix))
	b := k.cdc.MustMarshal(&meid)
	store.Set(types.MeidKey(meid.Account), b)
	storeReg := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidRegionKeyPrefix+meid.RegionId))
	storeReg.Set(types.MeidKey(meid.Account), []byte(meid.Account))
}

// GetMeid returns a meid from its index
func (k Keeper) GetMeid(ctx sdk.Context, account string) (val types.Meid, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidKeyPrefix))
	b := store.Get(types.MeidKey(account))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveMeid removes a meid from the store
func (k Keeper) RemoveMeid(ctx sdk.Context, account, regionid string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidKeyPrefix))
	store.Delete(types.MeidKey(account))
	storeReg := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidRegionKeyPrefix+regionid))
	storeReg.Delete(types.MeidKey(account))
}

// GetAllMeid returns all meid
func (k Keeper) GetAllMeid(ctx sdk.Context) (list []types.Meid) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Meid
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

func (k Keeper) GetMeidByRegion(ctx sdk.Context, regionId string) (list []types.Meid) {
	storeReg := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidRegionKeyPrefix+regionId))
	iterator := sdk.KVStorePrefixIterator(storeReg, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		account := iterator.Value()
		meid, found := k.GetMeid(ctx, string(account))
		if !found {
			panic("get meid by region fatal error")
		}
		list = append(list, meid)
	}
	return
}

func (k Keeper) GetValOwnerAddress(ctx sdk.Context, regionId string) (string, error) {
	region, ok := k.GetRegion(ctx, regionId)
	if !ok {
		return "", sdkerrors.Wrapf(types.ErrRegionNotExist, "region(%s) not found", regionId)
	}

	valAddr, err := sdk.ValAddressFromBech32(region.OperatorAddress)
	if err != nil {
		return "", sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "region bonded validator address(%s) invalid", region.OperatorAddress)
	}

	validator, ok := k.GetValidator(ctx, valAddr)
	if !ok {
		return "", sdkerrors.Wrapf(stakingtypes.ErrNoValidatorFound, "region bonded validator(%s) no found", valAddr.String())
	}
	return validator.OwnerAddress, nil
}
