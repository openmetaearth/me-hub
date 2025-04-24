package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/dao/types"
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
