package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSendToModule = "send_to_module"

var _ sdk.Msg = &MsgSendToModule{}

func NewMsgSendToModule(account, receiver string, amount sdk.Coins) *MsgSendToModule {
	return &MsgSendToModule{
		Sender:   account,
		Receiver: receiver,
		Amount:   amount,
	}
}

func (msg *MsgSendToModule) Route() string {
	return RouterKey
}

func (msg *MsgSendToModule) Type() string {
	return TypeMsgSendToModule
}

func (msg *MsgSendToModule) GetSigners() []sdk.AccAddress {
	account, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{account}
}

func (msg *MsgSendToModule) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendToModule) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address (%s)", err)
	}
	if msg.Receiver == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "receiver address cannot be empty")
	}
	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}
	return nil
}
