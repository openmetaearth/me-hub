package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// ZkContractDIDDecorator restricts direct calls to DID-configured zk contracts
// to senders whose DID exists and is active.
type ZkContractDIDDecorator struct {
	didKeeper DidKeeper
}

var _ sdk.AnteDecorator = ZkContractDIDDecorator{}

func NewZkContractDIDDecorator(didKeeper DidKeeper) ZkContractDIDDecorator {
	return ZkContractDIDDecorator{didKeeper: didKeeper}
}

func (d ZkContractDIDDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	for _, sdkMsg := range tx.GetMsgs() {
		msg, ok := sdkMsg.(*evmtypes.MsgEthereumTx)
		if !ok {
			continue
		}

		ethTx := msg.AsTransaction()
		if ethTx == nil {
			return ctx, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to decode Ethereum transaction")
		}
		target := ethTx.To()
		if target == nil || !d.didKeeper.IsZkContractAddress(ctx, *target) {
			continue
		}

		sender := msg.GetFrom()
		if sender.Empty() {
			return ctx, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Ethereum transaction sender is unavailable")
		}
		isAllowed, err := d.didKeeper.IsAllowToUseZkContractVerify(ctx, sender)
		if err != nil {
			return ctx, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "failed to check if sender is allowed to use zk contract verify," +
			" err = %s,caller = %s", err.Error(), sender.String())
		}
		if !isAllowed {
			return ctx, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "sender is not allowed to use zk contract verify.caller = %s", sender.String())
		}
		
	}

	return next(ctx, tx, simulate)
}
