package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewFixedDepositCfg = "add_fixed_deposit_cfg"

var _ sdk.Msg = &MsgNewFixedDepositCfg{}

func NewMsgNewFixedDepositCfg(dao string, regionId string, term int64, rate sdk.Dec) *MsgNewFixedDepositCfg {
	return &MsgNewFixedDepositCfg{
		Dao:      dao,
		RegionId: regionId,
		Term:     term,
		Rate:     rate,
	}
}

func (msg *MsgNewFixedDepositCfg) Route() string {
	return RouterKey
}

func (msg *MsgNewFixedDepositCfg) Type() string {
	return TypeMsgNewFixedDepositCfg
}

func (msg *MsgNewFixedDepositCfg) GetSigners() []sdk.AccAddress {
	admin, err := sdk.AccAddressFromBech32(msg.Dao)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{admin}
}

func (msg *MsgNewFixedDepositCfg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgNewFixedDepositCfg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Dao)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address (%s)", err)
	}
	return nil
}

const TypeMsgRemoveFixedDepositCfg = "remove_fixed_deposit_cfg"

var _ sdk.Msg = &MsgRemoveFixedDepositCfg{}

func NewMsgRemoveFixedDepositCfg(admin string, regionId string, term int64) *MsgRemoveFixedDepositCfg {
	return &MsgRemoveFixedDepositCfg{
		Admin:    admin,
		RegionId: regionId,
		Term:     term,
	}
}

func (msg *MsgRemoveFixedDepositCfg) Route() string {
	return RouterKey
}

func (msg *MsgRemoveFixedDepositCfg) Type() string {
	return TypeMsgRemoveFixedDepositCfg
}

func (msg *MsgRemoveFixedDepositCfg) GetSigners() []sdk.AccAddress {
	admin, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{admin}
}

func (msg *MsgRemoveFixedDepositCfg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveFixedDepositCfg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address (%s)", err)
	}
	return nil
}

const TypeMsgSetFixedDepositCfgRate = "set_fixed_deposit_cfg_rate"

var _ sdk.Msg = &MsgSetFixedDepositCfgRate{}

func NewMsgSetFixedDepositCfgRate(admin string, regionId string, term int64, rate sdk.Dec) *MsgSetFixedDepositCfgRate {
	return &MsgSetFixedDepositCfgRate{
		Admin:    admin,
		RegionId: regionId,
		Term:     term,
		Rate:     rate,
	}
}

func (msg *MsgSetFixedDepositCfgRate) Route() string {
	return RouterKey
}

func (msg *MsgSetFixedDepositCfgRate) Type() string {
	return TypeMsgSetFixedDepositCfgRate
}

func (msg *MsgSetFixedDepositCfgRate) GetSigners() []sdk.AccAddress {
	admin, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{admin}
}

func (msg *MsgSetFixedDepositCfgRate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetFixedDepositCfgRate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address (%s)", err)
	}
	return nil
}

const TypeMsgSetFixedDepositCfgStatus = "set_fixed_deposit_cfg_status"

var _ sdk.Msg = &MsgSetFixedDepositCfgStatus{}

func NewMsgSetFixedDepositCfgStatus(admin string, regionId string, term int64, status FIXED_DEPOSIT_CFG_STATUS) *MsgSetFixedDepositCfgStatus {
	return &MsgSetFixedDepositCfgStatus{
		Admin:    admin,
		RegionId: regionId,
		Term:     term,
		Status:   status,
	}
}

func (msg *MsgSetFixedDepositCfgStatus) Route() string {
	return RouterKey
}

func (msg *MsgSetFixedDepositCfgStatus) Type() string {
	return TypeMsgSetFixedDepositCfgStatus
}

func (msg *MsgSetFixedDepositCfgStatus) GetSigners() []sdk.AccAddress {
	admin, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{admin}
}

func (msg *MsgSetFixedDepositCfgStatus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetFixedDepositCfgStatus) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address (%s)", err)
	}
	return nil
}
