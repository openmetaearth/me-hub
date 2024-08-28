package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
)

type WminkHooks interface {
	GetPerBlockMintCoinAmount(ctx sdk.Context) (amount big.Int)
}
