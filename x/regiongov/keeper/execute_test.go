package keeper_test

import (
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (s *KeeperTestSuite) TestFinalizeProposalEvent() {
	for _, tc := range []struct {
		name          string
		operatorYes   bool
		daoYes        bool
		emptyTreasury bool
		tallyResult   string
		status        types.ProposalStatus
	}{
		{
			name:        "passed and executed",
			operatorYes: true,
			daoYes:      true,
			tallyResult: "passed",
			status:      types.ProposalStatusPassed,
		},
		{
			name:        "rejected without operator approval",
			daoYes:      true,
			tallyResult: "rejected",
			status:      types.ProposalStatusRejected,
		},
		{
			name:        "rejected without dao approval",
			operatorYes: true,
			tallyResult: "rejected",
			status:      types.ProposalStatusRejected,
		},
		{
			name:          "passed but execution failed with empty treasury",
			operatorYes:   true,
			daoYes:        true,
			emptyTreasury: true,
			tallyResult:   "passed",
			status:        types.ProposalStatusFailed,
		},
	} {
		s.Run(tc.name, func() {
			s.SetupTest()
			operator := operatorAccount(s.usaValidator.OperatorAddress)
			receiver := sdk.MustAccAddressFromBech32(sample.AccAddress())
			region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
			s.Require().True(found)
			treasury := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
			withdrawal := s.withdrawal("usa")
			withdrawal.Receiver = receiver.String()
			id := s.submit(operator, withdrawal)
			// Submission requires sufficient funds, but the balance can fall before execution.
			if tc.emptyTreasury {
				balance := s.App.BankKeeper.GetAllBalances(s.Ctx, treasury)
				s.Require().NoError(s.App.BankKeeper.SendCoins(s.Ctx, treasury, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), balance))
			}
			before := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)
			receiverBefore := s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom)

			if tc.operatorYes {
				s.Require().NoError(s.vote(operator, id, types.VoteOptionYes))
			}
			if tc.daoYes {
				s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
			}
			end := s.mustProposal(id).VotingEndTime
			s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
			s.endBlockAt(end)
			s.Require().Equal(tc.status, s.mustProposal(id).Status)

			var tallyEvents sdk.Events
			for _, event := range s.Ctx.EventManager().Events() {
				s.Require().NotEqual("execute_proposal", event.Type)
				if event.Type == types.EventTypeTallyProposal {
					tallyEvents = append(tallyEvents, event)
				}
				if tc.status != types.ProposalStatusPassed {
					s.Require().NotEqual(banktypes.EventTypeTransfer, event.Type)
				}
			}
			s.Require().Len(tallyEvents, 1)
			attributes := make(map[string]string)
			for _, attribute := range tallyEvents[0].Attributes {
				attributes[attribute.Key] = attribute.Value
			}
			expected := map[string]string{
				types.AttributeKeyProposalID:    strconv.FormatUint(id, 10),
				types.AttributeKeyTallyResult:   tc.tallyResult,
				types.AttributeKeyExecuteStatus: tc.status.String(),
			}
			if tc.emptyTreasury {
				s.Require().Contains(attributes[types.AttributeKeyError], "region usa")
				s.Require().Contains(attributes[types.AttributeKeyError], "does not have enough balance")
				expected[types.AttributeKeyError] = attributes[types.AttributeKeyError]
			}
			s.Require().Equal(expected, attributes)
			s.Require().Len(tallyEvents[0].Attributes, len(expected))

			after := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)
			receiverAfter := s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom)
			if tc.status == types.ProposalStatusPassed {
				s.Require().Equal(before.Amount.Sub(withdrawal.Amount[0].Amount), after.Amount)
				s.Require().Equal(receiverBefore.Amount.Add(withdrawal.Amount[0].Amount), receiverAfter.Amount)
			} else {
				s.Require().Equal(before, after)
				s.Require().Equal(receiverBefore, receiverAfter)
			}

			// Finalization and its event are not repeated in later blocks.
			s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())
			s.endBlockAt(end.Add(time.Second))
			s.Require().Empty(s.Ctx.EventManager().Events())
			s.Require().Equal(tc.status, s.mustProposal(id).Status)
		})
	}
}

func (s *KeeperTestSuite) TestEquityWithdrawCap() {
	operator := operatorAccount(s.usaValidator.OperatorAddress)
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
	s.Require().True(found)
	treasury := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
	balance := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)
	percent := int64(s.App.RegionGovKeeper.GetParams(s.Ctx).EquityWithdrawPercent)
	s.Require().Equal(int64(types.DefaultEquityWithdrawPercent), percent)
	maxAllowed := balance.Amount.MulRaw(percent).QuoRaw(100)
	s.Require().True(maxAllowed.IsPositive())

	equity := sample.AccAddress()
	s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity}})

	s.Run("equity above one percent rejects the submit transaction", func() {
		receiver := sdk.MustAccAddressFromBech32(sample.AccAddress())
		before := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)
		existing := len(s.App.RegionGovKeeper.GetAllProposals(s.Ctx))
		withdrawal := s.withdrawal("usa")
		withdrawal.Receiver = receiver.String()
		withdrawal.Amount = sdk.NewCoins(sdk.NewCoin(params.BaseDenom, maxAllowed.AddRaw(1)))

		_, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
			Proposer:     equity,
			Kind:         types.ProposalKindTreasuryWithdraw,
			Title:        "withdraw",
			VotingPeriod: 60,
			Treasury:     &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{withdrawal}},
		})
		s.Require().ErrorIs(err, types.ErrEquityWithdrawLimit)
		s.Require().Len(s.App.RegionGovKeeper.GetAllProposals(s.Ctx), existing)
		s.Require().Equal(before, s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom))
		s.Require().True(s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom).IsZero())
	})

	s.Run("equity at one percent is paid", func() {
		receiver := sdk.MustAccAddressFromBech32(sample.AccAddress())
		before := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)
		withdrawal := s.withdrawal("usa")
		withdrawal.Receiver = receiver.String()
		withdrawal.Amount = sdk.NewCoins(sdk.NewCoin(params.BaseDenom, maxAllowed))

		id := s.passTreasury(equity, withdrawal)
		s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)
		s.Require().Equal(before.Amount.Sub(maxAllowed), s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom).Amount)
		s.Require().Equal(maxAllowed, s.App.BankKeeper.GetBalance(s.Ctx, receiver, params.BaseDenom).Amount)
	})

	s.Run("region operator is not subject to the equity cap", func() {
		receiver := sdk.MustAccAddressFromBech32(sample.AccAddress())
		before := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)
		amount := maxAllowed.MulRaw(10)
		s.Require().True(before.Amount.GT(amount))
		withdrawal := s.withdrawal("usa")
		withdrawal.Receiver = receiver.String()
		withdrawal.Amount = sdk.NewCoins(sdk.NewCoin(params.BaseDenom, amount))

		id := s.passTreasury(operator, withdrawal)
		s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)
		s.Require().Equal(before.Amount.Sub(amount), s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom).Amount)
	})

	s.Run("dao is not subject to the equity cap", func() {
		receiver := sdk.MustAccAddressFromBech32(sample.AccAddress())
		before := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)
		amount := maxAllowed.MulRaw(10)
		s.Require().True(before.Amount.GT(amount))
		withdrawal := s.withdrawal("usa")
		withdrawal.Receiver = receiver.String()
		withdrawal.Amount = sdk.NewCoins(sdk.NewCoin(params.BaseDenom, amount))

		id := s.passTreasury(s.Dao.GlobalDao, withdrawal)
		s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)
		s.Require().Equal(before.Amount.Sub(amount), s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom).Amount)
	})

	s.Run("region operator who is also an equity address is not subject to the equity cap", func() {
		s.App.RegionGovKeeper.SetEquity(s.Ctx, types.RegionEquity{RegionId: "usa", Addresses: []string{equity, operator}})
		receiver := sdk.MustAccAddressFromBech32(sample.AccAddress())
		before := s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom)
		amount := maxAllowed.MulRaw(10)
		s.Require().True(before.Amount.GT(amount))
		withdrawal := s.withdrawal("usa")
		withdrawal.Receiver = receiver.String()
		withdrawal.Amount = sdk.NewCoins(sdk.NewCoin(params.BaseDenom, amount))

		id := s.passTreasury(operator, withdrawal)
		s.Require().Equal(types.ProposalStatusPassed, s.mustProposal(id).Status)
		s.Require().Equal(before.Amount.Sub(amount), s.App.BankKeeper.GetBalance(s.Ctx, treasury, params.BaseDenom).Amount)
	})
}

func (s *KeeperTestSuite) passTreasury(proposer string, withdrawal types.Withdrawal) uint64 {
	id := s.submit(proposer, withdrawal)
	if s.mustProposal(id).Electorate == types.ElectorateKindRegionEquity {
		for _, voter := range s.electorAccounts(id) {
			s.Require().NoError(s.vote(voter, id, types.VoteOptionYes))
		}
	}
	s.Require().NoError(s.vote(operatorAccount(s.usaValidator.OperatorAddress), id, types.VoteOptionYes))
	s.Require().NoError(s.vote(s.Dao.GlobalDao, id, types.VoteOptionYes))
	s.endBlockAt(s.mustProposal(id).VotingEndTime)
	return id
}
