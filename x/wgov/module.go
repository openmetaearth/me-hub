package wgov

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	wgovkeeper "github.com/st-chain/me-hub/x/wgov/keeper"
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

// EndBlock returns the end blocker for the gov module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx context.Context) error {
	c := sdk.UnwrapSDKContext(ctx)
	return EndBlocker(c, am.keeper)
}
