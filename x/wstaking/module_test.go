package wstaking

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"
)

// TestAppModuleOriginalStakingModule_GetConsensusVersion checks that the wrapped module still uses the save
// consensus version.
func TestAppModuleOriginalStakingModule_GetConsensusVersion(t *testing.T) {
	stakingModule := staking.NewAppModule(
		&codec.AminoCodec{}, nil, authkeeper.AccountKeeper{}, bankkeeper.BaseKeeper{}, nil,
	)
	require.Equal(t, uint64(4), stakingModule.ConsensusVersion())
}
