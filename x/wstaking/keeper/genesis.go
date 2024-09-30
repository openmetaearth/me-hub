package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
)

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, validators, and bonds found in
// the keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *wstakingtypes.GenesisState {
	var unbondingDelegations []types.UnbondingDelegation

	k.IterateUnbondingDelegations(ctx, func(_ int64, ubd types.UnbondingDelegation) (stop bool) {
		unbondingDelegations = append(unbondingDelegations, ubd)
		return false
	})

	var unbondingStakes []wstakingtypes.UnbondingStake

	k.IterateUnbondingStakes(ctx, func(_ int64, ubs wstakingtypes.UnbondingStake) (stop bool) {
		unbondingStakes = append(unbondingStakes, ubs)
		return false
	})

	var redelegations []types.Redelegation

	k.IterateRedelegations(ctx, func(_ int64, red types.Redelegation) (stop bool) {
		redelegations = append(redelegations, red)
		return false
	})

	var lastValidatorPowers []wstakingtypes.LastValidatorPower

	k.IterateLastValidatorPowers(ctx, func(addr sdk.ValAddress, power int64) (stop bool) {
		lastValidatorPowers = append(lastValidatorPowers, wstakingtypes.LastValidatorPower{Address: addr.String(), Power: power})
		return false
	})

	return &wstakingtypes.GenesisState{
		Params:               k.GetParams(ctx),
		LastTotalPower:       k.GetLastTotalPower(ctx),
		LastValidatorPowers:  lastValidatorPowers,
		Validators:           k.GetAllValidators(ctx),
		Delegations:          k.GetAllDelegations(ctx),
		UnbondingDelegations: unbondingDelegations,
		Redelegations:        redelegations,
		Stakes:               k.GetAllStakes(ctx),
		UnbondingStakes:      unbondingStakes,
		RegionList:           k.GetAllRegion(ctx),
		MeidList:             k.GetAllMeid(ctx),
		FixedDepositList:     k.GetAllFixedDeposit(ctx),
		FixedDepositCount:    k.GetFixedDepositCount(ctx),
		Exported:             true,
	}
}
