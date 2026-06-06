package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CheckBscUsdtUsdc(symbol, chainName string) bool {
	return (strings.ToLower(symbol) == "usdt" || strings.ToLower(symbol) == "usdc") && strings.ToLower(chainName) == "bsc"
}

func GetDecimals(claim *MsgBridgeTokenClaim) (decimals uint32) {
	decimals = uint32(claim.Decimals)
	if claim.Decimals > MaxBridgeTokenDecimals {
		decimals = uint32(MaxBridgeTokenDecimals)
	}
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
	if exponent, ok := bscStablecoinConversionExponent(chainName, bridgeToken); ok {
		convert := sdk.NewDec(10).Power(exponent).TruncateInt()
		amount = amount.Quo(convert)
	}
	return amount
}

func GetExternalUnlockAmount(amount sdk.Int, chainName string, bridgeToken *BridgeToken) sdk.Int {
	if exponent, ok := bscStablecoinConversionExponent(chainName, bridgeToken); ok {
		convert := sdk.NewDec(10).Power(exponent).TruncateInt()
		amount = amount.Mul(convert)
	}
	return amount
}

func bscStablecoinConversionExponent(chainName string, bridgeToken *BridgeToken) (uint64, bool) {
	if !CheckBscUsdtUsdc(bridgeToken.Symbol, chainName) || bridgeToken.Decimal <= 6 || bridgeToken.Decimal > MaxBridgeTokenDecimals {
		return 0, false
	}
	return bridgeToken.Decimal - 6, true
}
