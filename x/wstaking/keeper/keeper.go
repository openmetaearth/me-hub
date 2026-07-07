package keeper

import (
	"math/big"
	"sync"

	addresscodec "cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

type Keeper struct {
	*stakingkeeper.Keeper
	cdc             codec.BinaryCodec
	storeKey        storetypes.StoreKey
	authKeeper      banktypes.AccountKeeper
	bankKeeper      types.BankKeeper
	daoKeeper       types.DaoKeeper
	mintKeeper      types.MintKeeper
	nftKeeper       types.NFTKeeper
	wstakingHooks   types.WstakingHooks
	kycKeeper       types.KycKeeper
	didKeeper       types.DidKeeper
	groupKeeper     types.GroupKeeper
	slashingKeeper  slashingkeeper.Keeper
	regions         *sync.Map
	regionCacheOnce sync.Once
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	storeService corestore.KVStoreService,
	ak banktypes.AccountKeeper,
	bk types.BankKeeper,
	dk types.DaoKeeper,
	nk types.NFTKeeper,
	authority string,
	validatorAddressCodec addresscodec.Codec,
	consensusAddressCodec addresscodec.Codec,
) *Keeper {
	nativeKeeper := stakingkeeper.NewKeeper(cdc, storeService, ak, bk, authority, validatorAddressCodec, consensusAddressCodec)
	return &Keeper{
		Keeper:          nativeKeeper,
		cdc:             cdc,
		storeKey:        storeKey,
		authKeeper:      ak,
		bankKeeper:      bk,
		daoKeeper:       dk,
		nftKeeper:       nk,
		regions:         new(sync.Map),
		regionCacheOnce: sync.Once{},
	}
}

func (k *Keeper) SetMintKeeper(mintKeeper types.MintKeeper) {
	k.mintKeeper = mintKeeper
}

func (k *Keeper) SetKycKeeper(keeper types.KycKeeper) {
	k.kycKeeper = keeper
}

func (k *Keeper) SetGroupKeeper(keeper types.GroupKeeper) {
	k.groupKeeper = keeper
}

func (k *Keeper) SetDidKeeper(keeper types.DidKeeper) {
	k.didKeeper = keeper
}

func (k *Keeper) SetSlashingKeeper(keeper slashingkeeper.Keeper) {
	k.slashingKeeper = keeper
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k *Keeper) GetProposerOwnerAddress(ctx sdk.Context) (string, error) {
	header := ctx.BlockHeader()
	addr := header.GetProposerAddress()
	validator, err := k.GetValidatorByConsAddr(ctx, addr)
	if err != nil {
		return "", errorsmod.Wrapf(types.ErrParameter, "proposer not found")
	}
	return validator.OwnerAddress, nil
}

func (k *Keeper) GetStoreKey() storetypes.StoreKey {
	return k.storeKey
}

func (k *Keeper) GetCdc() codec.BinaryCodec {
	return k.cdc
}

func (k *Keeper) GetPerBlockMintCoinAmount(ctx sdk.Context) (amount big.Int) {
	return k.mintKeeper.GetPerBlockMintCoinAmount(ctx)
}

func (k *Keeper) GetValOwnerAddress(ctx sdk.Context, regionId string) (string, error) {
	region, ok := k.GetRegion(ctx, regionId)
	if !ok {
		return "", errorsmod.Wrapf(types.ErrRegionNotExist, "region(%s) not found", regionId)
	}

	valAddr, err := sdk.ValAddressFromBech32(region.OperatorAddress)
	if err != nil {
		return "", errorsmod.Wrapf(types.ErrParameter, "region bonded validator address(%s) invalid", region.OperatorAddress)
	}

	validator, err := k.GetValidator(ctx, valAddr)
	if err != nil {
		return "", errorsmod.Wrapf(stakingtypes.ErrNoValidatorFound, "region bonded validator(%s) not found", valAddr.String())
	}
	return validator.OwnerAddress, nil
}
