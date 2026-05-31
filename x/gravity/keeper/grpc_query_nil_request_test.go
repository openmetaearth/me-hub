package keeper_test

import (
	"context"
	"testing"

	"github.com/openmetaearth/me-hub/x/gravity/keeper"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRelayerSetRequestRejectsNilRequest(t *testing.T) {
	_, err := (keeper.QueryServer{}).RelayerSetRequest(context.Background(), nil)

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
