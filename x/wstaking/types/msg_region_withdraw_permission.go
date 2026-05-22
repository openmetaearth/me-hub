package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgGrantRegionWithdrawPermission  = "grant_region_withdraw_permission"
	TypeMsgRevokeRegionWithdrawPermission = "revoke_region_withdraw_permission"
)

var (
	_ sdk.Msg = &MsgGrantRegionWithdrawPermission{}
	_ sdk.Msg = &MsgRevokeRegionWithdrawPermission{}
)

func NewMsgGrantRegionWithdrawPermission(creator, regionId, address string) *MsgGrantRegionWithdrawPermission {
	return &MsgGrantRegionWithdrawPermission{
		Creator:  creator,
		RegionId: regionId,
		Address:  address,
	}
}

func (msg *MsgGrantRegionWithdrawPermission) Route() string { return RouterKey }

func (msg *MsgGrantRegionWithdrawPermission) Type() string {
	return TypeMsgGrantRegionWithdrawPermission
}

func (msg *MsgGrantRegionWithdrawPermission) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgGrantRegionWithdrawPermission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgGrantRegionWithdrawPermission) ValidateBasic() error {
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

func NewMsgRevokeRegionWithdrawPermission(creator, regionId string) *MsgRevokeRegionWithdrawPermission {
	return &MsgRevokeRegionWithdrawPermission{
		Creator:  creator,
		RegionId: regionId,
	}
}

func (msg *MsgRevokeRegionWithdrawPermission) Route() string { return RouterKey }

func (msg *MsgRevokeRegionWithdrawPermission) Type() string {
	return TypeMsgRevokeRegionWithdrawPermission
}

func (msg *MsgRevokeRegionWithdrawPermission) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRevokeRegionWithdrawPermission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRevokeRegionWithdrawPermission) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}
	if msg.RegionId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "region_id cannot be empty")
	}
	return nil
}
