package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/openmetaearth/me-hub/utils"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// UpdateValidator defines a method for editing an existing validator
func (k MsgServer) UpdateValidator(goCtx context.Context, msg *types.MsgUpdateValidator) (*types.MsgUpdateValidatorResponse, error) {
	ctx, write := sdk.UnwrapSDKContext(goCtx).CacheContext()

	if !k.daoKeeper.IsGlobalDao(ctx, msg.StakerAddress) {
		return nil, types.ErrCheckGlobalDao
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, err
	}

	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}
	oldRegionId := validator.Description.RegionID
	// replace all editable fields (clients should autofill existing values)
	description, err := validator.Description.UpdateDescription(msg.Description)
	if err != nil {
		return nil, err
	}
	if msg.Description.RegionID == stakingtypes.DoNotModifyDesc {
		description.RegionID = oldRegionId
		validator.Description = description
	} else {
		if _, err := utils.CheckRegionName(strings.ToUpper(msg.Description.RegionID)); err != nil {
			return nil, types.ErrRegionName
		}
		// remove duplication
		validators := k.GetAllValidators(ctx)
		for _, v := range validators {
			if v.Description.RegionID == msg.Description.RegionID {
				return nil, types.ErrValidatorRegionDuplication
			}
		}
		k.UnBondRegion(ctx, oldRegionId)
		description.RegionID = msg.Description.RegionID
		validator.Description = description
		region, f := k.GetRegion(ctx, msg.Description.RegionID)
		if f && region.OperatorAddress == "" {
			k.BondRegion(ctx, validator, validator.Tokens, true)
		}
	}

	//region, f := k.GetRegion(ctx, validator.Description.RegionID)
	// if !f {
	//	return nil, sdkerrors.Wrapf(types.ErrRegionNotExist, "please set region first")
	//}
	// if region.OperatorAddress != validator.OperatorAddress {
	//	return nil, fmt.Errorf("region id already bound to another validator(%s), please set region first", region.OperatorAddress)
	//}

	if msg.CommissionRate != nil {
		commission, err := k.UpdateValidatorCommission(ctx, validator, *msg.CommissionRate)
		if err != nil {
			return nil, err
		}

		// call the before-modification hook since we're about to update the commission
		if err := k.Hooks().BeforeValidatorModified(ctx, valAddr); err != nil {
			return nil, err
		}

		validator.Commission = commission
	}

	if msg.OwnerAddress != "" {
		owner, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
		if err != nil {
			return nil, err
		}
		if err := k.resetValidator(ctx, owner, validator); err != nil {
			return nil, sdkerrors.Wrap(types.ErrResetValidator, err.Error())
		}
	} else {
		k.SetValidator(ctx, validator)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateValidator,
			sdk.NewAttribute(stakingtypes.AttributeKeyCommissionRate, validator.Commission.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyMinSelfDelegation, validator.MinSelfDelegation.String()),
			sdk.NewAttribute(types.AttributeKeyRegionId, validator.Description.RegionID),
		),
	})
	write()
	return &types.MsgUpdateValidatorResponse{}, nil
}
