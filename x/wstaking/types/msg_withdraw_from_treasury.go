package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdrawFromTreasury = "withdraw_from_treasury"

var _ sdk.Msg = &MsgWithdrawFromTreasury{}

func NewMsgWithdrawFromTreasury(withdrawer, receiver string, Amount sdk.Coins) *MsgWithdrawFromTreasury {
	return &MsgWithdrawFromTreasury{
		Withdrawer: withdrawer,
		Receiver:   receiver,
		Amount:     Amount,
	}
}

func (msg *MsgWithdrawFromTreasury) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawFromTreasury) Type() string {
	return TypeMsgWithdrawFromTreasury
}

func (msg *MsgWithdrawFromTreasury) GetSigners() []sdk.AccAddress {
	creator, _ := sdk.AccAddressFromBech32(msg.Withdrawer)
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawFromTreasury) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawFromTreasury) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	return nil
}
