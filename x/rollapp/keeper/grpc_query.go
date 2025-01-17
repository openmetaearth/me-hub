package keeper

import (
	"github.com/st-chain/me-hub/x/rollapp/types"
)

var _ types.QueryServer = Keeper{}
