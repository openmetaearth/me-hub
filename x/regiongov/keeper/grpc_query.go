package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) Proposal(goCtx context.Context, req *types.QueryProposalRequest) (*types.QueryProposalResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	proposal, found := k.GetProposal(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "proposal not found")
	}
	return &types.QueryProposalResponse{Proposal: proposal}, nil
}

func (k Keeper) Proposals(goCtx context.Context, req *types.QueryProposalsRequest) (*types.QueryProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	proposals := filterProposals(k.GetAllProposals(ctx), req.Status, req.Kind)
	page, pageRes := paginate(proposals, req.Pagination)
	return &types.QueryProposalsResponse{Proposals: page, Pagination: pageRes}, nil
}

func (k Keeper) Votes(goCtx context.Context, req *types.QueryVotesRequest) (*types.QueryVotesResponse, error) {
	if req == nil || req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	if _, found := k.GetProposal(ctx, req.ProposalId); !found {
		return nil, status.Error(codes.NotFound, "proposal not found")
	}
	votes, pageRes := paginate(k.GetVotes(ctx, req.ProposalId), req.Pagination)
	return &types.QueryVotesResponse{Votes: votes, Pagination: pageRes}, nil
}

func (k Keeper) Electors(goCtx context.Context, req *types.QueryElectorsRequest) (*types.QueryElectorsResponse, error) {
	if req == nil || req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	if _, found := k.GetProposal(ctx, req.ProposalId); !found {
		return nil, status.Error(codes.NotFound, "proposal not found")
	}
	return &types.QueryElectorsResponse{Electors: k.GetElectors(ctx, req.ProposalId)}, nil
}

func (k Keeper) Equity(goCtx context.Context, req *types.QueryEquityRequest) (*types.QueryEquityResponse, error) {
	if req == nil || req.RegionId == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	equity, found := k.GetEquity(ctx, req.RegionId)
	if !found {
		equity = types.RegionEquity{RegionId: req.RegionId}
	}
	return &types.QueryEquityResponse{Equity: equity}, nil
}

func (k Keeper) RegionFreeze(goCtx context.Context, req *types.QueryRegionFreezeRequest) (*types.QueryRegionFreezeResponse, error) {
	if req == nil || req.RegionId == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	freeze, found := k.GetRegionFreeze(ctx, req.RegionId)
	return &types.QueryRegionFreezeResponse{Frozen: found && freeze.Frozen}, nil
}

func (k Keeper) EquityFreezes(goCtx context.Context, req *types.QueryEquityFreezesRequest) (*types.QueryEquityFreezesResponse, error) {
	if req == nil || req.RegionId == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.QueryEquityFreezesResponse{Freezes: k.GetEquityFreezes(ctx, req.RegionId)}, nil
}

func (k Keeper) PendingAction(goCtx context.Context, req *types.QueryPendingActionRequest) (*types.QueryPendingActionResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	action, found := k.GetPendingAction(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "pending action not found")
	}
	return &types.QueryPendingActionResponse{Action: action}, nil
}

func (k Keeper) PendingActions(goCtx context.Context, req *types.QueryPendingActionsRequest) (*types.QueryPendingActionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	actions := filterActions(k.GetAllPendingActions(ctx), req.Status)
	page, pageRes := paginate(actions, req.Pagination)
	return &types.QueryPendingActionsResponse{Actions: page, Pagination: pageRes}, nil
}
