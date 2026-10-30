package regiongov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/regiongov/keeper"
)

// EndBlocker tallies proposals whose voting period has ended.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.EndBlocker(ctx)
}
