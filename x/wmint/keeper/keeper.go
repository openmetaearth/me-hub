package keeper

import (
	"math/big"

	storetypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/openmetaearth/me-hub/x/wmint/types"
)

// Wrapper wraps the original mint keeper and intercepts its original methods if needed.
type Keeper struct {
	mintkeeper.Keeper
	storeService          storetypes.KVStoreService
	bankKeeper            minttypes.BankKeeper
	treasuryModuleAccount string
}

// NewWrappedMint returns a new instance of the WrappedNFTKeeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	sk minttypes.StakingKeeper,
	ak minttypes.AccountKeeper,
	bk minttypes.BankKeeper,
	treasuryModuleAccount string,
	authority string,
) Keeper {
	return Keeper{
		Keeper:                mintkeeper.NewKeeper(cdc, storeService, sk, ak, bk, treasuryModuleAccount, authority),
		storeService:          storeService,
		bankKeeper:            bk,
		treasuryModuleAccount: treasuryModuleAccount,
	}
}

// SetMintedCoinAmount sets the current total minted coins.
func (k Keeper) SetMintedCoinAmount(ctx sdk.Context, amount big.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.CoinAmountKey, amount.Bytes())
}

// GetMintedCoinAmount returns the current total minted coins.
func (k Keeper) GetMintedCoinAmount(ctx sdk.Context) (amount big.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := store.Get(types.CoinAmountKey)
	if b == nil {
		b = []byte{0x00}
	}
	amount.SetBytes(b)
	return
}

// SetPerBlockMintCoinAmount sets the every block mint coins amount.
func (k Keeper) SetPerBlockMintCoinAmount(ctx sdk.Context, amount big.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.PerBlockCoinAmountKey, amount.Bytes())
}

// GetPerBlockMintCoinAmount returns the current block mint coins amount.
func (k Keeper) GetPerBlockMintCoinAmount(ctx sdk.Context) (amount big.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := store.Get(types.PerBlockCoinAmountKey)
	if b == nil {
		b = []byte{0x00}
	}
	amount.SetBytes(b)
	return
}

// SendCoinsToTreasury to be used in BeginBlocker. send coins to me treasury module account
func (k Keeper) SendCoinsToTreasury(ctx sdk.Context, coins sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, k.treasuryModuleAccount, coins)
}
