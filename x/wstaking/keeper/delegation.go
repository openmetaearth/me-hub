package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (k Keeper) UnMeidDelegate(
	ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int, validatorAddr sdk.ValAddress) (newShares sdk.Dec, err error) {
	// Get or create the delegation object
	delegation, found := k.GetDelegation(ctx, delAddr, validatorAddr)
	if !found {
		delegation = stakingtypes.NewDelegation(delAddr, validatorAddr, sdk.ZeroDec())
	}
	//if found {
	//	err = k.hooks.BeforeDelegationSharesModified(ctx, delAddr, validatorAddr)
	//}

	if err != nil {
		return sdk.ZeroDec(), err
	}

	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), bondAmt))

	err = k.BankKeeper.DelegateCoinsFromAccountToModule(ctx, delegatorAddress, types.BondedPoolName, coins)
	if err != nil {
		return sdk.Dec{}, err
	}
	delegation.UnMeidAmount = delegation.UnMeidAmount.Add(bondAmt)
	delegation.StartHeight = ctx.BlockHeight()
	allAmount := k.GetAllUnMeidDelegationAmount(ctx)
	if allAmount.Amount.IsNil() {
		allAmount.Amount = sdk.ZeroInt()
	}
	allAmount.Amount = allAmount.Amount.Add(bondAmt)
	allAmount.Denom = k.BondDenom(ctx)
	k.SetAllUnMEIDDelegationAmount(ctx, allAmount)
	k.SetDelegation(ctx, delegation)
	return sdk.NewDecFromInt(delegation.UnMeidAmount), nil
}

// GetAllUnMeidDelegationAmount
func (k Keeper) GetAllUnMeidDelegationAmount(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyPrefix(types.ExperienceRegion + types.PrefixUnMeid)
	var coin sdk.Coin
	value := store.Get(key)
	if value == nil {
		ctx.Logger().Debug("get all UnMeidDelegationAmount err", "not found value by key", key)
		return coin
	}
	err := coin.Unmarshal(value)
	if err != nil {
		ctx.Logger().Error("coin unmarshal err=", err.Error())
	}
	return coin
}

// SetAllUnMEIDDelegationAmount
func (k Keeper) SetAllUnMEIDDelegationAmount(ctx sdk.Context, coin sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	b, err := coin.Marshal()
	if err != nil {
		ctx.Logger().Error("coin marshal err=", err.Error())
		return
	}

	key := types.KeyPrefix(types.ExperienceRegion + types.PrefixUnMeid)
	store.Set(key, b)
}
