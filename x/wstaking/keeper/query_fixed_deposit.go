package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) FixedDepositAll(c context.Context, req *types.QueryAllFixedDepositRequest) (*types.QueryAllFixedDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	fixedDeposits, pageRes, err := k.GetAllFixedDepositWithPage(ctx, req)

	return &types.QueryAllFixedDepositResponse{FixedDeposit: fixedDeposits, Pagination: pageRes}, err
}

func (k Keeper) FixedDeposit(c context.Context, req *types.QueryGetFixedDepositRequest) (*types.QueryGetFixedDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	fixedDeposit, found := k.GetFixedDeposit(ctx, req.Id)
	if !found {
		return nil, types.ErrNoFixedDepositFound.Wrapf("addr:%s", req.Address)
	}

	return &types.QueryGetFixedDepositResponse{FixedDeposit: fixedDeposit}, nil
}

func (k Keeper) FixedDepositByAcct(goCtx context.Context, req *types.QueryFixedDepositByAcctRequest) (*types.QueryFixedDepositByAcctResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.QueryType != types.FixedDepositState_AllState && req.QueryType != types.FixedDepositState_NotExpired && req.QueryType != types.FixedDepositState_Expired {
		return nil, status.Error(codes.InvalidArgument, "invalid query type")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var fixedDeposits []types.FixedDeposit
	tmpList, err := k.GetFixedDepositByAcct(ctx, req.Account)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if req.QueryType == types.FixedDepositState_AllState {
		return &types.QueryFixedDepositByAcctResponse{FixedDeposit: tmpList}, nil
	}
	for _, v := range tmpList {
		switch req.QueryType {
		case types.FixedDepositState_NotExpired:
			if ctx.BlockTime().Before(v.EndTime) {
				fixedDeposits = append(fixedDeposits, v)
			}
		case types.FixedDepositState_Expired:
			if ctx.BlockTime().After(v.EndTime) {
				fixedDeposits = append(fixedDeposits, v)
			}
		}
	}
	return &types.QueryFixedDepositByAcctResponse{FixedDeposit: fixedDeposits}, nil
}

func (k Keeper) FixedDepositByRegion(goCtx context.Context, req *types.QueryFixedDepositByRegionRequest) (*types.QueryFixedDepositByRegionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.QueryType != types.FixedDepositState_AllState && req.QueryType != types.FixedDepositState_NotExpired && req.QueryType != types.FixedDepositState_Expired {
		return nil, status.Error(codes.InvalidArgument, "invalid query type")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	fixedDeposits, pageRes, err := k.queryFixedDepositByRegionRecursively(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	return &types.QueryFixedDepositByRegionResponse{FixedDeposit: fixedDeposits, Pagination: pageRes}, nil
}

func (k Keeper) queryFixedDepositByRegionRecursively(ctx sdk.Context, req *types.QueryFixedDepositByRegionRequest, accumulated []types.FixedDeposit) ([]types.FixedDeposit, *query.PageResponse, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FixedDepositKey))
	fixedDeposits := make([]types.FixedDeposit, 0)

	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var fd types.FixedDeposit
		if err := k.cdc.Unmarshal(value, &fd); err != nil {
			return err
		}

		regionId, err := k.MustGetKycRegionIdByAccount(ctx, fd.Account)
		if err != nil {
			return err
		}

		if regionId == req.RegionId {
			switch req.QueryType {
			case types.FixedDepositState_AllState:
				fixedDeposits = append(fixedDeposits, fd)
			case types.FixedDepositState_NotExpired:
				if ctx.BlockTime().Before(fd.EndTime) {
					fixedDeposits = append(fixedDeposits, fd)
				}
			case types.FixedDepositState_Expired:
				if ctx.BlockTime().After(fd.EndTime) {
					fixedDeposits = append(fixedDeposits, fd)
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	accumulated = append(accumulated, fixedDeposits...)

	if len(accumulated) < int(req.Pagination.Limit) && pageRes.NextKey != nil {
		req.Pagination.Key = pageRes.NextKey
		return k.queryFixedDepositByRegionRecursively(ctx, req, accumulated)
	}

	return accumulated, pageRes, nil
}

func (k Keeper) FixedDepositTotalAmount(goCtx context.Context, req *types.QueryFixedDepositTotalAmountRequest) (*types.QueryFixedDepositTotalAmountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	amount, found := k.GetFixedDepositTotalAmount(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "total fixed Deposit amount not found")
	}
	return &types.QueryFixedDepositTotalAmountResponse{Amount: amount.Amount}, nil
}

func (k Keeper) FixedDepositAmountByMeid(goCtx context.Context, req *types.QueryFixedDepositAmountByMeidRequest) (*types.QueryFixedDepositAmountByMeidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	tmpList, err := k.GetFixedDepositByAcct(ctx, req.Account)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	totalAmount := sdk.NewCoin(params.BaseDenom, sdk.NewInt(0))
	for _, v := range tmpList {
		totalAmount = totalAmount.Add(v.Principal)
	}

	return &types.QueryFixedDepositAmountByMeidResponse{Amount: totalAmount}, nil
}
