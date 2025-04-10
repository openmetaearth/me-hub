package cmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethcmd "github.com/evmos/ethermint/cmd/config"
)

const DefaultMaxTxs = 5000

func initSDKConfig() {
	// Set additional config. prefix and denoms registered on app init
	config := sdk.GetConfig()
	ethcmd.SetBip44CoinType(config)
	config.Seal()
}
