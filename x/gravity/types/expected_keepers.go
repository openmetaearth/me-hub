package types

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stktypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

type StakingKeeper interface {
	GetRegion(ctx sdk.Context, regionId string) (val stktypes.Region, found bool)
	GetAllRegion(ctx sdk.Context) (regions []stktypes.Region)
	KycReward(ctx sdk.Context, account sdk.AccAddress, regionId, creator string) error
	RemoveKycReward(ctx sdk.Context, account sdk.AccAddress, regionId string) error
	TransferKycRegion(ctx sdk.Context, address sdk.AccAddress, creator, fromRegionId, toRegionId string) error
	SendInviteReward(ctx sdk.Context, inviter, invitee, regionId string) error
}

// BankKeeper defines the expected bank keeper methods
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx context.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, name string, amt sdk.Coins) error
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	GetSupply(ctx context.Context, denom string) sdk.Coin
	IterateAllDenomMetaData(ctx context.Context, cb func(banktypes.Metadata) bool)
	SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata)
}

type (
	ParamSet = paramtypes.ParamSet
	// Subspace defines an interface that implements the legacy x/params Subspace
	// type.
	//
	// NOTE: This is used solely for migration of x/params managed parameters.
	Subspace interface {
		GetParamSet(ctx sdk.Context, ps ParamSet)
		HasKeyTable() bool
		WithKeyTable(table paramtypes.KeyTable) paramtypes.Subspace
	}
)

type DaoKeeper interface {
	IsDao(ctx sdk.Context, address string) bool
}
