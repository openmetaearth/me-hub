package keeper

import (
	"encoding/json"
	"fmt"
	//"github.com/Workiva/go-datastructures/threadsafe/err"
	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/st-chain/me-hub/x/rollup/types"
	"sort"
	"strconv"
)

// Keeper struct
type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	bk         types.BankKeeper
	rk         types.RollappKeeper
	dk         types.DaoKeeper
	paramStore paramtypes.Subspace
	//lastElectionTime uint64
	rollAppID string
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(storeKey storetypes.StoreKey, cdc codec.BinaryCodec, paramSpace paramtypes.Subspace,
	bKeeper types.BankKeeper, rKeeper types.RollappKeeper, dKeeper types.DaoKeeper) *Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	return &Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramStore: paramSpace,
		rollAppID:  "",
		bk:         bKeeper,
		rk:         rKeeper,
		dk:         dKeeper,
	}
}

// Logger returns a logger instance for the incentives module.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.MODULE_NAME))
}

func (k *Keeper) ProcElection(ctx sdk.Context) error {
	blkTime := ctx.BlockTime().Unix()
	//获取上一次选举的时间
	kvStore := ctx.KVStore(k.storeKey)
	rollupStore := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(k.rollAppID))
	lastElectTime := int64(0)
	bIsNeedElect := false
	if electTimeVal := rollupStore.Get([]byte(types.KEY_LAST_ELECTION_TIME)); electTimeVal != nil {
		lastElectTime = types.BytesToInt64(electTimeVal)
		timeInterval := blkTime - lastElectTime
		electionInterval := int64(k.GetElectionPeriod(ctx)) * types.MinuteSeconds
		if timeInterval >= electionInterval {
			bIsNeedElect = true
		}

	} else { //找不到lastElectTime的话，则表示还没竞选过
		if 1 == ctx.BlockHeight() { //如果是第一个数据区块的话，则计算首次竞选的时间并保存
			firstElectTime := blkTime + int64(k.GetFirstElectionInterval(ctx))*types.MinuteSeconds
			rollupStore.Set([]byte(types.KEY_FIRST_ELECTION_TIME), types.Int64ToBytes(firstElectTime))
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EvtFirstElectionTime,
					sdk.NewAttribute("moduleName", types.MODULE_NAME),
					sdk.NewAttribute("firstElectTime", strconv.FormatInt(firstElectTime, 10)),
				),
			)
			return nil
		} else if ctx.BlockHeight() > 1 {
			firstElectTime := int64(0)
			if timeVal := rollupStore.Get([]byte(types.KEY_FIRST_ELECTION_TIME)); timeVal != nil {
				firstElectTime = types.BytesToInt64(timeVal)
				if blkTime >= firstElectTime {
					bIsNeedElect = true
				} else {
					return nil
				}
			} else {
				return errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("can not get firstElectTime. blkHeight = %d",
					ctx.BlockHeight()))
			}

		}

	}
	if bIsNeedElect { //开始竞选
		electList, err := k.startElection(ctx, k.GetMinStakeAmount(ctx)*types.MecPrecision)
		if err != nil {
			return err
		}
		var res []byte
		if res, err = json.Marshal(electList); err != nil {
			return errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("Marshal(electList) error.err = %s", err.Error()))
		}
		rollupStore.Set([]byte(types.KEY_LAST_ELECTION_TIME), types.Int64ToBytes(blkTime))
		//设置
		electResult := types.QueryElectionResponse{
			ElectionTime:   uint64(blkTime),
			BlockHeight:    uint64(ctx.BlockHeight()),
			NodeStatusList: electList,
		}
		electData := k.cdc.MustMarshal(&electResult)
		//保存上一次的竞选信息
		if preElectData := rollupStore.Get([]byte(types.KEY_LAST_ELECTION_INFO)); preElectData != nil {
			rollupStore.Set([]byte(types.KEY_PREVIOUS_ELECTION_INFO), preElectData)
		}
		rollupStore.Set([]byte(types.KEY_LAST_ELECTION_INFO), electData)
		//
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EvtElection,
				sdk.NewAttribute("moduleName", types.MODULE_NAME),
				sdk.NewAttribute("result", string(res)),
			),
		)
		return nil
	} else {
		return nil
	}

}

func (k *Keeper) ProcUnstake(ctx sdk.Context) error {
	blkTime := ctx.BlockTime().Unix()
	//获取上一次选举的时间
	kvStore := ctx.KVStore(k.storeKey)
	rollupStore := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(k.rollAppID))

	if ElectVal := rollupStore.Get([]byte(types.KEY_LAST_ELECTION_TIME)); ElectVal != nil {
		lastElectTime := types.BytesToInt64(ElectVal)
		lastUnStakeTime := int64(0)
		if val := rollupStore.Get([]byte(types.KEY_LAST_UNSTAKE_TIME)); val != nil {
			lastUnStakeTime = types.BytesToInt64(val)
		}
		if lastUnStakeTime < lastElectTime { //这里才需要进行解质押的处理
			number, err := k.startUnstake(ctx)
			if err != nil {
				return err
			}
			rollupStore.Set([]byte(types.KEY_LAST_UNSTAKE_TIME), types.Int64ToBytes(blkTime))
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EvtProcUnStakeStatistics,
					sdk.NewAttribute("moduleName", types.MODULE_NAME),
					sdk.NewAttribute("unstake_number", strconv.Itoa(int(number))),
					sdk.NewAttribute("time", strconv.FormatInt(blkTime, 10)),
				),
			)
			return nil

		} else {
			return nil
		}

	} else { //如果还没开始过选举，则也不操作解质押
		return nil
	}

}

func (k *Keeper) startUnstake(ctx sdk.Context) (int32, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupAppStakeKeyPrefix(k.rollAppID))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close() // nolint: errcheck
	var totalUnstakeAddr [][]byte
	procNumber := int32(0)
	for ; iterator.Valid(); iterator.Next() {
		var val types.MsgStakeInfo
		if err := k.cdc.Unmarshal(iterator.Value(), &val); err != nil {
			return 0, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Unmarshal stakeInfo error.err = %s", err.Error()))
		}
		if val.ApplyUnStakeAmount > 0 {
			if val.ApplyUnStakeAmount > val.StakeAmount {
				return 0, errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("ApplyUnStakeAmount(%d) > StakeAmount(%d),addr = %s",
					val.ApplyUnStakeAmount, val.StakeAmount, string(iterator.Key())))
			} else {
				val.StakeAmount -= val.ApplyUnStakeAmount
				recvAddr, err := sdk.AccAddressFromBech32(string(iterator.Key()))
				if err != nil {
					return 0, errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("AccAddressFromBech32 error,err = %s,addr = %s",
						err.Error(), string(iterator.Key())))
				}
				unStakeCoin := sdk.NewCoin("umec", sdk.NewInt(int64(val.ApplyUnStakeAmount)))
				if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.MODULE_NAME, recvAddr, sdk.NewCoins(unStakeCoin)); err != nil {
					return 0, errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("unstake coin form module error,err = %s,addr = %s,amount = %d",
						err.Error(), string(iterator.Key()), val.ApplyUnStakeAmount))

				}
				unStakeAmount := val.ApplyUnStakeAmount
				val.ApplyUnStakeAmount = 0
				resData := k.cdc.MustMarshal(&val)
				store.Set(iterator.Key(), resData)

				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EvtProcUnStake,
						sdk.NewAttribute("moduleName", types.MODULE_NAME),
						sdk.NewAttribute("address", string(iterator.Key())),
						sdk.NewAttribute("amount", strconv.FormatUint(unStakeAmount, 10)),
					),
				)
				if 0 == val.StakeAmount { //如果全部赎回了，则将该质押信息进行删除
					totalUnstakeAddr = append(totalUnstakeAddr, iterator.Key())
				}
				procNumber++
			}
		} else {
			continue
		}
	}
	if len(totalUnstakeAddr) > 0 {
		for _, unStakeVal := range totalUnstakeAddr {
			store.Delete(unStakeVal)
		}
	}
	return procNumber, nil
}

func (k *Keeper) startElection(ctx sdk.Context, minStakeAmount uint64) ([]*types.ElectionNodeStatus, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupAppStakeKeyPrefix(k.rollAppID))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close() // nolint: errcheck

	var electorList types.ElectionsList
	for ; iterator.Valid(); iterator.Next() {
		var val types.MsgStakeInfo
		if err := k.cdc.Unmarshal(iterator.Value(), &val); err != nil {
			return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Unmarshal stakeInfo error.err = %s", err.Error()))
		}
		//这里进行 val.StakeAmount - val.ApplyUnStakeAmount的作用是为了让解质押对于竞选的影响的也能锁仓一个周期
		//假设在竞选前一天进行解质押，如果不相减的话，则就相当解质押对于竞选的影响几乎没有
		stakeAmount := val.StakeAmount - val.ApplyUnStakeAmount
		if stakeAmount < minStakeAmount { //不满足最小质押要求，则不能参加竞选
			continue
		}
		electInfo := types.ElectionInfo{
			StakeAmount: stakeAmount,
			Address:     string(iterator.Key()),
		}
		electorList = append(electorList, electInfo)
	}
	//进行降序排序
	sort.Sort(electorList)
	SeqNumber := k.GetSequencerNumber(ctx)
	BackNumber := k.GetBackupNumber(ctx)
	if uint32(electorList.Len()) < SeqNumber {
		return nil, errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("electorList len(%d) < sequencer number(%d)",
			electorList.Len(), SeqNumber))
	}
	totalNumber := SeqNumber + BackNumber
	var res []*types.ElectionNodeStatus

	for i := 0; i < electorList.Len(); i++ {
		index := uint32(i)
		nodeElect := &types.ElectionNodeStatus{
			Address:     electorList[i].Address,
			StakeAmount: electorList[i].StakeAmount,
		}
		if index < SeqNumber {
			nodeElect.Status = types.NodeSequencer
		} else if index < totalNumber {
			nodeElect.Status = types.NodeBackup
		} else {
			break
		}
		res = append(res, nodeElect)

	}
	return res, nil

}

// GetModuleAddress returns the staking module account address
func (k Keeper) GetModuleAddress() sdk.AccAddress {
	return sdk.AccAddress([]byte(types.MODULE_NAME))
}

func (k Keeper) GetElectionPeriod(ctx sdk.Context) (res uint32) {
	k.paramStore.Get(ctx, []byte(types.KeyElectionPeriod), &res)
	return
}

func (k Keeper) GetMinStakeAmount(ctx sdk.Context) (res uint64) {
	k.paramStore.Get(ctx, []byte(types.KeyMinStakeAmount), &res)
	return
}

func (k Keeper) GetSequencerNumber(ctx sdk.Context) (res uint32) {
	k.paramStore.Get(ctx, []byte(types.KeySequencerNumber), &res)
	return
}

func (k Keeper) GetBackupNumber(ctx sdk.Context) (res uint32) {
	k.paramStore.Get(ctx, []byte(types.KeyBackupNumber), &res)
	return
}

func (k Keeper) GetFirstElectionInterval(ctx sdk.Context) (res uint32) {
	k.paramStore.Get(ctx, []byte(types.KeyFirstElectInterval), &res)
	return
}

func (k Keeper) GetAllowApplyElectionTime(ctx sdk.Context) (res uint32) {
	k.paramStore.Get(ctx, []byte(types.KeyApplyElectionTime), &res)
	return
}

func (k Keeper) GetElectionInterimTime(ctx sdk.Context) (res uint32) {
	k.paramStore.Get(ctx, []byte(types.KeyElectionInterimTime), &res)
	return
}

func (k *Keeper) Punishment(ctx sdk.Context, address, rollappID string, rate uint32, amount uint64) error {
	punishmentAmount := uint64(0)
	kvStore := ctx.KVStore(k.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppStakeKeyPrefix(rollappID))
	data := store.Get([]byte(address))
	if data == nil {
		return errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("can not found stake info. addr = %s", address))
	}
	resp := &types.MsgStakeInfo{
		StakeAmount:        0,
		ApplyUnStakeAmount: 0,
	}
	k.cdc.MustUnmarshal(data, resp)

	if 0 == rate {
		punishmentAmount = amount
	} else {
		if rate > 100 {
			return errorsmod.Wrapf(types.ErrInputDataErr, fmt.Sprintf("input rate error. rate = %d", rate))
		} else {
			punishmentAmount = (resp.StakeAmount * uint64(rate)) / 100
		}
	}
	if punishmentAmount > 0 {
		accAddr, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return errorsmod.Wrapf(types.ErrInputDataErr, fmt.Sprintf(" AccAddressFromBech32 error. err = %s,addr = %s",
				err.Error(), address))
		}
		stakeCoin := sdk.NewCoin("umec", sdk.NewInt(int64(punishmentAmount)))
		//如果金额不够的话，SendCoinsFromAccountToModule这里就已经会判断处理了
		if err = k.bk.SendCoinsFromAccountToModule(ctx, accAddr, types.MODULE_NAME, sdk.NewCoins(stakeCoin)); err != nil {
			return errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("transfer  coin to module error.err = %s,addr = %s",
				err.Error(), address))
		}
		resp.StakeAmount -= punishmentAmount
		resData := k.cdc.MustMarshal(resp)
		store.Set(types.GetRollupAppStakeKeyPrefix(rollappID), resData)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EvtPunishment,
				sdk.NewAttribute("moduleName", types.MODULE_NAME),
				sdk.NewAttribute("address", address),
				sdk.NewAttribute("amount", strconv.FormatUint(punishmentAmount, 10)),
			),
		)
		return nil
	} else {
		return nil
	}
}

/*
对该地址进行资质重估，涉及的流程：
1、查询此时该地址的质押金额
2、如果质押金额小于最小的质押进行，则查看该地址是否属于选举的sequencer或者backup
3、如果是sequencer，则将该地址踢出sequencer，并且取一个backup作为sequencer，然后踢出选举的节点信息列表
4、如果是backup，踢出选举的节点信息列表
5、发出相应的状态变更事件通知
*/
func (k Keeper) RevaluateSequencer(ctx sdk.Context, address, rollappID string) error {
	kvStore := ctx.KVStore(k.storeKey)
	stakeStore := prefix.NewStore(kvStore, types.GetRollupAppStakeKeyPrefix(rollappID))
	data := stakeStore.Get([]byte(address))
	if data == nil {
		return errorsmod.Wrapf(types.ErrNotFound, fmt.Sprintf("can not found stake info. addr = %s", address))
	}
	stakeInfo := &types.MsgStakeInfo{
		StakeAmount:        0,
		ApplyUnStakeAmount: 0,
	}
	k.cdc.MustUnmarshal(data, stakeInfo)
	if stakeInfo.StakeAmount < k.GetMinStakeAmount(ctx)*types.MecPrecision {
		//如果小于最小质押金额，则踢出
		store := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(rollappID))
		electionData := store.Get([]byte(types.KEY_LAST_ELECTION_INFO))

		resp := &types.QueryElectionResponse{
			ElectionTime:   0,
			BlockHeight:    0,
			NodeStatusList: nil,
		}
		if nil == electionData {
			return errorsmod.Wrapf(types.ErrNotFound, fmt.Sprintf("can not found election info."))
		}
		if err := k.cdc.Unmarshal(electionData, resp); err != nil {
			return errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Unmarshal error. err = %s", err.Error()))
		}
		bIsProcSequencer := false
		bIsNeedRewriteData := false
		deleteKey := int(0)
		for key, val := range resp.NodeStatusList { //这么操作的前提是NodeStatusList是按照金额从大到小排序的
			if val.Address == address {
				beforeStatus := ""
				afterStatus := ""
				if types.NodeSequencer == val.Status {
					bIsProcSequencer = true
					beforeStatus = strconv.Itoa(int(types.NodeSequencer))
					afterStatus = strconv.Itoa(int(types.NodeNormal))
					val.Status = types.NodeNormal
					bIsNeedRewriteData = true
				} else if types.NodeBackup == val.Status {
					beforeStatus = strconv.Itoa(int(types.NodeBackup))
					afterStatus = strconv.Itoa(int(types.NodeNormal))
					val.Status = types.NodeNormal
					bIsNeedRewriteData = true
				}
				if bIsNeedRewriteData { //产生了状态变更事件
					deleteKey = key
					ctx.EventManager().EmitEvent(
						sdk.NewEvent(
							types.EvtSequencerChange,
							sdk.NewAttribute("moduleName", types.MODULE_NAME),
							sdk.NewAttribute("address", address),
							sdk.NewAttribute("beforeStatus", beforeStatus),
							sdk.NewAttribute("afterStatus", afterStatus),
						),
					)
				}
				if !bIsProcSequencer {
					//如果处理的不是sequencer的话，则可以跳出循环了,因为只有处理的是sequencer，才需要让备用节点顶上
					break
				}

			} else {
				if bIsProcSequencer { //如果对Sequencer进行了状态变更，这个实际则需要一个备选节点顶替
					if types.NodeBackup == val.Status {
						//这里选择第一个备选节点作为sequencer，然后调出循环
						val.Status = types.NodeSequencer
						ctx.EventManager().EmitEvent(
							sdk.NewEvent(
								types.EvtSequencerChange,
								sdk.NewAttribute("moduleName", types.MODULE_NAME),
								sdk.NewAttribute("address", address),
								sdk.NewAttribute("beforeStatus", strconv.Itoa(int(types.NodeBackup))),
								sdk.NewAttribute("afterStatus", strconv.Itoa(int(val.Status))),
							),
						)
						break
					}
				}
			}
		}
		if bIsNeedRewriteData {
			//删除质押金额小于最小的节点
			if len(resp.NodeStatusList) > 1 {
				resp.NodeStatusList = append(resp.NodeStatusList[:deleteKey], resp.NodeStatusList[deleteKey+1:]...)
			} else {
				resp.NodeStatusList = nil
			}
			resData := k.cdc.MustMarshal(resp)
			store.Set([]byte(types.KEY_LAST_ELECTION_INFO), resData)
		}

	}
	return nil
}
