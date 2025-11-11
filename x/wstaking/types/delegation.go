package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// MustUnmarshalDelegation return the unmarshaled delegation from bytes.
// Panics if fails.
func MustUnmarshalDelegation(cdc codec.BinaryCodec, value []byte) stakingtypes.Delegation {
	delegation, err := UnmarshalDelegation(cdc, value)
	if err != nil {
		panic(err)
	}
	return delegation
}

// return the delegation
func UnmarshalDelegation(cdc codec.BinaryCodec, value []byte) (delegation stakingtypes.Delegation, err error) {
	err = cdc.Unmarshal(value, &delegation)
	return delegation, err
}

// NewDelegation creates a new delegation object
//
//nolint:interfacer
func NewDelegation(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, shares sdk.Dec) stakingtypes.Delegation {
	return stakingtypes.Delegation{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           sdk.ZeroDec(),
		StartHeight:      0,
		Amount:           sdk.ZeroInt(),
		Unmovable:        sdk.ZeroInt(),
		UnMeidAmount:     sdk.ZeroInt(),
	}
}
