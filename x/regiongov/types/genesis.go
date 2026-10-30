package types

import "fmt"

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("params: %w", err)
	}

	proposals := make(map[uint64]Proposal, len(gs.Proposals))
	var maxProposalID uint64
	for _, proposal := range gs.Proposals {
		if err := validateProposal(proposal, gs.Params); err != nil {
			return fmt.Errorf("proposal %d: %w", proposal.Id, err)
		}
		if _, ok := proposals[proposal.Id]; ok {
			return fmt.Errorf("duplicate proposal id %d", proposal.Id)
		}
		proposals[proposal.Id] = proposal
		if proposal.Id > maxProposalID {
			maxProposalID = proposal.Id
		}
	}
	if gs.ProposalSeq < maxProposalID {
		return fmt.Errorf("proposal_seq %d is below the highest proposal id %d", gs.ProposalSeq, maxProposalID)
	}

	electors := make(map[string]struct{}, len(gs.Electors))
	for _, elector := range gs.Electors {
		if _, ok := proposals[elector.ProposalId]; !ok {
			return fmt.Errorf("elector references missing proposal %d", elector.ProposalId)
		}
		if err := validateOperatorAddress(elector.Operator); err != nil {
			return fmt.Errorf("elector: %w", err)
		}
		if err := validateRegionID(elector.RegionId); err != nil {
			return fmt.Errorf("elector region: %w", err)
		}
		key := electorKey(elector.ProposalId, elector.Operator)
		if _, ok := electors[key]; ok {
			return fmt.Errorf("duplicate elector %s on proposal %d", elector.Operator, elector.ProposalId)
		}
		electors[key] = struct{}{}
	}

	votes := make(map[string]struct{}, len(gs.Votes))
	for _, vote := range gs.Votes {
		if _, ok := proposals[vote.ProposalId]; !ok {
			return fmt.Errorf("vote references missing proposal %d", vote.ProposalId)
		}
		if err := validateAccAddress(vote.Voter); err != nil {
			return fmt.Errorf("vote: %w", err)
		}
		if vote.Option != VoteOptionYes && vote.Option != VoteOptionNo {
			return fmt.Errorf("vote option must be yes or no")
		}
		key := electorKey(vote.ProposalId, vote.Voter)
		if _, ok := votes[key]; ok {
			return fmt.Errorf("duplicate vote from %s on proposal %d", vote.Voter, vote.ProposalId)
		}
		votes[key] = struct{}{}
	}

	regions := make(map[string]struct{}, len(gs.Equities))
	for _, equity := range gs.Equities {
		if err := validateRegionID(equity.RegionId); err != nil {
			return fmt.Errorf("equity: %w", err)
		}
		if _, ok := regions[equity.RegionId]; ok {
			return fmt.Errorf("duplicate equity list for region %s", equity.RegionId)
		}
		regions[equity.RegionId] = struct{}{}
		if err := validateUniqueAddresses(equity.Addresses); err != nil {
			return fmt.Errorf("equity %s: %w", equity.RegionId, err)
		}
	}

	freezes := make(map[string]struct{}, len(gs.RegionFreezes))
	for _, freeze := range gs.RegionFreezes {
		if err := validateRegionID(freeze.RegionId); err != nil {
			return fmt.Errorf("region freeze: %w", err)
		}
		if _, ok := freezes[freeze.RegionId]; ok {
			return fmt.Errorf("duplicate region freeze %s", freeze.RegionId)
		}
		freezes[freeze.RegionId] = struct{}{}
	}

	equityFreezes := make(map[string]struct{}, len(gs.EquityFreezes))
	for _, freeze := range gs.EquityFreezes {
		if err := validateRegionID(freeze.RegionId); err != nil {
			return fmt.Errorf("equity freeze: %w", err)
		}
		if err := validateAccAddress(freeze.Address); err != nil {
			return fmt.Errorf("equity freeze: %w", err)
		}
		key := freeze.RegionId + "/" + freeze.Address
		if _, ok := equityFreezes[key]; ok {
			return fmt.Errorf("duplicate equity freeze %s in %s", freeze.Address, freeze.RegionId)
		}
		equityFreezes[key] = struct{}{}
	}

	actions := make(map[uint64]struct{}, len(gs.PendingActions))
	var maxActionID uint64
	for _, action := range gs.PendingActions {
		if err := validatePendingAction(action); err != nil {
			return fmt.Errorf("pending action %d: %w", action.Id, err)
		}
		if _, ok := actions[action.Id]; ok {
			return fmt.Errorf("duplicate pending action id %d", action.Id)
		}
		actions[action.Id] = struct{}{}
		if action.Id > maxActionID {
			maxActionID = action.Id
		}
	}
	if gs.PendingActionSeq < maxActionID {
		return fmt.Errorf("pending_action_seq %d is below the highest pending action id %d", gs.PendingActionSeq, maxActionID)
	}
	return nil
}

func validateProposal(proposal Proposal, params Params) error {
	if proposal.Id == 0 {
		return fmt.Errorf("id is empty")
	}
	if err := validateAccAddress(proposal.Proposer); err != nil {
		return err
	}
	if proposal.Title == "" {
		return fmt.Errorf("title is empty")
	}
	if proposal.Status == ProposalStatusUnspecified {
		return fmt.Errorf("status is empty")
	}
	if proposal.VotingPeriod < params.MinVotingPeriod || proposal.VotingPeriod > MaxVotingPeriod {
		return fmt.Errorf("voting_period %d is outside [%d,%d]", proposal.VotingPeriod, params.MinVotingPeriod, MaxVotingPeriod)
	}
	if proposal.SubmitTime.IsZero() || !proposal.VotingEndTime.After(proposal.SubmitTime) {
		return fmt.Errorf("voting end time must be after submit time")
	}
	if proposal.Electorate == ElectorateKindUnspecified {
		return fmt.Errorf("electorate is empty")
	}
	if proposal.DaoVote != DaoVoteUnspecified && proposal.DaoVote != DaoVoteYes && proposal.DaoVote != DaoVoteNo {
		return fmt.Errorf("unknown dao vote")
	}

	switch proposal.Kind {
	case ProposalKindTreasuryWithdraw:
		if proposal.ReplaceOperator != nil {
			return fmt.Errorf("treasury proposal cannot include a replace-operator payload")
		}
		if proposal.Treasury == nil {
			return fmt.Errorf("treasury payload is empty")
		}
		if err := validateWithdrawals(proposal.Treasury.Withdrawals); err != nil {
			return err
		}
	case ProposalKindReplaceRegionOperator:
		if proposal.Treasury != nil {
			return fmt.Errorf("replace-operator proposal cannot include a treasury payload")
		}
		if err := validateReplaceOperator(proposal.ReplaceOperator); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown proposal kind")
	}

	if proposal.Electorate == ElectorateKindRegionEquity {
		if err := validateOperatorAddress(proposal.RequiredOperator); err != nil {
			return fmt.Errorf("required operator: %w", err)
		}
		return nil
	}
	if proposal.RequiredOperator != "" {
		return fmt.Errorf("required operator is only set for an equity electorate")
	}
	return nil
}

func validatePendingAction(action PendingAction) error {
	if action.Id == 0 {
		return fmt.Errorf("id is empty")
	}
	if action.Kind == PendingActionKindUnspecified {
		return fmt.Errorf("kind is empty")
	}
	if action.Status == PendingActionStatusUnspecified {
		return fmt.Errorf("status is empty")
	}
	if err := validateRegionID(action.RegionId); err != nil {
		return err
	}
	if err := validateAccAddress(action.Proposer); err != nil {
		return err
	}
	if action.SubmitTime.IsZero() {
		return fmt.Errorf("submit time is empty")
	}
	switch action.Kind {
	case PendingActionKindSetEquity:
		return validateUniqueAddresses(action.Addresses)
	case PendingActionKindFreezeEquity, PendingActionKindUnfreezeEquity:
		return validateEquityFreezeAddresses(action.Addresses)
	default:
		return fmt.Errorf("unknown pending action kind")
	}
}

func electorKey(proposalID uint64, address string) string {
	return fmt.Sprintf("%d/%s", proposalID, address)
}
