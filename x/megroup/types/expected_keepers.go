package types

import (
	//"me-hub/mocks/mock"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/st-chain/me-hub/x/kyc/handler"
	stakingTypes "github.com/st-chain/me-hub/x/wstaking/types"
)

type StakingKeeper interface {
	// Methods imported from staking should be defined here

	//FIXME: wait wstaking keep and types.Meid ;Temporarily use MockMeid instead of MEID and MockRegion instead types.Region
	// GetMeid(ctx sdk.Context, account string) (val types.Meid, found bool)
	GetMeid(ctx sdk.Context, account string) (val stakingTypes.Meid, found bool)
	// SetMeid(ctx sdk.Context, meid types.Meid)
	SetMeid(ctx sdk.Context, meid stakingTypes.Meid)

	// GetRegion(ctx sdk.Context, regionId string) (val types.Region, found bool)
	//	GetRegion(ctx sdk.Context, regionId string) (val mock.MockRegion, found bool)

	GetRegion(ctx sdk.Context, regionId string) (region stakingTypes.Region, found bool)
	CheckRegionName(name string) (string, error)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here

	// GetAccountAddressByID(ctx sdk.Context, int2 uint64) string
	GetModulePermissions() map[string]types.PermissionsForAddress
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetTreasuryPoolName() string
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type DAOKeeper interface {
	IsGlobalDao(ctx sdk.Context, address string) bool
	GetMeidDao(ctx sdk.Context) sdk.AccAddress
}

type KycKeeper interface {
	RegisterEventHandler(eventType string, priority int, module string, handler handler.HandlerFunc)
}
