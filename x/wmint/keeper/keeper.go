package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/openmetaearth/me-hub/x/wmint/types"
	"math/big"
)

// Wrapper wraps the original mint keeper and intercepts its original methods if needed.
type Keeper struct {
	mintkeeper.Keeper
	cdc                   codec.BinaryCodec
	storeKey              storetypes.StoreKey
	bankKeeper            minttypes.BankKeeper
	treasuryModuleAccount string
}

// NewWrappedMint returns a new instance of the WrappedNFTKeeper.
func NewKeeper(cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	sk minttypes.StakingKeeper,
	ak minttypes.AccountKeeper,
	bk minttypes.BankKeeper,
	treasuryModuleAccount string,
	authority string,
) Keeper {
	return Keeper{
		Keeper:                mintkeeper.NewKeeper(cdc, key, sk, ak, bk, treasuryModuleAccount, authority),
		storeKey:              key,
		bankKeeper:            bk,
		treasuryModuleAccount: treasuryModuleAccount,
	}
}

// SetMintedCoinAmount sets the current total minted coins.
func (k Keeper) SetMintedCoinAmount(ctx sdk.Context, amount big.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CoinAmountKey, amount.Bytes())
}

// GetMintedCoinAmount returns the current total minted coins.
func (k Keeper) GetMintedCoinAmount(ctx sdk.Context) (amount big.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.CoinAmountKey)
	if b == nil {
		b = []byte{0x00}
	}
	amount.SetBytes(b)
	return
}

// SetPerBlockMintCoinAmount sets the every block mint coins amount.
func (k Keeper) SetPerBlockMintCoinAmount(ctx sdk.Context, amount big.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PerBlockCoinAmountKey, amount.Bytes())
}

// GetPerBlockMintCoinAmount returns the current block mint coins amount.
func (k Keeper) GetPerBlockMintCoinAmount(ctx sdk.Context) (amount big.Int) {
	store := ctx.KVStore(k.storeKey)
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
