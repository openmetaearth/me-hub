package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func NewMsgUpdateValidator(valAddr sdk.ValAddress, description stakingtypes.Description, newRate *sdkmath.LegacyDec, newMinSelfDelegation *sdkmath.Int) *MsgUpdateValidator {
	return &MsgUpdateValidator{
		Description:       description,
		CommissionRate:    newRate,
		StakerAddress:     valAddr.String(),
		MinSelfDelegation: newMinSelfDelegation,
	}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateValidator) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.StakerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if msg.Description == (stakingtypes.Description{}) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.MinSelfDelegation != nil && !msg.MinSelfDelegation.IsPositive() {
		return errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"minimum self delegation must be a positive integer",
		)
	}

	if msg.CommissionRate != nil {
		if msg.CommissionRate.GT(sdkmath.LegacyOneDec()) || msg.CommissionRate.IsNegative() {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "commission rate must be between 0 and 1 (inclusive)")
		}
	}

	return nil
}
