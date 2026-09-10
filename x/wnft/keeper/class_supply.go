package keeper

import (
	"context"
	"encoding/binary"
)

// ClassTotalSupplyCap is kept outside the upstream nft.Class, which no longer
// contains the chain-specific supply cap field.
var ClassTotalSupplyCap = []byte{0x10}

func classTotalSupplyCapKey(classID string) []byte {
	key := make([]byte, len(ClassTotalSupplyCap)+len(classID))
	copy(key, ClassTotalSupplyCap)
	copy(key[len(ClassTotalSupplyCap):], classID)
	return key
}

func (k Keeper) SetClassTotalSupplyCap(ctx context.Context, classID string, supply uint64) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, supply)
	return store.Set(classTotalSupplyCapKey(classID), bz)
}

func (k Keeper) GetClassTotalSupplyCap(ctx context.Context, classID string) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(classTotalSupplyCapKey(classID))
	if err != nil || len(bz) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}
