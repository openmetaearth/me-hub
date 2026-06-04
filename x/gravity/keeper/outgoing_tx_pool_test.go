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

func (s *KeeperTestSuite) TestAddUnbatchedTxBridgeFeeUsesExternalUnits() {
	sender := helpers.GenerateAddress().Bytes()
	receiver := helpers.GenerateAddress().Hex()
	denom := "uusdt"
	tokenContract := helpers.GenerateAddress().Hex()
	bridgeToken := types.BridgeToken{
		ContractAddress: tokenContract,
		Denom:           denom,
		Symbol:          "USDT",
		Decimal:         18,
		Supply:          sdkmath.NewInt(20_000_000),
	}
	initialBalance := sdk.NewCoin(denom, bridgeToken.Supply)

	s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, s.chainName, sdk.NewCoins(initialBalance)))
	s.Require().NoError(s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, s.chainName, sender, sdk.NewCoins(initialBalance)))
	s.Keeper().SetBridgeToken(s.Ctx, &bridgeToken)

	amount := sdk.NewCoin(denom, sdkmath.NewInt(5_000_000))
	initialFee := sdk.NewCoin(denom, sdkmath.NewInt(1_000_000))
	txID, err := s.Keeper().AddToOutgoingPool(s.Ctx, sender, receiver, amount, initialFee)
	s.Require().NoError(err)

	tx, err := s.Keeper().GetUnbatchedTxById(s.Ctx, txID)
	s.Require().NoError(err)
	s.Require().Equal(sdkmath.NewIntWithDecimal(1, 18).String(), tx.Fee.Amount.String())

	addBridgeFee := sdk.NewCoin(denom, sdkmath.NewInt(2_000_000))
	s.Require().NoError(s.Keeper().AddUnbatchedTxBridgeFee(s.Ctx, txID, sender, addBridgeFee))

	tx, err = s.Keeper().GetUnbatchedTxById(s.Ctx, txID)
	s.Require().NoError(err)
	s.Require().Equal(sdkmath.NewIntWithDecimal(3, 18).String(), tx.Fee.Amount.String())
}
