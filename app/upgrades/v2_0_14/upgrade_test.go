package v2_0_14

import (
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"
)

type consensusVersionModule struct {
	version uint64
}

func (m consensusVersionModule) ConsensusVersion() uint64 {
	return m.version
}

func TestCreateUpgradeHandlerPreservesFromVersionForMigrations(t *testing.T) {
	const moduleName = "tracked"

	mm := &module.Manager{
		Modules: map[string]interface{}{
			moduleName: consensusVersionModule{version: 2},
		},
		OrderMigrations: []string{moduleName},
	}

	configurator := module.NewConfigurator(nil, nil, nil)
	migrated := false
	require.NoError(t, configurator.RegisterMigration(moduleName, 1, func(sdk.Context) error {
		migrated = true
		return nil
	}))

	fromVM := module.VersionMap{moduleName: 1}
	handler := CreateUpgradeHandler(mm, configurator, nil, nil)
	ctx := sdk.NewContext(nil, tmproto.Header{}, false, log.NewNopLogger())

	updatedVM, err := handler(ctx, upgradetypes.Plan{}, fromVM)

	require.NoError(t, err)
	require.True(t, migrated)
	require.Equal(t, uint64(1), fromVM[moduleName])
	require.Equal(t, uint64(2), updatedVM[moduleName])
}
