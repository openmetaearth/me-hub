package bridgingfee

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	transferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	commontypes "github.com/st-chain/me-hub/x/common/types"
	delayedackkeeper "github.com/st-chain/me-hub/x/delayedack/keeper"
	rollappkeeper "github.com/st-chain/me-hub/x/rollapp/keeper"
	rollapptypes "github.com/st-chain/me-hub/x/rollapp/types"
	"strings"
)

const (
	ModuleName = "bridging_fee"
)

// IBCModule is responsible for charging a bridging fee on transfers coming from rollapps
// The actual charge happens on the packet finalization
// based on ADR: https://www.notion.so/dymension/ADR-x-Bridging-Fee-Middleware-7ba8c191373f43ce81782fc759913299?pvs=4
type IBCModule struct {
	ibctransfer.IBCModule

	rollappKeeper    rollappkeeper.Keeper
	delayedAckKeeper delayedackkeeper.Keeper
	transferKeeper   transferkeeper.Keeper
	feeModuleAddr    sdk.AccAddress
	bankKeeper       bankkeeper.Keeper
}

func NewIBCModule(
	next ibctransfer.IBCModule,
	keeper delayedackkeeper.Keeper,
	transferKeeper transferkeeper.Keeper,
	feeModuleAddr sdk.AccAddress,
	rollappKeeper rollappkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) *IBCModule {
	return &IBCModule{
		IBCModule:        next,
		delayedAckKeeper: keeper,
		transferKeeper:   transferKeeper,
		feeModuleAddr:    feeModuleAddr,
		rollappKeeper:    rollappKeeper,
		bankKeeper:       bankKeeper,
	}
}

func (w IBCModule) logger(
	ctx sdk.Context,
	packet channeltypes.Packet,
	method string,
) log.Logger {
	return ctx.Logger().With(
		"module", ModuleName,
		"packet_source_port", packet.SourcePort,
		"packet_destination_port", packet.DestinationPort,
		"packet_sequence", packet.Sequence,
		"method", method,
	)
}

func (w *IBCModule) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	l := w.logger(ctx, packet, "OnRecvPacket")

	if commontypes.SkipRollappMiddleware(ctx) || !w.delayedAckKeeper.IsRollappsEnabled(ctx) {
		return w.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	transfer, err := w.rollappKeeper.GetValidTransfer(ctx, packet.GetData(), packet.GetDestPort(), packet.GetDestChannel())
	if err != nil {
		l.Error("Get valid transfer.", "err", err)
		err = errorsmod.Wrapf(err, "%s: get valid transfer", ModuleName)
		return channeltypes.NewErrorAcknowledgement(err)
	}

	// check if the token is a returned native token
	isReturnedNative, originalDenom, escrowAddress := w.checkIfReturnedNativeToken(ctx, transfer, packet)
	if isReturnedNative {
		err := w.unlockEscrowAndTransfer(ctx, escrowAddress, transfer, originalDenom, packet)
		if err != nil {
			l.Error("Unlock escrow failed.", "err", err)
			return channeltypes.NewErrorAcknowledgement(err)
		}
		return channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	}

	if !transfer.IsRollapp() {
		return w.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	// Use the packet as a basis for a fee transfer
	feeData := transfer
	fee := w.delayedAckKeeper.BridgingFeeFromAmt(ctx, transfer.MustAmountInt())
	if fee.IsZero() {
		return w.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}
	feeData.Amount = fee.String()
	feeData.Receiver = w.feeModuleAddr.String()

	// No event emitted, as we called the transfer keeper directly (vs the transfer middleware)
	err = w.transferKeeper.OnRecvPacket(ctx, packet, feeData.FungibleTokenPacketData)
	if err != nil {
		l.Error("Charge bridging fee.", "err", err)
		// we continue as we don't want the fee charge to fail the transfer in any case
		fee = sdk.ZeroInt()
	} else {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				EventTypeBridgingFee,
				sdk.NewAttribute(AttributeKeyFee, fee.String()),
				sdk.NewAttribute(sdk.AttributeKeySender, transfer.Sender),
				sdk.NewAttribute(transfertypes.AttributeKeyReceiver, transfer.Receiver),
				sdk.NewAttribute(transfertypes.AttributeKeyDenom, transfer.Denom),
				sdk.NewAttribute(transfertypes.AttributeKeyAmount, transfer.Amount),
			),
		)
	}

	// transfer the rest to the original recipient
	transfer.Amount = transfer.MustAmountInt().Sub(fee).String()
	packet.Data = transfer.GetBytes()
	return w.IBCModule.OnRecvPacket(ctx, packet, relayer)
}

// checkIfReturnedNativeToken checks whether the received denom is a "returned native token" that was originally from this chain
func (w *IBCModule) checkIfReturnedNativeToken(
	ctx sdk.Context,
	transfer rollapptypes.TransferData,
	packet channeltypes.Packet,
) (bool, string, sdk.AccAddress) {
	denomTrace := transfertypes.ParseDenomTrace(transfer.Denom)

	if denomTrace.Path == "" {
		return false, "", nil
	}

	baseDenom := denomTrace.BaseDenom
	paths := denomTrace.GetPath()

	if len(paths) < 2 {
		return false, "", nil
	}

	pathsSplit := strings.Split(paths, "/")
	originalSourceChannel := ""
	for i := len(pathsSplit) - 2; i >= 0; i -= 2 {
		if pathsSplit[i] == "transfer" && i+1 < len(paths) {
			originalSourceChannel = pathsSplit[i+1]
			break
		}
	}

	if originalSourceChannel == "" {
		return false, "", nil
	}

	escrowAddress := transfertypes.GetEscrowAddress(transfertypes.PortID, originalSourceChannel)
	escrowBalance := w.bankKeeper.GetBalance(ctx, escrowAddress, baseDenom)

	if escrowBalance.IsZero() {
		return false, "", nil
	}

	return true, baseDenom, escrowAddress
}

func (w *IBCModule) unlockEscrowAndTransfer(
	ctx sdk.Context,
	escrowAddress sdk.AccAddress,
	transfer rollapptypes.TransferData,
	originalDenom string,
	packet channeltypes.Packet,
) error {
	sourceChannel := packet.GetSourceChannel()
	receiver, err := sdk.AccAddressFromBech32(transfer.Receiver)
	if err != nil {
		return errorsmod.Wrap(err, "invalid receiver address")
	}

	amount, ok := sdk.NewIntFromString(transfer.Amount)
	if !ok {
		return errorsmod.Wrap(transfertypes.ErrInvalidAmount, transfer.Amount)
	}

	fee := w.delayedAckKeeper.BridgingFeeFromAmt(ctx, transfer.MustAmountInt())
	token := sdk.NewCoin(originalDenom, amount.Sub(fee))

	w.chargeBridgingFee(ctx, transfer, packet, fee)

	err = w.bankKeeper.SendCoins(ctx, escrowAddress, receiver, sdk.NewCoins(token))
	if err != nil {
		return errorsmod.Wrap(err, "escrow unlock failed")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"unlock_escrow",
			sdk.NewAttribute("receiver", transfer.Receiver),
			sdk.NewAttribute("denom", originalDenom),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("channel", sourceChannel),
		),
	)
	return nil
}

func (w *IBCModule) chargeBridgingFee(
	ctx sdk.Context,
	transfer rollapptypes.TransferData,
	packet channeltypes.Packet,
	fee sdk.Int,
) {
	feeData := transfer
	feeData.Amount = fee.String()
	feeData.Receiver = w.feeModuleAddr.String()

	err := w.transferKeeper.OnRecvPacket(ctx, packet, feeData.FungibleTokenPacketData)
	if err != nil {
		return
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeBridgingFee,
			sdk.NewAttribute(AttributeKeyFee, fee.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, feeData.Sender),
			sdk.NewAttribute(transfertypes.AttributeKeyReceiver, feeData.Receiver),
			sdk.NewAttribute(transfertypes.AttributeKeyDenom, feeData.Denom),
			sdk.NewAttribute(transfertypes.AttributeKeyAmount, feeData.Amount),
		),
	)
}
