package transfergenesis

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/openmetaearth/me-hub/utils/gerrc"
	"github.com/openmetaearth/me-hub/utils/uibc"

	transferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

type GetRollapp func(ctx sdk.Context, rollappId string) (val types.Rollapp, found bool)

// TransferEnabledDecorator only allows ibc transfers to a rollapp if that rollapp has finished
// the transfer genesis protocol.
type TransferEnabledDecorator struct {
	getRollapp            GetRollapp
	getChannelClientState ChannelKeeper
}

func NewTransferEnabledDecorator(getRollapp GetRollapp, getChannelClientState ChannelKeeper) *TransferEnabledDecorator {
	return &TransferEnabledDecorator{
		getRollapp:            getRollapp,
		getChannelClientState: getChannelClientState,
	}
}

func (h TransferEnabledDecorator) transfersEnabled(ctx sdk.Context, transfer *transferTypes.MsgTransfer) (bool, error) {
	chainID, err := uibc.ChainIDFromPortChannel(ctx, h.getChannelClientState, transfer.SourcePort, transfer.SourceChannel)
	if err != nil {
		return false, errorsmod.Wrap(err, "chain id from port channel")
	}
	ra, ok := h.getRollapp(ctx, chainID)
	if !ok {
		return true, nil
	}
	return ra.GenesisState.TransfersEnabled, nil
}

func (h TransferEnabledDecorator) validateMsgTransfersEnabled(ctx sdk.Context, msg sdk.Msg) error {
	switch msg := msg.(type) {
	case *transferTypes.MsgTransfer:
		return h.validateTransferEnabled(ctx, msg)
	case *authztypes.MsgExec:
		return h.validateAuthzExecTransfersEnabled(ctx, *msg)
	default:
		return nil
	}
}

func (h TransferEnabledDecorator) validateAuthzExecTransfersEnabled(ctx sdk.Context, msg authztypes.MsgExec) error {
	msgs, err := msg.GetMessages()
	if err != nil {
		return errorsmod.Wrap(err, "authz exec messages")
	}
	for i, innerMsg := range msgs {
		if err := h.validateMsgTransfersEnabled(ctx, innerMsg); err != nil {
			return errorsmod.Wrapf(err, "authz exec message %d", i)
		}
	}
	return nil
}

func (h TransferEnabledDecorator) validateTransferEnabled(ctx sdk.Context, msg *transferTypes.MsgTransfer) error {
	if msg == nil {
		return errorsmod.Wrap(gerrc.ErrUnknown, "nil transfer message")
	}
	enabled, err := h.transfersEnabled(ctx, msg)
	if err != nil {
		return errorsmod.Wrap(err, "transfer genesis: transfers enabled")
	}
	if !enabled {
		return errorsmod.Wrap(gerrc.ErrFailedPrecondition, "transfers to/from rollapp are disabled")
	}
	return nil
}

// AnteHandle will return an error if the tx contains an ibc transfer message to a rollapp that has not finished the transfer genesis protocol.
func (h TransferEnabledDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		if err := h.validateMsgTransfersEnabled(ctx, msg); err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}
