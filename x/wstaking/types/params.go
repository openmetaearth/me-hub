package types

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
)

const (
	MinDelegateAmount = "0.01"
	GlobalRegion      = "ME_EARTH"
)

func CheckMinDelegate(amount math.Int) error {
	//if amount.Denom == sdk.BaseMEDenom {
	delAmount := sdk.NewDecFromInt(amount).Mul(sdk.NewDecWithPrec(1, params.BaseDenomUnit))
	minAmount, _ := sdk.NewDecFromStr(MinDelegateAmount)
	if delAmount.LT(minAmount) {
		errStr := fmt.Sprintf("minimum delegate amount is %s, delegate value is %s",
			MinDelegateAmount, amount.String())
		return errors.New(errStr)
	}
	//}
	return nil
}
