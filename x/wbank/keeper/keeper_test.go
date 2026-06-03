package keeper_test

import (
	"testing"

	cometbftproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	wbanktypes "github.com/openmetaearth/me-hub/x/wbank/types"
	"github.com/stretchr/testify/require"
)

func TestFeeToReceiversRejectsMismatchedReceiverTypesBeforeTransfer(t *testing.T) {
	app := apptesting.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, cometbftproto.Header{}).WithChainID(apptesting.TestChainID)

	sender := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	receiver := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	fee := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)))

	require.NoError(t, bankutil.FundAccount(app.BankKeeper, ctx, sender, fee))

	err := app.BankKeeper.FeeToReceivers(
		ctx,
		[]banktypes.Input{{Address: sender.String(), Coins: fee}},
		[]banktypes.Output{{Address: receiver.String(), Coins: fee}},
		[]wbanktypes.FeeReceiverType{},
	)

	require.Error(t, err)
	require.Equal(t, fee, app.BankKeeper.GetAllBalances(ctx, sender))
	require.True(t, app.BankKeeper.GetAllBalances(ctx, receiver).IsZero())
}
