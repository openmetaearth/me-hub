package keeper

import (
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

var _ types.QueryServer = Keeper{}
