package wmint

import (
	"math"
	"math/big"
	"time"

	"github.com/st-chain/me-hub/app/params"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/st-chain/me-hub/x/wmint/keeper"
	"github.com/st-chain/me-hub/x/wmint/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper, ic mintypes.InflationCalculationFn) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	logger := k.Logger(ctx)

	mintedAmount := k.GetMintedCoinAmount(ctx)
	blockHeight := ctx.BlockHeight()
	mul := (blockHeight - 1) / types.OneYearTotalBlocks
	amount := types.InitOneYearMintAmount / types.OneYearTotalBlocks / math.Exp2(float64(mul))
	mintingMEAmount := RoundUpToFourDecimals(amount)
	mintingUMEAmount := mintingMEAmount * math.Pow(10, params.BaseDenomUnit)

	// Compare the currently mined coins with the total amount of coins
	// -1 means that the current accumulated amount of mined is smaller than the total amount
	result := mintedAmount.Cmp(big.NewInt(types.TotalMintCoinsAmount))
	if result == -1 {
		// Accumulate the mined coins
		mintedAmount.Add(&mintedAmount, big.NewInt(int64(mintingUMEAmount)))
		k.SetMintedCoinAmount(ctx, mintedAmount)
	} else {
		mintingUMEAmount = 0
	}

	k.SetPerBlockMintCoinAmount(ctx, *big.NewInt(int64(mintingUMEAmount)))

	mintedCoin := sdk.NewCoin(params.BaseDenom, sdk.NewInt(int64(mintingUMEAmount)))
	mintedCoins := sdk.NewCoins(mintedCoin)
	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// send the minted coins to me treasury module account
	err = k.SendCoinsToTreasury(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyAmount, mintedCoin.String()),
		),
	)

	logger.Info(
		"minted coin success",
		"block height", blockHeight,
		"minted umec amount", mintedCoin.Amount.String(),
	)
}

func RoundUpToFourDecimals(x float64) float64 {
	return math.Ceil(x*10000) / 10000
}

func CalculateCoinFromHeightToHeight() {
}

// getMintCoinsByHeight Get coins through the block height range
func getMintCoinsByHeight(fromHeight int64, toHeight int64) (coin sdk.Dec) {
	denomeUnit := 8
	baseDenom := "umec"
	var totalCoins int64
	lowMul := (fromHeight - 1) / types.OneYearTotalBlocks
	lowAmount := types.InitOneYearMintAmount / types.OneYearTotalBlocks / math.Exp2(float64(lowMul))
	lowMintMEAmount := RoundUpToFourDecimals(lowAmount)
	lowMintUMEAmount := lowMintMEAmount * math.Pow(10, float64(denomeUnit))

	highMul := (toHeight - 1) / types.OneYearTotalBlocks
	highAmount := types.InitOneYearMintAmount / types.OneYearTotalBlocks / math.Exp2(float64(highMul))
	highMintMEAmount := RoundUpToFourDecimals(highAmount)
	highMintUMEAmount := highMintMEAmount * math.Pow(10, float64(denomeUnit))

	for i := lowMul; i <= highMul; i++ {
		// If the range of from and to are in the same reduction height
		if i == lowMul && lowMul == highMul {
			totalCoins = totalCoins + (toHeight-fromHeight)*int64(lowMintUMEAmount)
			continue
			// Calculate the number of tokens between from and its first cut height
		} else if i == lowMul {
			totalCoins = totalCoins + int64(types.OneYearTotalBlocks*(lowMul+1)-(fromHeight)+1)*int64(lowMintUMEAmount)
			continue
			// Calculate the number of tokens between the last production reduction height and to
		} else if i == highMul {
			totalCoins = totalCoins + int64(toHeight-types.OneYearTotalBlocks*(i)-1)*int64(highMintUMEAmount)
			continue
		}

		// Calculate the number of tokens for each full cut interval
		mintAmount := types.InitOneYearMintAmount / types.OneYearTotalBlocks / math.Exp2(float64(i))
		mintMEAmount := RoundUpToFourDecimals(mintAmount)
		mintUMEAmount := mintMEAmount * math.Pow(10, float64(denomeUnit))
		totalCoins = totalCoins + int64(types.OneYearTotalBlocks)*int64(mintUMEAmount)
	}

	mintedUMECoin := sdk.NewCoin(baseDenom, sdk.NewInt(totalCoins))
	coin = sdk.NewDecFromInt(mintedUMECoin.Amount)

	return
}
