package network

import (
	"testing"

	"cosmossdk.io/store/pruning/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	"github.com/openmetaearth/me-hub/app"
)

type (
	Network = network.Network
	Config  = network.Config
)

// New creates instance with fully configured cosmos network.
// Accepts optional config, that will be used in place of the DefaultConfig() if provided.
func New(t *testing.T, configs ...network.Config) *network.Network {
	if len(configs) > 1 {
		panic("at most one config should be provided")
	}
	var cfg network.Config
	if len(configs) == 0 {
		cfg = DefaultConfig()
	} else {
		cfg = configs[0]
	}
	net, err := network.New(t, t.TempDir(), cfg)
	require.NoError(t, err)

	t.Cleanup(net.Cleanup)
	return net
}

// DefaultConfig will initialize config for the network with custom application,
// genesis and single validator. All other parameters are inherited from cosmos-sdk/testutil/network.DefaultConfig
func DefaultConfig() network.Config {
	encoding := app.MakeEncodingConfig()
	cfg := network.DefaultConfig(func() network.TestFixture {
		return network.TestFixture{
			EncodingConfig: moduletestutil.TestEncodingConfig{
				InterfaceRegistry: encoding.InterfaceRegistry,
				Codec:             encoding.Codec, TxConfig: encoding.TxConfig, Amino: encoding.Amino,
			},
			GenesisState: app.NewDefaultGenesisState(encoding.Codec),
			AppConstructor: func(val network.ValidatorI) servertypes.Application {
				return app.New(val.GetCtx().Logger, dbm.NewMemDB(), nil, true, map[int64]bool{},
					app.DefaultNodeHome, 0, encoding, val.GetCtx().Viper,
					baseapp.SetPruning(types.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
					baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
					baseapp.SetChainID("me_1000-1"))
			},
		}
	})
	cfg.ChainID = "me_1000-1"
	cfg.NumValidators = 1
	return cfg
}
