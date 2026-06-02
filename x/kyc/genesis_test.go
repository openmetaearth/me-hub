package kyc_test

import (
	"bytes"
	"testing"

	cometbftdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/openmetaearth/me-hub/x/did"
	didkeeper "github.com/openmetaearth/me-hub/x/did/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc"
	kyckeeper "github.com/openmetaearth/me-hub/x/kyc/keeper"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/stretchr/testify/require"
)

type noopNFTKeeper struct{}

var _ kyctypes.NFTKeeper = noopNFTKeeper{}

func (noopNFTKeeper) GetNFT(sdk.Context, string, string) (nft.NFT, bool) {
	return nft.NFT{}, false
}

func (noopNFTKeeper) HasNFT(sdk.Context, string, string) bool {
	return false
}

func (noopNFTKeeper) GetOwner(sdk.Context, string, string) sdk.AccAddress {
	return nil
}

func (noopNFTKeeper) Mint(sdk.Context, nft.NFT, sdk.AccAddress) error {
	return nil
}

func (noopNFTKeeper) Update(sdk.Context, nft.NFT) error {
	return nil
}

func (noopNFTKeeper) Burn(sdk.Context, string, string) error {
	return nil
}

func (noopNFTKeeper) SaveClass(sdk.Context, nft.Class) error {
	return nil
}

func newDidKycKeepers(t *testing.T) (*didkeeper.Keeper, *kyckeeper.Keeper, sdk.Context) {
	t.Helper()

	didStoreKey := sdk.NewKVStoreKey(didtypes.StoreKey)
	kycStoreKey := sdk.NewKVStoreKey(kyctypes.StoreKey)
	db := cometbftdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(didStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(kycStoreKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	didKeeper := didkeeper.NewKeeper(cdc, didStoreKey, nil)
	kycKeeper := kyckeeper.NewKeeper(cdc, kycStoreKey, nil, authkeeper.AccountKeeper{}, didKeeper, noopNFTKeeper{})
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return didKeeper, kycKeeper, ctx
}

func testIssuerAddr() sdk.AccAddress {
	return sdk.AccAddress(bytes.Repeat([]byte{0x01}, 20))
}

func testIssuerInfo() didtypes.DidInfo {
	addr := testIssuerAddr()
	return didtypes.DidInfo{
		Did:     "1000000000001",
		Address: addr.String(),
		Pubkey:  "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AyfZ/7fojbKMioe5Oaw378EH4F8w2CGvZ7SwOCRvlCH8\"}",
		Status:  didtypes.DID_STATUS_ACTIVE,
	}
}

func testKycGenesis(issuerInfo didtypes.DidInfo) kyctypes.GenesisState {
	return kyctypes.GenesisState{Issuers: []didtypes.DidInfo{issuerInfo}}
}

func TestDidKycExportedGenesisImportsWithoutDuplicateIssuerPanic(t *testing.T) {
	sourceDidKeeper, sourceKycKeeper, sourceCtx := newDidKycKeepers(t)
	issuerAddr := testIssuerAddr()
	issuerInfo := testIssuerInfo()
	service := didtypes.Service{
		Sid:         kyctypes.ModuleName,
		Name:        kyctypes.ModuleName,
		Description: "The KYC verifiable credential issuer based The DID(Decentralized Identity).",
		Issuers:     []string{issuerInfo.Did},
		Status:      didtypes.SERVICE_STATUS_ACTIVE,
	}

	sourceDidKeeper.SetDID(sourceCtx, issuerAddr, issuerInfo.Did)
	sourceDidKeeper.SetDidInfo(sourceCtx, issuerInfo.Did, issuerInfo)
	sourceDidKeeper.SetService(sourceCtx, kyctypes.ModuleName, service)
	didGenesis := did.ExportGenesis(sourceCtx, sourceDidKeeper)
	kycGenesis := kyc.ExportGenesis(sourceCtx, *sourceKycKeeper)

	targetDidKeeper, targetKycKeeper, targetCtx := newDidKycKeepers(t)
	did.InitGenesis(targetCtx, targetDidKeeper, *didGenesis)
	require.NotPanics(t, func() {
		kyc.InitGenesis(targetCtx, *targetKycKeeper, *kycGenesis)
	})

	importedDid, found := targetKycKeeper.GetDID(targetCtx, issuerAddr)
	require.True(t, found)
	require.Equal(t, issuerInfo.Did, importedDid)

	importedInfo, found := targetKycKeeper.GetDidInfo(targetCtx, issuerInfo.Did)
	require.True(t, found)
	require.Equal(t, issuerInfo, importedInfo)

	importedService, found := targetKycKeeper.GetService(targetCtx)
	require.True(t, found)
	require.Equal(t, []string{issuerInfo.Did}, importedService.Issuers)
}

func TestKycGenesisPanicsWhenIssuerAddressMapsToDifferentDid(t *testing.T) {
	didKeeper, kycKeeper, ctx := newDidKycKeepers(t)
	issuerInfo := testIssuerInfo()

	didKeeper.SetDID(ctx, testIssuerAddr(), "2000000000002")

	require.Panics(t, func() {
		kyc.InitGenesis(ctx, *kycKeeper, testKycGenesis(issuerInfo))
	})
}

func TestKycGenesisPanicsWhenRestoredIssuerDidInfoIsMissing(t *testing.T) {
	didKeeper, kycKeeper, ctx := newDidKycKeepers(t)
	issuerInfo := testIssuerInfo()

	didKeeper.SetDID(ctx, testIssuerAddr(), issuerInfo.Did)

	require.Panics(t, func() {
		kyc.InitGenesis(ctx, *kycKeeper, testKycGenesis(issuerInfo))
	})
}

func TestKycGenesisPanicsWhenRestoredIssuerDidInfoDiffers(t *testing.T) {
	didKeeper, kycKeeper, ctx := newDidKycKeepers(t)
	issuerInfo := testIssuerInfo()
	conflictingInfo := issuerInfo
	conflictingInfo.Pubkey = "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A83z2Fnur8jc+tGvkCJjkZTBeJDLSObk8nVKOpY9P679\"}"

	didKeeper.SetDID(ctx, testIssuerAddr(), issuerInfo.Did)
	didKeeper.SetDidInfo(ctx, issuerInfo.Did, conflictingInfo)

	require.Panics(t, func() {
		kyc.InitGenesis(ctx, *kycKeeper, testKycGenesis(issuerInfo))
	})
}

func TestKycGenesisPanicsWhenUnmappedIssuerDidInfoDiffers(t *testing.T) {
	didKeeper, kycKeeper, ctx := newDidKycKeepers(t)
	issuerInfo := testIssuerInfo()
	conflictingInfo := issuerInfo
	conflictingInfo.Address = sdk.AccAddress(bytes.Repeat([]byte{0x02}, 20)).String()

	didKeeper.SetDidInfo(ctx, issuerInfo.Did, conflictingInfo)

	require.Panics(t, func() {
		kyc.InitGenesis(ctx, *kycKeeper, testKycGenesis(issuerInfo))
	})
}
