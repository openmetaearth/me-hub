package types

import (
	"errors"
	"fmt"
)

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
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Check for duplicated index in sequencer
	sequencerIndexMap := make(map[string]struct{})
	proposerByRollapp := make(map[string]string)

	for _, elem := range gs.SequencerList {
		// FIXME: should run validation on the sequencer objects

		index := string(SequencerKey(elem.SequencerAddress))
		if _, ok := sequencerIndexMap[index]; ok {
			return errors.New("duplicated index for sequencer")
		}
		sequencerIndexMap[index] = struct{}{}

		if elem.Status == Bonded || elem.Status == Unbonding {
			if !elem.Tokens.IsValid() || elem.Tokens.Empty() {
				return fmt.Errorf("sequencer %s has no slashable bond", elem.SequencerAddress)
			}

			if elem.Tokens.AmountOf(gs.Params.MinBond.Denom).LT(gs.Params.MinBond.Amount) {
				return fmt.Errorf(
					"sequencer %s bond is below min bond: got %s, expected at least %s",
					elem.SequencerAddress,
					elem.Tokens,
					gs.Params.MinBond,
				)
			}
		}

		if elem.IsBonded() && elem.IsProposer() {
			if proposer, ok := proposerByRollapp[elem.RollappId]; ok {
				return fmt.Errorf("multiple bonded proposers for rollapp %s: %s and %s", elem.RollappId, proposer, elem.SequencerAddress)
			}
			proposerByRollapp[elem.RollappId] = elem.SequencerAddress
		}
	}

	return nil
}
