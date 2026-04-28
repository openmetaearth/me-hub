package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/utils"

	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

const TypeMsgIbcTransferFromRegionTreasure = "ibc_transfer_from_region_treasure"

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

func (msg *MsgIbcTransferFromRegionTreasure) Route() string {
	return RouterKey
}

func (msg *MsgIbcTransferFromRegionTreasure) Type() string {
	return TypeMsgIbcTransferFromRegionTreasure
}

func (msg *MsgIbcTransferFromRegionTreasure) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgIbcTransferFromRegionTreasure) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgIbcTransferFromRegionTreasure) ValidateBasic() error {
	if err := host.PortIdentifierValidator(msg.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(msg.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	if _, err := utils.CheckRegionName(strings.ToUpper(msg.RegionId)); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidType, err.Error())
	}
	if !msg.Token.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Token.String())
	}
	if !msg.Token.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, msg.Token.String())
	}
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
