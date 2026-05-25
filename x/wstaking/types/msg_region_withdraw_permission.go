package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgGrantRegionWithdraw  = "grant_region_withdraw"
	TypeMsgRevokeRegionWithdraw = "revoke_region_withdraw"
)

var (
	_ sdk.Msg = &MsgGrantRegionWithdraw{}
	_ sdk.Msg = &MsgRevokeRegionWithdraw{}
)

func NewMsgGrantRegionWithdraw(creator, regionId, address string) *MsgGrantRegionWithdraw {
	return &MsgGrantRegionWithdraw{
		Creator:  creator,
		RegionId: regionId,
		Address:  address,
	}
}

func (msg *MsgGrantRegionWithdraw) Route() string { return RouterKey }

func (msg *MsgGrantRegionWithdraw) Type() string {
	return TypeMsgGrantRegionWithdraw
}

func (msg *MsgGrantRegionWithdraw) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgGrantRegionWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgGrantRegionWithdraw) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid granted address: %s", err)
	}
	if msg.RegionId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "region_id cannot be empty")
	}
	return nil
}

func NewMsgRevokeRegionWithdraw(creator, regionId string) *MsgRevokeRegionWithdraw {
	return &MsgRevokeRegionWithdraw{
		Creator:  creator,
		RegionId: regionId,
	}
}

func (msg *MsgRevokeRegionWithdraw) Route() string { return RouterKey }

func (msg *MsgRevokeRegionWithdraw) Type() string {
	return TypeMsgRevokeRegionWithdraw
}

func (msg *MsgRevokeRegionWithdraw) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRevokeRegionWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRevokeRegionWithdraw) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}
	if msg.RegionId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "region_id cannot be empty")
	}
	return nil
}
