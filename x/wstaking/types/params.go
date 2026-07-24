package types

import (
	"errors"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/app/params"
)

const (
	MinDelegateAmount = "0.01"
	GlobalRegion      = "ME_EARTH"
)

func CheckMinDelegate(amount math.Int) error {
	// if amount.Denom == sdk.BaseMEDenom {
	delAmount := sdkmath.LegacyNewDecFromInt(amount).Mul(sdkmath.LegacyNewDecWithPrec(1, params.BaseDenomUnit))
	minAmount, _ := sdkmath.LegacyNewDecFromStr(MinDelegateAmount)
	if delAmount.LT(minAmount) {
		errStr := fmt.Sprintf("minimum delegate amount is %s, delegate value is %s",
			MinDelegateAmount, amount.String())
		return errors.New(errStr)
	}
	//}
	return nil
}
