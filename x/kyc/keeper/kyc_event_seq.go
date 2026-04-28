package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

// SetKycEventSeq set kycEventSeq in the store
func (k Keeper) SetKycEventSeq(ctx sdk.Context, kycEventSeq types.KycEventSeq) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.KycEventSeqKey))
	b := k.cdc.MustMarshal(&kycEventSeq)
	store.Set([]byte{0}, b)
}

// GetKycEventSeq returns kycEventSeq
func (k Keeper) GetKycEventSeq(ctx sdk.Context) (val types.KycEventSeq, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.KycEventSeqKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveKycEventSeq removes kycEventSeq from the store
func (k Keeper) RemoveKycEventSeq(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.KycEventSeqKey))
	store.Delete([]byte{0})
}

func (k Keeper) takeSeq(ctx sdk.Context) uint64 {
	val, _ := k.GetKycEventSeq(ctx)
	next := val.Seq + 1
	k.SetKycEventSeq(ctx, types.KycEventSeq{Seq: next})
	return val.Seq
}
