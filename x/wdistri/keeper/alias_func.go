package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/st-chain/me-hub/x/wdistri/types"
)

// GetDistributionAccount returns the distribution ModuleAccount
func (k Keeper) GetDistributionMintAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.authKeeper.GetModuleAccount(ctx, types.ReceiveMintReward)
}
