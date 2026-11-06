package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

// EndBlocker finalizes proposals whose voting period has ended and prunes old
// completed pending actions. Passed proposals execute; otherwise they are Rejected.
func (k Keeper) EndBlocker(ctx sdk.Context) {
	for _, id := range k.DequeueMatureProposals(ctx, ctx.BlockTime()) {
		proposal, found := k.GetProposal(ctx, id)
		if !found || proposal.Status != types.ProposalStatusVoting {
			continue
		}
		k.finalizeProposal(ctx, proposal, k.tallyPassed(ctx, proposal))
	}
	k.pruneCompletedPendingActions(ctx)
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
