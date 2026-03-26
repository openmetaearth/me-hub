package cmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	appparams "github.com/st-chain/me-hub/app/params"
)

// Set additional config
// prefix and denoms registered on app init
func initSDKConfig() {
	config := sdk.GetConfig()

	appparams.SetAddressPrefixes(config)
	SetBip44CoinType(config)
	config.Seal()

	appparams.RegisterDenoms()
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(ethermint.Bip44CoinType)
	config.SetPurpose(sdk.Purpose) // Shared
}
