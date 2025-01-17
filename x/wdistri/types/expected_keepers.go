package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	dt "github.com/cosmos/cosmos-sdk/x/distribution/types"
	wstakingtypes "github.com/st-chain/me-hub/x/wstaking/types"
)

type DistrKeeper interface {
	// Methods imported from distr should be defined here
}

type StakingKeeper interface {
	// Methods imported from staking should be defined here
	dt.StakingKeeper

	GetAllRegionI(ctx sdk.Context) (list []wstakingtypes.RegionI)
}

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
	dt.AccountKeeper
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
	dt.BankKeeper
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
