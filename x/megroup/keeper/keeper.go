func (k Keeper) procKycRegionChange(sdkCtx sdk.Context, address, preRegionID, nowRegionID string) error {
    newGrpId, found := k.GetGroupIdByRegion(sdkCtx, nowRegionID)
    if !found {
        newGrpId = 0
    }
    joined, JoinGroupFound := k.GetMemberJoined(sdkCtx, address)
    preJoinedGroupID := uint64(0)
    if JoinGroupFound && joined.GroupId > 0 {
        if newGrpId == joined.GroupId {
            k.Logger(sdkCtx).Error("newGrpId == joined.GroupId in procKycRegionChange.", "preJoinedGroupId = ", preJoinedGroupID)
            return nil
        }
        preJoinedGroupID = joined.GroupId
    }

    // Check if user has joined a group before
    if !JoinGroupFound {
        // Do not send rewards if user has never joined a group before
        return nil
    }

    // Check if user is joining a new group
    if preJoinedGroupID != 0 && preJoinedGroupID != newGrpId {
        rewardsCoin := sdk.NewCoin(params.BaseDenom, math.NewInt(1000000))
        // Send reward to user
        if err := k.bankKeeper.Extend().SendCoinsWithTag(sdkCtx, k.accountKeeper.GetModuleAddress(types.ModuleName), address, rewardsCoin, "RegionChange_SendUserReward_"+nowRegionID); err != nil {
            return err
        }
        // Send reward to admin
        adminAddr, found := k.GetAdminAddress(sdkCtx, newGrpId)
        if found {
            if err := k.bankKeeper.Extend().SendCoinsWithTag(sdkCtx, k.accountKeeper.GetModuleAddress(types.ModuleName), adminAddr, rewardsCoin, "RegionChange_SendAdminReward_"+nowRegionID); err != nil {
                // Rollback user reward if admin reward fails
                if err := k.bankKeeper.Extend().SendCoinsWithTag(sdkCtx, address, k.accountKeeper.GetModuleAddress(types.ModuleName), rewardsCoin, "RegionChange_RollbackUserReward_"+nowRegionID); err != nil {
                    return err
                }
                return err
            }
        }
    }
    return nil
}