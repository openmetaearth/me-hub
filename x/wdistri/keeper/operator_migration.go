package keeper

import (
	"bytes"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// MigrateValidatorOperator rekeys an existing validator's distribution state without
// settling rewards or invoking validator lifecycle hooks. The caller must use a
// cached context to make this atomic with the rest of the operator migration.
// Account withdrawal preferences and consensus-address state are not modified.
func (k Keeper) MigrateValidatorOperator(ctx sdk.Context, oldAddr, newAddr sdk.ValAddress) error {
	if len(oldAddr) == 0 || len(newAddr) == 0 || oldAddr.Equals(newAddr) {
		return errors.New("distribution operator migration requires distinct, nonempty addresses")
	}

	// Each category is keyed by module prefix | length-prefixed validator address
	// | optional record suffix (delegator, period, or height/period). Only these
	// validator-scoped ranges are visited; unrelated validators are never scanned.
	categories := [][]byte{
		distrtypes.ValidatorOutstandingRewardsPrefix,
		distrtypes.ValidatorAccumulatedCommissionPrefix,
		distrtypes.ValidatorCurrentRewardsPrefix,
		distrtypes.ValidatorHistoricalRewardsPrefix,
		distrtypes.ValidatorSlashEventPrefix,
		distrtypes.DelegatorStartingInfoPrefix,
	}
	store := ctx.KVStore(k.storeKey)
	type record struct{ oldKey, newKey, value []byte }
	var records []record
	for _, category := range categories {
		oldPrefix := append(bytes.Clone(category), address.MustLengthPrefix(oldAddr)...)
		newPrefix := append(bytes.Clone(category), address.MustLengthPrefix(newAddr)...)
		if validatorDistributionPrefixExists(store, newPrefix) {
			return fmt.Errorf("distribution state already exists for destination validator %s", newAddr)
		}

		iterator := sdk.KVStorePrefixIterator(store, oldPrefix)
		for ; iterator.Valid(); iterator.Next() {
			oldKey := bytes.Clone(iterator.Key())
			newKey := append(bytes.Clone(newPrefix), oldKey[len(oldPrefix):]...)
			records = append(records, record{oldKey, newKey, bytes.Clone(iterator.Value())})
		}
		iterator.Close()
	}

	// Preserve raw values and key suffixes exactly, including empty records and
	// all period/reference metadata. All destination checks and iterator reads
	// finish before any write, so a conflict cannot leave a partial migration.
	for _, record := range records {
		store.Set(record.newKey, record.value)
		store.Delete(record.oldKey)
	}
	return nil
}

func validatorDistributionPrefixExists(store sdk.KVStore, prefix []byte) bool {
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	return iterator.Valid()
}
