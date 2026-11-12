package v2_0_18_rc1_test //nolint:revive

import (
	"testing"
	"time"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/upgrades/v2_0_18_rc1"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

func TestUpgradePreservesRegionGovState(t *testing.T) {
	application := apptesting.Setup(t, false)
	now := time.Unix(1_700_000_000, 0).UTC()
	ctx := application.BaseApp.NewContext(false, cometbftproto.Header{Height: 10, Time: now})
	k := application.RegionGovKeeper
	account := sample.AccAddress()
	params := types.DefaultParams()
	params.MinVotingPeriod = 121
	require.NoError(t, k.SetParams(ctx, params))
	proposal := types.Proposal{
		Id: 7, Kind: types.ProposalKindUpdateMinVotingPeriod, Proposer: account,
		Title: "existing proposal", Status: types.ProposalStatusVoting,
		SubmitTime: now, VotingEndTime: now.Add(time.Hour), VotingPeriod: 3600,
		Electorate: types.ElectorateKindCountryRegionOperators, DaoVote: types.DaoVoteYes,
		UpdateMinVotingPeriod: &types.UpdateMinVotingPeriodPayload{MinVotingPeriod: 180},
	}
	k.SetProposal(ctx, proposal)
	k.SetProposalSeq(ctx, proposal.Id)
	k.InsertVotingQueue(ctx, proposal.VotingEndTime, proposal.Id)
	// Seed the old valoper representation directly, as persisted by v2.0.18.
	elector := types.Elector{
		ProposalId: proposal.Id, RegionId: "chn",
		Operator: sdk.ValAddress(sdk.MustAccAddressFromBech32(account)).String(),
	}
	store := ctx.KVStore(application.GetKey(types.StoreKey))
	store.Set(types.ElectorKey(proposal.Id, account), application.AppCodec().MustMarshal(&elector))
	k.SetVote(ctx, types.Vote{ProposalId: proposal.Id, Voter: account, Option: types.VoteOptionYes})
	k.SetEquity(ctx, types.RegionEquity{RegionId: "chn", Addresses: []string{account}})
	k.SetRegionFreeze(ctx, types.RegionFreeze{RegionId: "usa", Frozen: true})
	k.SetEquityFreeze(ctx, types.EquityFreeze{RegionId: "chn", Address: account, Frozen: true})
	k.SetPendingAction(ctx, types.PendingAction{
		Id: 3, Kind: types.PendingActionKindSetEquity, Status: types.PendingActionStatusPending,
		RegionId: "chn", Proposer: account, Addresses: []string{account}, SubmitTime: now,
	})
	k.SetPendingActionSeq(ctx, 3)
	// A legacy equity proposal keeps its original approver even if the live
	// region now points at a different validator operator.
	operator := sdk.MustAccAddressFromBech32(sample.AccAddress())
	equityProposal := types.Proposal{
		Id: 8, Kind: types.ProposalKindTreasuryWithdraw, Proposer: account,
		Title: "legacy equity proposal", Status: types.ProposalStatusVoting,
		SubmitTime: now, VotingEndTime: now.Add(2 * time.Hour), VotingPeriod: 7200,
		Electorate: types.ElectorateKindRegionEquity, DaoVote: types.DaoVoteYes,
		RequiredOperator: sdk.ValAddress(operator).String(),
		Treasury: &types.TreasuryWithdrawPayload{Withdrawals: []types.Withdrawal{{
			RegionId: "chn", Receiver: account, Amount: sdk.NewCoins(sdk.NewInt64Coin("umec", 1)),
		}}},
	}
	store.Set(types.ProposalKey(equityProposal.Id), application.AppCodec().MustMarshal(&equityProposal))
	k.SetProposalSeq(ctx, equityProposal.Id)
	k.InsertVotingQueue(ctx, equityProposal.VotingEndTime, equityProposal.Id)
	k.SetElector(ctx, types.Elector{ProposalId: equityProposal.Id, Operator: account, RegionId: "chn"})
	k.SetVote(ctx, types.Vote{ProposalId: equityProposal.Id, Voter: operator.String(), Option: types.VoteOptionYes})
	application.StakingKeeper.SetRegion(ctx, wstakingtypes.Region{
		RegionId: "chn", OperatorAddress: sdk.ValAddress(sdk.MustAccAddressFromBech32(sample.AccAddress())).String(),
	})
	regionsBefore := application.StakingKeeper.GetAllRegion(ctx)
	before := regiongov.ExportGenesis(ctx, k)
	require.NoError(t, before.Validate())
	versionMap := application.UpgradeKeeper.GetModuleVersionMap(ctx)
	require.Contains(t, versionMap, types.ModuleName, "upgrade starts from an initialized regiongov module")

	require.Equal(t, "v2.0.18.rc1", v2_0_18_rc1.Upgrade.Name)
	require.Empty(t, v2_0_18_rc1.Upgrade.StoreUpgrades.Added)
	require.Empty(t, v2_0_18_rc1.Upgrade.StoreUpgrades.Deleted)
	require.Empty(t, v2_0_18_rc1.Upgrade.StoreUpgrades.Renamed)
	require.True(t, application.UpgradeKeeper.HasHandler(v2_0_18_rc1.UpgradeName))
	require.NotPanics(t, func() {
		application.UpgradeKeeper.ApplyUpgrade(ctx, upgradetypes.Plan{Name: v2_0_18_rc1.UpgradeName, Height: ctx.BlockHeight()})
	})
	require.Equal(t, versionMap, application.UpgradeKeeper.GetModuleVersionMap(ctx))
	for i := range before.Proposals {
		if before.Proposals[i].Id == equityProposal.Id {
			before.Proposals[i].RequiredOperator = operator.String()
		}
	}
	require.Equal(t, before, regiongov.ExportGenesis(ctx, k), "only required_operator encoding changes")
	require.Equal(t, regionsBefore, application.StakingKeeper.GetAllRegion(ctx), "staking regions are untouched")
	require.NoError(t, k.MigrateRequiredOperatorAccounts(ctx))
	require.Equal(t, before, regiongov.ExportGenesis(ctx, k), "migration is idempotent")
	storedProposal, found := k.GetProposal(ctx, equityProposal.Id)
	require.True(t, found)
	require.Equal(t, operator.String(), storedProposal.RequiredOperator)
	require.NoError(t, regiongov.ExportGenesis(ctx, k).Validate())
	response, err := k.Electors(sdk.WrapSDKContext(ctx), &types.QueryElectorsRequest{ProposalId: proposal.Id})
	require.NoError(t, err)
	require.Len(t, response.Electors, 1)
	require.Equal(t, account, response.Electors[0].Operator)
	require.Equal(t, application.AppCodec().MustMarshal(&elector), store.Get(types.ElectorKey(proposal.Id, account)))
	require.Equal(t, []uint64{proposal.Id}, k.DequeueMatureProposals(ctx, proposal.VotingEndTime))
}
