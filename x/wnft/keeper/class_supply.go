package keeper

import (
	"context"
	"encoding/binary"
)

// ClassTotalSupplyCap is the KV store key prefix for storing the class total supply cap.
// Uses 0x10 to avoid collision with upstream nft module keys (0x01–0x05).
var ClassTotalSupplyCap = []byte{0x10}

func classTotalSupplyCapKey(classID string) []byte {
	key := make([]byte, len(ClassTotalSupplyCap)+len(classID))
	copy(key, ClassTotalSupplyCap)
	copy(key[len(ClassTotalSupplyCap):], classID)
	return key
}

// SetClassTotalSupplyCap stores the maximum supply cap for the given class.
func (k Keeper) SetClassTotalSupplyCap(ctx context.Context, classID string, supply uint64) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, supply)
	return store.Set(classTotalSupplyCapKey(classID), bz)
}

// GetClassTotalSupplyCap returns the maximum supply cap for the given class.
// Returns 0 if not set.
func (k Keeper) GetClassTotalSupplyCap(ctx context.Context, classID string) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(classTotalSupplyCapKey(classID))
	if err != nil || len(bz) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}
