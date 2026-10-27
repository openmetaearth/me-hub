package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDoFixedDeposit = "do_fixed_deposit"

var _ sdk.Msg = &MsgDoFixedDeposit{}

func NewMsgDoFixedDeposit(account string, principal sdk.Coin, term int64) *MsgDoFixedDeposit {
	return &MsgDoFixedDeposit{
		Account:   account,
		Principal: principal,
		Term:      term,
	}
}

func (msg *MsgDoFixedDeposit) Route() string {
	return RouterKey
}

func (msg *MsgDoFixedDeposit) Type() string {
	return TypeMsgDoFixedDeposit
}

func (msg *MsgDoFixedDeposit) GetSigners() []sdk.AccAddress {
	account, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{account}
}

func (msg *MsgDoFixedDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDoFixedDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address (%s)", err)
	}
	return nil
}

const TypeMsgWithdrawFixedDeposit = "do_fixed_withdraw"

var _ sdk.Msg = &MsgWithdrawFixedDeposit{}

func NewMsgWithdrawFixedDeposit(account string, id uint64) *MsgWithdrawFixedDeposit {
	return &MsgWithdrawFixedDeposit{
		Account: account,
		Id:      id,
	}
}

func (msg *MsgWithdrawFixedDeposit) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawFixedDeposit) Type() string {
	return TypeMsgWithdrawFixedDeposit
}

func (msg *MsgWithdrawFixedDeposit) GetSigners() []sdk.AccAddress {
	account, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{account}
}

func (msg *MsgWithdrawFixedDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawFixedDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address (%s)", err)
	}
	return nil
}



const TypeMsgBondDepositToPromptSct = "bond_deposit_to_prompt_sct"

var _ sdk.Msg = &MsgBondDepositToPromptSctRequest{}

func NewMsgBondDepositToPromptSctRequest(creator string, contractAddress string, id uint64) *MsgBondDepositToPromptSctRequest {
	return &MsgBondDepositToPromptSctRequest{
		Creator:          creator,
		Address: contractAddress,
		Id:               id,
	}
}

func (msg *MsgBondDepositToPromptSctRequest) Route() string {
	return RouterKey
}

func (msg *MsgBondDepositToPromptSctRequest) Type() string {
	return TypeMsgBondDepositToPromptSct
}

func (msg *MsgBondDepositToPromptSctRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBondDepositToPromptSctRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBondDepositToPromptSctRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

	

