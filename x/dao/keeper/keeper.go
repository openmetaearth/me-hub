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

func (k Keeper) SetDaoAddresses(ctx sdk.Context, daoAddresses types.DaoAddresses) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&daoAddresses)
	store.Set(types.DaoAddressesPrefix, b)
}

// GetDaoAddresses returns dao addresses
func (k Keeper) GetDaoAddresses(ctx sdk.Context) (dao types.DaoAddresses, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.DaoAddressesPrefix)
	if b == nil {
		return dao, false
	}
	k.cdc.MustUnmarshal(b, &dao)
	return dao, true
}

func (k Keeper) GetGlobalDao(ctx sdk.Context) string {
	dao, found := k.GetDaoAddresses(ctx)
	if found {
		return dao.GlobalDao
	}
	return ""
}

func (k Keeper) GetMeidDao(ctx sdk.Context) string {
	dao, found := k.GetDaoAddresses(ctx)
	if found {
		return dao.MeidDao
	}
	return ""
}

func (k Keeper) GetDevOperator(ctx sdk.Context) string {
	dao, found := k.GetDaoAddresses(ctx)
	if found {
		return dao.DevOperator
	}
	return ""
}

func (k Keeper) GetAirdropAddress(ctx sdk.Context) string {
	dao, found := k.GetDaoAddresses(ctx)
	if found {
		return dao.AirdropAddress
	}
	return ""
}

func (k Keeper) GetValidatorAddress(ctx sdk.Context) string {
	dao, found := k.GetDaoAddresses(ctx)
	if found {
		return dao.ValidatorAddress
	}
	return ""
}

func (k Keeper) IsGlobalDao(ctx sdk.Context, address string) bool {
	dao, found := k.GetDaoAddresses(ctx)
	if !found {
		return false
	}
	return dao.GlobalDao == address
}

func (k Keeper) IsMeidDao(ctx sdk.Context, address string) bool {
	dao, found := k.GetDaoAddresses(ctx)
	if !found {
		return false
	}
	return dao.MeidDao == address
}

func (k Keeper) IsValidatorDao(ctx sdk.Context, address string) bool {
	dao, found := k.GetDaoAddresses(ctx)
	if !found {
		return false
	}
	return dao.ValidatorAddress == address
}
