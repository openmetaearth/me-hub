package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openmetaearth/me-hub/testutil/nullify"
	"github.com/openmetaearth/me-hub/x/sequencer/keeper"
	"github.com/openmetaearth/me-hub/x/sequencer/types"
)

func (suite *SequencerTestSuite) TestSequencersByRollappQuery3() {
	suite.SetupTest()

	rollappId := suite.CreateDefaultRollapp()
	rollappId2 := suite.CreateDefaultRollapp()

	// create 2 sequencer
	suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	seq1Response := types.QueryGetSequencersByRollappResponse{
		Sequencers: suite.App.SequencerKeeper.GetSequencersByRollapp(suite.Ctx, rollappId),
	}

	suite.CreateDefaultSequencer(suite.Ctx, rollappId2)
	suite.CreateDefaultSequencer(suite.Ctx, rollappId2)
	seq2Response := types.QueryGetSequencersByRollappResponse{
		Sequencers: suite.App.SequencerKeeper.GetSequencersByRollapp(suite.Ctx, rollappId2),
	}

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetSequencersByRollappRequest
		response *types.QueryGetSequencersByRollappResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetSequencersByRollappRequest{
				RollappId: rollappId,
			},
			response: &seq1Response,
		},
		{
			desc: "Second",
			request: &types.QueryGetSequencersByRollappRequest{
				RollappId: rollappId2,
			},
			response: &seq2Response,
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetSequencersByRollappRequest{
				RollappId: strconv.Itoa(100000),
			},
			err: types.ErrUnknownRollappID,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		suite.T().Run(tc.desc, func(t *testing.T) {
			response, err := suite.App.SequencerKeeper.SequencersByRollapp(suite.Ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response.Sequencers),
					nullify.Fill(response.Sequencers),
				)
				require.NotNil(t, response.Pagination)
				require.Equal(t, uint64(len(tc.response.Sequencers)), response.Pagination.Total)
			}
		})
	}
}

func (suite *SequencerTestSuite) TestSequencersByRollappQueryPagination() {
	suite.SetupTest()

	rollappId := suite.CreateDefaultRollapp()
	suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	suite.CreateDefaultSequencer(suite.Ctx, rollappId)

	response, err := suite.App.SequencerKeeper.SequencersByRollapp(suite.Ctx, &types.QueryGetSequencersByRollappRequest{
		RollappId: rollappId,
		Pagination: &query.PageRequest{
			Limit:      1,
			CountTotal: true,
		},
	})
	require.NoError(suite.T(), err)
	require.Len(suite.T(), response.Sequencers, 1)
	require.NotNil(suite.T(), response.Pagination)
	require.Equal(suite.T(), uint64(3), response.Pagination.Total)
	require.NotEmpty(suite.T(), response.Pagination.NextKey)
}

func (suite *SequencerTestSuite) TestSequencersByRollappByStatusQuery() {
	suite.SetupTest()

	msgserver := keeper.NewMsgServerImpl(suite.App.SequencerKeeper)

	rollappId := suite.CreateDefaultRollapp()
	// create 2 sequencers on rollapp1
	addr1_1 := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	addr2_1 := suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	_, err := msgserver.Unbond(suite.Ctx, &types.MsgUnbond{
		Creator: addr2_1,
	})
	require.NoError(suite.T(), err)

	// create 2 sequencers on rollapp2
	rollappId2 := suite.CreateDefaultRollapp()
	addr1_2 := suite.CreateDefaultSequencer(suite.Ctx, rollappId2)
	addr2_2 := suite.CreateDefaultSequencer(suite.Ctx, rollappId2)

	for _, tc := range []struct {
		desc          string
		request       *types.QueryGetSequencersByRollappByStatusRequest
		response_addr []string
		err           error
	}{
		{
			desc: "First - Bonded",
			request: &types.QueryGetSequencersByRollappByStatusRequest{
				RollappId: rollappId,
				Status:    types.Bonded,
			},
			response_addr: []string{addr1_1},
		},
		{
			desc: "First - Unbonding",
			request: &types.QueryGetSequencersByRollappByStatusRequest{
				RollappId: rollappId,
				Status:    types.Unbonding,
			},
			response_addr: []string{addr2_1},
		},
		{
			desc: "First - Unbonded",
			request: &types.QueryGetSequencersByRollappByStatusRequest{
				RollappId: rollappId,
				Status:    types.Unbonded,
			},
			response_addr: []string{},
		},
		{
			desc: "Second",
			request: &types.QueryGetSequencersByRollappByStatusRequest{
				RollappId: rollappId2,
				Status:    types.Bonded,
			},
			response_addr: []string{addr1_2, addr2_2},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetSequencersByRollappByStatusRequest{
				RollappId: strconv.Itoa(100000),
			},
			err: types.ErrUnknownRollappID,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		suite.T().Run(tc.desc, func(t *testing.T) {
			response, err := suite.App.SequencerKeeper.SequencersByRollappByStatus(suite.Ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Len(t, response.Sequencers, len(tc.response_addr))

				for _, seqAddr := range tc.response_addr {
					seq, found := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, seqAddr)
					require.True(t, found)
					require.Contains(t, response.Sequencers, seq)
				}
			}
		})
	}
}

func (suite *SequencerTestSuite) TestSequencersByRollappByStatusQueryPagination() {
	suite.SetupTest()

	rollappId := suite.CreateDefaultRollapp()
	suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	suite.CreateDefaultSequencer(suite.Ctx, rollappId)
	suite.CreateDefaultSequencer(suite.Ctx, rollappId)

	response, err := suite.App.SequencerKeeper.SequencersByRollappByStatus(suite.Ctx, &types.QueryGetSequencersByRollappByStatusRequest{
		RollappId: rollappId,
		Status:    types.Bonded,
		Pagination: &query.PageRequest{
			Limit:      2,
			CountTotal: true,
		},
	})
	require.NoError(suite.T(), err)
	require.Len(suite.T(), response.Sequencers, 2)
	require.NotNil(suite.T(), response.Pagination)
	require.Equal(suite.T(), uint64(3), response.Pagination.Total)
	require.NotEmpty(suite.T(), response.Pagination.NextKey)
}
