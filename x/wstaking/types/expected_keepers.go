package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/wbank/keeper"
)

type DaoKeeper interface {
	IsDao(ctx sdk.Context, address string) bool
	IsGlobalDao(ctx sdk.Context, address string) bool
	IsMeidDao(ctx sdk.Context, address string) bool
	GetAirdropAddress(ctx sdk.Context) string
	GetDevOperator(ctx sdk.Context) string
	GetGlobalDaoFeePoolAddr(ctx sdk.Context) sdk.AccAddress
	GetGlobalDao(ctx sdk.Context) string
	GetMeidDao(ctx sdk.Context) string
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	GetSupply(ctx sdk.Context, denom string) sdk.Coin

	Extend() keeper.BankKeeperExtend

	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	InputOutputCoins(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error

	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error

	StakeCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	UnstakeCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
}

type MintKeeper interface {
	// GetPerBlockMintCoinAmount returns the current block mint coins amount.
	GetPerBlockMintCoinAmount(ctx sdk.Context) (amount big.Int)
}

type NFTKeeper interface {
	SaveClass(ctx sdk.Context, class nft.Class) error
	GetClass(ctx sdk.Context, classID string) (nft.Class, bool)
}

type KycKeeper interface {
	GetDID(ctx sdk.Context, addr sdk.AccAddress) (string, bool)
	GetKYC(ctx sdk.Context, did string) (kyc didtypes.Credential, found bool)
}

type DidKeeper interface {
	IteratorCredentialsByFilter(ctx sdk.Context, sid string, filter []byte, cb func(delegation didtypes.Credential) (stop bool))
	GetDidInfo(ctx sdk.Context, did string) (info didtypes.DidInfo, found bool)
}

type GroupKeeper interface {
	CreateGroupByRegion(sdkCtx sdk.Context, regionInfo Region) (uint64, error)
	UpdateGroupAdmin(ctx sdk.Context, regionID string, admin string)
}
