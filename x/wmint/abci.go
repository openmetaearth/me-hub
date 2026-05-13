package wmint

import (
	"math"
	"math/big"
	"time"

	"github.com/openmetaearth/me-hub/app/params"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/openmetaearth/me-hub/x/wmint/keeper"
	"github.com/openmetaearth/me-hub/x/wmint/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper, ic mintypes.InflationCalculationFn) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	logger := k.Logger(ctx)

	mintedAmount := k.GetMintedCoinAmount(ctx)
	blockHeight := ctx.BlockHeight()
	mul := (blockHeight - 1) / types.OneYearTotalBlocks

	// amount (MEC) = InitOneYearMintAmount / OneYearTotalBlocks / 2^mul
	halvingDivisor := sdkmath.LegacyNewDecFromBigInt(new(big.Int).Lsh(big.NewInt(1), uint(mul)))
	amount := sdkmath.LegacyNewDec(int64(types.InitOneYearMintAmount)).
		Quo(sdkmath.LegacyNewDec(int64(types.OneYearTotalBlocks))).
		Quo(halvingDivisor)

	// RoundUpToFourDecimals: Ceil(amount * 10000) / 10000
	// Then convert MEC to umec: multiply by 10^BaseDenomUnit (=100_000_000), truncate to integer
	mintingUMECAmount := RoundUpToFourDecimalsDec(amount).MulInt64(100_000_000).TruncateInt()

	// Compare the currently mined coins with the total amount of coins
	// Cmp returns -1 (below cap), 0 (exactly at cap), or 1 (above cap)
	totalCap := big.NewInt(types.TotalMintCoinsAmount)
	switch mintedAmount.Cmp(totalCap) {
	case -1:
		newMinted := new(big.Int).Add(&mintedAmount, mintingUMECAmount.BigInt())
		// Clamp to total cap: if adding this block's amount would exceed the cap,
		// only mint the remaining amount to avoid over-issuance.
		if newMinted.Cmp(totalCap) > 0 {
			remaining := new(big.Int).Sub(totalCap, &mintedAmount)
			mintingUMECAmount = sdkmath.NewIntFromBigInt(remaining)
			mintedAmount.Set(totalCap)
		} else {
			mintedAmount.Set(newMinted)
		}
		k.SetMintedCoinAmount(ctx, mintedAmount)
	default:
		mintingUMECAmount = sdkmath.ZeroInt()
	}

	k.SetPerBlockMintCoinAmount(ctx, *mintingUMECAmount.BigInt())

	mintedCoin := sdk.NewCoin(params.BaseDenom, mintingUMECAmount)
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

// RoundUpToFourDecimalsDec rounds x up to 4 decimal places using sdkmath.LegacyDec arithmetic.
func RoundUpToFourDecimalsDec(x sdkmath.LegacyDec) sdkmath.LegacyDec {
	return x.MulInt64(10000).Ceil().QuoInt64(10000)
}

// getMintCoinsByHeight Get coins through the block height range
func getMintCoinsByHeight(fromHeight int64, toHeight int64) (coin sdkmath.LegacyDec) {
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

	mintedUMECoin := sdk.NewCoin(baseDenom, sdkmath.NewInt(totalCoins))
	coin = sdkmath.LegacyNewDecFromInt(mintedUMECoin.Amount)

	return
}
