package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strings"
	"time"
)

func (k Keeper) UnMeidDelegate(
	ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int, validator stakingtypes.Validator) (newShares sdk.Dec, err error) {
	// Get or create the delegation object
	delegation, found := k.GetDelegation(ctx, delAddr, validator.GetOperator())
	if !found {
		delegation = types.NewDelegation(delAddr, validator.GetOperator(), sdk.ZeroDec())
	}

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
			return time.Time{}, anmout, stakingtypes.ErrMaxUnbondingDelegationEntries
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
		return amount, stakingtypes.ErrNoDelegatorForAddress
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

// Delegate performs a delegation, set/update everything necessary within the store.
// tokenSrc indicates the bond status of the incoming funds.
func (k Keeper) Delegate(
	ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus,
	validator stakingtypes.Validator, subtractAccount bool,
) (newShares sdk.Dec, err error) {
	// In some situations, the exchange rate becomes invalid, e.g. if
	// Validator loses all tokens due to slashing. In this case,
	// make all future delegations invalid.
	if validator.InvalidExRate() {
		return math.LegacyZeroDec(), stakingtypes.ErrDelegatorShareExRateInvalid
	}

	// Get or create the delegation object
	delegation, found := k.GetDelegation(ctx, delAddr, sdk.ValAddress{})
	if !found {
		delegation = stakingtypes.NewDelegation(delAddr, sdk.ValAddress{}, math.LegacyZeroDec())
		delegation.Amount = sdk.ZeroInt()
	}

	// call the appropriate hook if present
	if found {
		_, err = k.WithdrawDelegationRewards(ctx, delAddr, sdk.ValAddress{})
	}
	if err != nil {
		return math.LegacyZeroDec(), err
	}
	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
	if tokenSrc == stakingtypes.Bonded {
		panic("delegation token source cannot be bonded")
	}

	var sendName string

	switch {
	case validator.IsBonded():
		sendName = types.BondedPoolName
	case validator.IsUnbonding(), validator.IsUnbonded():
		sendName = types.NotBondedPoolName
	default:
		panic("invalid validator status")
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), bondAmt))
	if err = k.BankKeeper.DelegateCoinsFromAccountToModule(ctx, delegatorAddress, sendName, coins); err != nil {
		return sdk.Dec{}, err
	}
	_, newShares = k.AddValidatorTokensAndShares(ctx, validator, bondAmt)
	// Update delegation
	delegation.Amount = delegation.Amount.Add(bondAmt)
	delegation.StartHeight = ctx.BlockHeight()
	delegation.Shares = delegation.Shares.Add(newShares)
	k.SetDelegation(ctx, delegation)

	return newShares, nil
}

// WithdrawDelegationRewards withdraw rewards from a delegation
func (k Keeper) WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	meid, isMeid := k.GetMeid(ctx, delAddr.String())
	regionID := strings.ToLower(types.ExperienceRegion)
	if isMeid {
		regionID = meid.RegionId
	}
	region, hasRegion := k.GetRegion(ctx, regionID)
	if !hasRegion {
		return nil, fmt.Errorf("%s region not exist", regionID)
	}
	rewards, err := k.internalWithdrawDelegationRewards(ctx, delAddr, region)
	if err != nil {
		return nil, err
	}
	if rewards.IsZero() {
		baseDenom, _ := sdk.GetBaseDenom()
		rewards = sdk.Coins{sdk.Coin{
			Denom:  baseDenom,
			Amount: sdk.ZeroInt(),
		}}
	}
	return rewards, nil
}

func (k Keeper) internalWithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, region types.Region) (sdk.Coins, error) {
	valAddr, valErr := sdk.ValAddressFromBech32(region.OperatorAddress)
	if valErr != nil {
		k.Logger(ctx).Error("internalWithdrawDelegationRewards err=", valErr.Error())
		return nil, valErr
	}
	del := k.Delegation(ctx, delAddr, valAddr)
	if del == nil {
		return nil, types.ErrEmptyDelegationDistInfo
	}

	delegation, isOK := del.(stakingtypes.Delegation)
	if !isOK {
		panic("withdrawDelegationRewards err:type Delegation assertion failed")
		return nil, types.ErrAssertionFailed
	}

	rewards, err := k.CalculateInterest(ctx, delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable), delegation.StartHeight)
	if err != nil {
		return nil, types.ErrCalculateInterest.Wrap(err.Error())
	}
	if region.DelegateInterest.GTE(rewards) {
		region.DelegateInterest = region.DelegateInterest.Sub(rewards)
	} else {
		return nil, types.ErrCalculateInterest.Wrap(fmt.Sprintf("distribution reward.region(%s) total interest not enough.need pay %s,only have %s",
			region.RegionId, rewards.String(), region.DelegateInterest.String()))
	}
	// truncate reward dec coins, return remainder to community pool
	//finalRewards, remainder := rewards.TruncateDecimal()
	coin := sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())
	coins := sdk.NewCoins(coin)
	// add coins to user account
	if !coin.Amount.IsZero() {
		err = k.BankKeeper.SendCoins(ctx, sdk.MustAccAddressFromBech32(region.RegionTreasureAddr), del.GetDelegatorAddr(), coins)
		if err != nil {
			return nil, err
		}
		delegation.StartHeight = ctx.BlockHeight()
		k.SetDelegation(ctx, delegation)
		k.SetRegion(ctx, region)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawDelegatorReward,
			sdk.NewAttribute(types.AttributeKeyValidator, region.OperatorAddress),
			sdk.NewAttribute(types.AttributeKeyDelegator, delAddr.String()),
			sdk.NewAttribute(types.AttributeKeyRegionTreasuryAddress, region.RegionTreasureAddr),
			sdk.NewAttribute(types.AttributeKeyRegionId, region.RegionId),
			sdk.NewAttribute(types.AttributeKeyAmountDelegateInterest, region.DelegateInterest.String()+params.BaseDenom),
			sdk.NewAttribute(types.AttributeKeyPersonalDelegateInterest, rewards.TruncateInt().String()+params.BaseDenom),
		),
	})
	return coins, nil
}

func NewDelegationResp_new(del stakingtypes.Delegation, balance sdk.Coin) stakingtypes.DelegationResponse {
	return stakingtypes.DelegationResponse{
		Delegation: del,
		Balance:    balance,
	}
}
