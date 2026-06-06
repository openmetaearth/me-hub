package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

func (s *KeeperTestSuite) TestKeeper_OutgoingAncCancel() {
	sender := helpers.GenerateAddress().Bytes()
	bridgeToken := helpers.GenerateAddress().Hex()

	denom := "test"
	s.Equal(sdk.NewCoin(denom, sdkmath.ZeroInt()), s.App.BankKeeper.GetSupply(s.Ctx, denom))

	sendAmount := sdk.NewCoin(denom, sdkmath.NewInt(int64(tmrand.Uint32()*2)))
	err := s.App.BankKeeper.MintCoins(s.Ctx, s.chainName, sdk.NewCoins(sendAmount))
	s.NoError(err)
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, s.chainName, sender, sdk.NewCoins(sendAmount))
	s.NoError(err)
	s.Equal(sendAmount, s.App.BankKeeper.GetSupply(s.Ctx, denom))

	s.Keeper().SetBridgeToken(s.Ctx, &types.BridgeToken{ContractAddress: bridgeToken, Denom: denom, Supply: sendAmount.Amount})
	s.Equal(s.App.BankKeeper.GetAllBalances(s.Ctx, sender).AmountOf(denom).String(), sendAmount.Amount.String())

	receiver := helpers.GenerateAddress().Hex()
	amount := sdk.NewCoin(denom, sendAmount.Amount.QuoRaw(2))
	txId, err := s.Keeper().AddToOutgoingPool(s.Ctx, sender, receiver, amount, amount)
	s.NoError(err)
	s.Equal(s.App.BankKeeper.GetAllBalances(s.Ctx, sender).AmountOf(denom).String(), sdkmath.NewInt(0).String())

	s.Equal(sdk.NewCoin(denom, sdkmath.ZeroInt()), s.App.BankKeeper.GetSupply(s.Ctx, denom))

	_, err = s.Keeper().RemoveFromOutgoingPoolAndRefund(s.Ctx, txId, sender)
	s.NoError(err)
	s.Equal(s.App.BankKeeper.GetAllBalances(s.Ctx, sender).AmountOf(denom).String(), sendAmount.Amount.String())
	s.Equal(sendAmount, s.App.BankKeeper.GetSupply(s.Ctx, denom))
}

func (s *KeeperTestSuite) TestAddToOutgoingPoolAllowsMultipleWithdrawalsAfterBurnedPendingSupply() {
	sender := helpers.GenerateAddress().Bytes()
	denom := "testpending"
	initialAmount := sdk.NewCoin(denom, sdkmath.NewInt(100))
	bridgeToken := s.NewBridgeToken(sender, initialAmount)
	receiver := helpers.GenerateAddress().Hex()

	_, err := s.Keeper().AddToOutgoingPool(
		s.Ctx,
		sender,
		receiver,
		sdk.NewCoin(denom, sdkmath.NewInt(30)),
		sdk.NewCoin(denom, sdkmath.NewInt(10)),
	)
	s.Require().NoError(err)

	tokenAfterFirst, err := s.Keeper().GetBridgeTokenByContract(s.Ctx, bridgeToken.ContractAddress)
	s.Require().NoError(err)
	s.Require().Equal(sdkmath.NewInt(60), tokenAfterFirst.Supply)
	s.Require().Equal(sdkmath.NewInt(60), s.App.BankKeeper.GetBalance(s.Ctx, sender, denom).Amount)

	_, err = s.Keeper().AddToOutgoingPool(
		s.Ctx,
		sender,
		receiver,
		sdk.NewCoin(denom, sdkmath.NewInt(20)),
		sdk.NewCoin(denom, sdkmath.NewInt(5)),
	)
	s.Require().NoError(err)

	tokenAfterSecond, err := s.Keeper().GetBridgeTokenByContract(s.Ctx, bridgeToken.ContractAddress)
	s.Require().NoError(err)
	s.Require().Equal(sdkmath.NewInt(35), tokenAfterSecond.Supply)
	s.Require().Equal(sdkmath.NewInt(35), s.App.BankKeeper.GetBalance(s.Ctx, sender, denom).Amount)
	s.Require().Len(s.Keeper().GetUnbatchedTransactions(s.Ctx), 2)
}
