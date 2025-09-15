package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stktypes "github.com/st-chain/me-hub/x/wstaking/types"
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
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	IterateAllDenomMetaData(ctx sdk.Context, cb func(banktypes.Metadata) bool)
}
