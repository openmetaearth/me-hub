package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
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
