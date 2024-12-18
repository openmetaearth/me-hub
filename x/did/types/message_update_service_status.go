package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateServiceStatus = "update_service_status"
)

func NewMsgUpdateServiceStatus(creator, sid string, status ServiceStatus) *MsgUpdateServiceStatus {
	return &MsgUpdateServiceStatus{
		Creator: creator,
		Sid:     sid,
		Status:  status,
	}
}

// Route implements the sdk.Msg interface.
func (m *MsgUpdateServiceStatus) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgUpdateServiceStatus) Type() string { return TypeMsgCreateDid }

func (m *MsgUpdateServiceStatus) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgUpdateServiceStatus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgUpdateServiceStatus) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the creator is not a valid bech32 address")
	}
	if len(m.Sid) < 2 || len(m.Sid) > 8 {
		return errors.Wrap(sdkerrors.ErrInvalidType, "sid length must be between 2 and 8")
	}
	if _, ok := ServiceStatus_name[int32(m.Status)]; !ok {
		return errors.Wrap(sdkerrors.ErrInvalidType, "service status must be ACTIVE or INACTIVE")
	}

	return nil
}
