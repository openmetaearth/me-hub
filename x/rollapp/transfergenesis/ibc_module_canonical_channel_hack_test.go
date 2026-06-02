package transfergenesis

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/openmetaearth/me-hub/x/rollapp/keeper"
	"github.com/stretchr/testify/require"
)

type passThroughIBCModule struct {
	porttypes.IBCModule
	called bool
}

func (m *passThroughIBCModule) OnRecvPacket(
	_ sdk.Context,
	_ channeltypes.Packet,
	_ sdk.AccAddress,
) exported.Acknowledgement {
	m.called = true
	return channeltypes.NewResultAcknowledgement([]byte("ok"))
}

type panicChannelKeeper struct{}

func (panicChannelKeeper) GetChannelClientState(sdk.Context, string, string) (string, exported.ClientState, error) {
	panic("canonical channel hack must not derive rollapp channel from incoming packets")
}

func TestCanonicalChannelHackDoesNotSetChannelFromFirstPacket(t *testing.T) {
	next := &passThroughIBCModule{}
	module := NewIBCModuleCanonicalChannelHack(next, keeper.Keeper{}, panicChannelKeeper{})

	ack := module.OnRecvPacket(
		sdk.Context{},
		channeltypes.Packet{
			DestinationPort:    "transfer",
			DestinationChannel: "channel-attacker",
		},
		sdk.AccAddress{},
	)

	require.True(t, next.called)
	require.True(t, ack.Success())
}
