package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// bondedTokensToNotBonded handles the conversion of bonded tokens to not bonded tokens
func (k Keeper) bondedTokensToNotBonded(ctx sdk.Context, validator sdk.ValAddress, amount sdkmath.Int) error {
	if amount.LT(sdk.ZeroInt()) {
		k.Logger(ctx).Error("invalid amount", "error", "amount is less than zero")
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeError,
				sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
				sdk.NewAttribute(sdk.AttributeKeyError, "amount is less than zero"),
			),
		)
		return errorsmod.Wrap(types.ErrInvalid, "amount is less than zero")
	}
	// implementation remains the same
}