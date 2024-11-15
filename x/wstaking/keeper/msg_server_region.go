package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/st-chain/me-hub/utils"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

func (k MsgServer) NewRegion(goCtx context.Context, msg *types.MsgNewRegion) (*types.MsgNewRegionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := utils.CheckRegionName(msg.Name)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrRegionName, err.Error())
	}

	if !k.DaoKeeper.IsGlobalDao(ctx, msg.Creator) {
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
	uri := "https://docs.cosmos.network/main/modules/nft"
	hasher := sha256.New()
	_, err = hasher.Write(utils.UnsafeStrToBytes(uri))
	errors.AssertNil(err)
	uriHash := hasher.Sum(nil)

	nftClass := nft.Class{
		Id:          types.GetClassId(msg.Name),
		Name:        types.GetClassName(msg.Name),
		Symbol:      types.GetClassSymbol(msg.Name),
		Description: types.GetClassDescription(regionId),
		Uri:         uri,
		UriHash:     hex.EncodeToString(uriHash),
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
	event4Nft := utils.GenEventCompactAttr(types.EventNewNftClass, nftClass)
	k.SetRegion(ctx, region)
	ctx.EventManager().EmitEvent(event4Nft)
	event4NewRegion := utils.GenEventCompactAttr(types.EventNewRegion, region)
	ctx.EventManager().EmitEvent(event4NewRegion)
	return &types.MsgNewRegionResponse{}, nil
}

func (k MsgServer) RemoveRegion(goCtx context.Context, msg *types.MsgRemoveRegion) (*types.MsgRemoveRegionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.DaoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return nil, types.ErrCheckGlobalDao
	}

	_, found := k.GetRegion(ctx, msg.RegionId)
	if !found {
		return nil, types.ErrRegionNotExist
	}

	err := k.WstakingHooks().BeforeValidatorStakingModified(ctx, sdk.ValAddress{})
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrHooks, "before remove region :error :%+v", err)
	}
	k.Keeper.RemoveRegion(ctx, msg.RegionId)
	return &types.MsgRemoveRegionResponse{}, nil
}

func (k MsgServer) WithdrawFromRegion(goCtx context.Context, msg *types.MsgWithdrawFromRegion) (*types.MsgWithdrawFromRegionResp, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.DaoKeeper.IsGlobalDao(ctx, msg.Withdrawer) {
		return nil, types.ErrCheckGlobalDao
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

	err = k.BankKeeper.SendCoins(
		ctx,
		fromAddr,
		toAddr,
		msg.Amount)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "region %s have enough balance", region.RegionTreasureAddr)
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
