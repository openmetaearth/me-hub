package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/st-chain/me-hub/utils"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strings"
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
		return nil, sdkerrors.Wrapf(types.ErrRegionValidatorNotExist, "region bonded validator no found")
	}
	if strings.ToLower(validator.Description.RegionId) != strings.ToLower(regionId) {
		return nil, types.ErrRegion.Wrapf("only the validator with region id  %s can be bound,not bound %s region id", validator.Description.RegionId, msg.RegionId)
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

	//uri := "https://docs.cosmos.network/main/modules/nft"
	//hasher := sha256.New()
	//_, err = hasher.Write(conv.UnsafeStrToBytes(uri))
	//errors.AssertNil(err)
	//uriHash := hasher.Sum(nil)
	//
	//ntfClassId := msg.Name + "-NFT-CLASS-ID-"
	//nftClass := nft.Class{
	//	Id:          ntfClassId,
	//	Name:        msg.Name + "-NFT-CLASS-NAME",
	//	Symbol:      msg.Name + "-NFT-CLASS-SYMBOL",
	//	Description: "nft class for region: " + msg.RegionId,
	//	Uri:         uri,
	//	UriHash:     string(uriHash[:]),
	//}
	//err = k.nftKeeper.SaveClass(ctx, nftClass)
	//if err != nil {
	//	return nil, sdkerrors.Wrapf(types.ErrRegionAlreadyExist, "nft classe save error")
	//}

	region := types.Region{
		RegionId:        msg.RegionId,
		Creator:         msg.Creator,
		Name:            msg.Name,
		OperatorAddress: msg.OperatorAddress,
		//NftClassId:          ntfClassId,
		RegionTreasureAddr:  k.CreateRegionAccount(ctx, types.RegionAccountTypeBase, msg.RegionId).String(),
		DepositInterestAddr: k.CreateRegionAccount(ctx, types.RegionAccountTypeDepositInterest, msg.RegionId).String(),
		RegionShare:         validator.Tokens,
	}
	if msg.RegionId == strings.ToLower(types.ExperienceRegion) {
		region.DepositInterestAddr = ""
	}
	k.SetRegion(ctx, region)
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

	k.Keeper.RemoveRegion(ctx, msg.RegionId)
	return &types.MsgRemoveRegionResponse{}, nil
}

func (k MsgServer) RetrieveCoinsFromRegion(ctx context.Context, region *types.MsgRetrieveCoinsFromRegion) (*types.MsgRetrieveCoinsFromRegionResp, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) TransferRegion(ctx context.Context, region *types.MsgTransferRegion) (*types.MsgTransferRegionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) RetrieveFeeFromGlobalAdminFeePool(ctx context.Context, pool *types.MsgRetrieveFeeFromGlobalAdminFeePool) (*types.MsgRetrieveFeeFromGlobalAdminFeePoolResp, error) {
	//TODO implement me
	panic("implement me")
}
