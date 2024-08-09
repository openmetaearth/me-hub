package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// CustomParamsKeeper defines the custom params keeper interface required for the module.
type DaoKeeper interface {
	IsGlobalDao(ctx sdk.Context, address string) bool
	IsMeidDao(ctx sdk.Context, address string) bool
}
