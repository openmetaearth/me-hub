package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrikeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	Wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
)

type Hooks struct {
	distrikeeper.Hooks
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}
var _ Wstakingtypes.WstakingHooks = Hooks{}

// overwrite
// Create new distribution hooks
func (k Keeper) Hooks() Hooks {
	return Hooks{Hooks: k.Keeper.Hooks(), k: k}
}

// overwrite
// withdraw delegation rewards (which also increments period)
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

func (h Hooks) BeforeValidatorStakingModified(ctx sdk.Context, valAddr sdk.ValAddress) error {
	//TODO: distribution the block reward
	return nil
}
