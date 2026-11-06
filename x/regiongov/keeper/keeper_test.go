package keeper_test

import (
	"strings"
	"testing"
	"time"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov"
	regiongovkeeper "github.com/openmetaearth/me-hub/x/regiongov/keeper"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
	wstakingkeeper "github.com/openmetaearth/me-hub/x/wstaking/keeper"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	msgServer           types.MsgServer
	stakingMsgServer    wstakingkeeper.MsgServer
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
	s.stakingMsgServer = wstakingkeeper.NewMsgServerImpl(s.App.StakingKeeper, s.App.TransferKeeper, stakingMsg)
	for _, region := range []struct {
		name     string
		operator string
	}{
		{name: wstakingtypes.MeEarthRegionName, operator: s.meEarthValidator.OperatorAddress},
		{name: wstakingtypes.ExperienceRegionName, operator: s.experienceValidator.OperatorAddress},
		{name: "USA", operator: s.usaValidator.OperatorAddress},
	} {
		s.createFundedRegion(region.name, region.operator)
	}

	s.msgServer = regiongovkeeper.NewMsgServerImpl(s.App.RegionGovKeeper)
}

func (s *KeeperTestSuite) createFundedRegion(name, operator string) {
	s.T().Helper()
	_, err := s.stakingMsgServer.NewRegion(s.Ctx, &wstakingtypes.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            name,
		OperatorAddress: operator,
	})
	s.Require().NoError(err)
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, strings.ToLower(name))
	s.Require().True(found)
	treasury := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
	fund := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1_000_000))
	s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, mintypes.ModuleName, fund))
	s.Require().NoError(s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, mintypes.ModuleName, treasury, fund))
}

// addCountryRegion creates a real validator and its bound, funded region. Tests
// with multiple electors must not bypass wstaking's validator-region binding.
func (s *KeeperTestSuite) addCountryRegion(regionID string) string {
	s.T().Helper()
	operator := sdk.MustAccAddressFromBech32(sample.AccAddress())
	pubkey, err := codectypes.NewAnyWithValue(ed25519.GenPrivKey().PubKey())
	s.Require().NoError(err)
	amount := sdk.NewInt64Coin(params.BaseDenom, 1_000_000)
	s.FundModuleAcc(wstakingtypes.StakePoolName, sdk.NewCoins(amount))
	_, err = s.stakingMsgServer.CreateValidator(s.Ctx, &stakingtypes.MsgCreateValidator{
		Description: stakingtypes.Description{Moniker: regionID, RegionID: regionID},
		Commission: stakingtypes.NewCommissionRates(
			s.App.StakingKeeper.MinCommissionRate(s.Ctx), sdk.OneDec(), sdk.ZeroDec(),
		),
		MinSelfDelegation: sdk.OneInt(),
		DelegatorAddress:  s.Dao.GlobalDao,
		ValidatorAddress:  sdk.ValAddress(operator).String(),
		Pubkey:            pubkey,
		Value:             amount,
	})
	s.Require().NoError(err)
	s.createFundedRegion(strings.ToUpper(regionID), sdk.ValAddress(operator).String())
	return operator.String()
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

func (s *KeeperTestSuite) submit(proposer string, withdrawals ...types.Withdrawal) uint64 {
	res, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
		Proposer:     proposer,
		Kind:         types.ProposalKindTreasuryWithdraw,
		Title:        "withdraw",
		VotingPeriod: 60,
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
