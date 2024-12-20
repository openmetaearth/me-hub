package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "me-hub/testutil/keeper"
	"me-hub/testutil/nullify"

	"github.com/st-chain/me-hub/x/megroup/types"
)

func TestGroupMemberQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNGroupMember(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetGroupMemberRequest
		response *types.QueryGetGroupMemberResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetGroupMemberRequest{Address: msgs[0].Member.Address},
			response: &types.QueryGetGroupMemberResponse{GroupMember: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetGroupMemberRequest{Address: msgs[0].Member.Address},
			response: &types.QueryGetGroupMemberResponse{GroupMember: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetGroupMemberRequest{Address: msgs[0].Member.Address},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.GroupMember(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestGroupMemberQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNGroupMember(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllGroupMemberRequest {
		return &types.QueryAllGroupMemberRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.GroupMemberAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.GroupMember), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.GroupMember),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.GroupMemberAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.GroupMember), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.GroupMember),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.GroupMemberAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.GroupMember),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.GroupMemberAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
