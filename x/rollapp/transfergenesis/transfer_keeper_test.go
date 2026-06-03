package transfergenesis

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/utils/gerrc"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/stretchr/testify/require"
)

func TestTransferEnabledTransferKeeperRejectsDisabledRollappTransfer(t *testing.T) {
	const rollappID = "rollapp_100-1"
	next := &recordingPacketForwardTransferKeeper{}
	keeper := NewTransferEnabledTransferKeeper(
		next,
		staticRollappGetter(rollappID, false),
		staticChannelKeeper{chainID: rollappID},
	)

	_, err := keeper.Transfer(sdk.WrapSDKContext(sdk.Context{}), testMsgTransfer())

	require.ErrorIs(t, err, gerrc.ErrFailedPrecondition)
	require.False(t, next.called)
}

func TestTransferEnabledTransferKeeperAllowsEnabledRollappTransfer(t *testing.T) {
	const rollappID = "rollapp_100-1"
	next := &recordingPacketForwardTransferKeeper{}
	keeper := NewTransferEnabledTransferKeeper(
		next,
		staticRollappGetter(rollappID, true),
		staticChannelKeeper{chainID: rollappID},
	)

	_, err := keeper.Transfer(sdk.WrapSDKContext(sdk.Context{}), testMsgTransfer())

	require.NoError(t, err)
	require.True(t, next.called)
}

func testMsgTransfer() *types.MsgTransfer {
	return types.NewMsgTransfer(
		types.PortID,
		"channel-0",
		sdk.NewCoin(params.BaseDenom, sdk.NewInt(1)),
		"me1l2lt9e4w6qj0ra6jq3dn4csx3z2tyjwk6zu5ac",
		"rollapp1l2lt9e4w6qj0ra6jq3dn4csx3z2tyjwk7t4za7",
		clienttypes.Height{},
		0,
		"",
	)
}

func staticRollappGetter(rollappID string, transfersEnabled bool) GetRollapp {
	return func(ctx sdk.Context, id string) (rollapptypes.Rollapp, bool) {
		if id != rollappID {
			return rollapptypes.Rollapp{}, false
		}
		return rollapptypes.Rollapp{
			RollappId: rollappID,
			GenesisState: rollapptypes.RollappGenesisState{
				TransfersEnabled: transfersEnabled,
			},
		}, true
	}
}

type staticChannelKeeper struct {
	chainID string
}

func (k staticChannelKeeper) GetChannelClientState(
	ctx sdk.Context,
	portID string,
	channelID string,
) (string, exported.ClientState, error) {
	return "07-tendermint-0", &ibctmtypes.ClientState{ChainId: k.chainID}, nil
}

type recordingPacketForwardTransferKeeper struct {
	called bool
}

func (k *recordingPacketForwardTransferKeeper) Transfer(
	ctx context.Context,
	msg *types.MsgTransfer,
) (*types.MsgTransferResponse, error) {
	k.called = true
	return &types.MsgTransferResponse{Sequence: 1}, nil
}

func (k *recordingPacketForwardTransferKeeper) DenomPathFromHash(ctx sdk.Context, denom string) (string, error) {
	return denom, nil
}

func (k *recordingPacketForwardTransferKeeper) GetTotalEscrowForDenom(ctx sdk.Context, denom string) sdk.Coin {
	return sdk.NewCoin(denom, sdk.ZeroInt())
}

func (k *recordingPacketForwardTransferKeeper) SetTotalEscrowForDenom(ctx sdk.Context, coin sdk.Coin) {
}
