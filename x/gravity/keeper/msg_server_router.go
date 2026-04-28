package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

type msgServer struct {
	routerKeeper RouterKeeper
}

// NewMsgServerRouterImpl returns an implementation of the crosschain router MsgServer interface
// for the provided Keeper.
func NewMsgServerRouterImpl(routerKeeper RouterKeeper) types.MsgServer {
	return &msgServer{routerKeeper: routerKeeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) BondedRelayer(ctx context.Context, msg *types.MsgBondedRelayer) (*types.MsgBondedRelayerResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.BondedRelayer(ctx, msg)
	}
}

func (k msgServer) AddDelegate(ctx context.Context, msg *types.MsgAddDelegate) (*types.MsgAddDelegateResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.AddDelegate(ctx, msg)
	}
}

func (k msgServer) UnbondedRelayer(ctx context.Context, msg *types.MsgUnbondedRelayer) (*types.MsgUnbondedRelayerResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.UnbondedRelayer(ctx, msg)
	}
}

func (k msgServer) RelayerSetConfirm(ctx context.Context, msg *types.MsgRelayerSetConfirm) (*types.MsgRelayerSetConfirmResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.RelayerSetConfirm(ctx, msg)
	}
}

func (k msgServer) RelayerSetUpdateClaim(ctx context.Context, msg *types.MsgRelayerSetUpdateClaim) (*types.MsgRelayerSetUpdateClaimResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.RelayerSetUpdateClaim(ctx, msg)
	}
}

func (k msgServer) BridgeTokenClaim(ctx context.Context, msg *types.MsgBridgeTokenClaim) (*types.MsgBridgeTokenClaimResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.BridgeTokenClaim(ctx, msg)
	}
}

func (k msgServer) SendToMeClaim(ctx context.Context, msg *types.MsgSendToMeClaim) (*types.MsgSendToMeClaimResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.SendToMeClaim(ctx, msg)
	}
}

func (k msgServer) SendToExternal(ctx context.Context, msg *types.MsgSendToExternal) (*types.MsgSendToExternalResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.SendToExternal(ctx, msg)
	}
}

func (k msgServer) CancelSendToExternal(ctx context.Context, msg *types.MsgCancelSendToExternal) (*types.MsgCancelSendToExternalResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.CancelSendToExternal(ctx, msg)
	}
}

func (k msgServer) IncreaseBridgeFee(ctx context.Context, msg *types.MsgIncreaseBridgeFee) (*types.MsgIncreaseBridgeFeeResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.IncreaseBridgeFee(ctx, msg)
	}
}

func (k msgServer) SendToExternalClaim(ctx context.Context, msg *types.MsgSendToExternalClaim) (*types.MsgSendToExternalClaimResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.SendToExternalClaim(ctx, msg)
	}
}

func (k msgServer) RequestBatch(ctx context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.RequestBatch(ctx, msg)
	}
}

func (k msgServer) ConfirmBatch(ctx context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.ConfirmBatch(ctx, msg)
	}
}

func (k msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.UpdateParams(ctx, msg)
	}
}

func (k msgServer) ProposalRelayers(ctx context.Context, msg *types.MsgProposalRelayers) (*types.MsgProposalRelayersResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.ProposalRelayers(ctx, msg)
	}
}

func (k msgServer) getMsgServerByChainName(chainName string) (types.MsgServer, error) {
	msgServerRouter := k.routerKeeper.Router()
	if !msgServerRouter.HasRoute(chainName) {
		return nil, errorsmod.Wrap(errortypes.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", chainName))
	}
	return msgServerRouter.GetRoute(chainName).MsgServer, nil
}
