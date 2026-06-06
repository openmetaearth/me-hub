package keeper

import (
	"context"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VirtualFrontierBankContractByDenom guards the upstream query from empty input.
func (k Keeper) VirtualFrontierBankContractByDenom(ctx context.Context, request *evmtypes.QueryVirtualFrontierBankContractByDenomRequest) (*evmtypes.QueryVirtualFrontierBankContractByDenomResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if request.MinDenom == "" {
		return nil, status.Error(codes.InvalidArgument, "empty min denom")
	}

	return k.Keeper.VirtualFrontierBankContractByDenom(ctx, request)
}
