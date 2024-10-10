package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/rollup/types"
	"math/big"
	"strconv"
)

func (k Keeper) StakeForChallengeDaFraud(goCtx context.Context, rollappID, blockSubmitter, challenger string, challengeKey []byte) error {
	challengerAddr, err := sdk.AccAddressFromBech32(challenger)
	if err != nil {
		return errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("AccAddressFromBech32 error. err = %s", err.Error()))
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.IsInBlackList(challenger) {
		return errorsmod.Wrapf(types.ErrInBlackList, "")
	}
	//一次DA欺诈挑战只质押一次，此时的Key为challengeKey，如果之前已经存在，则认为之前已经质押，此时出错
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetStakeForChallengeDaFraudPrefix(rollappID))
	data := store.Get(challengeKey)
	if data != nil {
		return errorsmod.Wrapf(types.ErrStakeDaFraudRepeated, fmt.Sprintf("challengeKey = %s", hex.EncodeToString(challengeKey)))
	}
	/*
		balanceCoin := k.bk.GetBalance(ctx, challengerAddr, params.BaseDenom)
		if balanceCoin.IsLT(stakeCoin) {
			return errorsmod.Wrapf(types.ErrInsufficientBalance, fmt.Sprintf(",user's balanceCoin = %s,but need stake= %dumec",
				balanceCoin.String(), stakeForChallenge))
		}*/
	stakeForChallenge := uint64(k.GetDaFraudChallengeStake(ctx)) * types.MecPrecision
	stakeCoin := sdk.NewCoin(params.BaseDenom, sdk.NewInt(int64(stakeForChallenge)))
	//如果金额不够的话，SendCoinsFromAccountToModule这里就已经会判断处理了
	if err = k.bk.SendCoinsFromAccountToModule(ctx, challengerAddr, types.MODULE_NAME, sdk.NewCoins(stakeCoin)); err != nil {
		return errorsmod.Wrapf(types.ErrStakeDataErr, fmt.Sprintf("stake coin to module error.err = %s", err.Error()))
	}
	stakeMsg := &types.MsgStakeChallengeDaFraud{
		Challenger:     challenger,
		BlockSubmitter: blockSubmitter,
		Denom:          params.BaseDenom,
		Amount:         stakeForChallenge,
	}
	store.Set(challengeKey, k.cdc.MustMarshal(stakeMsg))
	return nil

}

func (k Keeper) ProcChallengeDaFraud(goCtx context.Context, rollappID string, challengeKey []byte, result int32) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetStakeForChallengeDaFraudPrefix(rollappID))
	data := store.Get(challengeKey)
	if nil == data {
		return errorsmod.Wrapf(types.ErrNotFound, fmt.Sprintf("can not found stake info in ProcChallengeDaFraud.rollappID = %s, challengeKey = %s",
			rollappID, string(challengeKey)))
	}
	stakeMsg := new(types.MsgStakeChallengeDaFraud)
	k.cdc.MustUnmarshal(data, stakeMsg)

	if stakeMsg.Amount < 1*types.MecPrecision {
		return errorsmod.Wrapf(types.ErrInsufficientBalance, fmt.Sprintf("challenger stake balances insufficient.stakeAmount = %s",
			strconv.FormatUint(stakeMsg.Amount, 10)))
	}
	challengerAccAddr, err := sdk.AccAddressFromBech32(stakeMsg.Challenger)
	if err != nil {
		return errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("AccAddressFromBech32 error. err = %s,addr = %s,key = %s",
			err.Error(), stakeMsg.Challenger, string(challengeKey)))
	}

	stakeCoin := sdk.NewCoin(params.BaseDenom, sdk.NewInt(int64(stakeMsg.Amount)))
	if types.RESULT_CHG_FAIL == result || types.RESULT_CHG_SUCCESS_SUBMIT_DATA_FAIL == result {
		//如果挑战失败，或者提交错误数据，则扣除挑战者的质押资金(由于之前已经质押了资金，这里只要清除质押记录即可)
		//并将质押金额的20%奖励给验证器地址
		rewardValidatorCoin := sdk.NewCoin(params.BaseDenom, sdk.NewInt(int64(stakeMsg.Amount/5)))
		validatorAccAddr, errP := sdk.AccAddressFromBech32(k.dk.GetValidatorAddress(ctx))
		if errP != nil {
			return errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("AccAddressFromBech32 validator_address error.err = %s,addr = %s ",
				errP.Error(), k.dk.GetValidatorAddress(ctx)))
		}
		if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.MODULE_NAME, validatorAccAddr, sdk.NewCoins(rewardValidatorCoin)); err != nil {
			return errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("reward coin to validator_address error.err = %s, addr = %s, amount = %s",
				err.Error(), validatorAccAddr.String(), rewardValidatorCoin.String()))
		}
		//记录事件
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EvtRewardDaFraudValidator,
				sdk.NewAttribute("moduleName", types.MODULE_NAME),
				sdk.NewAttribute("rollappID", rollappID),
				sdk.NewAttribute("from", stakeMsg.Challenger),
				sdk.NewAttribute("validatorAddr", validatorAccAddr.String()),
				sdk.NewAttribute("rewardAmount", rewardValidatorCoin.String()),
			),
		)
		store.Delete(challengeKey)
		ctx.Logger().Info("verify result: challenge fraud")
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EvtPunishDaChallengerFraud,
				sdk.NewAttribute("moduleName", types.MODULE_NAME),
				sdk.NewAttribute("rollappID", rollappID),
				sdk.NewAttribute("challenger", stakeMsg.Challenger),
				sdk.NewAttribute("challengeKey", string(challengeKey)),
				sdk.NewAttribute("amount", strconv.FormatUint(stakeMsg.Amount, 10)),
			),
		)
	} else if types.RESULT_CHG_SUCCESS_SUBMIT_DATA_SUCESS == result {
		//如果挑战成功，并且提交的数据也没问题，则先返回挑战者质押的资金，
		//同时如果blockSubmitter是sequencer的话，则提出sequencer队列
		ctx.Logger().Info("verify result: blockSubmitter da fraud")
		//删除质押记录，返回质押资金
		store.Delete(challengeKey)
		if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.MODULE_NAME, challengerAccAddr, sdk.NewCoins(stakeCoin)); err != nil {
			return errorsmod.Wrapf(types.ErrUnStakeProc, fmt.Sprintf("unStake challenger coin error. err = %s", err.Error()))
		}

		//加入黑名单，进行DaFraud的统计
		if err = k.AddToBlackList(ctx, stakeMsg.BlockSubmitter); err != nil {
			return err
		}
		//这里设置欺诈者罚款的金额和 欺诈的区块=>挑战者 的信息

		punishAmount := uint64(0)
		daFraudStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetDaFraudStaticsPrefix(rollappID))
		needPunishmentData := daFraudStore.Get([]byte(types.KeyNeedPunishment))
		var mapPunishment map[string]uint64

		if nil != needPunishmentData {
			if err = json.Unmarshal(needPunishmentData, &mapPunishment); err != nil {
				return errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Unmarshal data error.err = %s,key = %s",
					err.Error(), types.KeyNeedPunishment))
			}
		} else {
			mapPunishment = make(map[string]uint64)
		}
		if _, ok := mapPunishment[stakeMsg.BlockSubmitter]; !ok { //查看之前是否惩罚过了,由于一旦证实作恶，就会被加到黑名单，所以只需要扣除一次即可
			ctx.Logger().Info(fmt.Sprintf("punishment:send stake coin to module,amount = %d", punishAmount))

			//扣除之前的blockSubmitter所质押的资金来作为奖励池
			punishAmount, err = k.Punishment(ctx, stakeMsg.BlockSubmitter, rollappID, 100, 0)
			if err != nil {
				return err
			}
			mapPunishment[stakeMsg.BlockSubmitter] = punishAmount
			resData, err := json.Marshal(mapPunishment)
			if err != nil { //对交易的执行待验证
				return errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Marshal mapPunishment data error.err = %s,"+
					"blockSubmitter = %s,challengeKey = %s", err.Error(), stakeMsg.BlockSubmitter, string(challengeKey)))
			}
			daFraudStore.Set([]byte(types.KeyNeedPunishment), resData)
			//daFraudStore.Set(types.GetDaFraudPunishmentKey(stakeMsg.BlockSubmitter), []byte(strconv.FormatUint(punishAmount, 10)))

			//有可能是Sequencer，所以此时要对sequencer进行重估,这里不应该返回错误
			if err = k.RevaluateSequencer(ctx, stakeMsg.BlockSubmitter, rollappID); err != nil {
				return errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("RevaluateSequencer in ProcChallengeDaFraud error."+
					"err = %s,BlockSubmitter = %s,rollappID = %s", err.Error(), stakeMsg.BlockSubmitter, rollappID))
			}
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EvtPunishBlockDaSubmitter,
					sdk.NewAttribute("moduleName", types.MODULE_NAME),
					sdk.NewAttribute("rollappID", rollappID),
					sdk.NewAttribute("challenger", stakeMsg.Challenger),
					sdk.NewAttribute("challengeKey", string(challengeKey)),
					sdk.NewAttribute("blockSubmitter", stakeMsg.BlockSubmitter),
					sdk.NewAttribute("amount", strconv.FormatUint(punishAmount, 10)),
				),
			)

		}
		//挑战的具体细则 keyPrefix+区块提交者+challengeKey ==> 挑战者
		daFraudKey := append(types.GetProvedDaFraudKeyPrefix(stakeMsg.BlockSubmitter), challengeKey...)
		daFraudStore.Set(daFraudKey, []byte(stakeMsg.Challenger))

	}
	return nil

}
func (k *Keeper) InitPunishInfo(ctx sdk.Context) error {
	k.mapPunishInfo = nil
	for key, _ := range k.mapRollappInfoMng {
		rollAppID := key
		daFraudStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetDaFraudStaticsPrefix(rollAppID))
		needPunishmentData := daFraudStore.Get([]byte(types.KeyNeedPunishment))
		if nil != needPunishmentData {
			var punishInfo map[string]uint64
			if err := json.Unmarshal(needPunishmentData, &punishInfo); err != nil {
				k.mapPunishInfo = nil
				return fmt.Errorf("Unmarshal data to punishmentInfo error.err = %s,rollappID = %s", err.Error(), rollAppID)
			}
			if nil == k.mapPunishInfo {
				k.mapPunishInfo = make(map[string]map[string]uint64)
			}
			k.mapPunishInfo[rollAppID] = punishInfo
		} else {
			continue
		}
	}
	return nil

}

func (k Keeper) RewardsChallengeDaFraud(ctx sdk.Context) error {
	if nil == k.mapPunishInfo { //这里表示没有惩罚数据
		return nil
	}
	validatorAccAddr, errP := sdk.AccAddressFromBech32(k.dk.GetValidatorAddress(ctx))
	if errP != nil {
		return errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("AccAddressFromBech32 validator_address error.err = %s,addr = %s ",
			errP.Error(), k.dk.GetValidatorAddress(ctx)))
	}
	for rollappId, punishInfo := range k.mapPunishInfo {
		var procFraudster []string
		for key, val := range punishInfo {
			submitTime := k.rk.GetSubmitterLastSubmitTime(ctx, rollappId, key)
			blkTime := ctx.BlockTime().Unix()
			//7200秒为额外附加冗余的时间
			if blkTime >= (submitTime + int64(types.SubmitDaFraudTime)*types.HourSeconds + 7200) {
				//开始根据统计信息发放奖励
				rewardsAmount := (val * 4) / 5
				err := k.rewardBaseDeomToDaFraudChallenge(ctx, rollappId, key, rewardsAmount)
				if err != nil {
					return err
				}
				//剩下的部分发送给链下验证器地址
				remainCoin := sdk.NewCoin(params.BaseDenom, sdk.NewIntFromBigInt(big.NewInt(int64(val/5))))
				if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.MODULE_NAME, validatorAccAddr, sdk.NewCoins(remainCoin)); err != nil {
					return errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("reward coin to validator_address error.err = %s, addr = %s, amount = %s",
						err.Error(), validatorAccAddr.String(), remainCoin.String()))
				}
				//记录事件
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EvtRewardDaFraudValidator,
						sdk.NewAttribute("moduleName", types.MODULE_NAME),
						sdk.NewAttribute("rollappID", rollappId),
						sdk.NewAttribute("from", key),
						sdk.NewAttribute("validatorAddr", validatorAccAddr.String()),
						sdk.NewAttribute("rewardAmount", remainCoin.String()),
					),
				)
				store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetDaFraudStaticsPrefix(rollappId))
				store.Set(types.GetFinishDaFraudPunishKey(key), []byte(strconv.FormatUint(val, 10)))
				procFraudster = append(procFraudster, key)
			}
		}
		isNeedWrite := false
		//修改mapPunishInfo信息,写入ival
		for _, v := range procFraudster {
			delete(punishInfo, v)
			isNeedWrite = true
		}
		if isNeedWrite {
			resData, err := json.Marshal(punishInfo)
			if err != nil {
				return errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Marshal punishmentInfo in RewardsChallengeDaFraud error."+
					" err = %s", err.Error()))
			}
			daFraudStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetDaFraudStaticsPrefix(rollappId))
			daFraudStore.Set([]byte(types.KeyNeedPunishment), resData)
		}
	}

	return nil

}

func (k Keeper) rewardBaseDeomToDaFraudChallenge(ctx sdk.Context, rollappID, fraudsterAddr string, rewardAmount uint64) error {
	ctx.Logger().Info(fmt.Sprintf("enter rewardBaseDeomToDaFraudChallenge.rollappID = %s,fraudster = %s,amount = %d",
		rollappID, fraudsterAddr, rewardAmount))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetDaFraudStaticsPrefix(rollappID))
	iterator := sdk.KVStorePrefixIterator(store, types.GetProvedDaFraudKeyPrefix(fraudsterAddr))
	defer iterator.Close() // nolint: errcheck
	mapChallengeStatics := make(map[string]uint32)
	totalFraudNumber := uint32(0)
	for ; iterator.Valid(); iterator.Next() {
		challenger := string(iterator.Value())
		totalFraudNumber++
		if val, ok := mapChallengeStatics[challenger]; ok {
			val++
			mapChallengeStatics[challenger] = val
		} else {
			mapChallengeStatics[challenger] = 1
		}
	}
	amount := big.NewInt(int64(rewardAmount))
	totalFraudNumberVal := big.NewInt(int64(totalFraudNumber))
	totalSend := big.NewInt(0)
	for addr, v := range mapChallengeStatics {
		accAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("AccAddressFromBech32  in rewardBaseDeomToDaFraudChallenge error."+
				" err = %s,addr = %s", err.Error(), addr))
		}

		tmp := big.NewInt(0).Mul(amount, big.NewInt(int64(v)))
		profits := tmp.Div(tmp, totalFraudNumberVal)
		profitCoin := sdk.NewCoin(params.BaseDenom, sdk.NewIntFromBigInt(profits))
		totalSend.Add(totalSend, profits)
		//发放奖励，如果金额不够的话，SendCoinsFromModuleToAccount这里就已经会判断处理了
		if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.MODULE_NAME, accAddr, sdk.NewCoins(profitCoin)); err != nil {
			return errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("reward coin to challenger error.err = %s, addr = %s, amount = %s",
				err.Error(), accAddr.String(), profitCoin.String()))
		}
		//记录事件
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EvtRewardDaFraudChallenger,
				sdk.NewAttribute("moduleName", types.MODULE_NAME),
				sdk.NewAttribute("rollappID", rollappID),
				sdk.NewAttribute("fraudster", fraudsterAddr),
				sdk.NewAttribute("challenger", addr),
				sdk.NewAttribute("rewardAmount", profitCoin.String()),
				sdk.NewAttribute("totalAmount", amount.String()),
			),
		)
	}
	if totalSend.Cmp(amount) > 0 {
		panic(fmt.Errorf("totalSend rewards > total rewards.totalSend = %s,totalReward = %s",
			totalSend.String(), amount.String()))
	}
	return nil

}
