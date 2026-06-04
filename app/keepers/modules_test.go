package keepers

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
	"github.com/stretchr/testify/require"
)

func TestGravitySlashingModuleAccountCannotMint(t *testing.T) {
	perms := MaccPerms[gravitytypes.SlashingModuleAccount]
	require.NotContains(t, perms, authtypes.Minter)
	require.Contains(t, perms, authtypes.Burner)
}
