package ante

import (
	"bytes"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

// submitResponseSelector identifies IVerifier.submitResponse(AuthResponse, Response[], bytes).
var submitResponseSelector = crypto.Keccak256([]byte("submitResponse((string,bytes),(uint256,bytes,bytes)[],bytes)"))[:4]

// ZkContractDIDDecorator checks sender permission for direct submitResponse calls
// to the DID-configured universal verifier.
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
		if target == nil {
			continue
		}

		zkContractInfo := d.didKeeper.GetZkContractAddress(ctx, didtypes.PROXY_UNIVERSAL_VERIFIER)
		if zkContractInfo == nil {
			continue
		}
		if !common.IsHexAddress(zkContractInfo.ContractAddress) {
			return ctx, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "zk contract address is invalid,contractType = %s,contractAddress = %s",
				didtypes.PROXY_UNIVERSAL_VERIFIER.String(), zkContractInfo.ContractAddress)
		}
		if *target != common.HexToAddress(zkContractInfo.ContractAddress) {
			continue
		}

		data := ethTx.Data()
		if len(data) < len(submitResponseSelector) || !bytes.Equal(data[:len(submitResponseSelector)], submitResponseSelector) {
			continue
		}

		sender := msg.GetFrom()
		if sender.Empty() {
			return ctx, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Ethereum transaction sender is unavailable")
		}

		isAllowed, err := d.didKeeper.IsAllowToUseZkContractVerify(ctx, sender)
		if err != nil {
			return ctx, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "failed to check if sender is allowed to use zk contract verify,"+
				" err = %s,caller = %s", err.Error(), sender.String())
		}
		if !isAllowed {
			return ctx, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "sender is not allowed to use zk contract verify.caller = %s", sender.String())
		}
	}

	return next(ctx, tx, simulate)
}
