package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	commontypes "github.com/openmetaearth/me-hub/x/common/types"
	"github.com/openmetaearth/me-hub/x/delayedack/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Querier{}

type Querier struct {
	Keeper
}

// NewQuerier creates a new Querier struct.
func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

func (q Querier) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: q.GetParams(ctx)}, nil
}

// GetPackets implements types.QueryServer.
func (q Querier) GetPackets(goCtx context.Context, req *types.QueryRollappPacketsRequest) (*types.QueryRollappPacketListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	res := &types.QueryRollappPacketListResponse{}
	var filter types.RollappPacketListFilter

	if req.RollappId == "" {
		// query by status (PENDING by default) and type (if not UNDEFINED)
		filter = types.ByTypeByStatus(req.Type, req.Status)
	} else {
		// query by rollapp id and status (PENDING by default) and type (if not UNDEFINED)
		filter = types.ByRollappIDByTypeByStatus(req.RollappId, req.Type, req.Status)
	}

	packetStore := prefix.NewStore(ctx.KVStore(q.storeKey), filter.Prefixes[0].Start)
	pageRes, err := query.FilteredPaginate(packetStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var packet commontypes.RollappPacket
		q.cdc.MustUnmarshal(value, &packet)
		if !filter.FilterFunc(packet) {
			return false, nil
		}
		if accumulate {
			res.RollappPackets = append(res.RollappPackets, packet)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res.Pagination = pageRes

	return res, nil
}
