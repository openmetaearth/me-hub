package megroup

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/megroup/keeper"
	"github.com/openmetaearth/me-hub/x/megroup/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	lastGroupID := genState.GroupCount
	for _, elem := range genState.Groups {
		k.SetGroupInfo(ctx, elem)
		if elem.RegionID != "" {
			k.SetGroupToRegion(ctx, elem.RegionID, elem.Id)
		}
		if elem.Id > lastGroupID {
			lastGroupID = elem.Id
		}
	}
	k.SetLastGroupID(ctx, lastGroupID)

	for _, elem := range genState.GroupMembers {
		k.SetGroupMember(ctx, elem)
	}

	for _, elem := range genState.MemberJoinedList {
		k.SetMemberJoined(ctx, elem)
	}

	for _, elem := range genState.GroupMemberCountList {
		k.SetGroupMemberCount(ctx, elem.GroupId, elem.Num)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.Groups = k.GetAllGroup(ctx)
	genesis.GroupCount = k.GetLastGroupID(ctx)
	genesis.GroupMembers = k.GetAllGroupMember(ctx)
	genesis.MemberJoinedList = k.GetAllMemberJoined(ctx)
	genesis.GroupMemberCountList = k.GetAllGroupMemberCount(ctx)

	return genesis
}
