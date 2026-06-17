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

	nextDelegationAmount := validator.DelegationAmount.Add(delegation.Amount)
	if validator.Tokens.LT(nextDelegationAmount) {
		return types.ErrNodeLimitExceeded
	}
	nextMeidAmount := validator.MeidAmount.Add(types.Bonus)
	if nextMeidAmount.GT(validator.Tokens) {
		return types.ErrTransferRegion.Wrap(fmt.Sprintf("meid bonded validator can not hold this meid user, reach meid limit"))
	}

	cacheCtx, writeCache := ctx.CacheContext()

	// Handling fixed deposits
	err := k.transferDeposit(cacheCtx, &fromRegion, &toRegion, address.String())
	if err != nil {
		return types.ErrTransferRegion.Wrap(err.Error())
	}

	err = k.transferRemoveMeid(cacheCtx, address.String(), &fromRegion, delegation)
	if err != nil {
		return types.ErrTransferRegion.Wrap(err.Error())
	}

	validator.DelegationAmount = nextDelegationAmount
	validator.MeidAmount = nextMeidAmount
	k.SetValidator(cacheCtx, validator)

	err = k.transferNewMeid(cacheCtx, &toRegion, address.String(), valAddr, delegation)
	if err != nil {
		return types.ErrTransferRegion.Wrap(err.Error())
	}

	k.SetRegion(cacheCtx, fromRegion)
	k.SetRegion(cacheCtx, toRegion)

	cacheCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTransferRegion,
			sdk.NewAttribute(sdk.AttributeKeySender, creator),
			sdk.NewAttribute(types.AttributeKeyTransferAddress, address.String()),
			sdk.NewAttribute(types.AttributeKeyFromRegion, fromRegionId),
			sdk.NewAttribute(types.AttributeKeyToRegion, toRegionId),
			sdk.NewAttribute(types.AttributeKeyRewards, types.Bonus.String()+params.BaseDenom),
		),
	})
	writeCache()
	return nil
}
