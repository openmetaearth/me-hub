package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/x/rollup/types"
	"strconv"
)

func (t *rollupServer) StakeForSequencer(stakeCtx context.Context, req *types.MsgSeqStaking) (*types.MsgStakingResponse, error) {
	if req.RollappId != t.rollAppID {
		return nil, errorsmod.Wrapf(types.ErrRollappIDMismatch, fmt.Sprintf("rollupServer's rollappID = %s", t.rollAppID))
	}
	owner, err := sdk.AccAddressFromBech32(req.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("AccAddressFromBech32 error. err = %s", err.Error()))
	}
	ownerAddr := owner.String()
	ctx := sdk.UnwrapSDKContext(stakeCtx)
	if !t.Keeper.rk.RollappsEnabled(ctx) {
		return nil, types.ErrRollappDisable
	}
	found := t.Keeper.rk.IsRollappExist(ctx, req.RollappId)
	if !found {
		return nil, types.ErrRollappNotExist
	}
	/*
		if req.Version != rollapp.Version {
			return nil, fmt.Errorf("%s, rollappVersion = %d,reqVersion = %d", types.ErrRollappVersionMismatch.Error(),
				rollapp.Version, req.Version)
		}
	*/

	if !t.isAllowStake(ctx, ctx.BlockTime().Unix()) {
		return nil, types.ErrStakeTimeoutLimit
	}

	if req.Amount < 1 {
		return nil, types.ErrInputDataErr
	}
	store := prefix.NewStore(ctx.KVStore(t.Keeper.storeKey), types.GetRollupAppStakeKeyPrefix(t.rollAppID))
	stakeInfo := &types.MsgStakeInfo{
		StakeAmount:        0,
		ApplyUnStakeAmount: 0,
	}
	if val := store.Get([]byte(ownerAddr)); val != nil {
		if err = t.Keeper.cdc.Unmarshal(val, stakeInfo); err != nil {
			return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("err = Unmarshal msgStakeInfo error.err = %s", err.Error()))
		}
	}
	stakeInfo.StakeAmount += req.Amount

	stakeCoin := sdk.NewCoin("UMEC", sdk.NewInt(int64(req.Amount)))
	//如果金额不够的话，SendCoinsFromAccountToModule这里就已经会判断处理了
	if err = t.bk.SendCoinsFromAccountToModule(ctx, owner, types.MODULE_NAME, sdk.NewCoins(stakeCoin)); err != nil {
		return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("stake coin to module error.err = %s", err.Error()))
	}
	//verify stake balance
	qryRes := t.bk.GetBalance(ctx, owner, "UMEC")
	if !qryRes.Amount.Equal(sdk.NewInt(int64(stakeInfo.StakeAmount))) {
		return nil, errorsmod.Wrapf(types.ErrStakeDataErr, fmt.Sprintf("stake amount mismatch.statics's ammount = %s, module's balance = %s",
			strconv.FormatUint(stakeInfo.StakeAmount, 10), qryRes.Amount.String()))
	}

	stakeVal, err := t.Keeper.cdc.Marshal(stakeInfo)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("err = Marshal msgStakeInfo error.err = %s", err.Error()))
	}
	store.Set([]byte(ownerAddr), stakeVal)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EvtStaking,
			sdk.NewAttribute("moduleName", types.MODULE_NAME),
			sdk.NewAttribute("delegator", ownerAddr),
			sdk.NewAttribute("amount", strconv.FormatUint(req.Amount, 10)),
		),
	)
	return &types.MsgStakingResponse{}, nil

}
func (t rollupServer) isAllowStake(sdkCtx sdk.Context, curTime int64) bool {
	store := prefix.NewStore(sdkCtx.KVStore(t.Keeper.storeKey), types.GetRollupAppKeyPrefix(t.rollAppID))
	if val := store.Get([]byte(types.KEY_LAST_ELECTION_TIME)); val != nil {
		lastElectTime := types.BytesToInt64(val)
		stakeEndTime := lastElectTime + int64(t.Keeper.GetAllowApplyElectionTime(sdkCtx))*types.HourSeconds
		if curTime > stakeEndTime {
			return false
		} else {
			return true
		}

	} else { //如果还没有开始第一次选举，则不受这个限制
		return true
	}

}

func (t *rollupServer) UnStake(stakeCtx context.Context, req *types.MsgSeqUnStaking) (*types.MsgUnStakingResponse, error) {
	if req.RollappId != t.rollAppID {
		return nil, errorsmod.Wrapf(types.ErrRollappIDMismatch, fmt.Sprintf("rollupServer's rollappID = %s", t.rollAppID))
	}
	owner, err := sdk.AccAddressFromBech32(req.Creator)

	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("AccAddressFromBech32 error. err = %s", err.Error()))
	}
	ownerAddr := owner.String()
	ctx := sdk.UnwrapSDKContext(stakeCtx)
	kvStore := ctx.KVStore(t.Keeper.storeKey)
	if !t.Keeper.rk.RollappsEnabled(ctx) {
		return nil, types.ErrRollappDisable
	}
	found := t.Keeper.rk.IsRollappExist(ctx, req.RollappId)
	if !found {
		return nil, types.ErrRollappNotExist
	}
	/*
		if req.Version != rollapp.Version {
			return nil, fmt.Errorf("%s, rollappVersion = %d,reqVersion = %d", types.ErrRollappVersionMismatch.Error(),
				rollapp.Version, req.Version)
		}

	*/
	if req.Amount < 1 {
		return nil, types.ErrInputDataErr
	}

	store := prefix.NewStore(kvStore, types.GetRollupAppStakeKeyPrefix(t.rollAppID))
	stakeInfo := &types.MsgStakeInfo{
		StakeAmount:        0,
		ApplyUnStakeAmount: 0,
	}
	if val := store.Get([]byte(ownerAddr)); val != nil {
		if err = t.Keeper.cdc.Unmarshal(val, stakeInfo); err != nil {
			return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("err = Unmarshal msgStakeInfo error.err = %s", err.Error()))
		}
	}

	amount := req.Amount
	if amount > stakeInfo.StakeAmount {
		return nil, types.ErrInsufficientBalance
	}
	//这里一个周期内指允许一次
	if stakeInfo.ApplyUnStakeAmount > 0 {
		return nil, types.ErrUnStakeLimit
	}

	//获取上一次选举的时间
	rollupStore := prefix.NewStore(kvStore, types.GetRollupAppKeyPrefix(t.rollAppID))
	electTime := int64(0)
	if electBlkVal := rollupStore.Get([]byte(types.KEY_LAST_ELECTION_TIME)); electBlkVal != nil {
		electTime = types.BytesToInt64(electBlkVal)

	} else { //找不到lastElectBlock的话，则表示还没竞选过,此时依然不允许取回
		return nil, errorsmod.Wrapf(types.ErrUnStakeLimit, "please wait for election start")
	}

	/*
		electTimeVal := prefix.NewStore(kvStore, []byte(types.RollupBlockTimePrefix)).Get(types.Int64ToBytes(lastElectBlkHeight))
		if nil == electTimeVal {
			//这里有一种可能会出现这种情况，当在BeginBlock时，发现选举应该在这个区块完成，此时会先将改区块高度记录下来，但是此时该高度的区块还没共识完成，
			//也就是说解质押的操作和选举的操作处于同一个区块中，所以此时应该禁止解质押
			return nil, fmt.Errorf("%s,block is in consensusblkHeight = %d", types.ErrUnStakeLimit.Error(), lastElectBlkHeight)
		}
	*/

	curTime := ctx.BlockTime().Unix()
	if curTime > electTime {
		stakeInfo.ApplyUnStakeAmount += amount
		if stateVal, err := t.Keeper.cdc.Marshal(stakeInfo); err != nil {
			return nil, errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Marshal(stakeInfo) error,err = %s", err.Error()))
		} else {
			store.Set([]byte(ownerAddr), stateVal)
		}

	} else {
		return nil, errorsmod.Wrapf(types.ErrUnStakeLimit, "stake time is not enough")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EvtUnStaking,
			sdk.NewAttribute("moduleName", types.MODULE_NAME),
			sdk.NewAttribute("delegator", ownerAddr),
			sdk.NewAttribute("amount", strconv.FormatUint(req.Amount, 10)),
		),
	)
	return &types.MsgUnStakingResponse{}, nil

}

func (t *rollupServer) RegisterRollappID(ctx context.Context, req *types.RegisterRollappIDRequest) (*types.RegisterRollappIDResponse, error) {
	if t.rollAppID == "" {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		t.rollAppID = req.RollappID
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EvtRegisterRollappID,
				sdk.NewAttribute("moduleName", types.MODULE_NAME),
				sdk.NewAttribute("rollappID", req.RollappID),
			),
		)
		return &types.RegisterRollappIDResponse{}, nil
	}
	return nil, types.ErrRollappIdRegisterRepeated
}
