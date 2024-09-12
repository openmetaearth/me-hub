package types

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error

	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type RollappKeeper interface {
	IsRollappExist(ctx sdk.Context, rollappId string) bool
	RollappsEnabled(ctx sdk.Context) bool
	GetSubmitterLastSubmitTime(goCtx context.Context, rollappId, submitter string) int64
	GetBlockDaCommitTime(goCtx context.Context, rollappId string, startHeight uint64, numbers uint32) int64
}

type DaoKeeper interface {
	GetGlobalDao(ctx sdk.Context) string
}
