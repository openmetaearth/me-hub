package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/wdistri module sentinel errors
var (
	ErrRegionNotExist          = sdkerrors.Register(ModuleName, 1100, "region not exist")
	ErrMeidNotExists           = sdkerrors.Register(ModuleName, 1101, "MEID not exist")
	ErrEmptyDelegationDistInfo = sdkerrors.Register(ModuleName, 1102, "Empty Delegation Dist Info")
	ErrAssertionFailed         = sdkerrors.Register(ModuleName, 1103, "Assertion Failed")
	ErrCalculateInterest       = sdkerrors.Register(ModuleName, 1104, "delegator calculate interest err.")
	ErrAssertDelegation        = sdkerrors.Register(ModuleName, 1105, "The delegation structure assertion error.")
	ErrUnknownAccount          = sdkerrors.Register(ModuleName, 1106, "Unknown account")
	ErrDistributionIncome      = sdkerrors.Register(ModuleName, 1107, "distribution income err.")
	ErrDistributionOther       = sdkerrors.Register(ModuleName, 1108, "distribution err.")
	ErrPermissionDenied         = sdkerrors.Register(ModuleName, 1111, "permission denied")
	ErrInvalidParams            = sdkerrors.Register(ModuleName, 1112, "invalid params")
)
