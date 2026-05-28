package types

import (

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

func NewMsgDeleteSBT(issuer, did string) *MsgDeleteSBT {
	return &MsgDeleteSBT{
		Issuer: issuer,
		Did:    did,
	}
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
