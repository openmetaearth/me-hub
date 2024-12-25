package keeper

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/st-chain/me-hub/app/params"
	"github.com/st-chain/me-hub/x/megroup/types"
)

type kycHookFunc func(ctx sdk.Context, eventType string, beforeData interface{}, afterData interface{}) error

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
		daoKeeper     types.DAOKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	daoKeeper types.DAOKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		daoKeeper:     daoKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) KycStatusChanged(ctx sdk.Context, msgType string, data interface{}) error {
	//if eventType
	/*
		ctx.EventManager().EmitEvent(sdk.NewEvent(types.EvtGrpMigrateByKyc,
			sdk.NewAttribute("group_id", fmt.Sprintf("%d", msg.GroupId)),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("user"),
			//1sdk.NewAttribute("metadata", msg.),
		))

	*/
	return nil

}

func (k Keeper) procKycRegionChange(sdkCtx sdk.Context, address, nowRegionID string) error {
	newGrpId, found := k.GetGroupIdByRegion(sdkCtx, nowRegionID)
	if !found {
		return errors.Wrapf(types.ErrGroupNotExist, fmt.Sprintf("can not found groupId in region.regionID = %s", nowRegionID))
	}
	if 0 == newGrpId {
		return errors.Wrapf(types.ErrProcData, fmt.Sprintf("groupId is 0 in new region.regionID = %s", nowRegionID))
	}
	joined, JoinGroupFound := k.GetMemberJoined(sdkCtx, address)
	if JoinGroupFound && joined.GroupId > 0 {
		if newGrpId == joined.GroupId {
			k.Logger(sdkCtx).Error("newGrpId == joined.GroupId in procKycRegionChange.", "preJoinedGroupId = ",
				joined.GroupId, "newGroupID = ", newGrpId)
			return nil
		}
		if err := k.deleteMemberFormGroup(sdkCtx, joined.GroupId, address); err != nil {
			return err
		}
	}

	newGrpInfo, found := k.GetGroup(sdkCtx, newGrpId)
	if !found {
		return errors.Wrapf(types.ErrGroupNotExist, fmt.Sprintf("can not found group by groupID.groupID = %d", newGrpId))
	}

	newJoin := types.MemberJoined{
		Address: address,
		GroupId: newGrpId,
	}

	//set member's join group info
	k.SetMemberJoined(sdkCtx, newJoin)
	//add to group_member
	err := k.addGroupMember(sdkCtx, &types.GroupMember{
		GroupID: newGrpId,
		Member: &types.Member{
			Address: address,
			AddedAt: sdkCtx.BlockTime()}})
	if err != nil {
		return err
	}
	if !JoinGroupFound { //send rewards if user has not joined group

		//get RegionTreasureAddr
		region, found := k.stakingKeeper.GetRegion(sdkCtx, nowRegionID)
		if !found {
			return errors.Wrapf(types.ErrRegionNotExist, fmt.Sprintf("group's region = %d", nowRegionID))
		}
		rewardsCoin := sdk.NewCoin(params.BaseDenom, math.NewInt(1000000))
		err = k.bankKeeper.SendCoins(sdkCtx, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()),
			sdk.MustAccAddressFromBech32(address), sdk.NewCoins(rewardsCoin))
		if err != nil {
			return errors.Wrapf(types.ErrProcData, fmt.Sprintf("transfer rewards coins error. err = %s,fromAddr = %s,toAddr = %s",
				err.Error(), region.GetRegionTreasureAddr(), address))
		}
		err = k.bankKeeper.SendCoins(sdkCtx, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()),
			sdk.MustAccAddressFromBech32(newGrpInfo.Admin), sdk.NewCoins(rewardsCoin))
		if err != nil {
			return errors.Wrapf(types.ErrProcData, fmt.Sprintf("transfer rewards coins error. err = %s,fromAddr = %s,toAddr = %s",
				err.Error(), region.GetRegionTreasureAddr(), newGrpInfo.Admin))
		}

	}
	return nil

}
