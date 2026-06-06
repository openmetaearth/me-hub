package keeper_test

import (
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

func TestKeeper_FilterPrefixCollisionDoesNotReturnWrongCredential(t *testing.T) {
	k, ctx := keeper.DidKeeper(t)

	collidingVC := didtypes.Credential{
		Did:  "did-b",
		Sid:  "abc",
		Hash: "FF000000000000000001",
		Uri:  "https://www.example.com/b",
		Data: []byte("credential b"),
	}

	k.AddFilters(ctx, collidingVC.Did, collidingVC.Sid, [][]byte{[]byte("X")}, collidingVC)

	pr := query.PageRequest{}
	vcs, _, err := k.GetCredentialsByFilter(ctx, "ab", []byte("cX"), &pr)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(vcs))

	vcs, _, err = k.GetCredentialsByFilter(ctx, "abc", []byte("X"), &pr)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(vcs))
	assert.Equal(t, collidingVC, vcs[0])
}
