package wnft

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftmodule "github.com/cosmos/cosmos-sdk/x/nft/module"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/wnft/client/cli"
	"github.com/st-chain/me-hub/x/wnft/keeper"
	"github.com/st-chain/me-hub/x/wnft/types"
)

var (
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModule           = AppModule{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic implements the basic application module for the wrapped nft module.
type AppModuleBasic struct {
	nftmodule.AppModuleBasic
}

// AppModule implements an application module for the wnft module.
type AppModule struct {
	nftmodule.AppModule
	keeper *keeper.Keeper
}

// NewAppModule creates a new wnft AppModule object.
func NewAppModule(
	cdc codec.Codec,
	keeper *keeper.Keeper,
	ak nft.AccountKeeper,
	bk nft.BankKeeper,
	registry codectypes.InterfaceRegistry,
) AppModule {
	nftModule := nftmodule.NewAppModule(cdc, keeper.Keeper, ak, bk, registry)
	return AppModule{
		AppModule: nftModule,
		keeper:    keeper,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper, am.keeper.Keeper))
	querier := keeper.Querier{Keeper: am.keeper}
	types.RegisterQueryServer(cfg.QueryServer(), querier)
}

// GetQueryCmd returns no root query command for the staking module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// RegisterInterfaces registers the module's interface types
func (a AppModuleBasic) RegisterInterfaces(reg codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
	nft.RegisterInterfaces(reg)
}
