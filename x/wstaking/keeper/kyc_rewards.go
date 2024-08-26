package keeper

import (
	"context"
	"cosmossdk.io/math"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strings"
)

func (k Keeper) KycReward(ctx sdk.Context, account sdk.AccAddress, inviteAddr, regionId, creator string) error {
	region, found := k.GetRegion(ctx, regionId)
	if !found {
		return types.ErrRegionNotExist
	}

	if regionId == strings.ToLower(types.ExperienceRegion) {
		return sdkerrors.Wrapf(types.ErrTransferRegion, fmt.Sprintf("cannot transfer to %s", regionId))
	}

	valAddr, err := sdk.ValAddressFromBech32(region.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, err.Error())
	}
	validator, ok := k.GetValidator(ctx, valAddr)
	if !ok {
		return types.ErrRegionValidatorNotExist
	}

	bonus := sdk.NewDec(1).Quo(sdk.NewDecWithPrec(1, params.BaseDenomUnit))
	if validator.MeidAmount.Add(bonus.RoundInt()).GT(validator.Tokens) {
		return sdkerrors.Wrapf(types.ErrMeidNew, fmt.Sprintf("meid bonded validator can not hold this meid user, reach meid limit"))
	}

	validator.MeidAmount = validator.MeidAmount.Add(bonus.RoundInt())

	err = k.RegisterKyc(ctx, account, bonus.RoundInt(), valAddr, inviteAddr, validator, region)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrMeidNew, err.Error())
	}

	meid := types.Meid{
		Account:    account.String(),
		Creator:    creator,
		RegionId:   regionId,
		RegionName: region.Name,
		RewardType: types.MeidJoinGroupNoReward,
	}
	k.SetMeid(ctx, meid)
	//validator rewards
	ownerAddress := validator.OwnerAddress
	if len(validator.OwnerAddress) <= 0 {
		ownerAddress = k.DaoKeeper.GetDevOperator(ctx)
	}
	ownerAddr, _ := sdk.AccAddressFromBech32(ownerAddress)
	inviteReward := sdk.NewDec(1).Quo(sdk.NewDecWithPrec(1, params.BaseDenomUnit)).Quo(sdk.NewDec(10)).RoundInt()
	rewards := sdk.NewDec(1).Quo(sdk.NewDecWithPrec(1, params.BaseDenomUnit)).Quo(sdk.NewDec(100)).RoundInt()
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeMeidNew,
			sdk.NewAttribute(types.AttributeKeyAccount, meid.Account),
			sdk.NewAttribute(types.AttributeKeyRegionId, meid.RegionId),
			sdk.NewAttribute(types.AttributeKeyCreator, meid.Creator),
			sdk.NewAttribute(types.AttributeKeyMeidInviteAddress, inviteAddr),
			sdk.NewAttribute(types.AttributeKeyMeidInviteReward, inviteReward.String()+params.BaseDenom),
			sdk.NewAttribute(types.AttributeKeySendMeidInviteAddress, region.RegionTreasureAddr),
			sdk.NewAttribute(types.AttributeKeyReceiveMeidInviteAddress_Society, k.DaoKeeper.GetDevOperator(ctx)),
			sdk.NewAttribute(types.AttributeKeyReceiveMeidInviteAddress_Node, ownerAddr.String()),
			sdk.NewAttribute(types.AttributeKeyMeidNumAddReward, rewards.String()+params.BaseDenom),
		),
	)
	return nil
}

func (k Keeper) RemoveKycReward(goCtx context.Context, account sdk.AccAddress) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	meid, found := k.GetMeid(ctx, account.String())
	if !found {
		return types.ErrMeidNotExists
	}

	region, found := k.GetRegion(ctx, meid.RegionId)
	if !found {
		return sdkerrors.Wrapf(types.ErrMeidNotExists, fmt.Sprintf("meid's region %s not exists", meid.RegionId))
	}

	valAddr, err := sdk.ValAddressFromBech32(region.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Sprintf("meid region bonded validator no found"))
	}

	validator, ok := k.GetValidator(ctx, valAddr)
	if !ok {
		return sdkerrors.Wrapf(types.ErrMeidNew, fmt.Sprintf("meid region bonded validator no found"))
	}

	_, err = k.UnRegisterMeid(ctx, account, valAddr, region, stakingtypes.Delegation{})
	if err != nil {
		return sdkerrors.Wrapf(types.ErrMeidNew, fmt.Sprintf("unregister meid airdop failed"))
	}

	bonus := sdk.NewDec(1).Quo(sdk.NewDecWithPrec(1, params.BaseDenomUnit))
	validator.MeidAmount = validator.MeidAmount.Sub(bonus.RoundInt())
	k.SetValidator(ctx, validator)

	k.RemoveMeid(ctx, meid.Account, region.RegionId)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeMeidRemove,
			sdk.NewAttribute(types.AttributeKeyAccount, meid.Account),
			sdk.NewAttribute(types.AttributeKeyRegionId, meid.RegionId),
			sdk.NewAttribute(types.AttributeKeyCreator, meid.Creator),
		),
	)
	return nil
}

func (k Keeper) RegisterKyc(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int,
	validatorAddr sdk.ValAddress, inviteAddr string, validator stakingtypes.Validator, region types.Region) (err error) {
	experienceRegion, hasRegion := k.GetRegion(ctx, strings.ToLower(types.ExperienceRegion))
	if !hasRegion {
		return types.ErrExpRegionNotExist
	}

	experienceValAddress, err := sdk.ValAddressFromBech32(experienceRegion.OperatorAddress)
	if err != nil {
		return err
	}

	delegation, found := k.GetDelegation(ctx, delAddr, sdk.ValAddress{})
	if !found {
		delegation = stakingtypes.NewDelegation(delAddr, validatorAddr, sdk.ZeroDec())
	}

	if delegation.Unmovable.GT(sdk.ZeroInt()) {
		return types.ErrMeidExists
	}

	// call the appropriate hook if present
	if found {
		rewards, err := k.CalculateInterest(ctx, delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable), delegation.StartHeight)
		if err != nil {
			return types.ErrCalculateInterest.Wrap(err.Error())
		}

		// add coins to user account
		if rewards.GT(sdk.ZeroDec()) {
			err = k.BankKeeper.SendCoins(ctx,
				sdk.MustAccAddressFromBech32(region.RegionTreasureAddr),
				sdk.MustAccAddressFromBech32(delegation.DelegatorAddress),
				sdk.NewCoins(sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())))
			if err != nil {
				return err
			}
			if experienceRegion.DelegateInterest.GTE(rewards) {
				experienceRegion.DelegateInterest = experienceRegion.DelegateInterest.Sub(rewards)
			}
			experienceRegion.DelegateAmount = experienceRegion.DelegateAmount.Sub(delegation.UnMeidAmount)
			k.SetRegion(ctx, experienceRegion)

			experienceVal, ok := k.GetValidator(ctx, experienceValAddress)
			if !ok {
				return fmt.Errorf("experience region validator no found")
			}
			if experienceVal.DelegationAmount.GTE(delegation.Amount) {
				experienceVal.DelegationAmount = experienceVal.DelegationAmount.Sub(delegation.Amount)
				k.SetValidator(ctx, experienceVal)
			}
		}
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeSettleDelRewardsForKyc,
				sdk.NewAttribute(types.AttributeKeyValidator, experienceRegion.OperatorAddress),
				sdk.NewAttribute(types.AttributeKeyDelegator, delegation.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeyRegionTreasure, experienceRegion.RegionTreasureAddr),
				sdk.NewAttribute(types.AttributeKeyRegionId, experienceRegion.RegionId),
				sdk.NewAttribute(types.AttributeKeyAmountDelegateInterest, experienceRegion.DelegateInterest.String()+params.BaseDenom),
				sdk.NewAttribute(types.AttributeKeyPersonalDelegateInterest, rewards.TruncateInt().String()+params.BaseDenom),
			),
		})
	}

	// Update delegation
	delegation.Unmovable = bondAmt //delegation.Unmovable.Add(bondAmt)
	delegation.StartHeight = ctx.BlockHeight()
	delegation.ValidatorAddress = validatorAddr.String()
	treasureAddr := k.GetRegionAccount(ctx, types.RegionAccountTypeBase, region.RegionId)
	if len(inviteAddr) > 0 {
		addr, err := sdk.AccAddressFromBech32(inviteAddr)
		if err != nil {
			return err
		}
		inviteReward := sdk.NewDec(1).Quo(sdk.NewDecWithPrec(1, params.BaseDenomUnit)).Quo(sdk.NewDec(10)).RoundInt()
		err = k.BankKeeper.SendCoins(ctx, treasureAddr.GetAddress(), addr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, inviteReward)))
		if err != nil {
			return err
		}
	}

	rewards := sdk.NewDec(1).Quo(sdk.NewDecWithPrec(1, params.BaseDenomUnit)).Quo(sdk.NewDec(100)).RoundInt()
	//validator rewards
	if len(validator.OwnerAddress) <= 0 {
		validator.OwnerAddress = k.DaoKeeper.GetDevOperator(ctx)
	}

	ownerAddr, err := sdk.AccAddressFromBech32(validator.OwnerAddress)
	if err != nil {
		return err
	}

	err = k.BankKeeper.SendCoins(ctx, treasureAddr.GetAddress(), ownerAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, rewards)))
	if err != nil {
		return err
	}

	//committee rewards
	err = k.BankKeeper.SendCoins(ctx, treasureAddr.GetAddress(),
		sdk.MustAccAddressFromBech32(k.DaoKeeper.GetDevOperator(ctx)),
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, rewards)))
	if err != nil {
		return err
	}

	delegation.Amount = delegation.Amount.Add(delegation.UnMeidAmount)
	delegation.UnMeidAmount = sdk.ZeroInt()
	k.SetDelegation(ctx, delegation)

	bonus := sdk.NewDec(1).Quo(sdk.NewDecWithPrec(1, params.BaseDenomUnit))
	region.DelegateAmount = region.DelegateAmount.Add(delegation.Amount).Add(bonus.RoundInt())
	k.SetRegion(ctx, region)

	validator.DelegationAmount = validator.DelegationAmount.Add(delegation.Amount)
	k.SetValidator(ctx, validator)
	return nil
}

func (k Keeper) UnRegisterMeid(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, region types.Region, delegation stakingtypes.Delegation) (amount math.Int, err error) {
	if delegation == (stakingtypes.Delegation{}) {
		var found bool
		// check if a delegation object exists in the store
		delegation, found = k.GetDelegation(ctx, delAddr, valAddr)
		if !found {
			return amount, stakingtypes.ErrNoDelegatorForAddress
		}
	}

	bonus := sdk.NewDec(1).Quo(sdk.NewDecWithPrec(1, params.BaseDenomUnit))
	region.DelegateAmount = region.DelegateAmount.Sub(bonus.RoundInt())
	if region.DelegateAmount.LT(sdk.ZeroInt()) {
		return amount, errors.New("UnRegisterMeid err: region DelegationAmount < 0")
	}
	rewards, err := k.CalculateInterest(ctx, delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable), delegation.StartHeight)
	if err != nil {
		return amount, types.ErrCalculateInterest.Wrap(err.Error())
	}
	regionTreasureAddr, err := sdk.AccAddressFromBech32(region.RegionTreasureAddr)
	if err != nil {
		return amount, err
	}
	err = k.BankKeeper.SendCoins(ctx, regionTreasureAddr, delAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())))
	if err != nil {
		return amount, err
	}
	if region.DelegateInterest.GTE(rewards) {
		region.DelegateInterest = region.DelegateInterest.Sub(rewards)
	}
	if delegation.Unmovable.LTE(sdk.ZeroInt()) {
		return amount, types.ErrMeidExists
	}
	delegation.Unmovable = sdk.ZeroInt()
	delegation.StartHeight = ctx.BlockHeight()
	valaddrStr := k.AuthKeeper.GetModuleAccount(ctx, stakingtypes.NotBondedPoolName).GetAddress().String()
	valStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), []byte(valaddrStr))
	if err != nil {
		return amount, err
	}
	delegation.ValidatorAddress = valStr
	if delegation.Amount.IsZero() && delegation.UnMeidAmount.IsZero() {
		err = k.RemoveDelegation(ctx, delegation)
		if err != nil {
			return amount, err
		}
	} else {
		k.SetDelegation(ctx, delegation)
	}
	k.SetRegion(ctx, region)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMeidRemove,
			sdk.NewAttribute(types.AttributeKeyRegionId, region.RegionId),
			sdk.NewAttribute(sdk.AttributeKeySender, region.RegionTreasureAddr),
			sdk.NewAttribute(types.AttributeKeyReceiver, delAddr.String()),
			sdk.NewAttribute(types.AttributeKeyPersonalDelegateInterest, rewards.String()+params.BaseDenom),
			sdk.NewAttribute(types.AttributeKeyExecTime, ctx.BlockTime().String()),
		),
	})
	return amount, nil
}
