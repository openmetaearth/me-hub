package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/rollup/types"
)

func (t Keeper) QueryParams(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	paramas := types.Params{
		ElectionPeriod:        t.GetElectionPeriod(sdkCtx),
		SequencerNumber:       t.GetSequencerNumber(sdkCtx),
		BackupSequencerNumber: t.GetBackupNumber(sdkCtx),
		MinStakeAmount:        t.GetMinStakeAmount(sdkCtx),
		//	FirstElectionInterval:  t.GetFirstElectionInterval(sdkCtx),
		AllowApplyElectionTime: t.GetAllowApplyElectionTime(sdkCtx),
		ElectionInterimTime:    t.GetElectionInterimTime(sdkCtx),
		DaFraudChallengeStake:  t.GetDaFraudChallengeStake(sdkCtx),
	}
	return &types.QueryParamsResponse{
		Params: paramas,
	}, nil

}

func (t Keeper) QueryElectionResult(ctx context.Context, req *types.QueryElectionRequest) (*types.QueryElectionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if _, ok := t.mapRollappInfoMng[req.RollappId]; !ok {
		return nil, errorsmod.Wrapf(types.ErrNotFound, fmt.Sprintf("can not found rollapp Info, rollappID = %s", req.RollappId))
	}

	kvStore := sdkCtx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(req.RollappId))
	data := store.Get([]byte(types.KEY_LAST_ELECTION_INFO))

	resp := &types.QueryElectionResponse{
		ElectionTime:   0,
		BlockHeight:    0,
		NodeStatusList: nil,
	}
	if nil == data {
		return resp, nil
	}
	if err := t.cdc.Unmarshal(data, resp); err != nil {
		return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Unmarshal error. err = %s", err.Error()))
	}
	return resp, nil
}

func (t Keeper) GetPreviousElectionResult(ctx context.Context, rollappID string) (*types.QueryElectionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if _, ok := t.mapRollappInfoMng[rollappID]; !ok {
		return nil, errorsmod.Wrapf(types.ErrNotFound, fmt.Sprintf("can not found rollapp Info, rollappID = %s", rollappID))
	}
	kvStore := sdkCtx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(rollappID))
	data := store.Get([]byte(types.KEY_PREVIOUS_ELECTION_INFO))

	resp := &types.QueryElectionResponse{
		ElectionTime:   0,
		BlockHeight:    0,
		NodeStatusList: nil,
	}
	if nil == data {
		return resp, nil
	}
	if err := t.cdc.Unmarshal(data, resp); err != nil {
		return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Unmarshal error. err = %s", err.Error()))
	}
	return resp, nil
}

func (t Keeper) QueryStake(ctx context.Context, req *types.QueryStakeRequest) (*types.QueryStakeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if _, ok := t.mapRollappInfoMng[req.RollappId]; !ok {
		return nil, errorsmod.Wrapf(types.ErrNotFound, fmt.Sprintf("can not found rollapp Info, rollappID = %s", req.RollappId))
	}

	sdkCtx.Logger().Info(fmt.Sprintf("QueryStake,rollappID = %s", req.RollappId))
	kvStore := sdkCtx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppStakeKeyPrefix(req.RollappId))
	data := store.Get([]byte(req.Address))
	resp := &types.MsgStakeInfo{
		StakeAmount:        0,
		ApplyUnStakeAmount: 0,
	}
	if data != nil {
		if err := t.cdc.Unmarshal(data, resp); err != nil {
			return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Unmarshal error. err = %s", err.Error()))
		}

	}

	return &types.QueryStakeResponse{
		StakeInfo: resp,
	}, nil
}

func (t Keeper) queryStakeData(sdkCtx sdk.Context, rollappId, address string) (*types.MsgStakeInfo, error) {
	kvStore := sdkCtx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppStakeKeyPrefix(rollappId))
	data := store.Get([]byte(address))
	resp := &types.MsgStakeInfo{
		StakeAmount:        0,
		ApplyUnStakeAmount: 0,
	}
	if data != nil {
		if err := t.cdc.Unmarshal(data, resp); err != nil {
			return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Unmarshal stake data error. err = %s", err.Error()))
		}

	}

	return resp, nil
}

func (t Keeper) QueryStakeBondNode(ctx context.Context, req *types.QueryStakeBondNodeRequest) (*types.QueryStakeBondNodeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := sdkCtx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, types.GetDelegatorStakeNodePrefix(req.RollappId))
	data := store.Get([]byte(req.StakeAddr))
	return &types.QueryStakeBondNodeResponse{BondNodeAddr: data}, nil

}
