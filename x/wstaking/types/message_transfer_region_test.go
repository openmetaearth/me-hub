package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgTransferRegion_ValidateBasic(t *testing.T) {
	validAddress := sample.AccAddress()

	tests := []struct {
		name string
		msg  MsgTransferRegion
		err  error
	}{
		{
			name: "valid message",
			msg: MsgTransferRegion{
				Creator:    sample.AccAddress(),
				FromRegion: "me_earth",
				ToRegion:   "usa",
				Address:    []string{validAddress},
			},
		},
		{
			name: "invalid creator",
			msg: MsgTransferRegion{
				Creator:    "invalid_address",
				FromRegion: "me_earth",
				ToRegion:   "usa",
				Address:    []string{validAddress},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "same regions",
			msg: MsgTransferRegion{
				Creator:    sample.AccAddress(),
				FromRegion: "usa",
				ToRegion:   "USA",
				Address:    []string{validAddress},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "missing address",
			msg: MsgTransferRegion{
				Creator:    sample.AccAddress(),
				FromRegion: "me_earth",
				ToRegion:   "usa",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid transfer address",
			msg: MsgTransferRegion{
				Creator:    sample.AccAddress(),
				FromRegion: "me_earth",
				ToRegion:   "usa",
				Address:    []string{"invalid_address"},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
