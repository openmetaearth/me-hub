package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	"sort"
)

// GetBatchFeesByTokenType gets the fee the next batch of a given token type would
// have if created right now. This info is both presented to relayers for the purpose of determining
// when to request batches and also used by the batch creation process to decide not to create
// a new batch (fees must be increasing)
func (k Keeper) GetBatchFeesByTokenType(ctx sdk.Context, tokenContract string, maxElements uint, baseFee sdkmath.Int) *types.BatchFees {
	batchFee := &types.BatchFees{TokenContract: tokenContract, TotalFees: sdkmath.NewInt(0), TotalAmount: sdkmath.NewInt(0)}
	k.IterateUnbatchedTransactions(ctx, tokenContract, func(tx *types.OutgoingTransferTx) bool {
		if tx.Fee.Contract != tokenContract {
			//panic(fmt.Errorf("unexpected fee contract %s when getting batch fees for contract %s", tx.Fee.Contract, tokenContract))
			// Log error and skip transaction instead of panicking
			ctx.Logger().Error("unexpected fee contract", "got", tx.Fee.Contract, "expected", tokenContract)
			return false
		}
		if tx.Fee.Amount.LT(baseFee) {
			// sort by fee and use ReverseIterator, so the fee behind is less than base fee
			return true
		}
		batchFee.TotalFees = batchFee.TotalFees.Add(tx.Fee.Amount)
		batchFee.TotalAmount = batchFee.TotalAmount.Add(tx.Token.Amount)
		batchFee.TotalTxs += 1
		return batchFee.TotalTxs == uint64(maxElements)
	})
	return batchFee
}

// GetAllBatchFees creates a fee entry for every batch type currently in the store
// this can be used by relayers to determine what batch types are desirable to request
func (k Keeper) GetAllBatchFees(ctx sdk.Context, maxElements uint, minBatchFees []types.MinBatchFee) (batchFees []*types.BatchFees) {
	batchFeesMap := k.createBatchFees(ctx, maxElements, minBatchFees)
	// create array of batchFees
	for _, batchFee := range batchFeesMap {
		batchFees = append(batchFees, batchFee)
	}

	// quick sort by token to make this function safe for use
	// in consensus computations
	sort.Slice(batchFees, func(i, j int) bool {
		return batchFees[i].TokenContract < batchFees[j].TokenContract
	})

	return batchFees
}

// createBatchFees iterates over the unbatched transaction pool and creates batch token fee map
// Implicitly creates batches with the highest potential fee because the transaction keys enforce an order which goes
// fee contract address -> fee amount -> transaction nonce
func (k Keeper) createBatchFees(ctx sdk.Context, maxElements uint, minBatchFees []types.MinBatchFee) map[string]*types.BatchFees {
	batchFeesMap := make(map[string]*types.BatchFees)
	txCountMap := make(map[string]int)
	baseFees := types.MinBatchFeeToBaseFees(minBatchFees)

	k.IterateUnbatchedTransactions(ctx, "", func(tx *types.OutgoingTransferTx) bool {
		fee := tx.Fee

		baseFee, ok := baseFees[fee.Contract]
		if ok && fee.Amount.LT(baseFee) {
			return false // sort by token address and fee, behind have other token
		}

		if txCountMap[fee.Contract] < int(maxElements) {
			addFeeToMap(tx.Token, fee, batchFeesMap, txCountMap)
		}
		return false
	})

	return batchFeesMap
}

func (k Keeper) AddUnbatchedTxBridgeFee(ctx sdk.Context, txId uint64, sender sdk.AccAddress, addBridgeFee sdk.Coin) error {
	if ctx.IsZero() || txId < 1 || sender.Empty() || addBridgeFee.IsZero() {
		return errorsmod.Wrap(types.ErrInvalid, "arguments")
	}

	// check that we actually have a tx with that id and what it's details are
	tx, err := k.GetUnbatchedTxById(ctx, txId)
	if err != nil {
		return errorsmod.Wrapf(types.ErrInvalid, "txId %d not in unbatched index! Must be in a batch!", txId)
	}

	bridgeToken, err := k.GetBridgeTokenByDenom(ctx, addBridgeFee.Denom)
	if err != nil {
		return errorsmod.Wrapf(types.ErrInvalid, "get bridge token: %v", err)
	}

	if tx.Fee.Contract != bridgeToken.ContractAddress {
		return errorsmod.Wrap(types.ErrInvalid, "token not equal tx fee token")
	}

	// If it is an external blockchain asset we burn it send coins to module in prep for burn
	{
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, k.moduleName, sdk.NewCoins(addBridgeFee)); err != nil {
			return err
		}
		// burn vouchers to send them back to external blockchain
		if err := k.bankKeeper.BurnCoins(ctx, k.moduleName, sdk.NewCoins(addBridgeFee)); err != nil {
			return err
		}
	}

	if err := k.DelUnbatchedTx(ctx, tx.Fee, txId); err != nil {
		return err
	}

	// add bridge fee amount
	tx.Fee.Amount = tx.Fee.Amount.Add(addBridgeFee.Amount)
	if err := k.AddUnbatchedTx(ctx, tx); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeIncreaseBridgeFee,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxID, fmt.Sprint(tx.Id)),
		sdk.NewAttribute(types.AttributeKeyIncreaseFee, addBridgeFee.String()),
	))
	return nil
}

// Helper method for creating batch fees
func addFeeToMap(amt, fee types.ERC20Token, batchFeesMap map[string]*types.BatchFees, txCountMap map[string]int) {
	txCountMap[fee.Contract] = txCountMap[fee.Contract] + 1

	// add fee amount
	if _, ok := batchFeesMap[fee.Contract]; ok {
		batchFees := batchFeesMap[fee.Contract]
		batchFees.TotalFees = batchFees.TotalFees.Add(fee.Amount)
		batchFees.TotalTxs = batchFees.TotalTxs + 1
		batchFees.TotalAmount = batchFees.TotalAmount.Add(amt.Amount)
		batchFeesMap[fee.Contract] = batchFees
	} else {
		batchFeesMap[fee.Contract] = &types.BatchFees{
			TokenContract: fee.Contract,
			TotalFees:     fee.Amount,
			TotalTxs:      1,
			TotalAmount:   amt.Amount,
		}
	}
}
