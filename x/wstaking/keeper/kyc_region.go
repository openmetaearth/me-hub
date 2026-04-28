package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"strings"
)

func (k Keeper) GetRegionIdByAccount(ctx sdk.Context, address sdk.AccAddress) string {
	regionId := strings.ToLower(types.ExperienceRegionName)
	did, ok := k.kycKeeper.GetDID(ctx, address)
	if !ok {
		return regionId
	}
	kycData, ok := k.kycKeeper.GetKYC(ctx, did)
	if !ok {
		return regionId
	}
	return string(kycData.Data)
}

func (k Keeper) MustGetKycRegionIdByAccount(ctx sdk.Context, account string) (string, error) {
	did, ok := k.kycKeeper.GetDID(ctx, sdk.MustAccAddressFromBech32(account))
	if !ok {
		return "", sdkerrors.Wrapf(types.ErrDidNotExists, "did with account %s not exist", account)
	}
	kycData, ok := k.kycKeeper.GetKYC(ctx, did)
	if !ok {
		return "", sdkerrors.Wrapf(types.ErrKycNotExists, "kyc with account %s not exist", account)
	}
	return string(kycData.Data), nil
}

func (k Keeper) TransferKycRegion(ctx sdk.Context, address sdk.AccAddress, creator, fromRegionId, toRegionId string) error {

	fromRegion, found := k.GetRegion(ctx, fromRegionId)
	if !found {
		return types.ErrRegionNotExist
	}

	fromValAddr, valErr := sdk.ValAddressFromBech32(fromRegion.OperatorAddress)
	if valErr != nil {
		return valErr
	}

	toRegion, found := k.GetRegion(ctx, toRegionId)
	if !found {
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

	delegation, found := k.GetDelegation(ctx, address, fromValAddr)
	if !found {
		return types.ErrNoDelegatorForAddress
	}
	delegation.ValidatorAddress = toRegion.OperatorAddress
	k.SetDelegation(ctx, delegation)

	// Handling fixed deposits
	err := k.transferDeposit(ctx, &fromRegion, &toRegion, address.String())
	if err != nil {
		return types.ErrTransferRegion.Wrap(err.Error())
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

	err = k.transferRemoveMeid(ctx, address.String(), &fromRegion, delegation)
	if err != nil {
		return types.ErrTransferRegion.Wrap(err.Error())
	}

	err = k.transferNewMeid(ctx, &toRegion, address.String(), valAddr, delegation)
	if err != nil {
		return types.ErrTransferRegion.Wrap(err.Error())
	}

	k.SetRegion(ctx, fromRegion)
	k.SetRegion(ctx, toRegion)

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
