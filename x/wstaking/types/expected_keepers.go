package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	"math/big"
)

type DaoKeeper interface {
	IsGlobalDao(ctx sdk.Context, address string) bool
	IsMeidDao(ctx sdk.Context, address string) bool
	GetAirdropAddress(ctx sdk.Context) string
	GetDevOperator(ctx sdk.Context) string
	GetGlobalDaoFeePoolAddr(ctx sdk.Context) sdk.AccAddress
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	GetSupply(ctx sdk.Context, denom string) sdk.Coin

	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	InputOutputCoins(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error

	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error

	// Methods imported from bank should be defined here
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error

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

type GroupKeeper interface {
	CreateGroupByRegion(sdkCtx sdk.Context, regionInfo Region) (uint64, error)
}
