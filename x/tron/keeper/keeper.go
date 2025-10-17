package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/tron/types"

	gravitykeeper "github.com/st-chain/me-hub/x/gravity/keeper"
)

type Keeper struct {
	gravitykeeper.Keeper
}

func NewKeeper(keeper gravitykeeper.Keeper) Keeper {
	return Keeper{
		Keeper: keeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func NewModuleHandler(keeper Keeper) *gravitykeeper.ModuleHandler {
	return &gravitykeeper.ModuleHandler{
		QueryServer: gravitykeeper.NewQueryServerImpl(keeper.Keeper),
		MsgServer:   NewMsgServerImpl(keeper),
	}
}
