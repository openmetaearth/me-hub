package ante_test

import (
	"regexp"
	"strconv"
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authantetestutil "github.com/cosmos/cosmos-sdk/x/auth/ante/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/golang/mock/gomock"
	"github.com/openmetaearth/me-hub/app/ante"
	"github.com/openmetaearth/me-hub/app/params"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/ante/mock"
	"github.com/stretchr/testify/require"
)

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

func TestMockBankKeeper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBankKeeper := mock.NewMockBankKeeper(ctrl)

	ctx := sdk.Context{}
	addr := NewAccount().GetAddress()
	expectedBalances := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100)))

	mockBankKeeper.EXPECT().GetAllBalances(ctx, addr).Return(expectedBalances)

	balances := mockBankKeeper.GetAllBalances(ctx, addr)
	require.Equal(t, expectedBalances, balances)
}

func TestCheckFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := sdk.Context{}
	mockBankKeeper := mock.NewMockBankKeeper(ctrl)
	mockAccountKeeper := authantetestutil.NewMockAccountKeeper(ctrl)
	mockFeegrantKeeper := authantetestutil.NewMockFeegrantKeeper(ctrl)
	mockStakingKeeper := mock.NewMockStakingKeeper(ctrl)
	mockKycKeeper := mock.NewMockKycKeeper(ctrl)
	mockDaoKeeper := mock.NewMockDaoKeeper(ctrl)
	mockWasmKeeper := mock.NewMockWasmKeeper(ctrl)

	decorator := ante.NewDeductFeeDecorator(
		mockAccountKeeper,
		mockBankKeeper,
		mockFeegrantKeeper,
		mockDaoKeeper,
		mockStakingKeeper,
		mockKycKeeper,
		nil,
		mockWasmKeeper,
	)

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
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(200))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(50))),
				},
			},
			expectError: false,
		},
		{
			name:     "MsgSend with insufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(50))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(50))),
				},
			},
			expectError:  true,
			expectAmount: 150,
		},
		{
			name:     "MsgSend with sufficient funds, different fee payer",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				sender.Address:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(50))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(50))),
				},
			},
			expectError: false,
		},
		{
			name:     "MsgSend with insufficient funds, different fee payer, fee payer is no enough",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(200))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				sender.Address:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
			},
			expectError:  true,
			expectAmount: 200,
		},
		{
			name:     "Multi MsgSend with insufficient funds, different fee payer",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(50))),
				sender.Address:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(400))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: sender.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
			},
			expectError:  true,
			expectAmount: 100,
		},
		{
			name:     "Multi MsgSend with sufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(500))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
			},
			expectError: false,
		},
		{
			name:     "Multi MsgSend with insufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(400))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
				&banktypes.MsgSend{
					FromAddress: feePayer.Address,
					ToAddress:   receiver.Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
				},
			},
			expectError:  true,
			expectAmount: 500,
		},
		{
			name:     "MsgDelegate with sufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(300))),
			},
			messages: []sdk.Msg{
				&stakingtypes.MsgDelegate{
					DelegatorAddress: feePayer.Address,
					ValidatorAddress: sdk.ValAddress(receiver.GetAddress()).String(),
					Amount:           sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(150)),
				},
			},
			expectError: false,
		},
		{
			name:     "MsgMultiSend with sufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(300))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgMultiSend{
					Inputs: []banktypes.Input{
						{
							Address: feePayer.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(200))),
						},
					},
					Outputs: []banktypes.Output{
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
						},
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:     "MsgMultiSend with insufficient funds",
			feePayer: feePayer.Address,
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(200))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgMultiSend{
					Inputs: []banktypes.Input{
						{
							Address: feePayer.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(200))),
						},
					},
					Outputs: []banktypes.Output{
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(200))),
						},
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
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
			fees:     sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
			balances: map[string]sdk.Coins{
				feePayer.Address: sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(200))),
			},
			messages: []sdk.Msg{
				&banktypes.MsgMultiSend{
					Inputs: []banktypes.Input{
						{
							Address: feePayer.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(200))),
						},
					},
					Outputs: []banktypes.Output{
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
						},
						{
							Address: receiver.Address,
							Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100))),
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
			// Mock the balances for all involved addresses
			for address, balance := range tc.balances {
				mockBankKeeper.EXPECT().
					GetAllBalances(gomock.Any(), sdk.MustAccAddressFromBech32(address)).
					Return(balance)
			}

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
			}
		})
	}
}

func newDeductFeeDecoratorWithWasm(t *testing.T) (*ante.DeductFeeDecorator, *mock.MockWasmKeeper, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockWasmKeeper := mock.NewMockWasmKeeper(ctrl)
	d := ante.NewDeductFeeDecorator(
		authantetestutil.NewMockAccountKeeper(ctrl),
		mock.NewMockBankKeeper(ctrl),
		authantetestutil.NewMockFeegrantKeeper(ctrl),
		mock.NewMockDaoKeeper(ctrl),
		mock.NewMockStakingKeeper(ctrl),
		mock.NewMockKycKeeper(ctrl),
		nil,
		mockWasmKeeper,
	)
	return &d, mockWasmKeeper, ctrl.Finish
}

func TestParseWasmMsgContractCreator(t *testing.T) {
	ctx := sdk.Context{}.WithEventManager(sdk.NewEventManager())

	contractAddr := NewAccount()
	creatorAddr := NewAccount()

	tests := []struct {
		name        string
		msgs        []sdk.Msg
		setupMock   func(wk *mock.MockWasmKeeper)
		wantCreator string
		wantOk      bool
	}{
		{
			name: "single MsgExecuteContract returns creator",
			msgs: []sdk.Msg{
				&wasmtypes.MsgExecuteContract{
					Sender:   NewAccount().Address,
					Contract: contractAddr.Address,
					Msg:      []byte(`{}`),
				},
			},
			setupMock: func(wk *mock.MockWasmKeeper) {
				wk.EXPECT().
					GetContractInfo(gomock.Any(), contractAddr.GetAddress()).
					Return(&wasmtypes.ContractInfo{Creator: creatorAddr.Address})
			},
			wantCreator: creatorAddr.Address,
			wantOk:      true,
		},
		{
			name: "two wasm messages rejected — issue #23 regression",
			msgs: []sdk.Msg{
				&wasmtypes.MsgExecuteContract{
					Sender:   NewAccount().Address,
					Contract: contractAddr.Address,
					Msg:      []byte(`{"target":{}}`),
				},
				&wasmtypes.MsgExecuteContract{
					Sender:   NewAccount().Address,
					Contract: NewAccount().Address,
					Msg:      []byte(`{"noop":{}}`),
				},
			},
			setupMock: func(wk *mock.MockWasmKeeper) {
				// GetContractInfo must never be called when more than one msg
				wk.EXPECT().GetContractInfo(gomock.Any(), gomock.Any()).Times(0)
			},
			wantCreator: "",
			wantOk:      false,
		},
		{
			name: "zero messages rejected",
			msgs: []sdk.Msg{},
			setupMock: func(wk *mock.MockWasmKeeper) {
				wk.EXPECT().GetContractInfo(gomock.Any(), gomock.Any()).Times(0)
			},
			wantCreator: "",
			wantOk:      false,
		},
		{
			name: "single non-wasm message rejected",
			msgs: []sdk.Msg{
				&banktypes.MsgSend{
					FromAddress: NewAccount().Address,
					ToAddress:   NewAccount().Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(1))),
				},
			},
			setupMock: func(wk *mock.MockWasmKeeper) {
				wk.EXPECT().GetContractInfo(gomock.Any(), gomock.Any()).Times(0)
			},
			wantCreator: "",
			wantOk:      false,
		},
		{
			name: "single wasm + one non-wasm message rejected",
			msgs: []sdk.Msg{
				&wasmtypes.MsgExecuteContract{
					Sender:   NewAccount().Address,
					Contract: contractAddr.Address,
					Msg:      []byte(`{}`),
				},
				&banktypes.MsgSend{
					FromAddress: NewAccount().Address,
					ToAddress:   NewAccount().Address,
					Amount:      sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(1))),
				},
			},
			setupMock: func(wk *mock.MockWasmKeeper) {
				wk.EXPECT().GetContractInfo(gomock.Any(), gomock.Any()).Times(0)
			},
			wantCreator: "",
			wantOk:      false,
		},
		{
			name: "single wasm message with nil contract info returns false",
			msgs: []sdk.Msg{
				&wasmtypes.MsgExecuteContract{
					Sender:   NewAccount().Address,
					Contract: contractAddr.Address,
					Msg:      []byte(`{}`),
				},
			},
			setupMock: func(wk *mock.MockWasmKeeper) {
				wk.EXPECT().
					GetContractInfo(gomock.Any(), contractAddr.GetAddress()).
					Return(nil)
			},
			wantCreator: "",
			wantOk:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, mockWasm, finish := newDeductFeeDecoratorWithWasm(t)
			defer finish()
			tc.setupMock(mockWasm)

			tx := &mock.MockTx{Msgs: tc.msgs}
			creator, ok := d.ParseWasmMsgContractCreator(ctx, tx)
			require.Equal(t, tc.wantOk, ok)
			require.Equal(t, tc.wantCreator, creator)
		})
	}
}
