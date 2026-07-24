package keeper

import (
	"errors"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	tmdb "github.com/cometbft/cometbft-db"
	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	coretypes "github.com/cosmos/ibc-go/v8/modules/core/types"
	"github.com/stretchr/testify/require"

	commontypes "github.com/openmetaearth/me-hub/x/common/types"
)

func TestReplayRecvPacket(t *testing.T) {
	tests := []struct {
		name             string
		ack              exported.Acknowledgement
		expectCommit     bool
		expectErrorEvent bool
	}{
		{
			name:             "success ack commits state",
			ack:              channeltypes.NewResultAcknowledgement([]byte{1}),
			expectCommit:     true,
			expectErrorEvent: false,
		},
		{
			name:             "error ack rolls back state",
			ack:              channeltypes.NewErrorAcknowledgement(errors.New("recv failed")),
			expectCommit:     false,
			expectErrorEvent: true,
		},
		{
			name:             "nil ack commits async state",
			ack:              nil,
			expectCommit:     true,
			expectErrorEvent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, storeKey := newReplayTestContext(t)
			rollappPacket := commontypes.RollappPacket{
				Packet:    &channeltypes.Packet{},
				Relayer:   sdk.AccAddress("relayer"),
				Status:    commontypes.Status_PENDING,
				RollappId: "rollapp_123",
			}

			ibc := &mockRecvIBCModule{
				ack: tt.ack,
				onRecv: func(ctx sdk.Context, _ channeltypes.Packet, _ sdk.AccAddress) {
					ctx.KVStore(storeKey).Set([]byte("committed"), []byte("yes"))
					ctx.EventManager().EmitEvent(sdk.NewEvent("recv_test", sdk.NewAttribute("status", "called")))
				},
			}

			ack := replayRecvPacket(ctx, ibc, rollappPacket)
			if tt.ack == nil {
				require.Nil(t, ack)
			} else {
				require.NotNil(t, ack)
				require.Equal(t, tt.ack.Acknowledgement(), ack.Acknowledgement())
			}

			committed := ctx.KVStore(storeKey).Get([]byte("committed"))
			if tt.expectCommit {
				require.Equal(t, []byte("yes"), committed)
			} else {
				require.Nil(t, committed)
			}

			events := ctx.EventManager().Events()
			require.Len(t, events, 1)
			if tt.expectErrorEvent {
				require.Equal(t, coretypes.ErrorAttributeKeyPrefix+"recv_test", events[0].Type)
				require.Equal(t, coretypes.ErrorAttributeKeyPrefix+"status", events[0].Attributes[0].Key)
			} else {
				require.Equal(t, "recv_test", events[0].Type)
				require.Equal(t, "status", events[0].Attributes[0].Key)
			}
		})
	}
}

func TestConvertToErrorEvents(t *testing.T) {
	events := sdk.Events{
		sdk.NewEvent("test_event",
			sdk.NewAttribute("key1", "value1"),
			sdk.NewAttribute("key2", "value2"),
		),
	}

	converted := convertToErrorEvents(events)
	require.Len(t, converted, 1)
	require.Equal(t, coretypes.ErrorAttributeKeyPrefix+"test_event", converted[0].Type)
	require.Len(t, converted[0].Attributes, 2)
	require.Equal(t, coretypes.ErrorAttributeKeyPrefix+"key1", converted[0].Attributes[0].Key)
	require.Equal(t, "value1", converted[0].Attributes[0].Value)
	require.Equal(t, coretypes.ErrorAttributeKeyPrefix+"key2", converted[0].Attributes[1].Key)
	require.Equal(t, "value2", converted[0].Attributes[1].Value)
}

type mockRecvIBCModule struct {
	porttypes.IBCModule
	ack    exported.Acknowledgement
	onRecv func(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress)
}

func (m *mockRecvIBCModule) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	if m.onRecv != nil {
		m.onRecv(ctx, packet, relayer)
	}

	return m.ack
}

func newReplayTestContext(t *testing.T) (sdk.Context, *storetypes.KVStoreKey) {
	t.Helper()

	storeKey := sdk.NewKVStoreKey("replay-test")
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cometbftproto.Header{}, false, log.NewNopLogger())
	return ctx, storeKey
}
