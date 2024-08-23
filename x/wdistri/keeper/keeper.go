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

type (
	Keeper struct {
		*WrapDistrKeeper
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
		MecToUmec       int64
	}
)

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
		WrapDistrKeeper: &WrapDistrKeeper{
			&DistrKeeper,
		},
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
		MecToUmec:       int64(math.Pow(10, ME_EXPONENT)),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) AllocateBlockRewards(ctx sdk.Context, req abci.RequestEndBlock) {
	if ctx.BlockHeight()%oneDayTotalBlocks == 0 {
		fromHeight := req.Height - oneDayTotalBlocks + 1

		toHeight := req.Height + 1
		totalMintCoins := k.getMintCoinsByHeight(fromHeight, toHeight)

		regions := k.stakingKeeper.GetAllRegionI(ctx)
		totalRegionShare := sdkmath.NewInt(0)
		for _, region := range regions {
			totalRegionShare = region.GetRegionShare().Add(totalRegionShare)
		}
		totalRegionShareDec := sdk.NewDecFromInt(totalRegionShare)
		for _, region := range regions {
			// calculate every region coins: RegionShare * totalMintCoins / totalRegionShare
			amount := sdk.NewDecFromInt(region.GetRegionShare()).Mul(totalMintCoins).Quo(totalRegionShareDec)
			regionAmount := amount.TruncateInt()
			regionCoins := sdk.NewCoins(sdk.NewCoin(k.baseDenom, sdk.NewInt(regionAmount.Int64())))
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()), regionCoins)
			if err != nil {
				ctx.Logger().Error(err.Error())
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
	}
}

// getMintCoinsByHeight Get coins through the block height range
func (k Keeper) getMintCoinsByHeight(fromHeight int64, toHeight int64) (coin sdk.Dec) {
	var totalCoins int64
	lowMul := (fromHeight - 1) / oneYearTotalBlocks
	lowAmount := initOneYearMintAmount / oneYearTotalBlocks / math.Exp2(float64(lowMul))
	lowMintMEAmount := RoundUpToFourDecimals(lowAmount)
	lowMintUMEAmount := lowMintMEAmount * float64(k.MecToUmec)

	highMul := (toHeight - 1) / oneYearTotalBlocks
	highAmount := initOneYearMintAmount / oneYearTotalBlocks / math.Exp2(float64(highMul))
	highMintMEAmount := RoundUpToFourDecimals(highAmount)
	highMintUMEAmount := highMintMEAmount * float64(k.MecToUmec)

	for i := lowMul; i <= highMul; i++ {
		// If the range of from and to are in the same reduction height
		if i == lowMul && lowMul == highMul {
			totalCoins = totalCoins + (toHeight-fromHeight)*int64(lowMintUMEAmount)
			continue
			// Calculate the number of tokens between from and its first cut height
		} else if i == lowMul {
			totalCoins = totalCoins + int64(oneYearTotalBlocks*(lowMul+1)-(fromHeight)+1)*int64(lowMintUMEAmount)
			continue
			// Calculate the number of tokens between the last production reduction height and to
		} else if i == highMul {
			totalCoins = totalCoins + int64(toHeight-oneYearTotalBlocks*(i)-1)*int64(highMintUMEAmount)
			continue
		}

		// Calculate the number of tokens for each full cut interval
		mintAmount := initOneYearMintAmount / oneYearTotalBlocks / math.Exp2(float64(i))
		mintMEAmount := RoundUpToFourDecimals(mintAmount)
		mintUMEAmount := mintMEAmount * float64(k.MecToUmec)
		totalCoins = totalCoins + int64(oneYearTotalBlocks)*int64(mintUMEAmount)
	}

	mintedUMECoin := sdk.NewCoin(k.baseDenom, sdk.NewInt(totalCoins))
	coin = sdk.NewDecFromInt(mintedUMECoin.Amount)

	return
}

func RoundUpToFourDecimals(x float64) float64 {
	return math.Ceil(x*10000) / 10000
}
