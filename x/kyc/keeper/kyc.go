package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

/*
KYC
*/

func (k *Keeper) HasKYC(ctx sdk.Context, did string) (found bool) {
	return k.didKeeper.HasCredential(ctx, did, types.ModuleName)
}

func (k *Keeper) GetKYC(ctx sdk.Context, did string) (kyc didtypes.Credential, found bool) {
	return k.didKeeper.GetCredential(ctx, did, types.ModuleName)
}

func (k *Keeper) GetKYCsByRegion(
	ctx sdk.Context,
	regionId string,
	pageReq *query.PageRequest,
) (KYCs []didtypes.Credential, pageRes *query.PageResponse, err error) {
	return k.didKeeper.GetCredentialsByFilter(ctx, types.ModuleName, []byte(regionId), pageReq)
}

func (k *Keeper) SetKYC(ctx sdk.Context, did string, kyc didtypes.Credential) {
	k.didKeeper.SetCredential(ctx, did, types.ModuleName, kyc)
}

func (k *Keeper) DeleteKYC(ctx sdk.Context, did string) {
	k.didKeeper.DeleteCredential(ctx, did, types.ModuleName)
}
