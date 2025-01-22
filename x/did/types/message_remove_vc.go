package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRemoveVC = "remove_vc"
)

func NewMsgRemoveVC(issuer, did, sid string) *MsgRemoveVC {
	return &MsgRemoveVC{
		Issuer: issuer,
		Did:    did,
		Sid:    sid,
	}
}

// Route implements the sdk.Msg interface.
func (m *MsgRemoveVC) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgRemoveVC) Type() string { return TypeMsgRemoveVC }

func (m *MsgRemoveVC) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Issuer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgRemoveVC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRemoveVC) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Issuer); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the creator is not a valid bech32 address")
	}
	if len(m.Did) != DidLength {
		return errors.Wrap(sdkerrors.ErrInvalidType, fmt.Sprintf("DID length must be equal to %d", DidLength))
	}
	if len(m.Sid) < 2 || len(m.Sid) > 8 {
		return errors.Wrap(sdkerrors.ErrInvalidType, "sid length must be between 2 and 8")
	}

	return nil
}
