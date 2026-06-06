package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgIbcTransferFromRegionTreasure_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgIbcTransferFromRegionTreasure
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgIbcTransferFromRegionTreasure{
				SourcePort:    "transfer",
				SourceChannel: "channel-0",
				RegionId:      MeEarthRegionName,
				Token:         sdk.NewCoin(params.BaseDenom, sdk.NewInt(1)),
				Receiver:      "cosmos1destination",
				Creator:       "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "empty receiver",
			msg: MsgIbcTransferFromRegionTreasure{
				SourcePort:    "transfer",
				SourceChannel: "channel-0",
				RegionId:      MeEarthRegionName,
				Token:         sdk.NewCoin(params.BaseDenom, sdk.NewInt(1)),
				Creator:       sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid message",
			msg: MsgIbcTransferFromRegionTreasure{
				SourcePort:    "transfer",
				SourceChannel: "channel-0",
				RegionId:      MeEarthRegionName,
				Token:         sdk.NewCoin(params.BaseDenom, sdk.NewInt(1)),
				Receiver:      "cosmos1destination",
				Creator:       sample.AccAddress(),
			},
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

func TestMsgIbcTransferFromRegionTreasureMarshalPreservesReceiver(t *testing.T) {
	msg := NewMsgIbcTransferFromRegionTreasure(
		"transfer",
		"channel-0",
		MeEarthRegionName,
		"cosmos1destination",
		sdk.NewCoin(params.BaseDenom, sdk.NewInt(1)),
		Height{},
		0,
		"",
		sample.AccAddress(),
	)

	bz, err := msg.Marshal()
	require.NoError(t, err)

	var decoded MsgIbcTransferFromRegionTreasure
	require.NoError(t, decoded.Unmarshal(bz))
	require.Equal(t, msg.Receiver, decoded.Receiver)
}
