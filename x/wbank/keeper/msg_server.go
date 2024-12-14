package keeper

import (
	"context"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/st-chain/me-hub/x/wbank/types"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
)

// MsgServer is wrapper staking customParamsKeeper message server.
type MsgServer struct {
	banktypes.MsgServer
	BaseKeeperWrapper
}

var _ types.MsgServer = MsgServer{}

// NewMsgServerImpl returns an implementation of the staking wrapped MsgServer.
func NewMsgServerImpl(
	keeper BaseKeeperWrapper,
	msgSrv banktypes.MsgServer,
) MsgServer {
	return MsgServer{
		BaseKeeperWrapper: keeper,
		MsgServer:         msgSrv,
	}
}

func (k MsgServer) WithdrawTreasury(goCtx context.Context, msg *types.MsgWithdrawTreasury) (*types.MsgWithdrawTreasuryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	if !k.dk.IsGlobalDao(ctx, msg.FromAddress) {
		return nil, types.ErrNotGlobalDao
	}

	if k.BlockedAddr(receiver) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.FromAddress)
	}

	err = k.SendCoinsFromModuleToAccount(ctx, types.TreasuryPoolName, receiver, msg.Amount)
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range msg.Amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "sendToAdmin"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSendToGlobalDao,
			sdk.NewAttribute(banktypes.AttributeKeySender, msg.FromAddress),
			sdk.NewAttribute(banktypes.AttributeKeyReceiver, msg.Receiver),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	)
	return &types.MsgWithdrawTreasuryResponse{}, nil
}

func (k MsgServer) SendToAirdrop(goCtx context.Context, msg *types.MsgSendToAirdrop) (*types.MsgSendToAirdropResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, err
	}

	if !k.dk.IsGlobalDao(ctx, msg.FromAddress) {
		return nil, types.ErrNotGlobalDao
	}

	airdropAddress := k.dk.GetAirdropAddress(ctx)
	to, err := sdk.AccAddressFromBech32(airdropAddress)
	if err != nil {
		return nil, err
	}

	if k.BlockedAddr(to) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", to)
	}

	regionBaseAccount := wstakingtypes.GetRegionAccountAddr(wstakingtypes.RegionAccountTypeBase, msg.RegionId)
	err = k.SendCoins(ctx, regionBaseAccount, to, msg.Amount)
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range msg.Amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "sendToAirdrop"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeSendToAirdrop,
			sdk.NewAttribute(types.AttributeKeyRegionId, msg.RegionId),
			sdk.NewAttribute(sdk.AttributeKeySender, regionBaseAccount.String()),
			sdk.NewAttribute(banktypes.AttributeKeyReceiver, to.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	)
	return &types.MsgSendToAirdropResponse{}, nil
}

func (k MsgServer) SendToTreasury(goCtx context.Context, msg *types.MsgSendToTreasury) (*types.MsgSendToTreasuryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, err
	}

	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}

	if !k.dk.IsGlobalDao(ctx, msg.FromAddress) {
		return nil, types.ErrNotGlobalDao
	}

	err = k.SendCoinsFromAccountToModule(ctx, from, types.TreasuryPoolName, msg.Amount)
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range msg.Amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "sendToTreasury"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSendToTreasury,
			sdk.NewAttribute(banktypes.AttributeKeySender, msg.FromAddress),
			sdk.NewAttribute(banktypes.AttributeKeyReceiver, authtypes.NewModuleAddress(types.TreasuryPoolName).String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	)

	return &types.MsgSendToTreasuryResponse{}, nil
}
