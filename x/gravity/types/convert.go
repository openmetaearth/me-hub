package types

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var BaseConvert = sdkmath.NewInt(1_000_000_000_000) // 10^12

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
	if CheckBscUsdtUsdc(bridgeToken.Symbol, chainName) {
		amount = amount.Quo(BaseConvert)
	}
	coin := sdk.NewCoin(bridgeToken.Denom, amount)
	return coin
}

func GetExternalUnlockAmount(amount sdk.Int, chainName string, bridgeToken *BridgeToken) sdk.Int {
	if CheckBscUsdtUsdc(bridgeToken.Symbol, chainName) {
		amount = amount.Mul(BaseConvert)
	}
	return amount
}
