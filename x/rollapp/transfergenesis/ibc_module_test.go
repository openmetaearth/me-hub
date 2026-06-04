package transfergenesis

import (
	"errors"
	"testing"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

func TestRecvPacketSucceeded(t *testing.T) {
	require.True(t, recvPacketSucceeded(nil))
	require.True(t, recvPacketSucceeded(channeltypes.NewResultAcknowledgement([]byte("ok"))))
	require.False(t, recvPacketSucceeded(channeltypes.NewErrorAcknowledgement(errors.New("failed"))))
}
