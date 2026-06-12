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

func (s *KeeperTestSuite) TestKeeper_IncreaseBridgeFeeRejectsNonSender() {
	sender := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	attacker := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	bridgeToken := helpers.GenerateAddress().Hex()

	denom := "testfee"
	initialSupply := sdk.NewCoin(denom, sdkmath.NewInt(150))
	s.Equal(sdk.NewCoin(denom, sdkmath.ZeroInt()), s.App.BankKeeper.GetSupply(s.Ctx, denom))

	err := s.App.BankKeeper.MintCoins(s.Ctx, s.chainName, sdk.NewCoins(initialSupply))
	s.NoError(err)
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, s.chainName, sender, sdk.NewCoins(sdk.NewCoin(denom, sdkmath.NewInt(100))))
	s.NoError(err)
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, s.chainName, attacker, sdk.NewCoins(sdk.NewCoin(denom, sdkmath.NewInt(50))))
	s.NoError(err)

	s.Keeper().SetBridgeToken(s.Ctx, &types.BridgeToken{ContractAddress: bridgeToken, Denom: denom, Supply: initialSupply.Amount})

	receiver := helpers.GenerateAddress().Hex()
	amount := sdk.NewCoin(denom, sdkmath.NewInt(40))
	fee := sdk.NewCoin(denom, sdkmath.NewInt(10))
	txId, err := s.Keeper().AddToOutgoingPool(s.Ctx, sender, receiver, amount, fee)
	s.NoError(err)

	tx, err := s.Keeper().GetUnbatchedTxById(s.Ctx, txId)
	s.Require().NoError(err)
	s.Require().EqualValues(fee.Amount, tx.Fee.Amount)

	attackerBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, attacker, denom)
	senderBalanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, sender, denom)

	_, err = s.MsgServer().IncreaseBridgeFee(sdk.WrapSDKContext(s.Ctx), &types.MsgIncreaseBridgeFee{
		ChainName:     s.chainName,
		Sender:        attacker.String(),
		TransactionId: txId,
		AddBridgeFee:  sdk.NewCoin(denom, sdkmath.NewInt(5)),
	})
	s.Require().Error(err)
	s.Require().ErrorIs(err, types.ErrInvalid)
	s.Require().Contains(err.Error(), "is not tx sender")
	s.Require().Contains(err.Error(), sender.String())
	s.Require().Contains(err.Error(), attacker.String())

	s.Require().Equal(attackerBalanceBefore, s.App.BankKeeper.GetBalance(s.Ctx, attacker, denom))
	s.Require().Equal(senderBalanceBefore, s.App.BankKeeper.GetBalance(s.Ctx, sender, denom))

	_, err = s.MsgServer().IncreaseBridgeFee(sdk.WrapSDKContext(s.Ctx), &types.MsgIncreaseBridgeFee{
		ChainName:     s.chainName,
		Sender:        sender.String(),
		TransactionId: txId,
		AddBridgeFee:  sdk.NewCoin(denom, sdkmath.NewInt(5)),
	})
	s.Require().NoError(err)

	tx, err = s.Keeper().GetUnbatchedTxById(s.Ctx, txId)
	s.Require().NoError(err)
	s.Require().EqualValues(sdkmath.NewInt(15), tx.Fee.Amount)
	s.Require().EqualValues(sdkmath.NewInt(45), s.App.BankKeeper.GetBalance(s.Ctx, sender, denom).Amount)
	s.Require().EqualValues(attackerBalanceBefore, s.App.BankKeeper.GetBalance(s.Ctx, attacker, denom))
}

func (s *KeeperTestSuite) TestKeeper_IncreaseBridgeFeeRejectsInvalidStoredSender() {
	sender := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	bridgeToken := helpers.GenerateAddress().Hex()

	denom := "testbadfee"
	initialSupply := sdk.NewCoin(denom, sdkmath.NewInt(100))
	s.Equal(sdk.NewCoin(denom, sdkmath.ZeroInt()), s.App.BankKeeper.GetSupply(s.Ctx, denom))

	err := s.App.BankKeeper.MintCoins(s.Ctx, s.chainName, sdk.NewCoins(initialSupply))
	s.NoError(err)
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, s.chainName, sender, sdk.NewCoins(initialSupply))
	s.NoError(err)

	s.Keeper().SetBridgeToken(s.Ctx, &types.BridgeToken{ContractAddress: bridgeToken, Denom: denom, Supply: initialSupply.Amount})

	tx := &types.OutgoingTransferTx{
		Id:          1,
		Sender:      "not-a-bech32-address",
		DestAddress: helpers.GenerateAddress().Hex(),
		Token:       types.NewERC20Token(sdkmath.NewInt(40), bridgeToken),
		Fee:         types.NewERC20Token(sdkmath.NewInt(10), bridgeToken),
	}
	err = s.Keeper().AddUnbatchedTx(s.Ctx, tx)
	s.Require().NoError(err)

	balanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, sender, denom)

	_, err = s.MsgServer().IncreaseBridgeFee(sdk.WrapSDKContext(s.Ctx), &types.MsgIncreaseBridgeFee{
		ChainName:     s.chainName,
		Sender:        sender.String(),
		TransactionId: tx.Id,
		AddBridgeFee:  sdk.NewCoin(denom, sdkmath.NewInt(5)),
	})
	s.Require().Error(err)
	s.Require().ErrorIs(err, types.ErrInvalid)
	s.Require().Contains(err.Error(), "invalid tx sender")

	s.Require().Equal(balanceBefore, s.App.BankKeeper.GetBalance(s.Ctx, sender, denom))
	storedTx, err := s.Keeper().GetUnbatchedTxById(s.Ctx, tx.Id)
	s.Require().NoError(err)
	s.Require().EqualValues(tx.Fee.Amount, storedTx.Fee.Amount)
}

func (s *KeeperTestSuite) TestKeeper_IncreaseBridgeFeeDecrementsBridgeTokenSupply() {
	sender := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	bridgeToken := helpers.GenerateAddress().Hex()

	denom := "testsupply"
	initialSupply := sdk.NewCoin(denom, sdkmath.NewInt(150))
	err := s.App.BankKeeper.MintCoins(s.Ctx, s.chainName, sdk.NewCoins(initialSupply))
	s.NoError(err)
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, s.chainName, sender, sdk.NewCoins(initialSupply))
	s.NoError(err)

	bridgeTokenState := &types.BridgeToken{ContractAddress: bridgeToken, Denom: denom, Supply: initialSupply.Amount}
	s.Keeper().SetBridgeToken(s.Ctx, bridgeTokenState)

	receiver := helpers.GenerateAddress().Hex()
	amount := sdk.NewCoin(denom, sdkmath.NewInt(40))
	fee := sdk.NewCoin(denom, sdkmath.NewInt(10))
	txId, err := s.Keeper().AddToOutgoingPool(s.Ctx, sender, receiver, amount, fee)
	s.NoError(err)

	s.Require().EqualValues(sdkmath.NewInt(100), s.App.BankKeeper.GetSupply(s.Ctx, denom).Amount)
	bridgeTokenState, err = s.Keeper().GetBridgeTokenByDenom(s.Ctx, denom)
	s.Require().NoError(err)
	s.Require().EqualValues(sdkmath.NewInt(100), bridgeTokenState.Supply)

	_, err = s.MsgServer().IncreaseBridgeFee(sdk.WrapSDKContext(s.Ctx), &types.MsgIncreaseBridgeFee{
		ChainName:     s.chainName,
		Sender:        sender.String(),
		TransactionId: txId,
		AddBridgeFee:  sdk.NewCoin(denom, sdkmath.NewInt(5)),
	})
	s.Require().NoError(err)

	s.Require().EqualValues(sdkmath.NewInt(95), s.App.BankKeeper.GetSupply(s.Ctx, denom).Amount)
	bridgeTokenState, err = s.Keeper().GetBridgeTokenByDenom(s.Ctx, denom)
	s.Require().NoError(err)
	s.Require().EqualValues(sdkmath.NewInt(95), bridgeTokenState.Supply)

	tx, err := s.Keeper().GetUnbatchedTxById(s.Ctx, txId)
	s.Require().NoError(err)
	s.Require().EqualValues(sdkmath.NewInt(15), tx.Fee.Amount)
}
