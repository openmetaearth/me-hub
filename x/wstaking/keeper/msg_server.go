package keeper

import (
	"context"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// MsgServer is wrapper staking customParamsKeeper message server.
type MsgServer struct {
	stakingtypes.MsgServer
}

// NewMsgServerImpl returns an implementation of the staking wrapped MsgServer.
func NewMsgServerImpl(
	stakingMsgSrv stakingtypes.MsgServer,
) stakingtypes.MsgServer {
	return MsgServer{
		MsgServer: stakingMsgSrv,
	}
}

// CreateValidator defines wrapped method for creating a new validator.
func (s MsgServer) CreateValidator(
	goCtx context.Context, msg *stakingtypes.MsgCreateValidator,
) (*stakingtypes.MsgCreateValidatorResponse, error) {
	return s.MsgServer.CreateValidator(goCtx, msg)
}
