package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (k Keeper) castVote(ctx sdk.Context, msg *types.MsgVote) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	voter, err := canonicalAccAddress(msg.Voter)
	if err != nil {
		return err
	}
	proposal, found := k.GetProposal(ctx, msg.ProposalId)
	if !found {
		return sdkerrors.Wrapf(types.ErrProposalNotFound, "%d", msg.ProposalId)
	}
	if proposal.Status != types.ProposalStatusVoting || !ctx.BlockTime().Before(proposal.VotingEndTime) {
		return sdkerrors.Wrapf(types.ErrVotingClosed, "proposal %d", proposal.Id)
	}

	if proposal.Kind != types.ProposalKindTreasuryWithdraw && proposal.Kind != types.ProposalKindReplaceRegionOperator {
		return sdkerrors.Wrap(types.ErrInvalidProposal, "unknown proposal kind")
	}
	if k.isDao(ctx, voter) {
		return k.castDaoVote(ctx, proposal, voter, msg.Option)
	}

	if proposal.Electorate == types.ElectorateKindRegionEquity {
		if err := k.authorizeEquityVote(ctx, proposal, voter); err != nil {
			return err
		}
	} else {
		elector, found := k.GetElector(ctx, proposal.Id, voter)
		if !found {
			return sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not an elector of proposal %d", voter, proposal.Id)
		}
		if k.isRegionFrozen(ctx, elector.RegionId) {
			return sdkerrors.Wrapf(types.ErrRegionFrozen, "region %s", elector.RegionId)
		}
		// Keep the snapshot for tally, but require current authority to cast a vote.
		region, found := k.lookupRegion(ctx, elector.RegionId)
		if !found {
			return sdkerrors.Wrapf(types.ErrUnknownRegion, "%s", elector.RegionId)
		}
		operator, ok := operatorAccount(region.OperatorAddress)
		if !ok || operator.String() != voter {
			return sdkerrors.Wrapf(types.ErrUnauthorized, "%s is no longer the operator of region %s", voter, region.RegionId)
		}
	}
	if _, found := k.GetVote(ctx, proposal.Id, voter); found {
		return sdkerrors.Wrapf(types.ErrAlreadyVoted, "%s", voter)
	}
	k.SetVote(ctx, types.Vote{
		ProposalId: proposal.Id,
		Voter:      voter,
		Option:     msg.Option,
	})
	k.emitVote(ctx, proposal.Id, voter, msg.Option)
	return nil
}

// castDaoVote records the dao ballot outside the counted electorate. Both yes
// and no can be overwritten until voting closes; EndBlocker uses the final ballot.
func (k Keeper) castDaoVote(ctx sdk.Context, proposal types.Proposal, voter string, option types.VoteOption) error {
	proposal.DaoVote = daoVote(option)
	k.SetProposal(ctx, proposal)
	k.emitVote(ctx, proposal.Id, voter, option)
	return nil
}

func (k Keeper) emitVote(ctx sdk.Context, proposalID uint64, voter string, option types.VoteOption) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeVote,
		sdk.NewAttribute(types.AttributeKeyProposalID, strconv.FormatUint(proposalID, 10)),
		sdk.NewAttribute(types.AttributeKeyVoter, voter),
		sdk.NewAttribute(types.AttributeKeyOption, option.String()),
	))
}

func daoVote(option types.VoteOption) types.DaoVote {
	switch option {
	case types.VoteOptionYes:
		return types.DaoVoteYes
	case types.VoteOptionNo:
		return types.DaoVoteNo
	default:
		return types.DaoVoteUnspecified
	}
}
