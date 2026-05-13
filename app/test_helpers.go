package app

import (
	pruningtypes "cosmossdk.io/store/pruning/types"
	db "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	appparams "github.com/openmetaearth/me-hub/app/params"
)

// NewTestNetworkFixture returns a new network.TestFixture for integration tests.
func NewTestNetworkFixture() network.TestFixture {
	testEncodingConfig := appparams.MakeEncodingConfig()
	encodingConfig := appparams.EncodingConfig{
		InterfaceRegistry: testEncodingConfig.InterfaceRegistry,
		Codec:             testEncodingConfig.Codec,
		TxConfig:          testEncodingConfig.TxConfig,
		Amino:             testEncodingConfig.Amino,
	}

	return network.TestFixture{
		AppConstructor: func(val network.ValidatorI) servertypes.Application {
			return New(
				val.GetCtx().Logger,
				db.NewMemDB(),
				nil,
				true,
				sims.EmptyAppOptions{},
				baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
				baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			)
		},
		GenesisState: NewDefaultGenesisState(encodingConfig.Codec),
	}
}
