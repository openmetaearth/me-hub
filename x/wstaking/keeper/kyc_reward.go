package keeper

import (
	"errors"
	"fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) KycReward(ctx sdk.Context, account sdk.AccAddress, regionId, creator string) error {
	if regionId == strings.ToLower(types.ExperienceRegionName) {
		return sdkerrors.Wrapf(types.ErrSendKycReward, fmt.Sprintf("cannot set kyc to %s region", regionId))
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
		return sdkerrors.Wrapf(types.ErrSendKycReward, fmt.Sprintf("validator reach meid limit"))
	}

	validator.MeidAmount = validator.MeidAmount.Add(types.Bonus)

	err = k.sendKycRewards(ctx, account, valAddr, validator, region)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrSendKycReward, err.Error())
	}

	//validator rewards
	ownerAddress := validator.OwnerAddress
	if len(validator.OwnerAddress) <= 0 {
		ownerAddress = k.daoKeeper.GetDevOperator(ctx)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeNewKyc,
			sdk.NewAttribute(types.AttributeKeyAccount, account.String()),
			sdk.NewAttribute(types.AttributeKeyRegionId, regionId),
			sdk.NewAttribute(types.AttributeKeyCreator, creator),
			sdk.NewAttribute(types.AttributeKeyKycRewardCommitteeAddress, k.daoKeeper.GetDevOperator(ctx)),
			sdk.NewAttribute(types.AttributeKeyKycRewardNodeAddress, ownerAddress),
			sdk.NewAttribute(types.AttributeKeyMeidNumAddReward, types.ValidatorReward.String()+params.BaseDenom),
		),
	)
	return nil
}

func (k Keeper) RemoveKycReward(ctx sdk.Context, account sdk.AccAddress, regionId string) error {
	region, found := k.GetRegion(ctx, regionId)
	if !found {
		return sdkerrors.Wrapf(types.ErrRegionNotExist, fmt.Sprintf("%s not exists", regionId))
	}

	valAddr, err := sdk.ValAddressFromBech32(region.OperatorAddress)
	if err != nil {
		return fmt.Errorf("invalid region operator address")
	}

	validator, ok := k.GetValidator(ctx, valAddr)
	if !ok {
		return fmt.Errorf("region bonded validator not found")
	}

	delegation, found := k.GetDelegation(ctx, account, valAddr)
	if !found {
		return types.ErrNoDelegatorForAddress
	}

	if delegation.Amount.Add(delegation.UnMeidAmount).GT(sdk.ZeroInt()) {
		return types.ErrRemoveKyc.Wrap(fmt.Sprintf("The current user(%s) have delegate, need to withdraw.", account))
	}

	fixedDeposits, err := k.GetFixedDepositByAcct(ctx, account.String())
	if err != nil {
		return err
	}
	if len(fixedDeposits) > 0 {
		return types.ErrRemoveKyc.Wrap(fmt.Sprintf("The current user(%s) have fixed deposit, need to withdraw.", account))
	}

	region.DelegateAmount = region.DelegateAmount.Sub(types.Bonus).Sub(delegation.Amount)
	if region.DelegateAmount.LT(sdk.ZeroInt()) {
		return errors.New("remove kyc error: region delegation amount less than 0")
	}

	rewards, err := k.CalculateInterest(ctx, delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable), delegation.StartHeight)
	if err != nil {
		return types.ErrCalculateInterest.Wrap(err.Error())
	}

	// settle interest
	err = k.bankKeeper.Extend().SendCoinsWithTag(ctx,
		sdk.MustAccAddressFromBech32(region.RegionTreasureAddr),
		account,
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())),
		fmt.Sprintf("RemoveKyc_SettlementInterest_%s", region.RegionId),
	)
	if err != nil {
		return fmt.Errorf("settle interest error: %v", err)
	}

	if region.DelegateInterest.GTE(rewards) {
		region.DelegateInterest = region.DelegateInterest.Sub(rewards)
	}

	if delegation.Unmovable.LTE(sdk.ZeroInt()) {
		return types.ErrDidExists
	}

	delegation.Unmovable = sdk.ZeroInt()
	delegation.StartHeight = ctx.BlockHeight()

	experienceRegion, found := k.GetRegion(ctx, strings.ToLower(types.ExperienceRegionName))
	if !found {
		return types.ErrExperienceRegionNotExist
	}
	delegation.ValidatorAddress = experienceRegion.OperatorAddress
	if delegation.Amount.IsZero() && delegation.UnMeidAmount.IsZero() {
		err = k.removeDelegation(ctx, delegation)
		if err != nil {
			return err
		}
	} else {
		k.SetDelegation(ctx, delegation)
	}
	k.SetRegion(ctx, region)

	validator.DelegationAmount = validator.DelegationAmount.Sub(delegation.Amount)
	validator.MeidAmount = validator.MeidAmount.Sub(types.Bonus)
	k.SetValidator(ctx, validator)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveKyc,
			sdk.NewAttribute(types.AttributeKeyRegionId, region.RegionId),
			sdk.NewAttribute(sdk.AttributeKeySender, region.RegionTreasureAddr),
			sdk.NewAttribute(types.AttributeKeyReceiver, account.String()),
			sdk.NewAttribute(types.AttributeKeyPersonalDelegateInterest, rewards.String()+params.BaseDenom),
		),
	})
	return nil
}

func (k Keeper) sendKycRewards(ctx sdk.Context, delAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	validator stakingtypes.Validator, region types.Region) (err error) {
	experienceRegion, hasRegion := k.GetRegion(ctx, strings.ToLower(types.ExperienceRegionName))
	if !hasRegion {
		return types.ErrExpRegionNotExist
	}

	experienceValAddress, err := sdk.ValAddressFromBech32(experienceRegion.OperatorAddress)
	if err != nil {
		return err
	}

	delegation, found := k.GetDelegation(ctx, delAddr, experienceValAddress)
	if found {
		if delegation.Unmovable.GT(sdk.ZeroInt()) {
			return types.ErrDidExists
		}
		interest, err := k.CalculateInterest(ctx, delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable), delegation.StartHeight)
		if err != nil {
			return types.ErrCalculateInterest.Wrap(err.Error())
		}
		// add coins to user account
		if interest.GT(sdk.ZeroDec()) {
			err = k.bankKeeper.Extend().SendCoinsWithTag(ctx,
				sdk.MustAccAddressFromBech32(experienceRegion.RegionTreasureAddr),
				sdk.MustAccAddressFromBech32(delegation.DelegatorAddress),
				sdk.NewCoins(sdk.NewCoin(params.BaseDenom, interest.TruncateInt())),
				fmt.Sprintf("ApproveKyc_SettlementInterest_%s", region.RegionId),
			)
			if err != nil {
				return err
			}
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
		if experienceVal.DelegationAmount.GTE(delegation.UnMeidAmount) {
			experienceVal.DelegationAmount = experienceVal.DelegationAmount.Sub(delegation.UnMeidAmount)
			k.SetValidator(ctx, experienceVal)
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
	delegation.ValidatorAddress = validator.OperatorAddress
	treasureAddr := k.GetRegionAccount(ctx, types.RegionAccountTypeBase, region.RegionId)

	// validator rewards
	ownerAddress := validator.OwnerAddress
	if len(ownerAddress) == 0 {
		ownerAddress = k.daoKeeper.GetDevOperator(ctx)
	}
	err = k.bankKeeper.Extend().SendCoinsWithTag(ctx,
		treasureAddr.GetAddress(),
		sdk.MustAccAddressFromBech32(ownerAddress),
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, types.ValidatorReward)),
		fmt.Sprintf("ValidatorKycReward_%s", region.RegionId),
	)
	if err != nil {
		return fmt.Errorf("send kyc reward to validator, %v", err)
	}

	//committee rewards
	err = k.bankKeeper.Extend().SendCoinsWithTag(ctx,
		treasureAddr.GetAddress(),
		sdk.MustAccAddressFromBech32(k.daoKeeper.GetDevOperator(ctx)),
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, types.CommitteeReward)),
		fmt.Sprintf("CommitteeKycReward_%s", region.RegionId),
	)
	if err != nil {
		return fmt.Errorf("send kyc reward to committee, %v", err)
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

func (k Keeper) transferDeposit(ctx sdk.Context, fromRegion, toRegion *types.Region, userAddr string) error {
	// GetFixedDepositByAcct returns the list of fixedDeposits of an account
	fixedDeposits, err := k.GetFixedDepositByAcct(ctx, userAddr)
	if err != nil {
		return err
	}
	if len(fixedDeposits) == 0 {
		//if no have deposit，no need to execute the following logic
		return nil
	}

	fromTreasureAddr, accErr := sdk.AccAddressFromBech32(fromRegion.RegionTreasureAddr)
	if accErr != nil {
		return accErr
	}
	fromDepositInterestAddr, accErr := sdk.AccAddressFromBech32(fromRegion.DepositInterestAddr)
	if accErr != nil {
		return accErr
	}
	toTreasureAddr, accErr := sdk.AccAddressFromBech32(toRegion.RegionTreasureAddr)
	if accErr != nil {
		return accErr
	}
	toDepositInterestAddr, accErr := sdk.AccAddressFromBech32(toRegion.DepositInterestAddr)
	if accErr != nil {
		return accErr
	}
	//It is a regional rule used to define parameters such as fixed deposit term and interest rate for a certain region.
	depositConfig := k.GetAllFixedDepositCfg(ctx, toRegion.RegionId)
	depositConfigMap := make(map[int64]sdk.Dec)
	for _, cfg := range depositConfig {
		if cfg.Status == types.RegionFixedDepositCfgStatusInactive {
			return errors.New("fixed deposit cfg status is inactive")
		}
		depositConfigMap[cfg.Term] = cfg.Rate
	}
	totalFixedDepositByAcc := sdk.ZeroInt()
	totalFixedInterestCoin := sdk.ZeroInt()
	for _, fixed := range fixedDeposits {
		totalFixedDepositByAcc = totalFixedDepositByAcc.Add(fixed.Principal.Amount)
		totalFixedInterestCoin = totalFixedInterestCoin.Add(fixed.Interest.Amount)
		//check toRegion deposit config is exist and deposit rate is equal
		rate, exists := depositConfigMap[fixed.Term]
		if !exists || !rate.Equal(fixed.Rate) {
			return errors.New(fmt.Sprintf("deposit cfg not same.rate=%s,fixed.Rate=%s,exists=%v,fixed.Term=%v", rate.String(), fixed.Rate.String(), exists, fixed.Term))
		}

		err := k.IncreaseFixedDepositCountOfCfg(ctx, toRegion.RegionId, fixed.Term)
		if err != nil {
			return err
		}
		err = k.DecreaseFixedDepositCountOfCfg(ctx, fromRegion.RegionId, fixed.Term)
		if err != nil {
			return err
		}
	}
	fromRegion.FixedDepositAmount = fromRegion.FixedDepositAmount.Sub(totalFixedDepositByAcc)
	toRegion.FixedDepositAmount = toRegion.FixedDepositAmount.Add(totalFixedDepositByAcc)
	treasuryBalances := k.bankKeeper.GetBalance(ctx, toTreasureAddr, params.BaseDenom)
	// check toRegion treasury  when subtract original delegation interest,is the balance sufficient.
	if treasuryBalances.Amount.LT(toRegion.DelegateInterest.RoundInt().Add(totalFixedInterestCoin)) {
		return errors.New(fmt.Sprintf("the target region's treasury balance is insufficient,can not pay deposit interest.treasury balance: %s, delegation interest:%s, current user deposit interest:%s",
			treasuryBalances.Amount.String(), toRegion.DelegateInterest.String(),
			totalFixedInterestCoin.String()))
	}
	//pay deposit interest of toRegion
	err = k.bankKeeper.Extend().SendCoinsWithTag(ctx,
		toTreasureAddr,
		toDepositInterestAddr,
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, totalFixedInterestCoin)),
		fmt.Sprintf("TransferFixedInterest_%s", toRegion.RegionId),
	)
	if err != nil {
		return errors.New(fmt.Sprintf("pay deposit interest of toRegion:%s", err.Error()))
	}

	//recovering deposit interest
	err = k.bankKeeper.Extend().SendCoinsWithTag(ctx,
		fromDepositInterestAddr,
		fromTreasureAddr,
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, totalFixedInterestCoin)),
		fmt.Sprintf("RecoverFixedInterest_%s", fromRegion.RegionId),
	)
	if err != nil {
		return errors.New(fmt.Sprintf("recovering deposit interest of fromRegion:%s", err.Error()))
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTransferRegionSettlementFixedDeposit,
			sdk.NewAttribute(types.AttributeKeyTransferAddress, userAddr),
			sdk.NewAttribute(types.AttributeKeyFromRegion, fromRegion.RegionId),
			sdk.NewAttribute(types.AttributeKeyToRegion, toRegion.RegionId),
			sdk.NewAttribute(types.AttributeKeyFixedDeposit, totalFixedInterestCoin.String()+params.BaseDenom),
		),
	})
	return nil
}

func (k Keeper) transferNewMeid(ctx sdk.Context, region *types.Region, address string, valAddr sdk.ValAddress, delegation stakingtypes.Delegation) error {
	accAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return errors.New(fmt.Sprintf("account format error (%s)", err))
	}
	has := k.authKeeper.HasAccount(ctx, accAddr)
	if !has {
		newAccount := k.authKeeper.NewAccountWithAddress(ctx, accAddr)
		k.authKeeper.SetAccount(ctx, newAccount)
	}
	bonus := sdk.NewDec(1).Quo(sdk.NewDecWithPrec(1, params.BaseDenomUnit))
	region.DelegateAmount = region.DelegateAmount.Add(delegation.Amount).Add(bonus.RoundInt())
	delegation.StartHeight = ctx.BlockHeight()
	delegation.ValidatorAddress = valAddr.String()
	k.SetDelegation(ctx, delegation)
	return nil
}

func (k Keeper) transferRemoveMeid(ctx sdk.Context, address string, region *types.Region, delegation stakingtypes.Delegation) error {
	accAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	valAddr, err := sdk.ValAddressFromBech32(region.OperatorAddress)
	if err != nil {
		return err
	}

	validator, ok := k.GetValidator(ctx, valAddr)
	if !ok {
		return errors.New("validator no found")
	}

	_, err = k.transferUnRegisterMeid(ctx, accAddr, region, delegation)
	if err != nil {
		return err
	}

	bonus := sdk.NewDec(1).Quo(sdk.NewDecWithPrec(1, params.BaseDenomUnit))
	validator.MeidAmount = validator.MeidAmount.Sub(bonus.RoundInt())
	validator.DelegationAmount = validator.DelegationAmount.Sub(delegation.Amount)
	k.SetValidator(ctx, validator)
	return nil
}

func (k Keeper) transferUnRegisterMeid(ctx sdk.Context, delAddr sdk.AccAddress, region *types.Region, delegation stakingtypes.Delegation) (amount math.Int, err error) {

	bonus := sdk.NewDec(1).Quo(sdk.NewDecWithPrec(1, params.BaseDenomUnit))
	region.DelegateAmount = region.DelegateAmount.Sub(bonus.RoundInt()).Sub(delegation.Amount)
	if region.DelegateAmount.LT(sdk.ZeroInt()) {
		return amount, errors.New("UnRegisterMeid err: region DelegationAmount < 0")
	}
	rewards, err := k.CalculateInterest(ctx, delegation.Amount.Add(delegation.UnMeidAmount).Add(delegation.Unmovable), delegation.StartHeight)
	if err != nil {
		return amount, err
	}
	regionTreasureAddr, err := sdk.AccAddressFromBech32(region.RegionTreasureAddr)
	if err != nil {
		return amount, err
	}

	if region.DelegateInterest.GTE(rewards) {
		region.DelegateInterest = region.DelegateInterest.Sub(rewards)
	}

	if delegation.Unmovable.LTE(sdk.ZeroInt()) {
		return amount, errors.New("UnRegisterMeid err: delegation UnMovable <= 0")
	}

	err = k.bankKeeper.Extend().SendCoinsWithTag(ctx, regionTreasureAddr, delAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, rewards.TruncateInt())),
		fmt.Sprintf("UpdateKyc_SettlementInterest_%s", region.RegionId),
	)
	if err != nil {
		return amount, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveKyc,
			sdk.NewAttribute(types.AttributeKeyRegionId, region.RegionId),
			sdk.NewAttribute(sdk.AttributeKeySender, region.RegionTreasureAddr),
			sdk.NewAttribute(types.AttributeKeyReceiver, delAddr.String()),
			sdk.NewAttribute(types.AttributeKeyPersonalDelegateInterest, rewards.String()+params.BaseDenom),
			sdk.NewAttribute(types.AttributeKeyExecTime, ctx.BlockTime().String()),
		),
	})
	return amount, nil
}
