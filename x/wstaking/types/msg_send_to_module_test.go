package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgSendToModuleValidateBasicRejectsDisallowedTarget(t *testing.T) {
	msg := NewMsgSendToModule(
		sample.AccAddress(),
		distrtypes.ModuleName,
		sdk.NewCoins(sdk.NewInt64Coin("umec", 1)),
	)

	err := msg.ValidateBasic()

	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	require.ErrorContains(t, err, "is not an allowed SendToModule target")
}

func TestMsgSendToModuleValidateBasicAllowsApprovedTargets(t *testing.T) {
	for _, target := range []string{
		StakePoolName,
		FixedDepositPrincipalPool,
		BridgeFeePool,
	} {
		msg := NewMsgSendToModule(
			sample.AccAddress(),
			target,
			sdk.NewCoins(sdk.NewInt64Coin("umec", 1)),
		)

		require.NoError(t, msg.ValidateBasic())
	}
}
