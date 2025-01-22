package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateVC = "update_vc"
)

func NewMsgUpdateVC(issuer, did, sid, hash, uri string, data []byte, filters [][]byte) *MsgUpdateVC {
	return &MsgUpdateVC{
		Issuer:  issuer,
		Did:     did,
		Sid:     sid,
		Hash:    hash,
		Uri:     uri,
		Data:    data,
		Filters: filters,
	}
}

// Route implements the sdk.Msg interface.
func (m *MsgUpdateVC) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgUpdateVC) Type() string { return TypeMsgUpdateVC }

func (m *MsgUpdateVC) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Issuer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgUpdateVC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgUpdateVC) GetCredential() Credential {
	return NewCredential(m.Did, m.Sid, m.Hash, m.Uri, m.Data)
}

func (m *MsgUpdateVC) ValidateBasic() error {
	// check issuer
	if _, err := sdk.AccAddressFromBech32(m.Issuer); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the creator is not a valid bech32 address")
	}
	if len(m.Did) != DidLength {
		return errors.Wrap(sdkerrors.ErrInvalidType, fmt.Sprintf("DID length must be equal to %d", DidLength))
	}
	if len(m.Sid) < 2 || len(m.Sid) > 8 {
		return errors.Wrap(sdkerrors.ErrInvalidType, "sid length must be between 2 and 8")
	}
	if len(m.Hash) == 0 || len(m.Hash) > 128 {
		return errors.Wrap(sdkerrors.ErrInvalidType, "hash length must be between 0 and 128")
	}
	if len(m.Uri) > 1024 {
		return errors.Wrap(sdkerrors.ErrInvalidType, "uri length exceeds 1024")
	}

	for _, filter := range m.Filters {
		if len(filter) > 1024 {
			return errors.Wrap(sdkerrors.ErrInvalidType, "filter length exceeds 1024")
		}
	}

	return nil
}
