package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/utils/gerrc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openmetaearth/me-hub/x/rollapp/types"
)

func (k Keeper) LatestHeight(c context.Context, req *types.QueryGetLatestHeightRequest) (*types.QueryGetLatestHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	var val types.StateInfoIndex
	var found bool
	if req.Finalized {
		val, found = k.GetLatestFinalizedStateIndex(
			ctx,
			req.RollappId,
		)
	} else {
		val, found = k.GetLatestStateInfoIndex(
			ctx,
			req.RollappId,
		)
	}
	if !found {
		return nil, errorsmod.Wrapf(gerrc.ErrNotFound, "latest index: finalized: %t", req.Finalized)
	}

	state := k.MustGetStateInfo(ctx, req.RollappId, val.Index)

	return &types.QueryGetLatestHeightResponse{
		Height: state.GetLatestHeight(),
	}, nil
}
