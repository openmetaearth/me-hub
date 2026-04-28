package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/x/did/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Did(goCtx context.Context, req *types.QueryDid) (*types.QueryDidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address: "+err.Error())
	}

	did, found := k.GetDID(ctx, addr)
	if !found {
		return &types.QueryDidResponse{}, types.ErrDidNotFound
	}

	info, found := k.GetDidInfo(ctx, did)
	if !found {
		return &types.QueryDidResponse{}, types.ErrDidNotFound
	}

	return &types.QueryDidResponse{Info: info}, nil
}

func (k Keeper) DidInfo(goCtx context.Context, req *types.QueryDidInfo) (*types.QueryDidInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	info, found := k.GetDidInfo(ctx, req.Did)
	if !found {
		return &types.QueryDidInfoResponse{}, types.ErrDidNotFound
	}

	return &types.QueryDidInfoResponse{Info: info}, nil
}

func (k Keeper) DidInfos(goCtx context.Context, req *types.QueryDidInfos) (*types.QueryDidInfosResponse, error) {
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

func (k Keeper) DidDocument(goCtx context.Context, req *types.QueryDidDocument) (*types.QueryDidDocumentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	info, found := k.GetDidInfo(ctx, req.Did)
	if !found {
		return &types.QueryDidDocumentResponse{}, types.ErrDidNotFound
	}

	vcs := k.GetCredentialsByDid(ctx, req.Did)

	return &types.QueryDidDocumentResponse{Doc: types.DidDocument{Info: info, Vcs: vcs}}, nil
}

func (k Keeper) Service(goCtx context.Context, req *types.QueryService) (*types.QueryServiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	svc, found := k.GetService(ctx, req.Sid)
	if !found {
		return &types.QueryServiceResponse{}, types.ErrServiceNotFound
	}

	return &types.QueryServiceResponse{Service: svc}, nil
}

func (k Keeper) Services(goCtx context.Context, req *types.QueryServices) (*types.QueryServicesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var svcs []types.Service

	store := ctx.KVStore(k.storeKey)
	filterStore := prefix.NewStore(store, types.ServicePrefix)

	pageRes, err := query.Paginate(filterStore, req.Pagination, func(key []byte, value []byte) error {
		var svc types.Service
		if err := k.cdc.Unmarshal(value, &svc); err != nil {
			return err
		}

		svcs = append(svcs, svc)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryServicesResponse{Services: svcs, Pagination: pageRes}, nil
}

func (k Keeper) Credential(goCtx context.Context, req *types.QueryCredential) (*types.QueryCredentialResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	vc, found := k.GetCredential(ctx, req.Did, req.Sid)
	if !found {
		return &types.QueryCredentialResponse{}, types.ErrCredentialNotFound
	}

	return &types.QueryCredentialResponse{Credential: vc}, nil
}

func (k Keeper) Credentials(goCtx context.Context, req *types.QueryCredentials) (*types.QueryCredentialsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var vcs []types.Credential

	store := ctx.KVStore(k.storeKey)
	filterStore := prefix.NewStore(store, types.GetFilterPrefixBySidAndFilter(req.Sid, req.Filter))

	pageRes, err := query.Paginate(filterStore, req.Pagination, func(key []byte, value []byte) error {
		var vc types.Credential
		if err := k.cdc.Unmarshal(value, &vc); err != nil {
			return err
		}

		vcs = append(vcs, vc)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCredentialsResponse{Credentials: vcs, Pagination: pageRes}, nil
}
