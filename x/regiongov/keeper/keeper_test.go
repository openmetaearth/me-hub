package keeper_test

import (
	"testing"
	"time"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/regiongov"
	regiongovkeeper "github.com/openmetaearth/me-hub/x/regiongov/keeper"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
	wstakingkeeper "github.com/openmetaearth/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	msgServer           types.MsgServer
	meEarthValidator    stakingtypes.Validator
	experienceValidator stakingtypes.Validator
	usaValidator        stakingtypes.Validator
	startTime           time.Time
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T(), false)
	s.startTime = time.Unix(1_700_000_000, 0).UTC()
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{Time: s.startTime})

	s.App = app
	s.Ctx = ctx
	s.InitializeDao()

	validators := s.App.StakingKeeper.GetValidators(s.Ctx, 10)
	s.Require().GreaterOrEqual(len(validators), 3)
	for _, validator := range validators {
		switch validator.Description.RegionID {
		case wstakingtypes.MeEarthRegionId:
			s.meEarthValidator = validator
		case wstakingtypes.ExperienceRegionId:
			s.experienceValidator = validator
		case "usa":
			s.usaValidator = validator
		}
	}
	s.Require().NotEmpty(s.meEarthValidator.OperatorAddress)
	s.Require().NotEmpty(s.experienceValidator.OperatorAddress)
	s.Require().NotEmpty(s.usaValidator.OperatorAddress)

	govParams := types.DefaultParams()
	govParams.MinVotingPeriod = 60
	s.Require().NoError(s.App.RegionGovKeeper.SetParams(s.Ctx, govParams))

	stakingMsg := stakingkeeper.NewMsgServerImpl(s.App.StakingKeeper.Keeper)
	regionMsg := wstakingkeeper.NewMsgServerImpl(s.App.StakingKeeper, s.App.TransferKeeper, stakingMsg)
	for _, region := range []struct {
		name     string
		operator string
	}{
		{name: wstakingtypes.MeEarthRegionName, operator: s.meEarthValidator.OperatorAddress},
		{name: wstakingtypes.ExperienceRegionName, operator: s.experienceValidator.OperatorAddress},
		{name: "USA", operator: s.usaValidator.OperatorAddress},
	} {
		_, err := regionMsg.NewRegion(s.Ctx, &wstakingtypes.MsgNewRegion{
			Creator:         s.Dao.GlobalDao,
			Name:            region.name,
			OperatorAddress: region.operator,
		})
		s.Require().NoError(err)
	}

	s.msgServer = regiongovkeeper.NewMsgServerImpl(s.App.RegionGovKeeper)
}

func operatorAccount(operator string) string {
	val, err := sdk.ValAddressFromBech32(operator)
	if err != nil {
		panic(err)
	}
	return sdk.AccAddress(val.Bytes()).String()
}

func (s *KeeperTestSuite) withdrawal(regionID string) types.Withdrawal {
	return types.Withdrawal{
		RegionId: regionID,
		Receiver: s.Dao.GlobalDao,
		Amount:   sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1)),
	}
}

func (s *KeeperTestSuite) submit(proposer string, period uint64, withdrawals ...types.Withdrawal) uint64 {
	res, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     proposer,
		Kind:         types.ProposalKindTreasuryWithdraw,
		Title:        "withdraw",
		VotingPeriod: period,
		Treasury:     &types.TreasuryWithdrawPayload{Withdrawals: withdrawals},
	})
	s.Require().NoError(err)
	return res.Id
}

func (s *KeeperTestSuite) vote(voter string, id uint64, option types.VoteOption) error {
	_, err := s.msgServer.Vote(s.Ctx, &types.MsgVote{
		Voter:      voter,
		ProposalId: id,
		Option:     option,
	})
	return err
}

func (s *KeeperTestSuite) electorAccounts(id uint64) []string {
	electors := s.App.RegionGovKeeper.GetElectors(s.Ctx, id)
	accounts := make([]string, 0, len(electors))
	for _, elector := range electors {
		accounts = append(accounts, operatorAccount(elector.Operator))
	}
	return accounts
}

func (s *KeeperTestSuite) endBlockAt(t time.Time) {
	s.Ctx = s.Ctx.WithBlockTime(t)
	regiongov.EndBlocker(s.Ctx, s.App.RegionGovKeeper)
}
