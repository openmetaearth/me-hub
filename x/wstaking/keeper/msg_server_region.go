package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/st-chain/me-hub/utils"
	"github.com/st-chain/me-hub/x/wstaking/types"
	"strings"
)

func (k MsgServer) CurrentDeposit(ctx context.Context, deposit *types.MsgCurrentDeposit) (*types.MsgCurrentDepositResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) CurrentWithdraw(ctx context.Context, withdraw *types.MsgCurrentWithdraw) (*types.MsgCurrentWithdrawResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) FixedDeposit(ctx context.Context, deposit *types.MsgFixedDeposit) (*types.MsgFixedDepositResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) FixedWithdraw(ctx context.Context, withdraw *types.MsgFixedWithdraw) (*types.MsgFixedWithdrawResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) NewFixedDepositCfg(ctx context.Context, cfg *types.MsgFixedDepositCfg) (*types.MsgFixedDepositCfgResp, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) RemoveFixedDepositCfg(ctx context.Context, cfg *types.MsgRemoveFixedDepositCfg) (*types.MsgRemoveFixedDepositCfgResp, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) SetFixedDepositCfgStatus(ctx context.Context, status *types.MsgSetFixedDepositCfgStatus) (*types.MsgSetFixedDepositCfgStatusResp, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) SetFixedDepositCfgRate(ctx context.Context, rate *types.MsgSetFixedDepositCfgRate) (*types.MsgSetFixedDepositCfgRateResp, error) {
	//TODO implement me
	panic("implement me")
}

func (k MsgServer) RemoveRegion(ctx context.Context, region *types.MsgRemoveRegion) (*types.MsgRemoveRegionResponse, error) {
	//TODO implement me
	panic("implement me")
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

func (k MsgServer) NewRegion(goCtx context.Context, msg *types.MsgNewRegion) (*types.MsgNewRegionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, err := utils.CheckRegionName(msg.Name)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrRegionName, err.Error())
	}

	if !k.DaoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return nil, sdkerrors.Wrapf(types.ErrMeidRemove, "only global admin can  create region")
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
		return nil, sdkerrors.Wrapf(types.ErrRegionValidatorNotExist, "meid region bonded validator no found")
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
	//
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
