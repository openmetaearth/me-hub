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
