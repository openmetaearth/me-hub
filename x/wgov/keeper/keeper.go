package keeper

import (
	corestoretypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/st-chain/me-hub/x/wgov/types"
)

type Keeper struct {
	govkeeper.Keeper
	stakingKeeper types.StakingKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	accountKeeper govtypes.AccountKeeper,
	bankKeeper govtypes.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distributionKeeper govtypes.DistributionKeeper,
	router *baseapp.MsgServiceRouter,
	config govtypes.Config,
	authority string,
) *Keeper {
	return &Keeper{
		Keeper:        *govkeeper.NewKeeper(cdc, storeService, accountKeeper, bankKeeper, stakingKeeper, distributionKeeper, router, config, authority),
		stakingKeeper: stakingKeeper,
	}
}
