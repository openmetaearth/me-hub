package keeper

import (
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distriKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/st-chain/me-hub/x/wdistri/types"
)

type Keeper struct {
	distriKeeper.Keeper
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	memKey        storetypes.StoreKey
	paramstore    paramtypes.Subspace
	authKeeper    types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	feeCollectorName string // name of the FeeCollector ModuleAccount
	baseDenom        string
	MecToUmec        int64
}

type WrapDistrKeeper struct {
	*distriKeeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	feeCollectorName string,
	authority string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	DistrKeeper := distriKeeper.NewKeeper(
		cdc,
		storeKey,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		feeCollectorName,
		authority,
	)
	baseDenom, err := sdk.GetBaseDenom()
	if err != nil {
		panic("GetBaseDenom failed")
	}
	return &Keeper{
		Keeper:           DistrKeeper,
		cdc:              cdc,
		storeKey:         storeKey,
		memKey:           memKey,
		paramstore:       ps,
		authKeeper:       accountKeeper,
		bankKeeper:       bankKeeper,
		stakingKeeper:    stakingKeeper,
		authority:        authority,
		feeCollectorName: feeCollectorName,
		baseDenom:        baseDenom,
		MecToUmec:        int64(math.Pow(10, ME_EXPONENT)),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) AllocateBlockRewardEveryday(ctx sdk.Context, req abci.RequestEndBlock) error {
	if ctx.BlockHeight()%oneDayTotalBlocks == 0 {
		return k.AllocateBlockReward(ctx)
	}
	return nil
}

func (k Keeper) AllocateBlockReward(ctx sdk.Context) error {
	feeCollectorAddr := k.authKeeper.GetModuleAddress(k.feeCollectorName)
	totalMintCoin := k.bankKeeper.GetBalance(ctx, feeCollectorAddr, k.baseDenom)

	regions := k.stakingKeeper.GetAllRegionI(ctx)
	totalRegionShare := sdkmath.NewInt(0)
	for _, region := range regions {
		totalRegionShare = region.GetRegionShare().Add(totalRegionShare)
	}
	totalRegionShareDec := sdk.NewDecFromInt(totalRegionShare)
	if totalRegionShare.IsZero() {
		return fmt.Errorf("totalRegionShare is zero, cannot divide by zero")
	}
	for _, region := range regions {
		// calculate every region coins: RegionShare * totalMintCoins / totalRegionShare

		amount := sdk.NewDecFromInt(region.GetRegionShare()).Mul(totalMintCoin.Amount.ToLegacyDec()).Quo(totalRegionShareDec)
		regionAmount := amount.TruncateInt()
		regionCoins := sdk.NewCoins(sdk.NewCoin(k.baseDenom, sdk.NewInt(regionAmount.Int64())))
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()), regionCoins)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeRegionTreasuryReword,
				sdk.NewAttribute(types.AttributeKeyRegionTreasuryAddress, region.GetRegionTreasureAddr()),
				sdk.NewAttribute(types.AttributeKeyRegionId, region.GetRegionId()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, regionCoins.String()),
			),
		})
	}
	return nil
}
