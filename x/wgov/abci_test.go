package wgov_test

import (
	"testing"
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/x/wgov"
)

type ABCITestSuite struct {
	apptesting.KeeperTestHelper
}

func TestABCITestSuite(t *testing.T) {
	suite.Run(t, new(ABCITestSuite))
}

func (s *ABCITestSuite) SetupTest() {
	app := apptesting.Setup(s.T())
	ctx := app.GetBaseApp().NewContext(false)
	s.App = app
	s.Ctx = ctx
}

// TestEndBlocker_InactiveProposal tests that inactive proposals are deleted after deposit period ends
func (s *ABCITestSuite) TestEndBlocker_InactiveProposal() {
	ctx := s.Ctx
	govKeeper := s.App.GovKeeper

	// Create a proposal
	proposer := s.TestAccs()[0]
	proposal, err := v1.NewProposal(
		[]sdk.Msg{},
		1,
		ctx.BlockTime(),
		ctx.BlockTime().Add(48*time.Hour),
		"",
		"Test Proposal",
		"Test proposal description",
		proposer,
		false,
	)
	require.NoError(s.T(), err)
	proposal.Status = v1.StatusDepositPeriod
	depositEndTime := ctx.BlockTime().Add(24 * time.Hour)
	proposal.DepositEndTime = &depositEndTime

	err = govKeeper.SetProposal(ctx, proposal)
	require.NoError(s.T(), err)

	// Add to inactive proposals queue
	err = govKeeper.InactiveProposalsQueue.Set(ctx, collections.Join(depositEndTime, proposal.Id), proposal.Id)
	require.NoError(s.T(), err)

	// Move time forward past deposit end time
	ctx = ctx.WithBlockTime(depositEndTime.Add(1 * time.Second))

	// Run EndBlocker
	err = wgov.EndBlocker(ctx, govKeeper)
	require.NoError(s.T(), err)

	// Verify proposal is deleted
	_, err = govKeeper.Proposals.Get(ctx, proposal.Id)
	require.Error(s.T(), err)
}

// TestEndBlocker_ActiveProposalPassed tests that passed proposals are executed
func (s *ABCITestSuite) TestEndBlocker_ActiveProposalPassed() {
	ctx := s.Ctx
	govKeeper := s.App.GovKeeper

	// Get validators for voting  - need multiple validators to reach quorum
	validators, err := s.App.StakingKeeper.GetAllValidators(ctx)
	require.NoError(s.T(), err)
	require.True(s.T(), len(validators) >= 2, "need at least 2 validators")

	votingEndTime := ctx.BlockTime().Add(24 * time.Hour)

	// Create a proposal with empty messages (safe to execute)
	proposer := sdk.AccAddress([]byte("proposer"))
	proposal, err := v1.NewProposal(
		[]sdk.Msg{},
		1,
		ctx.BlockTime(),
		ctx.BlockTime().Add(48*time.Hour),
		"",
		"Test Proposal",
		"Test proposal description",
		proposer,
		false,
	)
	require.NoError(s.T(), err)

	proposal.VotingEndTime = &votingEndTime
	votingStartTime := ctx.BlockTime()
	proposal.VotingStartTime = &votingStartTime
	proposal.Status = v1.StatusVotingPeriod

	err = govKeeper.SetProposal(ctx, proposal)
	require.NoError(s.T(), err)

	// Add to active proposals queue
	err = govKeeper.ActiveProposalsQueue.Set(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id), proposal.Id)
	require.NoError(s.T(), err)

	// Cast yes votes from all validators to reach quorum
	for _, validator := range validators {
		valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		require.NoError(s.T(), err)
		voterAddr := sdk.AccAddress(valAddr)

		err = govKeeper.AddVote(ctx, proposal.Id, voterAddr, v1.NewNonSplitVoteOption(v1.OptionYes), "")
		require.NoError(s.T(), err)
	}

	// Move time forward past voting end time
	ctx = ctx.WithBlockTime(votingEndTime.Add(1 * time.Second))

	// Run EndBlocker
	err = wgov.EndBlocker(ctx, govKeeper)
	require.NoError(s.T(), err)

	// Verify proposal passed
	p, err := govKeeper.Proposals.Get(ctx, proposal.Id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), v1.StatusPassed, p.Status)
}

// TestEndBlocker_ActiveProposalRejected tests that rejected proposals are marked correctly
func (s *ABCITestSuite) TestEndBlocker_ActiveProposalRejected() {
	ctx := s.Ctx
	govKeeper := s.App.GovKeeper

	// Get validator for voting
	validators, err := s.App.StakingKeeper.GetAllValidators(ctx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), validators)
	validator := validators[0]

	valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
	require.NoError(s.T(), err)
	voterAddr := sdk.AccAddress(valAddr)

	// Create a proposal
	proposal, err := v1.NewProposal(
		[]sdk.Msg{},
		1,
		ctx.BlockTime(),
		ctx.BlockTime().Add(48*time.Hour),
		"",
		"Test Proposal",
		"Test proposal description",
		voterAddr,
		false,
	)
	require.NoError(s.T(), err)

	votingEndTime := ctx.BlockTime().Add(24 * time.Hour)
	proposal.VotingEndTime = &votingEndTime
	votingStartTime := ctx.BlockTime()
	proposal.VotingStartTime = &votingStartTime
	proposal.Status = v1.StatusVotingPeriod

	err = govKeeper.SetProposal(ctx, proposal)
	require.NoError(s.T(), err)

	// Add to active proposals queue
	err = govKeeper.ActiveProposalsQueue.Set(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id), proposal.Id)
	require.NoError(s.T(), err)

	// Cast a no vote from validator
	err = govKeeper.AddVote(ctx, proposal.Id, voterAddr, v1.NewNonSplitVoteOption(v1.OptionNo), "")
	require.NoError(s.T(), err)

	// Move time forward past voting end time
	ctx = ctx.WithBlockTime(votingEndTime.Add(1 * time.Second))

	// Run EndBlocker
	err = wgov.EndBlocker(ctx, govKeeper)
	require.NoError(s.T(), err)

	// Verify proposal rejected
	p, err := govKeeper.Proposals.Get(ctx, proposal.Id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), v1.StatusRejected, p.Status)
	require.Equal(s.T(), "proposal did not get enough votes to pass", p.FailedReason)
}

// TestEndBlocker_ExpeditedProposalConvertsToRegular tests expedited proposal conversion
func (s *ABCITestSuite) TestEndBlocker_ExpeditedProposalConvertsToRegular() {
	ctx := s.Ctx
	govKeeper := s.App.GovKeeper

	// Get validator for voting
	validators, err := s.App.StakingKeeper.GetAllValidators(ctx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), validators)
	validator := validators[0]

	valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
	require.NoError(s.T(), err)
	voterAddr := sdk.AccAddress(valAddr)

	// Create an expedited proposal
	proposal, err := v1.NewProposal(
		[]sdk.Msg{},
		1,
		ctx.BlockTime(),
		ctx.BlockTime().Add(24*time.Hour), // Shorter voting period for expedited
		"",
		"Expedited Test Proposal",
		"Expedited test proposal description",
		voterAddr,
		true, // Expedited
	)
	require.NoError(s.T(), err)

	votingEndTime := ctx.BlockTime().Add(24 * time.Hour)
	proposal.VotingEndTime = &votingEndTime
	votingStartTime := ctx.BlockTime()
	proposal.VotingStartTime = &votingStartTime
	proposal.Status = v1.StatusVotingPeriod

	err = govKeeper.SetProposal(ctx, proposal)
	require.NoError(s.T(), err)

	// Add to active proposals queue
	err = govKeeper.ActiveProposalsQueue.Set(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id), proposal.Id)
	require.NoError(s.T(), err)

	// Cast a no vote (expedited proposal fails)
	err = govKeeper.AddVote(ctx, proposal.Id, voterAddr, v1.NewNonSplitVoteOption(v1.OptionNo), "")
	require.NoError(s.T(), err)

	// Move time forward past voting end time
	ctx = ctx.WithBlockTime(votingEndTime.Add(1 * time.Second))

	// Run EndBlocker
	err = wgov.EndBlocker(ctx, govKeeper)
	require.NoError(s.T(), err)

	// Verify proposal converted to regular
	p, err := govKeeper.Proposals.Get(ctx, proposal.Id)
	require.NoError(s.T(), err)
	require.False(s.T(), p.Expedited, "proposal should no longer be expedited")
	require.Equal(s.T(), v1.StatusVotingPeriod, p.Status, "proposal should still be voting")
	require.True(s.T(), p.VotingEndTime.After(votingEndTime), "voting period should be extended")
}

// TestEndBlocker_NoProposals tests that EndBlocker works with no proposals
func (s *ABCITestSuite) TestEndBlocker_NoProposals() {
	ctx := s.Ctx
	govKeeper := s.App.GovKeeper

	// Run EndBlocker with no proposals
	err := wgov.EndBlocker(ctx, govKeeper)
	require.NoError(s.T(), err)
}

// TestEndBlocker_ProposalWithVeto tests proposal rejection with veto
func (s *ABCITestSuite) TestEndBlocker_ProposalWithVeto() {
	ctx := s.Ctx
	govKeeper := s.App.GovKeeper

	// Get validator for voting
	validators, err := s.App.StakingKeeper.GetAllValidators(ctx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), validators)
	validator := validators[0]

	valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
	require.NoError(s.T(), err)
	voterAddr := sdk.AccAddress(valAddr)

	// Create a proposal
	proposal, err := v1.NewProposal(
		[]sdk.Msg{},
		1,
		ctx.BlockTime(),
		ctx.BlockTime().Add(48*time.Hour),
		"",
		"Test Proposal",
		"Test proposal description",
		voterAddr,
		false,
	)
	require.NoError(s.T(), err)

	votingEndTime := ctx.BlockTime().Add(24 * time.Hour)
	proposal.VotingEndTime = &votingEndTime
	votingStartTime := ctx.BlockTime()
	proposal.VotingStartTime = &votingStartTime
	proposal.Status = v1.StatusVotingPeriod

	err = govKeeper.SetProposal(ctx, proposal)
	require.NoError(s.T(), err)

	// Add to active proposals queue
	err = govKeeper.ActiveProposalsQueue.Set(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id), proposal.Id)
	require.NoError(s.T(), err)

	// Cast a no with veto vote
	err = govKeeper.AddVote(ctx, proposal.Id, voterAddr, v1.NewNonSplitVoteOption(v1.OptionNoWithVeto), "")
	require.NoError(s.T(), err)

	// Move time forward past voting end time
	ctx = ctx.WithBlockTime(votingEndTime.Add(1 * time.Second))

	// Run EndBlocker
	err = wgov.EndBlocker(ctx, govKeeper)
	require.NoError(s.T(), err)

	// Verify proposal rejected
	p, err := govKeeper.Proposals.Get(ctx, proposal.Id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), v1.StatusRejected, p.Status)
}

// TestEndBlocker_MultipleProposals tests handling multiple proposals at once
func (s *ABCITestSuite) TestEndBlocker_MultipleProposals() {
	ctx := s.Ctx
	govKeeper := s.App.GovKeeper

	// Get validators for voting - need all to vote to reach quorum
	validators, err := s.App.StakingKeeper.GetAllValidators(ctx)
	require.NoError(s.T(), err)
	require.True(s.T(), len(validators) >= 2)

	votingEndTime := ctx.BlockTime().Add(24 * time.Hour)
	proposer := sdk.AccAddress([]byte("proposer"))

	// Create multiple proposals
	for i := uint64(1); i <= 3; i++ {
		proposal, err := v1.NewProposal(
			[]sdk.Msg{},
			i,
			ctx.BlockTime(),
			ctx.BlockTime().Add(48*time.Hour),
			"",
			"Test Proposal",
			"Test proposal description",
			proposer,
			false,
		)
		require.NoError(s.T(), err)

		proposal.VotingEndTime = &votingEndTime
		votingStartTime := ctx.BlockTime()
		proposal.VotingStartTime = &votingStartTime
		proposal.Status = v1.StatusVotingPeriod

		err = govKeeper.SetProposal(ctx, proposal)
		require.NoError(s.T(), err)

		// Add to active proposals queue
		err = govKeeper.ActiveProposalsQueue.Set(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id), proposal.Id)
		require.NoError(s.T(), err)

		// All validators vote yes on proposal to reach quorum
		for _, validator := range validators {
			valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
			require.NoError(s.T(), err)
			voterAddr := sdk.AccAddress(valAddr)

			err = govKeeper.AddVote(ctx, proposal.Id, voterAddr, v1.NewNonSplitVoteOption(v1.OptionYes), "")
			require.NoError(s.T(), err)
		}
	}

	// Move time forward past voting end time
	ctx = ctx.WithBlockTime(votingEndTime.Add(1 * time.Second))

	// Run EndBlocker
	err = wgov.EndBlocker(ctx, govKeeper)
	require.NoError(s.T(), err)

	// Verify all proposals passed
	for i := uint64(1); i <= 3; i++ {
		p, err := govKeeper.Proposals.Get(ctx, i)
		require.NoError(s.T(), err)
		require.Equal(s.T(), v1.StatusPassed, p.Status)
	}
}

// TestEndBlocker_InsufficientDeposit tests proposal deletion when deposit is insufficient
func (s *ABCITestSuite) TestEndBlocker_InsufficientDeposit() {
	ctx := s.Ctx
	govKeeper := s.App.GovKeeper

	proposer := s.TestAccs()[0]

	// Create a proposal with insufficient deposit
	proposal, err := v1.NewProposal(
		[]sdk.Msg{},
		1,
		ctx.BlockTime(),
		ctx.BlockTime().Add(48*time.Hour),
		"",
		"Test Proposal",
		"Test proposal description",
		proposer,
		false,
	)
	require.NoError(s.T(), err)
	proposal.Status = v1.StatusDepositPeriod
	depositEndTime := ctx.BlockTime().Add(24 * time.Hour)
	proposal.DepositEndTime = &depositEndTime

	// Set total deposit to zero (insufficient)
	proposal.TotalDeposit = sdk.NewCoins()

	err = govKeeper.SetProposal(ctx, proposal)
	require.NoError(s.T(), err)

	// Add to inactive proposals queue
	err = govKeeper.InactiveProposalsQueue.Set(ctx, collections.Join(*proposal.DepositEndTime, proposal.Id), proposal.Id)
	require.NoError(s.T(), err)

	// Move time forward past deposit end time
	ctx = ctx.WithBlockTime(depositEndTime.Add(1 * time.Second))

	// Run EndBlocker
	err = wgov.EndBlocker(ctx, govKeeper)
	require.NoError(s.T(), err)

	// Verify proposal is deleted
	_, err = govKeeper.Proposals.Get(ctx, proposal.Id)
	require.Error(s.T(), err)
}

func (s *ABCITestSuite) TestAccs() []sdk.AccAddress {
	if len(s.KeeperTestHelper.TestAccs) == 0 {
		s.KeeperTestHelper.TestAccs = s.NewAccounts(3)
	}
	return s.KeeperTestHelper.TestAccs
}
