package wgov

import (
	"context"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/openmetaearth/me-hub/x/wgov/client/cli"
	wgovkeeper "github.com/openmetaearth/me-hub/x/wgov/keeper"
	"github.com/openmetaearth/me-hub/x/wgov/types"
	"github.com/spf13/cobra"
)

// AppModuleBasic implements the basic application module for the wrapped nft module.
type AppModuleBasic struct {
	gov.AppModuleBasic
}

func NewAppModuleBasic(legacyProposalHandlers []govclient.ProposalHandler) AppModuleBasic {
	return AppModuleBasic{
		AppModuleBasic: gov.NewAppModuleBasic(legacyProposalHandlers),
	}
}

// AppModule implements an application module for the wnft module.
type AppModule struct {
	gov.AppModule
	keeper         *wgovkeeper.Keeper
	accountKeeper  govtypes.AccountKeeper
	legacySubspace govtypes.ParamSubspace
}

func NewAppModule(cdc codec.Codec, keeper *wgovkeeper.Keeper, ak govtypes.AccountKeeper, bk govtypes.BankKeeper, ss govtypes.ParamSubspace) AppModule {
	return AppModule{
		AppModule:      gov.NewAppModule(cdc, &keeper.Keeper, ak, bk, ss),
		keeper:         keeper,
		accountKeeper:  ak,
		legacySubspace: ss,
	}
}

func (a AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	a.AppModuleBasic.RegisterLegacyAminoCodec(cdc)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the staking module.
func (a AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	a.AppModuleBasic.RegisterGRPCGatewayRoutes(clientCtx, mux)
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	msgServer := govkeeper.NewMsgServerImpl(&am.keeper.Keeper)
	v1beta1.RegisterMsgServer(cfg.MsgServer(), govkeeper.NewLegacyMsgServerImpl(am.accountKeeper.GetModuleAddress(govtypes.ModuleName).String(), msgServer))
	v1.RegisterMsgServer(cfg.MsgServer(), msgServer)

	legacyQueryServer := govkeeper.NewLegacyQueryServer(&am.keeper.Keeper)
	v1beta1.RegisterQueryServer(cfg.QueryServer(), legacyQueryServer)
	v1.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	wgovQuerier := wgovkeeper.Querier{Keeper: am.keeper}
	types.RegisterQueryServer(cfg.QueryServer(), wgovQuerier)

	m := govkeeper.NewMigrator(&am.keeper.Keeper, am.legacySubspace)
	err := cfg.RegisterMigration(govtypes.ModuleName, 1, m.Migrate1to2)
	if err != nil {
		panic(fmt.Sprintf("failed to migrate x/gov from version 1 to 2: %v", err))
	}
	err = cfg.RegisterMigration(govtypes.ModuleName, 2, m.Migrate2to3)
	if err != nil {
		panic(fmt.Sprintf("failed to migrate x/gov from version 2 to 3: %v", err))
	}
	err = cfg.RegisterMigration(govtypes.ModuleName, 3, m.Migrate3to4)
	if err != nil {
		panic(fmt.Sprintf("failed to migrate x/gov from version 3 to 4: %v", err))
	}
}
