package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

// RollappKeeper defines the expected rollapp keeper used for retrieve rollapp.
type RollappKeeper interface {
	GetRollapp(ctx sdk.Context, rollappId string) (val rollapptypes.Rollapp, found bool)
	GetAllRollapps(ctx sdk.Context) (list []rollapptypes.Rollapp)
	GetLatestStateInfoIndex(ctx sdk.Context, rollappId string) (rollapptypes.StateInfoIndex, bool)
	GetStateInfo(ctx sdk.Context, rollappId string, index uint64) (rollapptypes.StateInfo, bool)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}
