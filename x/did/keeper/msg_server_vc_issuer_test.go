package keeper_test

import (
	"testing"

	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/testutil/keeper"
	didkeeper "github.com/openmetaearth/me-hub/x/did/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/stretchr/testify/require"
)

func didTestAddress(seed byte) sdk.AccAddress {
	priv := secp256k1.GenPrivKeyFromSecret([]byte{seed})
	return sdk.AccAddress(priv.PubKey().Address())
}

func TestUpdateVCRejectsDifferentActiveIssuer(t *testing.T) {
	k, ctx := keeper.DidKeeper(t)
	msgServer := didkeeper.NewMsgServerImpl(k)

	issuerAAddr := didTestAddress(1)
	issuerBAddr := didTestAddress(2)
	issuerADid := "1000000000001"
	issuerBDid := "1000000000002"
	holderDid := "1000000000003"
	sid := "kyc"

	k.SetDID(ctx, issuerAAddr, issuerADid)
	k.SetDID(ctx, issuerBAddr, issuerBDid)
	k.SetDidInfo(ctx, issuerADid, didtypes.NewDidInfo(issuerADid, issuerAAddr.String(), "", didtypes.DID_STATUS_ACTIVE))
	k.SetDidInfo(ctx, issuerBDid, didtypes.NewDidInfo(issuerBDid, issuerBAddr.String(), "", didtypes.DID_STATUS_ACTIVE))
	k.SetDidInfo(ctx, holderDid, didtypes.NewDidInfo(holderDid, "", "", didtypes.DID_STATUS_ACTIVE))
	k.SetService(ctx, sid, didtypes.NewService(sid, "KYC", "test service", didtypes.SERVICE_STATUS_ACTIVE, []string{issuerADid, issuerBDid}))

	goCtx := sdk.WrapSDKContext(ctx)
	_, err := msgServer.CreateVC(goCtx, didtypes.NewMsgCreateVC(issuerAAddr.String(), holderDid, sid, "hash-from-issuer-a", "ipfs://issuer-a", nil, nil))
	require.NoError(t, err)

	_, err = msgServer.UpdateVC(goCtx, didtypes.NewMsgUpdateVC(issuerBAddr.String(), holderDid, sid, "hash-from-issuer-b", "ipfs://issuer-b", nil, nil))
	require.ErrorIs(t, err, didtypes.ErrInvalidIssuer)

	vc, found := k.GetCredential(ctx, holderDid, sid)
	require.True(t, found)
	require.Equal(t, "hash-from-issuer-a", vc.Hash)
	require.Equal(t, "ipfs://issuer-a", vc.Uri)
}

func TestUpdateVCAllowsOriginalIssuer(t *testing.T) {
	k, ctx := keeper.DidKeeper(t)
	msgServer := didkeeper.NewMsgServerImpl(k)

	issuerAddr := didTestAddress(1)
	issuerDid := "1000000000001"
	holderDid := "1000000000003"
	sid := "kyc"

	k.SetDID(ctx, issuerAddr, issuerDid)
	k.SetDidInfo(ctx, issuerDid, didtypes.NewDidInfo(issuerDid, issuerAddr.String(), "", didtypes.DID_STATUS_ACTIVE))
	k.SetDidInfo(ctx, holderDid, didtypes.NewDidInfo(holderDid, "", "", didtypes.DID_STATUS_ACTIVE))
	k.SetService(ctx, sid, didtypes.NewService(sid, "KYC", "test service", didtypes.SERVICE_STATUS_ACTIVE, []string{issuerDid}))

	goCtx := sdk.WrapSDKContext(ctx)
	_, err := msgServer.CreateVC(goCtx, didtypes.NewMsgCreateVC(issuerAddr.String(), holderDid, sid, "hash-original", "ipfs://original", nil, nil))
	require.NoError(t, err)

	_, err = msgServer.UpdateVC(goCtx, didtypes.NewMsgUpdateVC(issuerAddr.String(), holderDid, sid, "hash-updated", "ipfs://updated", nil, nil))
	require.NoError(t, err)

	vc, found := k.GetCredential(ctx, holderDid, sid)
	require.True(t, found)
	require.Equal(t, "hash-updated", vc.Hash)
	require.Equal(t, "ipfs://updated", vc.Uri)
}
