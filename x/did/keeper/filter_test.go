package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/testutil/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/stretchr/testify/assert"
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

func TestMsgServer_UpdateVC(t *testing.T) {
	k, ctx := keeper.DidKeeper(t)
	srv := keeper.NewMsgServerImpl(k)

	did := "1000000000000001"
	issuer := "1000000000000002"
	sid := "test-service"

	// setup state
	k.SetDidInfo(ctx, did, didtypes.DidInfo{Did: did, Status: didtypes.DID_STATUS_ACTIVE})
	k.SetDidInfo(ctx, issuer, didtypes.DidInfo{Did: issuer, Status: didtypes.DID_STATUS_ACTIVE})
	issuerAddr := sdk.MustAccAddressFromBech32("me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4")
	k.SetDID(ctx, issuerAddr, issuer)
	k.SetService(ctx, sid, didtypes.Service{
		Sid:     sid,
		Status:  didtypes.SERVICE_STATUS_ACTIVE,
		Issuers: []string{issuer},
	})

	// create initial credential and filters
	f1 := []byte("filter1")
	f2 := []byte("filter2")

	msgCreate := &didtypes.MsgCreateVC{
		Did:     did,
		Sid:     sid,
		Issuer:  issuerAddr.String(),
		Filters: [][]byte{f1, f2},
	}
	_, err := srv.CreateVC(sdk.WrapSDKContext(ctx), msgCreate)
	assert.NoError(t, err)

	// verify filters added
	fs, found := k.GetFilters(ctx, did, sid)
	assert.True(t, found)
	assert.Equal(t, 2, len(fs))

	// update VC with new filters (f3, f4)
	f3 := []byte("filter3")
	f4 := []byte("filter4")
	msgUpdate := &didtypes.MsgUpdateVC{
		Did:     did,
		Sid:     sid,
		Issuer:  issuerAddr.String(),
		Filters: [][]byte{f3, f4},
	}
	_, err = srv.UpdateVC(sdk.WrapSDKContext(ctx), msgUpdate)
	assert.NoError(t, err)

	// verify new filters are active and old ones are deleted
	fs, found = k.GetFilters(ctx, did, sid)
	assert.True(t, found)
	assert.Equal(t, 2, len(fs))
	assert.Contains(t, fs, f3)
	assert.Contains(t, fs, f4)

	// verify old credentials by filter returns 0
	pr := query.PageRequest{}
	vcs, _, _ := k.GetCredentialsByFilter(ctx, sid, f1, &pr)
	assert.Equal(t, 0, len(vcs))
	vcs, _, _ = k.GetCredentialsByFilter(ctx, sid, f2, &pr)
	assert.Equal(t, 0, len(vcs))

	// verify new credentials by filter returns 1
	vcs, _, _ = k.GetCredentialsByFilter(ctx, sid, f3, &pr)
	assert.Equal(t, 1, len(vcs))
}
