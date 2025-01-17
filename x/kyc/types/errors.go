package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrSbtExists   = errors.Register(ModuleName, 100, "SBT already exists")
	ErrSbtNotFound = errors.Register(ModuleName, 101, "SBT not found")
)
