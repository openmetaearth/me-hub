package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalid                 = errorsmod.Register(ModuleName, 2, "invalid")
	ErrEmpty                   = errorsmod.Register(ModuleName, 3, "empty")
	ErrUnknown                 = errorsmod.Register(ModuleName, 4, "unknown")
	ErrDuplicate               = errorsmod.Register(ModuleName, 5, "duplicate")
	ErrNonContinuousEventNonce = errorsmod.Register(ModuleName, 6, "non continuous event nonce")
	ErrNotFound                = errorsmod.Register(ModuleName, 7, "not found")

	ErrNotProposedRelayer      = errorsmod.Register(ModuleName, 8, "not a proposed relayer")
	ErrNotFoundRelayer         = errorsmod.Register(ModuleName, 9, "not found relayer")
	ErrExternalAddressNotMatch = errorsmod.Register(ModuleName, 10, "external address not match relayer")
	ErrRelayerNotOnLine        = errorsmod.Register(ModuleName, 11, "relayer not on line")

	ErrDelegateAmountBelowMinimum  = errorsmod.Register(ModuleName, 12, "delegate amount must be greater than relayer stake threshold")
	ErrDelegateAmountAboveMaximum  = errorsmod.Register(ModuleName, 13, "delegate amount must be less than relayer stake threshold")
	ErrMaxChangePowerLimitExceeded = errorsmod.Register(ModuleName, 14, "max change power limit exceeded")
	ErrDuplicateRelayerConfirms    = errorsmod.Register(ModuleName, 15, "relayer already confirmed")
)
