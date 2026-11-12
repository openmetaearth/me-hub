package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

// snapshotEquityElectors fixes the counted accounts and the separately required
// operator at submission. Later equity/operator replacements do not rewrite them.
func (k Keeper) snapshotEquityElectors(ctx sdk.Context, proposalID uint64, regionID string) ([]types.Elector, string, error) {
	region, found := k.lookupRegion(ctx, regionID)
	if !found {
		return nil, "", sdkerrors.Wrapf(types.ErrUnknownRegion, "%s", regionID)
	}
	if k.isRegionFrozen(ctx, region.RegionId) {
		return nil, "", sdkerrors.Wrapf(types.ErrRegionFrozen, "%s", region.RegionId)
	}
	operator, ok := operatorAccount(region.OperatorAddress)
	if !ok {
		return nil, "", sdkerrors.Wrap(types.ErrInvalidProposal, "an equity withdrawal requires a region operator")
	}
	equity, found := k.GetEquity(ctx, region.RegionId)
	if !found {
		return nil, "", types.ErrEmptyElectorate
	}
	seen := make(map[string]struct{}, len(equity.Addresses))
	electors := make([]types.Elector, 0, len(equity.Addresses))
	for _, address := range equity.Addresses {
		account, err := canonicalAccAddress(address)
		if err != nil {
			return nil, "", err
		}
		// These approvals are independent of the counted equity ballots.
		if account == operator.String() || k.isDao(ctx, account) || k.isEquityFrozen(ctx, region.RegionId, account) {
			continue
		}
		if _, exists := seen[account]; exists {
			continue
		}
		seen[account] = struct{}{}
		electors = append(electors, types.Elector{
			ProposalId: proposalID,
			Operator:   account,
			RegionId:   region.RegionId,
		})
	}
	if len(electors) == 0 {
		return nil, "", types.ErrEmptyElectorate
	}
	return electors, operator.String(), nil
}

// equityProposalRegion also guards imported proposals before indexing the payload.
func equityProposalRegion(proposal types.Proposal) (string, bool) {
	if proposal.Kind != types.ProposalKindTreasuryWithdraw || proposal.Treasury == nil || len(proposal.Treasury.Withdrawals) != 1 {
		return "", false
	}
	return proposal.Treasury.Withdrawals[0].RegionId, true
}

func (k Keeper) authorizeEquityVote(ctx sdk.Context, proposal types.Proposal, voter string) error {
	regionID, ok := equityProposalRegion(proposal)
	if !ok {
		return sdkerrors.Wrap(types.ErrInvalidProposal, "an equity proposal must withdraw from exactly one region")
	}
	region, found := k.lookupRegion(ctx, regionID)
	if !found {
		return sdkerrors.Wrapf(types.ErrUnknownRegion, "%s", regionID)
	}
	if k.isRegionFrozen(ctx, region.RegionId) {
		return sdkerrors.Wrapf(types.ErrRegionFrozen, "%s", region.RegionId)
	}
	operator, err := proposal.RequiredOperatorAccount()
	if err == nil && operator.String() == voter {
		return nil
	}
	elector, found := k.GetElector(ctx, proposal.Id, voter)
	if !found || !containsRegion([]string{region.RegionId}, elector.RegionId) {
		return sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not an elector of proposal %d", voter, proposal.Id)
	}
	if k.isEquityFrozen(ctx, region.RegionId, voter) {
		return sdkerrors.Wrapf(types.ErrEquityFrozen, "%s", voter)
	}
	return nil
}

func (k Keeper) equityOperatorApproved(ctx sdk.Context, proposal types.Proposal) bool {
	if _, ok := equityProposalRegion(proposal); !ok {
		return false
	}
	operator, err := proposal.RequiredOperatorAccount()
	if err != nil {
		return false
	}
	vote, found := k.GetVote(ctx, proposal.Id, operator.String())
	return found && vote.Option == types.VoteOptionYes
}
