package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// AddToOutgoingPool
// - checks a counterpart denominator exists for the given voucher type
// - burns the voucher for transfer amount and fees
// - persists an OutgoingTx
// - adds the TX to the `available` TX pool via a second index
func (k Keeper) AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, receiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error) {
	bridgeToken, err := k.GetBridgeTokenByDenom(ctx, amount.Denom)
	if err != nil {
		return 0, errorsmod.Wrapf(types.ErrInvalid, "get bridge token: %v", err)
	}

	totalInVouchers := amount.Add(fee)
	if totalInVouchers.Amount.GT(bridgeToken.Supply) {
		return 0, errorsmod.Wrapf(types.ErrInvalid, "%s exceeds bridge token supply %s in %s chain",
			totalInVouchers.Amount.String(), bridgeToken.Supply.String(), k.moduleName)
	}

	totalPending := k.GetOutgoingPendingTxTotal(ctx, k.moduleName, bridgeToken)
	if totalInVouchers.Amount.Add(totalPending).GT(bridgeToken.Supply) {
		return 0, errorsmod.Wrapf(types.ErrInvalid, "total pending amount %s plus current amount %s exceeds bridge token supply %s in %s chain",
			totalPending.String(), totalInVouchers.Amount.String(), bridgeToken.Supply.String(), k.moduleName)
	}

	sendCoins := sdk.NewCoins(totalInVouchers)
	// If it is an external blockchain asset we burn it send coins to module in prep for burn
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, k.moduleName, sendCoins); err != nil {
		return 0, err
	}

	// burn vouchers to send them back to external blockchain
	if err := k.bankKeeper.BurnCoins(ctx, k.moduleName, sendCoins); err != nil {
		return 0, err
	}

	bridgeToken.Supply = bridgeToken.Supply.Sub(totalInVouchers.Amount)
	k.SetBridgeToken(ctx, bridgeToken)

	// get next tx id from keeper
	nextTxID := k.AutoIncrementID(ctx, types.KeyLastTxPoolID)

	// construct outgoing tx, as part of this process we represent
	// the token as an ERC20 token since it is preparing to go to ETH
	// rather than the denom that is the input to this function.

	externalBurnAmount := types.GetExternalUnlockAmount(amount.Amount, k.moduleName, bridgeToken)
	externalFeeAmount := types.GetExternalUnlockAmount(fee.Amount, k.moduleName, bridgeToken)
	outgoing := &types.OutgoingTransferTx{
		Id:          nextTxID,
		Sender:      sender.String(),
		DestAddress: receiver,
		Token:       types.NewERC20Token(externalBurnAmount, bridgeToken.ContractAddress),
		Fee:         types.NewERC20Token(externalFeeAmount, bridgeToken.ContractAddress),
	}

	if err := k.AddUnbatchedTx(ctx, outgoing); err != nil {
		return 0, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSendToExternal,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxID, fmt.Sprint(nextTxID)),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeKeyReceiver, receiver),
		sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
		sdk.NewAttribute(types.AttributeKeyBridgeFee, fee.String()),
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
	if err = k.DelUnbatchedTx(ctx, tx.Fee, txId); err != nil {
		return sdk.Coin{}, errorsmod.Wrapf(types.ErrInvalid, "txId %d not in unbatched index! Must be in a batch!", txId)
	}
	// Make sure the tx was removed
	oldTx, oldTxErr := k.GetUnbatchedTxByFeeAndId(ctx, tx.Fee, tx.Id)
	if oldTx != nil || oldTxErr == nil {
		return sdk.Coin{}, errorsmod.Wrapf(types.ErrInvalid, "tx with id %d was not fully removed from the pool, a duplicate must exist", txId)
	}

	// query denom, if not exist, return error
	bridgeToken, err := k.GetBridgeTokenByContract(ctx, tx.Token.Contract)
	if err != nil {
		return sdk.Coin{}, errorsmod.Wrapf(types.ErrInvalid, "Invalid token, contract %s", tx.Token.Contract)
	}
	// reissue the amount and the fee
	totalRefund := types.GetMintCoin(tx.Token.Amount.Add(tx.Fee.Amount), k.moduleName, bridgeToken)
	totalRefundCoins := sdk.NewCoins(totalRefund)

	// check bridge denom is origin denom or converted alias
	if err = k.bankKeeper.MintCoins(ctx, k.moduleName, totalRefundCoins); err != nil {
		return sdk.Coin{}, errorsmod.Wrapf(err, "mint vouchers coins: %s", totalRefundCoins)
	}
	if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, sender, totalRefundCoins); err != nil {
		return sdk.Coin{}, errorsmod.Wrap(err, "transfer vouchers")
	}

	bridgeToken.Supply = bridgeToken.Supply.Add(totalRefund.Amount)
	k.SetBridgeToken(ctx, bridgeToken)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSendToExternalCanceled,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(types.AttributeKeyOutgoingTxID, fmt.Sprint(txId)),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeKeyRefundAmount, totalRefundCoins.String()),
	))
	return totalRefund, nil
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

// DelUnbatchedTxIndex removes the tx from the pool
func (k Keeper) DelUnbatchedTx(ctx sdk.Context, fee types.ERC20Token, txID uint64) error {
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

func (k Keeper) AutoIncrementID(ctx sdk.Context, idKey []byte) uint64 {
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

func (k Keeper) ClearAutoIncrementID(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyLastOutgoingBatchID)
	store.Delete(types.KeyLastTxPoolID)
	store.Delete(types.LastSlashedBatchBlock)
	store.Delete(types.LastSlashedRelayerSetNonce)
	return
}

// GetOutgoingPendingTxTotal returns the total amount of a given token pending in the outgoing pool and all batches
func (k Keeper) GetOutgoingPendingTxTotal(ctx sdk.Context, chainName string, bridgeToken *types.BridgeToken) sdk.Int {
	totalPending := sdk.ZeroInt()
	// Add all unbatched transactions
	k.IterateUnbatchedTransactions(ctx, bridgeToken.ContractAddress, func(tx *types.OutgoingTransferTx) bool {
		totalPending = totalPending.Add(types.GetMintAmount(tx.Token.Amount, chainName, bridgeToken))
		totalPending = totalPending.Add(types.GetMintAmount(tx.Fee.Amount, chainName, bridgeToken))
		return false
	})
	// Add all batched transactions
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		if batch.TokenContract == bridgeToken.ContractAddress {
			totalPending = totalPending.Add(types.GetMintAmount(batch.TotalAmount(), chainName, bridgeToken))
		}
		return false
	})
	return totalPending
}
