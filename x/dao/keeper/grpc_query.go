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
	daoAddresses, found := k.GetDaoAddresses(ctx)
	if !found {
		return &types.QueryGlobalDaoResponse{}, types.ErrNotFound
	}

	return &types.QueryGlobalDaoResponse{DaoAddresses: daoAddresses}, nil
}
