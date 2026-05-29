package v2_0_13

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/module"
	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
	"github.com/stretchr/testify/require"
)

type consensusVersionModule struct {
	version uint64
}

func (m consensusVersionModule) ConsensusVersion() uint64 {
	return m.version
}

func TestMarkInitializedModuleVersionOnlySetsAbsentModule(t *testing.T) {
	mm := &module.Manager{
		Modules: map[string]interface{}{
			bsctypes.ModuleName:  consensusVersionModule{version: 1},
			trontypes.ModuleName: consensusVersionModule{version: 1},
			"existing":           consensusVersionModule{version: 3},
		},
	}
	fromVM := module.VersionMap{
		bsctypes.ModuleName: 0,
		"existing":          2,
	}

	markInitializedModuleVersion(fromVM, mm, bsctypes.ModuleName)
	markInitializedModuleVersion(fromVM, mm, trontypes.ModuleName)

	require.Equal(t, uint64(0), fromVM[bsctypes.ModuleName])
	require.Equal(t, uint64(1), fromVM[trontypes.ModuleName])
	require.Equal(t, uint64(2), fromVM["existing"])
	require.Len(t, fromVM, 3)
}
