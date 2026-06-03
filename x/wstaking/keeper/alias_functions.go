package keeper

import (
	"math/big"

	"github.com/openmetaearth/me-hub/x/wmint"

	cmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	mintTypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) CalculateInterest(ctx sdk.Context, totalStaking cmath.Int, height int64) (rewards sdk.Dec, err error) {
	if height >= ctx.BlockHeight() {
		return sdk.ZeroDec(), nil
	}
	blockRewards := k.getRewardsByHeight(height, ctx.BlockHeight())
	return k.Calculate(blockRewards, totalStaking), nil
}

// getRewardsByHeight Get coins through the block height range
func (k Keeper) getRewardsByHeight(fromHeight int64, toHeight int64) (coin sdk.Dec) {
	totalCoins := sdk.ZeroInt()

	lowMul := (fromHeight - 1) / mintTypes.OneYearTotalBlocks
	highMul := (toHeight - 1) / mintTypes.OneYearTotalBlocks

	for i := lowMul; i <= highMul; i++ {
		halvingDivisor := sdk.NewDecFromBigInt(new(big.Int).Lsh(big.NewInt(1), uint(i)))
		amountDec := sdk.NewDec(int64(mintTypes.InitOneYearMintAmount)).
			Quo(sdk.NewDec(int64(mintTypes.OneYearTotalBlocks))).
			Quo(halvingDivisor)
		mintUMECAmount := wmint.RoundUpToFourDecimalsDec(amountDec).MulInt64(100_000_000).TruncateInt()

		var blockCount int64
		// If the range of from and to are in the same reduction period
		if i == lowMul && lowMul == highMul {
			blockCount = toHeight - fromHeight
			// Calculate the number of tokens between fromHeight and its first halving boundary
		} else if i == lowMul {
			blockCount = int64(mintTypes.OneYearTotalBlocks)*(lowMul+1) - fromHeight + 1
			// Calculate the number of tokens between the last halving boundary and toHeight
		} else if i == highMul {
			blockCount = toHeight - int64(mintTypes.OneYearTotalBlocks)*i - 1
		} else {
			// Calculate the number of tokens for each full halving interval
			blockCount = int64(mintTypes.OneYearTotalBlocks)
		}

		totalCoins = totalCoins.Add(mintUMECAmount.MulRaw(blockCount))
	}

	coin = sdk.NewDecFromInt(totalCoins)
	return
}

// Calculate computes the per-block staking reward for a given amount of staked tokens.
// It is a pure arithmetic function and does not read from chain state.
// Returns ZeroDec when totalStaking is zero or the result would be non-positive.
func (k Keeper) Calculate(blockRewards sdk.Dec, totalStaking cmath.Int) sdk.Dec {
	if totalStaking.IsZero() {
		return sdk.ZeroDec()
	}
	totalSupply := sdk.NewDec(types.CaclTotalSupply)
	rate := sdk.OneDec().Quo(totalSupply)
	rewards := blockRewards.Mul(sdk.NewDecFromInt(totalStaking).Mul(rate)).Mul(sdk.NewDecWithPrec(1, params.BaseDenomUnit))
	if rewards.IsNegative() {
		return sdk.ZeroDec()
	}
	return rewards
}

// Delegation get the delegation interface for a particular set of delegator and validator addresses
func (k Keeper) Delegation(ctx sdk.Context, addrDel sdk.AccAddress, addrVal sdk.ValAddress) stakingtypes.DelegationI {
	bond, ok := k.GetDelegation(ctx, addrDel, addrVal)
	if !ok {
		return nil
	}

	return bond
}
