package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrSbtExists      = errors.Register(ModuleName, 100, "SBT already exists")
	ErrSbtNotFound    = errors.Register(ModuleName, 101, "SBT not found")
	ErrInvalidPubkey  = errors.Register(ModuleName, 102, "invalid pubkey")
	ErrTransferRegion = errors.Register(ModuleName, 103, "transfer region")
	ErrInviteReward   = errors.Register(ModuleName, 104, "send inviter reward failed")
)
