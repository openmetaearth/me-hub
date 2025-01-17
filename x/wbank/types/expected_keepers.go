package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type DaoKeeper interface {
	IsGlobalDao(ctx sdk.Context, address string) bool
	IsMeidDao(ctx sdk.Context, address string) bool
	GetAirdropAddress(ctx sdk.Context) string
	GetDevOperator(ctx sdk.Context) string
}
