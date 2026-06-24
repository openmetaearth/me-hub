package keeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFreeGasQueriesRejectNilRequests(t *testing.T) {
	k := Keeper{}

	freeGasAccountsResp, err := k.FreeGasAccounts(context.Background(), nil)
	require.Nil(t, freeGasAccountsResp)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	isFreeGasAccountResp, err := k.IsFreeGasAccount(context.Background(), nil)
	require.Nil(t, isFreeGasAccountResp)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
