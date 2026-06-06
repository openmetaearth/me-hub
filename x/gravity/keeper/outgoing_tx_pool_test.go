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

func (s *KeeperTestSuite) TestKeeper_AddUnbatchedTxBridgeFeeRequiresTxSender() {
	sender := helpers.GenerateAddress().Bytes()
	attacker := helpers.GenerateAddress().Bytes()
	denom := "test"

	sendAmount := sdk.NewCoin(denom, sdkmath.NewInt(1_000))
	bridgeFee := sdk.NewCoin(denom, sdkmath.NewInt(100))
	addBridgeFee := sdk.NewCoin(denom, sdkmath.NewInt(50))

	s.NewBridgeToken(sender, sendAmount.Add(bridgeFee))
	err := s.App.BankKeeper.MintCoins(s.Ctx, s.chainName, sdk.NewCoins(addBridgeFee))
	s.NoError(err)
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, s.chainName, attacker, sdk.NewCoins(addBridgeFee))
	s.NoError(err)

	txID, err := s.Keeper().AddToOutgoingPool(s.Ctx, sender, helpers.GenerateAddress().Hex(), sendAmount, bridgeFee)
	s.NoError(err)
	txBefore, err := s.Keeper().GetUnbatchedTxById(s.Ctx, txID)
	s.NoError(err)

	err = s.Keeper().AddUnbatchedTxBridgeFee(s.Ctx, txID, attacker, addBridgeFee)
	s.Error(err)

	txAfter, err := s.Keeper().GetUnbatchedTxById(s.Ctx, txID)
	s.NoError(err)
	s.Equal(txBefore.Sender, txAfter.Sender)
	s.Equal(txBefore.Fee, txAfter.Fee)
	s.Equal(addBridgeFee.Amount.String(), s.App.BankKeeper.GetAllBalances(s.Ctx, attacker).AmountOf(denom).String())
}
