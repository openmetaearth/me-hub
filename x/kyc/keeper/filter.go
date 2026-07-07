package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

/*
KYC filter
*/

func (k *Keeper) GetFilters(ctx sdk.Context, did string) (filters [][]byte, found bool) {
	return k.didKeeper.GetFilters(ctx, did, types.ModuleName)
}

func (k *Keeper) AddFilters(ctx sdk.Context, did string, filters [][]byte, vc didtypes.Credential) {
	k.didKeeper.AddFilters(ctx, did, types.ModuleName, filters, vc)
}

func (k *Keeper) DeleteFilters(ctx sdk.Context, did string, filters [][]byte) {
	k.didKeeper.DeleteFilters(ctx, did, types.ModuleName, filters)
}
