package keeper

import (
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/st-chain/me-hub/x/gravity/types"
)

// --- ORACLE SET REQUESTS --- //

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
	oracleSetNonce := k.GetLatestRelayerSetNonce(ctx) + 1
	return types.NewRelayerSet(oracleSetNonce, uint64(ctx.BlockHeight()), bridgeValidators)
}

// AddRelayerSetRequest returns a new instance of the Relayer BridgeValidatorSet
func (k Keeper) AddRelayerSetRequest(ctx sdk.Context, currentRelayerSet *types.RelayerSet) {
	// if currentRelayerSet member is empty, not store RelayerSet.
	if len(currentRelayerSet.Members) == 0 {
		return
	}
	k.StoreRelayerSet(ctx, currentRelayerSet)
	k.SetLatestRelayerSetNonce(ctx, currentRelayerSet.Nonce)
	k.SetLastTotalPower(ctx)

	// checkpoint, err := currentRelayerSet.GetCheckpoint(k.GetRelayerID(ctx))
	// if err != nil {
	// 	panic(err)
	// }
	// k.SetPastExternalSignatureCheckpoint(ctx, checkpoint)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRelayerSetUpdate,
		sdk.NewAttribute(types.AttributeKeyRelayerSetNonce, fmt.Sprint(currentRelayerSet.Nonce)),
		sdk.NewAttribute(types.AttributeKeyRelayerSetLen, fmt.Sprint(len(currentRelayerSet.Members))),
	))
}

// StoreRelayerSet is for storing a oracle set at a given height
func (k Keeper) StoreRelayerSet(ctx sdk.Context, oracleSet *types.RelayerSet) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetRelayerSetKey(oracleSet.Nonce), k.cdc.MustMarshal(oracleSet))
}

// HasRelayerSetRequest returns true if a oracleSet defined by a nonce exists
func (k Keeper) HasRelayerSetRequest(ctx sdk.Context, nonce uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetRelayerSetKey(nonce))
}

// DeleteRelayerSet deletes the oracleSet at a given nonce from state
func (k Keeper) DeleteRelayerSet(ctx sdk.Context, nonce uint64) {
	ctx.KVStore(k.storeKey).Delete(types.GetRelayerSetKey(nonce))
}

// SetLatestRelayerSetNonce sets the latest oracleSet nonce
func (k Keeper) SetLatestRelayerSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LatestRelayerSetNonce, sdk.Uint64ToBigEndian(nonce))
}

// GetLatestRelayerSetNonce returns the latest oracleSet nonce
func (k Keeper) GetLatestRelayerSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LatestRelayerSetNonce)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}

// GetRelayerSet returns a oracleSet by nonce
func (k Keeper) GetRelayerSet(ctx sdk.Context, nonce uint64) *types.RelayerSet {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRelayerSetKey(nonce))
	if bz == nil {
		return nil
	}
	var oracleSet types.RelayerSet
	k.cdc.MustUnmarshal(bz, &oracleSet)
	return &oracleSet
}

// IterateRelayerSets returns all oracleSet
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
		var oracleSet types.RelayerSet
		k.cdc.MustUnmarshal(iter.Value(), &oracleSet)
		// cb returns true to stop early
		if cb(&oracleSet) {
			break
		}
	}
}

// GetRelayerSets used in testing
func (k Keeper) GetRelayerSets(ctx sdk.Context) (oracleSets types.RelayerSets) {
	k.IterateRelayerSets(ctx, false, func(set *types.RelayerSet) bool {
		oracleSets = append(oracleSets, set)
		return false
	})
	return
}

// GetLatestRelayerSet returns the latest oracle set in state
func (k Keeper) GetLatestRelayerSet(ctx sdk.Context) *types.RelayerSet {
	latestRelayerSetNonce := k.GetLatestRelayerSetNonce(ctx)
	return k.GetRelayerSet(ctx, latestRelayerSetNonce)
}

// SetLastSlashedRelayerSetNonce sets the latest slashed oracleSet nonce
func (k Keeper) SetLastSlashedRelayerSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedRelayerSetNonce, sdk.Uint64ToBigEndian(nonce))
}

// GetLastSlashedRelayerSetNonce returns the latest slashed oracleSet nonce
func (k Keeper) GetLastSlashedRelayerSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LastSlashedRelayerSetNonce)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}

// GetUnSlashedRelayerSets returns all the unSlashed oracle sets in state
func (k Keeper) GetUnSlashedRelayerSets(ctx sdk.Context, maxHeight uint64) (oracleSets types.RelayerSets) {
	lastSlashedRelayerSetNonce := k.GetLastSlashedRelayerSetNonce(ctx) + 1
	k.IterateRelayerSetByNonce(ctx, lastSlashedRelayerSetNonce, func(oracleSet *types.RelayerSet) bool {
		if maxHeight > oracleSet.Height {
			oracleSets = append(oracleSets, oracleSet)
			return false
		}
		return true
	})
	return
}

// IterateRelayerSetByNonce iterates through all oracleSet by nonce
func (k Keeper) IterateRelayerSetByNonce(ctx sdk.Context, startNonce uint64, cb func(*types.RelayerSet) bool) {
	store := ctx.KVStore(k.storeKey)
	startKey := append(types.RelayerSetRequestKey, sdk.Uint64ToBigEndian(startNonce)...)
	endKey := append(types.RelayerSetRequestKey, sdk.Uint64ToBigEndian(math.MaxUint64)...)
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		oracleSet := new(types.RelayerSet)
		k.cdc.MustUnmarshal(iter.Value(), oracleSet)
		// cb returns true to stop early
		if cb(oracleSet) {
			break
		}
	}
}

// --- ORACLE SET CONFIRMS --- //

// GetRelayerSetConfirm returns a oracleSet confirmation by a nonce and external address
func (k Keeper) GetRelayerSetConfirm(ctx sdk.Context, nonce uint64, oracleAddr sdk.AccAddress) *types.MsgRelayerSetConfirm {
	store := ctx.KVStore(k.storeKey)
	entity := store.Get(types.GetRelayerSetConfirmKey(nonce, oracleAddr))
	if entity == nil {
		return nil
	}
	confirm := types.MsgRelayerSetConfirm{}
	k.cdc.MustUnmarshal(entity, &confirm)
	return &confirm
}

// SetRelayerSetConfirm sets a oracleSet confirmation
func (k Keeper) SetRelayerSetConfirm(ctx sdk.Context, oracleAddr sdk.AccAddress, oracleSetConfirm *types.MsgRelayerSetConfirm) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetRelayerSetConfirmKey(oracleSetConfirm.Nonce, oracleAddr)
	store.Set(key, k.cdc.MustMarshal(oracleSetConfirm))
}

// IterateRelayerSetConfirmByNonce iterates through all oracleSet confirms by nonce
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

// GetLastObservedRelayerSet retrieves the last observed oracle set from the store
// WARNING: This value is not an up to date oracle set on Ethereum, it is a oracle set
// that AT ONE POINT was the one in the bridge on Ethereum. If you assume that it's up
// to date you may break the bridge
func (k Keeper) GetLastObservedRelayerSet(ctx sdk.Context) *types.RelayerSet {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedRelayerSetKey)

	if len(bytes) == 0 {
		return nil
	}
	oracleSet := types.RelayerSet{}
	k.cdc.MustUnmarshal(bytes, &oracleSet)
	return &oracleSet
}

// SetLastObservedRelayerSet updates the last observed oracle set in the store
func (k Keeper) SetLastObservedRelayerSet(ctx sdk.Context, oracleSet *types.RelayerSet) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedRelayerSetKey, k.cdc.MustMarshal(oracleSet))
}

func (k Keeper) GetLastRelayerSlashBlockHeight(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.LastRelayerSlashBlockHeight)
	if len(data) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(data)
}
