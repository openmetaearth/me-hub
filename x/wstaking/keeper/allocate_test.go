package keeper_test

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	wbanktypes "github.com/openmetaearth/me-hub/x/wbank/types"
	"github.com/openmetaearth/me-hub/x/wdistri"
	"github.com/openmetaearth/me-hub/x/wmint"
	wminttypes "github.com/openmetaearth/me-hub/x/wmint/types"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestEndBlock() {
	s.SetupTest()

	newRegion := types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            types.MeEarthRegionName,
		OperatorAddress: s.meEarthValidator.OperatorAddress,
	}
	_, err := s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	newRegion = types.MsgNewRegion{
		Creator:         s.Dao.GlobalDao,
		Name:            "USA",
		OperatorAddress: s.usaValidator.OperatorAddress,
	}
	_, err = s.msgServer.NewRegion(s.Ctx, &newRegion)
	s.Require().NoError(err)

	treasuryPoolAcc := s.App.AccountKeeper.GetModuleAccount(s.Ctx, wbanktypes.TreasuryPoolName)
	if treasuryPoolAcc == nil {
		panic(errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", wbanktypes.TreasuryPoolName))
	}

	regionAmount := sdkmath.ZeroInt()
	for i := 0; i < 10; i++ {
		blockNumber := (i + 1) * wminttypes.OneDayTotalBlocks
		s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(int64(blockNumber)).WithChainID(apptesting.TestChainID)

		wmint.BeginBlocker(s.Ctx, s.App.MintKeeper, nil)
		treasuryBalance := s.App.BankKeeper.GetBalance(s.Ctx, treasuryPoolAcc.GetAddress(), params.BaseDenom)
		// s.T().Log("after mint: ", treasuryBalance)

		amount := sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(1)).Mul(treasuryBalance.Amount.ToLegacyDec()).Quo(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(3))).TruncateInt()
		regionAmount = regionAmount.Add(amount)
		wdistri.EndBlock(s.Ctx, *s.App.DistrKeeper)
		treasuryBalance = s.App.BankKeeper.GetBalance(s.Ctx, treasuryPoolAcc.GetAddress(), params.BaseDenom)
		// s.T().Log("after distri: ", treasuryBalance)

		regions := s.App.StakingKeeper.GetAllRegionI(s.Ctx)
		for _, region := range regions {
			balance := s.App.BankKeeper.GetBalance(s.Ctx, sdk.MustAccAddressFromBech32(region.GetRegionTreasureAddr()), params.BaseDenom)
			// s.T().Log(regionAmount.String(), balance.Amount.String())
			s.Require().EqualValues(regionAmount.String(), balance.Amount.String())
		}
	}
}
