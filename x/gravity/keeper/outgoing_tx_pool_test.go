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

func (s *KeeperTestSuite) TestAddUnbatchedTxBridgeFeeDecrementsBridgeTokenSupply() {
	sender := helpers.GenerateAddress().Bytes()
	receiver := helpers.GenerateAddress().Hex()
	denom := "supplydrift"

	initialSupply := sdk.NewCoin(denom, sdkmath.NewInt(1000))
	s.NewBridgeToken(sender, initialSupply)

	amount := sdk.NewCoin(denom, sdkmath.NewInt(300))
	fee := sdk.NewCoin(denom, sdkmath.NewInt(100))
	txID, err := s.Keeper().AddToOutgoingPool(s.Ctx, sender, receiver, amount, fee)
	s.Require().NoError(err)

	bridgeToken, err := s.Keeper().GetBridgeTokenByDenom(s.Ctx, denom)
	s.Require().NoError(err)
	s.Require().Equal(sdkmath.NewInt(600), bridgeToken.Supply)

	addedFee := sdk.NewCoin(denom, sdkmath.NewInt(200))
	err = s.Keeper().AddUnbatchedTxBridgeFee(s.Ctx, txID, sender, addedFee)
	s.Require().NoError(err)

	bridgeToken, err = s.Keeper().GetBridgeTokenByDenom(s.Ctx, denom)
	s.Require().NoError(err)
	s.Require().Equal(sdkmath.NewInt(400), bridgeToken.Supply)

	tx, err := s.Keeper().GetUnbatchedTxById(s.Ctx, txID)
	s.Require().NoError(err)
	s.Require().Equal(sdkmath.NewInt(300), tx.Fee.Amount)
}
