package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/testutil/keeper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_DID(t *testing.T) {
	k, ctx := keeper.DidKeeper(t)
	addr := sdk.MustAccAddressFromBech32("me1kjnt3ypezt3yf58w8upujvejdtt5xsvkq5dpk4")
	did := "1000000000000001"
	// set
	k.SetDID(ctx, addr, did)
	// has
	assert.True(t, k.HasDID(ctx, addr))
	// get
	value, found := k.GetDID(ctx, addr)
	assert.True(t, found)
	assert.Equal(t, did, value)
	// delete
	k.DeleteDID(ctx, addr)
	assert.False(t, k.HasDID(ctx, addr))
	value, found = k.GetDID(ctx, addr)
	assert.False(t, found)
	assert.Equal(t, "", value)
}
