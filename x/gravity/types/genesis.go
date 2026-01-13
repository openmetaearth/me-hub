package types

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (m *GenesisState) ValidateBasic() error {
	if err := m.Params.ValidateBasic(); err != nil {
		return err
	}
	return nil
}
