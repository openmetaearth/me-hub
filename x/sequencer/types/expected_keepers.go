package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

// RollappKeeper defines the expected rollapp keeper used for retrieve rollapp.
type RollappKeeper interface {
	GetRollapp(ctx sdk.Context, rollappId string) (val rollapptypes.Rollapp, found bool)
	MinBond(ctx sdk.Context, rollappID string) sdk.Coin
	MustGetRollapp(ctx sdk.Context, rollappId string) rollapptypes.Rollapp
	GetAllRollapps(ctx sdk.Context) (list []rollapptypes.Rollapp)
	SetRollappAsLaunched(ctx sdk.Context, rollapp *rollapptypes.Rollapp) error
	HardForkToLatest(ctx sdk.Context, rollappId string) error
	ForkLatestAllowed(ctx sdk.Context, rollappId string) bool
	GetLatestStateInfoIndex(ctx sdk.Context, rollappId string) (val rollapptypes.StateInfoIndex, found bool)
	GetStateInfo(ctx sdk.Context, rollappId string, index uint64) (val rollapptypes.StateInfo, found bool)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	BurnCoins(ctx context.Context, name string, amt sdk.Coins) error
}
