package keeper

import (
	"encoding/json"
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/dymension/v3/x/rollup/types"
	"sort"
	"strconv"
)

// Keeper struct
type Keeper struct {
	storeKey         storetypes.StoreKey
	cdc              codec.BinaryCodec
	bk               types.BankKeeper
	rk               types.RollappKeeper
	paramStore       paramtypes.Subspace
	lastElectionTime uint64
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(storeKey storetypes.StoreKey, cdc codec.BinaryCodec, paramSpace paramtypes.Subspace) *Keeper {
	//if !paramSpace.HasKeyTable() {
	//	paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	//}
	return &Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramStore: paramSpace,
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
	rollupStore := prefix.NewStore(kvStore, []byte(types.RollupKeyPrefix))
	lastElectTime := int64(0)
	bIsNeedElect := false
	if electTimeVal := rollupStore.Get([]byte(types.KEY_LAST_ELECTION_TIME)); electTimeVal != nil {
		lastElectTime = types.BytesToInt64(electTimeVal)
		timeInterval := blkTime - lastElectTime
		electionInterval := int64(k.GetElectionPeriod(ctx)) * types.DaySeconds
		if timeInterval >= electionInterval {
			bIsNeedElect = true
		}

	} else { //找不到lastElectTime的话，则表示还没竞选过
		if 1 == ctx.BlockHeight() { //如果是第一个数据区块的话，则计算首次竞选的时间并保存
			firstElectTime := blkTime + int64(k.GetFirstElectionInterval(ctx))*types.HourSeconds
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
				return fmt.Errorf("%s,can not get firstElectTime. blkHeight = %d",
					types.ErrProcessErr.Error(), ctx.BlockHeight())
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
			return fmt.Errorf("%s,Marshal(electList) error.err = %s", types.ErrProcessErr.Error(), err.Error())
		}
		rollupStore.Set([]byte(types.KEY_LAST_ELECTION_TIME), types.Int64ToBytes(blkTime))
		//设置
		electResult := types.QueryElectionResponse{
			ElectionTime:   uint64(blkTime),
			BlockHeight:    uint64(ctx.BlockHeight()),
			NodeStatusList: electList,
		}
		electData := k.cdc.MustMarshal(&electResult)
		rollupStore.Set([]byte(types.KEY_LAST_ELECTION_INFO), electData)
		//
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EvtElection,
				sdk.NewAttribute("moduleName", types.MODULE_NAME),
				sdk.NewAttribute("Result", string(res)),
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
	rollupStore := prefix.NewStore(kvStore, []byte(types.RollupKeyPrefix))

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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.RollupStakeKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close() // nolint: errcheck
	var totalUnstakeAddr [][]byte
	procNumber := int32(0)
	for ; iterator.Valid(); iterator.Next() {
		var val types.MsgStakeInfo
		if err := k.cdc.Unmarshal(iterator.Value(), &val); err != nil {
			return 0, fmt.Errorf("%s,Unmarshal stakeInfo error.err = %s", types.ErrParserDataErr.Error(), err.Error())
		}
		if val.ApplyUnStakeAmount > 0 {
			if val.ApplyUnStakeAmount > val.StakeAmount {
				return 0, fmt.Errorf("%s,ApplyUnStakeAmount(%d) > StakeAmount(%d),addr = %s",
					types.ErrProcessErr.Error(), val.ApplyUnStakeAmount, val.StakeAmount, string(iterator.Key()))
			} else {
				val.StakeAmount -= val.ApplyUnStakeAmount
				recvAddr, err := sdk.AccAddressFromBech32(string(iterator.Key()))
				if err != nil {
					return 0, fmt.Errorf("%s,AccAddressFromBech32 error,err = %s,addr = %s",
						types.ErrProcessErr.Error(), err.Error(), string(iterator.Key()))
				}
				unStakeCoin := sdk.NewCoin("UMEC", sdk.NewInt(int64(val.ApplyUnStakeAmount)))
				if err = k.bk.SendCoinsFromModuleToAccount(ctx, types.MODULE_NAME, recvAddr, sdk.NewCoins(unStakeCoin)); err != nil {
					return 0, fmt.Errorf("%s,unstake coin form module error,err = %s,addr = %s,amount = %d",
						types.ErrProcessErr.Error(), err.Error(), string(iterator.Key()), val.ApplyUnStakeAmount)

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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.RollupStakeKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close() // nolint: errcheck

	var electorList types.ElectionsList
	for ; iterator.Valid(); iterator.Next() {
		var val types.MsgStakeInfo
		if err := k.cdc.Unmarshal(iterator.Value(), &val); err != nil {
			return nil, fmt.Errorf("%s,Unmarshal stakeInfo error.err = %s", types.ErrParserDataErr.Error(), err.Error())
		}
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
		return nil, fmt.Errorf("%s,electorList len(%d) < sequencer number(%d)",
			types.ErrProcessErr.Error(), electorList.Len(), SeqNumber)
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
