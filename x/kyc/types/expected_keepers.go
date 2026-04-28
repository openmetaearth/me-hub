package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/nft"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
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

type DIDKeeper interface {
	HasDID(ctx sdk.Context, addr sdk.AccAddress) bool
	GetDID(ctx sdk.Context, addr sdk.AccAddress) (did string, found bool)
	SetDID(ctx sdk.Context, addr sdk.AccAddress, did string)
	DeleteDID(ctx sdk.Context, addr sdk.AccAddress)

	HasDidInfo(ctx sdk.Context, did string) bool
	GetDidInfo(ctx sdk.Context, did string) (info didtypes.DidInfo, found bool)
	SetDidInfo(ctx sdk.Context, did string, info didtypes.DidInfo)

	GetService(ctx sdk.Context, sid string) (service didtypes.Service, found bool)
	SetService(ctx sdk.Context, sid string, svc didtypes.Service)

	HasCredential(ctx sdk.Context, did string, sid string) bool
	GetCredential(ctx sdk.Context, did, sid string) (vc didtypes.Credential, found bool)
	GetCredentialsByFilter(ctx sdk.Context, sid string, filter []byte, pageReq *query.PageRequest) ([]didtypes.Credential, *query.PageResponse, error)
	SetCredential(ctx sdk.Context, did, sid string, credential didtypes.Credential)
	DeleteCredential(ctx sdk.Context, did, sid string)

	GetFilters(ctx sdk.Context, did, sid string) (filters [][]byte, found bool)
	AddFilters(ctx sdk.Context, did, sid string, filters [][]byte, vc didtypes.Credential)
	DeleteFilters(ctx sdk.Context, did, sid string, filters [][]byte)
}

type NFTKeeper interface {
	GetNFT(ctx sdk.Context, classID, nftID string) (nft.NFT, bool)
	HasNFT(ctx sdk.Context, classID, id string) bool
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
	Mint(ctx sdk.Context, token nft.NFT, receiver sdk.AccAddress) error
	Update(ctx sdk.Context, token nft.NFT) error
	Burn(ctx sdk.Context, classID string, nftID string) error
	SaveClass(ctx sdk.Context, class nft.Class) error
}
