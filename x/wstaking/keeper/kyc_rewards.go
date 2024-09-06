package keeper

import (
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
	if regionId == strings.ToLower(types.ExperienceRegion) {
		return sdkerrors.Wrapf(types.ErrTransferRegion, fmt.Sprintf("cannot transfer to %s", regionId))
	}

	region, found := k.GetRegion(ctx, regionId)
	if !found {
		return types.ErrRegionNotExist
	}

	valAddr, err := sdk.ValAddressFromBech32(region.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, err.Error())
	}
	validator, ok := k.GetValidator(ctx, valAddr)
	if !ok {
		return types.ErrRegionValidatorNotExist
	}

	if validator.MeidAmount.Add(types.Bonus).GT(validator.Tokens) {
		return sdkerrors.Wrapf(types.ErrKycReward, fmt.Sprintf("meid bonded validator can not hold this meid user, reach meid limit"))
	}

	validator.MeidAmount = validator.MeidAmount.Add(types.Bonus)

	err = k.SendKycRewards(ctx, account, valAddr, inviteAddr, validator, region)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrKycReward, err.Error())
	}

	//validator rewards
	ownerAddress := validator.OwnerAddress
	if len(validator.OwnerAddress) <= 0 {
		ownerAddress = k.DaoKeeper.GetDevOperator(ctx)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeMeidNew,
			sdk.NewAttribute(types.AttributeKeyAccount, account.String()),
			sdk.NewAttribute(types.AttributeKeyRegionId, regionId),
			sdk.NewAttribute(types.AttributeKeyCreator, creator),
			sdk.NewAttribute(types.AttributeKeyMeidInviteAddress, inviteAddr),
			sdk.NewAttribute(types.AttributeKeyMeidInviteReward, types.InviteReward.String()+params.BaseDenom),
			sdk.NewAttribute(types.AttributeKeySendMeidInviteAddress, region.RegionTreasureAddr),
			sdk.NewAttribute(types.AttributeKeyReceiveMeidInviteAddress_Society, k.DaoKeeper.GetDevOperator(ctx)),
			sdk.NewAttribute(types.AttributeKeyReceiveMeidInviteAddress_Node, ownerAddress),
			sdk.NewAttribute(types.AttributeKeyMeidNumAddReward, types.ValidatorReward.String()+params.BaseDenom),
		),
	)
	return nil
}

func (k Keeper) RemoveKycReward(ctx sdk.Context, account sdk.AccAddress, regionId string) error {
	region, found := k.GetRegion(ctx, regionId)
	if !found {
		return sdkerrors.Wrapf(types.ErrMeidNotExists, fmt.Sprintf("meid's region %s not exists", regionId))
	}

	valAddr, err := sdk.ValAddressFromBech32(region.OperatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Sprintf("region bonded validator no found"))
	}

	validator, ok := k.GetValidator(ctx, valAddr)
	if !ok {
		return sdkerrors.Wrapf(types.ErrKycReward, fmt.Sprintf("region bonded validator no found"))
	}

	_, err = k.removeKycReward(ctx, account, valAddr, region, stakingtypes.Delegation{})
	if err != nil {
		return sdkerrors.Wrapf(types.ErrKycReward, fmt.Sprintf("remove kyc reward failed"))
	}

	validator.MeidAmount = validator.MeidAmount.Sub(types.Bonus)
	k.SetValidator(ctx, validator)
	return nil
}

func (k Keeper) SendKycRewards(ctx sdk.Context, delAddr sdk.AccAddress,
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
	if found {
		if delegation.Unmovable.GT(sdk.ZeroInt()) {
			return types.ErrMeidExists
		}
		interest, err := k.CalculateInterest(ctx, delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable), delegation.StartHeight)
		if err != nil {
			return types.ErrCalculateInterest.Wrap(err.Error())
		}

		// add coins to user account
		if interest.GT(sdk.ZeroDec()) {
			err = k.BankKeeper.SendCoins(ctx,
				sdk.MustAccAddressFromBech32(region.RegionTreasureAddr),
				sdk.MustAccAddressFromBech32(delegation.DelegatorAddress),
				sdk.NewCoins(sdk.NewCoin(params.BaseDenom, interest.TruncateInt())))
			if err != nil {
				return err
			}
			if experienceRegion.DelegateInterest.GTE(interest) {
				experienceRegion.DelegateInterest = experienceRegion.DelegateInterest.Sub(interest)
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
				sdk.NewAttribute(types.AttributeKeyPersonalDelegateInterest, interest.TruncateInt().String()+params.BaseDenom),
			),
		})
	} else {
		delegation = stakingtypes.NewDelegation(delAddr, validatorAddr, sdk.ZeroDec())
	}

	// Update delegation
	delegation.Unmovable = types.Bonus //delegation.Unmovable.Add(bondAmt)
	delegation.StartHeight = ctx.BlockHeight()
	delegation.ValidatorAddress = validatorAddr.String()
	treasureAddr := k.GetRegionAccount(ctx, types.RegionAccountTypeBase, region.RegionId)
	if len(inviteAddr) > 0 {
		addr, err := sdk.AccAddressFromBech32(inviteAddr)
		if err != nil {
			return err
		}
		err = k.BankKeeper.SendCoins(ctx, treasureAddr.GetAddress(), addr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, types.InviteReward)))
		if err != nil {
			return err
		}
	}

	// validator rewards
	ownerAddress := validator.OwnerAddress
	if len(ownerAddress) == 0 {
		ownerAddress = k.DaoKeeper.GetDevOperator(ctx)
	}
	err = k.BankKeeper.SendCoins(ctx,
		treasureAddr.GetAddress(),
		sdk.MustAccAddressFromBech32(ownerAddress),
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, types.ValidatorReward)))
	if err != nil {
		return err
	}

	//committee rewards
	err = k.BankKeeper.SendCoins(ctx,
		treasureAddr.GetAddress(),
		sdk.MustAccAddressFromBech32(k.DaoKeeper.GetDevOperator(ctx)),
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, types.CommitteeReward)))
	if err != nil {
		return err
	}

	delegation.Amount = delegation.Amount.Add(delegation.UnMeidAmount)
	delegation.UnMeidAmount = sdk.ZeroInt()
	k.SetDelegation(ctx, delegation)

	region.DelegateAmount = region.DelegateAmount.Add(delegation.Amount).Add(types.Bonus)
	k.SetRegion(ctx, region)

	validator.DelegationAmount = validator.DelegationAmount.Add(delegation.Amount)
	k.SetValidator(ctx, validator)
	return nil
}

func (k Keeper) removeKycReward(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, region types.Region, delegation stakingtypes.Delegation) (amount math.Int, err error) {
	if delegation == (stakingtypes.Delegation{}) {
		var found bool
		// check if a delegation object exists in the store
		delegation, found = k.GetDelegation(ctx, delAddr, valAddr)
		if !found {
			return amount, stakingtypes.ErrNoDelegatorForAddress
		}
	}

	region.DelegateAmount = region.DelegateAmount.Sub(types.Bonus)
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

	experienceRegion, _ := k.GetRegion(ctx, strings.ToLower(types.ExperienceRegion))
	delegation.ValidatorAddress = experienceRegion.OperatorAddress
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
		),
	})
	return amount, nil
}
