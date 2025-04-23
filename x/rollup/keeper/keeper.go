package keeper

import (
	"encoding/hex"
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
	"github.com/st-chain/me-hub/app/params"
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
	mapBlackList      map[string]struct{}
	mapRollappInfoMng map[string]*types.RollappInitExtVal
	mapPunishInfo     map[string]map[string]uint64 //rollappID => address => punishmentAmount
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
		bk:         bKeeper,
		rk:         rKeeper,
		dk:         dKeeper,
	}
}

func (k *Keeper) InitRollappID(ctx sdk.Context) error {

	k.mapRollappInfoMng = make(map[string]*types.RollappInitExtVal)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.RollupKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyRollappInitInfoPrefix))
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		rollappID := types.GetRollappIdFromInitInfoKey(iterator.Key())
		if nil == rollappID {
			return errorsmod.Wrapf(types.ErrParserDataErr,
				fmt.Sprintf("parser rollappID from rollappInitInfoKey failed.key = %s", iterator.Key()))
		}
		val := &types.RollappInitExtVal{
			IdInDA:              nil,
			FirstElectBlkHeight: 0,
		}
		if err := json.Unmarshal(iterator.Value(), val); err != nil {
			return errorsmod.Wrapf(types.ErrParserDataErr,
				fmt.Sprintf("Unmarshal rollappInitInfo error,err = %s,key = %s", err.Error(), iterator.Key()))
		}
		k.mapRollappInfoMng[string(rollappID)] = val
		ctx.Logger().Info(fmt.Sprintf("praser rollapp init infodata,rollappID = %s,initInfo = %+v", rollappID, val))

	}
	return nil

}

// Logger returns a logger instance for the incentives module.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.MODULE_NAME))
}

func (k *Keeper) ProcElection(ctx sdk.Context) error {
	for key, val := range k.mapRollappInfoMng {
		if ctx.BlockHeight() < int64(val.FirstElectBlkHeight) {
			continue
		}
		rollAppID := key

		blkTime := ctx.BlockTime().Unix()
		//获取上一次选举的时间
		kvStore := ctx.KVStore(k.storeKey)
		rollupStore := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(rollAppID))
		lastElectTime := int64(0)
		bIsNeedElect := false
		if electTimeVal := rollupStore.Get([]byte(types.KEY_LAST_ELECTION_TIME)); electTimeVal != nil {
			lastElectTime = types.BytesToInt64(electTimeVal)
			timeInterval := blkTime - lastElectTime
			electionInterval := int64(k.GetElectionPeriod(ctx)) * types.MinuteSeconds
			if timeInterval >= electionInterval {
				bIsNeedElect = true
			}

		} else { //找不到lastElectTime的话，则表示还没竞选过,但是由于当前的区块已经>=FirstElectBlkHeight, 所以要立即开始竞选
			bIsNeedElect = true
		}
		if bIsNeedElect { //开始竞选
			ctx.Logger().Info("ready to elect sequencer")
			electList, err := k.startElection(ctx, rollAppID, k.GetMinStakeAmount(ctx)*types.MecPrecision)
			if err != nil {
				return err
			}
			strRes := ""
			if (nil != electList) && (len(electList) > 0) { //只有选举后有节点才成为有效选举，然后是无效的选举的话，则不更新选举信息
				var res []byte
				//if res, err = json.Marshal(electList); err != nil {
				//	return errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("Marshal(electList) error.err = %s", err.Error()))
				//}

				rollupStore.Set([]byte(types.KEY_LAST_ELECTION_TIME), types.Int64ToBytes(blkTime))
				//设置
				electResult := types.ElectionResult{
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
				res, err = electResult.Marshal()
				if err != nil {
					return errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("Marshal(electResult) error.err = %s", err.Error()))
				}
				strRes = hex.EncodeToString(res)
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EvtElection,
						sdk.NewAttribute("moduleName", types.MODULE_NAME),
						sdk.NewAttribute("rollappID", rollAppID),
						sdk.NewAttribute("result", strRes),
					),
				)
			} else {
				ctx.Logger().Error(fmt.Sprintf("can not elect sequencer node.blkHeight = %d,blkTime = %d",
					ctx.BlockHeight(), blkTime))
			}

		}
	}
	return nil

}

func (k *Keeper) ProcUnstake(ctx sdk.Context) error {

	blkTime := ctx.BlockTime().Unix()
	//获取上一次选举的时间
	kvStore := ctx.KVStore(k.storeKey)
	for key, _ := range k.mapRollappInfoMng {
		rollAppID := key
		rollupStore := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(rollAppID))

		if ElectVal := rollupStore.Get([]byte(types.KEY_LAST_ELECTION_TIME)); ElectVal != nil {
			lastElectTime := types.BytesToInt64(ElectVal)
			lastUnStakeTime := int64(0)
			if val := rollupStore.Get([]byte(types.KEY_LAST_UNSTAKE_TIME)); val != nil {
				lastUnStakeTime = types.BytesToInt64(val)
			}
			if lastUnStakeTime < lastElectTime { //这里才需要进行解质押的处理
				//由于选举存在过渡时间，所以这里解质押也需要增加过渡时间，180秒是额外附加
				interimTime := k.GetElectionInterimTime(ctx) + 120
				if blkTime < (lastElectTime + int64(interimTime)) {
					continue
				}

				number, err := k.startUnstake(ctx, rollAppID)
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
				ctx.Logger().Info("complete proc unStake")

			}
		} else { //如果还没开始过选举，则也不操作解质押
			continue
		}
	}
	return nil

}

func (k *Keeper) startUnstake(ctx sdk.Context, rollappId string) (int32, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupAppStakeKeyPrefix(rollappId))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close() // nolint: errcheck
	var totalUnstakeAddr [][]byte
	procNumber := int32(0)
	ctx.Logger().Info("start proc unStake")
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

				unStakeCoin := sdk.NewCoin(params.BaseDenom, sdk.NewInt(int64(val.ApplyUnStakeAmount)))
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
						sdk.NewAttribute("rollappID", rollappId),
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

func (k *Keeper) startElection(ctx sdk.Context, rollAppID string, minStakeAmount uint64) ([]*types.ElectionNodeStatus, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupAppStakeKeyPrefix(rollAppID))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close() // nolint: errcheck

	var electorList types.ElectionsList
	for ; iterator.Valid(); iterator.Next() {
		var val types.MsgStakeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		//这里进行 val.StakeAmount - val.ApplyUnStakeAmount的作用是为了让解质押对于竞选的影响的也能锁仓一个周期
		//假设在竞选前一天进行解质押，如果不相减的话，则就相当解质押对于竞选 的影响几乎没有
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
	//经过讨论，即使满足资格的质押人数已经不足，但是仍然作为有效的选举，
	//所以这里不在判断
	/*
		if uint32(electorList.Len()) < SeqNumber {
			return nil, errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("electorList len(%d) < sequencer number(%d)",
				electorList.Len(), SeqNumber))
		}*/
	totalNumber := SeqNumber + BackNumber
	var res []*types.ElectionNodeStatus

	for i := 0; i < electorList.Len(); i++ {
		index := uint32(i)
		nodeElect := &types.ElectionNodeStatus{
			Address:         electorList[i].Address,
			StakeAmount:     electorList[i].StakeAmount,
			BondNodeAddress: k.getDelegatorBondNodeAddr(ctx, rollAppID, electorList[i].Address),
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

func (t *Keeper) RegisterRollappInitInfo(ctx sdk.Context, rollappID string, FirstElectBlkHeight uint64, IdInDa []byte) error {
	if _, ok := t.mapRollappInfoMng[rollappID]; ok {
		return errorsmod.Wrapf(types.ErrRollappIdRegisterRepeated, "")
	}

	initExtInfo := &types.RollappInitExtVal{
		IdInDA:              IdInDa,
		FirstElectBlkHeight: FirstElectBlkHeight,
	}
	val, err := json.Marshal(initExtInfo)
	if err != nil {
		return errorsmod.Wrapf(types.ErrParserDataErr, "Marshal(initExtInfo) error.err = %s", err.Error())
	}
	kvStore := ctx.KVStore(t.storeKey)
	store := prefix.NewStore(kvStore, []byte(types.RollupKeyPrefix))
	store.Set(types.GetRollupAppInitInfKey(rollappID), val)
	ctx.Logger().Info(fmt.Sprintf("RegisterRollappID = %s", rollappID))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EvtRegisterRollappID,
			sdk.NewAttribute("moduleName", types.MODULE_NAME),
			sdk.NewAttribute("rollappID", rollappID),
			sdk.NewAttribute("rollappInitInfo", string(val)),
		),
	)
	return nil

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

/*
	func (k Keeper) GetFirstElectionInterval(ctx sdk.Context) (res uint32) {
		k.paramStore.Get(ctx, []byte(types.KeyFirstElectInterval), &res)
		return
	}
*/
func (k Keeper) GetAllowApplyElectionTime(ctx sdk.Context) (res uint32) {
	k.paramStore.Get(ctx, []byte(types.KeyApplyElectionTime), &res)
	return
}

func (k Keeper) GetElectionInterimTime(ctx sdk.Context) (res uint32) {
	k.paramStore.Get(ctx, []byte(types.KeyElectionInterimTime), &res)
	return
}

func (k Keeper) GetDaFraudChallengeStake(ctx sdk.Context) (res uint32) {
	k.paramStore.Get(ctx, []byte(types.KeyDaFraudChallengeStake), &res)
	return
}

func (k *Keeper) Punishment(ctx sdk.Context, address, rollappID string, rate uint32, amount uint64) (uint64, error) {
	punishmentAmount := uint64(0)
	kvStore := ctx.KVStore(k.storeKey)
	store := prefix.NewStore(kvStore, types.GetRollupAppStakeKeyPrefix(rollappID))
	data := store.Get([]byte(address))
	if data == nil {
		return 0, errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("can not found stake info in Punishment. addr = %s", address))
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
			return 0, errorsmod.Wrapf(types.ErrInputDataErr, fmt.Sprintf("input rate error. rate = %d", rate))
		} else {
			punishmentAmount = (resp.StakeAmount * uint64(rate)) / 100
		}
	}
	if punishmentAmount > 0 {
		if punishmentAmount > resp.StakeAmount {
			return 0, errorsmod.Wrapf(types.ErrInputDataErr, fmt.Sprintf("punishmentAmount > totalStakeAmount."+
				"punishmentAmount = %d,totalStakeAmount = %d", punishmentAmount, resp.StakeAmount))
		}
		accAddr, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return 0, errorsmod.Wrapf(types.ErrInputDataErr, fmt.Sprintf(" AccAddressFromBech32 error in Punishment. err = %s,addr = %s",
				err.Error(), address))
		}
		stakeCoin := sdk.NewCoin(params.BaseDenom, sdk.NewInt(int64(punishmentAmount)))
		//如果金额不够的话，SendCoinsFromAccountToModule这里就已经会判断处理了
		if err = k.bk.SendCoinsFromAccountToModule(ctx, accAddr, types.MODULE_NAME, sdk.NewCoins(stakeCoin)); err != nil {
			return 0, errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("punishment transfer coin to module error.err = %s,addr = %s",
				err.Error(), address))
		}
		resp.StakeAmount -= punishmentAmount
		if resp.StakeAmount < 1 { //在惩罚措施中，如果质押的金额已经为0，则将ApplyUnStakeAmount也设置为0
			resp.ApplyUnStakeAmount = 0
		}
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
		return punishmentAmount, nil
	} else {
		return 0, nil
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
	if data == nil { //如果没找到质押信息，则直接返回
		ctx.Logger().Info(fmt.Sprintf("can not found stake info in RevaluateSequencer. addr = %s", address))
		//return errorsmod.Wrapf(types.ErrNotFound, f)
		return nil
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

		resp := &types.ElectionResult{
			ElectionTime:   0,
			BlockHeight:    0,
			NodeStatusList: nil,
		}
		if nil == electionData { //如果没找到选举信息，则返回，但是这里不应该没有选举信息，所以打印错误日志
			ctx.Logger().Error("can not found election info in RevaluateSequencer")
			return nil
		}
		k.cdc.MustUnmarshal(electionData, resp)
		bIsProcSequencer := false
		bIsNeedRewriteData := false
		deleteKey := int(0)
		nodeStatusModifyList := new(types.NodeStatusModifyList)
		for key, val := range resp.NodeStatusList { //这么操作的前提是NodeStatusList是按照金额从大到小排序的
			if val.Address == address {
				beforeStatus := types.NodeNormal
				afterStatus := types.NodeNormal
				if types.NodeSequencer == val.Status {
					bIsProcSequencer = true
					beforeStatus = types.NodeSequencer
					afterStatus = types.NodeNormal
					val.Status = types.NodeNormal
					bIsNeedRewriteData = true
				} else if types.NodeBackup == val.Status {
					beforeStatus = types.NodeBackup
					afterStatus = types.NodeNormal
					val.Status = types.NodeNormal
					bIsNeedRewriteData = true
				}
				if bIsNeedRewriteData { //产生了状态变更事件
					deleteKey = key
					//bondNodeAddress := ""
					bondAddrBytes := k.getDelegatorBondNodeAddr(ctx, rollappID, address)
					if bondAddrBytes == nil {
						panic(fmt.Errorf("can not found BondNodeAddress from address.stakerAddress = %s", address))
					}
					//bondNodeAddress = hex.EncodeToString(bondAddrBytes)
					nodeStatusModifyList.NodeStatusList = append(nodeStatusModifyList.NodeStatusList, &types.NodeStatusModify{
						StakerAddress:   val.Address,
						BeforeStatus:    beforeStatus,
						AfterStatus:     afterStatus,
						BondNodeAddress: bondAddrBytes,
					})

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
						bondAddrBytes := k.getDelegatorBondNodeAddr(ctx, rollappID, val.Address)
						if nil == bondAddrBytes {
							panic(fmt.Errorf("can not found BondNodeAddress from address.stakerAddress = %s", address))
						}
						//bondNodeAddress = hex.EncodeToString(bondAddrBytes)
						nodeStatusModifyList.NodeStatusList = append(nodeStatusModifyList.NodeStatusList, &types.NodeStatusModify{
							StakerAddress:   address,
							BeforeStatus:    types.NodeBackup,
							AfterStatus:     val.Status,
							BondNodeAddress: bondAddrBytes,
						})
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
			nodeListsData, err := k.cdc.Marshal(nodeStatusModifyList)
			if err != nil {
				return errorsmod.Wrapf(types.ErrProcessErr, fmt.Sprintf("Marshal(nodeStatusModifyList) error.err = %s,"+
					"nodeStatusModifyList = %s", err.Error(), nodeStatusModifyList.String()))
			}
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EvtSequencerChange,
					sdk.NewAttribute("moduleName", types.MODULE_NAME),
					sdk.NewAttribute(types.EvtAttrRollappID, rollappID),
					sdk.NewAttribute(types.EvtAttrBlockHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
					sdk.NewAttribute(types.EvtAttrBlockTime, strconv.FormatInt(ctx.BlockTime().Unix(), 10)),
					sdk.NewAttribute(types.EvtAttrNodeStatusModifyList, hex.EncodeToString(nodeListsData)),
				),
			)
		}

	}
	return nil
}

// 将质押者的地址和L2层的节点地址进行绑定，这里运行绑定的条件为:
// 1、该节点地址未被绑定过
// 2、该节点已经绑定的地址，但是存在满足条件就可重新绑定
//
//	1）该地址既不是Sequencer也不是备用节点
//	2）该地址目前质押金额的已经小于最小质押金额
func (k Keeper) bondNodeAddr(ctx sdk.Context, rollappID, creator string, bondAddress []byte, amount uint64) error {
	if nil == bondAddress {
		return errorsmod.Wrapf(types.ErrInputDataErr, " bondAddress can not be nil")
	}
	if len(bondAddress) != 20 {
		return errorsmod.Wrapf(types.ErrInputDataErr, fmt.Sprintf(" bondAddress length error.len = %d,must be 20", len(bondAddress)))
	}
	kvStore := ctx.KVStore(k.storeKey)
	bondStore := prefix.NewStore(kvStore, types.GetStakeBondNodeAddrPrefix(rollappID))
	data := bondStore.Get(bondAddress)
	var orgDelegator []byte
	if data != nil {
		delegatorAddr := string(data)
		if delegatorAddr == creator {
			return nil
		}
		orgDelegator = data
		//查看是否可以覆盖绑定
		stakeData, err := k.queryStakeData(ctx, rollappID, delegatorAddr)
		if err != nil {
			return err
		}
		if amount <= stakeData.StakeAmount {
			return errorsmod.Wrapf(types.ErrNotAllowBondNodeAddr, ",input stake amount is smaller")
		}
		//如果此时质押的金额小于最小金额，那么肯定就不是Sequencer或者备选节点
		if stakeData.StakeAmount >= k.GetMinStakeAmount(ctx)*types.MecPrecision {
			return errorsmod.Wrapf(types.ErrNotAllowBondNodeAddr, fmt.Sprintf("orgDelegator stake amount = %d", stakeData.StakeAmount))
		}

	}

	bondStore.Set(bondAddress, []byte(creator))
	//设定Delegator和NodeAddr的映射关系
	delegatorStore := prefix.NewStore(kvStore, types.GetDelegatorStakeNodePrefix(rollappID))
	delegatorStore.Set([]byte(creator), bondAddress)
	if orgDelegator != nil { //如果该节点地址之前存在绑定关系，则解除该关系
		delegatorStore.Delete(orgDelegator)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EvtBondRollappNodeAddress,
			sdk.NewAttribute("moduleName", types.MODULE_NAME),
			sdk.NewAttribute("delegator", creator),
			sdk.NewAttribute("nodeKeyAddr", hex.EncodeToString(bondAddress)),
		))

	return nil
}

// 获取委托者所绑定的节点的地址
func (k Keeper) getDelegatorBondNodeAddr(ctx sdk.Context, rollappID, delegatorAddress string) []byte {
	kvStore := ctx.KVStore(k.storeKey)
	delegatorStore := prefix.NewStore(kvStore, types.GetDelegatorStakeNodePrefix(rollappID))
	return delegatorStore.Get([]byte(delegatorAddress))
}

// 获取绑定节点所关联的委托者
func (k Keeper) GetBondNodeDelegator(ctx sdk.Context, rollappID string, bondAddress []byte) []byte {
	kvStore := ctx.KVStore(k.storeKey)
	bondStore := prefix.NewStore(kvStore, []byte(types.GetStakeBondNodeAddrPrefix(rollappID)))
	return bondStore.Get(bondAddress)
}

/*
func (k Keeper) resetStakeInfo(ctx sdk.Context, rollappID, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupAppStakeKeyPrefix(rollappID))
	stakeInfo := &types.MsgStakeInfo{
		StakeAmount:        0,
		ApplyUnStakeAmount: 0,
	}
	stakeVal := k.cdc.MustMarshal(stakeInfo)
	store.Set([]byte(address), stakeVal)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EvtRestStakeInfo,
			sdk.NewAttribute("moduleName", types.MODULE_NAME),
			sdk.NewAttribute("rollappID", rollappID),
			sdk.NewAttribute("address", address),
		),
	)

}

*/
