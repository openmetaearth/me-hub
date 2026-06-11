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
		k.Logger(ctx).Error("failed to store batch", "error", err)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeError,
				sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
				sdk.NewAttribute(sdk.AttributeKeyError, err.Error()),
			),
		)
		return nil, err
	}

	// checkpoint, err := batch.GetCheckpoint(k.GetGravityID(ctx))
	// if err != nil {
	// 	k.Logger(ctx).Error("failed to get checkpoint", "error", err)
	// 	ctx.EventManager().EmitEvent(
	// 		sdk.NewEvent(
	// 			sdk.EventTypeError,
	// 			sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
	// 			sdk.NewAttribute(sdk.AttributeKeyError, err.Error()),
	// 		),
	// 	)
	// 	return nil, err
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
		k.Logger(ctx).Error("average external block time is zero", "error", "invalid average external block time")
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeError,
				sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
				sdk.NewAttribute(sdk.AttributeKeyError, "invalid average external block time"),
			),
		)
		return 0, 0
	}
	// we store the last observed Cosmos and Ethereum heights, we do not concern ourselves if these values
	// are zero because no batch can be produced if the last Ethereum block height is not first populated by a deposit event.
	heights := k.GetLastObservedBlockHeight(ctx)
	if heights.ExternalBlockHeight == 0 {
		k.Logger(ctx).Error("external block height is zero", "error", "invalid external block height")
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeError,
				sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
				sdk.NewAttribute(sdk.AttributeKeyError, "invalid external block height"),
			),
		)
		return 0, 0
	}
	// we project how long it has been in milliseconds since the last Ethereum block height was observed
	projectedMillis := (uint64(currentMeHeight) - heights.CosmosBlockHeight) * params.AverageExternalBlockTime
	// we calculate the batch timeout height by adding the projected milliseconds to the last observed Ethereum block height
	// and then dividing by the average external block time
	batchTimeout := (projectedMillis / params.AverageExternalBlockTime) + heights.ExternalBlockHeight
	return projectedMillis / params.AverageExternalBlockTime, batchTimeout
}

// pickUnBatchedTx picks unbatched transactions from the outgoing transaction pool
func (k Keeper) pickUnBatchedTx(ctx sdk.Context, contractAddress string, maxElements uint, baseFee sdkmath.Int) ([]*types.OutgoingTransferTx, error) {
	// implementation remains the same
}

// OutgoingTxBatchExecuted handles the execution of an outgoing batch
func (k Keeper) OutgoingTxBatchExecuted(ctx sdk.Context, batch *types.OutgoingTxBatch) error {
	if batch == nil {
		k.Logger(ctx).Error("batch not found", "error", "batch not found")
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeError,
				sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
				sdk.NewAttribute(sdk.AttributeKeyError, "batch not found"),
			),
		)
		return errorsmod.Wrap(types.ErrInvalid, "batch not found")
	}
	// implementation remains the same
}

// CancelOutgoingTxBatch cancels an outgoing batch
func (k Keeper) CancelOutgoingTxBatch(ctx sdk.Context, batch *types.OutgoingTxBatch) error {
	if batch == nil {
		k.Logger(ctx).Error("batch not found", "error", "batch not found")
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeError,
				sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
				sdk.NewAttribute(sdk.AttributeKeyError, "batch not found"),
			),
		)
		return errorsmod.Wrap(types.ErrInvalid, "batch not found")
	}
	// implementation remains the same
}