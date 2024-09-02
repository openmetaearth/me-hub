package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/wstaking/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) FixedDepositCfg(goCtx context.Context, req *types.QueryFixedDepositCfgRequest) (*types.QueryFixedDepositCfgResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	configs := k.GetAllFixedDepositCfg(ctx, req.RegionId)
	return &types.QueryFixedDepositCfgResponse{FixedDepositCfgs: configs}, nil
}

func (k Keeper) FixedDepositCfgByTerm(goCtx context.Context, req *types.QueryFixedDepositCfgByTermRequest) (*types.QueryFixedDepositCfgByTermResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	config, ok := k.GetFixedDepositCfg(ctx, req.RegionId, req.Term)
	if !ok {
		return nil, status.Error(codes.NotFound, "config not found")
	}
	return &types.QueryFixedDepositCfgByTermResponse{FixedDepositCfg: config}, nil
}
