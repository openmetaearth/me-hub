package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

type DaoKeeper interface {
	IsGlobalDao(ctx sdk.Context, address string) bool
}

type EVMKeeper interface {
	StaticCall(ctx sdk.Context, caller, contract common.Address, input []byte, gas uint64) ([]byte, error)
}

type  AccountKeeper  interface {
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
}
