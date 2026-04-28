package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// Keeper is wrapper of did keeper and nft keeper.
type Keeper struct {
	moduleName string

	cdc           codec.Codec
	storeKey      storetypes.StoreKey
	bankKeeper    types.BankKeeper
	accountKeeper authkeeper.AccountKeeper
	daoKeeper     types.DaoKeeper
	authority     string
}

func NewKeeper(
	moduleName string,
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	accountKeeper authkeeper.AccountKeeper,
	daoKeeper types.DaoKeeper,
	authority string,
) Keeper {
	return Keeper{
		moduleName:    moduleName,
		cdc:           cdc,
		storeKey:      storeKey,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		daoKeeper:     daoKeeper,
		authority:     authority,
	}
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k *Keeper) Codec() codec.Codec {
	return k.cdc
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", k.moduleName))
}

func (k Keeper) ModuleName() string {
	return k.moduleName
}
