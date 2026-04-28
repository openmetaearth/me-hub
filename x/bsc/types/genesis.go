package types

import (
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
)

func DefaultGenesisState() *gravitytypes.GenesisState {
	params := gravitytypes.DefaultParams()
	params.GravityId = "me-bsc-bridge"
	params.AverageExternalBlockTime = 750
	return &gravitytypes.GenesisState{
		Params: params,
	}
}
