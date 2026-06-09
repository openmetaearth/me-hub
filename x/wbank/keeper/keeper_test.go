package keeper_test

import (
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	wbanktypes "github.com/openmetaearth/me-hub/x/wbank/types"
	"github.com/stretchr/testify/require"
)

func TestFeeToReceiversRejectsReceiverTypeMismatchBeforeTransfer(t *testing.T) {
	meApp := apptesting.Setup(t, false)
	ctx := meApp.BaseApp.NewContext(false, cometbftproto.Header{}).WithBlockHeight(1)

	sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	receiverA := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	receiverB := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	inputCoins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100))
	outputCoinsA := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 40))
	outputCoinsB := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 60))

	err := bankutil.FundAccount(meApp.BankKeeper, ctx, sender, inputCoins)
	require.NoError(t, err)

	inputs := []banktypes.Input{{
		Address: sender.String(),
		Coins:   inputCoins,
	}}
	outputs := []banktypes.Output{
		{Address: receiverA.String(), Coins: outputCoinsA},
		{Address: receiverB.String(), Coins: outputCoinsB},
	}
	receiverTypes := []wbanktypes.FeeReceiverType{wbanktypes.FeeReceiverDevOperator}

	senderBefore := meApp.BankKeeper.GetAllBalances(ctx, sender)
	receiverABefore := meApp.BankKeeper.GetAllBalances(ctx, receiverA)
	receiverBBefore := meApp.BankKeeper.GetAllBalances(ctx, receiverB)

	err = meApp.BankKeeper.FeeToReceivers(ctx, inputs, outputs, receiverTypes)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	require.ErrorContains(t, err, "fee receiver types and outputs are not equal")

	require.True(t, senderBefore.IsEqual(meApp.BankKeeper.GetAllBalances(ctx, sender)))
	require.True(t, receiverABefore.IsEqual(meApp.BankKeeper.GetAllBalances(ctx, receiverA)))
	require.True(t, receiverBBefore.IsEqual(meApp.BankKeeper.GetAllBalances(ctx, receiverB)))
}
