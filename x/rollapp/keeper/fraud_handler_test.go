package keeper_test

import (
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	tmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	common "github.com/openmetaearth/me-hub/x/common/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

// Happy Flow
// - frozen rollapp
// - slashed sequecner and unbonded all other sequencers
// - reverted states
// - cleared queue

func (suite *RollappTestSuite) TestHandleFraud() {
	var err error
	suite.SetupTest()
	ctx := &suite.Ctx
	keeper := suite.App.RollappKeeper
	initialheight := uint64(10)
	suite.Ctx = suite.Ctx.WithBlockHeight(int64(initialheight))

	numOfSequencers := uint64(3)
	numOfStates := uint64(100)
	numOfBlocks := uint64(10)
	fraudHeight := uint64(300)

	// unrelated rollapp just to validate it's unaffected
	rollapp_2 := suite.CreateDefaultRollapp()
	proposer_2 := suite.CreateDefaultSequencer(*ctx, rollapp_2)

	// create rollapp and sequencers before fraud evidence
	rollapp := suite.CreateDefaultRollapp()
	proposer := suite.CreateDefaultSequencer(*ctx, rollapp)
	for i := uint64(0); i < numOfSequencers-1; i++ {
		suite.CreateDefaultSequencer(*ctx, rollapp)
	}

	// send state updates
	var lastHeight uint64 = 1

	for i := uint64(0); i < numOfStates; i++ {
		_, err = suite.PostStateUpdate(*ctx, rollapp, proposer, lastHeight, numOfBlocks)
		suite.Require().Nil(err)

		lastHeight, err = suite.PostStateUpdate(*ctx, rollapp_2, proposer_2, lastHeight, numOfBlocks)
		suite.Require().Nil(err)

		suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeader().Height + 1)
	}

	// finalize some of the states
	suite.App.RollappKeeper.FinalizeRollappStates(suite.Ctx.WithBlockHeight(20))

	// assert before fraud submission
	suite.assertBeforeFraud(rollapp, fraudHeight)

	err = keeper.HandleFraud(*ctx, rollapp, "", fraudHeight, proposer)
	suite.Require().Nil(err)

	suite.assertFraudHandled(rollapp)
}

// Fail - Invalid rollapp
func (suite *RollappTestSuite) TestHandleFraud_InvalidRollapp() {
	suite.SetupTest()
	ctx := &suite.Ctx
	keeper := suite.App.RollappKeeper

	rollapp := suite.CreateDefaultRollapp()
	proposer := suite.CreateDefaultSequencer(*ctx, rollapp)
	_, err := suite.PostStateUpdate(*ctx, rollapp, proposer, 1, uint64(10))
	suite.Require().Nil(err)

	err = keeper.HandleFraud(*ctx, "invalidRollapp", "", 2, proposer)
	suite.Require().NotNil(err)
}

// Fail - Wrong height
func (suite *RollappTestSuite) TestHandleFraud_WrongHeight() {
	suite.SetupTest()
	ctx := &suite.Ctx
	keeper := suite.App.RollappKeeper

	rollapp := suite.CreateDefaultRollapp()
	proposer := suite.CreateDefaultSequencer(*ctx, rollapp)
	_, err := suite.PostStateUpdate(*ctx, rollapp, proposer, 1, uint64(10))
	suite.Require().Nil(err)

	err = keeper.HandleFraud(*ctx, rollapp, "", 100, proposer)
	suite.Require().NotNil(err)
}

// Fail - Wrong sequencer address
func (suite *RollappTestSuite) TestHandleFraud_WrongSequencer() {
	suite.SetupTest()
	ctx := &suite.Ctx
	keeper := suite.App.RollappKeeper

	rollapp := suite.CreateDefaultRollapp()
	proposer := suite.CreateDefaultSequencer(*ctx, rollapp)
	_, err := suite.PostStateUpdate(*ctx, rollapp, proposer, 1, uint64(10))
	suite.Require().Nil(err)

	err = keeper.HandleFraud(*ctx, rollapp, "", 2, "wrongSequencer")
	suite.Require().NotNil(err)
}

// Fail - Wrong channel-ID
func (suite *RollappTestSuite) TestHandleFraud_WrongChannelID() {
	suite.SetupTest()
	ctx := &suite.Ctx
	keeper := suite.App.RollappKeeper

	rollapp := suite.CreateDefaultRollapp()
	proposer := suite.CreateDefaultSequencer(*ctx, rollapp)
	_, err := suite.PostStateUpdate(*ctx, rollapp, proposer, 1, uint64(10))
	suite.Require().Nil(err)

	err = keeper.HandleFraud(*ctx, rollapp, "wrongChannelID", 2, proposer)
	suite.Require().NotNil(err)
}

// Fail - Disputing already reverted state
func (suite *RollappTestSuite) TestHandleFraud_AlreadyReverted() {
	suite.SetupTest()
	var err error
	ctx := &suite.Ctx
	keeper := suite.App.RollappKeeper
	numOfSequencers := uint64(3)
	numOfStates := uint64(10)

	rollapp := suite.CreateDefaultRollapp()
	proposer := suite.CreateDefaultSequencer(*ctx, rollapp)
	for i := uint64(0); i < numOfSequencers-1; i++ {
		suite.CreateDefaultSequencer(*ctx, rollapp)
	}

	// send state updates
	var lastHeight uint64 = 1
	for i := uint64(0); i < numOfStates; i++ {
		lastHeight, err = suite.PostStateUpdate(*ctx, rollapp, proposer, lastHeight, uint64(10))
		suite.Require().Nil(err)

		suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeader().Height + 1)
	}

	err = keeper.HandleFraud(*ctx, rollapp, "", 11, proposer)
	suite.Require().Nil(err)

	err = keeper.HandleFraud(*ctx, rollapp, "", 1, proposer)
	suite.Require().NotNil(err)
}

// Fail - Disputing already finalized state
func (suite *RollappTestSuite) TestHandleFraud_AlreadyFinalized() {
	suite.SetupTest()
	ctx := &suite.Ctx
	keeper := suite.App.RollappKeeper

	rollapp := suite.CreateDefaultRollapp()
	proposer := suite.CreateDefaultSequencer(*ctx, rollapp)
	_, err := suite.PostStateUpdate(*ctx, rollapp, proposer, 1, uint64(10))
	suite.Require().Nil(err)

	// finalize state
	suite.Ctx = suite.Ctx.WithBlockHeight(ctx.BlockHeight() + int64(keeper.DisputePeriodInBlocks(*ctx)))
	suite.App.RollappKeeper.FinalizeRollappStates(suite.Ctx)
	stateInfo, err := suite.App.RollappKeeper.FindStateInfoByHeight(suite.Ctx, rollapp, 1)
	suite.Require().Nil(err)
	suite.Require().Equal(common.Status_FINALIZED, stateInfo.Status)

	err = keeper.HandleFraud(*ctx, rollapp, "", 2, proposer)
	suite.Require().NotNil(err)
}

func (suite *RollappTestSuite) TestHandleFraudFreezesIBCClientAtLatestHeight() {
	suite.SetupTest()
	keeper := suite.App.RollappKeeper

	rollappId := suite.CreateDefaultRollapp()
	proposer := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	_, err := suite.PostStateUpdate(suite.Ctx, rollappId, proposer, 1, uint64(10))
	suite.Require().NoError(err)

	const (
		clientId     = "07-tendermint-0"
		connectionId = "connection-0"
		channelId    = "channel-0"
	)

	latestHeight := clienttypes.NewHeight(7, 12345)
	clientState := tmtypes.NewClientState(
		"rollapp_1234-7",
		tmtypes.DefaultTrustLevel,
		ibctesting.TrustingPeriod,
		ibctesting.UnbondingPeriod,
		ibctesting.MaxClockDrift,
		latestHeight,
		commitmenttypes.GetSDKSpecs(),
		ibctesting.UpgradePath,
	)
	suite.App.IBCKeeper.ClientKeeper.SetClientState(suite.Ctx, clientId, clientState)

	connection := connectiontypes.NewConnectionEnd(
		connectiontypes.OPEN,
		clientId,
		connectiontypes.NewCounterparty("counterparty-client-0", "connection-1", commitmenttypes.NewMerklePrefix([]byte("ibc"))),
		[]*connectiontypes.Version{connectiontypes.DefaultIBCVersion},
		0,
	)
	suite.App.IBCKeeper.ConnectionKeeper.SetConnection(suite.Ctx, connectionId, connection)

	channel := channeltypes.NewChannel(
		channeltypes.OPEN,
		channeltypes.UNORDERED,
		channeltypes.NewCounterparty("transfer", "channel-1"),
		[]string{connectionId},
		"ics20-1",
	)
	suite.App.IBCKeeper.ChannelKeeper.SetChannel(suite.Ctx, "transfer", channelId, channel)

	rollapp, found := suite.App.RollappKeeper.GetRollapp(suite.Ctx, rollappId)
	suite.Require().True(found)
	rollapp.ChannelId = channelId
	suite.App.RollappKeeper.SetRollapp(suite.Ctx, rollapp)

	err = keeper.HandleFraud(suite.Ctx, rollappId, clientId, 2, proposer)
	suite.Require().NoError(err)

	storedClientState, found := suite.App.IBCKeeper.ClientKeeper.GetClientState(suite.Ctx, clientId)
	suite.Require().True(found)
	tmClientState, ok := storedClientState.(*tmtypes.ClientState)
	suite.Require().True(ok)
	suite.Require().Equal(latestHeight, tmClientState.FrozenHeight)
}

/* ---------------------------------- utils --------------------------------- */

// assert before fraud submission, to validate the Test itself
func (suite *RollappTestSuite) assertBeforeFraud(rollappId string, height uint64) {
	rollapp, found := suite.App.RollappKeeper.GetRollapp(suite.Ctx, rollappId)
	suite.Require().True(found)
	suite.Require().False(rollapp.Frozen)

	// check sequencers
	sequencers := suite.App.SequencerKeeper.GetSequencersByRollapp(suite.Ctx, rollappId)
	for _, sequencer := range sequencers {
		suite.Require().Equal(types.Bonded, sequencer.Status)
	}

	// check states
	stateInfo, err := suite.App.RollappKeeper.FindStateInfoByHeight(suite.Ctx, rollappId, height)
	suite.Require().Nil(err)
	suite.Require().Equal(common.Status_PENDING, stateInfo.Status)

	// check queue
	expectedHeight := stateInfo.CreationHeight + suite.App.RollappKeeper.DisputePeriodInBlocks(suite.Ctx)
	queue, found := suite.App.RollappKeeper.GetBlockHeightToFinalizationQueue(suite.Ctx, expectedHeight)
	suite.Require().True(found)

	found = false
	for _, stateInfoIndex := range queue.FinalizationQueue {
		if stateInfoIndex.RollappId == rollappId {
			val, _ := suite.App.RollappKeeper.GetStateInfo(suite.Ctx, rollappId, stateInfoIndex.Index)
			suite.Require().Equal(common.Status_PENDING, val.Status)
			found = true
			break
		}
	}
	suite.Require().True(found)
}

func (suite *RollappTestSuite) assertFraudHandled(rollappId string) {
	rollapp, found := suite.App.RollappKeeper.GetRollapp(suite.Ctx, rollappId)
	suite.Require().True(found)
	suite.Require().True(rollapp.Frozen)

	// check sequencers
	sequencers := suite.App.SequencerKeeper.GetSequencersByRollapp(suite.Ctx, rollappId)
	for _, sequencer := range sequencers {
		suite.Require().Equal(types.Unbonded, sequencer.Status)
	}

	// check states
	finalIdx, _ := suite.App.RollappKeeper.GetLatestFinalizedStateIndex(suite.Ctx, rollappId)
	start := finalIdx.Index + 1
	endIdx, _ := suite.App.RollappKeeper.GetLatestStateInfoIndex(suite.Ctx, rollappId)
	end := endIdx.Index

	for i := start; i <= end; i++ {
		stateInfo, found := suite.App.RollappKeeper.GetStateInfo(suite.Ctx, rollappId, i)
		suite.Require().True(found)
		suite.Require().Equal(common.Status_REVERTED, stateInfo.Status, "state info for height %d is not reverted", stateInfo.StartHeight)
	}

	// check queue
	queue := suite.App.RollappKeeper.GetAllBlockHeightToFinalizationQueue(suite.Ctx)
	suite.Greater(len(queue), 0)
	for _, q := range queue {
		for _, stateInfoIndex := range q.FinalizationQueue {
			suite.Require().NotEqual(rollappId, stateInfoIndex.RollappId)
		}
	}
}
