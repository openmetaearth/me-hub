package keeper

import (
	"github.com/st-chain/me-hub/x/rollup/types"
)

type rollupServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &rollupServer{Keeper: keeper}
}

type rollupQueryServer struct {
	Keeper
}

func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &rollupQueryServer{Keeper: keeper}
}
