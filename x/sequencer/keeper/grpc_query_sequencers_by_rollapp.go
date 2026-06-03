package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SequencersByRollapp(c context.Context, req *types.QueryGetSequencersByRollappRequest) (*types.QueryGetSequencersByRollappResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	if _, ok := k.rollappKeeper.GetRollapp(ctx, req.RollappId); !ok {
		return nil, types.ErrUnknownRollappID
	}

	var sequencers []types.Sequencer
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SequencersByRollappKey(req.RollappId))
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var sequencer types.Sequencer
		if err := k.cdc.Unmarshal(value, &sequencer); err != nil {
			return err
		}
		sequencers = append(sequencers, sequencer)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetSequencersByRollappResponse{
		Sequencers: sequencers,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) SequencersByRollappByStatus(c context.Context, req *types.QueryGetSequencersByRollappByStatusRequest) (*types.QueryGetSequencersByRollappByStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	if _, ok := k.rollappKeeper.GetRollapp(ctx, req.RollappId); !ok {
		return nil, types.ErrUnknownRollappID
	}

	var sequencers []types.Sequencer
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SequencersByRollappByStatusKey(req.RollappId, req.Status))
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var sequencer types.Sequencer
		if err := k.cdc.Unmarshal(value, &sequencer); err != nil {
			return err
		}
		sequencers = append(sequencers, sequencer)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetSequencersByRollappByStatusResponse{
		Sequencers: sequencers,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) UnConfirmSequencerAddressByRollappByStatus(goCtx context.Context, req *types.QueryGetUnConfirmSequencersAddrByRollappRequest) (*types.QueryGetUnConfirmSequencersAddrByRollappResponse, error) {
	/*
		if req == nil {
			return nil, status.Error(codes.InvalidArgument, "invalid request")
		}
		ctx := sdk.UnwrapSDKContext(goCtx)

		if _, ok := k.rollappKeeper.GetRollapp(ctx, req.RollappId); !ok {
			return nil, types.ErrUnknownRollappID
		}
		replaceProposer, err := k.GetReplaceProposer(ctx, req.RollappId)
		if err != nil {
			return nil, fmt.Errorf("UnConfirmSequencerAddressByRollappByStatus: failed to get replace proposer info, "+
				"rollappId = %s ,err = %s", req.RollappId, err.Error())
		}
		if replaceProposer == nil {
			return nil, nil
		}
		if (req.BlockHeight >= replaceProposer.BlockHeight) &&
			(req.BlockHeight < (replaceProposer.BlockHeight + int64(k.replaceSequencerCacheHeight))) {
			val, found := k.GetSequencer(ctx, replaceProposer.NewProposer)
			if !found {
				k.Logger(ctx).Error("UnConfirmSequencerAddressByRollappByStatus: can not found new sequencer address in sequencer store",
					"rollappId", req.RollappId, " blockHeight", req.BlockHeight, "newSequencerAddr", replaceProposer.NewProposer)
				return nil, fmt.Errorf("can not found new sequencer address of replaceProposer. address = %s ,rollapp = %s",
					replaceProposer.NewProposer, req.RollappId)
			}

			return &types.QueryGetUnConfirmSequencersAddrByRollappResponse{
				NewSequencer:         val,
				StartHeight:          replaceProposer.BlockHeight,
				UnconfirmCacheHeight: int64(k.replaceSequencerCacheHeight),
			}, nil
		} else {
			k.Logger(ctx).Info("UnConfirmSequencerAddressByRollappByStatus: can not found sequencer at this height",
				"rollappId", req.RollappId, " blockHeight", req.BlockHeight, "replaceProposerHeight", replaceProposer.BlockHeight,
				"cacheHeight", k.replaceSequencerCacheHeight)
		}
	*/
	return nil, fmt.Errorf("unsupport function")
}
func (k Keeper) ReplaceProposerInfo(goCtx context.Context, req *types.QueryReplaceProposerInfoRequest) (*types.QueryReplaceProposerInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, found := k.rollappKeeper.GetRollapp(ctx, req.RollappId)
	if !found {
		return nil, types.ErrUnknownRollappID
	}
	replaceProposer, err := k.GetReplaceProposer(ctx, req.RollappId)
	if err != nil {
		return nil, fmt.Errorf("ReplaceProposerInfo: failed to get replace proposer info,rollappID = %s, err = %s",
			req.RollappId, err.Error())
	}
	if nil == replaceProposer {
		return nil, nil
	}
	return &types.QueryReplaceProposerInfoResponse{
		ReplaceProposer: *replaceProposer,
	}, nil
}
