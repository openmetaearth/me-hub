package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) validateDelegateInterestPayout(ctx sdk.Context, region types.Region, rewards sdk.Dec) (sdk.AccAddress, sdk.Coin, error) {
	regionTreasureAddr, err := sdk.AccAddressFromBech32(region.RegionTreasureAddr)
	if err != nil {
		return nil, sdk.Coin{}, err
	}

	if !region.DelegateInterest.GTE(rewards) {
		return nil, sdk.Coin{}, types.ErrPayInterest.Wrapf(
			"region(%s) total interest not enough. need pay %s, only have %s",
			region.RegionId,
			rewards.String(),
			region.DelegateInterest.String(),
		)
	}

	rewardCoin := sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())
	treasuryBalance := k.bankKeeper.GetBalance(ctx, regionTreasureAddr, params.BaseDenom)
	if rewardCoin.Amount.IsPositive() {
		if treasuryBalance.IsLT(rewardCoin) {
			return nil, sdk.Coin{}, types.ErrPayInterest.Wrapf(
				"region(%s) treasury balance not enough. need pay %s, only have %s",
				region.RegionId,
				rewardCoin.String(),
				treasuryBalance.String(),
			)
		}
	}

	remainingInterest := region.DelegateInterest.Sub(rewards)
	projectedTreasuryBalance := treasuryBalance.Amount.Sub(rewardCoin.Amount)
	if projectedTreasuryBalance.LT(remainingInterest.RoundInt()) {
		return nil, sdk.Coin{}, types.ErrPayInterest.Wrapf(
			"region(%s) treasury reserve not enough after payout. need reserve %s%s, projected balance %s%s",
			region.RegionId,
			remainingInterest.RoundInt().String(),
			params.BaseDenom,
			projectedTreasuryBalance.String(),
			params.BaseDenom,
		)
	}

	return regionTreasureAddr, rewardCoin, nil
}
