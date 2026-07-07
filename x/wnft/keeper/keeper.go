package keeper

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/x/nft"
	nftkeeper "cosmossdk.io/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
	nftkeeper.Keeper
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	ak nft.AccountKeeper,
	bk nft.BankKeeper,
) *Keeper {
	return &Keeper{
		Keeper:       nftkeeper.NewKeeper(storeService, cdc, ak, bk),
		cdc:          cdc,
		storeService: storeService,
	}
}
