package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgClaimRewards represents a message to claim rewards
type MsgClaimRewards struct {
	ValidatorAddr string `json:"validator_addr" yaml:"validator_addr"`
}

// NewMsgClaimRewards returns a new MsgClaimRewards
func NewMsgClaimRewards(validatorAddr string) MsgClaimRewards {
	return MsgClaimRewards{
		ValidatorAddr: validatorAddr,
	}
}

// Route implements the sdk.Msg interface
func (msg MsgClaimRewards) Route() string {
	return ModuleName
}

// Type implements the sdk.Msg interface
func (msg MsgClaimRewards) Type() string {
	return "claim_rewards"
}

// ValidateBasic implements the sdk.Msg interface
func (msg MsgClaimRewards) ValidateBasic() error {
	if len(msg.ValidatorAddr) == 0 {
		return sdk.ErrInvalidAddress("validator address cannot be empty")
	}

	return nil
}

// GetSignBytes implements the sdk.Msg interface
func (msg MsgClaimRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements the sdk.Msg interface
func (msg MsgClaimRewards) GetSigners() []sdk.AccAddress {
	// Only the validator can claim their rewards
	validatorAddr, err := sdk.AccAddressFromBech32(msg.ValidatorAddr)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{validatorAddr}
}