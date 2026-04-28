package keeper_test

import (
	"github.com/openmetaearth/me-hub/testutil/keeper"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_Service(t *testing.T) {
	k, ctx := keeper.DidKeeper(t)

	sid := "test"
	svc := didtypes.Service{
		Sid:         sid,
		Name:        "test",
		Description: "this is a test service.",
		Issuers:     []string{"0000000000001"},
		Status:      didtypes.SERVICE_STATUS_ACTIVE,
	}

	// set
	k.SetService(ctx, sid, svc)
	// get
	value, found := k.GetService(ctx, sid)
	assert.True(t, found)
	assert.Equal(t, svc, value)
	//gets
	values := k.GetServices(ctx)
	assert.Equal(t, 1, len(values))
	// delete
	k.DeleteService(ctx, sid)
	value, found = k.GetService(ctx, sid)
	assert.False(t, found)
	assert.Equal(t, didtypes.Service{}, value)
	values = k.GetServices(ctx)
	assert.Equal(t, 0, len(values))
}
