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

func TestMemberJoinedQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNMemberJoined(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetMemberJoinedRequest
		response *types.QueryGetMemberJoinedResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetMemberJoinedRequest{
				Address: msgs[0].Address,
			},
			response: &types.QueryGetMemberJoinedResponse{MemberJoined: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetMemberJoinedRequest{
				Address: msgs[1].Address,
			},
			response: &types.QueryGetMemberJoinedResponse{MemberJoined: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetMemberJoinedRequest{
				Address: strconv.Itoa(100000),
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
			response, err := keeper.MemberJoined(wctx, tc.request)
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

func TestMemberJoinedQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.MegroupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNMemberJoined(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllMemberJoinedRequest {
		return &types.QueryAllMemberJoinedRequest{
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
			resp, err := keeper.MemberJoinedAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.MemberJoined), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.MemberJoined),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.MemberJoinedAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.MemberJoined), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.MemberJoined),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.MemberJoinedAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.MemberJoined),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.MemberJoinedAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
