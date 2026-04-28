package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) SetRecord(ctx sdk.Context, record types.Record, acc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	userKey := types.GetRecordKey(acc)
	userKey = append(userKey, []byte(record.RecordNumber)...)
	bz := types.MustMarshalRecord(k.cdc, record)
	store.Set(userKey, bz)
}

func (k Keeper) GetRecordsByAddress(ctx sdk.Context, from sdk.AccAddress) []types.Record {
	store := ctx.KVStore(k.storeKey)

	prefix := types.GetRecordKey(from)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var records []types.Record
	for ; iterator.Valid(); iterator.Next() {
		recordData := iterator.Value()
		record := types.MustUnmarshalRecord(k.cdc, recordData)
		records = append(records, record)
	}
	return records
}
func (k Keeper) GetAllRecords(ctx sdk.Context) []types.Record {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.NewRecordKey)
	defer iterator.Close()

	var records []types.Record
	for ; iterator.Valid(); iterator.Next() {
		recordData := iterator.Value()
		record := types.MustUnmarshalRecord(k.cdc, recordData)
		records = append(records, record)
	}
	return records
}

func (k Keeper) SetReviewRecord(ctx sdk.Context, rr types.ReviewRecord) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalReviewRecord(k.cdc, rr)
	store.Set(types.GetReviewRecordKey(rr.ActionNumber), bz)
}
func (k Keeper) GetReviewRecordByID(ctx sdk.Context, recordNumber string) types.ReviewRecord {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetReviewRecordKey(recordNumber))
	return types.MustUnmarshalReviewRecord(k.cdc, b)
}
