package rollup

import (
	"encoding/json"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/rollup/client/cli"
	"github.com/st-chain/me-hub/x/rollup/keeper"
	"github.com/st-chain/me-hub/x/rollup/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the capability module.
type AppModuleBasic struct {
	cdc codec.BinaryCodec
}

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the capability module's name.
func (AppModuleBasic) Name() string {
	return types.MODULE_NAME
}

func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns the capability module's default genesis state.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis performs genesis state validation for the capability module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.MODULE_NAME, err)
	}
	return genState.Validate()
}

// RegisterRESTRoutes registers the capability module's REST service handlers.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// nolint: errcheck, gosec
	//TODO:后面需加上查询
	//types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

// GetTxCmd returns the capability module's root tx command.
// TODO:目前暂时设置为空，留待后续有需要增增加cli端
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// TODO:目前暂时设置为空，留待后续有需要增增加cli端
// GetQueryCmd returns the capability module's root query command.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	//return cli.GetQueryCmd(types.StoreKey)
	return cli.GetQueryCmd(types.StoreKey)
}

type AppModule struct {
	AppModuleBasic
	Keeper *keeper.Keeper
	bk     types.BankKeeper
}

func NewAppModule(cdc codec.Codec, keeper *keeper.Keeper, bankKeeper types.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		Keeper:         keeper,
		bk:             bankKeeper,
	}
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(*am.Keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(*am.Keeper))
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, *am.Keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, *am.Keeper)

	return cdc.MustMarshalJSON(gs)
}

// BeginBlock contains the logic that is automatically triggered at the beginning of each block
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am AppModule) EndBlock(ctx sdk.Context, endBlk abci.RequestEndBlock) []abci.ValidatorUpdate {
	logger := am.Keeper.Logger(ctx)
	err := am.Keeper.ProcUnstake(ctx)
	if err != nil {
		logger.Error(err.Error())
	}
	if err = am.Keeper.ProcElection(ctx); err != nil {
		logger.Error(err.Error())
	}
	return []abci.ValidatorUpdate{}
}
