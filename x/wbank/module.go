package wbank

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/api/tendermint/abci"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/st-chain/me-hub/x/wbank/types"

	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"

	"github.com/st-chain/me-hub/x/wbank/keeper"
	"github.com/st-chain/me-hub/x/wstaking/client/cli"
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

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState banktypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	bank.InitGenesis(am.keeper, ctx, &genesisState)
	return []abci.ValidatorUpdate{}
}

// RegisterInterfaces registers the module's interface types
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
	banktypes.RegisterInterfaces(reg)
}

// GetTxCmd returns the root tx command for the staking module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// copied the bank's RegisterServices to replace with the keeper wrapper
	bankMsgSrv := bankkeeper.NewMsgServerImpl(am.keeper.BaseKeeper)
	banktypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper, bankMsgSrv))
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper, bankMsgSrv))
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
