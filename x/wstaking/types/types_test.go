package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestPowerReduction(t *testing.T) {
	t.Log(sdk.DefaultPowerReduction.String())
}
