package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	lightclientkeeper "github.com/openmetaearth/me-hub/x/lightclient/keeper"
	lightclienttypes "github.com/openmetaearth/me-hub/x/lightclient/types"
	sequencertypes "github.com/openmetaearth/me-hub/x/sequencer/types"
)

const (
	DefaultRollapp = "rollapp_1234-1"
	CanonClientID  = "07-tendermint-0"
)

var Alice = sequencertypes.Sequencer{
	Address:   "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u",
	RollappId: DefaultRollapp,
}

func LightClientKeeper(t testing.TB) (*lightclientkeeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(lightclienttypes.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	k := lightclientkeeper.NewKeeper(cdc, storeKey, nil, nil, nil, nil, nil)
	ctx := sdk.NewContext(stateStore, cometbftproto.Header{}, false, log.NewNopLogger())
	return k, ctx
}
