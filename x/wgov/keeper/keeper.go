package keeper

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/openmetaearth/me-hub/x/wgov/types"
)

type Keeper struct {
	govkeeper.Keeper
	storeKey      storetypes.StoreKey
	stakingKeeper types.StakingKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey,
	accountKeeper govtypes.AccountKeeper, bankKeeper govtypes.BankKeeper, stakingKeeper types.StakingKeeper,
	router *baseapp.MsgServiceRouter,
	config govtypes.Config,
	authority string,
) *Keeper {
	return &Keeper{
		Keeper:        *govkeeper.NewKeeper(cdc, key, accountKeeper, bankKeeper, stakingKeeper, router, config, authority),
		storeKey:      key,
		stakingKeeper: stakingKeeper,
	}
}
