package types

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"

	rollupTypes "github.com/st-chain/me-hub/x/rollup/types"
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
	Punishment(ctx sdk.Context, address, rollappID string, rate uint32, amount uint64) (uint64, error)
	RevaluateSequencer(ctx sdk.Context, address, rollappID string) error
	RegisterRollappInitInfo(ctx sdk.Context, rollappID string, FirstElectBlkHeight uint64, IdInDa []byte) error
	IsInBlackList(addr string) bool
	StakeForChallengeDaFraud(goCtx context.Context, rollappID, blockSubmitter, challenger string, challengeKey []byte) error
	ProcChallengeDaFraud(goCtx context.Context, rollappID string, challengeKey []byte, result int32) error
	GetBondNodeDelegator(ctx sdk.Context, rollappID string, bondAddress []byte) []byte
}

type DaoKeeper interface {
	IsGlobalDao(ctx sdk.Context, address string) bool
	IsValidatorDao(ctx sdk.Context, address string) bool
}
