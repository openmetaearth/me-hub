package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgNewFixedDepositCfg{}

func NewMsgNewFixedDepositCfg(dao string, regionId string, term int64, rate sdkmath.LegacyDec) *MsgNewFixedDepositCfg {
	return &MsgNewFixedDepositCfg{
		Dao:      dao,
		RegionId: regionId,
		Term:     term,
		Rate:     rate,
	}
}

func (msg *MsgNewFixedDepositCfg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Dao)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgRemoveFixedDepositCfg{}

func NewMsgRemoveFixedDepositCfg(admin string, regionId string, term int64) *MsgRemoveFixedDepositCfg {
	return &MsgRemoveFixedDepositCfg{
		Admin:    admin,
		RegionId: regionId,
		Term:     term,
	}
}

func (msg *MsgRemoveFixedDepositCfg) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgSetFixedDepositCfgRate{}

func NewMsgSetFixedDepositCfgRate(admin string, regionId string, term int64, rate sdkmath.LegacyDec) *MsgSetFixedDepositCfgRate {
	return &MsgSetFixedDepositCfgRate{
		Admin:    admin,
		RegionId: regionId,
		Term:     term,
		Rate:     rate,
	}
}

func (msg *MsgSetFixedDepositCfgRate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgSetFixedDepositCfgStatus{}

func NewMsgSetFixedDepositCfgStatus(admin string, regionId string, term int64, status FIXED_DEPOSIT_CFG_STATUS) *MsgSetFixedDepositCfgStatus {
	return &MsgSetFixedDepositCfgStatus{
		Admin:    admin,
		RegionId: regionId,
		Term:     term,
		Status:   status,
	}
}

func (msg *MsgSetFixedDepositCfgStatus) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address (%s)", err)
	}
	return nil
}
