package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"time"
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

// Undelegate unbonds an amount of delegator shares from a given validator. It
// will verify that the unbonding entries between the delegator and validator
// are not exceeded and unbond the staked tokens (based on shares) by creating
// an unbonding object and inserting it into the unbonding queue which will be
// processed during the staking EndBlocker.
func (k Keeper) Undelegate(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, isMeid bool, anmout math.Int) (time.Time, math.Int, error) {
	//validator, found := k.GetValidator(ctx, valAddr)
	//if !found {
	//	return time.Time{}, types.ErrNoDelegatorForAddress
	//}
	if !isMeid {
		if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
			return time.Time{}, anmout, types.ErrMaxUnbondingDelegationEntries
		}
	}
	returnAmount, err := k.Unbond(ctx, delAddr, valAddr, anmout, isMeid)
	if err != nil {
		return time.Time{}, anmout, err
	}
	completionTime := time.Time{}
	// transfer the validator tokens to the not bonded pool
	if !isMeid {
		k.bondedTokensToNotBonded(ctx, returnAmount)
		completionTime = ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
		ubd := k.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
		k.InsertUBDQueue(ctx, ubd, completionTime)

		allAmount := k.GetAllUnMeidDelegationAmount(ctx)
		if allAmount.Amount.IsNil() {
			allAmount.Amount = sdk.ZeroInt()
		}
		allAmount.Amount = allAmount.Amount.Sub(returnAmount)
		allAmount.Denom = k.BondDenom(ctx)
		k.SetAllUnMEIDDelegationAmount(ctx, allAmount)
	} else {
		amt := sdk.NewCoin(params.BaseDenom, returnAmount)
		err = k.BankKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.BondedPoolName, delAddr, sdk.NewCoins(amt))
		if err != nil {
			return completionTime, returnAmount, err
		}
	}

	return completionTime, returnAmount, nil
}

// Unbond unbonds a particular delegation and perform associated store operations.
func (k Keeper) Unbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, delAmount math.Int, isMeid bool) (amount math.Int, err error) {
	// check if a delegation object exists in the store
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return amount, types.ErrNoDelegatorForAddress
	}

	overAmount := sdk.ZeroInt()
	if isMeid {
		if delegation.Amount.LTE(sdk.ZeroInt()) {
			return amount, types.ErrNotEnoughDelegationAmount
		}
		if delAmount.GTE(delegation.Amount) {
			delAmount = delegation.Amount
			delegation.Amount = sdk.ZeroInt()
		} else {
			delegation.Amount = delegation.Amount.Sub(delAmount)
		}
		overAmount = delegation.Amount
	} else {
		if delegation.UnMeidAmount.LTE(sdk.ZeroInt()) {
			return amount, types.ErrNotEnoughDelegationAmount
		}
		if delAmount.GTE(delegation.UnMeidAmount) {
			delAmount = delegation.UnMeidAmount
			delegation.UnMeidAmount = sdk.ZeroInt()
		} else {
			delegation.UnMeidAmount = delegation.UnMeidAmount.Sub(delAmount)
		}
		overAmount = delegation.UnMeidAmount
	}
	err = types.CheckMinDelegate(overAmount)
	if err != nil {
		amount = delAmount.Add(overAmount)
		if isMeid {
			delegation.Amount = sdk.ZeroInt()
		} else {
			delegation.UnMeidAmount = sdk.ZeroInt()
		}
	} else {
		amount = delAmount
	}
	// call the before-delegation-modified hook
	if err = k.Hooks().BeforeDelegationSharesModified(ctx, delAddr, valAddr); err != nil {
		return amount, err
	}
	delegation.StartHeight = ctx.BlockHeight()
	k.SetDelegation(ctx, delegation)
	return amount, nil
}

// bondedTokensToNotBonded transfers coins from the bonded to the not bonded pool within staking
func (k Keeper) bondedTokensToNotBonded(ctx sdk.Context, tokens math.Int) {
	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), tokens))
	if err := k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.BondedPoolName, types.NotBondedPoolName, coins); err != nil {
		panic(err)
	}
}
