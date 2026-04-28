package wmint

import (
	"github.com/openmetaearth/me-hub/app/params"
	"math/big"
	"time"

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
	halvingDivisor := sdk.NewDecFromBigInt(new(big.Int).Lsh(big.NewInt(1), uint(mul)))
	amount := sdk.NewDec(int64(types.InitOneYearMintAmount)).
		Quo(sdk.NewDec(int64(types.OneYearTotalBlocks))).
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
			mintingUMECAmount = sdk.NewIntFromBigInt(remaining)
			mintedAmount.Set(totalCap)
		} else {
			mintedAmount.Set(newMinted)
		}
		k.SetMintedCoinAmount(ctx, mintedAmount)
	default:
		mintingUMECAmount = sdk.ZeroInt()
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

// RoundUpToFourDecimalsDec rounds x up to 4 decimal places using sdk.Dec arithmetic.
func RoundUpToFourDecimalsDec(x sdk.Dec) sdk.Dec {
	return x.MulInt64(10000).Ceil().QuoInt64(10000)
}
