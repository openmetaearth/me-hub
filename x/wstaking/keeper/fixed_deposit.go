package keeper

import (
	"encoding/binary"

	"github.com/openmetaearth/me-hub/x/wstaking/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

// GetFixedDepositCount get the total number of fixedDeposit
func (k Keeper) GetFixedDepositCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.FixedDepositCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetFixedDepositCount set the total number of fixedDeposit
func (k Keeper) SetFixedDepositCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.FixedDepositCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendFixedDeposit appends a fixedDeposit in the store with a new id and update the count
func (k Keeper) AppendFixedDeposit(
	ctx sdk.Context,
	fixedDeposit types.FixedDeposit,
) uint64 {
	// Create the fixedDeposit
	count := k.GetFixedDepositCount(ctx)

	// Set the ID of the appended value
	fixedDeposit.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositKey))
	appendedValue := k.cdc.MustMarshal(&fixedDeposit)
	store.Set(GetFixedDepositIDBytes(fixedDeposit.Id), appendedValue)

	storeAcct := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositKeyAcct+fixedDeposit.Account))
	storeAcct.Set(GetFixedDepositIDBytes(fixedDeposit.Id), GetFixedDepositIDBytes(fixedDeposit.Id))

	// Update fixedDeposit count
	k.SetFixedDepositCount(ctx, count+1)

	return count
}

// SetFixedDeposit set a specific fixedDeposit in the store
func (k Keeper) SetFixedDeposit(ctx sdk.Context, fixedDeposit types.FixedDeposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositKey))
	b := k.cdc.MustMarshal(&fixedDeposit)
	store.Set(GetFixedDepositIDBytes(fixedDeposit.Id), b)

	storeAcct := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositKeyAcct+fixedDeposit.Account))
	storeAcct.Set(GetFixedDepositIDBytes(fixedDeposit.Id), GetFixedDepositIDBytes(fixedDeposit.Id))
}

// GetFixedDeposit returns a fixedDeposit from its id
func (k Keeper) GetFixedDeposit(ctx sdk.Context, id uint64) (val types.FixedDeposit, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositKey))
	b := store.Get(GetFixedDepositIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetFixedDepositByAcct returns the list of fixedDeposits of an account
func (k Keeper) GetFixedDepositByAcct(ctx sdk.Context, acct string) ([]types.FixedDeposit, error) {
	var list []types.FixedDeposit
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositKeyAcct+acct))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := iterator.Value()
		fixedDeposit, found := k.GetFixedDeposit(ctx, GetFixedDepositIDFromBytes(id))
		if !found {
			return nil, types.ErrNoFixedDepositFound.Wrapf("index inconsistency: fixed deposit id %d not found in store", GetFixedDepositIDFromBytes(id))
		}
		list = append(list, fixedDeposit)
	}

	return list, nil
}

// RemoveFixedDeposit removes a fixedDeposit from the store
func (k Keeper) RemoveFixedDeposit(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositKey))
	var fixedDeposit types.FixedDeposit
	b := store.Get(GetFixedDepositIDBytes(id))
	if b == nil {
		return
	}
	k.cdc.MustUnmarshal(b, &fixedDeposit)
	store.Delete(GetFixedDepositIDBytes(id))

	storeAcct := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositKeyAcct+fixedDeposit.Account))
	storeAcct.Delete(GetFixedDepositIDBytes(id))
}

// GetAllFixedDeposit returns all fixedDeposit
func (k Keeper) GetAllFixedDepositWithPage(ctx sdk.Context, req *types.QueryAllFixedDepositRequest) (list []types.FixedDeposit, pageRes *query.PageResponse, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositKey))
	pageRes, err = query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var vc types.FixedDeposit
		if err := k.cdc.Unmarshal(value, &vc); err != nil {
			return err // todo: warp error
		}
		list = append(list, vc)
		return nil
	})
	if err != nil {
		return []types.FixedDeposit{}, &query.PageResponse{}, err
	}
	return
}

// GetAllFixedDeposit returns all fixedDeposit
func (k Keeper) GetAllFixedDeposit(ctx sdk.Context) (list []types.FixedDeposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.FixedDeposit
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

// GetFixedDepositIDBytes returns the byte representation of the ID
func GetFixedDepositIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetFixedDepositIDFromBytes returns ID in uint64 format from a byte array
func GetFixedDepositIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetFixedDepositTotalAmount(ctx sdk.Context, fixedDepositTotal types.FixedDepositTotal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	b := k.cdc.MustMarshal(&fixedDepositTotal)
	store.Set([]byte(types.FixedDepositTotalAmountKey), b)
}

func (k Keeper) GetFixedDepositTotalAmount(
	ctx sdk.Context,
) (val types.FixedDepositTotal, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})

	b := store.Get([]byte(types.FixedDepositTotalAmountKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
