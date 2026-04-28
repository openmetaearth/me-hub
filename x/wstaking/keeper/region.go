package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// SetRegion set a specific region in the store from its index
func (k Keeper) SetRegion(ctx sdk.Context, region types.Region) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionKeyPrefix))
	b := k.cdc.MustMarshal(&region)
	store.Set(types.RegionKey(region.RegionId), b)
}

// GetRegion returns a region from its index
func (k Keeper) GetRegion(ctx sdk.Context, regionId string) (region types.Region, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionKeyPrefix))
	b := store.Get(types.RegionKey(regionId))
	if b == nil {
		return region, false
	}
	k.cdc.MustUnmarshal(b, &region)
	return region, true
}

// RemoveRegion removes a region from the store
func (k Keeper) RemoveRegion(ctx sdk.Context, regionId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionKeyPrefix))
	store.Delete(types.RegionKey(regionId))
}

// GetAllRegion returns all region
func (k Keeper) GetAllRegion(ctx sdk.Context) (list []types.Region) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Region
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

func (k Keeper) CreateRegionAccount(ctx sdk.Context, accountType types.REGION_ACCOUNT_TYPE, regionId string) sdk.AccAddress {
	regionAcc := k.GetRegionAccount(ctx, accountType, regionId)
	if regionAcc == nil {
		vaultAddr := types.GetRegionAccountAddr(accountType, regionId)
		k.authKeeper.SetAccount(ctx, k.authKeeper.NewAccountWithAddress(ctx, vaultAddr))
		return vaultAddr
	}
	return regionAcc.GetAddress()
}

func (k Keeper) GetRegionAccount(ctx sdk.Context, accountType types.REGION_ACCOUNT_TYPE, regionId string) authtypes.AccountI {
	vaultAddr := types.GetRegionAccountAddr(accountType, regionId)
	return k.authKeeper.GetAccount(ctx, vaultAddr)
}

// GetAllRegion returns all region
func (k Keeper) GetAllRegionI(ctx sdk.Context) (list []types.RegionI) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RegionKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Region
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		val.GetCreator()
		list = append(list, &val)
	}
	return
}

func (k Keeper) BondRegion(ctx sdk.Context, validator stakingtypes.Validator, tokens sdk.Int, changeOperator bool) {
	region, found := k.GetRegion(ctx, validator.Description.RegionID)
	if !found {
		return
	}
	if changeOperator {
		region.OperatorAddress = validator.OperatorAddress
		k.groupKeeper.UpdateGroupAdmin(ctx, validator.Description.RegionID, validator.OwnerAddress)
	}
	region.RegionShare = tokens
	k.SetRegion(ctx, region)
}

func (k Keeper) UnBondRegion(ctx sdk.Context, regionId string) {
	region, found := k.GetRegion(ctx, regionId)
	if !found {
		return
	}
	region.RegionShare = sdk.ZeroInt()
	region.OperatorAddress = ""
	k.SetRegion(ctx, region)
	k.groupKeeper.UpdateGroupAdmin(ctx, regionId, "")
}
