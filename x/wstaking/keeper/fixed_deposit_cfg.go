package keeper

import (
	"encoding/binary"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) SetFixedDepositCfg(ctx sdk.Context, cfg types.FixedDepositCfg) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositCfgKeyPrefix+cfg.RegionId))
	b := k.cdc.MustMarshal(&cfg)
	store.Set(types.FixedDepositCfgKey(cfg.Term), b)
}

func (k Keeper) GetFixedDepositCfg(ctx sdk.Context, regionId string, term int64) (val types.FixedDepositCfg, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositCfgKeyPrefix+regionId))
	b := store.Get(types.FixedDepositCfgKey(term))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveFixedDepositCfg(ctx sdk.Context, regionId string, term int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositCfgKeyPrefix+regionId))
	store.Delete(types.FixedDepositCfgKey(term))
}

func (k Keeper) GetAllFixedDepositCfg(ctx sdk.Context, regionId string) (list []types.FixedDepositCfg) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositCfgKeyPrefix+regionId))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.FixedDepositCfg
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) InitFixedDepositCountOfCfg(ctx sdk.Context, regionId string, term int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositCountOfCfgKeyPrefix+regionId))
	byteKey := types.KeyPrefix(strconv.FormatInt(term, 10))
	bz := store.Get(byteKey)
	if bz == nil {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(0))
		store.Set(byteKey, buf)
	}
	return
}

func (k Keeper) GetFixedDepositCountOfCfg(ctx sdk.Context, regionId string, term int64) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositCountOfCfgKeyPrefix+regionId))
	byteKey := types.KeyPrefix(strconv.FormatInt(term, 10))
	bz := store.Get(byteKey)

	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) IncreaseFixedDepositCountOfCfg(ctx sdk.Context, regionId string, term int64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositCountOfCfgKeyPrefix+regionId))
	byteKey := types.KeyPrefix(strconv.FormatInt(term, 10))
	bz := store.Get(byteKey)
	if bz == nil {
		return types.ErrNoFixedDepositCountOfCfgFound
	}
	count := binary.BigEndian.Uint64(bz)
	count += 1

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, count)
	store.Set(byteKey, buf)
	return nil
}

func (k Keeper) DecreaseFixedDepositCountOfCfg(ctx sdk.Context, regionId string, term int64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositCountOfCfgKeyPrefix+regionId))
	byteKey := types.KeyPrefix(strconv.FormatInt(term, 10))
	bz := store.Get(byteKey)
	if bz == nil {
		return types.ErrNoFixedDepositCountOfCfgFound
	}
	count := binary.BigEndian.Uint64(bz)
	if count == 0 {
		return types.ErrFixedDepositCountOfCfgIsZero
	}
	count -= 1

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, count)
	store.Set(byteKey, buf)
	return nil
}
