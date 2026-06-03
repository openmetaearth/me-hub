package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgTransferRegionValidateBasic(t *testing.T) {
	validCreator := sample.AccAddress()
	validAddress := sample.AccAddress()

	tests := []struct {
		name string
		msg  MsgTransferRegion
		err  error
	}{
		{
			name: "valid",
			msg: MsgTransferRegion{
				Creator:    validCreator,
				FromRegion: MeEarthRegionId,
				ToRegion:   "usa",
				Address:    []string{validAddress},
			},
		},
		{
			name: "empty address list",
			msg: MsgTransferRegion{
				Creator:    validCreator,
				FromRegion: MeEarthRegionId,
				ToRegion:   "usa",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid transfer address",
			msg: MsgTransferRegion{
				Creator:    validCreator,
				FromRegion: MeEarthRegionId,
				ToRegion:   "usa",
				Address:    []string{"not-an-address"},
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
