package keeper

import (
	"fmt"
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/handler"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/openmetaearth/me-hub/x/megroup/types"
	"github.com/stretchr/testify/require"
)

type kycMembershipTestKeeper struct {
	didByAddress map[string]string
	infoByDid    map[string]didtypes.DidInfo
}

func (k *kycMembershipTestKeeper) RegisterEventHandler(string, int, string, handler.HandlerFunc) {}

func (k *kycMembershipTestKeeper) GetDID(_ sdk.Context, address sdk.AccAddress) (string, bool) {
	did, found := k.didByAddress[address.String()]
	return did, found
}

func (k *kycMembershipTestKeeper) GetDidInfo(_ sdk.Context, did string) (didtypes.DidInfo, bool) {
	info, found := k.infoByDid[did]
	return info, found
}

func TestKycStatusChangedEnforcesFinalEligibility(t *testing.T) {
	tests := []struct {
		name          string
		preRegionID   string
		nowRegionID   string
		level         didtypes.KycLevel
		wantGroupID   uint64
		wantNewCount  uint64
		wantNewMember bool
	}{
		{
			name:        "region change downgrade removes membership",
			preRegionID: "meearth",
			nowRegionID: "usa",
			level:       didtypes.KYC_LEVEL_ONE,
		},
		{
			name:        "same region downgrade removes membership",
			preRegionID: "meearth",
			nowRegionID: "meearth",
			level:       didtypes.KYC_LEVEL_ONE,
		},
		{
			name:          "eligible region change migrates membership",
			preRegionID:   "meearth",
			nowRegionID:   "usa",
			level:         didtypes.KYC_LEVEL_TWO,
			wantGroupID:   2,
			wantNewCount:  1,
			wantNewMember: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, keeper, address := setupKycMembershipChangeTest(t, tt.nowRegionID, tt.level)
			event := sdk.NewEvent(kyctypes.EventTypeUpdate,
				sdk.NewAttribute(kyctypes.AttributeKeyAddress, address),
				sdk.NewAttribute(kyctypes.AttributeKeyRegionId, tt.preRegionID),
				sdk.NewAttribute(kyctypes.AttributeKeyRegionIdChanged, tt.nowRegionID),
			)

			require.NoError(t, keeper.KycStatusChanged(sdk.WrapSDKContext(ctx), kyctypes.EventTypeUpdate, event))

			joined, found := keeper.GetMemberJoined(ctx, address)
			require.True(t, found)
			require.Equal(t, tt.wantGroupID, joined.GroupId)
			oldCount, found := keeper.GetGroupMemberCount(ctx, 1)
			require.True(t, found)
			require.Equal(t, uint64(0), oldCount)
			newCount, found := keeper.GetGroupMemberCount(ctx, 2)
			require.True(t, found)
			require.Equal(t, tt.wantNewCount, newCount)
			require.False(t, hasGroupMember(ctx, keeper, 1, address))
			require.Equal(t, tt.wantNewMember, hasGroupMember(ctx, keeper, 2, address))
		})
	}
}

func setupKycMembershipChangeTest(t *testing.T, finalRegionID string, level didtypes.KycLevel) (sdk.Context, Keeper, string) {
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
	address := sample.AccAddress()
	did := "test-did"
	kycKeeper := &kycMembershipTestKeeper{
		didByAddress: map[string]string{address: did},
		infoByDid: map[string]didtypes.DidInfo{
			did: {
				Did:      did,
				Address:  address,
				Status:   didtypes.DID_STATUS_ACTIVE,
				RegionId: finalRegionID,
				KycLevel: level,
			},
		},
	}
	keeper := Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		kycKeeper: kycKeeper,
	}

	require.NoError(t, keeper.AppendGroup(ctx, &types.GroupInfo{Id: 1, Admin: sample.AccAddress(), RegionID: "meearth"}))
	require.NoError(t, keeper.AppendGroup(ctx, &types.GroupInfo{Id: 2, Admin: sample.AccAddress(), RegionID: "usa"}))
	keeper.SetGroupToRegion(ctx, "meearth", 1)
	keeper.SetGroupToRegion(ctx, "usa", 2)
	keeper.SetGroupMemberCount(ctx, 1, 1)
	keeper.SetGroupMemberCount(ctx, 2, 0)
	keeper.SetMemberJoined(ctx, types.MemberJoined{Address: address, GroupId: 1})
	require.NoError(t, keeper.AddGroupMember(ctx, &types.GroupMember{
		GroupId: 1,
		Member: &types.Member{
			Address: address,
			AddedAt: ctx.BlockTime(),
		},
	}))

	return ctx, keeper, address
}

func hasGroupMember(ctx sdk.Context, keeper Keeper, groupID uint64, address string) bool {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), []byte(fmt.Sprintf("%s%d/", types.GroupMemberKey, groupID)))
	return store.Get([]byte(address)) != nil
}
