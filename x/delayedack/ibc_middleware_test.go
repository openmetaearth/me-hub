package delayedack_test

import (
	"errors"
	"testing"

	cometbftdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/openmetaearth/me-hub/testutil/keeper"
	commontypes "github.com/openmetaearth/me-hub/x/common/types"
	"github.com/openmetaearth/me-hub/x/delayedack"
	delayedackkeeper "github.com/openmetaearth/me-hub/x/delayedack/keeper"
	delayedacktypes "github.com/openmetaearth/me-hub/x/delayedack/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

func TestOnRecvPacketRollsBackPendingPacketWhenEIBCHandlerFails(t *testing.T) {
	packetData := transfertypes.NewFungibleTokenPacketData(
		"stake",
		"100",
		"sender",
		"receiver",
		`{"eibc":{"fee":"1000"}}`,
	)
	packet := channeltypes.NewPacket(
		packetData.GetBytes(),
		7,
		transfertypes.PortID,
		"channel-0",
		transfertypes.PortID,
		"channel-1",
		clienttypes.Height{},
		0,
	)

	eibcKeeper := &failingEIBCKeeper{err: errors.New("demand order rejected")}
	keeper, ctx := newDelayedAckKeeper(t, enabledRollappKeeper{
		transferData: rollapptypes.TransferData{
			FungibleTokenPacketData: packetData,
			Rollapp:                 &rollapptypes.Rollapp{RollappId: "rollapp-1"},
		},
	}, eibcKeeper)

	packetID := commontypes.NewPacketUID(
		commontypes.RollappPacket_ON_RECV,
		packet.GetDestPort(),
		packet.GetDestChannel(),
		packet.Sequence,
	)
	ctx = delayedacktypes.CtxWithPacketProofHeight(ctx, packetID, clienttypes.NewHeight(0, 10))

	middleware := delayedack.NewIBCMiddleware(delayedack.WithKeeper(*keeper))
	ack := middleware.OnRecvPacket(ctx, packet, sdk.AccAddress("relayer"))

	require.NotNil(t, ack)
	require.True(t, eibcKeeper.called)
	require.Empty(t, keeper.GetAllRollappPackets(ctx))
}

type failingEIBCKeeper struct {
	err    error
	called bool
}

func (k *failingEIBCKeeper) EIBCDemandOrderHandler(
	ctx sdk.Context,
	rollappPacket commontypes.RollappPacket,
	data transfertypes.FungibleTokenPacketData,
) error {
	k.called = true
	return k.err
}

type enabledRollappKeeper struct {
	transferData rollapptypes.TransferData
}

func (k enabledRollappKeeper) GetParams(ctx sdk.Context) rollapptypes.Params {
	return rollapptypes.Params{RollappsEnabled: true}
}

func (k enabledRollappKeeper) GetStateInfo(
	ctx sdk.Context,
	rollappID string,
	index uint64,
) (rollapptypes.StateInfo, bool) {
	return rollapptypes.StateInfo{}, false
}

func (k enabledRollappKeeper) MustGetStateInfo(
	ctx sdk.Context,
	rollappID string,
	index uint64,
) rollapptypes.StateInfo {
	return rollapptypes.StateInfo{}
}

func (k enabledRollappKeeper) GetLatestFinalizedStateIndex(
	ctx sdk.Context,
	rollappID string,
) (rollapptypes.StateInfoIndex, bool) {
	return rollapptypes.StateInfoIndex{}, false
}

func (k enabledRollappKeeper) GetAllRollapps(ctx sdk.Context) []rollapptypes.Rollapp {
	return nil
}

func (k enabledRollappKeeper) GetValidTransfer(
	ctx sdk.Context,
	packetData []byte,
	raPortOnHub string,
	raChanOnHub string,
) (rollapptypes.TransferData, error) {
	return k.transferData, nil
}

func (k enabledRollappKeeper) IsSkipDelayRollapp(ctx sdk.Context, rollappID string) bool {
	return false
}

func newDelayedAckKeeper(
	t testing.TB,
	rollappKeeper delayedacktypes.RollappKeeper,
	eibcKeeper delayedacktypes.EIBCKeeper,
) (*delayedackkeeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(delayedacktypes.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(delayedacktypes.MemStoreKey)

	db := cometbftdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	paramsSubspace := typesparams.NewSubspace(
		cdc,
		delayedacktypes.Amino,
		storeKey,
		memStoreKey,
		"DelayedackParams",
	)

	keeper := delayedackkeeper.NewKeeper(
		cdc,
		storeKey,
		paramsSubspace,
		rollappKeeper,
		keepertest.ICS4WrapperStub{},
		keepertest.ChannelKeeperStub{},
		eibcKeeper,
	)
	ctx := sdk.NewContext(stateStore, cometbftproto.Header{}, false, log.NewNopLogger())
	keeper.SetParams(ctx, delayedacktypes.DefaultParams())

	return keeper, ctx
}
