package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sudo module sentinel errors
var (
	ErrCreatorNotDao              = sdkerrors.Register(ModuleName, 1, "creator is not dao error")
	ErrLastAddressEqualNewAddress = sdkerrors.Register(ModuleName, 2, "last address euqal new address error")
	ErrNotFound                   = sdkerrors.Register(ModuleName, 3, "not found")
	ErrSetKycIssuer               = sdkerrors.Register(ModuleName, 4, "set kyc issuer")
	ErrFreeGasAccountAlreadyExist = sdkerrors.Register(ModuleName, 5, "free gas account already exist")
	ErrAccountIsNotFree           = sdkerrors.Register(ModuleName, 6, "account is already not free")
)
