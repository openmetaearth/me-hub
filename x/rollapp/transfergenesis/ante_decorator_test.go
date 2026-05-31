package transfergenesis

import (
	"testing"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	rollapptypes "github.com/openmetaearth/me-hub/x/rollapp/types"
	"github.com/stretchr/testify/require"
)

type transferEnabledTestTx struct {
	msgs []sdk.Msg
}

func (tx transferEnabledTestTx) GetMsgs() []sdk.Msg {
	return tx.msgs
}

func (tx transferEnabledTestTx) ValidateBasic() error {
	return nil
}

type transferEnabledChannelKeeper struct {
	chainID string
}

func (k transferEnabledChannelKeeper) GetChannelClientState(
	sdk.Context,
	string,
	string,
) (string, exported.ClientState, error) {
	return "client-0", &ibctmtypes.ClientState{ChainId: k.chainID}, nil
}

func TestTransferEnabledDecoratorRejectsDirectTransfer(t *testing.T) {
	transfer := newTestTransfer()
	decorator := newTestTransferEnabledDecorator(false)

	nextCalled := false
	_, err := decorator.AnteHandle(sdk.Context{}, transferEnabledTestTx{msgs: []sdk.Msg{transfer}}, false, func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		nextCalled = true
		return ctx, nil
	})

	require.ErrorContains(t, err, "transfers to/from rollapp are disabled")
	require.False(t, nextCalled)
}

func TestTransferEnabledDecoratorRejectsAuthzWrappedTransfer(t *testing.T) {
	transfer := newTestTransfer()
	exec := newTestMsgExec(t, transfer)
	decorator := newTestTransferEnabledDecorator(false)

	nextCalled := false
	_, err := decorator.AnteHandle(sdk.Context{}, transferEnabledTestTx{msgs: []sdk.Msg{&exec}}, false, func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		nextCalled = true
		return ctx, nil
	})

	require.ErrorContains(t, err, "transfers to/from rollapp are disabled")
	require.False(t, nextCalled)
}

func TestTransferEnabledDecoratorAllowsAuthzWrappedTransferWhenEnabled(t *testing.T) {
	transfer := newTestTransfer()
	exec := newTestMsgExec(t, transfer)
	decorator := newTestTransferEnabledDecorator(true)

	nextCalled := false
	_, err := decorator.AnteHandle(sdk.Context{}, transferEnabledTestTx{msgs: []sdk.Msg{&exec}}, false, func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		nextCalled = true
		return ctx, nil
	})

	require.NoError(t, err)
	require.True(t, nextCalled)
}

func TestTransferEnabledDecoratorRejectsNestedAuthzWrappedTransfer(t *testing.T) {
	transfer := newTestTransfer()
	innerExec := newTestMsgExec(t, transfer)
	outerExec := newTestMsgExec(t, &innerExec)
	decorator := newTestTransferEnabledDecorator(false)

	nextCalled := false
	_, err := decorator.AnteHandle(sdk.Context{}, transferEnabledTestTx{msgs: []sdk.Msg{&outerExec}}, false, func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		nextCalled = true
		return ctx, nil
	})

	require.ErrorContains(t, err, "transfers to/from rollapp are disabled")
	require.False(t, nextCalled)
}

func newTestTransfer() *transfertypes.MsgTransfer {
	return transfertypes.NewMsgTransfer(
		transfertypes.PortID,
		"channel-0",
		sdk.NewCoin("foo", math.NewInt(1)),
		"sender",
		"receiver",
		clienttypes.NewHeight(1, 1),
		0,
		"",
	)
}

func newTestTransferEnabledDecorator(transfersEnabled bool) *TransferEnabledDecorator {
	const rollappID = "rollapp_1-1"
	return NewTransferEnabledDecorator(
		func(ctx sdk.Context, rollappID string) (rollapptypes.Rollapp, bool) {
			return rollapptypes.Rollapp{
				RollappId: rollappID,
				GenesisState: rollapptypes.RollappGenesisState{
					TransfersEnabled: transfersEnabled,
				},
			}, true
		},
		transferEnabledChannelKeeper{chainID: rollappID},
	)
}

func newTestMsgExec(t *testing.T, msgs ...sdk.Msg) authztypes.MsgExec {
	t.Helper()

	anyMsgs := make([]*codectypes.Any, len(msgs))
	for i, msg := range msgs {
		anyMsg, err := codectypes.NewAnyWithValue(msg)
		require.NoError(t, err)
		anyMsgs[i] = anyMsg
	}

	return authztypes.MsgExec{
		Grantee: "me1grantee",
		Msgs:    anyMsgs,
	}
}
