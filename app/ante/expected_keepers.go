package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	wbanktypes "github.com/st-chain/me-hub/x/wbank/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

type DaoKeeper interface {
	GetAirdropAddress(ctx sdk.Context) string
	GetDevOperator(ctx sdk.Context) string
	GetGlobalDao(ctx sdk.Context) string
	GetMeidDao(ctx sdk.Context) string
	GetGlobalDaoFeePoolAddr(ctx sdk.Context) sdk.AccAddress
}

type BankKeeper interface {
	FeeToReceivers(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output, receiverTypes []wbanktypes.FeeReceiverType) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type StakingKeeper interface {
	GetMeid(ctx sdk.Context, account string) (val types.Meid, found bool)
	GetValOwnerAddress(ctx sdk.Context, meidAddress string) (string, error)
	GetProposerOwnerAddress(ctx sdk.Context) (string, error)
}

type DidKeeper interface {
	GetDID(ctx sdk.Context, addr sdk.AccAddress) (string, bool)
}
