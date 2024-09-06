package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wmintkeeper "github.com/st-chain/me-hub/x/wmint/keeper"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"math/big"
)

type Keeper struct {
	*stakingkeeper.Keeper
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	AuthKeeper    banktypes.AccountKeeper
	BankKeeper    types.BankKeeper
	DaoKeeper     types.DaoKeeper
	WMintKeeper   wmintkeeper.Keeper
	nftKeeper     types.NFTKeeper
	wstakingHooks types.WstakingHooks
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak banktypes.AccountKeeper,
	bk types.BankKeeper,
	dk types.DaoKeeper,
	nk types.NFTKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		Keeper:        stakingkeeper.NewKeeper(cdc, storeKey, ak, bk, authority),
		cdc:           cdc,
		storeKey:      storeKey,
		AuthKeeper:    ak,
		BankKeeper:    bk,
		DaoKeeper:     dk,
		nftKeeper:     nk,
		wstakingHooks: nil,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) GetProposerOwnerAddress(ctx sdk.Context) (string, error) {
	header := ctx.BlockHeader()
	addr := header.GetProposerAddress()

	validator, ok := k.GetValidatorByConsAddr(ctx, addr)
	if !ok {
		return "", sdkerrors.Wrapf(types.ErrParameter, "proposer not found")
	}
	return validator.OwnerAddress, nil
}

func (k Keeper) GetStoreKey() storetypes.StoreKey {
	return k.storeKey
}

func (k Keeper) GetCdc() codec.BinaryCodec {
	return k.cdc
}

func (k Keeper) GetPerBlockMintCoinAmount(ctx sdk.Context) (amount big.Int) {
	return k.WMintKeeper.GetPerBlockMintCoinAmount(ctx)
}
