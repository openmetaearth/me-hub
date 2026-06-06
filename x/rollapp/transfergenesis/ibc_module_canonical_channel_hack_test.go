package transfergenesis_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/stretchr/testify/require"

	keepertest "github.com/openmetaearth/me-hub/testutil/keeper"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
	"github.com/openmetaearth/me-hub/x/rollapp/transfergenesis"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

type canonicalChannelTestChannelKeeper struct {
	chainID string
}

func (k canonicalChannelTestChannelKeeper) GetChannelClientState(
	sdk.Context,
	string,
	string,
) (string, exported.ClientState, error) {
	return "client-0", &ibctmtypes.ClientState{ChainId: k.chainID}, nil
}

type canonicalChannelTestIBCModule struct {
	porttypes.IBCModule
	keeper        rollappkeeper.Keeper
	rollappID     string
	ack           exported.Acknowledgement
	seenChannel   string
	recvWasCalled bool
}

func (m *canonicalChannelTestIBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	m.recvWasCalled = true
	rollapp, found := m.keeper.GetRollapp(ctx, m.rollappID)
	if found {
		m.seenChannel = rollapp.ChannelId
	}
	return m.ack
}

func TestCanonicalChannelHackCommitsChannelOnlyOnSuccessfulAck(t *testing.T) {
	const (
		rollappID = "rollapp_1234-1"
		channelID = "channel-attacker"
	)

	testCases := []struct {
		name              string
		ack               exported.Acknowledgement
		expectedCommitted string
		expectedSuccess   bool
	}{
		{
			name:              "error acknowledgement leaves canonical channel unset",
			ack:               channeltypes.NewErrorAcknowledgement(fmt.Errorf("downstream rejected")),
			expectedCommitted: "",
			expectedSuccess:   false,
		},
		{
			name:              "success acknowledgement commits canonical channel",
			ack:               channeltypes.NewResultAcknowledgement([]byte{1}),
			expectedCommitted: channelID,
			expectedSuccess:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keeperPtr, ctx := keepertest.RollappKeeper(t)
			keeper := *keeperPtr

			keeper.SetRollapp(ctx, rollapptypes.Rollapp{RollappId: rollappID})

			next := &canonicalChannelTestIBCModule{
				keeper:    keeper,
				rollappID: rollappID,
				ack:       tc.ack,
			}
			module := transfergenesis.NewIBCModuleCanonicalChannelHack(
				next,
				keeper,
				canonicalChannelTestChannelKeeper{chainID: rollappID},
			)

			ack := module.OnRecvPacket(ctx, channeltypes.Packet{
				DestinationPort:    "transfer",
				DestinationChannel: channelID,
			}, sdk.AccAddress{})

			require.Equal(t, tc.expectedSuccess, ack.Success())
			require.True(t, next.recvWasCalled)
			require.Equal(t, channelID, next.seenChannel)

			rollapp, found := keeper.GetRollapp(ctx, rollappID)
			require.True(t, found)
			require.Equal(t, tc.expectedCommitted, rollapp.ChannelId)
		})
	}
}
