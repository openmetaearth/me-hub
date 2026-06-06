package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
)

type staticAccountGetter struct {
	account authtypes.AccountI
}

func (g staticAccountGetter) GetAccount(sdk.Context, sdk.AccAddress) authtypes.AccountI {
	return g.account
}

func TestValidateDaoAddressNotModuleAccountRejectsModuleAccounts(t *testing.T) {
	moduleAccount := authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)

	err := validateDaoAddressNotModuleAccount(
		sdk.Context{},
		staticAccountGetter{account: moduleAccount},
		"DevOperator",
		moduleAccount.GetAddress().String(),
	)

	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
}

func TestValidateDaoAddressNotModuleAccountAllowsBaseAccounts(t *testing.T) {
	addr := sdk.AccAddress("regular-account-addr")
	baseAccount := authtypes.NewBaseAccountWithAddress(addr)

	err := validateDaoAddressNotModuleAccount(
		sdk.Context{},
		staticAccountGetter{account: baseAccount},
		"DevOperator",
		addr.String(),
	)

	require.NoError(t, err)
}
