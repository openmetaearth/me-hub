package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/x/dao/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) GlobalDao(goCtx context.Context, req *types.QueryGlobalDaoRequest) (*types.QueryGlobalDaoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	daoAddresses, found := k.GetDaoAddresses(ctx)
	if !found {
		return &types.QueryGlobalDaoResponse{}, types.ErrNotFound
	}

	return &types.QueryGlobalDaoResponse{DaoAddresses: daoAddresses}, nil
}

func (k Keeper) GlobalDaoFeePool(goCtx context.Context, req *types.QueryGlobalDaoFeePoolReq) (*types.QueryGlobalDaoFeePoolResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	account := k.GetGlobalDaoFeePoolAddr(ctx)
	return &types.QueryGlobalDaoFeePoolResp{GlobalDaoFeePool: account.String()}, nil
}

func (k Keeper) FreeGasAccounts(goCtx context.Context, req *types.QueryFreeGasAccountsReq) (*types.QueryFreeGasAccountsResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	pstore := prefix.NewStore(store, types.FreeGasAddressePrefix)

	var accounts []string
	pageRes, err := query.Paginate(pstore, req.Pagination, func(key []byte, value []byte) error {
		accounts = append(accounts, string(value))
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryFreeGasAccountsResp{Addresses: accounts, Pagination: pageRes}, nil
}

func (k Keeper) IsFreeGasAccount(goCtx context.Context, req *types.QueryIsFreeGasAccountReq) (*types.QueryIsFreeGasAccountResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.QueryIsFreeGasAccountResp{IsFree: k.CheckFreeGasAccount(ctx, req.Address)}, nil
}
