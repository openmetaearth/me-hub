package regiongov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/regiongov/keeper"
)

// EndBlocker tallies proposals at VotingEndTime and executes those that pass.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.EndBlocker(ctx)
}
