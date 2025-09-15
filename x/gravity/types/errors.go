package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalid                 = errorsmod.Register(ModuleName, 2, "invalid")
	ErrEmpty                   = errorsmod.Register(ModuleName, 3, "empty")
	ErrUnknown                 = errorsmod.Register(ModuleName, 4, "unknown")
	ErrDuplicate               = errorsmod.Register(ModuleName, 5, "duplicate")
	ErrNonContiguousEventNonce = errorsmod.Register(ModuleName, 6, "non contiguous event nonce")
	ErrNotFound                = errorsmod.Register(ModuleName, 7, "not found")

	ErrNotProposedRelayer = errorsmod.Register(ModuleName, 7, "not proposed relayer")
	ErrNotFoundRelayer    = errorsmod.Register(ModuleName, 8, "not found relayer")
	ErrRelayerNotOnLine   = errorsmod.Register(ModuleName, 9, "relayer not on line")

	ErrDelegateAmountBelowMinimum = errorsmod.Register(ModuleName, 10, "delegate amount must be greater than relayer stake threshold")
	ErrDelegateAmountAboveMaximum = errorsmod.Register(ModuleName, 11, "delegate amount must be less than double relayer stake threshold")
)
