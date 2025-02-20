package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgUpdateDao{}
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
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Creator)
	}
	return nil
}
