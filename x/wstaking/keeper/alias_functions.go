package keeper

import (
	"context"
	"math"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/params"
	mintTypes "github.com/st-chain/me-hub/x/wmint/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (k Keeper) CalculateInterest(ctx sdk.Context, totalStaking sdkmath.Int, height int64) (rewards sdkmath.LegacyDec, err error) {
	if height >= ctx.BlockHeight() {
		return sdkmath.LegacyZeroDec(), nil
	}
	blockRewards := k.getRewardsByHeight(height, ctx.BlockHeight())
	return k.Calculate(ctx, blockRewards, totalStaking)
}

// getRewardsByHeight Get coins through the block height range
func (k Keeper) getRewardsByHeight(fromHeight int64, toHeight int64) (coin sdkmath.LegacyDec) {
	var totalCoins int64

	lowMul := (fromHeight - 1) / mintTypes.OneYearTotalBlocks
	lowAmount := mintTypes.InitOneYearMintAmount / mintTypes.OneYearTotalBlocks / math.Exp2(float64(lowMul))
	lowMintMEAmount := RoundUpToFourDecimals(lowAmount)
	lowMintUMEAmount := lowMintMEAmount * math.Pow(10, params.BaseDenomUnit)

	highMul := (toHeight - 1) / mintTypes.OneYearTotalBlocks
	highAmount := mintTypes.InitOneYearMintAmount / mintTypes.OneYearTotalBlocks / math.Exp2(float64(highMul))
	highMintMEAmount := RoundUpToFourDecimals(highAmount)
	highMintUMEAmount := highMintMEAmount * math.Pow(10, params.BaseDenomUnit)

	for i := lowMul; i <= highMul; i++ {
		// If the range of from and to are in the same reduction height
		if i == lowMul && lowMul == highMul {
			totalCoins = totalCoins + (toHeight-fromHeight)*int64(lowMintUMEAmount)
			continue
			// Calculate the number of tokens between from and its first cut height
		} else if i == lowMul {
			totalCoins = totalCoins + (mintTypes.OneYearTotalBlocks*(lowMul+1)-fromHeight+1)*int64(lowMintUMEAmount)
			continue
			// Calculate the number of tokens between the last production reduction height and to
		} else if i == highMul {
			totalCoins = totalCoins + (toHeight-mintTypes.OneYearTotalBlocks*(i)-1)*int64(highMintUMEAmount)
			continue
		}

		// Calculate the number of tokens for each full cut interval
		mintAmount := mintTypes.InitOneYearMintAmount / mintTypes.OneYearTotalBlocks / math.Exp2(float64(i))
		mintMEAmount := RoundUpToFourDecimals(mintAmount)
		mintUMEAmount := mintMEAmount * math.Pow(10, params.BaseDenomUnit)
		totalCoins = totalCoins + mintTypes.OneYearTotalBlocks*int64(mintUMEAmount)
	}

	mintedUMECoin := sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(totalCoins))
	coin = sdkmath.LegacyNewDecFromInt(mintedUMECoin.Amount)

	return
}

func (k Keeper) Calculate(ctx sdk.Context, blockRewards sdkmath.LegacyDec, totalStaking sdkmath.Int) (rewards sdkmath.LegacyDec, err error) {
	totalSupply := sdkmath.LegacyNewDec(types.CaclTotalSupply)
	rate := sdkmath.LegacyOneDec().Quo(totalSupply)
	rewards = blockRewards.Mul(sdkmath.LegacyNewDecFromInt(totalStaking).Mul(rate)).Mul(sdkmath.LegacyNewDecWithPrec(1, params.BaseDenomUnit))
	if rewards.LT(sdkmath.LegacyZeroDec()) {
		k.Logger(ctx).Error("Calculate_Interest", "Failed to calculate user revenue！")
		return rewards, types.ErrCalculateInterest.Wrap("withdraw coins amount too small")
	}
	return
}

func RoundUpToFourDecimals(x float64) float64 {
	return math.Ceil(x*10000) / 10000
}

func (k Keeper) Delegation(ctx context.Context, addrDel sdk.AccAddress, addrVal sdk.ValAddress) (stakingtypes.DelegationI, error) {
	bond, err := k.GetDelegation(ctx, addrDel, addrVal)
	if err != nil {
		return nil, err
	}
	return bond, nil
}
