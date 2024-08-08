package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/dao/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) GlobalDao(goCtx context.Context, req *types.QueryGlobalDaoRequest) (*types.QueryGlobalDaoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	globalDao := k.GetGlobalDao(ctx)
	if globalDao == nil {
		return &types.QueryGlobalDaoResponse{}, nil
	}

	return &types.QueryGlobalDaoResponse{Address: globalDao.String()}, nil
}
