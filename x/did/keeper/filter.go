package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/did/types"
)

func (k Keeper) GetFilters(ctx sdk.Context, did, sid string) (filters [][]byte, found bool) {
	flog, found := k.GetFilterLogger(ctx, did, sid)
	if !found {
		return [][]byte{}, false
	}

	return flog.Filters, true
}

func (k Keeper) AddFilters(ctx sdk.Context, did, sid string, filters [][]byte, vc types.Credential) {
	store := ctx.KVStore(k.storeKey)

	flog, found := k.GetFilterLogger(ctx, did, sid)
	if !found {
		flog = types.FilterLogger{}
	}

	for _, filter := range filters {
		store.Set(types.GetFilterKey(sid, did, filter), k.cdc.MustMarshal(&vc))
		// record the filter to FilterLogger
		flog.Add(filter)
	}

	k.SetFilterLogger(ctx, did, sid, flog)
}

func (k Keeper) DeleteFilters(ctx sdk.Context, did, sid string, filters [][]byte) {
	store := ctx.KVStore(k.storeKey)
	flog, found := k.GetFilterLogger(ctx, did, sid)
	if !found {
		flog = types.FilterLogger{}
	}

	for _, filter := range filters {
		store.Delete(types.GetFilterKey(sid, did, filter))
		// delete the filter form FilterLogger
		flog.Delete(filter)
	}

	// delete empty FilterLogger
	if len(flog.Filters) == 0 {
		k.DeleteFilterLogger(ctx, did, sid)
	}
}

/*
filter logger
*/

func (k Keeper) GetFilterLogger(ctx sdk.Context, did, sid string) (flog types.FilterLogger, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetFilterLoggerKey(did, sid))
	if bz == nil {
		return types.FilterLogger{}, false
	}

	k.cdc.MustUnmarshal(bz, &flog)
	return flog, true
}

// SetFilterLogger set credential filter and store filter logger
func (k Keeper) SetFilterLogger(ctx sdk.Context, did, sid string, flog types.FilterLogger) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetFilterLoggerKey(did, sid), k.cdc.MustMarshal(&flog))
}

func (k Keeper) DeleteFilterLogger(ctx sdk.Context, did, sid string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetFilterLoggerKey(did, sid))
}
