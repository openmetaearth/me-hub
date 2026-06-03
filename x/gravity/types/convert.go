package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CheckBscUsdtUsdc(symbol, chainName string) bool {
	return (strings.ToLower(symbol) == "usdt" || strings.ToLower(symbol) == "usdc") && strings.ToLower(chainName) == "bsc"
}

func GetDecimals(claim *MsgBridgeTokenClaim) (decimals uint32) {
	decimals = uint32(claim.Decimals)
	if CheckBscUsdtUsdc(claim.Symbol, claim.ChainName) {
		decimals = uint32(6)
	}
	return decimals
}

func GetMintCoin(amount sdk.Int, chainName string, bridgeToken *BridgeToken) sdk.Coin {
	mintAmount := GetMintAmount(amount, chainName, bridgeToken)
	coin := sdk.NewCoin(bridgeToken.Denom, mintAmount)
	return coin
}

func GetMintAmount(amount sdk.Int, chainName string, bridgeToken *BridgeToken) sdk.Int {
	if requiresBscUsdtUsdcDownConversion(chainName, bridgeToken) {
		amount = amount.Quo(bscUsdtUsdcConversionFactor(chainName, bridgeToken))
	}
	return amount
}

func GetExternalUnlockAmount(amount sdk.Int, chainName string, bridgeToken *BridgeToken) sdk.Int {
	if requiresBscUsdtUsdcDownConversion(chainName, bridgeToken) {
		amount = amount.Mul(bscUsdtUsdcConversionFactor(chainName, bridgeToken))
	}
	return amount
}

func requiresBscUsdtUsdcDownConversion(chainName string, bridgeToken *BridgeToken) bool {
	return bridgeToken != nil && CheckBscUsdtUsdc(bridgeToken.Symbol, chainName) && bridgeToken.Decimal > 6
}

func bscUsdtUsdcConversionFactor(chainName string, bridgeToken *BridgeToken) sdk.Int {
	if !requiresBscUsdtUsdcDownConversion(chainName, bridgeToken) {
		return sdk.OneInt()
	}
	return sdk.NewDec(10).Power(bridgeToken.Decimal - 6).TruncateInt()
}

// ValidateMintAmount rejects BSC USDT/USDC deposits that would lose precision during 6-decimal minting.
func ValidateMintAmount(amount sdk.Int, chainName string, bridgeToken *BridgeToken) error {
	if !requiresBscUsdtUsdcDownConversion(chainName, bridgeToken) {
		return nil
	}

	convert := bscUsdtUsdcConversionFactor(chainName, bridgeToken)
	if amount.IsNil() || !amount.IsPositive() {
		return errorsmod.Wrap(ErrInvalid, "converted BSC USDT/USDC deposit amount must be positive")
	}
	if amount.LT(convert) {
		return errorsmod.Wrapf(
			ErrInvalid,
			"converted BSC USDT/USDC deposit amount %s rounds down to zero with factor %s",
			amount.String(),
			convert.String(),
		)
	}
	if !amount.Mod(convert).IsZero() {
		return errorsmod.Wrapf(
			ErrInvalid,
			"converted BSC USDT/USDC deposit amount %s is not divisible by factor %s",
			amount.String(),
			convert.String(),
		)
	}
	return nil
}
