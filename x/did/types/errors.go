package types

import (
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrParameter        = errors.Register(sdkerrors.RootCodespace, 100, "parameter error")
	ErrApiInactive      = errors.Register(ModuleName, 101, "api is inactive")
	ErrPermissionDenial = errors.Register(ModuleName, 102, "permission denial")

	ErrDidExists     = errors.Register(ModuleName, 110, "DID already exists")
	ErrDidNotFound   = errors.Register(ModuleName, 111, "DID not found")
	ErrDidNotActive  = errors.Register(ModuleName, 112, "DID not active")
	ErrSameDidStatus = errors.Register(ModuleName, 113, "same DID status")

	ErrServiceExists     = errors.Register(ModuleName, 120, "credential service already exists")
	ErrServiceNotFound   = errors.Register(ModuleName, 121, "credential service not found")
	ErrServiceNotActive  = errors.Register(ModuleName, 122, "credential service not active")
	ErrSameServiceStatus = errors.Register(ModuleName, 123, "same credential service status")

	ErrIssuerExists    = errors.Register(ModuleName, 130, "issuer already exists")
	ErrIssuerNotFound  = errors.Register(ModuleName, 131, "issuer not found")
	ErrIssuerNotActive = errors.Register(ModuleName, 132, "issuer not active")
	ErrInvalidIssuer   = errors.Register(ModuleName, 133, "invalid issuer")

	ErrHolderExists    = errors.Register(ModuleName, 140, "holder already exists")
	ErrHolderNotFound  = errors.Register(ModuleName, 141, "holder not found")
	ErrHolderNotActive = errors.Register(ModuleName, 142, "holder not active")

	ErrCredentialExists   = errors.Register(ModuleName, 150, "credential already exists")
	ErrCredentialNotFound = errors.Register(ModuleName, 151, "credential not found")
)
