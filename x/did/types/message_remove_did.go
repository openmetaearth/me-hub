package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRemoveDid = "remove_did"
)

func NewMsgRemoveDid(creator, did string) *MsgRemoveDid {
	return &MsgRemoveDid{
		Creator: creator,
		Did:     did,
	}
}

// Route implements the sdk.Msg interface.
func (m *MsgRemoveDid) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgRemoveDid) Type() string { return TypeMsgRemoveDid }

func (m *MsgRemoveDid) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgRemoveDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRemoveDid) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the creator is not a valid bech32 address")
	}
	if len(m.Did) != DidLength {
		return errors.Wrap(sdkerrors.ErrInvalidType, fmt.Sprintf("DID length must be equal to %d", DidLength))
	}

	return nil
}
