package keeper_test

import (
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	mintypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wmintTypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
	"github.com/stretchr/testify/require"
	"strings"
)

func (s *KeeperTestSuite) createGlobalRegion() {
	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	msg := types.MsgNewFixedDepositCfg{
		Dao:      s.Dao.GlobalDao,
		RegionId: strings.ToLower(types.MeEarthRegionName),
		Term:     30,
		Rate:     sdk.MustNewDecFromStr("0.1"),
	}
	_, err = s.msgServer.NewFixedDepositCfg(s.Ctx, &msg)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) createUsaRegion() {
	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)
	msg := types.MsgNewFixedDepositCfg{
		Dao:      s.Dao.GlobalDao,
		RegionId: "usa",
		Term:     30,
		Rate:     sdk.MustNewDecFromStr("0.1"),
	}
	_, err = s.msgServer.NewFixedDepositCfg(s.Ctx, &msg)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) createFixedDeposits(count int, account string) {
	for i := 0; i < count; i++ {
		// Create a sample FixedDeposit
		fixedDeposit := types.MsgDoFixedDeposit{
			Account:   account,
			Principal: sdk.NewCoin(params.BaseDenom, sdk.NewInt(100000000)),
			Term:      30,
		}
		// Run DoFixedDeposit to set FixedDeposit data
		_, err := s.msgServer.DoFixedDeposit(s.Ctx, &fixedDeposit)
		require.NoError(s.T(), err)
	}
}

func (s *KeeperTestSuite) TestFixedDepositByRegionPagination() {
	s.SetupTest()
	s.createGlobalRegion()
	s.createUsaRegion()
	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{}).WithBlockHeight(wmintTypes.OneDayTotalBlocks).WithChainID(apptesting.TestChainID)
	wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
	wdistri.EndBlock(s.Ctx, abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}, *s.App.DistrKeeper)

	accounts := s.NewAccounts(3)
	for _, account := range accounts {
		wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
		err := s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx,
			mintypes.ModuleName,
			account,
			sdk.Coins{sdk.NewInt64Coin(params.BaseDenom, 10000000000)})
		s.Require().NoError(err)
	}
	s.InitKyc(accounts[0], "0000000000001", types.MeEarthRegionId)
	s.createFixedDeposits(10, accounts[0].String())
	s.InitKyc(accounts[1], "0000000000002", "usa")
	s.createFixedDeposits(10, accounts[1].String())
	s.InitKyc(accounts[2], "0000000000003", types.MeEarthRegionId)
	s.createFixedDeposits(10, accounts[2].String())

	var allFixedDeposits []types.FixedDeposit
	var nextKey []byte

	// Query for MeEarthRegionId
	for {
		req := &types.QueryFixedDepositByRegionRequest{
			RegionId: types.MeEarthRegionId,
			Pagination: &query.PageRequest{
				Key:        nextKey,
				Offset:     0,
				Limit:      5,
				CountTotal: false,
				Reverse:    false,
			},
			QueryType: types.FixedDepositState_AllState,
		}

		// Call the FixedDepositByRegion function
		res, err := s.queryClient.FixedDepositByRegion(s.Ctx, req)
		require.NoError(s.T(), err)
		require.NotNil(s.T(), res)

		// Accumulate the results
		allFixedDeposits = append(allFixedDeposits, res.FixedDeposit...)

		// Check if there is a next page
		if res.Pagination.NextKey == nil {
			break
		}
		nextKey = res.Pagination.NextKey
	}

	require.Equal(s.T(), 20, len(allFixedDeposits))

	// Query for USA region
	nextKey = nil
	for {
		req := &types.QueryFixedDepositByRegionRequest{
			RegionId: "usa",
			Pagination: &query.PageRequest{
				Key:        nextKey,
				Offset:     0,
				Limit:      5,
				CountTotal: false,
				Reverse:    false,
			},
			QueryType: types.FixedDepositState_AllState,
		}

		// Call the FixedDepositByRegion function
		res, err := s.queryClient.FixedDepositByRegion(s.Ctx, req)
		require.NoError(s.T(), err)
		require.NotNil(s.T(), res)

		// Accumulate the results
		allFixedDeposits = append(allFixedDeposits, res.FixedDeposit...)

		// Check if there is a next page
		if res.Pagination.NextKey == nil {
			break
		}
		nextKey = res.Pagination.NextKey
	}

	// Check the total number of FixedDeposits
	require.Equal(s.T(), 30, len(allFixedDeposits))
}
