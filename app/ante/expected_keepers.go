package ante

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	megrouptypes "github.com/openmetaearth/me-hub/x/megroup/types"
	wbanktypes "github.com/openmetaearth/me-hub/x/wbank/types"
)

type DaoKeeper interface {
	IsDao(ctx sdk.Context, addr string) bool
	GetAirdropAddress(ctx sdk.Context) string
	GetDevOperator(ctx sdk.Context) string
	GetGlobalDaoFeePoolAddr(ctx sdk.Context) sdk.AccAddress
	CheckFreeGasAccount(ctx sdk.Context, address string) bool
	// IsGlobalDao and GetMeidDao are used by JoinGroupValidateDecorator to
	// verify that a non-self creator is a privileged DAO account.
	IsGlobalDao(ctx sdk.Context, address string) bool
	GetMeidDao(ctx sdk.Context) string
}

type BankKeeper interface {
	FeeToReceivers(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output, receiverTypes []wbanktypes.FeeReceiverType) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

type StakingKeeper interface {
	GetValOwnerAddress(ctx sdk.Context, regionId string) (string, error)
	GetProposerOwnerAddress(ctx sdk.Context) (string, error)
}

type KycKeeper interface {
	GetDID(ctx sdk.Context, addr sdk.AccAddress) (string, bool)
	GetKYC(ctx sdk.Context, did string) (kyc didtypes.Credential, found bool)
}

type WasmKeeper interface {
	HasContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) bool
	GetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) *wasmtypes.ContractInfo
}

// MeGroupKeeper is the subset of the megroup keeper needed by JoinGroupValidateDecorator.
type MeGroupKeeper interface {
	GetGroupInfo(ctx sdk.Context, id uint64) (megrouptypes.GroupInfo, bool)
	GetDidAndKycActive(ctx sdk.Context, addr sdk.AccAddress, regionID string) (string, bool)
}
