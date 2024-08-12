package keeper

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

// MsgServer is wrapper staking customParamsKeeper message server.
type MsgServer struct {
	stakingtypes.MsgServer
	*Keeper
}

var _ types.MsgServer = MsgServer{}

// NewMsgServerImpl returns an implementation of the staking wrapped MsgServer.
func NewMsgServerImpl(
	keeper *Keeper,
	stakingMsgSrv stakingtypes.MsgServer,
) MsgServer {
	return MsgServer{
		Keeper:    keeper,
		MsgServer: stakingMsgSrv,
	}
}
