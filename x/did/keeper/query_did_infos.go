package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/st-chain/me-hub/x/did/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DidInfos(goCtx context.Context, req *types.QueryDidInfosRequest) (*types.QueryDidInfosResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	pstore := prefix.NewStore(store, types.DidInfoPrefix)
	infos := []types.DidInfo{}
	pageRes, err := query.Paginate(pstore, req.Pagination, func(key []byte, value []byte) error {
		var info types.DidInfo
		if err := k.cdc.Unmarshal(value, &info); err != nil {
			return err
		}

		infos = append(infos, info)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDidInfosResponse{
		Infos:      infos,
		Pagination: pageRes,
	}, nil
}
