package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRemoveService = "remove_service"
)

func NewMsgRemoveService(creator, sid string) *MsgRemoveService {
	return &MsgRemoveService{
		Creator: creator,
		Sid:     sid,
	}
}

// Route implements the sdk.Msg interface.
func (m *MsgRemoveService) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgRemoveService) Type() string { return TypeMsgRemoveService }

func (m *MsgRemoveService) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgRemoveService) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRemoveService) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the creator is not a valid bech32 address")
	}

	if len(m.Sid) < 2 || len(m.Sid) > 8 {
		return errors.Wrap(sdkerrors.ErrInvalidType, "sid length must be between 2 and 8")
	}

	return nil
}
