package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

type DaoKeeper interface {
	GetAirdropAddress(ctx sdk.Context) string
	GetDevOperator(ctx sdk.Context) string
	GetGlobalDao(ctx sdk.Context) string
	GetMeidDao(ctx sdk.Context) string
}

type BankKeeper interface {
	FeeToReceivers(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error
}

type StakingKeeper interface {
	GetMeid(ctx sdk.Context, account string) (val types.Meid, found bool)
	GetValOwnerAddress(ctx sdk.Context, meidAddress string) (string, error)
	GetProposerOwnerAddress(ctx sdk.Context) (string, error)

	GetGlobalAdminFeePoolAddr(ctx sdk.Context) sdk.AccAddress
}

type DidKeeper interface {
	GetDID(ctx sdk.Context, addr sdk.AccAddress) (string, bool)
}
