package ante_test

import (
	"regexp"
	"strconv"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authantetestutil "github.com/cosmos/cosmos-sdk/x/auth/ante/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/golang/mock/gomock"
	"github.com/openmetaearth/me-hub/app/ante"
	"github.com/openmetaearth/me-hub/app/ante/mock"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/stretchr/testify/require"
)

const anteMsgLimit = 1000

type feeTxWithMsgs struct {
	msgs       []sdk.Msg
	feePayer   sdk.AccAddress
	feeGranter sdk.AccAddress
}

func (tx feeTxWithMsgs) GetMsgs() []sdk.Msg {
	return tx.msgs
}

func (tx feeTxWithMsgs) ValidateBasic() error {
	return nil
}

func (tx feeTxWithMsgs) GetGas() uint64 {
	return 1
}

func (tx feeTxWithMsgs) GetFee() sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(10000)))
}

func (tx feeTxWithMsgs) FeePayer() sdk.AccAddress {
	return tx.feePayer
}

func (tx feeTxWithMsgs) FeeGranter() sdk.AccAddress {
	return tx.feeGranter
}

func NewAccount() *authtypes.BaseAccount {
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	return acc
}

func NewAccountWithEthPrivKey() (*authtypes.BaseAccount, *ethsecp256k1.PrivKey) {
	senderPrivKey, _ := ethsecp256k1.GenerateKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	return acc, senderPrivKey
}

func TestDeductFeeDecoratorCountsAuthzMessagesTowardLimit(t *testing.T) {
	grantee := NewAccount()
	sender := NewAccount()
	recipient := NewAccount()

	innerMsgs := makeMsgSends(sender.GetAddress(), recipient.GetAddress(), anteMsgLimit+1)
	nestedExec := newMsgExec(grantee.GetAddress(), []sdk.Msg{
		newMsgExec(grantee.GetAddress(), makeMsgSends(sender.GetAddress(), recipient.GetAddress(), anteMsgLimit)),
	})

	tests := []struct {
		name string
		msgs []sdk.Msg
	}{
		{
			name: "top-level messages",
			msgs: makeMsgSends(sender.GetAddress(), recipient.GetAddress(), anteMsgLimit+1),
		},
		{
			name: "authz inner messages",
			msgs: []sdk.Msg{newMsgExec(grantee.GetAddress(), innerMsgs)},
		},
		{
			name: "nested authz messages",
			msgs: []sdk.Msg{nestedExec},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			decorator := newDeductFeeDecoratorForMsgLimitTest(t)
			tx := feeTxWithMsgs{
				msgs:     tc.msgs,
				feePayer: grantee.GetAddress(),
			}

			_, err := decorator.AnteHandle(sdk.Context{}, tx, false, failAnteHandler(t))

			require.Error(t, err)
			require.True(t, sdkerrors.ErrInvalidRequest.Is(err))
			require.Contains(t, err.Error(), "messages should not exceed 1000")
		})
	}
}

func TestDeductFeeDecoratorRejectsUnpackedAuthzMessages(t *testing.T) {
	grantee := NewAccount()
	decorator := newDeductFeeDecoratorForMsgLimitTest(t)
	tx := feeTxWithMsgs{
		msgs: []sdk.Msg{
			&authz.MsgExec{
				Grantee: grantee.Address,
				Msgs: []*codectypes.Any{
					{TypeUrl: "/cosmos.bank.v1beta1.MsgSend"},
				},
			},
		},
		feePayer: grantee.GetAddress(),
	}

	_, err := decorator.AnteHandle(sdk.Context{}, tx, false, failAnteHandler(t))

	require.Error(t, err)
	require.True(t, sdkerrors.ErrUnpackAny.Is(err))
}

func TestDeductFeeDecoratorRejectsDeeplyNestedAuthzMessages(t *testing.T) {
	grantee := NewAccount()
	sender := NewAccount()
	recipient := NewAccount()
	decorator := newDeductFeeDecoratorForMsgLimitTest(t)

	var msg sdk.Msg = banktypes.NewMsgSend(
		sender.GetAddress(),
		recipient.GetAddress(),
		sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(1))),
	)
	for i := 0; i < 7; i++ {
		msg = newMsgExec(grantee.GetAddress(), []sdk.Msg{msg})
	}

	tx := feeTxWithMsgs{
		msgs:     []sdk.Msg{msg},
		feePayer: grantee.GetAddress(),
	}

	_, err := decorator.AnteHandle(sdk.Context{}, tx, false, failAnteHandler(t))

	require.Error(t, err)
	require.True(t, sdkerrors.ErrInvalidRequest.Is(err))
	require.Contains(t, err.Error(), "authz MsgExec nesting exceeds")
}

func TestDeductFeeDecoratorAllowsAuthzMessagesAtLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockDaoKeeper := mock.NewMockDaoKeeper(ctrl)
	decorator := ante.NewDeductFeeDecorator(
		authantetestutil.NewMockAccountKeeper(ctrl),
		mock.NewMockBankKeeper(ctrl),
		authantetestutil.NewMockFeegrantKeeper(ctrl),
		mockDaoKeeper,
		mock.NewMockStakingKeeper(ctrl),
		mock.NewMockKycKeeper(ctrl),
		nil,
		mock.NewMockWasmKeeper(ctrl),
	)

	grantee := NewAccount()
	sender := NewAccount()
	recipient := NewAccount()

	// Force the free-gas path after message counting so reaching next does not depend on fee-transfer keepers.
	mockDaoKeeper.EXPECT().IsDao(gomock.Any(), grantee.Address).Return(true)
	mockDaoKeeper.EXPECT().CheckFreeGasAccount(gomock.Any(), grantee.Address).Return(false)

	tx := feeTxWithMsgs{
		msgs: []sdk.Msg{
			newMsgExec(grantee.GetAddress(), makeMsgSends(sender.GetAddress(), recipient.GetAddress(), anteMsgLimit-1)),
		},
		feePayer: grantee.GetAddress(),
	}
	nextCalled := false

	_, err := decorator.AnteHandle(sdk.Context{}, tx, false, func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		nextCalled = true
		return ctx, nil
	})

	require.NoError(t, err)
	require.True(t, nextCalled)
}

func TestMockBankKeeper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankKeeper := mock.NewMockBankKeeper(ctrl)

	ctx := sdk.Context{}
	addr := NewAccount().GetAddress()
	expectedBalances := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)))

	mockBankKeeper.EXPECT().GetAllBalances(ctx, addr).Return(expectedBalances)

	balances := mockBankKeeper.GetAllBalances(ctx, addr)
	require.Equal(t, expectedBalances, balances)
}

func newDeductFeeDecoratorForMsgLimitTest(t *testing.T) ante.DeductFeeDecorator {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	return ante.NewDeductFeeDecorator(
		authantetestutil.NewMockAccountKeeper(ctrl),
		mock.NewMockBankKeeper(ctrl),
		authantetestutil.NewMockFeegrantKeeper(ctrl),
		mock.NewMockDaoKeeper(ctrl),
		mock.NewMockStakingKeeper(ctrl),
		mock.NewMockKycKeeper(ctrl),
		nil,
		mock.NewMockWasmKeeper(ctrl),
	)
}

func makeMsgSends(from, to sdk.AccAddress, count int) []sdk.Msg {
	msgs := make([]sdk.Msg, count)
	amount := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(1)))
	for i := range msgs {
		msgs[i] = banktypes.NewMsgSend(from, to, amount)
	}

	return msgs
}

func newMsgExec(grantee sdk.AccAddress, msgs []sdk.Msg) *authz.MsgExec {
	msg := authz.NewMsgExec(grantee, msgs)
	return &msg
}

func failAnteHandler(t *testing.T) sdk.AnteHandler {
	t.Helper()

	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		t.Fatalf("next ante handler should not run")
		return ctx, nil
	}
}

func TestCheckFunds(t *testing.T) {
	ctx := sdk.Context{}
	feePayer := NewAccount()
	receiver := NewAccount()
	sender := NewAccount()

	tests := []struct {
		name         string
		feePayer     string
		fees         sdk.Coins
		balances     map[string]sdk.Coins
		messages     []sdk.Msg
		expectError  bool
		expectAmount int64
	}{
		{
			name:     "MsgSend with sufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(200))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(50))),
				},
			},
			expectError: false,
		},
		{
			name:     "MsgSend with insufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(50))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(50))),
				},
			},
			expectError:  true,
			expectAmount: 150,
		},
		{
			name:     "MsgSend with sufficient funds, different fee payer",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				sender.Address:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(50))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(50))),
				},
			},
			expectError: false,
		},
		{
			name:     "MsgSend with insufficient funds, different fee payer, fee payer is no enough",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(200))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				sender.Address:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
			},
			expectError:  true,
			expectAmount: 200,
		},
		{
			name:     "Multi MsgSend with insufficient funds, different fee payer",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(50))),
				sender.Address:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(400))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
			},
			expectError:  true,
			expectAmount: 100,
		},
		{
			name:     "Multi MsgSend with sufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(500))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
			},
			expectError: false,
		},
		{
			name:     "Multi MsgSend with insufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(400))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
				},
			},
			expectError:  true,
			expectAmount: 500,
		},
		{
			name:     "MsgDelegate with sufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(300))),
			},
			messages: []sdk.Msg{
				&stakingtypes.MsgDelegate{
					DelegatorAddress: feePayer.Address,
					ValidatorAddress: sdk.ValAddress(receiver.GetAddress()).String(),
					Amount:           sdk.NewCoin(params.BaseDenom, sdk.NewInt(150)),
				},
			},
			expectError: false,
		},
		{
			name:     "MsgMultiSend with sufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(300))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgMultiSend{
					Inputs: []banktypes.Input{
						{
							Address: feePayer.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(200))),
						},
					},
					Outputs: []banktypes.Output{
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
						},
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:     "MsgMultiSend with insufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(200))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgMultiSend{
					Inputs: []banktypes.Input{
						{
							Address: feePayer.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(200))),
						},
					},
					Outputs: []banktypes.Output{
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(200))),
						},
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
						},
					},
				},
			},
			expectError:  true,
			expectAmount: 400,
		},
		{
			name:     "MsgMultiSend with insufficient funds, not enough for fees",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(200))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgMultiSend{
					Inputs: []banktypes.Input{
						{
							Address: feePayer.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(200))),
						},
					},
					Outputs: []banktypes.Output{
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
						},
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100))),
						},
					},
				},
			},
			expectError:  true,
			expectAmount: 300,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockBankKeeper := mock.NewMockBankKeeper(ctrl)
			decorator := ante.NewDeductFeeDecorator(
				authantetestutil.NewMockAccountKeeper(ctrl),
				mockBankKeeper,
				authantetestutil.NewMockFeegrantKeeper(ctrl),
				mock.NewMockDaoKeeper(ctrl),
				mock.NewMockStakingKeeper(ctrl),
				mock.NewMockKycKeeper(ctrl),
				nil,
				mock.NewMockWasmKeeper(ctrl),
			)

			calledBalances := map[string]struct{}{}
			mockBankKeeper.EXPECT().
				GetAllBalances(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ sdk.Context, addr sdk.AccAddress) sdk.Coins {
					address := addr.String()
					balance, ok := tc.balances[address]
					require.True(t, ok, "unexpected balance lookup for %s", address)
					calledBalances[address] = struct{}{}
					return balance
				}).
				MinTimes(1).
				MaxTimes(len(tc.balances))

			// Create a mock transaction with the provided messages
			tx := &mock.MockTx{Msgs: tc.messages}

			// Call CheckFunds
			err := decorator.CheckFunds(ctx, tx, tc.feePayer, tc.fees)

			// Assert the result
			if tc.expectError {
				require.Error(t, err)
				require.True(t, sdkerrors.ErrInsufficientFunds.Is(err))

				re := regexp.MustCompile(`required:\s(\d+)`)
				matches := re.FindStringSubmatch(err.Error())
				if len(matches) > 1 {
					requiredAmount, convErr := strconv.ParseInt(matches[1], 10, 64)
					require.NoError(t, convErr)
					require.Equal(t, tc.expectAmount, requiredAmount, "Required amount does not match expected value")
				} else {
					t.Errorf("Failed to extract required amount from error: %s", err.Error())
				}
			} else {
				require.NoError(t, err)
				require.Len(t, calledBalances, len(tc.balances))
			}
		})
	}
}
