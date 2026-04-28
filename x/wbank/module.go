package wbank

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/openmetaearth/me-hub/x/wbank/keeper"
)

// AppModuleBasic defines the basic application module used by the wrapped bank module.
type AppModuleBasic struct {
	bank.AppModuleBasic
}

// AppModule implements an application module for the wrapped bank module.
type AppModule struct {
	bank.AppModule
	keeper         keeper.BaseKeeperWrapper
	legacySubspace bankexported.Subspace
}

// NewAppModule creates a new bank AppModule object.
func NewAppModule(
	cdc codec.Codec, keeper keeper.BaseKeeperWrapper, accountKeeper banktypes.AccountKeeper, ss bankexported.Subspace,
) AppModule {
	bankModule := bank.NewAppModule(cdc, keeper.BaseKeeper, accountKeeper, ss)
	return AppModule{
		AppModule:      bankModule,
		keeper:         keeper,
		legacySubspace: ss,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// copied the bank's RegisterServices to replace with the keeper wrapper
	bankMsgSrv := bankkeeper.NewMsgServerImpl(am.keeper)
	banktypes.RegisterMsgServer(cfg.MsgServer(), bankMsgSrv)
	banktypes.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	m := bankkeeper.NewMigrator(am.keeper.BaseKeeper, am.legacySubspace)
	if err := cfg.RegisterMigration(banktypes.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/bank from version 1 to 2: %v", err))
	}

	if err := cfg.RegisterMigration(banktypes.ModuleName, 2, m.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/bank from version 2 to 3: %v", err))
	}

	if err := cfg.RegisterMigration(banktypes.ModuleName, 3, m.Migrate3to4); err != nil {
		panic(fmt.Sprintf("failed to migrate x/bank from version 3 to 4: %v", err))
	}
}
