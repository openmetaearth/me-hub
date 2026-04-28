package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = QueryServer{}

type QueryServer struct {
	Keeper
}

func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &QueryServer{Keeper: keeper}
}

func (k QueryServer) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := k.GetParams(sdk.UnwrapSDKContext(c))
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k QueryServer) ProposalRelayers(c context.Context, _ *types.QueryProposalRelayersRequest) (*types.QueryProposalRelayersResponse, error) {
	relays, _ := k.GetProposalRelayer(sdk.UnwrapSDKContext(c))
	return &types.QueryProposalRelayersResponse{ProposalRelayer: relays}, nil
}

func (k QueryServer) Relayer(c context.Context, req *types.QueryRelayerRequest) (*types.QueryRelayerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.RelayerAddress == "" && req.ExternalAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "either relayer address or external address must be provided.")
	}
	ctx := sdk.UnwrapSDKContext(c)
	relayer := types.Relayer{}
	found := false
	if req.RelayerAddress != "" {
		relayerAddress, err := sdk.AccAddressFromBech32(req.RelayerAddress)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "relayer address")
		}
		relayer, found = k.GetRelayer(ctx, relayerAddress)
		if !found {
			return nil, types.ErrNotFoundRelayer
		}
	}
	if req.ExternalAddress != "" {
		relayerAddress, found := k.GetRelayerByExternalAddress(ctx, req.ExternalAddress)
		if !found {
			return nil, types.ErrNotFoundRelayer
		}
		relayer, found = k.GetRelayer(ctx, relayerAddress)
		if !found {
			return nil, types.ErrNotFoundRelayer
		}
	}
	return &types.QueryRelayerResponse{Relayer: &relayer}, nil
}

func (k QueryServer) Relayers(c context.Context, _ *types.QueryRelayersRequest) (*types.QueryRelayersResponse, error) {
	relays := k.GetAllRelayers(sdk.UnwrapSDKContext(c), false)
	return &types.QueryRelayersResponse{Relayers: relays}, nil
}

func (k QueryServer) CurrentRelayerSet(c context.Context, _ *types.QueryCurrentRelayerSetRequest) (*types.QueryCurrentRelayerSetResponse, error) {
	return &types.QueryCurrentRelayerSetResponse{RelayerSet: k.GetCurrentRelayerSet(sdk.UnwrapSDKContext(c))}, nil
}

func (k QueryServer) RelayerSetRequest(c context.Context, req *types.QueryRelayerSetRequestRequest) (*types.QueryRelayerSetRequestResponse, error) {
	return &types.QueryRelayerSetRequestResponse{RelayerSet: k.GetRelayerSet(sdk.UnwrapSDKContext(c), req.Nonce)}, nil
}

func (k QueryServer) RelayerSetConfirm(c context.Context, req *types.QueryRelayerSetConfirmRequest) (*types.QueryRelayerSetConfirmResponse, error) {
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryRelayerSetConfirmResponse{Confirm: k.GetRelayerSetConfirm(ctx, req.Nonce, sdk.MustAccAddressFromBech32(req.RelayerAddress))}, nil
}

func (k QueryServer) RelayerSetConfirmsByNonce(c context.Context, req *types.QueryRelayerSetConfirmsByNonceRequest) (*types.QueryRelayerSetConfirmsByNonceResponse, error) {
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	var confirms []*types.MsgRelayerSetConfirm
	k.IterateRelayerSetConfirmByNonce(sdk.UnwrapSDKContext(c), req.Nonce, func(confirm *types.MsgRelayerSetConfirm) bool {
		confirms = append(confirms, confirm)
		return false
	})
	return &types.QueryRelayerSetConfirmsByNonceResponse{Confirms: confirms}, nil
}

func (k QueryServer) LastRelayerSetRequests(c context.Context, req *types.QueryLastRelayerSetRequestsRequest) (*types.QueryLastRelayerSetRequestsResponse, error) {
	var relayerSets []*types.RelayerSet
	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RelayerSetRequestKey)
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var relayerSet types.RelayerSet
		if err := k.cdc.Unmarshal(value, &relayerSet); err != nil {
			return status.Errorf(codes.Internal, "failed to unmarshal relayerSet: %v", err)
		}
		relayerSets = append(relayerSets, &relayerSet)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryLastRelayerSetRequestsResponse{RelayerSets: relayerSets, Pagination: pageRes}, nil
}

func (k QueryServer) LastPendingRelayerSetRequestByAddr(c context.Context, req *types.QueryLastPendingRelayerSetRequestByAddrRequest) (*types.QueryLastPendingRelayerSetRequestByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	relayer, ok := k.GetRelayer(ctx, sdk.MustAccAddressFromBech32(req.RelayerAddress))
	if !ok {
		return nil, types.ErrNotFoundRelayer
	}
	var pendingRelaySetReq []*types.RelayerSet
	k.IterateRelayerSets(ctx, false, func(relaySet *types.RelayerSet) bool {
		if relayer.StartHeight > int64(relaySet.Height) {
			return false
		}
		// found is true if the operatorAddr has signed the relayer set we are currently looking at
		// if this relayer set has NOT been signed by relayerAddr, store it in pendingRelayerSetReq and exit the loop
		if found := k.GetRelayerSetConfirm(ctx, relaySet.Nonce, relayer.GetRelayer()) != nil; !found {
			pendingRelaySetReq = append(pendingRelaySetReq, relaySet)
		}
		// if we have more than 100 unconfirmed requests in
		// our array we should exit, pagination
		return len(pendingRelaySetReq) == 100
	})
	return &types.QueryLastPendingRelayerSetRequestByAddrResponse{RelayerSets: pendingRelaySetReq}, nil
}

func (k QueryServer) LastPendingBatchRequestByAddr(c context.Context, req *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	RelayerAddress := sdk.MustAccAddressFromBech32(req.RelayerAddress)
	relayer, ok := k.GetRelayer(ctx, RelayerAddress)
	if !ok {
		return nil, types.ErrNotFoundRelayer
	}
	var pendingBatchReq *types.OutgoingTxBatch
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		// filter startHeight before confirm
		if relayer.StartHeight > int64(batch.Block) {
			return false
		}
		foundConfirm := k.GetBatchConfirm(ctx, batch.TokenContract, batch.BatchNonce, RelayerAddress) != nil
		if !foundConfirm {
			pendingBatchReq = batch
			return true
		}
		return false
	})
	return &types.QueryLastPendingBatchRequestByAddrResponse{Batch: pendingBatchReq}, nil
}

func (k QueryServer) LastEventNonceByAddr(c context.Context, req *types.QueryLastEventNonceByAddrRequest) (*types.QueryLastEventNonceByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	lastEventNonce := k.GetLastEventNonceByRelayer(ctx, sdk.MustAccAddressFromBech32(req.RelayerAddress))
	return &types.QueryLastEventNonceByAddrResponse{EventNonce: lastEventNonce}, nil
}

func (k QueryServer) LastEventBlockHeightByAddr(c context.Context, req *types.QueryLastEventBlockHeightByAddrRequest) (*types.QueryLastEventBlockHeightByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	lastEventBlockHeight := k.GetLastEventBlockHeightByRelayer(ctx, sdk.MustAccAddressFromBech32(req.RelayerAddress))
	return &types.QueryLastEventBlockHeightByAddrResponse{BlockHeight: lastEventBlockHeight}, nil
}

func (k QueryServer) LastObservedBlockHeight(c context.Context, _ *types.QueryLastObservedBlockHeightRequest) (*types.QueryLastObservedBlockHeightResponse, error) {
	blockHeight := k.GetLastObservedBlockHeight(sdk.UnwrapSDKContext(c))
	nonce := k.GetLastObservedEventNonce(sdk.UnwrapSDKContext(c))
	return &types.QueryLastObservedBlockHeightResponse{
		ExternalBlockHeight:    blockHeight.ExternalBlockHeight,
		BlockHeight:            blockHeight.BlockHeight,
		LastObservedEventNonce: nonce,
	}, nil
}

func (k QueryServer) OutgoingTxBatches(c context.Context, req *types.QueryOutgoingTxBatchesRequest) (*types.QueryOutgoingTxBatchesResponse, error) {
	var batches []*types.OutgoingTxBatch
	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OutgoingTxBatchKey)
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var batch types.OutgoingTxBatch
		if err := k.cdc.Unmarshal(value, &batch); err != nil {
			return status.Errorf(codes.Internal, "failed to unmarshal batch: %v", err)
		}
		batches = append(batches, &batch)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryOutgoingTxBatchesResponse{Batches: batches, Pagination: pageRes}, nil
}

func (k QueryServer) BatchRequestByNonce(c context.Context, req *types.QueryBatchRequestByNonceRequest) (*types.QueryBatchRequestByNonceResponse, error) {
	if err := types.ValidateExternalAddr(req.ChainName, req.GetTokenContract()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token contract address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	foundBatch := k.GetOutgoingTxBatch(sdk.UnwrapSDKContext(c), req.TokenContract, req.Nonce)
	if foundBatch == nil {
		return nil, status.Error(codes.NotFound, "tx batch")
	}
	return &types.QueryBatchRequestByNonceResponse{Batch: foundBatch}, nil
}

func (k QueryServer) BatchConfirm(c context.Context, req *types.QueryBatchConfirmRequest) (*types.QueryBatchConfirmResponse, error) {
	if err := types.ValidateExternalAddr(req.ChainName, req.GetTokenContract()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token contract address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	ctx := sdk.UnwrapSDKContext(c)
	RelayerAddress := sdk.MustAccAddressFromBech32(req.RelayerAddress)
	_, ok := k.GetRelayer(ctx, RelayerAddress)
	if !ok {
		return nil, types.ErrNotFoundRelayer
	}
	confirm := k.GetBatchConfirm(ctx, req.TokenContract, req.Nonce, RelayerAddress)
	return &types.QueryBatchConfirmResponse{Confirm: confirm}, nil
}

func (k QueryServer) BatchConfirms(c context.Context, req *types.QueryBatchConfirmsRequest) (*types.QueryBatchConfirmsResponse, error) {
	if err := types.ValidateExternalAddr(req.ChainName, req.GetTokenContract()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token contract address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	var confirms []*types.MsgConfirmBatch
	k.IterateBatchConfirmByNonceAndTokenContract(sdk.UnwrapSDKContext(c), req.Nonce, req.TokenContract, func(confirm *types.MsgConfirmBatch) bool {
		confirms = append(confirms, confirm)
		return false
	})
	return &types.QueryBatchConfirmsResponse{Confirms: confirms}, nil
}

func (k QueryServer) PendingOutgoingTxByAddr(c context.Context, req *types.QueryPendingOutgoingTxByAddrRequest) (*types.QueryPendingOutgoingTxByAddrResponse, error) {
	if _, err := sdk.AccAddressFromBech32(req.GetSenderAddress()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "sender address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var batches []*types.OutgoingTxBatch
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		batches = append(batches, batch)
		return false
	})
	res := &types.QueryPendingOutgoingTxByAddrResponse{
		TransfersInBatches: make([]*types.OutgoingTransferTx, 0),
		UnbatchedTransfers: make([]*types.OutgoingTransferTx, 0),
	}
	for _, batch := range batches {
		for _, tx := range batch.Transactions {
			if tx.Sender == req.SenderAddress {
				res.TransfersInBatches = append(res.TransfersInBatches, tx)
			}
		}
	}
	k.IterateUnbatchedTransactions(ctx, "", func(tx *types.OutgoingTransferTx) bool {
		if tx.Sender == req.SenderAddress {
			res.UnbatchedTransfers = append(res.UnbatchedTransfers, tx)
		}
		return false
	})
	return res, nil
}

func (k QueryServer) UnbatchedTxs(c context.Context, req *types.QueryUnbatchedTxsRequest) (*types.QueryUnbatchedTxsResponse, error) {
	var unbatchedTxs []*types.OutgoingTransferTx
	ctx := sdk.UnwrapSDKContext(c)
	prefixKey := types.GetOutgoingTxPoolContractPrefix(req.GetTokenContract())
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKey)
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var tx types.OutgoingTransferTx
		k.cdc.MustUnmarshal(value, &tx)
		unbatchedTxs = append(unbatchedTxs, &tx)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryUnbatchedTxsResponse{
		Txs:        unbatchedTxs,
		Pagination: pageRes,
	}, nil
}

func (k QueryServer) ProjectedBatchTimeoutHeight(c context.Context, _ *types.QueryProjectedBatchTimeoutHeightRequest) (*types.QueryProjectedBatchTimeoutHeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	projectedCurrentExternalHeight, batchTimeout := k.GetBatchTimeoutHeight(ctx)
	return &types.QueryProjectedBatchTimeoutHeightResponse{
		TimeoutHeight:                  batchTimeout,
		ProjectedCurrentExternalHeight: projectedCurrentExternalHeight}, nil
}

func (k QueryServer) BridgeTokens(c context.Context, req *types.QueryBridgeTokensRequest) (*types.QueryBridgeTokensResponse, error) {
	bridgeTokens := make([]*types.BridgeToken, 0)
	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgeTokenByDenomKey)
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var bridgeToken types.BridgeToken
		if err := k.cdc.Unmarshal(value, &bridgeToken); err != nil {
			return status.Errorf(codes.Internal, "failed to unmarshal bridgeToken: %v", err)
		}
		bridgeTokens = append(bridgeTokens, &bridgeToken)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryBridgeTokensResponse{BridgeTokens: bridgeTokens, Pagination: pageRes}, nil
}

func (k QueryServer) BridgeToken(c context.Context, req *types.QueryBridgeTokenRequest) (*types.QueryBridgeTokenResponse, error) {
	if len(req.GetDenom()) == 0 && len(req.GetContractAddress()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "bridge coin by denom request must contain a denom")
	}
	ctx := sdk.UnwrapSDKContext(c)

	var bridgeToken *types.BridgeToken
	var err error
	if len(req.GetContractAddress()) > 0 {
		bridgeToken, err = k.GetBridgeTokenByContract(ctx, req.ContractAddress)
		if err != nil {
			return nil, status.Error(codes.NotFound, "contract")
		}
	}
	if len(req.GetDenom()) > 0 {
		if bridgeToken != nil && bridgeToken.Denom != req.Denom {
			return nil, status.Error(codes.NotFound, "denom and contract do not match")
		}
		bridgeToken, err = k.GetBridgeTokenByDenom(ctx, req.Denom)
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrInvalid, "get bridge token: %v", err)
		}
	}
	if bridgeToken == nil {
		return nil, status.Error(codes.NotFound, "denom")
	}
	supply := k.bankKeeper.GetSupply(ctx, bridgeToken.Denom)
	return &types.QueryBridgeTokenResponse{BridgeToken: bridgeToken, TotalSupply: supply}, nil
}

func (k QueryServer) BridgeChainList(_ context.Context, _ *types.QueryBridgeChainListRequest) (*types.QueryBridgeChainListResponse, error) {
	return &types.QueryBridgeChainListResponse{ChainNames: types.GetSupportChains()}, nil
}

// BatchFees queries the batch fees from unbatched pool
func (k QueryServer) BatchFees(c context.Context, req *types.QueryBatchFeeRequest) (*types.QueryBatchFeeResponse, error) {
	if req.GetMinBatchFees() == nil {
		req.MinBatchFees = make([]types.MinBatchFee, 0)
	}
	for _, fee := range req.MinBatchFees {
		if fee.BaseFee.IsNil() || fee.BaseFee.IsNegative() {
			return nil, status.Error(codes.InvalidArgument, "base fee")
		}

		if err := types.ValidateExternalAddr(req.ChainName, fee.TokenContract); err != nil {
			return nil, status.Error(codes.InvalidArgument, "token contract")
		}
	}
	allBatchFees := k.GetAllBatchFees(sdk.UnwrapSDKContext(c), types.MaxResults, req.MinBatchFees)
	return &types.QueryBatchFeeResponse{BatchFees: allBatchFees}, nil
}

func (k QueryServer) ClaimsByEventNonce(c context.Context, req *types.QueryClaimsByEventNonceRequest) (*types.QueryClaimsByEventNonceResponse, error) {
	attestations := []types.Attestation{}
	k.IterateAttestationsByNonce(sdk.UnwrapSDKContext(c), req.EventNonce, func(attestation *types.Attestation) bool {
		attestations = append(attestations, *attestation)
		return false
	})
	return &types.QueryClaimsByEventNonceResponse{Claims: attestations}, nil
}

func (k QueryServer) LastObservedRelayer(c context.Context, req *types.QueryLastObservedRelayer) (*types.QueryLastObservedRelayerResponse, error) {
	set := k.GetLastObservedRelayerSet(sdk.UnwrapSDKContext(c))
	return &types.QueryLastObservedRelayerResponse{RelayerSet: set}, nil
}
