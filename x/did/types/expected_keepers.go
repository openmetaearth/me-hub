package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DaoKeeper interface {
	IsGlobalDao(ctx sdk.Context, address string) bool
}
