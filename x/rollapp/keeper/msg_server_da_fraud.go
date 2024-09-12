package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/rollapp/types"
	rollupTypes "github.com/st-chain/me-hub/x/rollup/types"
	"strconv"
)

/*

2、判断用户的当前的余额是否足够(每一次挑战都会质押一定的金额)
3、判断是否在挑战时间允许的范围内
*/

func (k msgServer) ChallengeDaFraud(goCtx context.Context, req *types.MsgSubmitDaFraudRequest) (*types.MsgSubmitDaFraudResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//ctx.Logger().Info("receviced SubmitBlockDAInfo", "msg", req.String())
	if !k.RollappsEnabled(ctx) {
		return nil, types.ErrRollappsDisabled
	}
	isFound := k.IsRollappExist(ctx, req.RollappId)
	if !isFound {
		return nil, types.ErrUnknownRollappID
	}
	//判断改区块是否已经被挑战
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetDaChallengeKeyPrefix(req.RollappId))
	challengeRollupBlkKey := types.GetRollupBlockKey(req.StartHeight, req.NumBlocks)
	data := store.Get(challengeRollupBlkKey)
	if data != nil {
		return nil, errorsmod.Wrapf(types.ErrDaFraudConflict, "da fraud challenge has been exist")
	}
	//判断是否在挑战时间允许的范围内
	cmtBlkInfo := k.getCommitBlockDaInfo(goCtx, req.RollappId, req.StartHeight, req.NumBlocks)
	if nil == cmtBlkInfo {
		return nil, errorsmod.Wrapf(types.ErrNotFound, fmt.Sprintf("commit block not found. startHeight = %d,number = %d",
			req.StartHeight, req.NumBlocks))
	}
	if ctx.BlockTime().Unix() > (int64(cmtBlkInfo.BlockTime) + int64(rollupTypes.SubmitDaFraudTime)*rollupTypes.HourSeconds) {
		return nil, errorsmod.Wrapf(types.ErrDaFraudTimeout, "da fraud challenge has been timeout")
	}
	//判断用户的当前的余额是否足够(每一次挑战都会质押一定的金额)
	err := k.rollupKeeper.StakeForDaFraud(ctx, req.RollappId, req.Creator, challengeRollupBlkKey)
	if err != nil {
		return nil, err
	}
	//
	challengeData := &types.MsgDaFraudChallengerDataStatus{
		Challenger:          req.Creator,
		SubmitCommitment:    req.Commitment,
		SubmitDaRoot:        req.DaRoot,
		SubmitDaBlockHeight: req.DaBlockHeight,
		Status:              types.STATUS_CHG_DA_FRAUD_ING,
	}
	store.Set(challengeRollupBlkKey, k.cdc.MustMarshal(challengeData))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventDAFraudChallenge,
			sdk.NewAttribute("moduleName", types.ModuleName),
			sdk.NewAttribute("challenger", req.Creator),
			sdk.NewAttribute("rollappId", req.RollappId),
			sdk.NewAttribute("startHeight", strconv.FormatUint(req.StartHeight, 10)),
			sdk.NewAttribute("number", strconv.FormatUint(uint64(req.NumBlocks), 10)),
			sdk.NewAttribute("namespace", string(req.Namespace)),
			sdk.NewAttribute("challengeBlockHeight", strconv.FormatUint(cmtBlkInfo.BlkDaInfo.DaBlockHeight, 10)),
			sdk.NewAttribute("challengeCommitment", hex.EncodeToString(cmtBlkInfo.BlkDaInfo.Commitment)),
			sdk.NewAttribute("submitCommitment", hex.EncodeToString(req.Commitment)),
			sdk.NewAttribute("submitBlockHeight", strconv.FormatUint(req.DaBlockHeight, 10)),
			sdk.NewAttribute("submitDaRoot", hex.EncodeToString(req.DaRoot)),
		),
	)
	return &types.MsgSubmitDaFraudResponse{}, nil

}

func (k msgServer) SubmitDaFraudVerifyData(context.Context, *types.MsgDaFraudVerifyResult) (*types.MsgDaFraudVerifyResult, error) {

}
