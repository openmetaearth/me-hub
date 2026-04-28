package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

const (
	TypeMsgCreateSBT = "create_sbt"
)

func NewMsgCreateSBT(issuer, did, uri, uriHash string, data []byte) *MsgCreateSBT {
	return &MsgCreateSBT{
		Issuer:  issuer,
		Did:     did,
		Uri:     uri,
		UriHash: uriHash,
		Data:    data,
	}
}

// Route implements the sdk.Msg interface.
func (m *MsgCreateSBT) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgCreateSBT) Type() string { return TypeMsgCreateSBT }

func (m *MsgCreateSBT) GetSigners() []sdk.AccAddress {
	issuer, err := sdk.AccAddressFromBech32(m.Issuer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{issuer}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgCreateSBT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgCreateSBT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Issuer); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the issuer is not a valid bech32 address")
	}
	if len(m.Did) != didtypes.DidLength {
		return errors.Wrap(sdkerrors.ErrInvalidType, fmt.Sprintf("DID length must be equal to %d", didtypes.DidLength))
	}
	if len(m.UriHash) == 0 || len(m.UriHash) > 128 {
		return errors.Wrap(sdkerrors.ErrInvalidType, "uri hash length must be between 0 and 128")
	}

	return nil
}
