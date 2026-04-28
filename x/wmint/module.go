package wmint

import (
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintexported "github.com/cosmos/cosmos-sdk/x/mint/exported"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/openmetaearth/me-hub/x/wmint/keeper"
)

// AppModuleBasic implements the basic application module for the wrapped nft module.
type AppModuleBasic struct {
	mint.AppModuleBasic
}

// AppModule implements an application module for the wnft module.
type AppModule struct {
	mint.AppModule
	keeper         keeper.Keeper
	legacySubspace mintexported.Subspace
}

// NewAppModule creates a new wmint AppModule object.
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	ak minttypes.AccountKeeper,
	ic minttypes.InflationCalculationFn,
	ss mintexported.Subspace,
) AppModule {
	if ic == nil {
		ic = minttypes.DefaultInflationCalculationFn
	}
	mintModule := mint.NewAppModule(cdc, keeper.Keeper, ak, ic, ss)

	return AppModule{
		AppModule:      mintModule,
		keeper:         keeper,
		legacySubspace: ss,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	minttypes.RegisterMsgServer(cfg.MsgServer(), mintkeeper.NewMsgServerImpl(am.keeper.Keeper))
	minttypes.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	m := mintkeeper.NewMigrator(am.keeper.Keeper, am.legacySubspace)

	if err := cfg.RegisterMigration(minttypes.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", minttypes.ModuleName, err))
	}
}

func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper, nil)
}
