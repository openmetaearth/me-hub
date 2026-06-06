package keeper_test

import (
	"testing"

	cometbftdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	didkeeper "github.com/openmetaearth/me-hub/x/did/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
)

func TestUpdateDidStatusInactiveRevokesServiceIssuerMemberships(t *testing.T) {
	keeper, ctx := didKeeperWithGlobalDao(t)
	msgServer := didkeeper.NewMsgServerImpl(keeper)

	issuerDID := "did:issuer"
	otherIssuerDID := "did:other-issuer"
	keeper.SetDidInfo(ctx, issuerDID, didtypes.DidInfo{
		Did:    issuerDID,
		Status: didtypes.DID_STATUS_ACTIVE,
	})
	keeper.SetService(ctx, "kyc", didtypes.NewService(
		"kyc",
		"KYC",
		"",
		didtypes.SERVICE_STATUS_ACTIVE,
		[]string{issuerDID, otherIssuerDID, issuerDID},
	))
	keeper.SetService(ctx, "identity", didtypes.NewService(
		"identity",
		"Identity",
		"",
		didtypes.SERVICE_STATUS_ACTIVE,
		[]string{otherIssuerDID},
	))

	_, err := msgServer.UpdateDidStatus(sdk.WrapSDKContext(ctx), &didtypes.MsgUpdateDidStatus{
		Creator: "dao",
		Did:     issuerDID,
		Status:  didtypes.DID_STATUS_INACTIVE,
	})
	require.NoError(t, err)

	info, found := keeper.GetDidInfo(ctx, issuerDID)
	require.True(t, found)
	require.Equal(t, didtypes.DID_STATUS_INACTIVE, info.Status)

	kycService, found := keeper.GetService(ctx, "kyc")
	require.True(t, found)
	require.Equal(t, []string{otherIssuerDID}, kycService.Issuers)

	identityService, found := keeper.GetService(ctx, "identity")
	require.True(t, found)
	require.Equal(t, []string{otherIssuerDID}, identityService.Issuers)
}

type globalDaoKeeper struct{}

func (globalDaoKeeper) IsGlobalDao(ctx sdk.Context, address string) bool {
	return true
}

func didKeeperWithGlobalDao(t testing.TB) (*didkeeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(didtypes.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(didtypes.MemStoreKey)

	db := cometbftdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	keeper := didkeeper.NewKeeper(cdc, storeKey, globalDaoKeeper{})
	ctx := sdk.NewContext(stateStore, cometbftproto.Header{}, false, log.NewNopLogger())

	return keeper, ctx
}
