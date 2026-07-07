package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateGroup_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateGroup
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateGroup{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "missing GroupInfo",
			msg: MsgCreateGroup{
				Creator: sample.AccAddress(),
				// GroupInfo is nil -> should fail with ErrInvalidAddress on Admin
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "empty RegionID",
			msg: MsgCreateGroup{
				Creator: sample.AccAddress(),
				GroupInfo: &GroupInfo{
					Admin:    sample.AccAddress(),
					RegionID: "",
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid",
			msg: MsgCreateGroup{
				Creator: sample.AccAddress(),
				GroupInfo: &GroupInfo{
					Admin:    sample.AccAddress(),
					RegionID: "ME_EARTH",
				},
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

func TestMsgUpdateGroup_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateGroup
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateGroup{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateGroup{
				Creator: sample.AccAddress(),
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

func TestMsgDeleteGroup_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteGroup
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteGroup{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteGroup{
				Creator: sample.AccAddress(),
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
