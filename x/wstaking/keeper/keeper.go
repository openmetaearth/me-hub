package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wmintkeeper "github.com/st-chain/me-hub/x/wmint/keeper"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

type Keeper struct {
	*stakingkeeper.Keeper
	cdc         codec.BinaryCodec
	storeKey    storetypes.StoreKey
	AuthKeeper  banktypes.AccountKeeper
	BankKeeper  types.BankKeeper
	DaoKeeper   types.DaoKeeper
	WMintKeeper wmintkeeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak banktypes.AccountKeeper,
	bk types.BankKeeper,
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

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}
