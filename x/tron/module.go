package tron

import (
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	gravitycli "github.com/openmetaearth/me-hub/x/gravity/client/cli"
	gravitykeeper "github.com/openmetaearth/me-hub/x/gravity/keeper"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
	"github.com/openmetaearth/me-hub/x/tron/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule         = AppModule{}
	_ module.AppModuleBasic    = AppModuleBasic{}
	_ module.EndBlockAppModule = AppModule{}
)

// AppModuleBasic object for module implementation
type AppModuleBasic struct{}

// Name implements app module basic
func (AppModuleBasic) Name() string { return types.ModuleName }

// RegisterLegacyAminoCodec implements app module basic
func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

// DefaultGenesis implements app module basic
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis implements app module basic
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, data json.RawMessage) error {
	var state gravitytypes.GenesisState
	if err := cdc.UnmarshalJSON(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return state.ValidateBasic()
}

// GetQueryCmd implements app module basic
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return gravitycli.GetQueryCmd(types.ModuleName)
}

// GetTxCmd implements app module basic
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return gravitycli.GetTxCmd(types.ModuleName)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway
func (AppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

// RegisterInterfaces implements app bmodule basic
func (AppModuleBasic) RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule object for module implementation
type AppModule struct {
	AppModuleBasic
	keeper gravitykeeper.Keeper
}

// NewAppModule creates a new AppModule Object
func NewAppModule(keeper gravitykeeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// RegisterInvariants implements app module
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	migrator := gravitykeeper.NewMigrator(am.keeper)
	if err := cfg.RegisterMigration(am.Name(), 1, migrator.Migrate); err != nil {
		panic(fmt.Sprintf("failed to register migration for %s: %v", am.Name(), err))
	}
}

// InitGenesis initializes the genesis state for this module and implements app module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState gravitytypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	gravitykeeper.InitGenesis(ctx, am.keeper, &genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports the current genesis state to a json.RawMessage
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	state := gravitykeeper.ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(state)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 {
	return 1
}

// EndBlock implements app module
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	am.keeper.EndBlocker(ctx)
	return []abci.ValidatorUpdate{}
}
