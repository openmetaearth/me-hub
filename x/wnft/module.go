package wnft

import (
	"context"
	"encoding/json"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/x/nft"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	nftmodule "cosmossdk.io/x/nft/module"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/wnft/client/cli"
	"github.com/st-chain/me-hub/x/wnft/keeper"
	"github.com/st-chain/me-hub/x/wnft/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.AppModule      = AppModule{}
	_ module.HasGenesis     = AppModule{}

	_ appmodule.AppModule   = AppModule{}
	_ appmodule.HasServices = AppModule{}
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
	keeper        *keeper.Keeper
	accountKeeper nft.AccountKeeper
	bankKeeper    nft.BankKeeper
	registry      codectypes.InterfaceRegistry
}

// NewAppModule creates a new wnft AppModule object.
func NewAppModule(
	cdc codec.Codec,
	k keeper.Keeper,
	ak nft.AccountKeeper,
	bk nft.BankKeeper,
	registry codectypes.InterfaceRegistry,
) AppModule {
	nftModule := nftmodule.NewAppModule(cdc, k.Keeper, ak, bk, registry)
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc, nftModule.AppModuleBasic),
		keeper:         &k,
		accountKeeper:  ak,
		bankKeeper:     bk,
		registry:       registry,
	}
}

// IsOnePerModuleType implements depinject.OnePerModuleType.
func (AppModule) IsOnePerModuleType() {}

// IsAppModule implements appmodule.AppModule.
func (AppModule) IsAppModule() {}

// RegisterServices registers module gRPC services.
// Implements appmodule.HasServices.
func (am AppModule) RegisterServices(registrar grpc.ServiceRegistrar) error {
	types.RegisterMsgServer(registrar, keeper.NewMsgServerImpl(am.keeper, am.keeper.Keeper))
	types.RegisterQueryServer(registrar, am.keeper)
	nft.RegisterQueryServer(registrar, am.keeper.Keeper)
	return nil
}

// InitGenesis performs genesis initialization for the wnft module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var genesisState nft.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.InitGenesis(ctx, &genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the wnft module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// GetQueryCmd returns the root query command for the wnft module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// GetTxCmd returns the root tx command for the wnft module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// RegisterInterfaces registers the module's interface types.
func (a AppModuleBasic) RegisterInterfaces(reg codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
	nft.RegisterInterfaces(reg)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := nft.RegisterQueryHandlerClient(context.Background(), mux, nft.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// RegisterLegacyAminoCodec registers legacy amino codec.
// Note: cosmossdk.io/x/nft v0.1.1+ no longer exposes RegisterCodec,
// so only wnft-specific types are registered here.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// nftKeeper returns the embedded nft keeper (for use by other modules if needed).
func (am AppModule) NftKeeper() nftkeeper.Keeper {
	return am.keeper.Keeper
}
