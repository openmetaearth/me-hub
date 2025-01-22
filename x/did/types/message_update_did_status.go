package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateDidStatus = "update_did_status"
)

func NewMsgUpdateDidStatus(creator, did string, status DidStatus) *MsgUpdateDidStatus {
	return &MsgUpdateDidStatus{
		Creator: creator,
		Did:     did,
		Status:  status,
	}
}

// Route implements the sdk.Msg interface.
func (m *MsgUpdateDidStatus) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgUpdateDidStatus) Type() string { return TypeMsgCreateDid }

func (m *MsgUpdateDidStatus) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgUpdateDidStatus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgUpdateDidStatus) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the creator is not a valid bech32 address")
	}
	if len(m.Did) != DidLength {
		return errors.Wrap(sdkerrors.ErrInvalidType, fmt.Sprintf("DID length must be equal to %d", DidLength))
	}
	if _, ok := DidStatus_name[int32(m.Status)]; !ok {
		return errors.Wrap(sdkerrors.ErrInvalidType, "DID status must be ACTIVE or INACTIVE")
	}

	return nil
}
