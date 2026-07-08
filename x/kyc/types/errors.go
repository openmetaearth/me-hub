package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrSbtExists                   = errors.Register(ModuleName, 100, "SBT already exists")
	ErrSbtNotFound                 = errors.Register(ModuleName, 101, "SBT not found")
	ErrInvalidPubkey               = errors.Register(ModuleName, 102, "invalid pubkey")
	ErrTransferRegion              = errors.Register(ModuleName, 103, "transfer region")
	ErrInviteReward                = errors.Register(ModuleName, 104, "send inviter reward failed")
	ErrSubAccountAlreadyExists     = errors.Register(ModuleName, 105, "sub account already exists")
	ErrSubAccountAlreadyRegistered = errors.Register(ModuleName, 106, "sub account already registered")
	ErrUnauthorized                = errors.Register(ModuleName, 107, "unauthorized")
	ErrEthAccountNotAllowed        = errors.Register(ModuleName, 108, "main account is an eth account and cannot create sub account")
	ErrMainAccountPubkeyNotSet     = errors.Register(ModuleName, 109, "main account pubkey is not set")
)
