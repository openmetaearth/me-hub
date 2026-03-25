package app

import (
	"time"

	pruningtypes "cosmossdk.io/store/pruning/types"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
)

// NewTestNetworkFixture returns a new network.TestFixture for integration tests.
func NewTestNetworkFixture() network.TestFixture {
	encodingConfig := MakeEncodingConfig()

	return network.TestFixture{
		AppConstructor: func(val network.ValidatorI) servertypes.Application {
			return New(
				val.GetCtx().Logger,
				dbm.NewMemDB(),
				nil,
				true,
				map[int64]bool{},
				val.GetCtx().Config.RootDir,
				0,
				encodingConfig,
				sims.EmptyAppOptions{},
				baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
				baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			)
		},
		GenesisState:  NewDefaultGenesisState(encodingConfig.Codec),
		TimeoutCommit: 2 * time.Second,
	}
}
