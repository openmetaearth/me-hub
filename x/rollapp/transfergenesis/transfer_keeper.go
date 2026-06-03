package transfergenesis

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/openmetaearth/me-hub/utils/gerrc"
)

type PacketForwardTransferKeeper interface {
	Transfer(ctx context.Context, msg *transferTypes.MsgTransfer) (*transferTypes.MsgTransferResponse, error)
	DenomPathFromHash(ctx sdk.Context, denom string) (string, error)
	GetTotalEscrowForDenom(ctx sdk.Context, denom string) sdk.Coin
	SetTotalEscrowForDenom(ctx sdk.Context, coin sdk.Coin)
}

type transferEnabledTransferKeeper struct {
	next                  PacketForwardTransferKeeper
	getRollapp            GetRollapp
	getChannelClientState ChannelKeeper
}

func NewTransferEnabledTransferKeeper(
	next PacketForwardTransferKeeper,
	getRollapp GetRollapp,
	getChannelClientState ChannelKeeper,
) PacketForwardTransferKeeper {
	return transferEnabledTransferKeeper{
		next:                  next,
		getRollapp:            getRollapp,
		getChannelClientState: getChannelClientState,
	}
}

func (k transferEnabledTransferKeeper) Transfer(
	goCtx context.Context,
	msg *transferTypes.MsgTransfer,
) (*transferTypes.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	gate := TransferEnabledDecorator{
		getRollapp:            k.getRollapp,
		getChannelClientState: k.getChannelClientState,
	}

	enabled, err := gate.transfersEnabled(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(err, "transfer genesis: transfers enabled")
	}
	if !enabled {
		return nil, errorsmod.Wrap(gerrc.ErrFailedPrecondition, "transfers to/from rollapp are disabled")
	}

	return k.next.Transfer(goCtx, msg)
}

func (k transferEnabledTransferKeeper) DenomPathFromHash(ctx sdk.Context, denom string) (string, error) {
	return k.next.DenomPathFromHash(ctx, denom)
}

func (k transferEnabledTransferKeeper) GetTotalEscrowForDenom(ctx sdk.Context, denom string) sdk.Coin {
	return k.next.GetTotalEscrowForDenom(ctx, denom)
}

func (k transferEnabledTransferKeeper) SetTotalEscrowForDenom(ctx sdk.Context, coin sdk.Coin) {
	k.next.SetTotalEscrowForDenom(ctx, coin)
}
