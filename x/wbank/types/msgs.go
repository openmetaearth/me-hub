package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// bank message types
const (
	TypeMsgWithdrawTreasury = "MsgWithdrawTreasury"
	TypeMsgSendToAirdrop    = "MsgSendToAirdrop"
	TypeMsgSendToTreasury   = "MsgSendToTreasury"
)

var (
	_ sdk.Msg = &MsgWithdrawTreasury{}
	_ sdk.Msg = &MsgSendToAirdrop{}
	_ sdk.Msg = &MsgSendToTreasury{}
)

// MsgSendToGlobalDao - construct a msg to send coins from treasury to global admin.
//
//nolint:interfacer
func NewMsgWithdrawTreasury(fromAddr, receiver sdk.AccAddress, amount sdk.Coins) *MsgWithdrawTreasury {
	return &MsgWithdrawTreasury{
		FromAddress: fromAddr.String(),
		Receiver:    receiver.String(),
		Amount:      amount,
	}
}

// Route Implements Msg.
func (msg MsgWithdrawTreasury) Route() string { return banktypes.RouterKey }

// Type Implements Msg.
func (msg MsgWithdrawTreasury) Type() string { return TypeMsgWithdrawTreasury }

// ValidateBasic Implements Msg.
func (msg MsgWithdrawTreasury) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.FromAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Receiver); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgWithdrawTreasury) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgWithdrawTreasury) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}

// NewMsgSendToAirdrop - construct a msg to send coins from region treasury to airdrop address.
//
//nolint:interfacer
func NewMsgSendToAirdrop(adminAddr sdk.AccAddress, regionID string, amount sdk.Coins) *MsgSendToAirdrop {
	return &MsgSendToAirdrop{FromAddress: adminAddr.String(), RegionId: regionID, Amount: amount}
}

// Route Implements Msg.
func (msg MsgSendToAirdrop) Route() string { return banktypes.RouterKey }

// Type Implements Msg.
func (msg MsgSendToAirdrop) Type() string { return TypeMsgSendToAirdrop }

// ValidateBasic Implements Msg.
func (msg MsgSendToAirdrop) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.FromAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid admin address: %s", err)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSendToAirdrop) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgSendToAirdrop) GetSigners() []sdk.AccAddress {
	adminAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{adminAddress}
}

// NewMsgSendToTreasury - construct a msg to send coins from global admin to treasury.
//
//nolint:interfacer
func NewMsgSendToTreasury(fromAddr sdk.AccAddress, amount sdk.Coins) *MsgSendToTreasury {
	return &MsgSendToTreasury{FromAddress: fromAddr.String(), Amount: amount}
}

// Route Implements Msg.
func (msg MsgSendToTreasury) Route() string { return banktypes.RouterKey }

// Type Implements Msg.
func (msg MsgSendToTreasury) Type() string { return TypeMsgSendToTreasury }

// ValidateBasic Implements Msg.
func (msg MsgSendToTreasury) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.FromAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSendToTreasury) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgSendToTreasury) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}
