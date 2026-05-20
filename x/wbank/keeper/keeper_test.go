package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/openmetaearth/me-hub/app/apptesting"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wbank/keeper"
	"github.com/openmetaearth/me-hub/x/wbank/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	bankKeeperWrapper keeper.BaseKeeperWrapper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	app := apptesting.Setup(s.T())
	ctx := app.GetBaseApp().NewContext(false)
	s.App = app
	s.Ctx = ctx

	// Create the wrapper
	s.bankKeeperWrapper = keeper.NewBankKeeperWrapper(
		app.BankKeeper,
		app.AccountKeeper,
		app.DaoKeeper,
	)
}

// TestStakeCoinsFromModuleToModule tests staking coins transfer between module accounts
func (s *KeeperTestSuite) TestStakeCoinsFromModuleToModule() {
	ctx := s.Ctx

	// Create test modules
	senderModule := "staking"
	recipientModule := "distribution"

	testCases := []struct {
		name         string
		setupFunc    func()
		senderMod    string
		recipientMod string
		amount       sdk.Coins
		expectError  bool
		errorMsg     string
	}{
		{
			name: "successful stake transfer",
			setupFunc: func() {
				// Fund sender module account

				coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1000000))
				err := s.App.BankKeeper.MintCoins(ctx, senderModule, coins)
				require.NoError(s.T(), err)
			},
			senderMod:    senderModule,
			recipientMod: recipientModule,
			amount:       sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100000)),
			expectError:  false,
		},
		{
			name:         "sender module does not exist",
			setupFunc:    func() {},
			senderMod:    "nonexistent",
			recipientMod: recipientModule,
			amount:       sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100000)),
			expectError:  true,
			errorMsg:     "module account nonexistent does not exist",
		},
		{
			name: "recipient module does not exist",
			setupFunc: func() {

				coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1000000))
				err := s.App.BankKeeper.MintCoins(ctx, senderModule, coins)
				require.NoError(s.T(), err)
			},
			senderMod:    senderModule,
			recipientMod: "nonexistent",
			amount:       sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100000)),
			expectError:  true,
			errorMsg:     "module account nonexistent does not exist",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Reset state for each test
			s.SetupTest()
			ctx := s.Ctx

			tc.setupFunc()

			if tc.expectError {
				require.Panics(s.T(), func() {
					_ = s.bankKeeperWrapper.StakeCoinsFromModuleToModule(ctx, tc.senderMod, tc.recipientMod, tc.amount)
				}, tc.errorMsg)
			} else {
				// Get initial balances
				recipientAddr := s.App.AccountKeeper.GetModuleAddress(tc.recipientMod)
				initialBalance := s.App.BankKeeper.GetBalance(ctx, recipientAddr, params.BaseDenom)

				// Execute transfer
				err := s.bankKeeperWrapper.StakeCoinsFromModuleToModule(ctx, tc.senderMod, tc.recipientMod, tc.amount)
				require.NoError(s.T(), err)

				// Verify recipient balance increased
				finalBalance := s.App.BankKeeper.GetBalance(ctx, recipientAddr, params.BaseDenom)
				expectedBalance := initialBalance.Add(tc.amount[0])
				require.Equal(s.T(), expectedBalance, finalBalance)
			}
		})
	}
}

// TestUnstakeCoinsFromModuleToModule tests unstaking coins transfer between module accounts
func (s *KeeperTestSuite) TestUnstakeCoinsFromModuleToModule() {
	ctx := s.Ctx

	senderModule := "staking"
	recipientModule := "distribution"

	// Fund sender module account using fee_collector (which has mint permission)
	coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1000000))
	err := s.App.BankKeeper.MintCoins(ctx, authtypes.FeeCollectorName, coins)
	require.NoError(s.T(), err)
	err = s.App.BankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, senderModule, coins)
	require.NoError(s.T(), err)

	amount := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100000))
	recipientAddr := s.App.AccountKeeper.GetModuleAddress(recipientModule)
	initialBalance := s.App.BankKeeper.GetBalance(ctx, recipientAddr, params.BaseDenom)

	err = s.bankKeeperWrapper.UnstakeCoinsFromModuleToModule(ctx, senderModule, recipientModule, amount)
	require.NoError(s.T(), err)

	finalBalance := s.App.BankKeeper.GetBalance(ctx, recipientAddr, params.BaseDenom)
	expectedBalance := initialBalance.Add(amount[0])
	require.Equal(s.T(), expectedBalance, finalBalance)
}

// TestFeeToReceivers tests fee distribution to multiple receivers
func (s *KeeperTestSuite) TestFeeToReceivers() {
	ctx := s.Ctx

	// Create test accounts
	sender := s.TestAccs()[0]
	receiver1 := s.TestAccs()[1]
	receiver2 := s.TestAccs()[2]

	testCases := []struct {
		name          string
		setupFunc     func()
		inputs        []banktypes.Input
		outputs       []banktypes.Output
		receiverTypes []types.FeeReceiverType
		expectError   bool
		errorMsg      string
	}{
		{
			name: "successful fee distribution",
			setupFunc: func() {
				// Fund sender account
				coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1000000))
				err := s.App.BankKeeper.MintCoins(ctx, authtypes.FeeCollectorName, coins)
				require.NoError(s.T(), err)
				err = s.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, sender, coins)
				require.NoError(s.T(), err)
			},
			inputs: []banktypes.Input{
				{
					Address: sender.String(),
					Coins:   sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100000)),
				},
			},
			outputs: []banktypes.Output{
				{
					Address: receiver1.String(),
					Coins:   sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 60000)),
				},
				{
					Address: receiver2.String(),
					Coins:   sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 40000)),
				},
			},
			receiverTypes: []types.FeeReceiverType{
				types.FeeReceiverGlobalDaoFeePool,
				types.FeeReceiverDevOperator,
			},
			expectError: false,
		},
		{
			name:      "empty inputs error",
			setupFunc: func() {},
			inputs:    []banktypes.Input{},
			outputs: []banktypes.Output{
				{
					Address: receiver1.String(),
					Coins:   sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100)),
				},
			},
			receiverTypes: []types.FeeReceiverType{types.FeeReceiverGlobalDaoFeePool},
			expectError:   true,
			errorMsg:      "inputs error",
		},
		{
			name: "mismatched receiver types and outputs",
			setupFunc: func() {
				coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1000000))
				err := s.App.BankKeeper.MintCoins(ctx, authtypes.FeeCollectorName, coins)
				require.NoError(s.T(), err)
				err = s.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, sender, coins)
				require.NoError(s.T(), err)
			},
			inputs: []banktypes.Input{
				{
					Address: sender.String(),
					Coins:   sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100)),
				},
			},
			outputs: []banktypes.Output{
				{
					Address: receiver1.String(),
					Coins:   sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 100)),
				},
			},
			receiverTypes: []types.FeeReceiverType{
				types.FeeReceiverGlobalDaoFeePool,
				types.FeeReceiverDevOperator,
			},
			expectError: true,
			errorMsg:    "fee receiver types and outputs are not equal",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			ctx := s.Ctx

			tc.setupFunc()

			err := s.bankKeeperWrapper.FeeToReceivers(ctx, tc.inputs, tc.outputs, tc.receiverTypes)

			if tc.expectError {
				require.Error(s.T(), err)
				if tc.errorMsg != "" {
					require.Contains(s.T(), err.Error(), tc.errorMsg)
				}
			} else {
				require.NoError(s.T(), err)

				// Verify event was emitted
				events := ctx.EventManager().Events()
				foundEvent := false
				for _, event := range events {
					if event.Type == types.EventTypeFeeToReceivers {
						foundEvent = true
						break
					}
				}
				require.True(s.T(), foundEvent, "FeeToReceivers event should be emitted")
			}
		})
	}
}

// TestSendCoinsWithTag tests sending coins with custom tags
func (s *KeeperTestSuite) TestSendCoinsWithTag() {
	ctx := s.Ctx

	sender := s.TestAccs()[0]
	recipient := s.TestAccs()[1]

	testCases := []struct {
		name      string
		setupFunc func()
		tags      []string
		amount    sdk.Coins
	}{
		{
			name: "send coins with single tag",
			setupFunc: func() {
				coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1000000))
				err := s.App.BankKeeper.MintCoins(ctx, authtypes.FeeCollectorName, coins)
				require.NoError(s.T(), err)
				err = s.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, sender, coins)
				require.NoError(s.T(), err)
			},
			tags:   []string{"test_tag"},
			amount: sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 10000)),
		},
		{
			name: "send coins with multiple tags",
			setupFunc: func() {
				coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1000000))
				err := s.App.BankKeeper.MintCoins(ctx, authtypes.FeeCollectorName, coins)
				require.NoError(s.T(), err)
				err = s.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, sender, coins)
				require.NoError(s.T(), err)
			},
			tags:   []string{"tag1", "tag2", "tag3"},
			amount: sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 10000)),
		},
		{
			name: "send coins without tags",
			setupFunc: func() {
				coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1000000))
				err := s.App.BankKeeper.MintCoins(ctx, authtypes.FeeCollectorName, coins)
				require.NoError(s.T(), err)
				err = s.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, sender, coins)
				require.NoError(s.T(), err)
			},
			tags:   []string{},
			amount: sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 10000)),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			ctx := s.Ctx
			extendedKeeper := s.bankKeeperWrapper.Extend()

			tc.setupFunc()

			initialBalance := s.App.BankKeeper.GetBalance(ctx, recipient, params.BaseDenom)

			err := extendedKeeper.SendCoinsWithTag(ctx, sender, recipient, tc.amount, tc.tags...)
			require.NoError(s.T(), err)

			// Verify balance transferred
			finalBalance := s.App.BankKeeper.GetBalance(ctx, recipient, params.BaseDenom)
			expectedBalance := initialBalance.Add(tc.amount[0])
			require.Equal(s.T(), expectedBalance, finalBalance)

			// Verify tag attributes in transfer event if tags provided
			if len(tc.tags) > 0 {
				events := ctx.EventManager().Events()
				foundTransfer := false
				for _, event := range events {
					if event.Type == banktypes.EventTypeTransfer {
						foundTransfer = true
						// Count tag attributes
						tagCount := 0
						for _, attr := range event.Attributes {
							if attr.Key == "tag" {
								tagCount++
							}
						}
						require.Equal(s.T(), len(tc.tags), tagCount, "should have correct number of tags")
						break
					}
				}
				require.True(s.T(), foundTransfer, "transfer event should exist")
			}
		})
	}
}

// TestSendCoinsFromModuleToAccountWithTag tests module to account transfer with tags
func (s *KeeperTestSuite) TestSendCoinsFromModuleToAccountWithTag() {
	ctx := s.Ctx
	extendedKeeper := s.bankKeeperWrapper.Extend()

	module := authtypes.FeeCollectorName
	recipient := s.TestAccs()[0]

	// Fund module account
	coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1000000))
	err := s.App.BankKeeper.MintCoins(ctx, module, coins)
	require.NoError(s.T(), err)

	amount := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 10000))
	initialBalance := s.App.BankKeeper.GetBalance(ctx, recipient, params.BaseDenom)

	err = extendedKeeper.SendCoinsFromModuleToAccountWithTag(ctx, module, recipient, amount, "module_transfer")
	require.NoError(s.T(), err)

	// Verify balance transferred
	finalBalance := s.App.BankKeeper.GetBalance(ctx, recipient, params.BaseDenom)
	expectedBalance := initialBalance.Add(amount[0])
	require.Equal(s.T(), expectedBalance, finalBalance)
}

// TestSendCoinsFromAccountToModuleWithTag tests account to module transfer with tags
func (s *KeeperTestSuite) TestSendCoinsFromAccountToModuleWithTag() {
	ctx := s.Ctx
	extendedKeeper := s.bankKeeperWrapper.Extend()

	sender := s.TestAccs()[0]
	module := authtypes.FeeCollectorName

	// Fund sender account
	coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1000000))
	err := s.App.BankKeeper.MintCoins(ctx, module, coins)
	require.NoError(s.T(), err)
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, module, sender, coins)
	require.NoError(s.T(), err)

	amount := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 10000))
	moduleAddr := s.App.AccountKeeper.GetModuleAddress(module)
	initialBalance := s.App.BankKeeper.GetBalance(ctx, moduleAddr, params.BaseDenom)

	err = extendedKeeper.SendCoinsFromAccountToModuleWithTag(ctx, sender, module, amount, "account_transfer")
	require.NoError(s.T(), err)

	// Verify balance transferred
	finalBalance := s.App.BankKeeper.GetBalance(ctx, moduleAddr, params.BaseDenom)
	expectedBalance := initialBalance.Add(amount[0])
	require.Equal(s.T(), expectedBalance, finalBalance)
}

// TestSendCoinsFromModuleToModuleWithTag tests module to module transfer with tags
func (s *KeeperTestSuite) TestSendCoinsFromModuleToModuleWithTag() {
	ctx := s.Ctx
	extendedKeeper := s.bankKeeperWrapper.Extend()

	senderModule := authtypes.FeeCollectorName
	recipientModule := "distribution"

	// Fund sender module account
	coins := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 1000000))
	err := s.App.BankKeeper.MintCoins(ctx, senderModule, coins)
	require.NoError(s.T(), err)

	amount := sdk.NewCoins(sdk.NewInt64Coin(params.BaseDenom, 10000))
	recipientAddr := s.App.AccountKeeper.GetModuleAddress(recipientModule)
	initialBalance := s.App.BankKeeper.GetBalance(ctx, recipientAddr, params.BaseDenom)

	err = extendedKeeper.SendCoinsFromModuleToModuleWithTag(ctx, senderModule, recipientModule, amount, "module_module_transfer")
	require.NoError(s.T(), err)

	// Verify balance transferred
	finalBalance := s.App.BankKeeper.GetBalance(ctx, recipientAddr, params.BaseDenom)
	expectedBalance := initialBalance.Add(amount[0])
	require.Equal(s.T(), expectedBalance, finalBalance)
}

func (s *KeeperTestSuite) TestAccs() []sdk.AccAddress {
	if len(s.KeeperTestHelper.TestAccs) == 0 {
		s.KeeperTestHelper.TestAccs = s.NewAccounts(3)
	}
	return s.KeeperTestHelper.TestAccs
}
