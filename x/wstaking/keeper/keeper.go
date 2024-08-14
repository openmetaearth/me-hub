package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/st-chain/me-hub/x/wbank/keeper"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

type Keeper struct {
	*stakingkeeper.Keeper
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	AuthKeeper banktypes.AccountKeeper
	BankKeeper keeper.BaseKeeperWrapper
	DaoKeeper  types.DaoKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak banktypes.AccountKeeper,
	bk keeper.BaseKeeperWrapper,
	dk types.DaoKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		Keeper:     stakingkeeper.NewKeeper(cdc, storeKey, ak, bk, authority),
		cdc:        cdc,
		storeKey:   storeKey,
		AuthKeeper: ak,
		BankKeeper: bk,
		DaoKeeper:  dk,
	}
}
