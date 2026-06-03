package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/openmetaearth/me-hub/testutil/keeper"
	didkeeper "github.com/openmetaearth/me-hub/x/did/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServer_RemoveVCRejectsReservedKYCSID(t *testing.T) {
	k, ctx := testkeeper.DidKeeper(t)
	msgServer := didkeeper.NewMsgServerImpl(k)

	issuerAddr := sdk.MustAccAddressFromBech32("me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4")
	holderAddr := sdk.MustAccAddressFromBech32("me13w3mxrd9tvq3r6gzheqjuzf8pnaruvug5787yu")
	issuerDID := "1000000000001"
	holderDID := "1000000000002"
	const kycSID = "kyc"

	k.SetDID(ctx, issuerAddr, issuerDID)
	k.SetDidInfo(ctx, issuerDID, didtypes.NewDidInfo(issuerDID, issuerAddr.String(), "", didtypes.DID_STATUS_ACTIVE))
	k.SetDidInfo(ctx, holderDID, didtypes.NewDidInfo(holderDID, holderAddr.String(), "", didtypes.DID_STATUS_ACTIVE))
	k.SetService(ctx, kycSID, didtypes.NewService(kycSID, kycSID, "reserved KYC service", didtypes.SERVICE_STATUS_ACTIVE, []string{issuerDID}))

	credential := didtypes.NewCredential(holderDID, kycSID, "hash-before-remove", "https://example.com/kyc", []byte("me_earth"))
	k.SetCredential(ctx, holderDID, kycSID, credential)

	_, err := msgServer.RemoveVC(ctx, didtypes.NewMsgRemoveVC(issuerAddr.String(), holderDID, kycSID))
	require.ErrorIs(t, err, didtypes.ErrReservedCredential)

	got, found := k.GetCredential(ctx, holderDID, kycSID)
	require.True(t, found)
	require.Equal(t, credential, got)
}
