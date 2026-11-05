package types

import (
	"errors"
	"fmt"
)

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

	proposals, err := validateGenesisProposals(gs)
	if err != nil {
		return err
	}
	if err := validateGenesisElectors(gs, proposals); err != nil {
		return err
	}
	if err := validateGenesisVotes(gs, proposals); err != nil {
		return err
	}
	if err := validateGenesisEquities(gs); err != nil {
		return err
	}
	if err := validateGenesisRegionFreezes(gs); err != nil {
		return err
	}
	if err := validateGenesisEquityFreezes(gs); err != nil {
		return err
	}
	return validateGenesisPendingActions(gs)
}

func validateGenesisProposals(gs GenesisState) (map[uint64]Proposal, error) {
	proposals := make(map[uint64]Proposal, len(gs.Proposals))
	var maxProposalID uint64
	for _, proposal := range gs.Proposals {
		if err := validateProposal(proposal); err != nil {
			return nil, fmt.Errorf("proposal %d: %w", proposal.Id, err)
		}
		if _, ok := proposals[proposal.Id]; ok {
			return nil, fmt.Errorf("duplicate proposal id %d", proposal.Id)
		}
		proposals[proposal.Id] = proposal
		if proposal.Id > maxProposalID {
			maxProposalID = proposal.Id
		}
	}
	if gs.ProposalSeq < maxProposalID {
		return nil, fmt.Errorf("proposal_seq %d is below the highest proposal id %d", gs.ProposalSeq, maxProposalID)
	}
	return proposals, nil
}

func validateGenesisElectors(gs GenesisState, proposals map[uint64]Proposal) error {
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
		key := genesisRecordKey(elector.ProposalId, elector.Operator)
		if _, ok := electors[key]; ok {
			return fmt.Errorf("duplicate elector %s on proposal %d", elector.Operator, elector.ProposalId)
		}
		electors[key] = struct{}{}
	}
	return nil
}

func validateGenesisVotes(gs GenesisState, proposals map[uint64]Proposal) error {
	votes := make(map[string]struct{}, len(gs.Votes))
	for _, vote := range gs.Votes {
		if _, ok := proposals[vote.ProposalId]; !ok {
			return fmt.Errorf("vote references missing proposal %d", vote.ProposalId)
		}
		if err := validateAccAddress(vote.Voter); err != nil {
			return fmt.Errorf("vote: %w", err)
		}
		if vote.Option != VoteOptionYes && vote.Option != VoteOptionNo {
			return errors.New("vote option must be yes or no")
		}
		key := genesisRecordKey(vote.ProposalId, vote.Voter)
		if _, ok := votes[key]; ok {
			return fmt.Errorf("duplicate vote from %s on proposal %d", vote.Voter, vote.ProposalId)
		}
		votes[key] = struct{}{}
	}
	return nil
}

func validateGenesisEquities(gs GenesisState) error {
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
	return nil
}

func validateGenesisRegionFreezes(gs GenesisState) error {
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
	return nil
}

func validateGenesisEquityFreezes(gs GenesisState) error {
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
	return nil
}

func validateGenesisPendingActions(gs GenesisState) error {
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

func validateProposal(proposal Proposal) error {
	if err := validateProposalHeader(proposal); err != nil {
		return err
	}
	if err := validateProposalPayload(proposal); err != nil {
		return err
	}
	return validateProposalOperator(proposal)
}

func validateProposalHeader(proposal Proposal) error {
	if proposal.Id == 0 {
		return errors.New("id is empty")
	}
	if err := validateAccAddress(proposal.Proposer); err != nil {
		return err
	}
	if proposal.Title == "" {
		return errors.New("title is empty")
	}
	if proposal.Status == ProposalStatusUnspecified {
		return errors.New("status is empty")
	}
	// Stored proposals were validated against the minimum at submission time.
	// A later parameter increase must not invalidate an exported genesis.
	if proposal.VotingPeriod < 1 || proposal.VotingPeriod > MaxVotingPeriod {
		return fmt.Errorf("voting_period %d is outside [1,%d]", proposal.VotingPeriod, MaxVotingPeriod)
	}
	if proposal.SubmitTime.IsZero() || !proposal.VotingEndTime.After(proposal.SubmitTime) {
		return errors.New("voting end time must be after submit time")
	}
	if proposal.Electorate == ElectorateKindUnspecified {
		return errors.New("electorate is empty")
	}
	if proposal.DaoVote != DaoVoteUnspecified && proposal.DaoVote != DaoVoteYes && proposal.DaoVote != DaoVoteNo {
		return errors.New("unknown dao vote")
	}
	return nil
}

func validateProposalOperator(proposal Proposal) error {
	if proposal.Electorate == ElectorateKindRegionEquity {
		if err := validateOperatorAddress(proposal.RequiredOperator); err != nil {
			return fmt.Errorf("required operator: %w", err)
		}
		return nil
	}
	if proposal.RequiredOperator != "" {
		return errors.New("required operator is only set for an equity electorate")
	}
	return nil
}

func validatePendingAction(action PendingAction) error {
	if action.Id == 0 {
		return errors.New("id is empty")
	}
	if action.Kind == PendingActionKindUnspecified {
		return errors.New("kind is empty")
	}
	if action.Status == PendingActionStatusUnspecified {
		return errors.New("status is empty")
	}
	if err := validateRegionID(action.RegionId); err != nil {
		return err
	}
	if err := validateAccAddress(action.Proposer); err != nil {
		return err
	}
	if action.SubmitTime.IsZero() {
		return errors.New("submit time is empty")
	}
	switch action.Kind {
	case PendingActionKindSetEquity:
		return validateUniqueAddresses(action.Addresses)
	case PendingActionKindFreezeEquity, PendingActionKindUnfreezeEquity:
		return validateEquityFreezeAddresses(action.Addresses)
	default:
		return errors.New("unknown pending action kind")
	}
}

func genesisRecordKey(proposalID uint64, address string) string {
	return fmt.Sprintf("%d/%s", proposalID, address)
}
