package keeper

import (
	"context"
	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

var _ types.QueryServer = RouterKeeper{}

func (k RouterKeeper) getQueryServerByChainName(chainName string) (types.QueryServer, error) {
	if !k.router.HasRoute(chainName) {
		return nil, status.Error(codes.InvalidArgument, "chain name not found:"+chainName)
	}
	return k.router.GetRoute(chainName).QueryServer, nil
}

func (k RouterKeeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.Params(c, req)
	}
}

func (k RouterKeeper) ProposalRelayers(c context.Context, req *types.QueryProposalRelayersRequest) (*types.QueryProposalRelayersResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.ProposalRelayers(c, req)
	}
}

func (k RouterKeeper) Relayer(c context.Context, req *types.QueryRelayerRequest) (*types.QueryRelayerResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.Relayer(c, req)
	}
}

func (k RouterKeeper) Relayers(c context.Context, req *types.QueryRelayersRequest) (*types.QueryRelayersResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.Relayers(c, req)
	}
}

func (k RouterKeeper) CurrentRelayerSet(c context.Context, req *types.QueryCurrentRelayerSetRequest) (*types.QueryCurrentRelayerSetResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.CurrentRelayerSet(c, req)
	}
}

func (k RouterKeeper) RelayerSetRequest(c context.Context, req *types.QueryRelayerSetRequestRequest) (*types.QueryRelayerSetRequestResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.RelayerSetRequest(c, req)
	}
}

func (k RouterKeeper) RelayerSetConfirm(c context.Context, req *types.QueryRelayerSetConfirmRequest) (*types.QueryRelayerSetConfirmResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.RelayerSetConfirm(c, req)
	}
}

func (k RouterKeeper) RelayerSetConfirmsByNonce(c context.Context, req *types.QueryRelayerSetConfirmsByNonceRequest) (*types.QueryRelayerSetConfirmsByNonceResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.RelayerSetConfirmsByNonce(c, req)
	}
}

func (k RouterKeeper) LastRelayerSetRequests(c context.Context, req *types.QueryLastRelayerSetRequestsRequest) (*types.QueryLastRelayerSetRequestsResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastRelayerSetRequests(c, req)
	}
}

func (k RouterKeeper) LastPendingRelayerSetRequestByAddr(c context.Context, req *types.QueryLastPendingRelayerSetRequestByAddrRequest) (*types.QueryLastPendingRelayerSetRequestByAddrResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastPendingRelayerSetRequestByAddr(c, req)
	}
}

func (k RouterKeeper) LastPendingBatchRequestByAddr(c context.Context, req *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastPendingBatchRequestByAddr(c, req)
	}
}

func (k RouterKeeper) OutgoingTxBatches(c context.Context, req *types.QueryOutgoingTxBatchesRequest) (*types.QueryOutgoingTxBatchesResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.OutgoingTxBatches(c, req)
	}
}

func (k RouterKeeper) BatchRequestByNonce(c context.Context, req *types.QueryBatchRequestByNonceRequest) (*types.QueryBatchRequestByNonceResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BatchRequestByNonce(c, req)
	}
}

func (k RouterKeeper) BatchConfirm(c context.Context, req *types.QueryBatchConfirmRequest) (*types.QueryBatchConfirmResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BatchConfirm(c, req)
	}
}

func (k RouterKeeper) BatchConfirms(c context.Context, req *types.QueryBatchConfirmsRequest) (*types.QueryBatchConfirmsResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BatchConfirms(c, req)
	}
}

func (k RouterKeeper) LastEventNonceByAddr(c context.Context, req *types.QueryLastEventNonceByAddrRequest) (*types.QueryLastEventNonceByAddrResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastEventNonceByAddr(c, req)
	}
}

func (k RouterKeeper) PendingOutgoingTxByAddr(c context.Context, req *types.QueryPendingOutgoingTxByAddrRequest) (*types.QueryPendingOutgoingTxByAddrResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.PendingOutgoingTxByAddr(c, req)
	}
}

func (k RouterKeeper) UnbatchedTxs(c context.Context, req *types.QueryUnbatchedTxsRequest) (*types.QueryUnbatchedTxsResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.UnbatchedTxs(c, req)
	}
}

func (k RouterKeeper) LastObservedBlockHeight(c context.Context, req *types.QueryLastObservedBlockHeightRequest) (*types.QueryLastObservedBlockHeightResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastObservedBlockHeight(c, req)
	}
}

func (k RouterKeeper) LastEventBlockHeightByAddr(c context.Context, req *types.QueryLastEventBlockHeightByAddrRequest) (*types.QueryLastEventBlockHeightByAddrResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastEventBlockHeightByAddr(c, req)
	}
}

func (k RouterKeeper) ProjectedBatchTimeoutHeight(c context.Context, req *types.QueryProjectedBatchTimeoutHeightRequest) (*types.QueryProjectedBatchTimeoutHeightResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.ProjectedBatchTimeoutHeight(c, req)
	}
}

func (k RouterKeeper) BridgeTokens(c context.Context, req *types.QueryBridgeTokensRequest) (*types.QueryBridgeTokensResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeTokens(c, req)
	}
}

func (k RouterKeeper) BridgeToken(c context.Context, req *types.QueryBridgeTokenRequest) (*types.QueryBridgeTokenResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeToken(c, req)
	}
}

func (k RouterKeeper) BatchFees(c context.Context, req *types.QueryBatchFeeRequest) (*types.QueryBatchFeeResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BatchFees(c, req)
	}
}

func (k RouterKeeper) ClaimsByEventNonce(c context.Context, req *types.QueryClaimsByEventNonceRequest) (*types.QueryClaimsByEventNonceResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.ClaimsByEventNonce(c, req)
	}
}

func (k RouterKeeper) BridgeChainList(c context.Context, req *types.QueryBridgeChainListRequest) (*types.QueryBridgeChainListResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(bsctypes.ModuleName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeChainList(c, req)
	}
}

func (k RouterKeeper) LastObservedRelayer(c context.Context, req *types.QueryLastObservedRelayer) (*types.QueryLastObservedRelayerResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(bsctypes.ModuleName); err != nil {
		return nil, err
	} else {
		return queryServer.LastObservedRelayer(c, req)
	}
}
