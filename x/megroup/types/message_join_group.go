package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgJoinGroup = "join_group"

var _ sdk.Msg = &MsgJoinGroup{}

func NewMsgJoinGroup(creator string, groupId uint64, applicantAddress string) *MsgJoinGroup {
	return &MsgJoinGroup{
		Creator:          creator,
		GroupId:          groupId,
		ApplicantAddress: applicantAddress,
	}
}

func (msg *MsgJoinGroup) Route() string {
	return RouterKey
}

func (msg *MsgJoinGroup) Type() string {
	return TypeMsgJoinGroup
}

func (msg *MsgJoinGroup) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgJoinGroup) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgJoinGroup) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
