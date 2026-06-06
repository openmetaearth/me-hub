package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/x/did/types"
)

func (k Keeper) HasCredential(ctx sdk.Context, did, sid string) (found bool) {
	if _, found := k.GetCredential(ctx, did, sid); found {
		return true
	}

	return false
}

func (k Keeper) GetCredential(ctx sdk.Context, did, sid string) (vc types.Credential, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCredentialKey(did, sid))
	if bz == nil {
		return types.Credential{}, false
	}

	k.cdc.MustUnmarshal(bz, &vc)
	return vc, true
}

func (k Keeper) GetCredentials(ctx sdk.Context) (vcs []types.Credential) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CredentialPrefix)
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		var vc types.Credential
		k.cdc.MustUnmarshal(iterator.Value(), &vc)
		vcs = append(vcs, vc)
	}

	return vcs
}

func (k Keeper) GetCredentialsByDid(ctx sdk.Context, did string) (vcs []types.Credential) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.GetCredentialPrefixByDid(did))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var vc types.Credential
		k.cdc.MustUnmarshal(iterator.Value(), &vc)
		vcs = append(vcs, vc)
	}

	return vcs
}

func (k Keeper) credentialMatchesFilter(ctx sdk.Context, sid string, filter []byte, vc types.Credential) bool {
	if vc.Sid != sid {
		return false
	}

	flog, found := k.GetFilterLogger(ctx, vc.Did, sid)
	return found && flog.Contains(filter)
}

func (k Keeper) getCredentialsByExactFilterPrefix(
	ctx sdk.Context,
	keyPrefix []byte,
	pageReq *query.PageRequest,
) (vcs []types.Credential, pageRes *query.PageResponse, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), keyPrefix)

	pageRes, err = query.Paginate(store, pageReq, func(key []byte, value []byte) error {
		var vc types.Credential
		if err := k.cdc.Unmarshal(value, &vc); err != nil {
			return err
		}
		vcs = append(vcs, vc)
		return nil
	})
	if err != nil {
		return []types.Credential{}, &query.PageResponse{}, err
	}

	return vcs, pageRes, nil
}

func (k Keeper) getLegacyCredentialsByFilter(
	ctx sdk.Context,
	sid string,
	filter []byte,
	pageReq *query.PageRequest,
) (vcs []types.Credential, pageRes *query.PageResponse, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetLegacyFilterPrefixBySidAndFilter(sid, filter))

	pageRes, err = query.FilteredPaginate(store, pageReq, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var vc types.Credential
		if err := k.cdc.Unmarshal(value, &vc); err != nil {
			return false, err
		}

		if !k.credentialMatchesFilter(ctx, sid, filter, vc) {
			return false, nil
		}

		if accumulate {
			vcs = append(vcs, vc)
		}

		return true, nil
	})
	if err != nil {
		return []types.Credential{}, &query.PageResponse{}, err
	}

	return vcs, pageRes, nil
}

func (k Keeper) GetCredentialsByFilter(
	ctx sdk.Context,
	sid string,
	filter []byte,
	pageReq *query.PageRequest,
) (vcs []types.Credential, pageRes *query.PageResponse, err error) {
	vcs, pageRes, err = k.getCredentialsByExactFilterPrefix(ctx, types.GetFilterPrefixBySidAndFilter(sid, filter), pageReq)
	if err != nil {
		return []types.Credential{}, &query.PageResponse{}, err
	}
	if len(vcs) > 0 {
		return vcs, pageRes, nil
	}

	return k.getLegacyCredentialsByFilter(ctx, sid, filter, pageReq)
}

func (k Keeper) SetCredential(ctx sdk.Context, did, sid string, credential types.Credential) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetCredentialKey(did, sid), k.cdc.MustMarshal(&credential))
}

func (k Keeper) DeleteCredential(ctx sdk.Context, did, sid string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetCredentialKey(did, sid))
}

func (k Keeper) IteratorCredentialsByFilter(ctx sdk.Context, sid string, filter []byte, cb func(delegation types.Credential) (stop bool)) {
	seen := map[string]struct{}{}
	stop := k.iterateCredentialsByFilterPrefix(ctx, types.GetFilterPrefixBySidAndFilter(sid, filter), nil, seen, cb)
	if stop {
		return
	}

	k.iterateCredentialsByFilterPrefix(ctx, types.GetLegacyFilterPrefixBySidAndFilter(sid, filter), func(vc types.Credential) bool {
		return k.credentialMatchesFilter(ctx, sid, filter, vc)
	}, seen, cb)
}

func (k Keeper) iterateCredentialsByFilterPrefix(
	ctx sdk.Context,
	keyPrefix []byte,
	match func(types.Credential) bool,
	seen map[string]struct{},
	cb func(delegation types.Credential) (stop bool),
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), keyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		var vc types.Credential
		k.cdc.MustUnmarshal(iterator.Value(), &vc)
		if match != nil && !match(vc) {
			continue
		}

		credentialID := vc.Did + "\x00" + vc.Sid
		if _, ok := seen[credentialID]; ok {
			continue
		}
		seen[credentialID] = struct{}{}

		if cb(vc) {
			return true
		}
	}

	return false
}
