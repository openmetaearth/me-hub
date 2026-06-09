package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParamsValidateBasicRejectsZeroPowerMinDelegate(t *testing.T) {
	params := DefaultParams()
	require.NoError(t, params.ValidateBasic())

	params.MinDelegate = sdk.DefaultPowerReduction.Sub(sdkmath.NewInt(1))
	params.MaxDelegate = sdk.DefaultPowerReduction
	require.ErrorContains(t, params.ValidateBasic(), "min delegate threshold must produce non-zero relayer power")

	params.MinDelegate = sdk.DefaultPowerReduction
	params.MaxDelegate = sdk.DefaultPowerReduction
	require.NoError(t, params.ValidateBasic())
}
