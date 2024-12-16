package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

// GetBondedStakePool returns the bonded stake tokens pool's module account
func (k Keeper) GetBondedStakePool(ctx sdk.Context) (bondedStakePool authtypes.ModuleAccountI) {
	return k.AuthKeeper.GetModuleAccount(ctx, types.BondedStakePoolName)
}

func (k Keeper) TotalBondedStakePool(ctx sdk.Context) math.Int {
	bondedStakePool := k.GetBondedStakePool(ctx)
	return k.BankKeeper.GetBalance(ctx, bondedStakePool.GetAddress(), k.BondDenom(ctx)).Amount
}
