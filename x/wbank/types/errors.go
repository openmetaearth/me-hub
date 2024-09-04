package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var (
	ErrNotGlobalDao = sdkerrors.Register(banktypes.ModuleName, 10, "only global dao can do this")
)
