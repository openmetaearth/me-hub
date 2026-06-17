package delayedack_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/openmetaearth/me-hub/x/common/types"
	delayedack "github.com/openmetaearth/me-hub/x/delayedack"
	delayedacktypes "github.com/openmetaearth/me-hub/x/delayedack/types"
)

type proofHeightTx struct {
	msgs []sdk.Msg
}

func (tx proofHeightTx) GetMsgs() []sdk.Msg {
	return tx.msgs
}

func (tx proofHeightTx) ValidateBasic() error {
	return nil
}

func TestIBCProofHeightDecoratorRecordsTimeoutOnClose(t *testing.T) {
	decorator := delayedack.NewIBCProofHeightDecorator()
	ctx := sdk.NewContext(nil, tmproto.Header{}, false, nil)
	proofHeight := clienttypes.NewHeight(7, 99)
	packet := channeltypes.Packet{
		Sequence:      1,
		SourcePort:    "transfer",
		SourceChannel: "channel-0",
	}
	packetID := types.NewPacketUID(
		types.RollappPacket_ON_TIMEOUT,
		packet.SourcePort,
		packet.SourceChannel,
		packet.Sequence,
	)

	tx := proofHeightTx{
		msgs: []sdk.Msg{
			&channeltypes.MsgTimeoutOnClose{
				Packet:      packet,
				ProofHeight: proofHeight,
			},
		},
	}

	_, err := decorator.AnteHandle(ctx, tx, false, func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		got, ok := delayedacktypes.PacketProofHeightFromCtx(ctx, packetID)
		require.True(t, ok)
		require.Equal(t, proofHeight, got)
		return ctx, nil
	})

	require.NoError(t, err)
}
