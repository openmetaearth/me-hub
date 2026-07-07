package types

import (
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
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
