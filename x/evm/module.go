package evm

import (
	"encoding/json"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/evmos/ethermint/x/evm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/openmetaearth/me-hub/x/evm/keeper"
	"github.com/spf13/cobra"
)

// AppModuleBasic implements the basic application module for the wrapped nft module.
type AppModuleBasic struct {
	evm.AppModuleBasic
}

// DefaultGenesis returns default genesis state as raw bytes for the evm
// module.
func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return b.AppModuleBasic.DefaultGenesis(cdc)
}

// ValidateGenesis is the validation check of the Genesis
func (b AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, txConfig client.TxEncodingConfig, bz json.RawMessage) error {
	return b.AppModuleBasic.ValidateGenesis(cdc, txConfig, bz)
}

// RegisterLegacyAminoCodec registers the evm module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	evmtypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers interfaces and implementations of the evm module.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	evmtypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(c client.Context, serveMux *runtime.ServeMux) {
	b.AppModuleBasic.RegisterGRPCGatewayRoutes(c, serveMux)
}

// GetTxCmd returns the root tx command for the evm module.
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return b.AppModuleBasic.GetTxCmd()
}

// GetQueryCmd returns no root query command for the evm module.
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return b.AppModuleBasic.GetQueryCmd()
}

// AppModule implements an application module for the wnft module.
type AppModule struct {
	AppModuleBasic
	evm.AppModule
	keeper         *keeper.Keeper
	accountKeeper  evmtypes.AccountKeeper
	bankKeeper     evmtypes.BankKeeper
	legacySubspace evmtypes.Subspace
}

func NewAppModule(k *keeper.Keeper, accountKeeper evmtypes.AccountKeeper, bankKeeper evmtypes.BankKeeper, legacySubspace evmtypes.Subspace) AppModule {
	return AppModule{
		AppModule:      evm.NewAppModule(k.Keeper, accountKeeper, bankKeeper, legacySubspace),
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		legacySubspace: legacySubspace,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	am.AppModule.RegisterServices(cfg)
}

func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	// not reset chain-id on the begin-block
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState evmtypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.accountKeeper, am.bankKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}
