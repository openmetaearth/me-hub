package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

// EndBlocker tallies proposals whose voting period has ended.
// Passed means the count and the dao yes both succeeded. Execution is a later phase.
func (k Keeper) EndBlocker(ctx sdk.Context) {
	for _, id := range k.DequeueMatureProposals(ctx, ctx.BlockTime()) {
		proposal, found := k.GetProposal(ctx, id)
		if !found || proposal.Status != types.ProposalStatusVoting {
			continue
		}
		if k.tallyPassed(ctx, proposal) {
			proposal.Status = types.ProposalStatusPassed
		} else {
			proposal.Status = types.ProposalStatusRejected
		}
		k.SetProposal(ctx, proposal)
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeTallyProposal,
			sdk.NewAttribute(types.AttributeKeyProposalID, strconv.FormatUint(proposal.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyStatus, proposal.Status.String()),
		))
		k.Logger(ctx).Info("tallied regiongov proposal", "id", proposal.Id, "status", proposal.Status.String())
	}
}

// tallyPassed is true when the dao has voted yes and
// yes * 100 >= eligible * pass_threshold_percent. Absent and no votes are not yes.
// An empty electorate does not pass.
func (k Keeper) tallyPassed(ctx sdk.Context, proposal types.Proposal) bool {
	if proposal.DaoVote != types.DaoVoteYes {
		return false
	}
	electors := k.GetElectors(ctx, proposal.Id)
	eligible := uint64(len(electors))
	if eligible == 0 {
		return false
	}
	allowed := make(map[string]struct{}, len(electors))
	for _, elector := range electors {
		acc, ok := operatorAccount(elector.Operator)
		if !ok {
			continue
		}
		allowed[acc.String()] = struct{}{}
	}
	var yes uint64
	for _, vote := range k.GetVotes(ctx, proposal.Id) {
		if _, ok := allowed[vote.Voter]; !ok {
			continue
		}
		if vote.Option == types.VoteOptionYes {
			yes++
		}
	}
	threshold := uint64(k.GetParams(ctx).PassThresholdPercent)
	return yes*100 >= eligible*threshold
}
