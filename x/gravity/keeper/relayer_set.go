package keeper

import (
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// --- RELAYER SET REQUESTS --- //

// GetCurrentRelayerSet gets powers from the store and normalizes them
// into an integer percentage with a resolution of uint32 Max meaning
// a given validators 'Relayer power' is computed as
// Cosmos power / total cosmos power = x / uint32 Max
// where x is the voting power on the Relayer contract. This allows us
// to only use integer division which produces a known rounding error
// from truncation equal to the ratio of the validators
// Cosmos power / total cosmos power ratio, leaving us at uint32 Max - 1
// total voting power. This is an acceptable rounding error since floating
// point may cause consensus problems if different floating point unit
// implementations are involved.
func (k Keeper) GetCurrentRelayerSet(ctx sdk.Context) *types.RelayerSet {
	allRelayers := k.GetAllRelayers(ctx, true)
	bridgeValidators := make([]types.BridgeValidator, 0, len(allRelayers))
	var totalPower uint64

	for _, relayer := range allRelayers {
		power := relayer.GetPower()
		if power.LTE(sdkmath.ZeroInt()) {
			continue
		}
		totalPower += power.Uint64()
		bridgeValidators = append(bridgeValidators, types.BridgeValidator{
			Power:           power.Uint64(),
			ExternalAddress: relayer.ExternalAddress,
		})
	}
	for i := range bridgeValidators {
		// normalize power, use 10000 as the base, meaning 50.01% is 5001.
		bridgeValidators[i].Power = sdkmath.NewUint(bridgeValidators[i].Power).MulUint64(types.PowerBase).QuoUint64(totalPower).Uint64()
	}
	relayerSetNonce := k.GetLastRelayerSetNonce(ctx)
	return types.CurrentRelayerSet(relayerSetNonce, uint64(ctx.BlockHeight()), bridgeValidators)
}

// AddRelayerSetChangeRequest returns a new instance of the Relayer BridgeValidatorSet
func (k Keeper) AddRelayerSetChangeRequest(ctx sdk.Context, currentRelayerSet *types.RelayerSet) {
	// if CurrentRelayerSet member is empty, not store RelayerSet.
	if len(currentRelayerSet.Members) == 0 {
		return
	}
	currentRelayerSet.Nonce = k.GetLastRelayerSetNonce(ctx) + 1
	k.StoreRelayerSet(ctx, currentRelayerSet)
	k.SetLastRelayerSetNonce(ctx, currentRelayerSet.Nonce)
	k.SetLastTotalPower(ctx)

	// checkpoint, err := CurrentRelayerSet.GetCheckpoint(k.GetRelayerID(ctx))
	// if err != nil {
	// 	panic(err)
	// }
	// k.SetPastExternalSignatureCheckpoint(ctx, checkpoint)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRelayerSetUpdate,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyRelayerSetNonce, fmt.Sprint(currentRelayerSet.Nonce)),
		sdk.NewAttribute(types.AttributeKeyRelayerSetLen, fmt.Sprint(len(currentRelayerSet.Members))),
	))
}

// StoreRelayerSet is for storing a relayer set at a given height
func (k Keeper) StoreRelayerSet(ctx sdk.Context, relayerSet *types.RelayerSet) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetRelayerSetKey(relayerSet.Nonce), k.cdc.MustMarshal(relayerSet))
}

// HasRelayerSetRequest returns true if a relayerSet defined by a nonce exists
func (k Keeper) HasRelayerSetRequest(ctx sdk.Context, nonce uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetRelayerSetKey(nonce))
}

// DeleteRelayerSet deletes the relayerSet at a given nonce from state
func (k Keeper) DeleteRelayerSet(ctx sdk.Context, nonce uint64) {
	ctx.KVStore(k.storeKey).Delete(types.GetRelayerSetKey(nonce))
}

// SetLatestRelayerSetNonce sets the latest relayerSet nonce
func (k Keeper) SetLastRelayerSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LatestRelayerSetNonce, sdk.Uint64ToBigEndian(nonce))
}

// GetLastRelayerSetNonce returns the latest relayerSet nonce
func (k Keeper) GetLastRelayerSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LatestRelayerSetNonce)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}

// GetRelayerSet returns a relayerSet by nonce
func (k Keeper) GetRelayerSet(ctx sdk.Context, nonce uint64) *types.RelayerSet {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRelayerSetKey(nonce))
	if bz == nil {
		return nil
	}
	var relayerSet types.RelayerSet
	k.cdc.MustUnmarshal(bz, &relayerSet)
	return &relayerSet
}

// IterateRelayerSets returns all relayerSet
func (k Keeper) IterateRelayerSets(ctx sdk.Context, reverse bool, cb func(*types.RelayerSet) bool) {
	store := ctx.KVStore(k.storeKey)
	var iter sdk.Iterator
	if reverse {
		iter = sdk.KVStoreReversePrefixIterator(store, types.RelayerSetRequestKey)
	} else {
		iter = sdk.KVStorePrefixIterator(store, types.RelayerSetRequestKey)
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var relayerSet types.RelayerSet
		k.cdc.MustUnmarshal(iter.Value(), &relayerSet)
		// cb returns true to stop early
		if cb(&relayerSet) {
			break
		}
	}
}

// GetRelayerSets used in testing
func (k Keeper) GetRelayerSets(ctx sdk.Context) (relayerSets types.RelayerSets) {
	k.IterateRelayerSets(ctx, false, func(relayer *types.RelayerSet) bool {
		relayerSets = append(relayerSets, relayer)
		return false
	})
	return
}

// GetLastRelayerSet returns the latest relayer set in state
func (k Keeper) GetLastRelayerSet(ctx sdk.Context) *types.RelayerSet {
	latestRelayerSetNonce := k.GetLastRelayerSetNonce(ctx)
	return k.GetRelayerSet(ctx, latestRelayerSetNonce)
}

// SetLastSlashedRelayerSetNonce sets the latest slashed relayerSet nonce
func (k Keeper) SetLastSlashedRelayerSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedRelayerSetNonce, sdk.Uint64ToBigEndian(nonce))
}

// GetLastSlashedRelayerSetNonce returns the latest slashed relayerSet nonce
func (k Keeper) GetLastSlashedRelayerSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LastSlashedRelayerSetNonce)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}

// GetUnSlashedRelayerSets returns all the unSlashed relayer sets in state
func (k Keeper) GetUnSlashedRelayerSets(ctx sdk.Context, maxHeight uint64) (relayerSets types.RelayerSets) {
	lastSlashedRelayerSetNonce := k.GetLastSlashedRelayerSetNonce(ctx) + 1
	k.IterateRelayerSetByNonce(ctx, lastSlashedRelayerSetNonce, func(relayerSet *types.RelayerSet) bool {
		if maxHeight > relayerSet.Height {
			relayerSets = append(relayerSets, relayerSet)
			return false
		}
		return true
	})
	return
}

// IterateRelayerSetByNonce iterates through all relayerSet by nonce
func (k Keeper) IterateRelayerSetByNonce(ctx sdk.Context, startNonce uint64, cb func(*types.RelayerSet) bool) {
	store := ctx.KVStore(k.storeKey)
	startKey := append(types.RelayerSetRequestKey, sdk.Uint64ToBigEndian(startNonce)...)
	endKey := append(types.RelayerSetRequestKey, sdk.Uint64ToBigEndian(math.MaxUint64)...)
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		relayerSet := new(types.RelayerSet)
		k.cdc.MustUnmarshal(iter.Value(), relayerSet)
		// cb returns true to stop early
		if cb(relayerSet) {
			break
		}
	}
}

// --- RELAYER SET CONFIRMS --- //

// GetRelayerSetConfirm returns a relayerSet confirmation by a nonce and external address
func (k Keeper) GetRelayerSetConfirm(ctx sdk.Context, nonce uint64, relayerAddr sdk.AccAddress) *types.MsgRelayerSetConfirm {
	store := ctx.KVStore(k.storeKey)
	entity := store.Get(types.GetRelayerSetConfirmKey(nonce, relayerAddr))
	if entity == nil {
		return nil
	}
	confirm := types.MsgRelayerSetConfirm{}
	k.cdc.MustUnmarshal(entity, &confirm)
	return &confirm
}

// SetRelayerSetConfirm sets a relayerSet confirmation
func (k Keeper) SetRelayerSetConfirm(ctx sdk.Context, relayerAddr sdk.AccAddress, relayerSetConfirm *types.MsgRelayerSetConfirm) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRelayerSetConfirmKey(relayerSetConfirm.Nonce, relayerAddr)
	store.Set(key, k.cdc.MustMarshal(relayerSetConfirm))
}

// IterateRelayerSetConfirmByNonce iterates through all relayerSet confirms by nonce
func (k Keeper) IterateRelayerSetConfirmByNonce(ctx sdk.Context, nonce uint64, cb func(*types.MsgRelayerSetConfirm) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetRelayerSetConfirmKey(nonce, sdk.AccAddress{}))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := new(types.MsgRelayerSetConfirm)
		k.cdc.MustUnmarshal(iter.Value(), confirm)
		// cb returns true to stop early
		if cb(confirm) {
			break
		}
	}
}

func (k Keeper) DeleteRelayerSetConfirm(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetRelayerSetConfirmKey(nonce, sdk.AccAddress{}))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// GetLastObservedRelayerSet retrieves the last observed relayer set from the store
// WARNING: This value is not an up to date relayer set on Ethereum, it is a relayer set
// that AT ONE POINT was the one in the bridge on Ethereum. If you assume that it's up
// to date you may break the bridge
func (k Keeper) GetLastObservedRelayerSet(ctx sdk.Context) *types.RelayerSet {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedRelayerSetKey)

	if len(bytes) == 0 {
		return nil
	}
	relayerSet := types.RelayerSet{}
	k.cdc.MustUnmarshal(bytes, &relayerSet)
	return &relayerSet
}

// SetLastObservedRelayerSet updates the last observed relayer set in the store
func (k Keeper) SetLastObservedRelayerSet(ctx sdk.Context, relayerSet *types.RelayerSet) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedRelayerSetKey, k.cdc.MustMarshal(relayerSet))
}

func (k Keeper) DelLastObservedRelayerSet(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.LastObservedRelayerSetKey)
}

func (k Keeper) GetLastRelayerSlashBlockHeight(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LastRelayerSlashBlockHeight)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}
