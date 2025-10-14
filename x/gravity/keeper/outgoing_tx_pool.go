package keeper

import (
	"encoding/binary"
	"fmt"
	"sort"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/st-chain/me-hub/x/gravity/types"
)

// AddToOutgoingPool
// - checks a counterpart denominator exists for the given voucher type
// - burns the voucher for transfer amount and fees
// - persists an OutgoingTx
// - adds the TX to the `available` TX pool via a second index
func (k Keeper) AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, receiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error) {
	bridgeToken, _ := k.GetBridgeTokenByDenom(ctx, amount.Denom)
	if bridgeToken == nil {
		return 0, errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
	}

	// TODO: also check other outgoing txs to make sure we don't exceed supply
	if amount.Amount.GT(bridgeToken.Supply) {
		return 0, errorsmod.Wrapf(types.ErrInvalid, "transfer to %s chain, amount %s exceeds bridge token supply %s",
			k.moduleName, amount.Amount.String(), bridgeToken.Supply.String())
	}

	totalInVouchers := sdk.NewCoins(amount.Add(fee))

	// If it is an external blockchain asset we burn it send coins to module in prep for burn
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, k.moduleName, totalInVouchers); err != nil {
		return 0, err
	}

	// burn vouchers to send them back to external blockchain
	if err := k.bankKeeper.BurnCoins(ctx, k.moduleName, totalInVouchers); err != nil {
		return 0, err
	}

	// get next tx id from keeper
	nextTxID := k.autoIncrementID(ctx, types.KeyLastTxPoolID)

	// construct outgoing tx, as part of this process we represent
	// the token as an ERC20 token since it is preparing to go to ETH
	// rather than the denom that is the input to this function.
	outgoing := &types.OutgoingTransferTx{
		Id:          nextTxID,
		Sender:      sender.String(),
		DestAddress: receiver,
		Token:       types.NewERC20Token(amount.Amount, bridgeToken.Contract),
		Fee:         types.NewERC20Token(fee.Amount, bridgeToken.Contract),
	}

	if err := k.AddUnbatchedTx(ctx, outgoing); err != nil {
		return 0, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSendToExternal,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxID, fmt.Sprint(nextTxID)),
	))

	return nextTxID, nil
}

// RemoveFromOutgoingPoolAndRefund
// - checks that the provided tx actually exists
// - deletes the unbatched tx from the pool
// - issues the tokens back to the sender
//
//gocyclo:ignore
func (k Keeper) RemoveFromOutgoingPoolAndRefund(ctx sdk.Context, txId uint64, sender sdk.AccAddress) (sdk.Coin, error) {
	if ctx.IsZero() || txId < 1 || sender.Empty() {
		return sdk.Coin{}, errorsmod.Wrap(types.ErrInvalid, "arguments")
	}
	// check that we actually have a tx with that id and what it's details are
	tx, err := k.GetUnbatchedTxById(ctx, txId)
	if err != nil {
		return sdk.Coin{}, err
	}

	// Check that this user actually sent the transaction, this prevents someone from refunding someone
	// else transaction to themselves.
	txSender := sdk.MustAccAddressFromBech32(tx.Sender)
	if !txSender.Equals(sender) {
		return sdk.Coin{}, errorsmod.Wrapf(types.ErrInvalid, "Sender %s did not send Id %d", sender, txId)
	}

	// An inconsistent entry should never enter the store, but this is the ideal place to exploit
	// it such a bug if it did ever occur, so we should double check to be really sure
	if tx.Fee.Contract != tx.Token.Contract {
		return sdk.Coin{}, errorsmod.Wrapf(types.ErrInvalid, "Inconsistent tokens to cancel!: %s %s", tx.Fee.Contract, tx.Token.Contract)
	}

	// delete this tx from the pool
	if err = k.removeUnbatchedTx(ctx, tx.Fee, txId); err != nil {
		return sdk.Coin{}, errorsmod.Wrapf(types.ErrInvalid, "txId %d not in unbatched index! Must be in a batch!", txId)
	}
	// Make sure the tx was removed
	oldTx, oldTxErr := k.GetUnbatchedTxByFeeAndId(ctx, tx.Fee, tx.Id)
	if oldTx != nil || oldTxErr == nil {
		return sdk.Coin{}, errorsmod.Wrapf(types.ErrInvalid, "tx with id %d was not fully removed from the pool, a duplicate must exist", txId)
	}

	// query denom, if not exist, return error
	bridgeToken, _ := k.GetBridgeTokenByContract(ctx, tx.Token.Contract)
	if bridgeToken == nil {
		return sdk.Coin{}, errorsmod.Wrapf(types.ErrInvalid, "Invalid token, contract %s", tx.Token.Contract)
	}
	// reissue the amount and the fee
	totalToRefund := sdk.NewCoin(bridgeToken.Denom, tx.Token.Amount)
	totalToRefund.Amount = totalToRefund.Amount.Add(tx.Fee.Amount)
	totalToRefundCoins := sdk.NewCoins(totalToRefund)

	// check bridge denom is origin denom or converted alias
	if err = k.bankKeeper.MintCoins(ctx, k.moduleName, totalToRefundCoins); err != nil {
		return sdk.Coin{}, errorsmod.Wrapf(err, "mint vouchers coins: %s", totalToRefundCoins)
	}
	if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, sender, totalToRefundCoins); err != nil {
		return sdk.Coin{}, errorsmod.Wrap(err, "transfer vouchers")
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSendToExternalCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxID, fmt.Sprint(txId)),
	))
	return totalToRefund, nil
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
	bridgeToken, _ := k.GetBridgeTokenByDenom(ctx, addBridgeFee.Denom)
	if bridgeToken == nil {
		return errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
	}

	if tx.Fee.Contract != bridgeToken.Contract {
		return errorsmod.Wrap(types.ErrInvalid, "token not equal tx fee token")
	}

	// If it is an external blockchain asset we burn it send coins to module in prep for burn
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(addBridgeFee)); err != nil {
		return err
	}

	// burn vouchers to send them back to external blockchain
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(addBridgeFee)); err != nil {
		return err
	}

	if err := k.removeUnbatchedTx(ctx, tx.Fee, txId); err != nil {
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

// AddUnbatchedTx creates a new transaction in the pool
func (k Keeper) AddUnbatchedTx(ctx sdk.Context, outgoingTransferTx *types.OutgoingTransferTx) error {
	store := ctx.KVStore(k.storeKey)
	idxKey := types.GetOutgoingTxPoolKey(outgoingTransferTx.Fee, outgoingTransferTx.Id)
	if store.Has(idxKey) {
		return errorsmod.Wrap(types.ErrDuplicate, "transaction already in pool")
	}

	store.Set(idxKey, k.cdc.MustMarshal(outgoingTransferTx))
	return nil
}

// removeUnbatchedTXIndex removes the tx from the pool
func (k Keeper) removeUnbatchedTx(ctx sdk.Context, fee types.ERC20Token, txID uint64) error {
	store := ctx.KVStore(k.storeKey)
	idxKey := types.GetOutgoingTxPoolKey(fee, txID)
	if !store.Has(idxKey) {
		return errorsmod.Wrap(types.ErrUnknown, "pool transaction")
	}
	store.Delete(idxKey)
	return nil
}

// GetUnbatchedTxByFeeAndId grabs a tx from the pool given its fee and txID
func (k Keeper) GetUnbatchedTxByFeeAndId(ctx sdk.Context, fee types.ERC20Token, txID uint64) (*types.OutgoingTransferTx, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOutgoingTxPoolKey(fee, txID))
	if bz == nil {
		return nil, errorsmod.Wrap(types.ErrUnknown, "pool transaction")
	}
	var r types.OutgoingTransferTx
	err := k.cdc.Unmarshal(bz, &r)
	return &r, err
}

// GetUnbatchedTxById grabs a tx from the pool given only the txID
// note that due to the way unbatched txs are indexed, the GetUnbatchedTxByFeeAndId method is much faster
func (k Keeper) GetUnbatchedTxById(ctx sdk.Context, txID uint64) (*types.OutgoingTransferTx, error) {
	var r *types.OutgoingTransferTx = nil
	k.IterateUnbatchedTransactions(ctx, "", func(tx *types.OutgoingTransferTx) bool {
		if tx.Id == txID {
			r = tx
			return true
		}
		return false
	})

	if r == nil {
		// We have no return tx, it was either batched or never existed
		return nil, errorsmod.Wrap(types.ErrUnknown, "pool transaction")
	}
	return r, nil
}

// GetUnbatchedTransactions used in testing
func (k Keeper) GetUnbatchedTransactions(ctx sdk.Context) []*types.OutgoingTransferTx {
	var txs []*types.OutgoingTransferTx
	k.IterateUnbatchedTransactions(ctx, "", func(tx *types.OutgoingTransferTx) bool {
		txs = append(txs, tx)
		return false
	})
	return txs
}

// IterateUnbatchedTransactions iterates through all unbatched transactions
func (k Keeper) IterateUnbatchedTransactions(ctx sdk.Context, tokenContract string, cb func(tx *types.OutgoingTransferTx) bool) {
	store := ctx.KVStore(k.storeKey)
	prefixKey := types.GetOutgoingTxPoolContractPrefix(tokenContract)
	iter := sdk.KVStoreReversePrefixIterator(store, prefixKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var transact types.OutgoingTransferTx
		k.cdc.MustUnmarshal(iter.Value(), &transact)
		// cb returns true to stop early
		if cb(&transact) {
			break
		}
	}
}

// GetBatchFeesByTokenType gets the fee the next batch of a given token type would
// have if created right now. This info is both presented to relayers for the purpose of determining
// when to request batches and also used by the batch creation process to decide not to create
// a new batch (fees must be increasing)
func (k Keeper) GetBatchFeesByTokenType(ctx sdk.Context, tokenContract string, maxElements uint, baseFee sdkmath.Int) *types.BatchFees {
	batchFee := &types.BatchFees{TokenContract: tokenContract, TotalFees: sdkmath.NewInt(0), TotalAmount: sdkmath.NewInt(0)}
	k.IterateUnbatchedTransactions(ctx, tokenContract, func(tx *types.OutgoingTransferTx) bool {
		if tx.Fee.Contract != tokenContract {
			panic(fmt.Errorf("unexpected fee contract %s when getting batch fees for contract %s", tx.Fee.Contract, tokenContract))
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

func (k Keeper) autoIncrementID(ctx sdk.Context, idKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(idKey)
	var id uint64 = 1
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	bz = sdk.Uint64ToBigEndian(id + 1)
	store.Set(idKey, bz)
	return id
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
