package types

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (m *GenesisState) ValidateBasic() error {
	return m.Params.ValidateBasic()
}
