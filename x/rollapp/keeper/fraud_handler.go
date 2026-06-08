package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func (k Keeper) HandleFraud(ctx sdk.Context, appID string) error {
	// Freeze the app
	k.SetFrozen(ctx, appID, true)

	// Add unfreeze mechanism
	k.SetUnfreezeTime(ctx, appID, ctx.BlockTime().Add(types.UnfreezeDuration))
	return nil
}

func (k Keeper) UnfreezeApp(ctx sdk.Context, appID string) error {
	// Check if the app is frozen and the unfreeze time has passed
	if k.GetFrozen(ctx, appID) && ctx.BlockTime().After(k.GetUnfreezeTime(ctx, appID)) {
		// Unfreeze the app
		k.SetFrozen(ctx, appID, false)
		return nil
	}
	return types.ErrAppNotFrozen
}