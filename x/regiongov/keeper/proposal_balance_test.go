package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestTreasuryWithdrawBalanceAtSubmission() {
	for _, role := range []string{"dao", "operator"} {
		for _, tc := range []struct {
			name             string
			amount           sdk.Coins
			emptyTreasury    bool
			withGlobalRegion bool
			wantErr          bool
		}{
			{
				name:   "below balance and above equity cap",
				amount: sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 999_999), sdk.NewInt64Coin("utest", 499)),
			},
			{
				name:   "exact balance of each denom",
				amount: sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1_000_000), sdk.NewInt64Coin("utest", 500)),
			},
			{
				name:    "one above balance",
				amount:  sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1_000_001)),
				wantErr: true,
			},
			{
				name:    "second denom above balance",
				amount:  sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1), sdk.NewInt64Coin("utest", 501)),
				wantErr: true,
			},
			{
				name:    "missing denom has zero balance",
				amount:  sdk.NewCoins(sdk.NewInt64Coin("uunknown", 1)),
				wantErr: true,
			},
			{
				name:          "empty treasury",
				amount:        sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1)),
				emptyTreasury: true,
				wantErr:       true,
			},
			{
				name:             "one insufficient region rejects entire proposal",
				amount:           sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1_000_001)),
				withGlobalRegion: true,
				wantErr:          true,
			},
			{
				name:             "multiple regions with sufficient balances",
				amount:           sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1_000_000)),
				withGlobalRegion: true,
			},
		} {
			s.Run(role+"/"+tc.name, func() {
				s.SetupTest()
				proposer := s.Dao.GlobalDao
				if role == "operator" {
					proposer = operatorAccount(s.usaValidator.OperatorAddress)
				}
				receiver := sdk.MustAccAddressFromBech32(sample.AccAddress())
				region, found := s.App.StakingKeeper.GetRegion(s.Ctx, "usa")
				s.Require().True(found)
				treasury := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
				additional := sdk.NewCoins(sdk.NewInt64Coin("utest", 500))
				s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, additional))
				s.Require().NoError(s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, treasury, additional))
				if tc.emptyTreasury {
					balance := s.App.BankKeeper.GetAllBalances(s.Ctx, treasury)
					s.Require().NoError(s.App.BankKeeper.SendCoins(s.Ctx, treasury, sdk.MustAccAddressFromBech32(s.Dao.GlobalDao), balance))
				}

				withdrawals := []types.Withdrawal{{RegionId: "usa", Receiver: receiver.String(), Amount: tc.amount}}
				if tc.withGlobalRegion {
					// The earlier valid withdrawal must not mask a later insufficient balance.
					withdrawals = append([]types.Withdrawal{s.withdrawal(wstakingtypes.MeEarthRegionId)}, withdrawals...)
				}
				balances := make(map[string]sdk.Coins)
				for _, withdrawal := range withdrawals {
					region, found := s.App.StakingKeeper.GetRegion(s.Ctx, withdrawal.RegionId)
					s.Require().True(found)
					balances[region.RegionTreasureAddr] = s.App.BankKeeper.GetAllBalances(s.Ctx, sdk.MustAccAddressFromBech32(region.RegionTreasureAddr))
				}
				seq := s.App.RegionGovKeeper.GetProposalSeq(s.Ctx)
				events := len(s.Ctx.EventManager().Events())
				res, err := s.msgServer.SubmitProposal(s.Ctx, &types.MsgSubmitProposal{
					Proposer:     proposer,
					Kind:         types.ProposalKindTreasuryWithdraw,
					Title:        "check treasury balance at submission",
					VotingPeriod: 60,
					Treasury:     &types.TreasuryWithdrawPayload{Withdrawals: withdrawals},
				})
				if tc.wantErr {
					s.Require().ErrorIs(err, sdkerrors.ErrInsufficientFunds)
					s.Require().Contains(err.Error(), "region usa spendable treasury balance")
					s.Require().Nil(res)
					s.Require().Equal(seq, s.App.RegionGovKeeper.GetProposalSeq(s.Ctx))
					s.Require().Empty(s.App.RegionGovKeeper.GetAllProposals(s.Ctx))
					s.Require().Empty(s.App.RegionGovKeeper.GetElectors(s.Ctx, seq+1))
					s.Require().Empty(s.App.RegionGovKeeper.DequeueMatureProposals(s.Ctx, s.startTime.Add(time.Hour)))
					s.Require().Len(s.Ctx.EventManager().Events(), events)
				} else {
					s.Require().NoError(err)
					s.Require().Equal(seq+1, res.Id)
					s.Require().Equal(types.ProposalStatusVoting, s.mustProposal(res.Id).Status)
				}
				// Submission neither transfers nor reserves treasury funds.
				for address, balance := range balances {
					s.Require().Equal(balance, s.App.BankKeeper.GetAllBalances(s.Ctx, sdk.MustAccAddressFromBech32(address)))
				}
				s.Require().Empty(s.App.BankKeeper.GetAllBalances(s.Ctx, receiver))
			})
		}
	}
}
