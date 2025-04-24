package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateGroup = "create_group"
	TypeMsgUpdateGroup = "update_group"
	TypeMsgDeleteGroup = "delete_group"
)

var _ sdk.Msg = &MsgCreateGroup{}

func NewMsgCreateGroup(creator string, groupInfo *GroupInfo) *MsgCreateGroup {
	return &MsgCreateGroup{
		Creator:   creator,
		GroupInfo: groupInfo,
	}
}

func (msg *MsgCreateGroup) Route() string {
	return RouterKey
}

func (msg *MsgCreateGroup) Type() string {
	return TypeMsgCreateGroup
}

func (msg *MsgCreateGroup) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateGroup) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateGroup) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.GroupInfo.Admin)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid Group Admin address (%s)", err)
	}
	if "" == msg.GroupInfo.RegionID {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "RegionID can not be empty")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateGroup{}

func NewMsgUpdateGroup(creator string, id uint64, groupInfo *GroupInfo) *MsgUpdateGroup {
	return &MsgUpdateGroup{
		Id:        id,
		Creator:   creator,
		GroupInfo: groupInfo,
	}
}

func (msg *MsgUpdateGroup) Route() string {
	return RouterKey
}

func (msg *MsgUpdateGroup) Type() string {
	return TypeMsgUpdateGroup
}

func (msg *MsgUpdateGroup) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateGroup) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateGroup) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteGroup{}

func NewMsgDeleteGroup(creator string, id uint64) *MsgDeleteGroup {
	return &MsgDeleteGroup{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteGroup) Route() string {
	return RouterKey
}

func (msg *MsgDeleteGroup) Type() string {
	return TypeMsgDeleteGroup
}

func (msg *MsgDeleteGroup) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteGroup) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteGroup) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
