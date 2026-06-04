package transfergenesis

import (
	"errors"
	"testing"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/stretchr/testify/require"
)

func TestShouldEnableTransfersOnRecvAck(t *testing.T) {
	tests := []struct {
		name             string
		ack              exported.Acknowledgement
		transfersEnabled bool
		want             bool
	}{
		{
			name: "pending nil ack does not enable transfers",
			ack:  nil,
			want: false,
		},
		{
			name: "success ack enables transfers",
			ack:  channeltypes.NewResultAcknowledgement([]byte{1}),
			want: true,
		},
		{
			name: "error ack does not enable transfers",
			ack:  channeltypes.NewErrorAcknowledgement(errors.New("recv failed")),
			want: false,
		},
		{
			name:             "already enabled stays unchanged",
			ack:              channeltypes.NewResultAcknowledgement([]byte{1}),
			transfersEnabled: true,
			want:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, shouldEnableTransfersOnRecvAck(tt.ack, tt.transfersEnabled))
		})
	}
}
