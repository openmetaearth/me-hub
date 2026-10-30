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

	if k.isDao(ctx, voter) {
		return k.castDaoVote(ctx, proposal, voter, msg.Option)
	}

	elector, found := k.GetElector(ctx, proposal.Id, voter)
	if !found {
		return sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not an elector of proposal %d", voter, proposal.Id)
	}
	if k.isRegionFrozen(ctx, elector.RegionId) {
		return sdkerrors.Wrapf(types.ErrRegionFrozen, "region %s", elector.RegionId)
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

// castDaoVote records the dao ballot. It is not a 51% vote.
// No rejects the proposal immediately. Yes can be overwritten with No until then.
func (k Keeper) castDaoVote(ctx sdk.Context, proposal types.Proposal, voter string, option types.VoteOption) error {
	proposal.DaoVote = daoVote(option)
	if option == types.VoteOptionNo {
		proposal.Status = types.ProposalStatusRejected
		k.DeleteVotingQueue(ctx, proposal.VotingEndTime, proposal.Id)
		k.SetProposal(ctx, proposal)
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeVetoProposal,
			sdk.NewAttribute(types.AttributeKeyProposalID, strconv.FormatUint(proposal.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyVoter, voter),
			sdk.NewAttribute(types.AttributeKeyStatus, proposal.Status.String()),
		))
		return nil
	}
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
