package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNotGlobalDao = sdkerrors.Register(ModuleName, 10, "only global dao can do this")
)
