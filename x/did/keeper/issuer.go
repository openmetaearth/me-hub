package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/did/types"
)

func (k Keeper) GetIssuer(ctx sdk.Context, did string) (issuer types.Issuer, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetIssuerKey(did))
	if bz == nil {
		return types.Issuer{}, false
	}

	k.cdc.MustUnmarshal(bz, &issuer)
	return issuer, true
}

func (k Keeper) SetIssuer(ctx sdk.Context, issuer types.Issuer) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIssuerKey(issuer.Did), k.cdc.MustMarshal(&issuer))
}
