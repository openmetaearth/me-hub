package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "me-hub/testutil/keeper"
	"me-hub/testutil/nullify"

	"github.com/st-chain/me-hub/x/megroup/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestGroupMemberCountQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNGroupMemberCount(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetGroupMemberCountRequest
		response *types.QueryGetGroupMemberCountResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetGroupMemberCountRequest{
				GroupId: msgs[0].GroupId,
			},
			response: &types.QueryGetGroupMemberCountResponse{GroupMemberCount: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetGroupMemberCountRequest{
				GroupId: msgs[1].GroupId,
			},
			response: &types.QueryGetGroupMemberCountResponse{GroupMemberCount: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetGroupMemberCountRequest{
				GroupId: 100000,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.GroupMemberCount(wctx, tc.request)
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

func TestGroupMemberCountQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNGroupMemberCount(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllGroupMemberCountRequest {
		return &types.QueryAllGroupMemberCountRequest{
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
			resp, err := keeper.GroupMemberCountAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.GroupMemberCount), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.GroupMemberCount),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.GroupMemberCountAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.GroupMemberCount), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.GroupMemberCount),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.GroupMemberCountAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.GroupMemberCount),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.GroupMemberCountAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
