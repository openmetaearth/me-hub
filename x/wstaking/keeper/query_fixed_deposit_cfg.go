package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) FixedDepositCfg(goCtx context.Context, req *types.QueryFixedDepositCfgRequest) (*types.QueryFixedDepositCfgResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(req.RegionIds) == 0 {
		regions := k.GetAllRegion(ctx)
		for _, region := range regions {
			req.RegionIds = append(req.RegionIds, region.RegionId)
		}
	}
	configs := make([]types.RegionAllFixedDepositCfg, 0, len(req.RegionIds))

	for _, regionId := range req.RegionIds {
		regionConfigs := k.GetAllFixedDepositCfg(ctx, regionId)
		regionFixedDepositCfgs := make([]types.RegionFixedDepositCfg, 0, len(regionConfigs))
		for _, config := range regionConfigs {
			regionFixedDepositCfgs = append(regionFixedDepositCfgs, types.RegionFixedDepositCfg{
				Term:   config.Term,
				Rate:   config.Rate,
				Status: config.Status,
			})
		}
		configs = append(configs, types.RegionAllFixedDepositCfg{
			RegionId:              regionId,
			RegionFixedDepositCfg: regionFixedDepositCfgs,
		})
	}

	return &types.QueryFixedDepositCfgResponse{RegionFixedDepositCfgs: configs}, nil
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
