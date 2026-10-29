package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateKycIssuer       = "update_kyc_issuer"
	KycIssuerAdd           int32 = 0
	KycIssuerRemove        int32 = 1
)

var _ sdk.Msg = &MsgUpdateKycIssuer{}

func NewMsgUpdateKycIssuer(creator string, addresses []string, optionType int32) *MsgUpdateKycIssuer {
	return &MsgUpdateKycIssuer{Creator: creator, Addresses: addresses, OptionType: optionType}
}

func (m *MsgUpdateKycIssuer) Route() string { return RouterKey }
func (m *MsgUpdateKycIssuer) Type() string  { return TypeMsgUpdateKycIssuer }

func (m *MsgUpdateKycIssuer) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (m *MsgUpdateKycIssuer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgUpdateKycIssuer) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if len(m.Addresses) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "addresses must not be empty")
	}
	if m.OptionType != KycIssuerAdd && m.OptionType != KycIssuerRemove {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "option_type must be 0 (add) or 1 (remove)")
	}
	for _, address := range m.Addresses {
		if _, err := sdk.AccAddressFromBech32(address); err != nil {
			return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid issuer address: %s", address)
		}
	}
	return nil
}
