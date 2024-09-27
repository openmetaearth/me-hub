package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/st-chain/me-hub/app/params"
	types2 "github.com/st-chain/me-hub/x/kyc/types"
	"github.com/st-chain/me-hub/x/wstaking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) FixedDepositAll(c context.Context, req *types.QueryAllFixedDepositRequest) (*types.QueryAllFixedDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	fixedDeposits := k.GetAllFixedDeposit(ctx)

	return &types.QueryAllFixedDepositResponse{FixedDeposit: fixedDeposits}, nil
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
	tmpList := k.GetFixedDepositByAcct(ctx, req.Account)
	if req.QueryType == types.FixedDepositState_AllState {
		return &types.QueryFixedDepositByAcctResponse{FixedDeposit: tmpList}, nil
	}
	for _, v := range tmpList {
		switch req.QueryType {
		//case types.AllState:
		//	fixedDeposits = append(fixedDeposits, v)
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
	var fixedDeposits []types.FixedDeposit
	meidList, err := k.KycKeeper.DIDs(ctx, &types2.QueryDIDs{
		RegionId: req.Regionid,
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      9999,
			CountTotal: false,
			Reverse:    false,
		},
	})
	if err != nil {
		return nil, err
	}
	for _, did := range meidList.Infos {
		tmpList := k.GetFixedDepositByAcct(ctx, did.Address)
		if req.QueryType == types.FixedDepositState_AllState {
			fixedDeposits = append(fixedDeposits, tmpList...)
			continue
		} else {
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
		}

	}

	return &types.QueryFixedDepositByRegionResponse{FixedDeposit: fixedDeposits}, nil
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
	tmpList := k.GetFixedDepositByAcct(ctx, req.Account)

	totalAmount := sdk.NewCoin(params.BaseDenom, sdk.NewInt(0))
	for _, v := range tmpList {
		totalAmount = totalAmount.Add(v.Principal)
	}

	return &types.QueryFixedDepositAmountByMeidResponse{Amount: totalAmount}, nil
}
