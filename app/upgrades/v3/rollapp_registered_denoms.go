package v3

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	legacyrollapp "github.com/openmetaearth/me-hub/app/upgrades/v3/types/rollapp"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
)

// migrateRollappRegisteredDenoms copies v2 Rollapp.registeredDenoms (proto field 10)
// into the v3 registered denoms keyset. Must run before migrateRollapps, which
// rewrites rollapp objects and drops the reserved field.
func migrateRollappRegisteredDenoms(ctx sdk.Context, rk *rollappkeeper.Keeper) error {
	var migrateErr error
	rk.IterateRollappBytes(ctx, func(value []byte) bool {
		var legacy legacyrollapp.LegacyRollapp
		if err := legacy.Unmarshal(value); err != nil {
			migrateErr = fmt.Errorf("unmarshal legacy rollapp: %w", err)
			return true
		}
		if legacy.RollappId == "" || len(legacy.RegisteredDenoms) == 0 {
			return false
		}
		for _, denom := range legacy.RegisteredDenoms {
			if denom == "" {
				continue
			}
			if err := rk.SetRegisteredDenom(ctx, legacy.RollappId, denom); err != nil {
				migrateErr = fmt.Errorf("set registered denom %s/%s: %w", legacy.RollappId, denom, err)
				return true
			}
		}
		ctx.Logger().Info("migrated rollapp registered denoms",
			"rollapp_id", legacy.RollappId,
			"count", len(legacy.RegisteredDenoms),
		)
		return false
	})
	return migrateErr
}
