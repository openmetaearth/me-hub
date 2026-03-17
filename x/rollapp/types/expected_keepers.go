package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

type IBCClientKeeper interface {
	GetClientState(ctx sdk.Context, clientID string) (exported.ClientState, bool)
	SetClientState(ctx sdk.Context, clientID string, clientState exported.ClientState)
}

type ChannelKeeper interface {
	GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error)
}

type DaoKeeper interface {
	IsGlobalDao(ctx sdk.Context, address string) bool
	IsDao(ctx sdk.Context, address string) bool
	GetGlobalDao(ctx sdk.Context) string
}

type SequencerKeeper interface {
	ProcSequencerByPendingStates(ctx sdk.Context, rollappId string, rollappState *StateInfo) error
	IsExceedAuthoredBlockHeight(ctx sdk.Context, rollappId string, creator string, startHeight uint64, numBlocks uint64) error
}
