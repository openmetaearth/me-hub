package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/rollup/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.Params{
		ElectionPeriod:        k.GetElectionPeriod(ctx),
		SequencerNumber:       k.GetSequencerNumber(ctx),
		BackupSequencerNumber: k.GetBackupNumber(ctx),
		MinStakeAmount:        k.GetMinStakeAmount(ctx),
		//	FirstElectionInterval:  k.GetFirstElectionInterval(ctx),
		AllowApplyElectionTime: k.GetAllowApplyElectionTime(ctx),
		ElectionInterimTime:    k.GetElectionInterimTime(ctx),
		DaFraudChallengeStake:  k.GetDaFraudChallengeStake(ctx),
	}

}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramStore.SetParamSet(ctx, &params)
}

func (k Keeper) SetRollupParams(ctx context.Context, req *types.MsgSetRollupParamsRequest) (*types.MsgSetRollupParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.Logger().Info(fmt.Sprintf("SetRollupParams,new Params = %s", req.String()))

	if req.Creator != k.dk.GetDevOperator(sdkCtx) {
		return nil, errorsmod.Wrapf(types.ErrInputDataErr, "creator has not right to set params")
	}
	k.paramStore.SetParamSet(sdkCtx, req.NewParams)
	return &types.MsgSetRollupParamsResponse{}, nil
}
