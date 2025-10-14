package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Constructor functions for message types
const TypeMsgUpdateBlacklist = "update_blacklist"

var _ sdk.Msg = &MsgUpdateBlacklist{}

// NewMsgUpdateBlacklist creates a new MsgUpdateBlacklist instance
func NewMsgUpdateBlacklist(
	creator string,
	addressesToRemove []string,
	addressesToAdd []string,
) *MsgUpdateBlacklist {
	return &MsgUpdateBlacklist{
		Creator:           creator,
		AddressesToRemove: addressesToRemove,
		AddressesToAdd:    addressesToAdd,
	}
}

func (msg *MsgUpdateBlacklist) Route() string {
	return RouterKey
}

func (msg *MsgUpdateBlacklist) Type() string {
	return TypeMsgUpdateBlacklist
}

func (msg *MsgUpdateBlacklist) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateBlacklist) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic runs stateless checks on the message
func (msg *MsgUpdateBlacklist) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	for _, addr := range msg.AddressesToRemove {
		_, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address to remove (%s)", err)
		}
	}
	for _, addr := range msg.AddressesToAdd {
		_, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address to add (%s)", err)
		}
	}
	return nil
}
