package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MaxDecimals = uint64(18)

func CheckBscUsdtUsdc(symbol, chainName string) bool {
	return (strings.ToLower(symbol) == "usdt" || strings.ToLower(symbol) == "usdc") && strings.ToLower(chainName) == "bsc"
}

func GetDecimals(claim *MsgBridgeTokenClaim) (decimals uint32) {
	decimals = uint32(claim.Decimals)
	if CheckBscUsdtUsdc(claim.Symbol, claim.ChainName) {
		decimals = uint32(6)
	} else if decimals > uint32(MaxDecimals) {
		decimals = uint32(MaxDecimals)
	}
	return decimals
}

func GetMintCoin(amount sdk.Int, chainName string, bridgeToken *BridgeToken) sdk.Coin {
	mintAmount := GetMintAmount(amount, chainName, bridgeToken)
	coin := sdk.NewCoin(bridgeToken.Denom, mintAmount)
	return coin
}

func GetMintAmount(amount sdk.Int, chainName string, bridgeToken *BridgeToken) sdk.Int {
	if CheckBscUsdtUsdc(bridgeToken.Symbol, chainName) && bridgeToken.Decimal > 6 {
		decimals := bridgeToken.Decimal
		if decimals > MaxDecimals {
			decimals = MaxDecimals
		}
		if decimals > 6 {
			convert := sdk.NewDec(10).Power(decimals - 6).TruncateInt()
			// Prevent truncation-to-zero: only convert if amount is large enough
			if amount.GTE(convert) {
				amount = amount.Quo(convert)
			}
		}
	}
	return amount
}

func GetExternalUnlockAmount(amount sdk.Int, chainName string, bridgeToken *BridgeToken) sdk.Int {
	if CheckBscUsdtUsdc(bridgeToken.Symbol, chainName) && bridgeToken.Decimal > 6 {
		decimals := bridgeToken.Decimal
		if decimals > MaxDecimals {
			decimals = MaxDecimals
		}
		if decimals > 6 {
			convert := sdk.NewDec(10).Power(decimals - 6).TruncateInt()
			amount = amount.Mul(convert)
		}
	}
	return amount
}
