package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

// finalizeProposal is called when VotingEndTime is reached. execute means the
// tally and dao yes both passed; otherwise the proposal is Rejected.
func (k Keeper) finalizeProposal(ctx sdk.Context, proposal types.Proposal, execute bool) {
	if proposal.Status != types.ProposalStatusVoting {
		return
	}

	k.DeleteVotingQueue(ctx, proposal.VotingEndTime, proposal.Id)

	if !execute {
		proposal.Status = types.ProposalStatusRejected
		k.SetProposal(ctx, proposal)
		k.emitTally(ctx, proposal, false, nil)
		return
	}

	cacheCtx, write := ctx.CacheContext()
	if err := k.executeProposal(cacheCtx, proposal); err != nil {
		proposal.Status = types.ProposalStatusFailed
		k.SetProposal(ctx, proposal)
		k.emitTally(ctx, proposal, true, err)
		return
	}
	write()
	ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
	proposal.Status = types.ProposalStatusPassed
	k.SetProposal(ctx, proposal)
	k.emitTally(ctx, proposal, true, nil)
}

// emitTally reports the vote outcome and final proposal status in one event.
// A rejected tally skips execution; a passed tally can still fail execution.
func (k Keeper) emitTally(ctx sdk.Context, proposal types.Proposal, passed bool, executionErr error) {
	result := "rejected"
	if passed {
		result = "passed"
	}
	attributes := []sdk.Attribute{
		sdk.NewAttribute(types.AttributeKeyProposalID, strconv.FormatUint(proposal.Id, 10)),
		sdk.NewAttribute(types.AttributeKeyTallyResult, result),
		sdk.NewAttribute(types.AttributeKeyExecuteStatus, proposal.Status.String()),
	}
	if executionErr != nil {
		attributes = append(attributes, sdk.NewAttribute(types.AttributeKeyError, executionErr.Error()))
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeTallyProposal, attributes...))

	logger := k.Logger(ctx).With("id", proposal.Id, "tally_result", result, "execute_status", proposal.Status.String())
	if executionErr != nil {
		logger.Error("regiongov proposal execution failed", "err", executionErr)
		return
	}
	logger.Info("finalized regiongov proposal")
}

func (k Keeper) executeProposal(ctx sdk.Context, proposal types.Proposal) error {
	switch proposal.Kind {
	case types.ProposalKindTreasuryWithdraw:
		return k.executeTreasury(ctx, proposal)
	case types.ProposalKindReplaceRegionOperator:
		return k.executeReplaceOperator(ctx, proposal)
	default:
		return sdkerrors.Wrap(types.ErrInvalidProposal, "unknown proposal kind")
	}
}

func (k Keeper) executeTreasury(ctx sdk.Context, proposal types.Proposal) error {
	if proposal.Treasury == nil || len(proposal.Treasury.Withdrawals) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidProposal, "treasury payload is empty")
	}
	for _, withdrawal := range proposal.Treasury.Withdrawals {
		if err := k.stakingKeeper.WithdrawRegionTreasure(
			ctx,
			withdrawal.RegionId,
			withdrawal.Receiver,
			withdrawal.Amount,
		); err != nil {
			return sdkerrors.Wrapf(types.ErrProposalExecution, "region %s: %s", withdrawal.RegionId, err)
		}
	}
	return nil
}

func (k Keeper) executeReplaceOperator(ctx sdk.Context, proposal types.Proposal) error {
	if proposal.ReplaceOperator == nil {
		return sdkerrors.Wrap(types.ErrInvalidProposal, "replace operator payload is empty")
	}
	if err := k.stakingKeeper.SetRegionOperator(
		ctx,
		proposal.ReplaceOperator.RegionId,
		proposal.ReplaceOperator.NewOperator,
	); err != nil {
		return sdkerrors.Wrapf(types.ErrProposalExecution, "%s", err)
	}
	return nil
}
