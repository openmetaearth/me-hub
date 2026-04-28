package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// SetMeidNFT set a specific meid in the store from its index
func (k Keeper) SetMeidNFT(ctx sdk.Context, meidNFT types.MeidNFT) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidNFTKeyPrefix))
	b := k.cdc.MustMarshal(&meidNFT)
	store.Set(types.MeidNFTKey(meidNFT.Account), b)

	storeReg := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidNFTAccountKeyPrefix+meidNFT.RegionId))
	storeReg.Set(types.MeidNFTKey(meidNFT.Account), []byte(meidNFT.Account))
}

// GetMeidNFT returns a meidNFT from its index
func (k Keeper) GetMeidNFT(ctx sdk.Context, meidNFT string) (val types.MeidNFT, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidNFTKeyPrefix))
	b := store.Get(types.MeidNFTKey(meidNFT))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveMeidNFT removes a meidNFT from the store
func (k Keeper) RemoveMeidNFT(ctx sdk.Context, account, regionId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidNFTKeyPrefix))
	store.Delete(types.MeidNFTKey(account))

	storeReg := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidNFTAccountKeyPrefix+account))
	storeReg.Delete(types.MeidNFTKey(account))
}

func (k Keeper) GetMeidNFTByAccount(ctx sdk.Context, account string) (val types.MeidNFT, found bool) {
	storeReg := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MeidNFTAccountKeyPrefix+account))
	iterator := sdk.KVStorePrefixIterator(storeReg, []byte{})
	defer iterator.Close()
	return k.GetMeidNFT(ctx, account)
}
