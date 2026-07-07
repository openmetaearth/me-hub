package types

import (
	"strings"

	"cosmossdk.io/errors"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/utils"

	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
)

var _ sdk.Msg = &MsgIbcTransferFromRegionTreasure{}

func NewMsgIbcTransferFromRegionTreasure(sourcePort, sourceChannel, regionId string, amount sdk.Coin, timeoutHeight Height, timeoutTimestamp uint64, momo, creator string) *MsgIbcTransferFromRegionTreasure {
	return &MsgIbcTransferFromRegionTreasure{
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		RegionId:         regionId,
		Token:            amount,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
		Memo:             momo,
		Creator:          creator,
	}
}

func (msg *MsgIbcTransferFromRegionTreasure) ValidateBasic() error {
	if err := host.PortIdentifierValidator(msg.SourcePort); err != nil {
		return errorsmod.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(msg.SourceChannel); err != nil {
		return errorsmod.Wrap(err, "invalid source channel ID")
	}
	if _, err := utils.CheckRegionName(strings.ToUpper(msg.RegionId)); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidType, err.Error())
	}
	if !msg.Token.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, msg.Token.String())
	}
	if !msg.Token.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, msg.Token.String())
	}
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
