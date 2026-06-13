package keeper_test

import (
	"bytes"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	commontypes "github.com/openmetaearth/me-hub/x/common/types"
	dacktypes "github.com/openmetaearth/me-hub/x/delayedack/types"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
)

var _ porttypes.IBCModule = (*countingIBCModule)(nil)

type countingIBCModule struct {
	ackCallbacks     int
	recvCallbacks    int
	timeoutCallbacks int
}

func (m *countingIBCModule) OnChanOpenInit(
	sdk.Context,
	channeltypes.Order,
	[]string,
	string,
	string,
	*capabilitytypes.Capability,
	channeltypes.Counterparty,
	string,
) (string, error) {
	return "", nil
}

func (m *countingIBCModule) OnChanOpenTry(
	sdk.Context,
	channeltypes.Order,
	[]string,
	string,
	string,
	*capabilitytypes.Capability,
	channeltypes.Counterparty,
	string,
) (string, error) {
	return "", nil
}

func (m *countingIBCModule) OnChanOpenAck(sdk.Context, string, string, string, string) error {
	return nil
}

func (m *countingIBCModule) OnChanOpenConfirm(sdk.Context, string, string) error {
	return nil
}

func (m *countingIBCModule) OnChanCloseInit(sdk.Context, string, string) error {
	return nil
}

func (m *countingIBCModule) OnChanCloseConfirm(sdk.Context, string, string) error {
	return nil
}

func (m *countingIBCModule) OnRecvPacket(
	sdk.Context,
	channeltypes.Packet,
	sdk.AccAddress,
) exported.Acknowledgement {
	m.recvCallbacks++
	return channeltypes.NewResultAcknowledgement([]byte{1})
}

func (m *countingIBCModule) OnAcknowledgementPacket(
	sdk.Context,
	channeltypes.Packet,
	[]byte,
	sdk.AccAddress,
) error {
	m.ackCallbacks++
	return nil
}

func (m *countingIBCModule) OnTimeoutPacket(
	sdk.Context,
	channeltypes.Packet,
	sdk.AccAddress,
) error {
	m.timeoutCallbacks++
	return nil
}

func (suite *DelayedAckTestSuite) TestFinalizeRollappPacketsRejectsConsensusRootMismatch() {
	keeper, ctx := suite.App.DelayedAckKeeper, suite.Ctx
	rollappID := "rollapp-root-mismatch-1"
	proofHeight := uint64(7)
	stateRoot := bytes.Repeat([]byte{0x11}, 32)
	consensusRoot := bytes.Repeat([]byte{0x22}, 32)
	packet := newRootCheckedPacket(rollappID, proofHeight, commontypes.RollappPacket_ON_ACK)
	stateInfo := newRootCheckedStateInfo(rollappID, proofHeight, stateRoot)
	suite.setPacketConsensusRoot(ctx, packet, consensusRoot)
	keeper.SetRollappPacket(ctx, packet)

	ibcModule := &countingIBCModule{}
	err := keeper.FinalizeRollappPackets(ctx, ibcModule, rollappID, stateInfo)

	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "state root mismatch")
	suite.Require().Zero(ibcModule.ackCallbacks)
	suite.Require().Len(keeper.ListRollappPackets(ctx, dacktypes.ByRollappIDByStatus(rollappID, commontypes.Status_PENDING)), 1)
	suite.Require().Empty(keeper.ListRollappPackets(ctx, dacktypes.ByRollappIDByStatus(rollappID, commontypes.Status_FINALIZED)))
}

func (suite *DelayedAckTestSuite) TestFinalizeRollappPacketsAllowsMatchingConsensusRoot() {
	keeper, ctx := suite.App.DelayedAckKeeper, suite.Ctx
	rollappID := "rollapp-root-match-1"
	proofHeight := uint64(8)
	stateRoot := bytes.Repeat([]byte{0x33}, 32)
	packet := newRootCheckedPacket(rollappID, proofHeight, commontypes.RollappPacket_ON_ACK)
	stateInfo := newRootCheckedStateInfo(rollappID, proofHeight, stateRoot)
	suite.setPacketConsensusRoot(ctx, packet, stateRoot)
	keeper.SetRollappPacket(ctx, packet)

	ibcModule := &countingIBCModule{}
	err := keeper.FinalizeRollappPackets(ctx, ibcModule, rollappID, stateInfo)

	suite.Require().NoError(err)
	suite.Require().Equal(1, ibcModule.ackCallbacks)
	suite.Require().Empty(keeper.ListRollappPackets(ctx, dacktypes.ByRollappIDByStatus(rollappID, commontypes.Status_PENDING)))
	suite.Require().Len(keeper.ListRollappPackets(ctx, dacktypes.ByRollappIDByStatus(rollappID, commontypes.Status_FINALIZED)), 1)
}

func newRootCheckedPacket(
	rollappID string,
	proofHeight uint64,
	packetType commontypes.RollappPacket_Type,
) commontypes.RollappPacket {
	return commontypes.RollappPacket{
		RollappId: rollappID,
		Packet: &channeltypes.Packet{
			SourcePort:         "transfer",
			SourceChannel:      "channel-root-check",
			DestinationPort:    "transfer",
			DestinationChannel: "channel-hub",
			Data:               []byte("packet-data"),
			Sequence:           proofHeight,
		},
		Status:      commontypes.Status_PENDING,
		ProofHeight: proofHeight,
		Type:        packetType,
		Relayer:     sdk.AccAddress(bytes.Repeat([]byte{0x44}, 20)),
	}
}

func newRootCheckedStateInfo(rollappID string, proofHeight uint64, stateRoot []byte) rollapptypes.StateInfo {
	return rollapptypes.StateInfo{
		StateInfoIndex: rollapptypes.StateInfoIndex{RollappId: rollappID, Index: 1},
		StartHeight:    proofHeight,
		NumBlocks:      1,
		Status:         commontypes.Status_FINALIZED,
		BDs: rollapptypes.BlockDescriptors{BD: []rollapptypes.BlockDescriptor{{
			Height:    proofHeight,
			StateRoot: stateRoot,
		}}},
	}
}

func (suite *DelayedAckTestSuite) setPacketConsensusRoot(
	ctx sdk.Context,
	packet commontypes.RollappPacket,
	consensusRoot []byte,
) {
	clientID := "07-tendermint-364"
	connectionID := "connection-364"
	revisionNumber := uint64(1)
	consensusHeight := clienttypes.NewHeight(revisionNumber, packet.ProofHeight)

	suite.App.IBCKeeper.ClientKeeper.SetClientState(ctx, clientID, ibctmtypes.NewClientState(
		"rollapp-root-check-1",
		ibctmtypes.DefaultTrustLevel,
		time.Hour,
		2*time.Hour,
		time.Minute,
		consensusHeight,
		commitmenttypes.GetSDKSpecs(),
		[]string{"upgrade", "upgradedIBCState"},
	))
	suite.App.IBCKeeper.ClientKeeper.SetClientConsensusState(ctx, clientID, consensusHeight, ibctmtypes.NewConsensusState(
		time.Unix(1, 0).UTC(),
		commitmenttypes.NewMerkleRoot(consensusRoot),
		bytes.Repeat([]byte{0x55}, 32),
	))
	suite.App.IBCKeeper.ConnectionKeeper.SetConnection(ctx, connectionID, connectiontypes.NewConnectionEnd(
		connectiontypes.OPEN,
		clientID,
		connectiontypes.Counterparty{},
		[]*connectiontypes.Version{connectiontypes.DefaultIBCVersion},
		0,
	))
	suite.App.IBCKeeper.ChannelKeeper.SetChannel(ctx,
		packet.Packet.SourcePort,
		packet.Packet.SourceChannel,
		channeltypes.NewChannel(
			channeltypes.OPEN,
			channeltypes.UNORDERED,
			channeltypes.Counterparty{},
			[]string{connectionID},
			"ics20-1",
		),
	)
}
