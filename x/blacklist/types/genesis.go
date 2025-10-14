package types

// DefaultGenesis returns the default genesis state for the blacklist module
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Blacklist: Blacklist{Addresses: []string{}},
	}
}

// Validate performs basic validation of genesis data
func (gs GenesisState) Validate() error {
	return nil
}
