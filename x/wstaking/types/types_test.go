package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func TestPowerReduction(t *testing.T) {
	t.Log(sdk.DefaultPowerReduction.String())

}
