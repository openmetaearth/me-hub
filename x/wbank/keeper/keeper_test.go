package keeper_test

import (
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	wbanktypes "github.com/openmetaearth/me-hub/x/wbank/types"
	"github.com/stretchr/testify/require"
)

func TestFeeToReceiversRejectsBlockedReceiverBeforeTransfer(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	feePayer := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	allowedReceiver := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	blockedReceiver := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	require.True(t, app.BankKeeper.BlockedAddr(blockedReceiver))

	fee := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)))
	err := banktestutil.FundAccount(app.BankKeeper, ctx, feePayer, fee)
	require.NoError(t, err)

	feePayerBalanceBefore := app.BankKeeper.GetAllBalances(ctx, feePayer)
	allowedReceiverBalanceBefore := app.BankKeeper.GetAllBalances(ctx, allowedReceiver)
	blockedReceiverBalanceBefore := app.BankKeeper.GetAllBalances(ctx, blockedReceiver)

	err = app.BankKeeper.FeeToReceivers(
		ctx,
		[]banktypes.Input{{
			Address: feePayer.String(),
			Coins:   fee,
		}},
		[]banktypes.Output{
			{
				Address: allowedReceiver.String(),
				Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(40))),
			},
			{
				Address: blockedReceiver.String(),
				Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(60))),
			},
		},
		[]wbanktypes.FeeReceiverType{
			wbanktypes.FeeReceiverDevOperator,
			wbanktypes.FeeReceiverGlobalDaoFeePool,
		},
	)

	require.Error(t, err)
	require.True(t, sdkerrors.ErrUnauthorized.Is(err))
	require.Equal(t, feePayerBalanceBefore, app.BankKeeper.GetAllBalances(ctx, feePayer))
	require.Equal(t, allowedReceiverBalanceBefore, app.BankKeeper.GetAllBalances(ctx, allowedReceiver))
	require.Equal(t, blockedReceiverBalanceBefore, app.BankKeeper.GetAllBalances(ctx, blockedReceiver))
}

func TestFeeToReceiversRejectsReceiverTypeMismatchBeforeTransfer(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	feePayer := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	receiverA := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	receiverB := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	fee := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)))
	err := banktestutil.FundAccount(app.BankKeeper, ctx, feePayer, fee)
	require.NoError(t, err)

	feePayerBalanceBefore := app.BankKeeper.GetAllBalances(ctx, feePayer)
	receiverABalanceBefore := app.BankKeeper.GetAllBalances(ctx, receiverA)
	receiverBBalanceBefore := app.BankKeeper.GetAllBalances(ctx, receiverB)

	err = app.BankKeeper.FeeToReceivers(
		ctx,
		[]banktypes.Input{{
			Address: feePayer.String(),
			Coins:   fee,
		}},
		[]banktypes.Output{
			{
				Address: receiverA.String(),
				Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(40))),
			},
			{
				Address: receiverB.String(),
				Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(60))),
			},
		},
		[]wbanktypes.FeeReceiverType{
			wbanktypes.FeeReceiverDevOperator,
		},
	)

	require.Error(t, err)
	require.True(t, sdkerrors.ErrInvalidRequest.Is(err))
	require.ErrorContains(t, err, "fee receiver types and outputs are not equal")
	require.Equal(t, feePayerBalanceBefore, app.BankKeeper.GetAllBalances(ctx, feePayer))
	require.Equal(t, receiverABalanceBefore, app.BankKeeper.GetAllBalances(ctx, receiverA))
	require.Equal(t, receiverBBalanceBefore, app.BankKeeper.GetAllBalances(ctx, receiverB))
}

func TestFeeToReceiversRejectsInvalidReceiverBeforeTransfer(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.GetBaseApp().NewContext(false, cometbftproto.Header{})

	feePayer := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	allowedReceiver := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	invalidReceiver := "not-a-bech32-address"

	fee := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)))
	err := banktestutil.FundAccount(app.BankKeeper, ctx, feePayer, fee)
	require.NoError(t, err)

	feePayerBalanceBefore := app.BankKeeper.GetAllBalances(ctx, feePayer)
	allowedReceiverBalanceBefore := app.BankKeeper.GetAllBalances(ctx, allowedReceiver)

	err = app.BankKeeper.FeeToReceivers(
		ctx,
		[]banktypes.Input{{
			Address: feePayer.String(),
			Coins:   fee,
		}},
		[]banktypes.Output{
			{
				Address: allowedReceiver.String(),
				Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(40))),
			},
			{
				Address: invalidReceiver,
				Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(60))),
			},
		},
		[]wbanktypes.FeeReceiverType{
			wbanktypes.FeeReceiverDevOperator,
			wbanktypes.FeeReceiverGlobalDaoFeePool,
		},
	)

	require.Error(t, err)
	require.True(t, sdkerrors.ErrInvalidAddress.Is(err))
	require.Equal(t, feePayerBalanceBefore, app.BankKeeper.GetAllBalances(ctx, feePayer))
	require.Equal(t, allowedReceiverBalanceBefore, app.BankKeeper.GetAllBalances(ctx, allowedReceiver))
}
