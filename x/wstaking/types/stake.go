package types

import (
	"cosmossdk.io/math"
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

// AddEntry - append entry to the unbonding stake
func (ubs *UnbondingStake) AddEntry(creationHeight int64, minTime time.Time, balance math.Int) {
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
