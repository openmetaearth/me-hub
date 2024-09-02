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
	if req.RollappId != t.rollAppID {
		return nil, errorsmod.Wrapf(types.ErrRollappIDMismatch, fmt.Sprintf("rollupServer's rollappID = %s", t.rollAppID))
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	paramas := types.Params{
		ElectionPeriod:         t.GetElectionPeriod(sdkCtx),
		SequencerNumber:        t.GetSequencerNumber(sdkCtx),
		BackupSequencerNumber:  t.GetBackupNumber(sdkCtx),
		MinStakeAmount:         t.GetMinStakeAmount(sdkCtx),
		FirstElectionInterval:  t.GetFirstElectionInterval(sdkCtx),
		AllowApplyElectionTime: t.GetAllowApplyElectionTime(sdkCtx),
		ElectionInterimTime:    t.GetElectionInterimTime(sdkCtx),
	}
	return &types.QueryParamsResponse{
		Params: paramas,
	}, nil

}

func (t Keeper) QueryElectionResult(ctx context.Context, req *types.QueryElectionRequest) (*types.QueryElectionResponse, error) {
	if req.RollappId != t.rollAppID {
		return nil, errorsmod.Wrapf(types.ErrRollappIDMismatch, fmt.Sprintf("rollupServer's rollappID = %s", t.rollAppID))
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := sdkCtx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(t.rollAppID))
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
	if rollappID != t.rollAppID {
		return nil, errorsmod.Wrapf(types.ErrRollappIDMismatch, fmt.Sprintf("rollupServer's rollappID = %s", t.rollAppID))
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := sdkCtx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(t.rollAppID))
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
	if req.RollappId != t.rollAppID {
		return nil, errorsmod.Wrapf(types.ErrRollappIDMismatch, fmt.Sprintf("rollupServer's rollappID = %s", t.rollAppID))
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := sdkCtx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppStakeKeyPrefix(t.rollAppID))
	data := store.Get([]byte(req.Address))
	resp := &types.MsgStakeInfo{
		StakeAmount:        0,
		ApplyUnStakeAmount: 0,
	}
	if err := t.cdc.Unmarshal(data, resp); err != nil {
		return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Unmarshal error. err = %s", err.Error()))
	}

	return &types.QueryStakeResponse{
		StakeInfo: resp,
	}, nil
}
