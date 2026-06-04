package keeper_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/openmetaearth/me-hub/testutil/keeper"
	didkeeper "github.com/openmetaearth/me-hub/x/did/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	kyctypes "github.com/openmetaearth/me-hub/x/kyc/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerRemoveVCRejectsKycReservedSid(t *testing.T) {
	k, ctx := keepertest.DidKeeper(t)
	msgServer := didkeeper.NewMsgServerImpl(k)

	issuerAddr := sdk.MustAccAddressFromBech32("me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4")
	issuerDid := "0000000000000001"
	holderDid := "1000000000000001"
	regionFilter := []byte("me_earth")

	k.SetDID(ctx, issuerAddr, issuerDid)
	k.SetDidInfo(ctx, issuerDid, didtypes.DidInfo{
		Did:    issuerDid,
		Status: didtypes.DID_STATUS_ACTIVE,
	})
	k.SetDidInfo(ctx, holderDid, didtypes.DidInfo{
		Did:    holderDid,
		Status: didtypes.DID_STATUS_ACTIVE,
	})
	k.SetService(ctx, kyctypes.ModuleName, didtypes.Service{
		Sid:     kyctypes.ModuleName,
		Issuers: []string{issuerDid},
		Status:  didtypes.SERVICE_STATUS_ACTIVE,
	})

	credential := didtypes.NewCredential(holderDid, kyctypes.ModuleName, "hash", "uri", regionFilter)
	k.SetCredential(ctx, holderDid, kyctypes.ModuleName, credential)
	k.AddFilters(ctx, holderDid, kyctypes.ModuleName, [][]byte{regionFilter}, credential)

	_, err := msgServer.RemoveVC(ctx, didtypes.NewMsgRemoveVC(issuerAddr.String(), holderDid, kyctypes.ModuleName))

	require.Error(t, err)
	require.True(t, errorsmod.IsOf(err, didtypes.ErrReservedCredentialService))
	require.True(t, k.HasCredential(ctx, holderDid, kyctypes.ModuleName))
	filters, found := k.GetFilters(ctx, holderDid, kyctypes.ModuleName)
	require.True(t, found)
	require.Equal(t, [][]byte{regionFilter}, filters)
}
