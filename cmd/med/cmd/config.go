package cmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethermint "github.com/evmos/ethermint/types"

	"github.com/openmetaearth/me-hub/app/params"
)

func initSDKConfig() {
	// Denoms are registered by params.init when the package is loaded.
	config := sdk.GetConfig()
	params.SetAddressPrefixes()
	config.SetCoinType(ethermint.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)
	config.Seal()
}
