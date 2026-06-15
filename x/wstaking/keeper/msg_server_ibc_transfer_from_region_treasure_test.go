package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/require"
)

func TestNewRegionTreasureIBCTransferMsgUsesDestinationReceiver(t *testing.T) {
	treasureAddress := "me1treasure0000000000000000000000000000006jz0xm"
	receiver := "cosmos1receiver00000000000000000000000000000dr3nrk"

	msg := &types.MsgIbcTransferFromRegionTreasure{
		SourcePort:       "transfer",
		SourceChannel:    "channel-7",
		RegionId:         "USA",
		Token:            sdk.NewCoin("umec", sdk.NewInt(123)),
		Receiver:         receiver,
		TimeoutHeight:    types.Height{RevisionNumber: 1, RevisionHeight: 99},
		TimeoutTimestamp: 123456789,
		Memo:             "region treasury payout",
		Creator:          "me1creator00000000000000000000000000000006u9k9e",
	}

	transferMsg := newRegionTreasureIBCTransferMsg(msg, treasureAddress)

	require.Equal(t, treasureAddress, transferMsg.Sender)
	require.Equal(t, receiver, transferMsg.Receiver)
	require.Equal(t, msg.SourcePort, transferMsg.SourcePort)
	require.Equal(t, msg.SourceChannel, transferMsg.SourceChannel)
	require.True(t, msg.Token.Equal(transferMsg.Token))
	require.Equal(t, ibcclienttypes.Height{RevisionNumber: 1, RevisionHeight: 99}, transferMsg.TimeoutHeight)
	require.Equal(t, msg.TimeoutTimestamp, transferMsg.TimeoutTimestamp)
	require.Equal(t, msg.Memo, transferMsg.Memo)
	require.NotEqual(t, transferMsg.Sender, transferMsg.Receiver)
}
