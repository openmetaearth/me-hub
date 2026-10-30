package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(k Keeper) types.MsgServer {
	return msgServer{Keeper: k}
}

func (k msgServer) SubmitProposal(goCtx context.Context, msg *types.MsgSubmitProposal) (*types.MsgSubmitProposalResponse, error) {
	id, err := k.submitProposal(sdk.UnwrapSDKContext(goCtx), msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgSubmitProposalResponse{Id: id}, nil
}

func (k msgServer) Vote(goCtx context.Context, msg *types.MsgVote) (*types.MsgVoteResponse, error) {
	if err := k.castVote(sdk.UnwrapSDKContext(goCtx), msg); err != nil {
		return nil, err
	}
	return &types.MsgVoteResponse{}, nil
}

func (k msgServer) SubmitEquityUpdate(goCtx context.Context, msg *types.MsgSubmitEquityUpdate) (*types.MsgSubmitEquityUpdateResponse, error) {
	id, err := k.submitEquityUpdate(sdk.UnwrapSDKContext(goCtx), msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgSubmitEquityUpdateResponse{ActionId: id}, nil
}

func (k msgServer) SubmitFreezeEquity(goCtx context.Context, msg *types.MsgSubmitFreezeEquity) (*types.MsgSubmitFreezeEquityResponse, error) {
	id, err := k.submitFreezeEquity(sdk.UnwrapSDKContext(goCtx), msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgSubmitFreezeEquityResponse{ActionId: id}, nil
}

func (k msgServer) SubmitUnfreezeEquity(goCtx context.Context, msg *types.MsgSubmitUnfreezeEquity) (*types.MsgSubmitUnfreezeEquityResponse, error) {
	id, err := k.submitUnfreezeEquity(sdk.UnwrapSDKContext(goCtx), msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgSubmitUnfreezeEquityResponse{ActionId: id}, nil
}

func (k msgServer) ApprovePendingAction(goCtx context.Context, msg *types.MsgApprovePendingAction) (*types.MsgApprovePendingActionResponse, error) {
	if err := k.approvePendingAction(sdk.UnwrapSDKContext(goCtx), msg); err != nil {
		return nil, err
	}
	return &types.MsgApprovePendingActionResponse{}, nil
}

func (k msgServer) RejectPendingAction(goCtx context.Context, msg *types.MsgRejectPendingAction) (*types.MsgRejectPendingActionResponse, error) {
	if err := k.rejectPendingAction(sdk.UnwrapSDKContext(goCtx), msg); err != nil {
		return nil, err
	}
	return &types.MsgRejectPendingActionResponse{}, nil
}

func (k msgServer) FreezeRegion(goCtx context.Context, msg *types.MsgFreezeRegion) (*types.MsgFreezeRegionResponse, error) {
	if err := k.freezeRegion(sdk.UnwrapSDKContext(goCtx), msg); err != nil {
		return nil, err
	}
	return &types.MsgFreezeRegionResponse{}, nil
}

func (k msgServer) UnfreezeRegion(goCtx context.Context, msg *types.MsgUnfreezeRegion) (*types.MsgUnfreezeRegionResponse, error) {
	if err := k.unfreezeRegion(sdk.UnwrapSDKContext(goCtx), msg); err != nil {
		return nil, err
	}
	return &types.MsgUnfreezeRegionResponse{}, nil
}

func (k msgServer) FreezeEquity(goCtx context.Context, msg *types.MsgFreezeEquity) (*types.MsgFreezeEquityResponse, error) {
	if err := k.freezeEquity(sdk.UnwrapSDKContext(goCtx), msg); err != nil {
		return nil, err
	}
	return &types.MsgFreezeEquityResponse{}, nil
}

func (k msgServer) UnfreezeEquity(goCtx context.Context, msg *types.MsgUnfreezeEquity) (*types.MsgUnfreezeEquityResponse, error) {
	if err := k.unfreezeEquity(sdk.UnwrapSDKContext(goCtx), msg); err != nil {
		return nil, err
	}
	return &types.MsgUnfreezeEquityResponse{}, nil
}
