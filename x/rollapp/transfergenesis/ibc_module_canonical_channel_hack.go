package transfergenesis

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	rollappkeeper "github.com/openmetaearth/me-hub/x/rollapp/keeper"
)

/*
TODO: this whole file is temporary
	Prior to this we relied on the whitelist addr to set the canonical channel, but that is no longer possible
	This currently file is a hack (not secure)
	The real solution will come in a followup PR
*/

type ChannelKeeper interface {
	GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error) // implemented by ibc channel keeper
}

type IBCModuleCanonicalChannelHack struct {
	porttypes.IBCModule // next one
}

func NewIBCModuleCanonicalChannelHack(
	next porttypes.IBCModule,
	_ rollappkeeper.Keeper,
	_ ChannelKeeper,
) *IBCModuleCanonicalChannelHack {
	return &IBCModuleCanonicalChannelHack{IBCModule: next}
}

func (w IBCModuleCanonicalChannelHack) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	// Canonical channels must be set by an authorized UpdateRollapp call, not by
	// the first inbound packet observed on an arbitrary channel.
	return w.IBCModule.OnRecvPacket(ctx, packet, relayer)
}
