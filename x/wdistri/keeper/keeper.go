package keeper

import (
	"cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distriKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wdistri/types"
)

type Keeper struct {
	distriKeeper.Keeper
	cdc           codec.BinaryCodec
	authKeeper    distritypes.AccountKeeper
	bankKeeper    distritypes.BankKeeper
	stakingKeeper distritypes.StakingKeeper
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
	storeService store.KVStoreService,
	accountKeeper distritypes.AccountKeeper,
	bankKeeper distritypes.BankKeeper,
	stakingKeeper distritypes.StakingKeeper,
	feeCollectorName string,
	authority string,
) *Keeper {
	distrKeeper := distriKeeper.NewKeeper(
		cdc,
		storeService,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		feeCollectorName,
		authority,
	)
	return &Keeper{
		Keeper:           distrKeeper,
		cdc:              cdc,
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

func (k Keeper) AllocateBlockRewardEveryday(ctx sdk.Context) error {
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
	totalRegionShareDec := sdkmath.LegacyNewDecFromInt(totalRegionShare)
	if totalRegionShare.IsZero() {
		return nil
	}
	for _, region := range regions {
		// calculate every region coins: RegionShare * totalMintCoins / totalRegionShare
		amount := sdkmath.LegacyNewDecFromInt(region.GetRegionShare()).Mul(totalMintCoin.AmountOf(params.BaseDenom).ToLegacyDec()).Quo(totalRegionShareDec)
		regionAmount := amount.TruncateInt()
		regionCoins := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(regionAmount.Int64())))
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
