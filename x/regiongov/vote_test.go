package regiongov_test

import (
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/openmetaearth/me-hub/x/regiongov/keeper"
	"github.com/openmetaearth/me-hub/x/regiongov/types"
	wstakingtypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

type stubStaking struct {
	regions []wstakingtypes.Region
}

func (s stubStaking) GetAllRegion(sdk.Context) []wstakingtypes.Region { return s.regions }

func (s stubStaking) GetRegion(_ sdk.Context, regionID string) (wstakingtypes.Region, bool) {
	for _, region := range s.regions {
		if region.RegionId == regionID || strings.EqualFold(region.RegionId, regionID) {
			return region, true
		}
	}
	return wstakingtypes.Region{}, false
}

type stubDao struct {
	dao string
}

func (d stubDao) IsGlobalDao(_ sdk.Context, address string) bool {
	return d.dao != "" && d.dao == address
}

func regionOf(id string, operatorAcc string) wstakingtypes.Region {
	return wstakingtypes.Region{
		RegionId:        id,
		OperatorAddress: sdk.ValAddress(sdk.MustAccAddressFromBech32(operatorAcc)).String(),
	}
}

func setupGov(t *testing.T, regions []wstakingtypes.Region, dao string) (keeper.Keeper, sdk.Context, types.MsgServer) {
	t.Helper()
	k, ctx := newKeeperWith(t, stubStaking{regions: regions}, stubDao{dao: dao})
	params := types.DefaultParams()
	params.MinVotingPeriod = 60
	require.NoError(t, k.SetParams(ctx, params))
	ctx = ctx.WithBlockTime(time.Unix(1_700_000_000, 0).UTC())
	return k, ctx, keeper.NewMsgServerImpl(k)
}

func withdrawal(regionID, receiver string) types.Withdrawal {
	return types.Withdrawal{
		RegionId: regionID,
		Receiver: receiver,
		Amount:   sdk.NewCoins(sdk.NewInt64Coin("umec", 1)),
	}
}

func submitTreasury(srv types.MsgServer, ctx sdk.Context, proposer string, period uint64, withdrawals ...types.Withdrawal) (*types.MsgSubmitProposalResponse, error) {
	return srv.SubmitProposal(sdk.WrapSDKContext(ctx), &types.MsgSubmitProposal{
		Proposer:     proposer,
		Kind:         types.ProposalKindTreasuryWithdraw,
		Title:        "withdraw",
		VotingPeriod: period,
		Treasury:     &types.TreasuryWithdrawPayload{Withdrawals: withdrawals},
	})
}

func cast(srv types.MsgServer, ctx sdk.Context, voter string, id uint64, option types.VoteOption) error {
	_, err := srv.Vote(sdk.WrapSDKContext(ctx), &types.MsgVote{
		Voter:      voter,
		ProposalId: id,
		Option:     option,
	})
	return err
}

func electorAccounts(t *testing.T, k keeper.Keeper, ctx sdk.Context, id uint64) []string {
	t.Helper()
	electors := k.GetElectors(ctx, id)
	accounts := make([]string, 0, len(electors))
	for _, elector := range electors {
		acc, err := sdk.ValAddressFromBech32(elector.Operator)
		require.NoError(t, err)
		accounts = append(accounts, sdk.AccAddress(acc.Bytes()).String())
	}
	return accounts
}

func TestTreasuryElectorates(t *testing.T) {
	cn := sample.AccAddress()
	us := sample.AccAddress()
	earth := sample.AccAddress()
	experience := sample.AccAddress()
	super := sample.AccAddress()
	receiver := sample.AccAddress()
	regions := []wstakingtypes.Region{
		regionOf(wstakingtypes.MeEarthRegionId, earth),
		regionOf(wstakingtypes.ExperienceRegionId, experience),
		regionOf("cn", cn),
		regionOf("us", us),
	}
	k, ctx, srv := setupGov(t, regions, super)

	global, err := submitTreasury(srv, ctx, cn, 60, withdrawal(wstakingtypes.MeEarthRegionId, receiver))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{cn, us}, electorAccounts(t, k, ctx, global.Id))
	stored, found := k.GetProposal(ctx, global.Id)
	require.True(t, found)
	require.Equal(t, types.ElectorateKindCountryRegionOperators, stored.Electorate)
	require.Equal(t, types.DaoVoteUnspecified, stored.DaoVote)
	require.Equal(t, ctx.BlockTime().Add(60*time.Second), stored.VotingEndTime)

	country, err := submitTreasury(srv, ctx, cn, 60, withdrawal("cn", receiver))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{cn, us}, electorAccounts(t, k, ctx, country.Id))

	_, err = submitTreasury(srv, ctx, cn, 60, withdrawal("us", receiver))
	require.ErrorIs(t, err, types.ErrUnauthorized)

	_, err = submitTreasury(srv, ctx, cn, 60, withdrawal("cn", receiver), withdrawal("us", receiver))
	require.ErrorIs(t, err, types.ErrUnauthorized)

	multi, err := submitTreasury(srv, ctx, super, 60, withdrawal("cn", receiver), withdrawal("us", receiver))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{cn, us}, electorAccounts(t, k, ctx, multi.Id))
	require.NotContains(t, electorAccounts(t, k, ctx, multi.Id), super)

	mixed, err := submitTreasury(srv, ctx, super, 60, withdrawal(wstakingtypes.MeEarthRegionId, receiver), withdrawal("cn", receiver))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{cn, us}, electorAccounts(t, k, ctx, mixed.Id))

	_, err = submitTreasury(srv, ctx, cn, 60, withdrawal(wstakingtypes.MeEarthRegionId, receiver), withdrawal("us", receiver))
	require.ErrorIs(t, err, types.ErrUnauthorized)

	_, err = submitTreasury(srv, ctx, super, 60, withdrawal("missing", receiver))
	require.ErrorIs(t, err, types.ErrUnknownRegion)

	_, err = submitTreasury(srv, ctx, super, 1, withdrawal("cn", receiver))
	require.ErrorIs(t, err, types.ErrInvalidVotingPeriod)

	_, err = submitTreasury(srv, ctx, sample.AccAddress(), 60, withdrawal("cn", receiver))
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestFrozenRegionIsOutOfTheSnapshot(t *testing.T) {
	cn := sample.AccAddress()
	us := sample.AccAddress()
	earth := sample.AccAddress()
	super := sample.AccAddress()
	receiver := sample.AccAddress()
	k, ctx, srv := setupGov(t, []wstakingtypes.Region{
		regionOf(wstakingtypes.MeEarthRegionId, earth),
		regionOf("cn", cn),
		regionOf("us", us),
	}, super)
	k.SetRegionFreeze(ctx, types.RegionFreeze{RegionId: "cn", Frozen: true})

	_, err := submitTreasury(srv, ctx, cn, 60, withdrawal(wstakingtypes.MeEarthRegionId, receiver))
	require.ErrorIs(t, err, types.ErrRegionFrozen)

	res, err := submitTreasury(srv, ctx, us, 60, withdrawal(wstakingtypes.MeEarthRegionId, receiver))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{us}, electorAccounts(t, k, ctx, res.Id))

	replaced, err := srv.SubmitProposal(sdk.WrapSDKContext(ctx), &types.MsgSubmitProposal{
		Proposer:     super,
		Kind:         types.ProposalKindReplaceRegionOperator,
		Title:        "replace",
		VotingPeriod: 60,
		ReplaceOperator: &types.ReplaceRegionOperatorPayload{
			RegionId:    "cn",
			NewOperator: sdk.ValAddress(sdk.MustAccAddressFromBech32(sample.AccAddress())).String(),
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, []string{us}, electorAccounts(t, k, ctx, replaced.Id))
}

func TestEmptyElectorateDoesNotOpen(t *testing.T) {
	super := sample.AccAddress()
	k, ctx, srv := setupGov(t, []wstakingtypes.Region{{
		RegionId: "cn",
	}}, super)
	_, err := submitTreasury(srv, ctx, super, 60, withdrawal("cn", sample.AccAddress()))
	require.ErrorIs(t, err, types.ErrEmptyElectorate)
	require.Empty(t, k.GetAllProposals(ctx))
}

func TestDuplicateOperatorHasOneVote(t *testing.T) {
	operatorAcc := sample.AccAddress()
	dao := sample.AccAddress()
	k, ctx, srv := setupGov(t, []wstakingtypes.Region{
		regionOf("cn", operatorAcc),
		regionOf("us", operatorAcc),
	}, dao)
	res, err := submitTreasury(srv, ctx, dao, 60, withdrawal("cn", sample.AccAddress()), withdrawal("us", sample.AccAddress()))
	require.NoError(t, err)
	require.Equal(t, []string{operatorAcc}, electorAccounts(t, k, ctx, res.Id))
}

func TestSuperVetoAndOverwrite(t *testing.T) {
	cn := sample.AccAddress()
	us := sample.AccAddress()
	super := sample.AccAddress()
	k, ctx, srv := setupGov(t, []wstakingtypes.Region{
		regionOf("cn", cn),
		regionOf("us", us),
	}, super)
	res, err := submitTreasury(srv, ctx, cn, 60, withdrawal("cn", sample.AccAddress()))
	require.NoError(t, err)

	require.NoError(t, cast(srv, ctx, super, res.Id, types.VoteOptionYes))
	proposal, _ := k.GetProposal(ctx, res.Id)
	require.Equal(t, types.DaoVoteYes, proposal.DaoVote)
	require.Equal(t, types.ProposalStatusVoting, proposal.Status)

	require.NoError(t, cast(srv, ctx, super, res.Id, types.VoteOptionNo))
	proposal, _ = k.GetProposal(ctx, res.Id)
	require.Equal(t, types.DaoVoteNo, proposal.DaoVote)
	require.Equal(t, types.ProposalStatusRejected, proposal.Status)

	require.ErrorIs(t, cast(srv, ctx, super, res.Id, types.VoteOptionYes), types.ErrVotingClosed)
	ctx = ctx.WithBlockTime(proposal.VotingEndTime)
	k.EndBlocker(ctx)
	proposal, _ = k.GetProposal(ctx, res.Id)
	require.Equal(t, types.ProposalStatusRejected, proposal.Status)
}

func TestBallotRulesAndTally(t *testing.T) {
	cn := sample.AccAddress()
	us := sample.AccAddress()
	jp := sample.AccAddress()
	super := sample.AccAddress()
	outsider := sample.AccAddress()
	k, ctx, srv := setupGov(t, []wstakingtypes.Region{
		regionOf("cn", cn),
		regionOf("us", us),
		regionOf("jp", jp),
	}, super)

	half, err := submitTreasury(srv, ctx, cn, 60, withdrawal("cn", sample.AccAddress()))
	require.NoError(t, err)
	require.ErrorIs(t, cast(srv, ctx, outsider, half.Id, types.VoteOptionYes), types.ErrUnauthorized)
	require.NoError(t, cast(srv, ctx, us, half.Id, types.VoteOptionYes))
	require.ErrorIs(t, cast(srv, ctx, us, half.Id, types.VoteOptionNo), types.ErrAlreadyVoted)
	require.NoError(t, cast(srv, ctx, cn, half.Id, types.VoteOptionNo))
	require.NoError(t, cast(srv, ctx, super, half.Id, types.VoteOptionYes))

	twoThirds, err := submitTreasury(srv, ctx, cn, 60, withdrawal("cn", sample.AccAddress()))
	require.NoError(t, err)
	require.NoError(t, cast(srv, ctx, us, twoThirds.Id, types.VoteOptionYes))
	require.NoError(t, cast(srv, ctx, jp, twoThirds.Id, types.VoteOptionYes))
	require.NoError(t, cast(srv, ctx, super, twoThirds.Id, types.VoteOptionYes))

	noSuper, err := submitTreasury(srv, ctx, cn, 60, withdrawal("cn", sample.AccAddress()))
	require.NoError(t, err)
	require.NoError(t, cast(srv, ctx, cn, noSuper.Id, types.VoteOptionYes))
	require.NoError(t, cast(srv, ctx, us, noSuper.Id, types.VoteOptionYes))
	require.NoError(t, cast(srv, ctx, jp, noSuper.Id, types.VoteOptionYes))

	open, _ := k.GetProposal(ctx, half.Id)
	k.EndBlocker(ctx)
	stillOpen, _ := k.GetProposal(ctx, half.Id)
	require.Equal(t, types.ProposalStatusVoting, stillOpen.Status)

	ctx = ctx.WithBlockTime(open.VotingEndTime.Add(-time.Second))
	require.NoError(t, cast(srv, ctx, jp, half.Id, types.VoteOptionNo))
	ctx = ctx.WithBlockTime(open.VotingEndTime)
	require.ErrorIs(t, cast(srv, ctx, sample.AccAddress(), half.Id, types.VoteOptionYes), types.ErrVotingClosed)

	k.EndBlocker(ctx)
	rejected, _ := k.GetProposal(ctx, half.Id)
	passed, _ := k.GetProposal(ctx, twoThirds.Id)
	missedSuper, _ := k.GetProposal(ctx, noSuper.Id)
	require.Equal(t, types.ProposalStatusRejected, rejected.Status)
	require.Equal(t, types.ProposalStatusPassed, passed.Status)
	require.Equal(t, types.ProposalStatusRejected, missedSuper.Status)

	k.EndBlocker(ctx)
	passedAgain, _ := k.GetProposal(ctx, twoThirds.Id)
	require.Equal(t, types.ProposalStatusPassed, passedAgain.Status)
}

func TestFreezeAfterSnapshotBlocksThatVote(t *testing.T) {
	cn := sample.AccAddress()
	us := sample.AccAddress()
	super := sample.AccAddress()
	k, ctx, srv := setupGov(t, []wstakingtypes.Region{
		regionOf("cn", cn),
		regionOf("us", us),
	}, super)
	res, err := submitTreasury(srv, ctx, super, 60, withdrawal("cn", sample.AccAddress()), withdrawal("us", sample.AccAddress()))
	require.NoError(t, err)
	k.SetRegionFreeze(ctx, types.RegionFreeze{RegionId: "cn", Frozen: true})

	require.ErrorIs(t, cast(srv, ctx, cn, res.Id, types.VoteOptionYes), types.ErrRegionFrozen)
	require.NoError(t, cast(srv, ctx, us, res.Id, types.VoteOptionYes))
	require.NoError(t, cast(srv, ctx, super, res.Id, types.VoteOptionYes))

	proposal, _ := k.GetProposal(ctx, res.Id)
	ctx = ctx.WithBlockTime(proposal.VotingEndTime)
	k.EndBlocker(ctx)
	proposal, _ = k.GetProposal(ctx, res.Id)
	require.Equal(t, types.ProposalStatusRejected, proposal.Status)
}

func TestOneOfTwoFailsAndTwoOfTwoPasses(t *testing.T) {
	cn := sample.AccAddress()
	us := sample.AccAddress()
	super := sample.AccAddress()
	k, ctx, srv := setupGov(t, []wstakingtypes.Region{
		regionOf("cn", cn),
		regionOf("us", us),
	}, super)

	failing, err := submitTreasury(srv, ctx, cn, 60, withdrawal("cn", sample.AccAddress()))
	require.NoError(t, err)
	require.NoError(t, cast(srv, ctx, us, failing.Id, types.VoteOptionYes))
	require.NoError(t, cast(srv, ctx, super, failing.Id, types.VoteOptionYes))

	passing, err := submitTreasury(srv, ctx, cn, 60, withdrawal("cn", sample.AccAddress()))
	require.NoError(t, err)
	require.NoError(t, cast(srv, ctx, cn, passing.Id, types.VoteOptionYes))
	require.NoError(t, cast(srv, ctx, us, passing.Id, types.VoteOptionYes))
	require.NoError(t, cast(srv, ctx, super, passing.Id, types.VoteOptionYes))

	proposal, _ := k.GetProposal(ctx, failing.Id)
	ctx = ctx.WithBlockTime(proposal.VotingEndTime)
	k.EndBlocker(ctx)

	failed, _ := k.GetProposal(ctx, failing.Id)
	passed, _ := k.GetProposal(ctx, passing.Id)
	require.Equal(t, types.ProposalStatusRejected, failed.Status)
	require.Equal(t, types.ProposalStatusPassed, passed.Status)
}
