package wstaking

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/dymensionxyz/dymension/v3/x/wstaking/keeper"
)

// AppModule implements an application module for the wrapped staking module.
type AppModule struct {
	staking.AppModule
	keeper         *keeper.KeeperWrapper
	legacySubspace stakingexported.Subspace
}

// NewAppModule creates a new AppModule object.
func NewAppModule(
	cdc codec.Codec,
	keeper *keeper.KeeperWrapper,
	ak stakingtypes.AccountKeeper,
	bk stakingtypes.BankKeeper,
	ls stakingexported.Subspace,
) AppModule {
	stakingAppModule := staking.NewAppModule(cdc, keeper.Keeper, ak, bk, ls)

	return AppModule{
		AppModule:      stakingAppModule,
		keeper:         keeper,
		legacySubspace: ls,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	stakingKeeperMsgSrv := stakingkeeper.NewMsgServerImpl(am.keeper.Keeper)
	// wrap the staking keeper message server to intersect the messages
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(stakingKeeperMsgSrv))
	querier := stakingkeeper.Querier{Keeper: am.keeper.Keeper}
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), querier)

	m := stakingkeeper.NewMigrator(am.keeper.Keeper, am.legacySubspace)
	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", stakingtypes.ModuleName, err))
	}
	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 2, m.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 2 to 3: %v", stakingtypes.ModuleName, err))
	}
	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 3, m.Migrate3to4); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 3 to 4: %v", stakingtypes.ModuleName, err))
	}
}
