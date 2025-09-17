package types

import (
	gravitytypes "github.com/st-chain/me-hub/x/gravity/types"
)

func DefaultGenesisState() *gravitytypes.GenesisState {
	params := gravitytypes.DefaultParams()
	params.GravityId = "me-bsc-bridge"
	params.AverageExternalBlockTime = 3_000
	return &gravitytypes.GenesisState{
		Params: params,
	}
}
