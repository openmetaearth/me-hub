package keeper

import (
	"context"
	"testing"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestVirtualFrontierBankContractByDenomRejectsEmptyRequest(t *testing.T) {
	var k Keeper

	_, err := k.VirtualFrontierBankContractByDenom(context.Background(), nil)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestVirtualFrontierBankContractByDenomRejectsEmptyMinDenom(t *testing.T) {
	var k Keeper

	_, err := k.VirtualFrontierBankContractByDenom(context.Background(), &evmtypes.QueryVirtualFrontierBankContractByDenomRequest{})
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
