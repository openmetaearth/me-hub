package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	wbanktypes "github.com/st-chain/me-hub/x/wbank/types"
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
	GetValOwnerAddress(ctx sdk.Context, regionId string) (string, error)
	GetProposerOwnerAddress(ctx sdk.Context) (string, error)
}

type KycKeeper interface {
	GetDID(ctx sdk.Context, addr sdk.AccAddress) (string, bool)
	GetKYC(ctx sdk.Context, did string) (kyc didtypes.Credential, found bool)
}
