package ucoin

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMulDec(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		d, _ := math.LegacyNewDecFromStr("0.5")
		coins := sdk.NewCoins(
			sdk.NewCoin("foo", math.NewInt(2)),
			sdk.NewCoin("bar", math.NewInt(3)),
		)
		res := MulDec(d, coins...)
		require.Equal(t, math.NewInt(1), res[0].Amount)
		require.Equal(t, math.NewInt(1), res[0].Amount)
	})
}
