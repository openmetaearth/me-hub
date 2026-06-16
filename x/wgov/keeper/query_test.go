package keeper_test

import (
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/x/wgov/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T(), false)
	ctx := app.GetBaseApp().NewContext(false, tmproto.Header{}).WithChainID(apptesting.TestChainID)

	s.Require().NoError(app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams()))
	s.Require().NoError(app.BankKeeper.SetParams(ctx, banktypes.DefaultParams()))
	s.Require().NoError(app.GovKeeper.SetParams(ctx, govv1.DefaultParams()))

	s.App = app
	s.Ctx = ctx
}

func (s *KeeperTestSuite) TestMeTallyResultDoesNotDeleteVotesForActiveProposal() {
	voter := s.NewAccounts(1)[0]
	now := s.Ctx.BlockTime()
	votingEnd := now.Add(time.Hour)
	emptyTally := govv1.EmptyTallyResult()
	const proposalID uint64 = 230

	s.App.GovKeeper.SetProposal(s.Ctx, govv1.Proposal{
		Id:               proposalID,
		Status:           govv1.StatusVotingPeriod,
		FinalTallyResult: &emptyTally,
		SubmitTime:       &now,
		DepositEndTime:   &votingEnd,
		TotalDeposit:     sdk.Coins{},
		VotingStartTime:  &now,
		VotingEndTime:    &votingEnd,
		Title:            "active proposal",
		Summary:          "active proposal",
		Proposer:         voter.String(),
	})

	vote := govv1.NewVote(proposalID, voter, govv1.NewNonSplitVoteOption(govv1.OptionYes), "")
	s.App.GovKeeper.SetVote(s.Ctx, vote)

	_, found := s.App.GovKeeper.GetVote(s.Ctx, proposalID, voter)
	s.Require().True(found)

	_, err := s.App.GovKeeper.MeTallyResult(sdk.WrapSDKContext(s.Ctx), &types.QueryMeTallyResultRequest{
		ProposalId: proposalID,
	})
	s.Require().NoError(err)

	_, found = s.App.GovKeeper.GetVote(s.Ctx, proposalID, voter)
	s.Require().True(found, "MeTallyResult is a query and must not delete stored votes")
}
