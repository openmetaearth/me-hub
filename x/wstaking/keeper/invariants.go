package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// RegisterInvariants registers all staking invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k *Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-accounts",
		ModuleAccountInvariants(k))
	ir.RegisterRoute(types.ModuleName, "nonnegative-power",
		stakingkeeper.NonNegativePowerInvariant(k.Keeper))
	ir.RegisterRoute(types.ModuleName, "positive-delegation",
		stakingkeeper.PositiveDelegationInvariant(k.Keeper))
	ir.RegisterRoute(types.ModuleName, "delegator-shares",
		DelegatorSharesInvariant(k))
}

// ModuleAccountInvariants checks that the bonded and notBonded ModuleAccounts pools
// reflects the tokens actively bonded and not bonded
func ModuleAccountInvariants(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		bonded := math.ZeroInt()
		notBonded := math.ZeroInt()
		bondedPool := k.GetBondedStakePool(ctx)
		notBondedPool := k.GetNotBondedStakePool(ctx)
		bondDenom := k.BondDenom(ctx)

		k.IterateValidators(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
			switch validator.GetStatus() {
			case stakingtypes.Bonded:
				bonded = bonded.Add(validator.GetTokens())
			case stakingtypes.Unbonding, stakingtypes.Unbonded:
				notBonded = notBonded.Add(validator.GetTokens())
			default:
				panic("invalid validator status")
			}
			return false
		})

		k.IterateUnbondingStakes(ctx, func(ubd types.UnbondingStake) bool {
			for _, entry := range ubd.Entries {
				notBonded = notBonded.Add(entry.Balance)
			}
			return false
		})

		poolBonded := k.bankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom)
		poolNotBonded := k.bankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom)
		broken := !poolBonded.Amount.Equal(bonded) || !poolNotBonded.Amount.Equal(notBonded)

		// Bonded tokens should equal sum of tokens with bonded validators
		// Not-bonded tokens should equal unbonding delegations	plus tokens on unbonded validators
		return sdk.FormatInvariant(types.ModuleName, "bonded and not bonded module account coins", fmt.Sprintf(
			"\tPool's bonded tokens: %v\n"+
				"\tsum of bonded tokens: %v\n"+
				"not bonded token invariance:\n"+
				"\tPool's not bonded tokens: %v\n"+
				"\tsum of not bonded tokens: %v\n"+
				"module accounts total (bonded + not bonded):\n"+
				"\tModule Accounts' tokens: %v\n"+
				"\tsum tokens:              %v\n",
			poolBonded, bonded, poolNotBonded, notBonded, poolBonded.Add(poolNotBonded), bonded.Add(notBonded))), broken
	}
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
