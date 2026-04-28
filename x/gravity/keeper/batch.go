package keeper

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// BuildOutgoingTxBatch starts the following process chain:
//   - find bridged denominator for given voucher type
//   - determine if a an unExecuted batch is already waiting for this token type, if so confirm the new batch would
//     have a higher total fees. If not exit without creating a batch
//   - select available transactions from the outgoing transaction pool sorted by fee desc
//   - persist an outgoing batch object with an incrementing ID = nonce
//   - emit an event
func (k Keeper) BuildOutgoingTxBatch(ctx sdk.Context, contractAddress, feeReceive string, maxElements uint, minimumFee, baseFee sdkmath.Int) (*types.OutgoingTxBatch, error) {
	if maxElements == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalid, "max elements value")
	}
	projectedCurrentExternalHeight, batchTimeout := k.GetBatchTimeoutHeight(ctx)
	if batchTimeout <= 0 {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "batch timeout height %d less than 0", batchTimeout)
	}

	// if there is a more profitable batch for this token type do not create a new batch
	if lastBatch := k.GetLastOutgoingBatchByTokenType(ctx, contractAddress); lastBatch != nil {
		if lastBatch.BatchTimeout < projectedCurrentExternalHeight {
			return nil, errorsmod.Wrap(types.ErrInvalid, "existing unexecuted batch, and the batch not timeout")
		}
	}
	selectedTx, err := k.pickUnBatchedTx(ctx, contractAddress, maxElements, baseFee)
	if err != nil {
		return nil, err
	}
	if len(selectedTx) == 0 {
		return nil, errorsmod.Wrap(types.ErrEmpty, "no batch tx selected")
	}
	if types.OutgoingTransferTxs(selectedTx).TotalFee().LT(minimumFee) {
		return nil, errorsmod.Wrap(types.ErrInvalid, "total fee less than minimum fee")
	}

	nextID := k.AutoIncrementID(ctx, types.KeyLastOutgoingBatchID)
	batch := &types.OutgoingTxBatch{
		BatchNonce:    nextID,
		BatchTimeout:  batchTimeout,
		Transactions:  selectedTx,
		TokenContract: contractAddress,
		FeeReceive:    feeReceive,
		Block:         uint64(ctx.BlockHeight()), // set the current block height when storing the batch
	}
	if err = k.StoreBatch(ctx, batch); err != nil {
		return nil, err
	}

	// checkpoint, err := batch.GetCheckpoint(k.GetGravityID(ctx))
	// if err != nil {
	// 	panic(err)
	// }
	// k.SetPastExternalSignatureCheckpoint(ctx, checkpoint)

	eventBatchNonceTxIds := strings.Builder{}
	eventBatchNonceTxIds.WriteString(fmt.Sprintf("%d", selectedTx[0].Id))
	for _, tx := range selectedTx[1:] {
		_, _ = eventBatchNonceTxIds.WriteString(fmt.Sprintf(",%d", tx.Id))
	}
	batchEvent := sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchNonce, fmt.Sprint(nextID)),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxIds, eventBatchNonceTxIds.String()),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchTimeout, fmt.Sprint(batch.BatchTimeout)),
		sdk.NewAttribute(types.AttributeKeyTokenContract, contractAddress),
	)
	ctx.EventManager().EmitEvent(batchEvent)
	return batch, nil
}

// GetBatchTimeoutHeight This gets the batch timeout height in External blocks.
func (k Keeper) GetBatchTimeoutHeight(ctx sdk.Context) (uint64, uint64) {
	currentMeHeight := ctx.BlockHeight()
	params := k.GetParams(ctx)
	if params.AverageExternalBlockTime == 0 {
		return 0, 0
	}
	// we store the last observed Cosmos and Ethereum heights, we do not concern ourselves if these values
	// are zero because no batch can be produced if the last Ethereum block height is not first populated by a deposit event.
	heights := k.GetLastObservedBlockHeight(ctx)
	if heights.ExternalBlockHeight == 0 {
		return 0, 0
	}
	// we project how long it has been in milliseconds since the last Ethereum block height was observed
	projectedMillis := (uint64(currentMeHeight) - heights.BlockHeight) * params.AverageBlockTime
	// we convert that projection into the current Ethereum height using the average Ethereum block time in millis
	projectedCurrentEthereumHeight := (projectedMillis / params.AverageExternalBlockTime) + heights.ExternalBlockHeight
	// we convert our target time for block timeouts (lets say 12 hours) into a number of blocks to
	// place on top of our projection of the current Ethereum block height.
	blocksToAdd := params.ExternalBatchTimeout / params.AverageExternalBlockTime
	return projectedCurrentEthereumHeight, projectedCurrentEthereumHeight + blocksToAdd
}

// OutgoingTxBatchExecuted is run when the Cosmos chain detects that a batch has been executed on Ethereum
// It frees all the transactions in the batch, then cancels all earlier batches
func (k Keeper) OutgoingTxBatchExecuted(ctx sdk.Context, contractAddress string, batchNonce uint64) {
	batch := k.GetOutgoingTxBatch(ctx, contractAddress, batchNonce)
	if batch == nil {
		panic(fmt.Sprintf("unknown batch nonce for outgoing tx batch %s %d", contractAddress, batchNonce))
	}

	// Iterate through remaining batches
	k.IterateOutgoingTxBatches(ctx, func(iterBatch *types.OutgoingTxBatch) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		if iterBatch.BatchNonce < batch.BatchNonce && iterBatch.TokenContract == contractAddress {
			if err := k.CancelOutgoingTxBatch(ctx, contractAddress, iterBatch.BatchNonce); err != nil {
				panic(fmt.Sprintf("Failed cancel out batch %s-%d while trying to execute failed: %s", batch.TokenContract, batch.BatchNonce, err))
			}
		}
		return false
	})

	// Delete batch since it is finished
	k.DeleteBatch(ctx, batch)
	k.DeleteBatchConfirm(ctx, batch.BatchNonce, batch.TokenContract)
}

// StoreBatch stores a transaction batch
func (k Keeper) StoreBatch(ctx sdk.Context, batch *types.OutgoingTxBatch) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce)

	blockKey := types.GetOutgoingTxBatchBlockKey(batch.Block)
	// Note: Only one OutgoingTxBatch can be submitted in a block
	if store.Has(blockKey) {
		return errorsmod.Wrap(types.ErrInvalid, fmt.Sprintf("block:[%v] has batch request", batch.Block))
	}

	value := k.cdc.MustMarshal(batch)
	store.Set(key, value)
	store.Set(blockKey, value)
	return nil
}

// DeleteBatch deletes an outgoing transaction batch
func (k Keeper) DeleteBatch(ctx sdk.Context, batch *types.OutgoingTxBatch) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce))
	store.Delete(types.GetOutgoingTxBatchBlockKey(batch.Block))
}

// pickUnBatchedTx find Tx in pool and remove from "available" second index
func (k Keeper) pickUnBatchedTx(ctx sdk.Context, contractAddress string, maxElements uint, baseFee sdkmath.Int) ([]*types.OutgoingTransferTx, error) {
	var selectedTx []*types.OutgoingTransferTx
	var err error
	k.IterateUnbatchedTransactions(ctx, contractAddress, func(tx *types.OutgoingTransferTx) bool {
		if tx.Fee.Amount.LT(baseFee) {
			return true
		}
		selectedTx = append(selectedTx, tx)
		err = k.DelUnbatchedTx(ctx, tx.Fee, tx.Id)
		oldTx, oldTxErr := k.GetUnbatchedTxByFeeAndId(ctx, tx.Fee, tx.Id)
		if oldTx != nil || oldTxErr == nil {
			panic("picked a duplicate transaction from the pool, duplicates should never exist!")
		}
		return err != nil || uint(len(selectedTx)) == maxElements
	})
	return selectedTx, err
}

// GetOutgoingTxBatch loads a batch object. Returns nil when not exists.
func (k Keeper) GetOutgoingTxBatch(ctx sdk.Context, contractAddress string, batchNonce uint64) *types.OutgoingTxBatch {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOutgoingTxBatchKey(contractAddress, batchNonce)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	batch := new(types.OutgoingTxBatch)
	k.cdc.MustUnmarshal(bz, batch)
	return batch
}

// CancelOutgoingTxBatch releases all TX in the batch and deletes the batch
func (k Keeper) CancelOutgoingTxBatch(ctx sdk.Context, contractAddress string, batchNonce uint64) error {
	batch := k.GetOutgoingTxBatch(ctx, contractAddress, batchNonce)
	if batch == nil {
		return types.ErrUnknown
	}
	for _, tx := range batch.Transactions {
		if err := k.AddUnbatchedTx(ctx, tx); err != nil {
			return errorsmod.Wrapf(err, "unable to add batched transaction back into pool %v", tx)
		}
	}

	// Delete batch since it is finished
	k.DeleteBatch(ctx, batch)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOutgoingBatchCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyTokenContract, contractAddress),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchNonce, fmt.Sprint(batchNonce)),
	))
	return nil
}

// IterateOutgoingTxBatches iterates through all outgoing batches
func (k Keeper) IterateOutgoingTxBatches(ctx sdk.Context, cb func(batch *types.OutgoingTxBatch) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStoreReversePrefixIterator(store, types.OutgoingTxBatchKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		batch := new(types.OutgoingTxBatch)
		k.cdc.MustUnmarshal(iter.Value(), batch)
		// cb returns true to stop early
		if cb(batch) {
			break
		}
	}
}

// GetOutgoingTxBatches used in testing
func (k Keeper) GetOutgoingTxBatches(ctx sdk.Context) (out []*types.OutgoingTxBatch) {
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		out = append(out, batch)
		return false
	})
	return
}

// GetLastOutgoingBatchByTokenType gets the latest outgoing tx batch by token type
func (k Keeper) GetLastOutgoingBatchByTokenType(ctx sdk.Context, token string) *types.OutgoingTxBatch {
	var lastBatch *types.OutgoingTxBatch = nil
	lastNonce := uint64(0)
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		if batch.TokenContract == token && batch.BatchNonce > lastNonce {
			lastBatch = batch
			lastNonce = batch.BatchNonce
		}
		return false
	})
	return lastBatch
}

// IterateBatchByBlockHeight iterates through all Batch by block in the half-open interval [start,end)
func (k Keeper) IterateBatchByBlockHeight(ctx sdk.Context, start uint64, end uint64, cb func(*types.OutgoingTxBatch) bool) {
	store := ctx.KVStore(k.storeKey)
	startKey := append(types.OutgoingTxBatchBlockKey, sdk.Uint64ToBigEndian(start)...)
	endKey := append(types.OutgoingTxBatchBlockKey, sdk.Uint64ToBigEndian(end)...)
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		batch := new(types.OutgoingTxBatch)
		k.cdc.MustUnmarshal(iter.Value(), batch)
		// cb returns true to stop early
		if cb(batch) {
			break
		}
	}
}
