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

func (s *KeeperTestSuite) TestSendToExternalPendingDoesNotBlockRemainingSupply() {
	attacker := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	victim := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	receiver := helpers.GenerateAddress().Hex()
	denom := "pendingdoublecount"

	bridgeToken := s.NewBridgeToken(attacker, sdk.NewCoin(denom, sdkmath.NewInt(60)))
	victimInitialBalance := sdk.NewCoin(denom, sdkmath.NewInt(40))
	s.Require().NoError(s.App.BankKeeper.MintCoins(s.Ctx, s.chainName, sdk.NewCoins(victimInitialBalance)))
	s.Require().NoError(s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, s.chainName, victim, sdk.NewCoins(victimInitialBalance)))

	bridgeToken.Supply = bridgeToken.Supply.Add(victimInitialBalance.Amount)
	s.Keeper().SetBridgeToken(s.Ctx, &bridgeToken)

	_, err := s.MsgServer().SendToExternal(sdk.WrapSDKContext(s.Ctx), &types.MsgSendToExternal{
		Sender:    attacker.String(),
		Dest:      receiver,
		Amount:    sdk.NewCoin(denom, sdkmath.NewInt(50)),
		BridgeFee: sdk.NewCoin(denom, sdkmath.NewInt(10)),
		ChainName: s.chainName,
	})
	s.Require().NoError(err)

	storedBridgeToken, err := s.Keeper().GetBridgeTokenByDenom(s.Ctx, denom)
	s.Require().NoError(err)
	s.Require().Equal(sdkmath.NewInt(40), storedBridgeToken.Supply)
	s.Require().Equal(sdkmath.NewInt(40), s.App.BankKeeper.GetBalance(s.Ctx, victim, denom).Amount)
	s.Require().Equal(sdkmath.NewInt(60), s.Keeper().GetOutgoingPendingTxTotal(s.Ctx, s.chainName, &bridgeToken))

	_, err = s.MsgServer().SendToExternal(sdk.WrapSDKContext(s.Ctx), &types.MsgSendToExternal{
		Sender:    victim.String(),
		Dest:      receiver,
		Amount:    sdk.NewCoin(denom, sdkmath.NewInt(1)),
		BridgeFee: sdk.NewCoin(denom, sdkmath.ZeroInt()),
		ChainName: s.chainName,
	})
	s.Require().NoError(err)

	storedBridgeToken, err = s.Keeper().GetBridgeTokenByDenom(s.Ctx, denom)
	s.Require().NoError(err)
	s.Require().Equal(sdkmath.NewInt(39), storedBridgeToken.Supply)
	s.Require().Equal(sdkmath.NewInt(39), s.App.BankKeeper.GetBalance(s.Ctx, victim, denom).Amount)
}
