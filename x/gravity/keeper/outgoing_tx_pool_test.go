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

func (s *KeeperTestSuite) TestAddUnbatchedTxRejectsDuplicateTxIDWithDifferentFee() {
	tokenContract := helpers.GenerateAddress().Hex()
	sender := sdk.AccAddress(helpers.GenerateAddress().Bytes()).String()
	dest := helpers.GenerateAddress().Hex()

	first := &types.OutgoingTransferTx{
		Id:          7,
		Sender:      sender,
		DestAddress: dest,
		Token: types.ERC20Token{
			Contract: tokenContract,
			Amount:   sdkmath.NewInt(100),
		},
		Fee: types.ERC20Token{
			Contract: tokenContract,
			Amount:   sdkmath.NewInt(1),
		},
	}
	duplicate := *first
	duplicate.Fee = types.ERC20Token{
		Contract: tokenContract,
		Amount:   sdkmath.NewInt(2),
	}

	s.NoError(s.Keeper().AddUnbatchedTx(s.Ctx, first))
	err := s.Keeper().AddUnbatchedTx(s.Ctx, &duplicate)
	s.Error(err)
	s.Contains(err.Error(), "transaction id 7 already in pool")

	unbatched := s.Keeper().GetUnbatchedTransactions(s.Ctx)
	s.Len(unbatched, 1)
	s.Equal(uint64(7), unbatched[0].Id)
	s.Equal(sdkmath.NewInt(1), unbatched[0].Fee.Amount)
}
