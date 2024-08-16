package types

import (
	cmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	dt "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/st-chain/me-hub/mocks/mock"
)

type DistrKeeper interface {
	// Methods imported from distr should be defined here
}

type StakingKeeper interface {
	// Methods imported from staking should be defined here
	dt.StakingKeeper

	//FIXME: some method are not need ,remove it
	CalculateInterest(ctx sdk.Context, totalStaking cmath.Int, height int64) (rewards sdk.Dec, err error)
	SetDelegation(ctx sdk.Context, delegation mock.Delegation)
	// GetALLDelegationAmount(ctx sdk.Context) sdk.Dec
	//FIXME: replace mock.MockRegion type
	GetAllRegion(ctx sdk.Context) (list []mock.MockRegion)
	// GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)

	GetMeid(ctx sdk.Context, account string) (val mock.MockMeid, found bool)
	//FIXME: replace mock.MockRegion type
	GetRegion(ctx sdk.Context, regionId string) (val mock.MockRegion, found bool)
	SetRegion(ctx sdk.Context, region mock.MockRegion)
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
