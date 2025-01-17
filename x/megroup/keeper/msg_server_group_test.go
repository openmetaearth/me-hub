package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/st-chain/me-hub/x/megroup/types"
)

func TestGroupMsgServerCreate(t *testing.T) {
	srv, ctx := setupMsgServer(t)
	creator := "cosmos1lugrmnrk3ngky85n3hsrxumr3ca7m643h59t72"
	for i := 0; i < 5; i++ {
		resp, err := srv.CreateGroup(ctx, &types.MsgCreateGroup{Creator: creator, GroupInfo: &types.GroupInfo{Admin: "cosmos1lugrmnrk3ngky85n3hsrxumr3ca7m643h59t72"}})
		require.NoError(t, err)
		require.Equal(t, i, int(resp.Id))
	}
}

func TestGroupMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdateGroup
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdateGroup{Creator: creator},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdateGroup{Creator: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdateGroup{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			srv, ctx := setupMsgServer(t)
			_, err := srv.CreateGroup(ctx, &types.MsgCreateGroup{Creator: creator})
			require.NoError(t, err)

			_, err = srv.UpdateGroup(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGroupMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteGroup
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDeleteGroup{Creator: creator},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDeleteGroup{Creator: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgDeleteGroup{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			srv, ctx := setupMsgServer(t)

			_, err := srv.CreateGroup(ctx, &types.MsgCreateGroup{Creator: creator})
			require.NoError(t, err)
			_, err = srv.DeleteGroup(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
