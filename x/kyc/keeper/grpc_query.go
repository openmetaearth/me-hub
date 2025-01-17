package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	"github.com/st-chain/me-hub/x/kyc/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Querier struct {
	*Keeper
}

var _ types.QueryServer = Keeper{}

func (k Keeper) Protocol(goCtx context.Context, req *types.QueryProtocol) (*types.QueryProtocolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	proto, found := k.GetProtocol(ctx)
	if !found {
		return nil, status.Error(codes.Internal, "proto not found")
	}

	return &types.QueryProtocolResponse{Protocol: proto}, nil
}

func (k Keeper) DID(goCtx context.Context, req *types.QueryDID) (*types.QueryDIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr := sdk.MustAccAddressFromBech32(req.Address)
	did, found := k.GetDID(ctx, addr)
	if !found {
		return nil, status.Error(codes.Internal, "did not found")
	}

	info, found := k.GetDidInfo(ctx, did)
	if !found {
		return nil, status.Error(codes.Internal, "did not found")
	}

	return &types.QueryDIDResponse{Info: info}, nil
}

func (k Keeper) DIDs(goCtx context.Context, req *types.QueryDIDs) (*types.QueryDIDsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	KYCs, pageRes, err := k.GetKYCsByRegion(ctx, req.RegionId, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var infos []didtypes.DidInfo
	for _, kyc := range KYCs {
		info, found := k.GetDidInfo(ctx, kyc.Did)
		if !found {
			return nil, status.Error(codes.Internal, fmt.Sprintf("kyc exist, but did %s is not found", kyc.Did))
		}

		infos = append(infos, info)
	}

	return &types.QueryDIDsResponse{Infos: infos, Pagination: pageRes}, nil
}

func (k Keeper) KYC(goCtx context.Context, req *types.QueryKYC) (*types.QueryKYCResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasDidInfo(ctx, req.Did) {
		return nil, status.Error(codes.Internal, "DID not found")
	}
	kyc, found := k.GetKYC(ctx, req.Did)
	if !found {
		return nil, status.Error(codes.Internal, "KYC not found")
	}

	return &types.QueryKYCResponse{Kyc: kyc}, nil
}

func (k Keeper) KYCs(goCtx context.Context, req *types.QueryKYCs) (*types.QueryKYCsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	KYCs, pageRes, err := k.GetKYCsByRegion(ctx, req.RegionId, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryKYCsResponse{KYCs: KYCs, Pagination: pageRes}, nil
}

func (k Keeper) SBT(goCtx context.Context, req *types.QuerySBT) (*types.QuerySBTResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	sbt, found := k.GetSBT(ctx, req.Did)
	if !found {
		return nil, status.Error(codes.Internal, "SBT not found")
	}

	return &types.QuerySBTResponse{Sbt: sbt}, nil
}
