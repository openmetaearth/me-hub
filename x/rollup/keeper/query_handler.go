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
	rollAppID := t.getRollappID(sdkCtx)
	if req.RollappId != rollAppID {
		return nil, errorsmod.Wrapf(types.ErrRollappIDMismatch, fmt.Sprintf("rollupServer's rollappID = %s", rollAppID))
	}

	paramas := types.Params{
		ElectionPeriod:         t.GetElectionPeriod(sdkCtx),
		SequencerNumber:        t.GetSequencerNumber(sdkCtx),
		BackupSequencerNumber:  t.GetBackupNumber(sdkCtx),
		MinStakeAmount:         t.GetMinStakeAmount(sdkCtx),
		FirstElectionInterval:  t.GetFirstElectionInterval(sdkCtx),
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
	rollAppID := t.getRollappID(sdkCtx)
	if req.RollappId != rollAppID {
		return nil, errorsmod.Wrapf(types.ErrRollappIDMismatch, fmt.Sprintf("rollupServer's rollappID = %s", rollAppID))
	}

	kvStore := sdkCtx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(rollAppID))
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
	tRollAppID := t.getRollappID(sdkCtx)
	if rollappID != tRollAppID {
		return nil, errorsmod.Wrapf(types.ErrRollappIDMismatch, fmt.Sprintf("rollupServer's rollappID = %s", tRollAppID))
	}
	kvStore := sdkCtx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(tRollAppID))
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
	rollAppID := t.getRollappID(sdkCtx)
	if req.RollappId != rollAppID {
		sdkCtx.Logger().Error(fmt.Sprintf("reqRollappID = %s,t.RollappID = %s", req.RollappId, rollAppID))
		return nil, errorsmod.Wrapf(types.ErrRollappIDMismatch, fmt.Sprintf("rollupServer's rollappID = %s", rollAppID))
	}

	sdkCtx.Logger().Info(fmt.Sprintf("QueryStake,rollappID = %s", rollAppID))
	kvStore := sdkCtx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppStakeKeyPrefix(rollAppID))
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

func (t Keeper) getRollappID(ctx sdk.Context) string {
	if t.rollAppID != "" {
		return t.rollAppID
	}
	kvStore := ctx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, []byte(types.RollupKeyPrefix))
	data := store.Get([]byte(types.KEY_ROLLAPP_ID))
	if data != nil {
		return string(data)
	}
	return ""
}
