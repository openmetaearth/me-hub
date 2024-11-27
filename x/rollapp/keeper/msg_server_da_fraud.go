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
	//todo: for test
	/*
		if !k.RollappsEnabled(ctx) {
			return nil, types.ErrRollappsDisabled
		}
		isFound := k.IsRollappExist(ctx, req.RollappId)
		if !isFound {
			return nil, types.ErrUnknownRollappID
		}
	*/
	//=================end
	if k.rollupKeeper.IsInBlackList(req.Creator) {
		return nil, errorsmod.Wrapf(rollupTypes.ErrInBlackList, "")
	}
	//判断改区块是否已经被挑战
	ctx.Logger().Info("enter ChallengeDaFraud")
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
	err := k.rollupKeeper.StakeForChallengeDaFraud(ctx, req.RollappId, cmtBlkInfo.BlkDaInfo.Creator, req.Creator, challengeRollupBlkKey)
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
	ctx.Logger().Info("end ChallengeDaFraud")
	return &types.MsgSubmitDaFraudResponse{}, nil

}

func (k msgServer) SubmitDaFraudVerifyData(goCtx context.Context, req *types.MsgDaFraudVerifyResult) (*types.MsgDaFraudVerifyResultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.daoKeeper.IsGlobalDao(ctx, req.Creator) || k.daoKeeper.IsValidatorDao(ctx, req.Creator) {
		//todo: for test
		/*
			isFound := k.IsRollappExist(ctx, req.RollappId)
			if !isFound {
				return nil, types.ErrUnknownRollappID
			}
		*/
		//=============end
		//判断改区块是否已经被挑战
		ctx.Logger().Info("enter SubmitDaFraudVerifyData")
		store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetDaChallengeKeyPrefix(req.RollappId))
		challengeRollupBlkKey := types.GetRollupBlockKey(req.StartHeight, req.NumBlocks)
		data := store.Get(challengeRollupBlkKey)
		if nil == data { //如果没有改挑战的记录
			return nil, errorsmod.Wrapf(types.ErrNotFound, fmt.Sprintf("da fraud challenge was not found.rollappId = %s,startBlockHeight = %d,numBlks = %d",
				req.RollappId, req.StartHeight, req.NumBlocks))
		}
		challengeData := new(types.MsgDaFraudChallengerDataStatus)
		err := k.cdc.Unmarshal(data, challengeData)
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrParserData, "Unmarshal data to MsgDaFraudChallengerDataStatus error.err = %s", err.Error())
		}
		if challengeData.Status != types.STATUS_CHG_DA_FRAUD_ING { //如果该挑战的状态不是正在处理，则返回错误
			return nil, errorsmod.Wrapf(types.ErrDaFraudConflict,
				fmt.Sprintf("challenge has been verified.status = %d", challengeData.Status))
		}
		//开始进行处理
		if err = k.rollupKeeper.ProcChallengeDaFraud(ctx, req.RollappId, challengeRollupBlkKey, req.Result); err != nil {
			return nil, err
		}
		//如果处理成功，则根据verify的验证状态进行判断
		if req.Result == rollupTypes.RESULT_CHG_FAIL { //挑战失败，则说明原来的数据为正确
			challengeData.Status = types.STATUS_CHG_DA_FRAUD_FAIL
			store.Set(challengeRollupBlkKey, k.cdc.MustMarshal(challengeData))

		} else if req.Result == rollupTypes.RESULT_CHG_SUCCESS_SUBMIT_DATA_SUCESS {
			//挑战成功，并且提交的数据也正确，则更新之前存储的数据DA信息
			challengeData.Status = types.STATUS_CHG_DA_FRAUD_SUCCESS
			store.Set(challengeRollupBlkKey, k.cdc.MustMarshal(challengeData))

			blkStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupBlockKeyPrefix(req.RollappId))
			cmtBlkData := blkStore.Get(challengeRollupBlkKey)
			if nil == cmtBlkData {
				return nil, errorsmod.Wrapf(types.ErrLogic, fmt.Sprintf("proc challengeDaFraud success,but can not found commit block data."+
					"commitBlockkey = %s", string(challengeRollupBlkKey)))
			}

			cmtBlkInfo := new(types.MsgCommitBlockDaInfo)
			k.cdc.MustUnmarshal(cmtBlkData, cmtBlkInfo)
			cmtBlkInfo.BlkDaInfo.Commitment = challengeData.SubmitCommitment
			cmtBlkInfo.BlkDaInfo.DaBlockHeight = challengeData.SubmitDaBlockHeight
			cmtBlkInfo.BlkDaInfo.DaRoot = challengeData.SubmitDaRoot
			//挑战之后就不在需要CommitmentProof
			cmtBlkInfo.BlkDaInfo.CommitmentProof = nil
			blkStore.Set(challengeRollupBlkKey, k.cdc.MustMarshal(cmtBlkInfo))
		} else { //如果是挑战成功，但是挑战者提交的数据是错误的，那就删除挑战记录，让这个可以继续被挑战
			store.Delete(challengeRollupBlkKey)
		}
		ctx.Logger().Info("end SubmitDaFraudVerifyData")
		return &types.MsgDaFraudVerifyResultResponse{}, nil

	} else {
		return nil, errorsmod.Wrapf(types.ErrInvalidAddress, "address is not allowed")
	}

}

func (k msgServer) QueryElectionResult(ctx context.Context, req *rollupTypes.QueryElectionRequest) (*rollupTypes.QueryElectionResponse, error) {
	return k.Keeper.rollupKeeper.QueryElectionResult(ctx, req)
}

func (k Keeper) GetSubmitBlockDaInfo(goCtx context.Context, req *types.MsgGetBlockDaInfoRequest) (*types.MsgGetBlockDaInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	blkStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupBlockKeyPrefix(req.RollappId))
	blkKey := types.GetRollupBlockKey(req.StartHeight, req.NumberBlocks)
	cmtBlkData := blkStore.Get(blkKey)
	if nil == cmtBlkData {
		return nil, errorsmod.Wrapf(types.ErrLogic, fmt.Sprintf("can not found commit block data."+
			"commitBlockkey = %s", string(blkKey)))
	}

	cmtBlkInfo := new(types.MsgCommitBlockDaInfo)
	k.cdc.MustUnmarshal(cmtBlkData, cmtBlkInfo)
	return &types.MsgGetBlockDaInfoResponse{
		Blocks:        cmtBlkInfo.BlkDaInfo.Blocks,
		Namespace:     k.mapRollappAssociateDa[req.RollappId],
		DaRoot:        cmtBlkInfo.BlkDaInfo.DaRoot,
		Commitment:    cmtBlkInfo.BlkDaInfo.Commitment,
		DaBlockHeight: cmtBlkInfo.BlkDaInfo.DaBlockHeight,
	}, nil
}

// 查询所有的为处理挑战信息
func (k Keeper) GetAppendingDaFraudChallenge(goCtx context.Context, req *types.MsgGetDaFraudChallengeRequest) (*types.MsgGetDaFraudChallengeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var resp []*types.MsgSubmitDaFraudRequest
	for key, val := range k.mapRollappAssociateDa {
		appendChallenge, err := k.getAppendDaFraudChallengeByRollappID(ctx, key, val)
		if err != nil {
			return nil, err

		}
		if nil == appendChallenge || len(appendChallenge) < 1 {
			continue
		}
		resp = append(resp, appendChallenge...)

	}
	return &types.MsgGetDaFraudChallengeResponse{DaFraudChallengeList: resp}, nil

}

func (k Keeper) getAppendDaFraudChallengeByRollappID(ctx sdk.Context, rollappID string, namespace []byte) ([]*types.MsgSubmitDaFraudRequest, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetDaChallengeKeyPrefix(rollappID))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close() // nolint: errcheck
	var res []*types.MsgSubmitDaFraudRequest

	for ; iterator.Valid(); iterator.Next() {
		if nil == iterator.Value() { //这里不应该存在这种情况
			return nil, errorsmod.Wrapf(types.ErrLogic, fmt.Sprintf("da fraud challenge's value is nil."+
				"rollappId = %s,key = %s", rollappID, string(iterator.Key())))
		}
		challengeData := new(types.MsgDaFraudChallengerDataStatus)
		k.cdc.MustUnmarshal(iterator.Value(), challengeData)
		if types.STATUS_CHG_DA_FRAUD_ING == challengeData.Status {
			startBlkHeight, numberBlk, err := types.ParserRollupKey(string(iterator.Key()))
			if err != nil {
				return nil, errorsmod.Wrapf(types.ErrParserData, fmt.Sprintf(" types.ParserRollupKey error,err = %s,key = %s",
					err.Error(), string(iterator.Key())))
			}
			daFraudReq := &types.MsgSubmitDaFraudRequest{
				Creator:       challengeData.Challenger,
				RollappId:     rollappID,
				StartHeight:   startBlkHeight,
				NumBlocks:     numberBlk,
				Namespace:     namespace,
				Commitment:    challengeData.SubmitCommitment,
				DaRoot:        challengeData.SubmitDaRoot,
				DaBlockHeight: challengeData.GetSubmitDaBlockHeight(),
			}
			res = append(res, daFraudReq)
		}
	}
	return res, nil
}
