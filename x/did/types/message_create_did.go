package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateDid = "create_did"
)

func NewMsgCreateDid(creator, did, pubkey string) *MsgCreateDid {
	return &MsgCreateDid{
		Creator: creator,
		Did:     did,
		Pubkey:  pubkey,
	}
}

// Route implements the sdk.Msg interface.
func (m *MsgCreateDid) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgCreateDid) Type() string { return TypeMsgCreateDid }

func (m *MsgCreateDid) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgCreateDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgCreateDid) GetDidInfo() DidInfo {
	return NewDidInfo(m.Did, m.Address, m.Pubkey, DID_STATUS_ACTIVE)
}

func (m *MsgCreateDid) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the creator is not a valid bech32 address")
	}
	if len(m.Did) != DidLength {
		return errors.Wrap(sdkerrors.ErrInvalidType, fmt.Sprintf("DID length must be equal to %d", DidLength))
	}
	if m.Pubkey == "" {
		return errors.Wrap(sdkerrors.ErrInvalidPubKey, "pubkey must not be nil")
	}

	return nil
}
