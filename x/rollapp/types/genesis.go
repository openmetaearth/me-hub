package types

import (
	"errors"
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		RollappList:                        []Rollapp{},
		StateInfoList:                      []StateInfo{},
		LatestStateInfoIndexList:           []StateInfoIndex{},
		LatestFinalizedStateIndexList:      []StateInfoIndex{},
		BlockHeightToFinalizationQueueList: []BlockHeightToFinalizationQueue{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in rollapp
	type rollappEIP155Entry struct {
		chainID ChainID
		rollapp Rollapp
	}

	rollappIndexMap := make(map[string]struct{})
	eip155RollappIndexMap := make(map[uint64]rollappEIP155Entry)

	for _, elem := range gs.RollappList {
		if err := elem.ValidateBasic(); err != nil {
			return err
		}
		index := string(RollappKey(elem.RollappId))
		if _, ok := rollappIndexMap[index]; ok {
			return errors.New("duplicated index for rollapp")
		}
		rollappIndexMap[index] = struct{}{}

		chainID, err := NewChainID(elem.RollappId)
		if err != nil {
			return err
		}
		if !chainID.IsEIP155() {
			continue
		}

		eip155 := chainID.GetEIP155ID()
		if previous, ok := eip155RollappIndexMap[eip155]; ok {
			if !previous.rollapp.Frozen {
				return fmt.Errorf("duplicated eip155 rollapp index %d", eip155)
			}
			if chainID.GetName() != previous.chainID.GetName() {
				return fmt.Errorf("eip155 rollapp %d name must be %s", eip155, previous.chainID.GetName())
			}
			nextRevision := previous.chainID.GetRevisionNumber() + 1
			if chainID.GetRevisionNumber() != nextRevision {
				return fmt.Errorf("eip155 rollapp %d revision number should be %d", eip155, nextRevision)
			}
		}
		eip155RollappIndexMap[eip155] = rollappEIP155Entry{
			chainID: chainID,
			rollapp: elem,
		}
	}
	// Check for duplicated index in stateInfo
	stateInfoIndexMap := make(map[string]struct{})

	for _, elem := range gs.StateInfoList {
		index := string(StateInfoKey(elem.StateInfoIndex))
		if _, ok := stateInfoIndexMap[index]; ok {
			return errors.New("duplicated index for stateInfo")
		}
		stateInfoIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in latestStateInfoIndex
	latestStateInfoIndexIndexMap := make(map[string]struct{})

	for _, elem := range gs.LatestStateInfoIndexList {
		index := string(LatestStateInfoIndexKey(elem.RollappId))
		if _, ok := latestStateInfoIndexIndexMap[index]; ok {
			return errors.New("duplicated index for latestStateInfoIndex")
		}
		latestStateInfoIndexIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in latestFinalizedStateIndex
	latestFinalizedStateIndexIndexMap := make(map[string]struct{})

	for _, elem := range gs.LatestFinalizedStateIndexList {
		index := string(LatestFinalizedStateIndexKey(elem.RollappId))
		if _, ok := latestFinalizedStateIndexIndexMap[index]; ok {
			return errors.New("duplicated index for latestFinalizedStateIndex")
		}
		latestFinalizedStateIndexIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in blockHeightToFinalizationQueue
	blockHeightToFinalizationQueueIndexMap := make(map[string]struct{})

	for _, elem := range gs.BlockHeightToFinalizationQueueList {
		index := string(BlockHeightToFinalizationQueueKey(elem.CreationHeight))
		if _, ok := blockHeightToFinalizationQueueIndexMap[index]; ok {
			return errors.New("duplicated index for blockHeightToFinalizationQueue")
		}
		blockHeightToFinalizationQueueIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	// TODO:

	return gs.Params.Validate()
}
