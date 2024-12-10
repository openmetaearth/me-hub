package types

import (
	"cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sigs.k8s.io/yaml"
	"time"
)

func (s Stake) GetValidatorAddr() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(s.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (s Stake) GetShares() sdk.Dec { return s.Shares }

func (s Stake) GetStakerAddr() sdk.AccAddress {
	stakerAddress := sdk.MustAccAddressFromBech32(s.StakerAddress)
	return stakerAddress
}

// NewStake creates a new stake object
//
//nolint:interfacer
func NewStake(stakerAddr sdk.AccAddress, validatorAddr sdk.ValAddress, shares sdk.Dec) Stake {
	return Stake{
		StakerAddress:    stakerAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           shares,
		StartHeight:      0,
		Rewards:          sdk.ZeroDec(),
		Amount:           sdk.ZeroInt(),
		Unmovable:        sdk.ZeroInt(),
	}
}

// MustUnmarshalStake return the unmarshaled stake from bytes.
// Panics if fails.
func MustUnmarshalStake(cdc codec.BinaryCodec, value []byte) Stake {
	stake, err := UnmarshalStake(cdc, value)
	if err != nil {
		panic(err)
	}

	return stake
}

// return the stake
func UnmarshalStake(cdc codec.BinaryCodec, value []byte) (stake Stake, err error) {
	err = cdc.Unmarshal(value, &stake)
	return stake, err
}

// MustMarshalStake returns the stake bytes. Panics if fails
func MustMarshalStake(cdc codec.BinaryCodec, stake Stake) []byte {
	return cdc.MustMarshal(&stake)
}

// unmarshal a unbonding stake from a store value
func MustUnmarshalUBS(cdc codec.BinaryCodec, value []byte) UnbondingStake {
	ubd, err := UnmarshalUBS(cdc, value)
	if err != nil {
		panic(err)
	}

	return ubd
}

// return the unbonding stake
func MustMarshalUBS(cdc codec.BinaryCodec, ubs UnbondingStake) []byte {
	return cdc.MustMarshal(&ubs)
}

// unmarshal a unbonding stake from a store value
func UnmarshalUBS(cdc codec.BinaryCodec, value []byte) (ubs UnbondingStake, err error) {
	err = cdc.Unmarshal(value, &ubs)
	return ubs, err
}

// String returns a human readable string representation of an UnbondingStake.
func (ubs UnbondingStake) String() string {
	out := fmt.Sprintf(`Unbonding Stakes between:
  staker:                 %s
  Validator:                 %s
	Entries:`, ubs.StakerAddress, ubs.ValidatorAddress)
	for i, entry := range ubs.Entries {
		out += fmt.Sprintf(`    Unbonding Stake %d:
      Creation Height:           %v
      Min time to unbond (unix): %v
      Expected balance:          %s`, i, entry.CreationHeight,
			entry.CompletionTime, entry.Balance)
	}

	return out
}

// AddEntry - append entry to the unbonding stake
func (ubs *UnbondingStake) AddEntry(creationHeight int64, minTime time.Time, balance math.Int) {
	entry := NewUnbondingStakeEntry(creationHeight, minTime, balance)
	ubs.Entries = append(ubs.Entries, entry)
}

// RemoveEntry - remove entry at index i to the unbonding stake
func (ubd *UnbondingStake) RemoveEntry(i int64) {
	ubd.Entries = append(ubd.Entries[:i], ubd.Entries[i+1:]...)
}

// String implements the stringer interface for a UnbondingStakeEntry.
func (e UnbondingStakeEntry) String() string {
	out, _ := yaml.Marshal(e)
	return string(out)
}

// IsMature - is the current entry mature
func (e UnbondingStakeEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

func NewUnbondingStakeEntry(creationHeight int64, completionTime time.Time, balance math.Int) UnbondingStakeEntry {
	return UnbondingStakeEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		Balance:        balance,
	}
}

// NewUnbondingStake - create a new unbonding stake object
//
//nolint:interfacer
func NewUnbondingStake(
	stakerAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, balance math.Int,
) UnbondingStake {
	return UnbondingStake{
		StakerAddress:    stakerAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Entries: []UnbondingStakeEntry{
			NewUnbondingStakeEntry(creationHeight, minTime, balance),
		},
	}
}
