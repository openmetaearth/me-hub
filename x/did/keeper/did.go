package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/did/types"
)

func (k Keeper) HasDID(ctx sdk.Context, addr sdk.AccAddress) bool {
	if _, found := k.GetDID(ctx, addr); found {
		return true
	}
	return false
}

func (k Keeper) GetDID(ctx sdk.Context, addr sdk.AccAddress) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDIDKey(addr))
	if bz == nil {
		return "", false
	}

	return string(bz), true
}

func (k Keeper) SetDID(ctx sdk.Context, addr sdk.AccAddress, did string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDIDKey(addr), []byte(did))
}

func (k Keeper) DeleteDID(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDIDKey(addr))
}

func (k Keeper) HasDidInfo(ctx sdk.Context, did string) bool {
	if _, found := k.GetDidInfo(ctx, did); found {
		return true
	}
	return false
}

func (k Keeper) GetDidInfo(ctx sdk.Context, did string) (info types.DidInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDidInfoKey(did))
	if bz == nil {
		return types.DidInfo{}, false
	}

	k.cdc.MustUnmarshal(bz, &info)
	return info, true
}

func (k Keeper) GetDidInfos(ctx sdk.Context) (infos []types.DidInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DidInfoPrefix)
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		var info types.DidInfo
		k.cdc.MustUnmarshal(iterator.Value(), &info)
		infos = append(infos, info)
	}

	return infos
}

func (k Keeper) SetDidInfo(ctx sdk.Context, did string, info types.DidInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDidInfoKey(did), k.cdc.MustMarshal(&info))
}

func (k Keeper) DeleteDidInfo(ctx sdk.Context, did string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDidInfoKey(did))
}

func (k Keeper) GetDidDocument(ctx sdk.Context, did string) (doc types.DidDocument, found bool) {
	baseInfo, found := k.GetDidInfo(ctx, did)
	if !found {
		return types.DidDocument{}, false
	}

	vcs := k.GetCredentialsByDid(ctx, did)
	return types.DidDocument{Info: baseInfo, Vcs: vcs}, true
}

// SetDidDocument is only used at genesis
func (k Keeper) SetDidDocument(ctx sdk.Context, did string, doc types.DidDocument) {
	addr := k.MustAccAddressFromPubkeyString(doc.Info.Pubkey)
	k.SetDID(ctx, addr, did)
	k.SetDidInfo(ctx, did, doc.Info)

	// store vcs
	for _, vc := range doc.Vcs {
		//if vc.Holder != did {
		//	panic("certificate does not belong to DID")
		//}
		k.SetCredential(ctx, did, vc.Sid, vc)
	}
}

func (k Keeper) DeleteDidDocument(ctx sdk.Context, did string) {
	store := ctx.KVStore(k.storeKey)

	if info, found := k.GetDidInfo(ctx, did); found {
		addr := k.MustAccAddressFromPubkeyString(info.Pubkey)
		k.DeleteDID(ctx, addr)
	}
	k.DeleteDidInfo(ctx, did)

	// delete all credentials
	iter := sdk.KVStorePrefixIterator(store, types.GetCredentialPrefixByDid(did))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
