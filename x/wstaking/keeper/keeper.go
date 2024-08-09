package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

type Keeper struct {
	*stakingkeeper.Keeper
	DaoKeeper types.DaoKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak banktypes.AccountKeeper,
	bk stakingtypes.BankKeeper,
	dk types.DaoKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		Keeper:    stakingkeeper.NewKeeper(cdc, storeKey, ak, bk, authority),
		DaoKeeper: dk,
	}
}
