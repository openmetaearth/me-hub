package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type WminkHooks interface {
	GetPerBlockMintCoinAmount(ctx sdk.Context) (amount big.Int)
}
