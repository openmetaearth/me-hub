package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

type Keeper struct {
	*stakingkeeper.Keeper
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	AuthKeeper banktypes.AccountKeeper
	BankKeeper types.BankKeeper
	DaoKeeper  types.DaoKeeper
	nftKeeper  types.NFTKeeper
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
		Keeper:     stakingkeeper.NewKeeper(cdc, storeKey, ak, bk, authority),
		cdc:        cdc,
		storeKey:   storeKey,
		AuthKeeper: ak,
		BankKeeper: bk,
		DaoKeeper:  dk,
		nftKeeper:  nk,
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
