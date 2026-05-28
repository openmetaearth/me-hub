package types

import (

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

func NewMsgRemove(issuer, did string) *MsgRemove {
	return &MsgRemove{
		Issuer: issuer,
		Did:    did,
	}
}

func (m *MsgRemove) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Issuer); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the issuer is not a valid bech32 address")
	}
	if len(m.Did) != didtypes.DidLength {
		return errors.Wrapf(sdkerrors.ErrInvalidPubKey, "DID length must be equal to %d", didtypes.DidLength)
	}

	return nil
}
