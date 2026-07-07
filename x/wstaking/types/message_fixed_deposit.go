package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgDoFixedDeposit{}

func NewMsgDoFixedDeposit(account string, principal sdk.Coin, term int64) *MsgDoFixedDeposit {
	return &MsgDoFixedDeposit{
		Account:   account,
		Principal: principal,
		Term:      term,
	}
}

func (msg *MsgDoFixedDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgWithdrawFixedDeposit{}

func NewMsgWithdrawFixedDeposit(account string, id uint64) *MsgWithdrawFixedDeposit {
	return &MsgWithdrawFixedDeposit{
		Account: account,
		Id:      id,
	}
}

func (msg *MsgWithdrawFixedDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address (%s)", err)
	}
	return nil
}
