package keeper

import (
	"fmt"
	"strconv"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// HistoricalInfoKey is inherited from Cosmos SDK staking module
	HistoricalInfoKey = []byte{0x50} // prefix for the historical info
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper *Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper *Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate4to5 migrates x/wstaking state from consensus version 4 to 5.
// This migration only handles HistoricalInfo keys migration.
// Unlike standard Cosmos SDK staking, wstaking uses custom "Stake" (0x61) instead of "Delegation" (0x31),
// so we skip the delegation index migration.
func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	store := ctx.KVStore(m.keeper.storeKey)
	return migrateHistoricalInfoKeys(store, ctx.Logger())
}

// migrateHistoricalInfoKeys migrates HistoricalInfo keys to binary format
// This is the same logic as Cosmos SDK v5 migration, but isolated for wstaking.
func migrateHistoricalInfoKeys(store storetypes.KVStore, logger log.Logger) error {
	// old key is of format:
	// prefix (0x50) || heightBytes (string representation of height in 10 base)
	// new key is of format:
	// prefix (0x50) || heightBytes (byte array representation using big-endian byte order)
	oldStore := prefix.NewStore(store, HistoricalInfoKey)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer sdk.LogDeferred(logger, func() error { return oldStoreIter.Close() })

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		strHeight := oldStoreIter.Key()

		intHeight, err := strconv.ParseInt(string(strHeight), 10, 64)
		if err != nil {
			return fmt.Errorf("can't parse height from key %q to int64: %v", strHeight, err)
		}

		newStoreKey := GetHistoricalInfoKey(intHeight)

		// Set new key on store. Values don't change.
		store.Set(newStoreKey, oldStoreIter.Value())
		oldStore.Delete(oldStoreIter.Key())
	}

	return nil
}

// GetHistoricalInfoKey returns a key prefix for indexing HistoricalInfo objects.
// This uses the same logic as Cosmos SDK v5.
func GetHistoricalInfoKey(height int64) []byte {
	heightBytes := sdk.Uint64ToBigEndian(uint64(height))
	return append(HistoricalInfoKey, heightBytes...)
}
