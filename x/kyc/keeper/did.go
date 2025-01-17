package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/st-chain/me-hub/x/did/types"
)

/*
These methods wrap the did methods without any other logic
*/

func (k *Keeper) HasDID(ctx sdk.Context, addr sdk.AccAddress) bool {
	return k.didKeeper.HasDID(ctx, addr)
}

func (k *Keeper) GetDID(ctx sdk.Context, addr sdk.AccAddress) (string, bool) {
	return k.didKeeper.GetDID(ctx, addr)
}

func (k *Keeper) SetDID(ctx sdk.Context, addr sdk.AccAddress, did string) {
	k.didKeeper.SetDID(ctx, addr, did)
}

func (k *Keeper) HasDidInfo(ctx sdk.Context, did string) bool {
	return k.didKeeper.HasDidInfo(ctx, did)
}

func (k *Keeper) GetDidInfo(ctx sdk.Context, did string) (didtypes.DidInfo, bool) {
	return k.didKeeper.GetDidInfo(ctx, did)
}

func (k *Keeper) SetDidInfo(ctx sdk.Context, did string, info didtypes.DidInfo) {
	k.didKeeper.SetDidInfo(ctx, did, info)
}
