package keeper_test

import (
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	commontypes "github.com/openmetaearth/me-hub/x/common/types"
	damodule "github.com/openmetaearth/me-hub/x/delayedack"
	"github.com/openmetaearth/me-hub/x/delayedack/types"
)

func (suite *DelayedAckTestSuite) TestHandleFraud() {
	keeper, ctx := suite.App.DelayedAckKeeper, suite.Ctx
	transferStack := damodule.NewIBCMiddleware(
		damodule.WithIBCModule(ibctransfer.NewIBCModule(suite.App.TransferKeeper)),
		damodule.WithKeeper(keeper),
		damodule.WithRollappKeeper(suite.App.RollappKeeper),
	)

	rollappId := "testRollappId"
	pkts := generatePackets(rollappId, 5)
	rollappId2 := "testRollappId2"
	pkts2 := generatePackets(rollappId2, 5)
	prefixPending1 := types.ByRollappIDByStatus(rollappId, commontypes.Status_PENDING)
	prefixPending2 := types.ByRollappIDByStatus(rollappId2, commontypes.Status_PENDING)
	prefixReverted := types.ByRollappIDByStatus(rollappId, commontypes.Status_REVERTED)
	prefixFinalized := types.ByRollappIDByStatus(rollappId, commontypes.Status_FINALIZED)
	prefixFinalized2 := types.ByRollappIDByStatus(rollappId, commontypes.Status_FINALIZED)

	for _, pkt := range append(pkts, pkts2...) {
		keeper.SetRollappPacket(ctx, pkt)
	}

	suite.Require().Equal(5, len(keeper.ListRollappPackets(ctx, prefixPending1)))
	suite.Require().Equal(5, len(keeper.ListRollappPackets(ctx, prefixPending2)))

	// finalize some packets
	_, err := keeper.UpdateRollappPacketWithStatus(ctx, pkts[0], commontypes.Status_FINALIZED)
	suite.Require().Nil(err)
	_, err = keeper.UpdateRollappPacketWithStatus(ctx, pkts2[0], commontypes.Status_FINALIZED)
	suite.Require().Nil(err)

	err = keeper.HandleFraud(ctx, rollappId, 0, transferStack)
	suite.Require().Nil(err)

	suite.Require().Equal(0, len(keeper.ListRollappPackets(ctx, prefixPending1)))
	suite.Require().Equal(4, len(keeper.ListRollappPackets(ctx, prefixPending2)))
	suite.Require().Equal(4, len(keeper.ListRollappPackets(ctx, prefixReverted)))
	suite.Require().Equal(1, len(keeper.ListRollappPackets(ctx, prefixFinalized)))
	suite.Require().Equal(1, len(keeper.ListRollappPackets(ctx, prefixFinalized2)))
}

func (suite *DelayedAckTestSuite) TestHandleFraudKeepsPacketsBeforeFraudHeight() {
	keeper, ctx := suite.App.DelayedAckKeeper, suite.Ctx
	transferStack := damodule.NewIBCMiddleware(
		damodule.WithIBCModule(ibctransfer.NewIBCModule(suite.App.TransferKeeper)),
		damodule.WithKeeper(keeper),
		damodule.WithRollappKeeper(suite.App.RollappKeeper),
	)

	rollappId := "testRollappId"
	fraudHeight := uint64(3)
	pkts := generatePackets(rollappId, 5)
	prefixPending := types.ByRollappIDByStatus(rollappId, commontypes.Status_PENDING)
	prefixReverted := types.ByRollappIDByStatus(rollappId, commontypes.Status_REVERTED)

	for _, pkt := range pkts {
		keeper.SetRollappPacket(ctx, pkt)
	}

	err := keeper.HandleFraud(ctx, rollappId, fraudHeight, transferStack)
	suite.Require().NoError(err)

	pendingPackets := keeper.ListRollappPackets(ctx, prefixPending)
	revertedPackets := keeper.ListRollappPackets(ctx, prefixReverted)
	suite.Require().Len(pendingPackets, 3)
	suite.Require().Len(revertedPackets, 2)

	for _, pkt := range pendingPackets {
		suite.Require().True(pkt.ProofHeight < fraudHeight)
	}
	for _, pkt := range revertedPackets {
		suite.Require().True(pkt.ProofHeight >= fraudHeight)
	}
}

func (suite *DelayedAckTestSuite) TestDeletionOfRevertedPackets() {
	keeper, ctx := suite.App.DelayedAckKeeper, suite.Ctx
	transferStack := damodule.NewIBCMiddleware(
		damodule.WithIBCModule(ibctransfer.NewIBCModule(suite.App.TransferKeeper)),
		damodule.WithKeeper(keeper),
		damodule.WithRollappKeeper(suite.App.RollappKeeper),
	)

	rollappId := "testRollappId"
	pkts := generatePackets(rollappId, 5)
	rollappId2 := "testRollappId2"
	pkts2 := generatePackets(rollappId2, 5)

	for _, pkt := range append(pkts, pkts2...) {
		keeper.SetRollappPacket(ctx, pkt)
	}

	err := keeper.HandleFraud(ctx, rollappId, 0, transferStack)
	suite.Require().Nil(err)

	suite.Require().Equal(10, len(keeper.GetAllRollappPackets(ctx)))

	keeper.SetParams(ctx, types.Params{EpochIdentifier: "minute", BridgingFee: keeper.BridgingFee(ctx)})
	epochHooks := keeper.GetEpochHooks()
	err = epochHooks.AfterEpochEnd(ctx, "minute", 1)
	suite.Require().NoError(err)

	suite.Require().Equal(5, len(keeper.GetAllRollappPackets(ctx)))
}

// TODO: test refunds of pending packets

/* ---------------------------------- utils --------------------------------- */

func generatePackets(rollappId string, num uint64) []commontypes.RollappPacket {
	var packets []commontypes.RollappPacket
	for i := uint64(0); i < num; i++ {
		packets = append(packets, commontypes.RollappPacket{
			RollappId: rollappId,
			Packet: &channeltypes.Packet{
				SourcePort:         "testSourcePort",
				SourceChannel:      "testSourceChannel",
				DestinationPort:    "testDestinationPort",
				DestinationChannel: "testDestinationChannel",
				Data:               []byte("testData"),
				Sequence:           i,
			},
			Status:      commontypes.Status_PENDING,
			ProofHeight: i,
		})
	}
	return packets
}
