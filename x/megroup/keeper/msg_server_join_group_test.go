package keeper

import (
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/megroup/types"
	"github.com/stretchr/testify/require"
)

func TestJoinGroupRejectsGroupAdminSelfJoin(t *testing.T) {
	ctx, keeper := setupJoinGroupTestKeeper(t)
	admin := sample.AccAddress()
	groupID := uint64(7)

	require.NoError(t, keeper.AppendGroup(ctx, &types.GroupInfo{
		Id:       groupID,
		Admin:    admin,
		RegionID: "usa",
	}))
	keeper.SetGroupMemberCount(ctx, groupID, 0)

	server := NewMsgServerImpl(keeper)
	_, err := server.JoinGroup(sdk.WrapSDKContext(ctx), &types.MsgJoinGroup{
		Creator:          admin,
		GroupId:          groupID,
		ApplicantAddress: admin,
	})

	require.ErrorIs(t, err, types.ErrPermissionDenied)
	require.Contains(t, err.Error(), "group admin does not need to join own group")

	_, found := keeper.GetMemberJoined(ctx, admin)
	require.False(t, found)
	memberCount, found := keeper.GetGroupMemberCount(ctx, groupID)
	require.True(t, found)
	require.Equal(t, uint64(0), memberCount)
}

func setupJoinGroupTestKeeper(t *testing.T) (sdk.Context, Keeper) {
	t.Helper()
	params.SetAddressPrefixes()

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	ctx := sdk.NewContext(stateStore, tmproto.Header{Time: time.Now()}, false, log.NewNopLogger())

	return ctx, Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}
