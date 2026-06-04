package transfergenesis

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
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/stretchr/testify/require"

	delayedackkeeper "github.com/openmetaearth/me-hub/x/delayedack/keeper"
	delayedacktypes "github.com/openmetaearth/me-hub/x/delayedack/types"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

const (
	skipDelayTestRollappID = "rollappevm_1234-1"
	skipDelayTestChannelID = "channel-0"
)

type transferGenesisChannelKeeper struct{}

func (transferGenesisChannelKeeper) GetChannelClientState(
	sdk.Context,
	string,
	string,
) (string, exported.ClientState, error) {
	return "07-tendermint-0", &ibctmtypes.ClientState{ChainId: skipDelayTestRollappID}, nil
}

func (transferGenesisChannelKeeper) LookupModuleByChannel(
	sdk.Context,
	string,
	string,
) (string, *capabilitytypes.Capability, error) {
	return "", nil, nil
}

type transferGenesisNextModule struct {
	porttypes.IBCModule
	recvCalls int
}

func (m *transferGenesisNextModule) OnRecvPacket(
	sdk.Context,
	channeltypes.Packet,
	sdk.AccAddress,
) exported.Acknowledgement {
	m.recvCalls++
	return channeltypes.NewResultAcknowledgement([]byte{1})
}

func TestSkipDelayRollappCannotBypassPostGenesisTransferMemoRejection(t *testing.T) {
	ctx, rollappKeeper, delayedackKeeper := setupSkipDelayTransferGenesisKeepers(t)
	rollappKeeper.SetRollapp(ctx, rollapptypes.Rollapp{
		RollappId: skipDelayTestRollappID,
		GenesisState: rollapptypes.RollappGenesisState{
			TransfersEnabled: true,
		},
		ChannelId: skipDelayTestChannelID,
	})
	rollappKeeper.SetSkipDelayRollapp(ctx, skipDelayTestRollappID, true)

	t.Run("ordinary transfer still passes through", func(t *testing.T) {
		next := &transferGenesisNextModule{}
		module := NewIBCModule(next, *delayedackKeeper, *rollappKeeper, nil, nil)

		ack := module.OnRecvPacket(ctx, skipDelayTransferGenesisPacket(""), nil)

		require.True(t, ack.Success())
		require.Equal(t, 1, next.recvCalls)
	})

	t.Run("post genesis genesis_transfer is rejected before next module", func(t *testing.T) {
		next := &transferGenesisNextModule{}
		module := NewIBCModule(next, *delayedackKeeper, *rollappKeeper, nil, nil)

		ack := module.OnRecvPacket(ctx, skipDelayTransferGenesisPacket(memoHappyPath), nil)

		require.False(t, ack.Success())
		require.Equal(t, 0, next.recvCalls)
	})
}

func setupSkipDelayTransferGenesisKeepers(
	t *testing.T,
) (sdk.Context, *rollappkeeper.Keeper, *delayedackkeeper.Keeper) {
	t.Helper()

	rollappStoreKey := sdk.NewKVStoreKey(rollapptypes.StoreKey)
	rollappMemStoreKey := storetypes.NewMemoryStoreKey(rollapptypes.MemStoreKey)
	delayedackStoreKey := sdk.NewKVStoreKey(delayedacktypes.StoreKey)
	delayedackMemStoreKey := storetypes.NewMemoryStoreKey(delayedacktypes.MemStoreKey)

	db := cometbftdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(rollappStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(rollappMemStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(delayedackStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(delayedackMemStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	channelKeeper := transferGenesisChannelKeeper{}
	rollappParams := paramtypes.NewSubspace(
		cdc,
		rollapptypes.Amino,
		rollappStoreKey,
		rollappMemStoreKey,
		"RollappParams",
	)
	rollappKeeper := rollappkeeper.NewKeeper(cdc, rollappStoreKey, rollappParams, channelKeeper, nil, nil)
	ctx := sdk.NewContext(stateStore, cometbftproto.Header{}, false, log.NewNopLogger())
	rollappKeeper.SetParams(ctx, rollapptypes.DefaultParams())

	delayedackParams := paramtypes.NewSubspace(
		cdc,
		delayedacktypes.Amino,
		delayedackStoreKey,
		delayedackMemStoreKey,
		"DelayedackParams",
	)
	delayedackKeeper := delayedackkeeper.NewKeeper(
		cdc,
		delayedackStoreKey,
		delayedackParams,
		rollappKeeper,
		nil,
		channelKeeper,
		nil,
	)
	delayedackKeeper.SetParams(ctx, delayedacktypes.DefaultParams())

	return ctx, rollappKeeper, delayedackKeeper
}

func skipDelayTransferGenesisPacket(memo string) channeltypes.Packet {
	data := transfertypes.NewFungibleTokenPacketData(
		"stake",
		"1",
		"sender",
		"receiver",
		memo,
	)

	return channeltypes.Packet{
		Sequence:           1,
		SourcePort:         transfertypes.PortID,
		SourceChannel:      "channel-1",
		DestinationPort:    transfertypes.PortID,
		DestinationChannel: skipDelayTestChannelID,
		Data:               data.GetBytes(),
		TimeoutHeight:      clienttypes.NewHeight(0, 1),
	}
}
