package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestSendToModuleRejectsDisallowedModuleTargets() {
	sender := sdk.MustAccAddressFromBech32(s.Dao.GlobalDao)
	amount := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1_000))

	for _, moduleName := range []string{
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
	} {
		senderBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, sender, params.BaseDenom)
		moduleAddr := s.App.AccountKeeper.GetModuleAddress(moduleName)
		s.Require().NotNil(moduleAddr)
		moduleBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, moduleAddr, params.BaseDenom)

		_, err := s.msgServer.SendToModule(
			s.Ctx,
			types.NewMsgSendToModule(s.Dao.GlobalDao, moduleName, amount),
		)

		s.Require().ErrorIs(err, sdkerrors.ErrUnauthorized)
		s.Require().ErrorContains(err, "is not an allowed SendToModule target")
		s.Require().Equal(senderBalanceBefore, s.App.BankKeeper.GetBalance(s.Ctx, sender, params.BaseDenom))
		s.Require().Equal(moduleBalanceBefore, s.App.BankKeeper.GetBalance(s.Ctx, moduleAddr, params.BaseDenom))
	}
}

func (s *KeeperTestSuite) TestSendToModuleAllowsApprovedModuleTarget() {
	sender := sdk.MustAccAddressFromBech32(s.Dao.GlobalDao)
	moduleAddr := s.App.AccountKeeper.GetModuleAddress(types.BridgeFeePool)
	s.Require().NotNil(moduleAddr)

	amount := sdk.NewInt64Coin(params.BaseDenom, 1_000)
	senderBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, sender, params.BaseDenom)
	moduleBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, moduleAddr, params.BaseDenom)

	_, err := s.msgServer.SendToModule(
		s.Ctx,
		types.NewMsgSendToModule(s.Dao.GlobalDao, types.BridgeFeePool, sdk.NewCoins(amount)),
	)

	s.Require().NoError(err)
	s.Require().Equal(senderBalanceBefore.Sub(amount), s.App.BankKeeper.GetBalance(s.Ctx, sender, params.BaseDenom))
	s.Require().Equal(moduleBalanceBefore.Add(amount), s.App.BankKeeper.GetBalance(s.Ctx, moduleAddr, params.BaseDenom))
}
