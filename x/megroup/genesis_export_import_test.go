package megroup

import (
	"testing"
	"time"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/handler"
	"github.com/openmetaearth/me-hub/x/megroup/keeper"
	"github.com/openmetaearth/me-hub/x/megroup/types"
	"github.com/stretchr/testify/require"
)

type megroupGenesisKycKeeper struct{}

func (megroupGenesisKycKeeper) RegisterEventHandler(string, int, string, handler.HandlerFunc) {}

func (megroupGenesisKycKeeper) GetDID(sdk.Context, sdk.AccAddress) (string, bool) {
	return "", false
}

func (megroupGenesisKycKeeper) GetDidInfo(sdk.Context, string) (didtypes.DidInfo, bool) {
	return didtypes.DidInfo{}, false
}

func setupMegroupGenesisKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	t.Helper()

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	paramsSubspace := paramtypes.NewSubspace(cdc, types.Amino, storeKey, memStoreKey, "MegroupParams")
	k := keeper.NewKeeper(cdc, storeKey, paramsSubspace, nil, nil, nil, nil, megroupGenesisKycKeeper{})
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}

func megroupGenesisFixture() types.GenesisState {
	addedAt := time.Unix(1700000000, 0).UTC()

	return types.GenesisState{
		Params: types.DefaultParams(),
		Groups: []types.GroupInfo{
			{
				Id:          7,
				Admin:       "admin-address",
				Metadata:    "metadata",
				Version:     1,
				TotalWeight: "1",
				CreatedAt:   addedAt,
				RegionID:    "region-7",
			},
		},
		GroupCount: 7,
		GroupMembers: []types.GroupMember{
			{
				GroupId: 7,
				Member: &types.Member{
					Address:  "member-address",
					Weight:   "1",
					Metadata: "member-metadata",
					AddedAt:  addedAt,
				},
			},
		},
		MemberJoinedList: []types.MemberJoined{
			{
				Address: "member-address",
				GroupId: 7,
			},
		},
		GroupMemberCountList: []types.GroupMemberCount{
			{
				GroupId: 7,
				Num:     1,
			},
		},
	}
}

func TestExportGenesisMustPreserveGroupRewardState(t *testing.T) {
	k, ctx := setupMegroupGenesisKeeper(t)
	fixture := megroupGenesisFixture()

	group := fixture.Groups[0]
	k.SetGroupInfo(ctx, group)
	k.SetGroupToRegion(ctx, group.RegionID, group.Id)
	k.SetLastGroupID(ctx, fixture.GroupCount)
	k.SetGroupMember(ctx, fixture.GroupMembers[0])
	k.SetMemberJoined(ctx, fixture.MemberJoinedList[0])
	k.SetGroupMemberCount(ctx, fixture.GroupMemberCountList[0].GroupId, fixture.GroupMemberCountList[0].Num)

	exported := ExportGenesis(ctx, *k)
	require.Equal(t, fixture.GroupCount, exported.GroupCount)
	require.Equal(t, fixture.Groups, exported.Groups)
	require.Equal(t, fixture.GroupMembers, exported.GroupMembers)
	require.Equal(t, fixture.MemberJoinedList, exported.MemberJoinedList)
	require.Equal(t, fixture.GroupMemberCountList, exported.GroupMemberCountList)
}

func TestInitGenesisMustRestoreGroupRewardState(t *testing.T) {
	k, ctx := setupMegroupGenesisKeeper(t)
	fixture := megroupGenesisFixture()

	InitGenesis(ctx, *k, fixture)

	group, found := k.GetGroupInfo(ctx, fixture.Groups[0].Id)
	require.True(t, found)
	require.Equal(t, fixture.Groups[0], group)

	regionGroupID, found := k.GetGroupIdByRegion(ctx, fixture.Groups[0].RegionID)
	require.True(t, found)
	require.Equal(t, fixture.Groups[0].Id, regionGroupID)
	require.Equal(t, fixture.GroupCount, k.GetLastGroupID(ctx))

	groupMembers := k.GetAllGroupMember(ctx)
	require.Equal(t, fixture.GroupMembers, groupMembers)

	memberJoined, found := k.GetMemberJoined(ctx, fixture.MemberJoinedList[0].Address)
	require.True(t, found)
	require.Equal(t, fixture.MemberJoinedList[0], memberJoined)

	groupMemberCount, found := k.GetGroupMemberCount(ctx, fixture.GroupMemberCountList[0].GroupId)
	require.True(t, found)
	require.Equal(t, fixture.GroupMemberCountList[0].Num, groupMemberCount)
	require.Equal(t, fixture.GroupMemberCountList, k.GetAllGroupMemberCount(ctx))

	restored := ExportGenesis(ctx, *k)
	require.Equal(t, fixture, *restored)
}

func TestInitGenesisMustAllowExportOnEmptyState(t *testing.T) {
	k, ctx := setupMegroupGenesisKeeper(t)

	InitGenesis(ctx, *k, *types.DefaultGenesis())

	exported := ExportGenesis(ctx, *k)
	require.Equal(t, types.DefaultParams(), exported.Params)
	require.Empty(t, exported.Groups)
	require.Empty(t, exported.GroupMembers)
	require.Empty(t, exported.MemberJoinedList)
	require.Empty(t, exported.GroupMemberCountList)
}
