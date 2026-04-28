package types

import (
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
)

func DefaultGenesisState() *gravitytypes.GenesisState {
	params := gravitytypes.DefaultParams()
	params.GravityId = "me-tron-bridge"
	params.AverageExternalBlockTime = 3_000
	return &gravitytypes.GenesisState{
		Params: params,
	}
}
