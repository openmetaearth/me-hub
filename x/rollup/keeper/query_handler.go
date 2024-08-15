package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/x/rollup/types"
)

func (t rollupQueryServer) QueryParams(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	paramas := types.Params{
		ElectionPeriod:         t.Keeper.GetElectionPeriod(sdkCtx),
		SequencerNumber:        t.Keeper.GetSequencerNumber(sdkCtx),
		BackupSequencerNumber:  t.Keeper.GetBackupNumber(sdkCtx),
		MinStakeAmount:         t.Keeper.GetMinStakeAmount(sdkCtx),
		FirstElectionInterval:  t.Keeper.GetFirstElectionInterval(sdkCtx),
		AllowApplyElectionTime: t.Keeper.GetAllowApplyElectionTime(sdkCtx),
	}
	return &types.QueryParamsResponse{
		Params: paramas,
	}, nil

}

func (t rollupQueryServer) QueryElectionResult(ctx context.Context, req *types.QueryElectionRequest) (*types.QueryElectionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := sdkCtx.KVStore(t.Keeper.storeKey)
	store := prefix.NewStore(kvStore, []byte(types.RollupKeyPrefix))
	data := store.Get([]byte(types.KEY_LAST_ELECTION_INFO))
	resp := &types.QueryElectionResponse{
		ElectionTime:   0,
		BlockHeight:    0,
		NodeStatusList: nil,
	}
	if err := t.Keeper.cdc.Unmarshal(data, resp); err != nil {
		return nil, fmt.Errorf("%s,Unmarshal error. err = %s", types.ErrParserDataErr, err.Error())
	}
	return resp, nil
}

func (t rollupQueryServer) QueryStake(ctx context.Context, req *types.QueryStakeRequest) (*types.QueryStakeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := sdkCtx.KVStore(t.Keeper.storeKey)
	store := prefix.NewStore(kvStore, []byte(types.RollupStakeKeyPrefix))
	data := store.Get([]byte(req.Address))
	resp := &types.MsgStakeInfo{
		StakeAmount:        0,
		ApplyUnStakeAmount: 0,
	}
	if err := t.Keeper.cdc.Unmarshal(data, resp); err != nil {
		return nil, fmt.Errorf("%s,Unmarshal error. err = %s", types.ErrParserDataErr, err.Error())
	}
	if resp.StakeAmount > 0 {
		resp.StakeAmount = resp.StakeAmount / types.MecPrecision
	}
	if resp.ApplyUnStakeAmount > 0 {
		resp.ApplyUnStakeAmount = resp.ApplyUnStakeAmount / types.MecPrecision
	}
	return &types.QueryStakeResponse{
		StakeInfo: resp,
	}, nil
}
