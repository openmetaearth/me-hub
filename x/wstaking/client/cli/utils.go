package cli

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func buildCommissionRates(rateStr, maxRateStr, maxChangeRateStr string) (commission types.CommissionRates, err error) {
	if rateStr == "" || maxRateStr == "" || maxChangeRateStr == "" {
		return commission, errors.New("must specify all validator commission parameters")
	}

	rate, err := sdkmath.LegacyNewDecFromStr(rateStr)
	if err != nil {
		return commission, err
	}

	maxRate, err := sdkmath.LegacyNewDecFromStr(maxRateStr)
	if err != nil {
		return commission, err
	}

	maxChangeRate, err := sdkmath.LegacyNewDecFromStr(maxChangeRateStr)
	if err != nil {
		return commission, err
	}

	commission = types.NewCommissionRates(rate, maxRate, maxChangeRate)

	return commission, nil
}
