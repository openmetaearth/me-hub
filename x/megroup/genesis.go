package megroup

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/megroup/keeper"
	"github.com/st-chain/me-hub/x/megroup/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the group
	for _, elem := range genState.GroupList {
		k.SetGroup(ctx, elem)
	}

	// Set group count
	k.SetGroupCount(ctx, genState.GroupCount)
	// Set all the groupMember
	for _, elem := range genState.GroupMemberList {
		k.LoadMemberStoreByGroupID(ctx, elem.GroupID).SetGroupMember(elem)
	}

	// Set all the memberJoined
	for _, elem := range genState.MemberJoinedList {
		k.SetMemberJoined(ctx, elem)
	}
	// Set all the groupMemberCount
	for _, elem := range genState.GroupMemberCountList {
		k.SetGroupMemberCount(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.GroupList = k.GetAllGroup(ctx)
	genesis.GroupCount = k.GetGroupCount(ctx)
	genesis.GroupMemberList = k.GetAllGroupMember(ctx)
	genesis.MemberJoinedList = k.GetAllMemberJoined(ctx)
	genesis.GroupMemberCountList = k.GetAllGroupMemberCount(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
