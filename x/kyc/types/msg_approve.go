package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/utils"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

const (
	TypeMsgApprove = "approve"
)

func NewMsgApprove(issuer, did, regionId, address, pubkey string, level didtypes.KycLevel, uri, hash, inviter string) *MsgApprove {
	return &MsgApprove{
		Issuer:   issuer,
		Did:      did,
		RegionId: regionId,
		Address:  address,
		Pubkey:   pubkey,
		Level:    level,
		Uri:      uri,
		Hash:     hash,
		Inviter:  inviter,
	}
}

// Route implements the sdk.Msg interface.
func (m *MsgApprove) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (m *MsgApprove) Type() string { return TypeMsgApprove }

func (m *MsgApprove) GetSigners() []sdk.AccAddress {
	issuer, err := sdk.AccAddressFromBech32(m.Issuer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{issuer}
}

// GetSignBytes returns the message bytes to sign over.
func (m *MsgApprove) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgApprove) GetKYC() didtypes.Credential {
	return didtypes.NewCredential(m.Did, ModuleName, m.Hash, m.Uri, []byte(m.RegionId))
}

func (m *MsgApprove) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Issuer); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the issuer is not a valid bech32 address")
	}
	if len(m.Did) != didtypes.DidLength {
		return errors.Wrapf(sdkerrors.ErrInvalidType, "DID length must be equal to %d", didtypes.DidLength)
	}
	if _, err := utils.CheckRegionName(strings.ToUpper(m.RegionId)); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidType, err.Error())
	}
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "the address is not a valid bech32 address")
	}
	if m.Pubkey == "" {
		return errors.Wrap(sdkerrors.ErrInvalidPubKey, "the pubkey is empty")
	}
	if _, ok := didtypes.KycLevel_name[int32(m.Level)]; !ok {
		return errors.Wrap(sdkerrors.ErrInvalidType, "the level is not valid")
	}
	//if len(m.Hash) == 0 || len(m.Hash) > 128 {
	//	return errors.Wrap(sdkerrors.ErrInvalidType, "hash length must be between 0 and 128")
	//}
	if m.Inviter != "" {
		if _, err := sdk.AccAddressFromBech32(m.Inviter); err != nil {
			return errors.Wrap(sdkerrors.ErrInvalidAddress, "the inviter is not a valid bech32 address")
		}
	}

	return nil
}
