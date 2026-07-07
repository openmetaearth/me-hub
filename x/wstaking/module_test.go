package wstaking

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"
)

// TestAppModuleOriginalStakingModule_GetConsensusVersion checks that the wrapped module still uses the same
// consensus version.
func TestAppModuleOriginalStakingModule_GetConsensusVersion(t *testing.T) {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	stakingModule := staking.NewAppModule(
		cdc, nil, authkeeper.AccountKeeper{}, bankkeeper.BaseKeeper{}, nil,
	)
	require.Equal(t, uint64(5), stakingModule.ConsensusVersion())
}
