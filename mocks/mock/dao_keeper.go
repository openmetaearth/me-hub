package mock

import sdk "github.com/cosmos/cosmos-sdk/types"

type daoKeeper interface {
	GetMeidDao(ctx sdk.Context) sdk.AccAddress
	IsGlobalDao(ctx sdk.Context, address string) bool
}

var _ daoKeeper = &MockDAOKeeper{}

type MockDAOKeeper struct{}

func (m *MockDAOKeeper) GetMeidDao(ctx sdk.Context) sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32("cosmos1lugrmnrk3ngky85n3hsrxumr3ca7m643h59t72")
	if err != nil {
		panic(err)
	}
	return acc
}
func (m *MockDAOKeeper) IsGlobalDao(ctx sdk.Context, address string) bool {
	return true
}
