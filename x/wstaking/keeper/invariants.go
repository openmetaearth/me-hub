package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// RegisterInvariants registers all staking invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k *Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-accounts",
		stakingkeeper.ModuleAccountInvariants(k.Keeper))
	ir.RegisterRoute(types.ModuleName, "nonnegative-power",
		stakingkeeper.NonNegativePowerInvariant(k.Keeper))
	ir.RegisterRoute(types.ModuleName, "positive-delegation",
		stakingkeeper.PositiveDelegationInvariant(k.Keeper))
	ir.RegisterRoute(types.ModuleName, "delegator-shares",
		DelegatorSharesInvariant(k))
}

// DelegatorSharesInvariant checks whether all the delegator shares which persist
// in the delegator object add up to the correct total delegator shares
// amount stored in each validator.
func DelegatorSharesInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			msg    string
			broken bool
		)

		validators := k.GetAllValidators(ctx)
		validatorsDelegationShares := map[string]sdk.Dec{}

		// initialize a map: validator -> its delegation shares
		for _, validator := range validators {
			validatorsDelegationShares[validator.GetOperator().String()] = math.LegacyZeroDec()
		}

		// iterate through all the delegations to calculate the total delegation shares for each validator
		stakes := k.GetAllStakes(ctx)
		for _, stake := range stakes {
			stakeValidatorAddr := stake.GetValidatorAddr().String()
			validatorDelegationShares := validatorsDelegationShares[stakeValidatorAddr]
			validatorsDelegationShares[stakeValidatorAddr] = validatorDelegationShares.Add(stake.Shares)
		}

		// for each validator, check if its total delegation shares calculated from the step above equals to its expected delegation shares
		for _, validator := range validators {
			expValTotalDelShares := validator.GetDelegatorShares()
			calculatedValTotalDelShares := validatorsDelegationShares[validator.GetOperator().String()]
			if !calculatedValTotalDelShares.Equal(expValTotalDelShares) {
				broken = true
				msg += fmt.Sprintf("broken delegator shares invariance:\n"+
					"\tvalidator.DelegatorShares: %v\n"+
					"\tsum of Delegator.Shares: %v\n", expValTotalDelShares, calculatedValTotalDelShares)
			}
		}

		return sdk.FormatInvariant(types.ModuleName, "delegator shares", msg), broken
	}
}
