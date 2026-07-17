package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateRollapp = "update_rollapp"

var _ sdk.Msg = &MsgUpdateRollapp{}

func NewMsgUpdateRollapp(creator, rollappId, channelId string, maxSequencers uint64, permissionedAddresses []string) *MsgUpdateRollapp {
	return &MsgUpdateRollapp{
		Creator:               creator,
		RollappId:             rollappId,
		MaxSequencers:         maxSequencers,
		PermissionedAddresses: permissionedAddresses,
		ChannelId:             channelId,
	}
}

func (msg *MsgUpdateRollapp) Route() string {
	return RouterKey
}

func (msg *MsgUpdateRollapp) Type() string {
	return TypeMsgUpdateRollapp
}

func (msg *MsgUpdateRollapp) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateRollapp) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateRollapp) ValidateBasic() error {
	if msg.RollappId == "" {
		return ErrInvalidRollappID
	}
	return nil
}
