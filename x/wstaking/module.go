package wstaking

import (
	"context"
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/x/wstaking/client/cli"
	"github.com/spf13/cobra"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/openmetaearth/me-hub/x/wstaking/keeper"
	"github.com/openmetaearth/me-hub/x/wstaking/types"

	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
)

// AppModuleBasic defines the basic application module used by the wrapped staking module.
type AppModuleBasic struct {
	staking.AppModuleBasic
}

// AppModule implements an application module for the wrapped staking module.
type AppModule struct {
	staking.AppModule
	keeper            *keeper.Keeper
	IbcTransferKeeper ibctransferkeeper.Keeper
	legacySubspace    stakingexported.Subspace
}

// NewAppModule creates a new AppModule object.
func NewAppModule(
	cdc codec.Codec,
	keeper *keeper.Keeper,
	ibcTransferKeeper ibctransferkeeper.Keeper,
	ak stakingtypes.AccountKeeper,
	bk stakingtypes.BankKeeper,
	ls stakingexported.Subspace,
) AppModule {
	stakingAppModule := staking.NewAppModule(cdc, keeper.Keeper, ak, bk, ls)

	return AppModule{
		AppModule:         stakingAppModule,
		keeper:            keeper,
		IbcTransferKeeper: ibcTransferKeeper,
		legacySubspace:    ls,
	}
}

func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
	stakingtypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
	stakingtypes.RegisterInterfaces(reg)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the staking module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := stakingtypes.RegisterQueryHandlerClient(context.Background(), mux, stakingtypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// DefaultGenesis returns default genesis state as raw bytes for the staking
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the staking module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return types.ValidateGenesis(&data)
}

// GetTxCmd returns the root tx command for the staking module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns no root query command for the staking module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// InitGenesis performs genesis initialization for the staking module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	return am.keeper.InitGenesis(ctx, &genesisState)
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(am.keeper.ExportGenesis(ctx))
}

// BeginBlock returns the begin blocker for the staking module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlock(ctx, am.keeper)
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlock(ctx, am.keeper)
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	stakingKeeperMsgSrv := stakingkeeper.NewMsgServerImpl(am.keeper.Keeper)
	// wrap the staking keeper message server to intersect the messages
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper, am.IbcTransferKeeper, stakingKeeperMsgSrv))
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper, am.IbcTransferKeeper, stakingKeeperMsgSrv))

	querier := stakingkeeper.Querier{Keeper: am.keeper.Keeper}
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), querier)
	wstakingQuerier := keeper.Querier{Keeper: am.keeper}
	types.RegisterQueryServer(cfg.QueryServer(), wstakingQuerier)

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

// RegisterInvariants registers the staking module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	keeper.RegisterInvariants(ir, am.keeper)
}
