package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const TypeMsgUpdateValidator = "update_validator"

var _ sdk.Msg = &MsgUpdateValidator{}

// NewMsgUpdateValidator creates a new MsgUpdateValidator instance
//
//nolint:interfacer
func NewMsgUpdateValidator(valAddr sdk.ValAddress, description stakingtypes.Description, newRate *sdk.Dec, newMinSelfDelegation *math.Int) *MsgUpdateValidator {
	return &MsgUpdateValidator{
		Description:       description,
		CommissionRate:    newRate,
		StakerAddress:     valAddr.String(),
		MinSelfDelegation: newMinSelfDelegation,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUpdateValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUpdateValidator) Type() string { return TypeMsgUpdateValidator }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateValidator) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.StakerAddress)
	return []sdk.AccAddress{sdk.AccAddress(addr)}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUpdateValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateValidator) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.StakerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if msg.Description == (stakingtypes.Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.MinSelfDelegation != nil && !msg.MinSelfDelegation.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"minimum self delegation must be a positive integer",
		)
	}

	if msg.CommissionRate != nil {
		if msg.CommissionRate.GT(math.LegacyOneDec()) || msg.CommissionRate.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "commission rate must be between 0 and 1 (inclusive)")
		}
	}

	return nil
}
