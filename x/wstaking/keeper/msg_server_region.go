package keeper

import (
	"context"
	"fmt"
	"strings"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	wnfttypes "github.com/openmetaearth/me-hub/x/wnft/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/openmetaearth/me-hub/utils"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k MsgServer) NewRegion(goCtx context.Context, msg *types.MsgNewRegion) (*types.MsgNewRegionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := utils.CheckRegionName(msg.Name)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrRegionName, err.Error())
	}

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return nil, types.ErrCheckGlobalDao
	}

	regionId := strings.ToLower(msg.Name)
	_, found := k.GetRegion(ctx, regionId)
	if found {
		return nil, sdkerrors.Wrapf(types.ErrRegionAlreadyExist, "region already exist")
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "region bonded validator no found")
	}

	validator, ok := k.GetValidator(ctx, valAddr)
	if !ok {
		return nil, types.ErrRegionValidatorNotExist
	}
	if strings.ToLower(validator.Description.RegionID) != strings.ToLower(regionId) {
		return nil, types.ErrRegion.Wrapf("only the validator with region id %s can be bound, not bound %s region", validator.Description.RegionID, regionId)
	}

	allRegions := k.Keeper.GetAllRegion(ctx)
	for _, reg := range allRegions {
		if reg.OperatorAddress == msg.OperatorAddress {
			return nil, sdkerrors.Wrapf(types.ErrRegionValidatorDuplicate, "meid region bonded validator duplicates")
		}
		if reg.RegionId == regionId {
			return nil, sdkerrors.Wrapf(types.ErrRegionNameDuplicate, "meid region name duplicates")
		}
	}
	err = k.WstakingHooks().BeforeValidatorStakingModified(ctx, valAddr)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrHooks, "before create new region :error :%+v", err)
	}

	uri := ""
	classMetadata := &wnfttypes.ClassMetadata{
		Creator: msg.Creator,
	}
	metadata, err := codectypes.NewAnyWithValue(classMetadata)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrLogic, "%v", err)
	}
	nftClass := nft.Class{
		Id:          types.GetClassId(msg.Name),
		Name:        types.GetClassName(msg.Name),
		Symbol:      types.GetClassSymbol(msg.Name),
		Description: types.GetClassDescription(regionId),
		Uri:         uri,
		UriHash:     utils.CalculateUriHash(uri),
		Data:        metadata,
	}

	_, nftClassFound := k.nftKeeper.GetClass(ctx, nftClass.Id)
	if !nftClassFound {
		err = k.nftKeeper.SaveClass(ctx, nftClass)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrRegionAlreadyExist, "save nft class: %v", err)
		}
	}

	region := types.Region{
		RegionId:            regionId,
		Creator:             msg.Creator,
		Name:                msg.Name,
		OperatorAddress:     msg.OperatorAddress,
		NftClassId:          types.GetClassId(msg.Name),
		RegionTreasureAddr:  k.CreateRegionAccount(ctx, types.RegionAccountTypeBase, regionId).String(),
		DepositInterestAddr: k.CreateRegionAccount(ctx, types.RegionAccountTypeDepositInterest, regionId).String(),
		RegionShare:         validator.Tokens,
	}
	if regionId == strings.ToLower(types.ExperienceRegionName) {
		region.DepositInterestAddr = ""
	}

	event4Nft := utils.GenEventCompactAttrWithBytes(types.EventNewNftClass, k.cdc.MustMarshal(&nftClass))
	k.SetRegion(ctx, region)
	//create megroup
	if regionId != strings.ToLower(types.ExperienceRegionName) {
		if _, err := k.groupKeeper.CreateGroupByRegion(ctx, region); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(event4Nft)
	event4NewRegion := utils.GenEventCompactAttr(types.EventNewRegion, region)
	ctx.EventManager().EmitEvent(event4NewRegion)
	return &types.MsgNewRegionResponse{}, nil
}

func (k MsgServer) RemoveRegion(goCtx context.Context, msg *types.MsgRemoveRegion) (*types.MsgRemoveRegionResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)

	// if !k.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
	// 	return nil, types.ErrCheckGlobalDao
	// }

	// _, found := k.GetRegion(ctx, msg.RegionId)
	// if !found {
	// 	return nil, types.ErrRegionNotExist
	// }

	// err := k.WstakingHooks().BeforeValidatorStakingModified(ctx, sdk.ValAddress{})
	// if err != nil {
	// 	return nil, sdkerrors.Wrapf(types.ErrHooks, "before remove region :error :%+v", err)
	// }
	// k.Keeper.RemoveRegion(ctx, msg.RegionId)
	// ctx.EventManager().EmitEvent(
	// 	sdk.NewEvent(
	// 		types.EventTypeRemoveRegion,
	// 		sdk.NewAttribute(types.AttributeKeyRegionId, msg.RegionId),
	// 	),
	// )
	return &types.MsgRemoveRegionResponse{}, nil
}

func (k MsgServer) WithdrawFromRegion(goCtx context.Context, msg *types.MsgWithdrawFromRegion) (*types.MsgWithdrawFromRegionResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isDao := k.daoKeeper.IsGlobalDao(ctx, msg.Withdrawer)
	isGranted := k.CanRegionWithdraw(ctx, msg.Withdrawer, msg.RegionId)
	if !isDao && !isGranted {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized,
			"address %s can not withdraw for region %s", msg.Withdrawer, msg.RegionId)
	}

	region, found := k.GetRegion(ctx, msg.RegionId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrRegionNotExist, "region not exist")
	}

	fromAddr, err := sdk.AccAddressFromBech32(region.RegionTreasureAddr)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownAccount, "region account %s format error %s", region.RegionTreasureAddr, err)
	}

	toAddr, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownAccount, "receiver account %s format error %s", msg.Receiver, err)
	}

	err = k.bankKeeper.Extend().SendCoinsWithTag(
		ctx,
		fromAddr,
		toAddr,
		msg.Amount,
		fmt.Sprintf("WithdrawFromRegionTreasure_%s", region.RegionId),
	)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "region treasure %s does not have enough balance", region.RegionTreasureAddr)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeWithdrawFromRegion,
			sdk.NewAttribute(types.AttributeKeyRegionId, msg.RegionId),
			sdk.NewAttribute(sdk.AttributeKeySender, fromAddr.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, toAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	)
	return &types.MsgWithdrawFromRegionResp{}, nil
}

func (k MsgServer) TransferRegion(goCtx context.Context, msg *types.MsgTransferRegion) (*types.MsgTransferRegionResponse, error) {
	return &types.MsgTransferRegionResponse{}, nil
}

// GrantRegionWithdraw grants (or overwrites) a address for who can withdraw from the region treasury.
// Only GlobalDao can call this. One region maps to exactly one address.
func (k MsgServer) GrantRegionWithdraw(goCtx context.Context, msg *types.MsgGrantRegionWithdraw) (*types.MsgGrantRegionWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return nil, types.ErrCheckGlobalDao
	}

	if _, found := k.GetRegion(ctx, msg.RegionId); !found {
		return nil, sdkerrors.Wrapf(types.ErrRegionNotExist, "region %s not found", msg.RegionId)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", err)
	}

	k.SetRegionWithdraw(ctx, msg.RegionId, msg.Address)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGrantRegionWithdraw,
		sdk.NewAttribute(types.AttributeKeyRegionId, msg.RegionId),
		sdk.NewAttribute(types.AttributeKeyGrantedAddress, msg.Address),
	))

	return &types.MsgGrantRegionWithdrawResponse{}, nil
}

// RevokeRegionWithdraw removes the withdraw address for a region.
// Only GlobalDao can call this.
func (k MsgServer) RevokeRegionWithdraw(goCtx context.Context, msg *types.MsgRevokeRegionWithdraw) (*types.MsgRevokeRegionWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return nil, types.ErrCheckGlobalDao
	}

	if _, found := k.GetRegion(ctx, msg.RegionId); !found {
		return nil, sdkerrors.Wrapf(types.ErrRegionNotExist, "region %s not found", msg.RegionId)
	}

	if _, found := k.GetRegionWithdraw(ctx, msg.RegionId); !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound,
			"no withdraw address found for region %s", msg.RegionId)
	}

	k.DeleteRegionWithdraw(ctx, msg.RegionId)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRevokeRegionWithdraw,
		sdk.NewAttribute(types.AttributeKeyRegionId, msg.RegionId),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
	))

	return &types.MsgRevokeRegionWithdrawResponse{}, nil
}
