package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/rollup/types"
)

func (k *Keeper) InitBlackList(ctx sdk.Context) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.RollupKeyPrefix))
	data := store.Get([]byte(types.KeyRollupBlackList))
	if data != nil {
		k.mapBlackList = nil
		err := json.Unmarshal(data, &k.mapBlackList)
		if err != nil {
			return fmt.Errorf("Unmarshal mapBlackList data error. err = %s", err.Error())
		}
	} else {
		k.mapBlackList = make(map[string]struct{}) //每个区块开始都要重新初始化
	}
	return nil

}

func (k Keeper) IsInBlackList(addr string) bool {
	if _, ok := k.mapBlackList[addr]; ok {
		return true
	}
	return false
}

// 黑名单是全局的，跨rollapp的
func (k *Keeper) AddToBlackList(ctx sdk.Context, addr string) error {
	if _, ok := k.mapBlackList[addr]; ok {
		return nil
	}
	k.mapBlackList[addr] = struct{}{}
	data, err := json.Marshal(k.mapBlackList)
	if err != nil {
		delete(k.mapBlackList, addr) //内存还原
		return errorsmod.Wrapf(types.ErrParserDataErr, fmt.Sprintf("Marshal mapBlackList error.err = %s", err.Error()))
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.RollupKeyPrefix))
	store.Set([]byte(types.KeyRollupBlackList), data)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EvtAddBlackList,
			sdk.NewAttribute("moduleName", types.MODULE_NAME),
			sdk.NewAttribute("address", addr),
		),
	)
	return nil
}
