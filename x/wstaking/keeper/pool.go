package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// GetBondedStakePool returns the bonded stake tokens pool's module account
func (k Keeper) GetBondedStakePool(ctx sdk.Context) (bondedStakePool authtypes.ModuleAccountI) {
	return k.authKeeper.GetModuleAccount(ctx, types.BondedStakePoolName)
}

// GetNotBondedStakePool returns the not bonded stake tokens pool's module account
func (k Keeper) GetNotBondedStakePool(ctx sdk.Context) (bondedStakePool authtypes.ModuleAccountI) {
	return k.authKeeper.GetModuleAccount(ctx, types.NotBondedStakePoolName)
}

func (k Keeper) TotalBondedStakePool(ctx sdk.Context) math.Int {
	bondedStakePool := k.GetBondedStakePool(ctx)
	return k.bankKeeper.GetBalance(ctx, bondedStakePool.GetAddress(), k.BondDenom(ctx)).Amount
}

// bondedStakeTokensToNotBonded transfers coins from the bonded to the not bonded pool within staking
func (k Keeper) BondedStakeTokensToNotBonded(ctx sdk.Context, tokens math.Int, regionID string) {
	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), tokens))
	if err := k.bankKeeper.Extend().SendCoinsFromModuleToModuleWithTag(ctx, types.BondedStakePoolName, types.NotBondedStakePoolName, coins,
		fmt.Sprintf("BondedStakeTokensToNotBonded_%s", regionID),
	); err != nil {
		panic(err)
	}
}

// notBondedStakeTokensToBonded transfers coins from the not bonded to the bonded pool within staking
func (k Keeper) NotBondedStakeTokensToBonded(ctx sdk.Context, tokens math.Int) {
	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), tokens))
	if err := k.bankKeeper.Extend().SendCoinsFromModuleToModuleWithTag(ctx, types.NotBondedStakePoolName, types.BondedStakePoolName, coins,
		"NotBondedStakeTokensToBonded",
	); err != nil {
		panic(err)
	}
}
