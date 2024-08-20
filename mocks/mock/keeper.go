package mock

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m MockStakingKeeper) GetAllRegion(ctx sdk.Context) (list []MockRegion) {
	return []MockRegion{{RegionId: "aaa", RegionShare: math.NewInt(1), RegionTreasureAddr: "cosmos1xvh8nef0tj5w00cntns3mxy43nxxg3jss9d4k4"}}
}

// func (m MockStakingKeeper) GetRegion(ctx sdk.Context, regionId string) (val testutil.Region, found bool) {
// 	region := testutil.Region{RegionId: "aaa", RegionShare: math.NewInt(1), RegionTreasureAddr: "cosmos1xvh8nef0tj5w00cntns3mxy43nxxg3jss9d4k4"}
// 	return region, true
// }

func (m MockStakingKeeper) CalculateInterest(ctx sdk.Context, totalStaking math.Int, height int64) (rewards sdk.Dec, err error) {
	return sdk.NewDec(1), nil
}

func (m *MockStakingKeeper) SetDelegation(ctx sdk.Context, delegation Delegation) {
	m.Keeper.SetDelegation(ctx, delegation.Delegation)
}

func (m *MockStakingKeeper) SetRegion(ctx sdk.Context, region MockRegion) {
	fmt.Println("set region")
}
