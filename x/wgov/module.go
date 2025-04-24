package wgov

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	wgovkeeper "github.com/st-chain/me-hub/x/wgov/keeper"
)

// AppModuleBasic implements the basic application module for the wrapped nft module.
type AppModuleBasic struct {
	gov.AppModuleBasic
}

// AppModule implements an application module for the wnft module.
type AppModule struct {
	gov.AppModule
	keeper *wgovkeeper.Keeper
}

func NewAppModule(cdc codec.Codec, keeper *wgovkeeper.Keeper, ak govtypes.AccountKeeper, bk govtypes.BankKeeper, ss govtypes.ParamSubspace) AppModule {
	return AppModule{
		AppModule: gov.NewAppModule(cdc, &keeper.Keeper, ak, bk, ss),
		keeper:    keeper,
	}
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}
