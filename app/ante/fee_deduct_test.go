package ante_test

import (
	"regexp"
	"strconv"
	"testing"

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
	megrouptypes "github.com/openmetaearth/me-hub/x/megroup/types"

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
	expectedBalances := sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)))

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
	mockMeGroupKeeper := mock.NewMockMeGroupKeeper(ctrl)

	decorator := ante.NewDeductFeeDecorator(
		mockAccountKeeper,
		mockBankKeeper,
		mockFeegrantKeeper,
		mockDaoKeeper,
		mockStakingKeeper,
		mockKycKeeper,
		nil,
		mockWasmKeeper,
		mockMeGroupKeeper,
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

func buildJoinMsg(creator, applicant string, groupID uint64) *megrouptypes.MsgJoinGroup {
	return &megrouptypes.MsgJoinGroup{
		Creator:          creator,
		ApplicantAddress: applicant,
		GroupId:          groupID,
	}
}

func newDeductFeeDecorator(
	ctrl *gomock.Controller,
	dk *mock.MockDaoKeeper,
	gk *mock.MockMeGroupKeeper,
) ante.DeductFeeDecorator {
	ak := authantetestutil.NewMockAccountKeeper(ctrl)
	fk := authantetestutil.NewMockFeegrantKeeper(ctrl)
	sk := mock.NewMockStakingKeeper(ctrl)
	kk := mock.NewMockKycKeeper(ctrl)
	bk := mock.NewMockBankKeeper(ctrl)
	wk := mock.NewMockWasmKeeper(ctrl)
	return ante.NewDeductFeeDecorator(ak, bk, fk, dk, sk, kk, nil, wk, gk)
}

func passThrough(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) { return ctx, nil }

func TestDeductFeeDecorator_JoinGroupValidation(t *testing.T) {
	regionID := "region-1"
	group := megrouptypes.GroupInfo{Id: 1, RegionID: regionID}

	creator := NewAccount()
	applicant := NewAccount()
	daoUser := NewAccount()

	tests := []struct {
		name      string
		feePayer  *authtypes.BaseAccount
		msgs      []sdk.Msg
		setup     func(dk *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper)
		expectErr bool
		errText   string
	}{
		{
			name:     "valid_creator_equals_applicant",
			feePayer: creator,
			msgs:     []sdk.Msg{buildJoinMsg(creator.Address, creator.Address, 1)},
			setup: func(_ *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper) {
				gk.EXPECT().GetGroupInfo(gomock.Any(), uint64(1)).Return(group, true)
				gk.EXPECT().GetDidAndKycActive(gomock.Any(), creator.GetAddress(), regionID).Return("did1", true)
			},
			expectErr: false,
		},
		{
			name:     "valid_dao_creator_for_other_applicant",
			feePayer: daoUser,
			msgs:     []sdk.Msg{buildJoinMsg(daoUser.Address, applicant.Address, 1)},
			setup: func(dk *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper) {
				dk.EXPECT().IsDao(gomock.Any(), daoUser.Address).Return(true)
				gk.EXPECT().GetGroupInfo(gomock.Any(), uint64(1)).Return(group, true)
				gk.EXPECT().GetDidAndKycActive(gomock.Any(), applicant.GetAddress(), regionID).Return("did2", true)
			},
			expectErr: false,
		},
		{
			name:      "reject_group_id_zero",
			feePayer:  creator,
			msgs:      []sdk.Msg{buildJoinMsg(creator.Address, creator.Address, 0)},
			setup:     func(_ *mock.MockDaoKeeper, _ *mock.MockMeGroupKeeper) {},
			expectErr: true,
			errText:   "group_id must be greater than 0",
		},
		{
			name:      "reject_empty_applicant",
			feePayer:  creator,
			msgs:      []sdk.Msg{buildJoinMsg(creator.Address, "", 1)},
			setup:     func(_ *mock.MockDaoKeeper, _ *mock.MockMeGroupKeeper) {},
			expectErr: true,
			errText:   "applicant_address is required",
		},
		{
			name:      "reject_invalid_applicant_address",
			feePayer:  creator,
			msgs:      []sdk.Msg{buildJoinMsg(creator.Address, "not-bech32", 1)},
			setup:     func(_ *mock.MockDaoKeeper, _ *mock.MockMeGroupKeeper) {},
			expectErr: true,
			errText:   "invalid applicant_address",
		},
		{
			name:     "reject_non_dao_creator_for_other_applicant",
			feePayer: creator,
			msgs:     []sdk.Msg{buildJoinMsg(creator.Address, applicant.Address, 1)},
			setup: func(dk *mock.MockDaoKeeper, _ *mock.MockMeGroupKeeper) {
				dk.EXPECT().IsDao(gomock.Any(), creator.Address).Return(false)
			},
			expectErr: true,
			errText:   "creator is neither the applicant nor a DAO admin",
		},
		{
			name:     "reject_group_not_found",
			feePayer: creator,
			msgs:     []sdk.Msg{buildJoinMsg(creator.Address, creator.Address, 99)},
			setup: func(_ *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper) {
				gk.EXPECT().GetGroupInfo(gomock.Any(), uint64(99)).Return(megrouptypes.GroupInfo{}, false)
			},
			expectErr: true,
			errText:   "group 99 does not exist",
		},
		{
			name:     "reject_inactive_kyc",
			feePayer: creator,
			msgs:     []sdk.Msg{buildJoinMsg(creator.Address, creator.Address, 1)},
			setup: func(_ *mock.MockDaoKeeper, gk *mock.MockMeGroupKeeper) {
				gk.EXPECT().GetGroupInfo(gomock.Any(), uint64(1)).Return(group, true)
				gk.EXPECT().GetDidAndKycActive(gomock.Any(), creator.GetAddress(), regionID).Return("", false)
			},
			expectErr: true,
			errText:   "does not have active KYC",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dk := mock.NewMockDaoKeeper(ctrl)
			gk := mock.NewMockMeGroupKeeper(ctrl)
			d := newDeductFeeDecorator(ctrl, dk, gk)

			dk.EXPECT().IsDao(gomock.Any(), tc.feePayer.Address).Return(false)
			dk.EXPECT().CheckFreeGasAccount(gomock.Any(), tc.feePayer.Address).Return(false)

			tc.setup(dk, gk)

			tx := &mock.MockFeeTx{
				Msgs:      tc.msgs,
				FeeAmount: sdk.Coins{},
				GasLimit:  200_000,
				Payer:     tc.feePayer.GetAddress(),
			}

			_, err := d.AnteHandle(sdk.Context{}, tx, false, passThrough)
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errText)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
