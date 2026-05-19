package keeper

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/openmetaearth/me-hub/app/params"
	didTypes "github.com/openmetaearth/me-hub/x/did/types"
	kycTypes "github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/openmetaearth/me-hub/x/megroup/types"
	stakingTypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

type kycHookFunc func(ctx sdk.Context, eventType string, beforeData interface{}, afterData interface{}) error

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		paramstore paramtypes.Subspace

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
		daoKeeper     types.DAOKeeper
		kycKeeper     types.KycKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	daoKeeper types.DAOKeeper,
	kycKeeper types.KycKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	keeperVal := &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		daoKeeper:     daoKeeper,
		kycKeeper:     kycKeeper,
	}
	keeperVal.kycKeeper.RegisterEventHandler(kycTypes.EventTypeUpdate, 0, types.ModuleName, keeperVal.KycStatusChanged)
	return keeperVal
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) KycStatusChanged(goCtx context.Context, msgType string, data interface{}) error {
	//if eventType
	ctx := sdk.UnwrapSDKContext(goCtx)
	if kycTypes.EventTypeUpdate != msgType {
		return nil
	}
	if val, ok := data.(sdk.Event); !ok {
		return fmt.Errorf("data's type is not sdk.Event.but msgType is update")
	} else {
		attrPreRegion, found := val.GetAttribute(kycTypes.AttributeKeyRegionId)
		if !found {
			return fmt.Errorf("can not found AttributeKeyRegionId.but EventType is update")
		}
		attrNewRegion, found := val.GetAttribute(kycTypes.AttributeKeyRegionIdChanged)
		if !found {
			return fmt.Errorf("can not found AttributeKeyRegionIdChanged.but EventType is update")
		}
		if attrPreRegion.Value == attrNewRegion.Value { //if region not changed,return
			k.Logger(ctx).Info("regionID was not changed in KycStatusChanged!!!")
			return nil
		}
		k.Logger(ctx).Info("start to proc KycStatusChanged!!!")
		attrAddress, found := val.GetAttribute(kycTypes.AttributeKeyAddress)
		if !found {
			return fmt.Errorf("can not found AttributeKeyAddress.but EventType is update")
		}

		if err := k.procKycRegionChange(ctx, attrAddress.Value, attrPreRegion.Value, attrNewRegion.Value); err != nil {
			return err
		}

	}

	return nil

}

func (k Keeper) procKycRegionChange(sdkCtx sdk.Context, address, preRegionID, nowRegionID string) error {
	newGrpId, found := k.GetGroupIdByRegion(sdkCtx, nowRegionID)
	if !found {
		newGrpId = 0
		//	return errors.Wrapf(types.ErrGroupNotExist, fmt.Sprintf("can not found groupId in region.regionID = %s", nowRegionID))
	}
	//if 0 == newGrpId {
	//	return errors.Wrapf(types.ErrProcData, fmt.Sprintf("groupId is 0 in new region.regionID = %s", nowRegionID))
	//}
	joined, JoinGroupFound := k.GetMemberJoined(sdkCtx, address)
	preJoinedGroupID := uint64(0)
	if JoinGroupFound && joined.GroupId > 0 {
		if newGrpId == joined.GroupId {
			k.Logger(sdkCtx).Error("newGrpId == joined.GroupId in procKycRegionChange.", "preJoinedGroupId = ",
				joined.GroupId, "newGroupID = ", newGrpId)
			return nil
		}
		preJoinedGroupID = joined.GroupId
		preGrpIdByRegion, found := k.GetGroupIdByRegion(sdkCtx, preRegionID)
		if !found {
			return errors.Wrapf(types.ErrGroupNotExist, "can not found groupId in previous region.preRegionID = %s."+
				"but user has been joined group.joinGroupID = %d", preRegionID, joined.GroupId)
		}
		if preGrpIdByRegion != joined.GroupId {
			return errors.Wrapf(types.ErrProcData, "preGrpIdByRegion != joined.GroupId.preGrpIdByRegion = %d."+
				"but user has been joined group.joinGroupID = %d", preGrpIdByRegion, joined.GroupId)
		}
		preGroupInfo, found := k.GetGroupInfo(sdkCtx, joined.GroupId)
		if !found {
			return errors.Wrapf(types.ErrGroupNotExist, "can not found joined previous gourp.groupID = %d", joined.GroupId)
		}
		//admin can not migrate
		if address == preGroupInfo.Admin { //admin can not leave group
			return errors.Wrapf(types.ErrExecute, "admin of group can not leave")
		}

		preGroupNumber, found := k.GetGroupMemberCount(sdkCtx, joined.GroupId)
		if !found {
			return fmt.Errorf("can not found preGroup number count while ready to levae preGourp in procKycRegionChange")
		}
		if preGroupNumber == 0 {
			return fmt.Errorf("preGroup number is 0 while ready to levae preGourp in procKycRegionChange")
		}
		if err := k.deleteMemberFormGroup(sdkCtx, joined.GroupId, address); err != nil {
			return err
		}
		k.SetGroupMemberCount(sdkCtx, joined.GroupId, preGroupNumber-1)

	} else {
		return nil
	}
	if 0 == newGrpId {
		if preJoinedGroupID > 0 {
			//set member's join group info
			k.SetMemberJoined(sdkCtx, types.MemberJoined{
				Address: address,
				GroupId: 0,
			})
			sdkCtx.EventManager().EmitEvent(sdk.NewEvent(types.EvtGrpMigrateByKyc,
				sdk.NewAttribute("applicant", address),
				sdk.NewAttribute("previous_region_id", preRegionID),
				sdk.NewAttribute("now_region_id", nowRegionID),
				sdk.NewAttribute("previous_group_id", strconv.FormatUint(preJoinedGroupID, 10)),
				sdk.NewAttribute("now_group_id", "0"),
				//1sdk.NewAttribute("metadata", msg.),
			))
		}
		return nil

	}

	newGrpInfo, found := k.GetGroupInfo(sdkCtx, newGrpId)
	if !found { //if new group has not been created,emit event and return
		return errors.Wrapf(types.ErrGroupNotExist, "can not found group by groupID.groupID = %d", newGrpId)
	}
	newGrpNumberCnt, found := k.GetGroupMemberCount(sdkCtx, newGrpId)
	if !found {
		return fmt.Errorf("can not found newGroup number count while ready to join newGourp in procKycRegionChange")
	}

	newJoin := types.MemberJoined{
		Address: address,
		GroupId: newGrpId,
	}

	//set member's join group info
	k.SetMemberJoined(sdkCtx, newJoin)
	//add to group_member
	err := k.AddGroupMember(sdkCtx, &types.GroupMember{
		GroupId: newGrpId,
		Member: &types.Member{
			Address: address,
			AddedAt: sdkCtx.BlockTime()}})
	if err != nil {
		return err
	}
	k.SetGroupMemberCount(sdkCtx, newGrpId, newGrpNumberCnt+1)
	if !JoinGroupFound { //send rewards if user has not joined group

		//get RegionTreasureAddr
		region, found := k.stakingKeeper.GetRegion(sdkCtx, nowRegionID)
		if !found {
			return errors.Wrapf(types.ErrRegionNotExist, "group's region = %s", nowRegionID)
		}
		rewardsCoin := sdk.NewCoin(params.BaseDenom, math.NewInt(1000000))
		err = k.bankKeeper.Extend().SendCoinsWithTag(sdkCtx, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()),
			sdk.MustAccAddressFromBech32(address), sdk.NewCoins(rewardsCoin), fmt.Sprintf("RegionChange_SendUserReward_%s", region.RegionId))
		if err != nil {
			return errors.Wrap(types.ErrProcData, fmt.Sprintf("transfer rewards coins error. err = %s,fromAddr = %s,toAddr = %s",
				err.Error(), region.GetRegionTreasureAddr(), address))
		}
		err = k.bankKeeper.Extend().SendCoinsWithTag(sdkCtx, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()),
			sdk.MustAccAddressFromBech32(newGrpInfo.Admin), sdk.NewCoins(rewardsCoin), fmt.Sprintf("RegionChange_SendAdminReward_%s", region.RegionId))
		if err != nil {
			return errors.Wrapf(types.ErrProcData, "transfer rewards coins error. err = %s,fromAddr = %s,toAddr = %s",
				err.Error(), region.GetRegionTreasureAddr(), newGrpInfo.Admin)
		}
		sdkCtx.EventManager().EmitEvent(sdk.NewEvent(types.EvtJoinGroupReward,
			sdk.NewAttribute("applicant", address),
			sdk.NewAttribute("admin", newGrpInfo.Admin),
			sdk.NewAttribute("regionTreasureAddress", region.GetRegionTreasureAddr()),
			sdk.NewAttribute("rewards", rewardsCoin.String()),
		))

	}
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(types.EvtGrpMigrateByKyc,
		sdk.NewAttribute("applicant", address),
		sdk.NewAttribute("previous_region_id", preRegionID),
		sdk.NewAttribute("now_region_id", nowRegionID),
		sdk.NewAttribute("previous_group_id", strconv.FormatUint(preJoinedGroupID, 10)),
		sdk.NewAttribute("now_group_id", strconv.FormatUint(newGrpId, 10)),
		//1sdk.NewAttribute("metadata", msg.),
	))
	return nil

}

func (k Keeper) GetDidAndKycActive(sdkCtx sdk.Context, address sdk.AccAddress, regionID string) (string, bool) {
	didVal, found := k.kycKeeper.GetDID(sdkCtx, address)
	if !found {
		return "", false
	}
	didInfo, found := k.kycKeeper.GetDidInfo(sdkCtx, didVal)
	if !found {
		return "", false
	}
	if didInfo.RegionId != regionID {
		return "", false
	}
	if didInfo.KycLevel != didTypes.KYC_LEVEL_TWO {
		return "", false
	}
	if didInfo.Status == didTypes.DID_STATUS_ACTIVE {
		return didVal, true
	}
	return didVal, false
}

// only called by wstaking/newRgion
func (k Keeper) CreateGroupByRegion(sdkCtx sdk.Context, regionInfo stakingTypes.Region) (uint64, error) {

	//check group has been created
	preGroupID, found := k.GetGroupIdByRegion(sdkCtx, regionInfo.RegionId)
	if found {
		return 0, errors.Wrap(types.ErrGroupExceededInRegion, fmt.Sprintf("group of region has been created.groupId = %d", preGroupID))
	}

	operValAddr, err := sdk.ValAddressFromBech32(regionInfo.OperatorAddress)
	if err != nil {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Sprintf("OperatorAddress can not convert to ValAddress."+
			"err = %s, OperatorAddress = %s", err.Error(), regionInfo.OperatorAddress))
	}
	accAddr := sdk.AccAddress(operValAddr.Bytes())

	newGroupID := k.GetLastGroupID(sdkCtx) + 1
	groupInfo := &types.GroupInfo{
		Id:          newGroupID,
		Admin:       accAddr.String(),
		Metadata:    "",
		Version:     1,
		TotalWeight: math.NewInt(0).String(),
		CreatedAt:   sdkCtx.BlockTime(),
		RegionID:    regionInfo.RegionId,
	}
	if err := k.AppendGroup(sdkCtx, groupInfo); err != nil {
		return 0, err
	}
	k.SetGroupToRegion(sdkCtx, regionInfo.RegionId, newGroupID)
	k.SetGroupMemberCount(sdkCtx, newGroupID, 0)
	//group's admin(region's operator) no need to add group

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(types.EvtGroupCreated,
		sdk.NewAttribute("group_id", fmt.Sprintf("%d", groupInfo.Id)),
		sdk.NewAttribute("admin", groupInfo.Admin),
		sdk.NewAttribute("region_id", groupInfo.RegionID),
		sdk.NewAttribute("metadata", groupInfo.Metadata),
	))
	return groupInfo.Id, nil

}
