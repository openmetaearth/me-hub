package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/testutil/keeper"
	didkeeper "github.com/openmetaearth/me-hub/x/did/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeeper_Filter(t *testing.T) {
	k, ctx := keeper.DidKeeper(t)

	did := "1000000000000001"
	sid := "test"
	vc := didtypes.Credential{
		Did:  did,
		Sid:  sid,
		Hash: "FF000000000000000000",
		Uri:  "https://www.example.com",
		Data: []byte("this is test data"),
	}

	f1 := []byte("fs1")
	f2 := []byte("fs2")
	filters := [][]byte{f1, f2}

	// add
	k.AddFilters(ctx, did, sid, filters, vc)
	// gets
	fs, found := k.GetFilters(ctx, did, sid)
	assert.True(t, found)
	assert.Equal(t, 2, len(fs))
	// gets vc by filter
	pr := query.PageRequest{}
	vcs, _, _ := k.GetCredentialsByFilter(ctx, sid, f1, &pr)
	assert.Equal(t, 1, len(vcs))
	vcs, _, _ = k.GetCredentialsByFilter(ctx, sid, f2, &pr)
	assert.Equal(t, 1, len(vcs))
	// delete
	k.DeleteFilters(ctx, did, sid, filters)
	fs, found = k.GetFilters(ctx, did, sid)
	assert.False(t, found)
	assert.Equal(t, 0, len(fs))
	vcs, _, _ = k.GetCredentialsByFilter(ctx, sid, f1, &pr)
	assert.Equal(t, 0, len(vcs))
	vcs, _, _ = k.GetCredentialsByFilter(ctx, sid, f2, &pr)
	assert.Equal(t, 0, len(vcs))
}

func TestMsgServer_UpdateVCRemovesStaleFilterIndex(t *testing.T) {
	k, ctx := keeper.DidKeeper(t)

	issuerAddr := sdk.MustAccAddressFromBech32("me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4")
	issuerDid := "1000000000001"
	holderDid := "1000000000002"
	sid := "test"
	oldFilter := []byte("old-kyc-scope")
	newFilter := []byte("new-kyc-scope")

	k.SetDID(ctx, issuerAddr, issuerDid)
	k.SetDidInfo(ctx, issuerDid, didtypes.NewDidInfo(issuerDid, issuerAddr.String(), "issuer-pubkey", didtypes.DID_STATUS_ACTIVE))
	k.SetDidInfo(ctx, holderDid, didtypes.NewDidInfo(holderDid, "holder-address", "holder-pubkey", didtypes.DID_STATUS_ACTIVE))
	k.SetService(ctx, sid, didtypes.NewService(sid, "test", "this is a test service.", didtypes.SERVICE_STATUS_ACTIVE, []string{issuerDid}))

	msgServer := didkeeper.NewMsgServerImpl(k)
	goCtx := sdk.WrapSDKContext(ctx)

	_, err := msgServer.CreateVC(goCtx, &didtypes.MsgCreateVC{
		Issuer:  issuerAddr.String(),
		Did:     holderDid,
		Sid:     sid,
		Hash:    "hash-before",
		Uri:     "https://www.example.com/before",
		Filters: [][]byte{oldFilter},
	})
	require.NoError(t, err)

	_, err = msgServer.UpdateVC(goCtx, &didtypes.MsgUpdateVC{
		Issuer:  issuerAddr.String(),
		Did:     holderDid,
		Sid:     sid,
		Hash:    "hash-after",
		Uri:     "https://www.example.com/after",
		Filters: [][]byte{newFilter},
	})
	require.NoError(t, err)

	pr := query.PageRequest{}
	vcs, _, err := k.GetCredentialsByFilter(ctx, sid, oldFilter, &pr)
	require.NoError(t, err)
	assert.Equal(t, 0, len(vcs))

	vcs, _, err = k.GetCredentialsByFilter(ctx, sid, newFilter, &pr)
	require.NoError(t, err)
	require.Equal(t, 1, len(vcs))
	assert.Equal(t, "hash-after", vcs[0].Hash)
}
