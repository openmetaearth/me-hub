package wnft

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftmodule "github.com/cosmos/cosmos-sdk/x/nft/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/openmetaearth/me-hub/x/wnft/client/cli"
	"github.com/openmetaearth/me-hub/x/wnft/keeper"
	"github.com/openmetaearth/me-hub/x/wnft/types"
	"github.com/spf13/cobra"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.AppModule      = AppModule{}
)

// AppModuleBasic implements the basic application module for the wrapped nft module.
type AppModuleBasic struct {
	nftmodule.AppModuleBasic
	cdc codec.Codec
}

func NewAppModuleBasic(cdc codec.Codec, basic nftmodule.AppModuleBasic) AppModuleBasic {
	return AppModuleBasic{cdc: cdc, AppModuleBasic: basic}
}

// AppModule implements an application module for the wnft module.
type AppModule struct {
	AppModuleBasic
	nftmodule.AppModule
	keeper *keeper.Keeper
}

// NewAppModule creates a new wnft AppModule object.
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	ak nft.AccountKeeper,
	bk nft.BankKeeper,
	registry codectypes.InterfaceRegistry,
) AppModule {
	nftModule := nftmodule.NewAppModule(cdc, keeper.Keeper, ak, bk, registry)
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc, nftModule.AppModuleBasic),
		AppModule:      nftModule,
		keeper:         &keeper,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper, am.keeper.Keeper))

	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	nft.RegisterQueryServer(cfg.QueryServer(), am.keeper.Keeper)
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

// RegisterRESTRoutes registers the capability module's REST service handlers.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := nft.RegisterQueryHandlerClient(context.Background(), mux, nft.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
	// nolint: errcheck, gosec
	types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {
	nft.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	nft.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
}
