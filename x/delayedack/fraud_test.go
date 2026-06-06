package delayedack_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/stretchr/testify/require"

	keepertest "github.com/openmetaearth/me-hub/testutil/keeper"
	commontypes "github.com/openmetaearth/me-hub/x/common/types"
	delayedacktypes "github.com/openmetaearth/me-hub/x/delayedack/types"
)

func TestHandleFraudKeepsOutgoingPacketPendingWhenRefundFails(t *testing.T) {
	keeper, ctx := keepertest.DelayedackKeeper(t)
	rollappID := "rollapp-1"
	rollappPacket := commontypes.RollappPacket{
		RollappId: rollappID,
		Packet: &channeltypes.Packet{
			SourcePort:         "transfer",
			SourceChannel:      "channel-0",
			DestinationPort:    "transfer",
			DestinationChannel: "channel-1",
			Data:               []byte("packet-data"),
			Sequence:           1,
		},
		Status:      commontypes.Status_PENDING,
		Type:        commontypes.RollappPacket_ON_TIMEOUT,
		Relayer:     sdk.AccAddress("relayer"),
		ProofHeight: 10,
	}
	keeper.SetRollappPacket(ctx, rollappPacket)

	ibc := &timeoutFailingIBCModule{err: errors.New("refund failed")}
	err := keeper.HandleFraud(ctx, rollappID, ibc)

	require.ErrorContains(t, err, "refund reverted packet")
	require.Equal(t, 1, ibc.timeoutCalls)
	require.Len(t, keeper.ListRollappPackets(ctx, delayedacktypes.ByRollappIDByStatus(rollappID, commontypes.Status_PENDING)), 1)
	require.Empty(t, keeper.ListRollappPackets(ctx, delayedacktypes.ByRollappIDByStatus(rollappID, commontypes.Status_REVERTED)))
}

var _ porttypes.IBCModule = (*timeoutFailingIBCModule)(nil)

type timeoutFailingIBCModule struct {
	err          error
	timeoutCalls int
}

func (m *timeoutFailingIBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	return version, nil
}

func (m *timeoutFailingIBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	return counterpartyVersion, nil
}

func (m *timeoutFailingIBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID string,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	return nil
}

func (m *timeoutFailingIBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID string,
	channelID string,
) error {
	return nil
}

func (m *timeoutFailingIBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID string,
	channelID string,
) error {
	return nil
}

func (m *timeoutFailingIBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID string,
	channelID string,
) error {
	return nil
}

func (m *timeoutFailingIBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	return channeltypes.NewResultAcknowledgement([]byte{1})
}

func (m *timeoutFailingIBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return nil
}

func (m *timeoutFailingIBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	m.timeoutCalls++
	return m.err
}

func (m *timeoutFailingIBCModule) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (uint64, error) {
	return 0, nil
}

func (m *timeoutFailingIBCModule) WriteAcknowledgement(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet exported.PacketI,
	ack exported.Acknowledgement,
) error {
	return nil
}

func (m *timeoutFailingIBCModule) GetAppVersion(
	ctx sdk.Context,
	portID string,
	channelID string,
) (string, bool) {
	return "", true
}
