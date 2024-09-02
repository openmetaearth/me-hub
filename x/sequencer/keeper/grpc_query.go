package keeper

import (
	"github.com/st-chain/me-hub/x/sequencer/types"
)

var _ types.QueryServer = Keeper{}
