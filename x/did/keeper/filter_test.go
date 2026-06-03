package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/openmetaearth/me-hub/testutil/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeeper_Filter(t *testing.T) {
	k, ctx := keeper.DidKeeper(t)

	did := "1000000000001"
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

func TestKeeper_GetCredentialsByFilterDoesNotMatchLongerRegionPrefix(t *testing.T) {
	k, ctx := keeper.DidKeeper(t)

	sid := "kyc"
	did := "1000000000002"
	parentRegion := []byte("me_earth")
	childRegion := []byte("me_earth-usa")
	vc := didtypes.Credential{
		Did:  did,
		Sid:  sid,
		Hash: "hash-earth-usa",
		Uri:  "https://www.example.com/earth-usa",
	}

	k.AddFilters(ctx, did, sid, [][]byte{childRegion}, vc)

	pr := query.PageRequest{}
	vcs, _, err := k.GetCredentialsByFilter(ctx, sid, parentRegion, &pr)
	require.NoError(t, err)
	require.Empty(t, vcs)

	vcs, _, err = k.GetCredentialsByFilter(ctx, sid, childRegion, &pr)
	require.NoError(t, err)
	require.Len(t, vcs, 1)
	assert.Equal(t, did, vcs[0].Did)
}
