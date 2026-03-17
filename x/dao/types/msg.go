package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgUpdateDao{}
	_ sdk.Msg = &MsgFreeGasAccount{}
)

func NewMsgUpdateDao(creator sdk.AccAddress, addresses DaoAddresses) *MsgUpdateDao {
	return &MsgUpdateDao{
		Creator:      creator.String(),
		DaoAddresses: addresses,
	}
}

func (msg *MsgUpdateDao) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDao) Type() string {
	return "UpdateDao"
}

func (msg *MsgUpdateDao) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic("invalid creator address")
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateDao) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDao) ValidateBasic() error {
	if len(msg.DaoAddresses.GlobalDao) > 0 {
		if _, err := sdk.AccAddressFromBech32(msg.DaoAddresses.GlobalDao); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.DaoAddresses.GlobalDao)
		}
	}
	if len(msg.DaoAddresses.MeidDao) > 0 {
		if _, err := sdk.AccAddressFromBech32(msg.DaoAddresses.MeidDao); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.DaoAddresses.MeidDao)
		}
	}
	if len(msg.DaoAddresses.DevOperator) > 0 {
		if _, err := sdk.AccAddressFromBech32(msg.DaoAddresses.DevOperator); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.DaoAddresses.DevOperator)
		}
	}
	if len(msg.DaoAddresses.AirdropAddress) > 0 {
		if _, err := sdk.AccAddressFromBech32(msg.DaoAddresses.AirdropAddress); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.DaoAddresses.AirdropAddress)
		}
	}
	return nil
}

func NewMsgFreeGasAccount(creator sdk.AccAddress, accounts []FreeGasAccount) *MsgFreeGasAccount {
	return &MsgFreeGasAccount{
		Creator:  creator.String(),
		Accounts: accounts,
	}
}

func (msg *MsgFreeGasAccount) Route() string {
	return RouterKey
}

func (msg *MsgFreeGasAccount) Type() string {
	return "UpdateDao"
}

func (msg *MsgFreeGasAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic("invalid creator address")
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgFreeGasAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFreeGasAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Creator)
	}
	if len(msg.Accounts) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "addresses is empty")
	}
	if len(msg.Accounts) > 100 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "addresses is too long, max 100")
	}
	for _, account := range msg.Accounts {
		if _, err := sdk.AccAddressFromBech32(account.Address); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("address %s", account.Address))
		}
	}
	return nil
}
