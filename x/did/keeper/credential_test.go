package keeper_test

import (
	"github.com/openmetaearth/me-hub/testutil/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_Credential(t *testing.T) {
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

	// set
	k.SetCredential(ctx, did, sid, vc)
	// has
	assert.True(t, k.HasCredential(ctx, did, sid))
	// get
	value, found := k.GetCredential(ctx, did, sid)
	assert.True(t, found)
	assert.Equal(t, vc, value)
	//gets
	values := k.GetCredentials(ctx)
	assert.Equal(t, 1, len(values))
	// gets by did
	values = k.GetCredentialsByDid(ctx, did)
	assert.Equal(t, 1, len(values))
	// delete
	k.DeleteCredential(ctx, did, sid)
	assert.False(t, k.HasDidInfo(ctx, did))
	value, found = k.GetCredential(ctx, did, sid)
	assert.False(t, found)
	assert.Equal(t, didtypes.Credential{}, value)
	values = k.GetCredentials(ctx)
	assert.Equal(t, 0, len(values))
	values = k.GetCredentialsByDid(ctx, did)
	assert.Equal(t, 0, len(values))
}
