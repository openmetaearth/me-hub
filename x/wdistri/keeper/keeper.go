package keeper

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distriKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wdistri/types"
)

type Keeper struct {
	distriKeeper.Keeper
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	paramstore    paramtypes.Subspace
	authKeeper    types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	feeCollectorName string // name of the FeeCollector ModuleAccount
}

type WrapDistrKeeper struct {
	*distriKeeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	feeCollectorName string,
	authority string,
) *Keeper {
	DistrKeeper := distriKeeper.NewKeeper(
		cdc,
		storeKey,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		feeCollectorName,
		authority,
	)
	return &Keeper{
		Keeper:           DistrKeeper,
		cdc:              cdc,
		storeKey:         storeKey,
		paramstore:       ps,
		authKeeper:       accountKeeper,
		bankKeeper:       bankKeeper,
		stakingKeeper:    stakingKeeper,
		authority:        authority,
		feeCollectorName: feeCollectorName,
	}
}

func (k Keeper) GetTreasuryModuleAccount() string {
	return k.feeCollectorName
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", distributiontypes.ModuleName))
}

func (k Keeper) AllocateBlockRewardEveryday(ctx sdk.Context, req abci.RequestEndBlock) error {
	if ctx.BlockHeight()%types.OneDayTotalBlocks == 0 {
		return k.AllocateBlockReward(ctx)
	}
	return nil
}

func (k Keeper) AllocateBlockReward(ctx sdk.Context) error {
	feeCollectorAddr := k.authKeeper.GetModuleAddress(k.GetTreasuryModuleAccount())
	totalMintCoin := k.bankKeeper.GetAllBalances(ctx, feeCollectorAddr)
	if totalMintCoin.AmountOf(params.BaseDenom).IsZero() {
		ctx.Logger().Info("totalMintCoin is zero, no need to allocate reward")
		return nil
	}
	regions := k.stakingKeeper.GetAllRegionI(ctx)
	totalRegionShare := sdkmath.NewInt(0)
	for _, region := range regions {
		totalRegionShare = region.GetRegionShare().Add(totalRegionShare)
	}
	totalRegionShareDec := sdk.NewDecFromInt(totalRegionShare)
	if totalRegionShare.IsZero() {
		return nil
	}
	for _, region := range regions {
		// calculate every region coins: RegionShare * totalMintCoins / totalRegionShare
		amount := sdk.NewDecFromInt(region.GetRegionShare()).Mul(totalMintCoin.AmountOf(params.BaseDenom).ToLegacyDec()).Quo(totalRegionShareDec)
		regionAmount := amount.TruncateInt()
		regionCoins := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(regionAmount.Int64())))
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.GetTreasuryModuleAccount(), sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()), regionCoins)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeRegionTreasuryReward,
				sdk.NewAttribute(types.AttributeKeyRegionTreasuryAddress, region.GetRegionTreasureAddr()),
				sdk.NewAttribute(types.AttributeKeyRegionId, region.GetRegionId()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, regionCoins.String()),
			),
		})
	}
	return nil
}
