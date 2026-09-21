package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgBondZkCoreID = "bond_zk_core_id"

func NewMsgBondZkCoreID(
	creator string,
	coreInfo *ZkAuthClaimCoreIdInfo,
) *MsgBondZkCoreIDRequest {
	return &MsgBondZkCoreIDRequest{Creator: creator, CoreInfo: coreInfo}
}

func (m *MsgBondZkCoreIDRequest) Route() string { return RouterKey }

func (m *MsgBondZkCoreIDRequest) Type() string { return TypeMsgBondZkCoreID }

func (m *MsgBondZkCoreIDRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (m *MsgBondZkCoreIDRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgBondZkCoreIDRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "creator is not a valid bech32 address")
	}
	if m.CoreInfo == nil {
		return errors.Wrap(ErrParameter, "core info is required")
	}
	if err := ValidateZkCoreID(m.CoreInfo.CoreId); err != nil {
		return errors.Wrap(ErrInvalidZkCoreID, err.Error())
	}
	if len(m.CoreInfo.AuthProof) == 0 {
		return ErrInvalidZkProof
	}
	if len(m.CoreInfo.AuthProof) > MaxZkAuthProofLength {
		return errors.Wrapf(
			ErrInvalidZkProof,
			"proof length %d exceeds %d",
			len(m.CoreInfo.AuthProof),
			MaxZkAuthProofLength,
		)
	}
	if m.CoreInfo.ExpiresAt <= 0 {
		return errors.Wrap(ErrParameter, "expires_at must be positive")
	}
	return nil
}
