package types

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"

	rollupTypes "github.com/dymensionxyz/dymension/v3/x/rollup/types"
)

type IBCClientKeeper interface {
	GetClientState(ctx sdk.Context, clientID string) (exported.ClientState, bool)
	SetClientState(ctx sdk.Context, clientID string, clientState exported.ClientState)
}

type ChannelKeeper interface {
	GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error)
}

type RollupKeeper interface {
	QueryElectionResult(ctx context.Context, req *rollupTypes.QueryElectionRequest) (*rollupTypes.QueryElectionResponse, error)
	GetPreviousElectionResult(ctx context.Context, rollappID string) (*rollupTypes.QueryElectionResponse, error)
	GetParams(ctx sdk.Context) rollupTypes.Params
	QueryStake(ctx context.Context, req *rollupTypes.QueryStakeRequest) (*rollupTypes.QueryStakeResponse, error)
	Punishment(ctx sdk.Context, address, rollappID string, rate uint32, amount uint64) error
}
