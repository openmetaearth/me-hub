package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type KeeperWrapper struct {
	*stakingkeeper.Keeper
}

// NewKeeper returns a new BaseKeeperWrapper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak banktypes.AccountKeeper,
	bk stakingtypes.BankKeeper,
	authority string,
) *KeeperWrapper {
	return &KeeperWrapper{
		Keeper: stakingkeeper.NewKeeper(cdc, storeKey, ak, bk, authority),
	}
}
