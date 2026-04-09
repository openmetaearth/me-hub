package keeper

import (
	"context"
	"fmt"
	"time"

	didtypes "github.com/st-chain/me-hub/x/did/types"
	kyctypes "github.com/st-chain/me-hub/x/kyc/types"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

const unbondingTime = time.Hour * 24 * 7

// Undelegate unbonds an amount of delegator shares from a given validator. It
// will verify that the unbonding entries between the delegator and validator
// are not exceeded and unbond the staked tokens (based on shares) by creating
// an unbonding object and inserting it into the unbonding queue which will be
// processed during the staking EndBlocker.
func (k Keeper) Undelegate(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, isMeid bool, anmout sdkmath.Int, delegation stakingtypes.Delegation) (time.Time, sdkmath.Int, error) {
	if has, _ := k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr); has {
		return time.Time{}, anmout, stakingtypes.ErrMaxUnbondingDelegationEntries
	}

	returnAmount, err := k.Unbond(ctx, anmout, isMeid, delegation)
	if err != nil {
		return time.Time{}, anmout, err
	}
	completionTime := time.Time{}
	// transfer the validator tokens to the not bonded pool
	if !isMeid {
		k.bondedTokensToNotBonded(ctx, returnAmount)
		completionTime = ctx.BlockHeader().Time.Add(unbondingTime)
		ubd, err := k.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
		if err != nil {
			return time.Time{}, anmout, err
		}
		err = k.InsertUBDQueue(ctx, ubd, completionTime)
		if err != nil {
			return time.Time{}, anmout, err
		}
	} else {
		amt := sdk.NewCoin(params.BaseDenom, returnAmount)
		err = k.bankKeeper.UndelegateCoinsFromModuleToAccount(ctx, stakingtypes.BondedPoolName, delAddr, sdk.NewCoins(amt))
		if err != nil {
			return completionTime, returnAmount, err
		}
	}
	return completionTime, returnAmount, nil
}

// Unbond unbonds a particular delegation and perform associated store operations.
func (k Keeper) Unbond(ctx sdk.Context, delAmount sdkmath.Int, isMeid bool, delegation stakingtypes.Delegation) (amount sdkmath.Int, err error) {
	// check if a delegation object exists in the store
	overAmount := sdkmath.ZeroInt()
	if isMeid {
		if delegation.Amount.LTE(sdkmath.ZeroInt()) {
			return amount, types.ErrNotEnoughDelegationAmount
		}
		if delAmount.GTE(delegation.Amount) {
			delAmount = delegation.Amount
			delegation.Amount = sdkmath.ZeroInt()
		} else {
			delegation.Amount = delegation.Amount.Sub(delAmount)
		}
		overAmount = delegation.Amount
	} else {
		if delegation.UnMeidAmount.LTE(sdkmath.ZeroInt()) {
			return amount, types.ErrNotEnoughDelegationAmount
		}
		if delAmount.GTE(delegation.UnMeidAmount) {
			delAmount = delegation.UnMeidAmount
			delegation.UnMeidAmount = sdkmath.ZeroInt()
		} else {
			delegation.UnMeidAmount = delegation.UnMeidAmount.Sub(delAmount)
		}
		overAmount = delegation.UnMeidAmount
	}
	err = types.CheckMinDelegate(overAmount)
	if err != nil {
		amount = delAmount.Add(overAmount)
		if isMeid {
			delegation.Amount = sdkmath.ZeroInt()
		} else {
			delegation.UnMeidAmount = sdkmath.ZeroInt()
		}
	} else {
		amount = delAmount
	}
	delegation.StartHeight = ctx.BlockHeight()
	if delegation.UnMeidAmount.Add(delegation.Amount).Add(delegation.Unmovable).Equal(sdkmath.ZeroInt()) {
		err = k.removeDelegation(ctx, delegation)
		if err != nil {
			return amount, err
		}
	} else {
		k.SetDelegation(ctx, delegation)
	}
	return amount, nil
}

// bondedTokensToNotBonded transfers coins from the bonded to the not bonded pool within staking
func (k Keeper) bondedTokensToNotBonded(ctx sdk.Context, tokens sdkmath.Int) {
	bondDenom, err := k.BondDenom(ctx)
	if err != nil {
		panic(fmt.Sprintf("unable to get bond denom: %s", err.Error()))
	}
	coins := sdk.NewCoins(sdk.NewCoin(bondDenom, tokens))
	if err := k.bankKeeper.Extend().SendCoinsFromModuleToModuleWithTag(ctx, stakingtypes.BondedPoolName, stakingtypes.NotBondedPoolName, coins, "BondedTokensToNotBonded"); err != nil {
		panic(err)
	}
}

// Delegate performs a delegation, set/update everything necessary within the store.
// tokenSrc indicates the bond status of the incoming funds.
func (k Keeper) Delegate(
	ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdkmath.Int, tokenSrc stakingtypes.BondStatus,
	validator stakingtypes.Validator, delegation stakingtypes.Delegation, valAddr sdk.ValAddress,
) (newShares sdkmath.LegacyDec, err error) {
	// In some situations, the exchange rate becomes invalid, e.g. if
	// Validator loses all tokens due to slashing. In this case,
	// make all future delegations invalid.
	if validator.InvalidExRate() {
		return sdkmath.LegacyZeroDec(), stakingtypes.ErrDelegatorShareExRateInvalid
	}
	if delegation.DelegatorAddress == "" {
		delegation = stakingtypes.NewDelegation(delAddr.String(), valAddr.String(), sdkmath.LegacyZeroDec())
	}
	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
	if tokenSrc == stakingtypes.Bonded {
		panic("delegation token source cannot be bonded")
	}

	var pool string

	switch {
	case validator.IsBonded():
		pool = stakingtypes.BondedPoolName
	case validator.IsUnbonding(), validator.IsUnbonded():
		pool = stakingtypes.NotBondedPoolName
	default:
		panic("invalid validator status")
	}
	denom, _ := k.BondDenom(ctx)
	gage := sdk.NewCoin(denom, bondAmt)
	coins := sdk.NewCoins(gage)
	if err = k.bankKeeper.DelegateCoinsFromAccountToModule(ctx, delegatorAddress, pool, coins); err != nil {
		return sdkmath.LegacyDec{}, err
	}
	poolAccI := k.authKeeper.GetModuleAccount(ctx, pool)
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeDelegateTransfer,
		sdk.NewAttribute(banktypes.AttributeKeySender, delegation.DelegatorAddress),
		sdk.NewAttribute(sdk.AttributeKeyAmount, gage.String()),
		sdk.NewAttribute(banktypes.AttributeKeyRecipient, poolAccI.GetAddress().String()),
	))
	// Update delegation
	if validator.Description.RegionID == types.ExperienceRegionId {
		delegation.UnMeidAmount = delegation.UnMeidAmount.Add(bondAmt)
	} else {
		delegation.Amount = delegation.Amount.Add(bondAmt)
	}
	delegation.StartHeight = ctx.BlockHeight()
	k.SetDelegation(ctx, delegation)

	return sdkmath.LegacyNewDecFromInt(delegation.Amount), nil
}

// WithdrawDelegationRewards withdraw rewards from a delegation
func (k Keeper) WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	regionId := k.GetRegionIdByAccount(ctx, delAddr)
	region, hasRegion := k.GetRegion(ctx, regionId)
	if !hasRegion {
		return nil, fmt.Errorf("%s region not exist", regionId)
	}
	rewards, err := k.internalWithdrawDelegationRewards(ctx, delAddr, region)
	if err != nil {
		return nil, err
	}
	if rewards.IsZero() {
		baseDenom, _ := sdk.GetBaseDenom()
		rewards = sdk.Coins{sdk.Coin{
			Denom:  baseDenom,
			Amount: sdkmath.ZeroInt(),
		}}
	}
	return rewards, nil
}

func (k Keeper) internalWithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, region types.Region) (sdk.Coins, error) {
	//valAddr, valErr := sdk.ValAddressFromBech32(region.OperatorAddress)
	//if valErr != nil {
	//	k.Logger(ctx).Error("internalWithdrawDelegationRewards err=", valErr.Error())
	//	return nil, valErr
	//}
	delegation, err := k.GetDelegation(ctx, delAddr, sdk.ValAddress{})
	if err != nil {
		return nil, fmt.Errorf("delegation not exist")
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
	// finalRewards, remainder := rewards.TruncateDecimal()
	coin := sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())
	coins := sdk.NewCoins(coin)
	// add coins to user account
	if !coin.Amount.IsZero() {
		err = k.bankKeeper.Extend().SendCoinsWithTag(ctx,
			sdk.MustAccAddressFromBech32(region.RegionTreasureAddr),
			sdk.MustAccAddressFromBech32(delegation.DelegatorAddress),
			coins,
			fmt.Sprintf("WithdrawDelegationRewards_%s", region.RegionId))
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

func NewDelegationResp(del stakingtypes.Delegation, balance sdk.Coin) stakingtypes.DelegationResponse {
	return stakingtypes.DelegationResponse{
		Delegation: del,
		Balance:    balance,
	}
}

func (k Keeper) GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	key := stakingtypes.GetDelegationKey(delAddr, sdk.ValAddress{})

	value := store.Get(key)
	if value == nil {
		return stakingtypes.Delegation{}, errorsmod.Wrap(sdkerrors.ErrNotFound, "delegation not found")
	}

	delegation := stakingtypes.MustUnmarshalDelegation(k.cdc, value)

	return delegation, nil
}

func (k Keeper) SetDelegation(ctx sdk.Context, delegation stakingtypes.Delegation) {
	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
	store := ctx.KVStore(k.storeKey)
	b := stakingtypes.MustMarshalDelegation(k.cdc, delegation)
	store.Set(stakingtypes.GetDelegationKey(delegatorAddress, sdk.ValAddress{}), b)
}

func (k Keeper) removeDelegation(ctx sdk.Context, delegation stakingtypes.Delegation) error {
	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
	store := ctx.KVStore(k.storeKey)
	store.Delete(stakingtypes.GetDelegationKey(delegatorAddress, sdk.ValAddress{}))
	return nil
}

func (k *Keeper) SetChangeDelegationValidator(ctx sdk.Context, regionId string) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.ChangeDelegationValidatorKey, []byte(regionId)...)
	store.Set(key, []byte(regionId))
}

func (k *Keeper) DeleteChangeDelegationValidator(ctx sdk.Context, regionId string) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.ChangeDelegationValidatorKey, []byte(regionId)...)
	store.Delete(key)
}

func (k *Keeper) GetAllChangeDelegationValidator(ctx sdk.Context) []string {
	store := ctx.KVStore(k.storeKey)

	iterator := storetypes.KVStorePrefixIterator(store, types.ChangeDelegationValidatorKey)
	defer iterator.Close()

	regionIds := []string{}
	for ; iterator.Valid(); iterator.Next() {
		regionIds = append(regionIds, string(iterator.Value()))
	}
	return regionIds
}

func (k *Keeper) ChangeDelegationValidator(ctx sdk.Context) {
	regionIds := k.GetAllChangeDelegationValidator(ctx)
	for _, regionId := range regionIds {
		region, found := k.GetRegion(ctx, regionId)
		if found {
			k.didKeeper.IteratorCredentialsByFilter(ctx, kyctypes.ModuleName, []byte(regionId), func(vc didtypes.Credential) (stop bool) {
				info, found := k.didKeeper.GetDidInfo(ctx, vc.Did)
				if found {
					delegation, err := k.GetDelegation(ctx, sdk.MustAccAddressFromBech32(info.Address), sdk.ValAddress{})
					if err == nil {
						delegation.ValidatorAddress = region.OperatorAddress
						k.SetDelegation(ctx, delegation)
					}
				}
				return false
			})
			k.DeleteChangeDelegationValidator(ctx, regionId)
		}
	}
}
