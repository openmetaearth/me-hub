package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (k Keeper) TransferKycRegion(ctx sdk.Context, address sdk.AccAddress, creator, fromRegionId, toRegionId string) error {
	toRegion, isFound := k.GetRegion(ctx, toRegionId)
	if !isFound {
		return types.ErrRegionNotExist
	}

	valAddr, valErr := sdk.ValAddressFromBech32(toRegion.OperatorAddress)
	if valErr != nil {
		return valErr
	}
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return stakingtypes.ErrNoValidatorFound
	}

	delegation, found := k.GetDelegation(ctx, address, sdk.ValAddress{})
	if !found {
		return types.ErrNoDelegatorForAddress
	}

	if delegation.Amount.Add(delegation.UnMeidAmount).GT(sdk.ZeroInt()) {
		return types.ErrTransferRegion.Wrap(fmt.Sprintf("The current user(%s) have delegate, need to withdrow.", address))
	}

	fixedCount := len(k.GetFixedDepositByAcct(ctx, address.String()))
	if fixedCount > 0 {
		return types.ErrTransferRegion.Wrap(fmt.Sprintf("The current user(%s) have fixed deposit, need to withdrow.", address))
	}

	// fix validator meid amount
	validator.DelegationAmount = validator.DelegationAmount.Add(delegation.Amount)
	if validator.Tokens.LT(validator.DelegationAmount) {
		return types.ErrNodeLimitExceeded
	}
	if validator.MeidAmount.Add(types.Bonus).GT(validator.Tokens) {
		return types.ErrTransferRegion.Wrap(fmt.Sprintf("meid bonded validator can not hold this meid user, reach meid limit"))
	}
	validator.MeidAmount = validator.MeidAmount.Add(types.Bonus)
	k.SetValidator(ctx, validator)

	err := k.RemoveKycReward(ctx, address, fromRegionId)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrDidReward, err.Error())
	}

	err = k.SendKycRewards(ctx, address, valAddr, "", validator, toRegion)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrDidReward, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTransferRegion,
			sdk.NewAttribute(sdk.AttributeKeySender, creator),
			sdk.NewAttribute(types.AttributeKeyTransferAddress, address.String()),
			sdk.NewAttribute(types.AttributeKeyFromRegion, fromRegionId),
			sdk.NewAttribute(types.AttributeKeyToRegion, toRegionId),
			sdk.NewAttribute(types.AttributeKeyRewards, types.Bonus.String()+params.BaseDenom),
		),
	})
	return nil
}
