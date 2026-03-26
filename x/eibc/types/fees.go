package types

import (
	sdkmath "cosmossdk.io/math"
)

// calculate the new price: transferTotal - fee - bridgingFee
func CalcPriceWithBridgingFee(amt sdkmath.Int, feeInt sdkmath.Int, bridgingFeeMultiplier sdkmath.LegacyDec) (sdkmath.Int, error) {
	bridgingFee := bridgingFeeMultiplier.MulInt(amt).TruncateInt()
	price := amt.Sub(feeInt).Sub(bridgingFee)
	// Check that the price is positive
	if !price.IsPositive() {
		return sdkmath.ZeroInt(), ErrFeeTooHigh
	}
	return price, nil
}
