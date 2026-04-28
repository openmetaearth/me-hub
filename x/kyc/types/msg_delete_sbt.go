package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

const (
	TypeMsgDeleteSBT = "delete_sbt"
)

func NewMsgDeleteSBT(issuer, did string) *MsgDeleteSBT {
	return &MsgDeleteSBT{
		Issuer: issuer,
		Did:    did,
	}
}

// Route implements the sdk.Msg interface.
func (m *MsgDeleteSBT) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgDeleteSBT) Type() string { return TypeMsgDeleteSBT }

func (m *MsgDeleteSBT) GetSigners() []sdk.AccAddress {
	issuer, err := sdk.AccAddressFromBech32(m.Issuer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{issuer}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgDeleteSBT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgDeleteSBT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Issuer); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the holder is not a valid bech32 address")
	}
	if len(m.Did) != didtypes.DidLength {
		return errors.Wrapf(sdkerrors.ErrInvalidPubKey, "DID length must be equal to %d", didtypes.DidLength)
	}

	return nil
}
