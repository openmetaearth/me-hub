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
		k.emitTally(ctx, proposal)
		return
	}

	cacheCtx, write := ctx.CacheContext()
	if err := k.executeProposal(cacheCtx, proposal); err != nil {
		proposal.Status = types.ProposalStatusFailed
		k.SetProposal(ctx, proposal)
		k.emitTally(ctx, proposal)
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeExecuteProposal,
			sdk.NewAttribute(types.AttributeKeyProposalID, strconv.FormatUint(proposal.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyStatus, proposal.Status.String()),
		))
		k.Logger(ctx).Error("regiongov proposal execution failed", "id", proposal.Id, "err", err)
		return
	}
	write()
	ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
	proposal.Status = types.ProposalStatusPassed
	k.SetProposal(ctx, proposal)
	k.emitTally(ctx, proposal)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeExecuteProposal,
		sdk.NewAttribute(types.AttributeKeyProposalID, strconv.FormatUint(proposal.Id, 10)),
		sdk.NewAttribute(types.AttributeKeyStatus, proposal.Status.String()),
	))
	k.Logger(ctx).Info("executed regiongov proposal", "id", proposal.Id, "status", proposal.Status.String())
}

func (k Keeper) emitTally(ctx sdk.Context, proposal types.Proposal) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeTallyProposal,
		sdk.NewAttribute(types.AttributeKeyProposalID, strconv.FormatUint(proposal.Id, 10)),
		sdk.NewAttribute(types.AttributeKeyStatus, proposal.Status.String()),
	))
	k.Logger(ctx).Info("tallied regiongov proposal", "id", proposal.Id, "status", proposal.Status.String())
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
	params := k.GetParams(ctx)
	for _, withdrawal := range proposal.Treasury.Withdrawals {
		maxPercent := uint32(0)
		if k.equityOnlyForRegion(ctx, proposal.Proposer, withdrawal.RegionId) {
			maxPercent = params.EquityWithdrawPercent
		}
		if err := k.stakingKeeper.WithdrawRegionTreasure(
			ctx,
			withdrawal.RegionId,
			withdrawal.Receiver,
			withdrawal.Amount,
			maxPercent,
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

// equityOnlyForRegion is true when the proposer is authorized for this region
// only as an equity address, not as dao or that region's operator.
func (k Keeper) equityOnlyForRegion(ctx sdk.Context, proposer, regionID string) bool {
	if k.isDao(ctx, proposer) {
		return false
	}
	region, found := k.lookupRegion(ctx, regionID)
	if !found {
		return false
	}
	if operator, ok := operatorAccount(region.OperatorAddress); ok && operator.String() == proposer {
		return false
	}
	equity, found := k.GetEquity(ctx, region.RegionId)
	if !found {
		return false
	}
	return containsAddress(equity.Addresses, proposer)
}
