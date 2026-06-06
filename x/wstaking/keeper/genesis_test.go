package keeper_test

import (
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestInitGenesisRestoresExportedRuntimeState() {
	app, _ := apptesting.SetupTestingApp()
	ctx := app.GetBaseApp().NewContext(false, tmproto.Header{}).WithBlockTime(time.Unix(1_700_000_000, 0).UTC())

	app.AccountKeeper.SetModuleAccount(ctx, authtypes.NewEmptyModuleAccount(types.BondedStakePoolName, authtypes.Burner, authtypes.Staking))
	app.AccountKeeper.SetModuleAccount(ctx, authtypes.NewEmptyModuleAccount(types.NotBondedStakePoolName, authtypes.Burner, authtypes.Staking))

	staker := sdk.AccAddress([]byte("staker_address_12345"))
	validator := sdk.ValAddress([]byte("validator_address123"))
	completionTime := ctx.BlockTime().Add(24 * time.Hour)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = params.BaseDenom
	ubs := types.NewUnbondingStake(staker, validator, 12, completionTime, sdk.NewInt(500))
	region := types.Region{
		RegionId:            "moon",
		Name:                "Moon",
		Creator:             staker.String(),
		OperatorAddress:     validator.String(),
		RegionTreasureAddr:  staker.String(),
		DepositInterestAddr: staker.String(),
		RegionShare:         sdk.NewInt(77),
		DelegateInterest:    sdk.MustNewDecFromStr("0.25"),
		DelegateAmount:      sdk.NewInt(33),
		FixedDepositAmount:  sdk.NewInt(44),
	}
	fixedDeposit := types.FixedDeposit{
		Id:        7,
		Account:   staker.String(),
		Principal: sdk.NewCoin(stakingParams.BondDenom, sdk.NewInt(1000)),
		Interest:  sdk.NewCoin(stakingParams.BondDenom, sdk.NewInt(50)),
		StartTime: ctx.BlockTime(),
		EndTime:   ctx.BlockTime().Add(30 * 24 * time.Hour),
		Term:      30,
		Rate:      sdk.MustNewDecFromStr("0.05"),
	}

	genesis := types.DefaultGenesisState()
	genesis.Params = stakingParams
	genesis.UnbondingStakes = []types.UnbondingStake{ubs}
	genesis.Regions = []types.Region{region}
	genesis.FixedDepositList = []types.FixedDeposit{fixedDeposit}
	genesis.FixedDepositCount = 8
	genesis.Exported = true

	updates := app.StakingKeeper.InitGenesis(ctx, genesis)
	s.Require().Empty(updates)

	restoredUBS, found := app.StakingKeeper.GetUnbondingStake(ctx, staker, validator)
	s.Require().True(found)
	s.Require().Equal(ubs, restoredUBS)

	queue := app.StakingKeeper.GetUBSQueueTimeSlice(ctx, completionTime)
	s.Require().Len(queue, 1)
	s.Require().Equal(ubs.StakerAddress, queue[0].StakerAddress)
	s.Require().Equal(ubs.ValidatorAddress, queue[0].ValidatorAddress)

	restoredRegion, found := app.StakingKeeper.GetRegion(ctx, region.RegionId)
	s.Require().True(found)
	s.Require().Equal(region, restoredRegion)

	restoredFixedDeposit, found := app.StakingKeeper.GetFixedDeposit(ctx, fixedDeposit.Id)
	s.Require().True(found)
	s.Require().Equal(fixedDeposit, restoredFixedDeposit)

	accountFixedDeposits, err := app.StakingKeeper.GetFixedDepositByAcct(ctx, staker.String())
	s.Require().NoError(err)
	s.Require().Len(accountFixedDeposits, 1)
	s.Require().Equal(fixedDeposit, accountFixedDeposits[0])
	s.Require().Equal(uint64(8), app.StakingKeeper.GetFixedDepositCount(ctx))
}
