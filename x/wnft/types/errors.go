package types

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

// x/nft module sentinel errors
var (
	ErrEmptyTotalSupply = errors.Register(nft.ModuleName, 9, "empty total supply")
	ErrEmptyTokenId     = errors.Register(nft.ModuleName, 10, "empty token id")
	ErrEmptyUri         = errors.Register(nft.ModuleName, 11, "empty uri")
)
