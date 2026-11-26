package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/x/nft"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/utils"
	wnfttypes "github.com/openmetaearth/me-hub/x/wnft/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k MsgServer) NewRegion(goCtx context.Context, msg *types.MsgNewRegion) (*types.MsgNewRegionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := utils.CheckRegionName(msg.Name)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrRegionName, err.Error())
	}

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return nil, types.ErrCheckGlobalDao
	}

	regionId := strings.ToLower(msg.Name)
	_, found := k.GetRegion(ctx, regionId)
	if found {
		return nil, errorsmod.Wrapf(types.ErrRegionAlreadyExist, "region already exist")
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "region bonded validator no found")
	}

	validator, err := k.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, types.ErrRegionValidatorNotExist
	}
	if !strings.EqualFold(validator.Description.RegionID, regionId) {
		return nil, types.ErrRegion.Wrapf("only the validator with region id %s can be bound, not bound %s region", validator.Description.RegionID, regionId)
	}

	allRegions := k.Keeper.GetAllRegion(ctx)
	for _, reg := range allRegions {
		if reg.OperatorAddress == msg.OperatorAddress {
			return nil, errorsmod.Wrapf(types.ErrRegionValidatorDuplicate, "meid region bonded validator duplicates")
		}
		if reg.RegionId == regionId {
			return nil, errorsmod.Wrapf(types.ErrRegionNameDuplicate, "meid region name duplicates")
		}
	}
	err = k.WstakingHooks().BeforeValidatorStakingModified(ctx, valAddr)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrHooks, "before create new region :error :%+v", err)
	}

	uri := ""
	classMetadata := &wnfttypes.ClassMetadata{
		Creator: msg.Creator,
	}
	metadata, err := codectypes.NewAnyWithValue(classMetadata)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrLogic, "%v", err)
	}
	nftClass := nft.Class{
		Id:          types.GetClassId(msg.Name),
		Name:        types.GetClassName(msg.Name),
		Symbol:      types.GetClassSymbol(msg.Name),
		Description: types.GetClassDescription(regionId),
		Uri:         uri,
		UriHash:     utils.CalculateURIHash(uri),
		Data:        metadata,
	}

	_, nftClassFound := k.nftKeeper.GetClass(ctx, nftClass.Id)
	if !nftClassFound {
		err = k.nftKeeper.SaveClass(ctx, nftClass)
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrRegionAlreadyExist, "save nft class: %v", err)
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
	// create megroup
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
	// RemoveRegion is currently a no-op, matching main.
	return &types.MsgRemoveRegionResponse{}, nil
}

func (k MsgServer) TransferRegion(goCtx context.Context, msg *types.MsgTransferRegion) (*types.MsgTransferRegionResponse, error) {
	return &types.MsgTransferRegionResponse{}, nil
}
