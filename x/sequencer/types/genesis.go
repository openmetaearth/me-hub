package types

import fmt "fmt"

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		SequencerList:       []Sequencer{},
		ReplaceProposerList: []MsgStoreReplaceProposer{},
		Params:              DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in sequencer
	sequencerIndexMap := make(map[string]struct{})

	for _, elem := range gs.SequencerList {

		// FIXME: should run validation on the sequencer objects

		index := string(SequencerKey(elem.SequencerAddress))
		if _, ok := sequencerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for sequencer")
		}
		sequencerIndexMap[index] = struct{}{}
	}

	replaceProposerIndexMap := make(map[string]struct{})
	for _, elem := range gs.ReplaceProposerList {
		index := elem.ReplaceProposer.RollappId
		if index == "" {
			return fmt.Errorf("empty index for replaceProposer")
		}
		if _, ok := replaceProposerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for replaceProposer")
		}
		replaceProposerIndexMap[index] = struct{}{}
	}

	// FIXME: validate single PROPOSER per rollapp

	return gs.Params.Validate()
}
