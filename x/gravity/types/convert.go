package types

import (
	"strings"

	sdkmath "cosmossdk.io/math"
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

func GetMintCoin(amount sdkmath.Int, chainName string, bridgeToken *BridgeToken) sdk.Coin {
	mintAmount := GetMintAmount(amount, chainName, bridgeToken)
	coin := sdk.NewCoin(bridgeToken.Denom, mintAmount)
	return coin
}

func GetMintAmount(amount sdkmath.Int, chainName string, bridgeToken *BridgeToken) sdkmath.Int {
	if CheckBscUsdtUsdc(bridgeToken.Symbol, chainName) && bridgeToken.Decimal > 6 {
		convert := sdkmath.LegacyNewDec(10).Power(bridgeToken.Decimal - 6).TruncateInt()
		amount = amount.Quo(convert)
	}
	return amount
}

func GetExternalUnlockAmount(amount sdkmath.Int, chainName string, bridgeToken *BridgeToken) sdkmath.Int {
	if CheckBscUsdtUsdc(bridgeToken.Symbol, chainName) && bridgeToken.Decimal > 6 {
		convert := sdkmath.LegacyNewDec(10).Power(bridgeToken.Decimal - 6).TruncateInt()
		amount = amount.Mul(convert)
	}
	return amount
}
