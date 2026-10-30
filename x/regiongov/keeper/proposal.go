package keeper

import (
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/openmetaearth/me-hub/x/regiongov/types"
)

func (k Keeper) submitProposal(ctx sdk.Context, msg *types.MsgSubmitProposal) (uint64, error) {
	if err := msg.ValidateBasic(); err != nil {
		return 0, err
	}
	params := k.GetParams(ctx)
	if msg.VotingPeriod < params.MinVotingPeriod {
		return 0, sdkerrors.Wrapf(
			types.ErrInvalidVotingPeriod,
			"voting period %d is below min_voting_period %d",
			msg.VotingPeriod, params.MinVotingPeriod,
		)
	}

	proposer, err := canonicalAccAddress(msg.Proposer)
	if err != nil {
		return 0, err
	}
	electorate, err := k.resolveElectorate(ctx, msg)
	if err != nil {
		return 0, err
	}
	if err := k.authorizeProposer(ctx, proposer, msg); err != nil {
		return 0, err
	}

	id := k.GetProposalSeq(ctx) + 1
	electors, err := k.snapshotElectors(ctx, id, electorate)
	if err != nil {
		return 0, err
	}

	now := ctx.BlockTime()
	proposal := types.Proposal{
		Id:              id,
		Kind:            msg.Kind,
		Proposer:        proposer,
		Title:           msg.Title,
		Summary:         msg.Summary,
		Status:          types.ProposalStatusVoting,
		SubmitTime:      now,
		VotingEndTime:   now.Add(time.Duration(msg.VotingPeriod) * time.Second),
		VotingPeriod:    msg.VotingPeriod,
		Electorate:      electorate,
		Treasury:        msg.Treasury,
		ReplaceOperator: msg.ReplaceOperator,
	}
	k.SetProposal(ctx, proposal)
	for _, elector := range electors {
		k.SetElector(ctx, elector)
	}
	k.SetProposalSeq(ctx, id)
	k.InsertVotingQueue(ctx, proposal.VotingEndTime, id)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSubmitProposal,
		sdk.NewAttribute(types.AttributeKeyProposalID, strconv.FormatUint(id, 10)),
		sdk.NewAttribute(types.AttributeKeyProposer, proposer),
		sdk.NewAttribute(types.AttributeKeyElectorate, electorate.String()),
	))
	return id, nil
}

func (k Keeper) resolveElectorate(ctx sdk.Context, msg *types.MsgSubmitProposal) (types.ElectorateKind, error) {
	switch msg.Kind {
	case types.ProposalKindTreasuryWithdraw:
		return k.resolveTreasuryElectorate(ctx, msg)
	case types.ProposalKindReplaceRegionOperator:
		if _, found := k.lookupRegion(ctx, msg.ReplaceOperator.RegionId); !found {
			return types.ElectorateKindUnspecified, sdkerrors.Wrapf(types.ErrUnknownRegion, "%s", msg.ReplaceOperator.RegionId)
		}
		return types.ElectorateKindCountryRegionOperators, nil
	default:
		return types.ElectorateKindUnspecified, sdkerrors.Wrap(types.ErrInvalidProposal, "unknown proposal kind")
	}
}

func (k Keeper) resolveTreasuryElectorate(ctx sdk.Context, msg *types.MsgSubmitProposal) (types.ElectorateKind, error) {
	for _, withdrawal := range msg.Treasury.Withdrawals {
		if _, found := k.lookupRegion(ctx, withdrawal.RegionId); !found {
			return types.ElectorateKindUnspecified, sdkerrors.Wrapf(types.ErrUnknownRegion, "%s", withdrawal.RegionId)
		}
	}
	return types.ElectorateKindCountryRegionOperators, nil
}

func (k Keeper) authorizeProposer(ctx sdk.Context, proposer string, msg *types.MsgSubmitProposal) error {
	dao, operatorRegionIDs, err := k.proposerAuthority(ctx, proposer)
	if err != nil {
		return err
	}
	if msg.Kind != types.ProposalKindTreasuryWithdraw || dao {
		return nil
	}
	for _, withdrawal := range msg.Treasury.Withdrawals {
		region, found := k.lookupRegion(ctx, withdrawal.RegionId)
		if !found {
			return sdkerrors.Wrapf(types.ErrUnknownRegion, "%s", withdrawal.RegionId)
		}
		// Any region operator may include the global or experience treasury.
		// A country withdrawal is limited to a region that operator runs.
		if isSpecialRegion(region.RegionId) {
			continue
		}
		if !containsRegion(operatorRegionIDs, region.RegionId) {
			return sdkerrors.Wrapf(
				types.ErrUnauthorized,
				"a region operator can withdraw only from regions they operate, got %s",
				withdrawal.RegionId,
			)
		}
	}
	return nil
}

// proposerAuthority accepts the dao or an operator of at least one unfrozen region.
// A frozen region's operator cannot propose unless they are the dao or still
// operate another unfrozen region.
func (k Keeper) proposerAuthority(ctx sdk.Context, proposer string) (bool, []string, error) {
	acc, err := sdk.AccAddressFromBech32(proposer)
	if err != nil {
		return false, nil, err
	}
	dao := k.isDao(ctx, proposer)
	var operatorRegionIDs []string
	frozenOperator := false
	for _, region := range k.stakingKeeper.GetAllRegion(ctx) {
		operator, ok := operatorAccount(region.OperatorAddress)
		if !ok || !operator.Equals(acc) {
			continue
		}
		if k.isRegionFrozen(ctx, region.RegionId) {
			frozenOperator = true
			continue
		}
		operatorRegionIDs = append(operatorRegionIDs, region.RegionId)
	}
	if dao || len(operatorRegionIDs) > 0 {
		return dao, operatorRegionIDs, nil
	}
	if frozenOperator {
		return false, nil, sdkerrors.Wrap(types.ErrRegionFrozen, "region operator is frozen")
	}
	return false, nil, sdkerrors.Wrap(types.ErrUnauthorized, "proposer is not the dao or a region operator")
}

// snapshotElectors records one vote per country-region operator account.
// The global region and the experience region do not vote. Frozen regions and regions
// without an operator are left out. A proposing country region is included.
func (k Keeper) snapshotElectors(ctx sdk.Context, proposalID uint64, kind types.ElectorateKind) ([]types.Elector, error) {
	if kind != types.ElectorateKindCountryRegionOperators {
		return nil, sdkerrors.Wrap(types.ErrInvalidProposal, "electorate is not a counted operator vote")
	}
	seen := make(map[string]struct{})
	var electors []types.Elector
	for _, region := range k.stakingKeeper.GetAllRegion(ctx) {
		if isSpecialRegion(region.RegionId) {
			continue
		}
		if k.isRegionFrozen(ctx, region.RegionId) {
			continue
		}
		operator, ok := operatorAccount(region.OperatorAddress)
		if !ok {
			continue
		}
		if _, exists := seen[operator.String()]; exists {
			continue
		}
		seen[operator.String()] = struct{}{}
		electors = append(electors, types.Elector{
			ProposalId: proposalID,
			Operator:   region.OperatorAddress,
			RegionId:   region.RegionId,
		})
	}
	if len(electors) == 0 {
		return nil, types.ErrEmptyElectorate
	}
	return electors, nil
}
