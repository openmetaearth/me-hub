package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/sudo module sentinel errors
var (
	ErrCreatorNotDao              = errorsmod.Register(ModuleName, 1, "creator is not dao error")
	ErrLastAddressEqualNewAddress = errorsmod.Register(ModuleName, 2, "last address euqal new address error")
	ErrNotFound                   = errorsmod.Register(ModuleName, 3, "not found")
	ErrSetKycIssuer               = errorsmod.Register(ModuleName, 4, "set kyc issuer")
	ErrFreeGasAccountAlreadyExist = errorsmod.Register(ModuleName, 5, "free gas account already exist")
	ErrAccountIsNotFree           = errorsmod.Register(ModuleName, 6, "account is already not free")
)
