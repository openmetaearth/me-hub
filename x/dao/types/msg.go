package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgUpdateGlobalDao{}
)

func NewMsgUpdateGlobalDao(creator sdk.AccAddress, address sdk.AccAddress) *MsgUpdateGlobalDao {
	return &MsgUpdateGlobalDao{
		Creator: creator.String(),
		Address: address.String(),
	}
}

func (msg *MsgUpdateGlobalDao) Route() string {
	return RouterKey
}

func (msg *MsgUpdateGlobalDao) Type() string {
	return "UpdateGlobalDao"
}

func (msg *MsgUpdateGlobalDao) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic("invalid creator address")
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateGlobalDao) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateGlobalDao) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Creator)
	}
	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Address)
	}
	return nil
}
