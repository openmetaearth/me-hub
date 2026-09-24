package keeper

import (
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/x/nft"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	nftkeeper.Keeper
	cdc          codec.BinaryCodec
	storeService corestoretypes.KVStoreService
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService corestoretypes.KVStoreService,
	ak nft.AccountKeeper,
	bk nft.BankKeeper,
) *Keeper {
	return &Keeper{
		Keeper:       nftkeeper.NewKeeper(storeService, cdc, ak, bk),
		cdc:          cdc,
		storeService: storeService,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/wnft")
}
