package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/dao/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetGlobalDao(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GlobalDaoPrefix, address)
}

func (k Keeper) GetGlobalDao(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.GlobalDaoPrefix)
}

func (k Keeper) SetMeidDao(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.MeidDaoPrefix, address)
}

func (k Keeper) GetMeidDao(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.MeidDaoPrefix)
}

func (k Keeper) SetDevOperator(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DevOperatorPrefix, address)
}

func (k Keeper) GetDevOperator(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.DevOperatorPrefix)
}

func (k Keeper) IsGlobalDao(ctx sdk.Context, address string) bool {
	admin := k.GetGlobalDao(ctx)
	return admin.String() == address
}

func (k Keeper) IsMeidDao(ctx sdk.Context, address string) bool {
	admin := k.GetGlobalDao(ctx)
	return admin.String() == address
}
