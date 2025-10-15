package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected bank keeper.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type DaoKeeper interface {
	IsGlobalDao(ctx sdk.Context, address string) bool
	IsMeidDao(ctx sdk.Context, address string) bool
	GetAirdropAddress(ctx sdk.Context) string
	GetDevOperator(ctx sdk.Context) string
}

type BlacklistKeeper interface {
	IsBlackList(ctx sdk.Context, addr sdk.AccAddress) bool
}
