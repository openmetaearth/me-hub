package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func (s Stake) GetValidatorAddr() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(s.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (s Stake) GetShares() sdkmath.LegacyDec { return s.Shares }

func (s Stake) GetStakerAddr() sdk.AccAddress {
	stakerAddress := sdk.MustAccAddressFromBech32(s.StakerAddress)
	return stakerAddress
}

// NewStake creates a new stake object
//
//nolint:interfacer
func NewStake(stakerAddr sdk.AccAddress, validatorAddr sdk.ValAddress, shares sdkmath.LegacyDec) Stake {
	return Stake{
		StakerAddress:    stakerAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Shares:           shares,
		StartHeight:      0,
		Rewards:          sdkmath.LegacyZeroDec(),
		Amount:           sdkmath.ZeroInt(),
		Unmovable:        sdkmath.ZeroInt(),
	}
}

// AddEntry - append entry to the unbonding stake
func (ubs *UnbondingStake) AddEntry(creationHeight int64, minTime time.Time, balance sdkmath.Int) {
	entry := NewUnbondingStakeEntry(creationHeight, minTime, balance)
	ubs.Entries = append(ubs.Entries, entry)
}

// RemoveEntry - remove entry at index i to the unbonding stake
func (ubd *UnbondingStake) RemoveEntry(i int64) {
	ubd.Entries = append(ubd.Entries[:i], ubd.Entries[i+1:]...)
}

// IsMature - is the current entry mature
func (e UnbondingStakeEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

func NewUnbondingStakeEntry(creationHeight int64, completionTime time.Time, balance sdkmath.Int) UnbondingStakeEntry {
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
	creationHeight int64, minTime time.Time, balance sdkmath.Int,
) UnbondingStake {
	return UnbondingStake{
		StakerAddress:    stakerAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Entries: []UnbondingStakeEntry{
			NewUnbondingStakeEntry(creationHeight, minTime, balance),
		},
	}
}
