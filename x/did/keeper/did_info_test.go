package keeper_test

import (
	"github.com/openmetaearth/me-hub/testutil/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_DidInfo(t *testing.T) {
	k, ctx := keeper.DidKeeper(t)

	did := "1000000000000001"
	info := didtypes.DidInfo{
		Did:    did,
		Pubkey: "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"Aw4xnG14NvISABc/JSY7NnkOjx1ApE92Xly/KJuZJ7Rm\"}",
		Status: didtypes.DID_STATUS_ACTIVE,
	}

	// set
	k.SetDidInfo(ctx, did, info)
	// has
	assert.True(t, k.HasDidInfo(ctx, did))
	// get
	value, found := k.GetDidInfo(ctx, did)
	assert.True(t, found)
	assert.Equal(t, info, value)
	//gets
	values := k.GetDidInfos(ctx)
	assert.Equal(t, 1, len(values))
	// delete
	k.DeleteDidInfo(ctx, did)
	assert.False(t, k.HasDidInfo(ctx, did))
	value, found = k.GetDidInfo(ctx, did)
	assert.False(t, found)
	assert.Equal(t, didtypes.DidInfo{}, value)
	values = k.GetDidInfos(ctx)
	assert.Equal(t, 0, len(values))

}
