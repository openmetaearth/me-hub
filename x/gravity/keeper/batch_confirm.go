package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// --- BATCH CONFIRMS --- //

// GetBatchConfirm returns a batch confirmation given its nonce, the token contract, and a relayer address
func (k Keeper) GetBatchConfirm(ctx sdk.Context, tokenContract string, batchNonce uint64, relayerAddr sdk.AccAddress) *types.MsgConfirmBatch {
	store := ctx.KVStore(k.storeKey)
	entity := store.Get(types.GetBatchConfirmKey(tokenContract, batchNonce, relayerAddr))
	if entity == nil {
		return nil
	}
	confirm := types.MsgConfirmBatch{}
	k.cdc.MustUnmarshal(entity, &confirm)
	return &confirm
}

// SetBatchConfirm sets a batch confirmation by a relayer
func (k Keeper) SetBatchConfirm(ctx sdk.Context, relayerAddr sdk.AccAddress, batch *types.MsgConfirmBatch) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetBatchConfirmKey(batch.TokenContract, batch.Nonce, relayerAddr)
	store.Set(key, k.cdc.MustMarshal(batch))
}

// IterateBatchConfirmByNonceAndTokenContract iterates through all batch confirmations
func (k Keeper) IterateBatchConfirmByNonceAndTokenContract(ctx sdk.Context, batchNonce uint64, tokenContract string, cb func(*types.MsgConfirmBatch) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBatchConfirmKey(tokenContract, batchNonce, []byte{}))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := new(types.MsgConfirmBatch)
		k.cdc.MustUnmarshal(iter.Value(), confirm)
		// cb returns true to stop early
		if cb(confirm) {
			break
		}
	}
}

func (k Keeper) HasBatchConfirmQuorum(ctx sdk.Context, batchNonce uint64, tokenContract string) bool {
	totalPower := k.GetLastTotalPower(ctx)
	if totalPower.IsZero() {
		return false
	}

	requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(sdk.NewIntFromUint64(types.PowerBase))
	confirmPower := sdkmath.ZeroInt()
	k.IterateBatchConfirmByNonceAndTokenContract(ctx, batchNonce, tokenContract, func(confirm *types.MsgConfirmBatch) bool {
		relayerAddr, err := sdk.AccAddressFromBech32(confirm.RelayerAddress)
		if err != nil {
			k.Logger(ctx).Error("HasBatchConfirmQuorum", "invalid relayer address", confirm.RelayerAddress, "error", err)
			return false
		}

		relayer, found := k.GetRelayer(ctx, relayerAddr)
		if !found || !relayer.Online {
			return false
		}

		confirmPower = confirmPower.Add(relayer.GetPower())
		return confirmPower.GTE(requiredPower)
	})
	return confirmPower.GTE(requiredPower)
}

func (k Keeper) DeleteBatchConfirm(ctx sdk.Context, batchNonce uint64, tokenContract string) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBatchConfirmKey(tokenContract, batchNonce, []byte{}))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// --- LAST SLASHED OUTGOING BATCH BLOCK --- //

// SetLastSlashedBatchBlock sets the latest slashed Batch block height
func (k Keeper) SetLastSlashedBatchBlock(ctx sdk.Context, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedBatchBlock, sdk.Uint64ToBigEndian(blockHeight))
}

// GetLastSlashedBatchBlock returns the latest slashed Batch block
func (k Keeper) GetLastSlashedBatchBlock(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastSlashedBatchBlock)
	if len(bytes) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(bytes)
}

// GetUnSlashedBatches returns all the unSlashed batches in state
func (k Keeper) GetUnSlashedBatches(ctx sdk.Context, maxHeight uint64) (outgoingTxBatches types.OutgoingTxBatches) {
	lastSlashedBatchBlock := k.GetLastSlashedBatchBlock(ctx) + 1
	k.IterateBatchByBlockHeight(ctx, lastSlashedBatchBlock, maxHeight, func(batch *types.OutgoingTxBatch) bool {
		outgoingTxBatches = append(outgoingTxBatches, batch)
		return false
	})
	return
}
