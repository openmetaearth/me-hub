package keeper

import (
	"github.com/dymensionxyz/dymension/v3/x/rollup/types"
)

type rollupServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &rollupServer{Keeper: keeper}
}
