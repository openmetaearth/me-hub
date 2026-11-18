package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (k Keeper) lookupRegion(ctx sdk.Context, regionID string) (wstakingtypes.Region, bool) {
	if region, found := k.stakingKeeper.GetRegion(ctx, regionID); found {
		return region, true
	}
	lowered := strings.ToLower(regionID)
	if lowered == regionID {
		return wstakingtypes.Region{}, false
	}
	return k.stakingKeeper.GetRegion(ctx, lowered)
}

func (k Keeper) isRegionFrozen(ctx sdk.Context, regionID string) bool {
	freeze, found := k.GetRegionFreeze(ctx, regionID)
	return found && freeze.Frozen
}

func (k Keeper) isDao(ctx sdk.Context, address string) bool {
	return k.daoKeeper.IsGlobalDao(ctx, address)
}

// operatorAccount converts Region.OperatorAddress into the account that signs a ballot.
func operatorAccount(operator string) (sdk.AccAddress, bool) {
	if operator == "" {
		return nil, false
	}
	val, err := sdk.ValAddressFromBech32(operator)
	if err != nil {
		return nil, false
	}
	return sdk.AccAddress(val.Bytes()), true
}

func isSpecialRegion(regionID string) bool {
	switch strings.ToLower(regionID) {
	case wstakingtypes.MeEarthRegionId, wstakingtypes.ExperienceRegionId:
		return true
	default:
		return false
	}
}

func containsRegion(regionIDs []string, regionID string) bool {
	for _, id := range regionIDs {
		if strings.EqualFold(id, regionID) {
			return true
		}
	}
	return false
}

func canonicalAccAddress(address string) (string, error) {
	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return "", err
	}
	return acc.String(), nil
}
