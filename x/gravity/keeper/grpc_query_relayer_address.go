package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parseQueryRelayerAddress(relayerAddress string) (sdk.AccAddress, error) {
	address, err := sdk.AccAddressFromBech32(relayerAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "relayer address")
	}
	return address, nil
}
