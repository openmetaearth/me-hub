package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
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
