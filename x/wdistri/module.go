package wdistri

import (
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/exported"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/openmetaearth/me-hub/x/wdistri/keeper"
	"github.com/openmetaearth/me-hub/x/wdistri/types"
)

// AppModuleBasic implements the basic application module for the wrapped nft module.
type AppModuleBasic struct {
	distribution.AppModuleBasic
}

// AppModule implements an application module for the wnft module.
type AppModule struct {
	distribution.AppModule
	keeper         keeper.Keeper
	legacySubspace exported.Subspace
}

// NewAppModule creates a new wmint AppModule object.
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.StakingKeeper,
	ss exported.Subspace,
) AppModule {
	distributionModule := distribution.NewAppModule(cdc, keeper.Keeper, ak, bk, sk, ss)

	return AppModule{
		AppModule:      distributionModule,
		keeper:         keeper,
		legacySubspace: ss,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	distributiontypes.RegisterMsgServer(cfg.MsgServer(), distributionkeeper.NewMsgServerImpl(am.keeper.Keeper))
	distributiontypes.RegisterQueryServer(cfg.QueryServer(), distributionkeeper.NewQuerier(am.keeper.Keeper))

	m := distributionkeeper.NewMigrator(am.keeper.Keeper, am.legacySubspace)
	if err := cfg.RegisterMigration(distributiontypes.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", distributiontypes.ModuleName, err))
	}

	if err := cfg.RegisterMigration(distributiontypes.ModuleName, 2, m.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 2 to 3: %v", distributiontypes.ModuleName, err))
	}
}

// BeginBlock returns the begin blocker for the distribution module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlock(ctx, req, am.keeper)
	return []abci.ValidatorUpdate{}
}
