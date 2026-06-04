package types

import fmt "fmt"

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		SequencerList: []Sequencer{},
		Params:        DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in sequencer
	sequencerIndexMap := make(map[string]struct{})
	bondedProposersByRollapp := make(map[string]string)

	for _, elem := range gs.SequencerList {

		// FIXME: should run validation on the sequencer objects

		index := string(SequencerKey(elem.SequencerAddress))
		if _, ok := sequencerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for sequencer")
		}
		sequencerIndexMap[index] = struct{}{}

		if elem.Status == Bonded && elem.Proposer {
			if existing, ok := bondedProposersByRollapp[elem.RollappId]; ok {
				return fmt.Errorf(
					"multiple bonded proposers for rollapp %s: %s and %s",
					elem.RollappId,
					existing,
					elem.SequencerAddress,
				)
			}
			bondedProposersByRollapp[elem.RollappId] = elem.SequencerAddress
		}
	}

	return gs.Params.Validate()
}
