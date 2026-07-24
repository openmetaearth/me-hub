package keeper

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/nft"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	nftkeeper.Keeper
	cdc codec.BinaryCodec
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak nft.AccountKeeper,
	bk nft.BankKeeper,
) *Keeper {
	return &Keeper{
		Keeper: nftkeeper.NewKeeper(storeKey, cdc, ak, bk),
		cdc:    cdc,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/wnft")
}
