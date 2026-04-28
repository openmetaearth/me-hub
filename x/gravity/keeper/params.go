package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// GetParams returns the parameters from the store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the parameters in the store
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
	if err := params.ValidateBasic(); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(params)
	store.Set(types.ParamsKey, bz)
	return nil
}

// GetGravityID returns the GravityID is essentially a salt value
// for bridge signatures, provided each chain running Gravity has a unique ID
// it won't be possible to play back signatures from one bridge onto another
// even if they share a relayer set.
//
// The lifecycle of the GravityID is that it is set in the Genesis file
// read from the live chain for the contract deployment, once a Gravity contract
// is deployed the GravityID CAN NOT BE CHANGED. Meaning that it can't just be the
// same as the chain id since the chain id may be changed many times with each
// successive chain in charge of the same bridge
func (k Keeper) GetGravityID(ctx sdk.Context) string {
	return k.GetParams(ctx).GravityId
}

func (k Keeper) GetGravityMinDelegate(ctx sdk.Context) sdk.Int {
	return k.GetParams(ctx).MinDelegate
}

func (k Keeper) GetGravityMaxDelegate(ctx sdk.Context) sdk.Int {
	return k.GetParams(ctx).MaxDelegate
}

func (k Keeper) GetSlashFraction(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).SlashFraction
}

func (k Keeper) GetSignedWindow(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).SignedWindow
}

func (k Keeper) GetMaxRelayers(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).MaxRelayers
}

func (k Keeper) GetRelayerSetUpdatePowerChangePercent(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).RelayerSetUpdatePowerChangePercent
}

func (k Keeper) MaxSlashTimes(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).MaxSlashTimes
}
