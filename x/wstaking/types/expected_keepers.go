package types

import (
	"context"
	"math/big"

	"cosmossdk.io/x/nft"
	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	"github.com/st-chain/me-hub/x/wbank/keeper"
)

type DaoKeeper interface {
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
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	LockedCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx context.Context, denom string) sdk.Coin
	Extend() keeper.BankKeeperExtend
	SendCoinsFromModuleToModule(ctx context.Context, senderPool, recipientPool string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, name string, amt sdk.Coins) error

	StakeCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	UnstakeCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
}

type MintKeeper interface {
	// GetPerBlockMintCoinAmount returns the current block mint coins amount.
	GetPerBlockMintCoinAmount(ctx sdk.Context) (amount big.Int)
}

type NFTKeeper interface {
	SaveClass(ctx context.Context, class nft.Class) error
	GetClass(ctx context.Context, classID string) (nft.Class, bool)
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
