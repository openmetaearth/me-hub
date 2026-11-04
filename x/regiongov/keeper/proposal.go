package keeper

import (
	"errors"
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
	if msg.Kind == types.ProposalKindTreasuryWithdraw {
		if err := k.enforceTreasuryWithdrawLimits(ctx, proposer, msg.Treasury.Withdrawals); err != nil {
			return 0, err
		}
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
	dao, operatorRegionIDs, operatorErr := k.proposerAuthority(ctx, proposer)
	equityRegionIDs, equityErr := k.activeEquityRegions(ctx, proposer)
	if operatorErr != nil && len(equityRegionIDs) == 0 {
		if errors.Is(operatorErr, types.ErrUnauthorized) && equityErr != nil {
			return equityErr
		}
		return operatorErr
	}
	if operatorErr != nil {
		dao = false
		operatorRegionIDs = nil
	}
	if msg.Kind != types.ProposalKindTreasuryWithdraw {
		if dao || len(operatorRegionIDs) > 0 {
			return nil
		}
		return sdkerrors.Wrap(types.ErrUnauthorized, "an equity address cannot replace a region operator")
	}
	if dao {
		return nil
	}
	for _, withdrawal := range msg.Treasury.Withdrawals {
		if err := k.authorizeWithdrawal(ctx, withdrawal.RegionId, operatorRegionIDs, equityRegionIDs); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) authorizeWithdrawal(ctx sdk.Context, regionID string, operatorRegionIDs, equityRegionIDs []string) error {
	region, found := k.lookupRegion(ctx, regionID)
	if !found {
		return sdkerrors.Wrapf(types.ErrUnknownRegion, "%s", regionID)
	}
	// An operator may include the global or experience treasury, or a country they run.
	if len(operatorRegionIDs) > 0 && (isSpecialRegion(region.RegionId) || containsRegion(operatorRegionIDs, region.RegionId)) {
		return nil
	}
	if containsRegion(equityRegionIDs, region.RegionId) {
		return nil
	}
	if len(operatorRegionIDs) > 0 {
		return sdkerrors.Wrapf(
			types.ErrUnauthorized,
			"a region operator can withdraw only from regions they operate, got %s",
			regionID,
		)
	}
	return sdkerrors.Wrapf(
		types.ErrUnauthorized,
		"an equity address can withdraw only from its own region, got %s",
		regionID,
	)
}

// enforceTreasuryWithdrawLimits checks each region and denom before creating a
// proposal. All withdrawals are limited to the spendable treasury balance;
// equity-only proposers additionally have the equity_withdraw_percent limit.
// These checks do not reserve funds, so execution can still fail if the balance falls.
func (k Keeper) enforceTreasuryWithdrawLimits(ctx sdk.Context, proposer string, withdrawals []types.Withdrawal) error {
	percent := k.GetParams(ctx).EquityWithdrawPercent
	for _, withdrawal := range withdrawals {
		spendable, err := k.stakingKeeper.SpendableRegionTreasure(ctx, withdrawal.RegionId)
		if err != nil {
			return err
		}
		equityOnly := k.equityOnlyForRegion(ctx, proposer, withdrawal.RegionId)
		for _, coin := range withdrawal.Amount {
			balance := spendable.AmountOf(coin.Denom)
			if equityOnly && coin.Amount.MulRaw(100).GT(balance.MulRaw(int64(percent))) {
				return sdkerrors.Wrapf(
					types.ErrEquityWithdrawLimit,
					"equity withdrawal %s exceeds %d%% of region %s treasury balance %s",
					coin, percent, withdrawal.RegionId, sdk.NewCoin(coin.Denom, balance),
				)
			}
			if coin.Amount.GT(balance) {
				return sdkerrors.Wrapf(
					sdkerrors.ErrInsufficientFunds,
					"withdrawal %s exceeds region %s spendable treasury balance %s",
					coin, withdrawal.RegionId, sdk.NewCoin(coin.Denom, balance),
				)
			}
		}
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

// activeEquityRegions returns regions whose equity list contains proposer and
// where that address can still propose. A frozen region or a frozen equity
// address is omitted. When every membership is blocked, blocked explains why.
func (k Keeper) activeEquityRegions(ctx sdk.Context, proposer string) ([]string, error) {
	equities := k.GetAllEquities(ctx)
	regionIDs := make([]string, 0, len(equities))
	regionFrozen := false
	equityFrozen := false
	for _, equity := range equities {
		if !containsAddress(equity.Addresses, proposer) {
			continue
		}
		if k.isRegionFrozen(ctx, equity.RegionId) {
			regionFrozen = true
			continue
		}
		if k.isEquityFrozen(ctx, equity.RegionId, proposer) {
			equityFrozen = true
			continue
		}
		regionIDs = append(regionIDs, equity.RegionId)
	}
	if len(regionIDs) > 0 {
		return regionIDs, nil
	}
	if regionFrozen {
		return nil, sdkerrors.Wrap(types.ErrRegionFrozen, "equity region is frozen")
	}
	if equityFrozen {
		return nil, sdkerrors.Wrap(types.ErrEquityFrozen, "equity address is frozen")
	}
	return nil, nil
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
	regions := k.stakingKeeper.GetAllRegion(ctx)
	operatorRegionIDs := make([]string, 0, len(regions))
	frozenOperator := false
	for _, region := range regions {
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
	regions := k.stakingKeeper.GetAllRegion(ctx)
	seen := make(map[string]struct{})
	electors := make([]types.Elector, 0, len(regions))
	for _, region := range regions {
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
