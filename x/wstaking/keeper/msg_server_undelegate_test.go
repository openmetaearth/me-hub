package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

// TestUndelegate tests the Undelegate message handler
func (s *KeeperTestSuite) TestUndelegate() {
	testCases := []struct {
		name      string
		setup     func() *stakingtypes.MsgUndelegate
		expectErr bool
		errMsg    string
	}{
		{
			name: "successful undelegation",
			setup: func() *stakingtypes.MsgUndelegate {
				account := s.TestAccs[0]
				// Initialize KYC for the account
				s.InitKyc(account, "did:metaearth:test1", types.ExperienceRegionName)
				// IMPORTANT: Initialize KYC for GlobalDao too (it's the actual delegator)
				daoAddr := sdk.MustAccAddressFromBech32(s.Dao.GlobalDao)
				s.InitKyc(daoAddr, "did:metaearth:globaldao", types.ExperienceRegionName)

				// Use GlobalDao for staking (MsgStake requires DAO permission)
				delegateMsg := &types.MsgStake{
					StakerAddress:    s.Dao.GlobalDao,
					ValidatorAddress: s.experienceValidator.OperatorAddress,
					Amount:           sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100000000)),
				}
				_, err := s.msgServer.Stake(s.Ctx, delegateMsg)
				s.Require().NoError(err)

				// Set up region with sufficient interest pool
				region, found := s.App.StakingKeeper.GetRegion(s.Ctx, types.ExperienceRegionName)
				if !found {
					// Create region if it doesn't exist
					region = types.Region{
						RegionId:           types.ExperienceRegionName,
						OperatorAddress:    s.experienceValidator.OperatorAddress,
						DelegateInterest:   sdkmath.LegacyNewDec(100000),
						DelegateAmount:     sdkmath.ZeroInt(),
						RegionShare:        sdkmath.ZeroInt(),
						RegionTreasureAddr: s.Dao.GlobalDao,
					}
				} else {
					region.DelegateInterest = sdkmath.LegacyNewDec(100000) // 100k umec interest pool
				}
				s.App.StakingKeeper.SetRegion(s.Ctx, region)
				// Fund region treasure account for reward payments
				regionTreasureAddr := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
				s.FundAcc(regionTreasureAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(10000000))))
				// Manually create delegation record for testing (normally done by KYC/other flows)
				// daoAddr already declared above
				s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.Delegation{
					DelegatorAddress: daoAddr.String(),
					ValidatorAddress: s.experienceValidator.OperatorAddress,
					Amount:           sdkmath.ZeroInt(),
					UnMeidAmount:     sdkmath.NewInt(100000000), // For Experience Region, use UnMeidAmount
					Unmovable:        sdkmath.ZeroInt(),
					StartHeight:      s.Ctx.BlockHeight(),
				})

				// Advance blocks to accumulate some rewards
				s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 100)
				s.Ctx = s.Ctx.WithBlockTime(s.Ctx.BlockTime().Add(time.Hour))

				return &stakingtypes.MsgUndelegate{
					DelegatorAddress: s.Dao.GlobalDao,
					ValidatorAddress: s.experienceValidator.OperatorAddress,
					Amount:           sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100000000)),
				}
			},
			expectErr: false,
		},
		{
			name: "undelegation without region",
			setup: func() *stakingtypes.MsgUndelegate {
				account := s.TestAccs[1]
				// No KYC initialization - no region assigned

				return &stakingtypes.MsgUndelegate{
					DelegatorAddress: account.String(),
					ValidatorAddress: s.experienceValidator.OperatorAddress,
					Amount:           sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100000000)),
				}
			},
			expectErr: true,
			errMsg:    "Validator DelegationAmount less than requested value", // Validator has no delegation amount set
		},
		{
			name: "undelegation with invalid denom",
			setup: func() *stakingtypes.MsgUndelegate {
				account := s.TestAccs[2]
				s.InitKyc(account, "did:metaearth:test-denom", types.ExperienceRegionName)

				return &stakingtypes.MsgUndelegate{
					DelegatorAddress: account.String(),
					ValidatorAddress: s.experienceValidator.OperatorAddress,
					Amount:           sdk.NewCoin("invalid", sdkmath.NewInt(100000)),
				}
			},
			expectErr: true,
			errMsg:    "invalid coin denomination", // Denom validation happens before region check
		},
		{
			name: "undelegation exceeding delegation",
			setup: func() *stakingtypes.MsgUndelegate {
				account := s.TestAccs[0]

				// Try to undelegate more than delegated
				return &stakingtypes.MsgUndelegate{
					DelegatorAddress: account.String(),
					ValidatorAddress: s.experienceValidator.OperatorAddress,
					Amount:           sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(900000000)),
				}
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			msg := tc.setup()
			resp, err := s.msgServer.Undelegate(s.Ctx, msg)

			if tc.expectErr {
				s.Require().Error(err)
				if tc.errMsg != "" {
					s.Require().Contains(err.Error(), tc.errMsg)
				}
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().NotNil(resp.CompletionTime)
				s.Require().True(resp.CompletionTime.After(s.Ctx.BlockTime()))

				// Verify events were emitted
				events := s.Ctx.EventManager().Events()
				foundUndelegateEvent := false
				for _, event := range events {
					if event.Type == types.EventTypeUnDelegate {
						foundUndelegateEvent = true
						break
					}
				}
				s.Require().True(foundUndelegateEvent, "Undelegate event should be emitted")
			}
		})
	}
}

// TestUndelegateWithRewards tests undelegation with accumulated rewards
func (s *KeeperTestSuite) TestUndelegateWithRewards() {
	account := s.TestAccs[0]
	s.InitKyc(account, "did:metaearth:reward-test", types.ExperienceRegionName)
	// IMPORTANT: Initialize KYC for GlobalDao too
	daoAddr := sdk.MustAccAddressFromBech32(s.Dao.GlobalDao)
	s.InitKyc(daoAddr, "did:metaearth:globaldao-reward", types.ExperienceRegionName)

	// Use GlobalDao for staking (MsgStake requires DAO permission)
	delegateAmount := sdkmath.NewInt(100000000)
	delegateMsg := &types.MsgStake{
		StakerAddress:    s.Dao.GlobalDao,
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, delegateAmount),
	}
	_, err := s.msgServer.Stake(s.Ctx, delegateMsg)
	s.Require().NoError(err)

	// Set up region with sufficient interest pool AFTER stake
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, types.ExperienceRegionName)
	if !found {
		region = types.Region{
			RegionId:           types.ExperienceRegionName,
			OperatorAddress:    s.experienceValidator.OperatorAddress,
			DelegateInterest:   sdkmath.LegacyNewDec(100000),
			DelegateAmount:     sdkmath.ZeroInt(),
			RegionShare:        sdkmath.ZeroInt(),
			RegionTreasureAddr: s.Dao.GlobalDao,
		}
	} else {
		region.DelegateInterest = sdkmath.LegacyNewDec(100000)
	}
	s.App.StakingKeeper.SetRegion(s.Ctx, region)

	// Fund region treasure account for reward payments
	regionTreasureAddr := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
	s.FundAcc(regionTreasureAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(10000000))))

	// Manually create delegation record for testing
	// daoAddr already declared above
	s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.Delegation{
		DelegatorAddress: daoAddr.String(),
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdkmath.ZeroInt(),
		UnMeidAmount:     sdkmath.NewInt(100000000), // For Experience Region
		Unmovable:        sdkmath.ZeroInt(),
		StartHeight:      s.Ctx.BlockHeight(),
	})

	// Get initial balance of GlobalDao
	balanceBefore := s.App.BankKeeper.GetBalance(s.Ctx, daoAddr, params.BaseDenom)

	// Advance blocks significantly to accumulate rewards
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1000)
	s.Ctx = s.Ctx.WithBlockTime(s.Ctx.BlockTime().Add(24 * time.Hour))

	// Undelegate
	undelegateAmount := sdkmath.NewInt(50000000)
	undelegateMsg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: s.Dao.GlobalDao,
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, undelegateAmount),
	}
	resp, err := s.msgServer.Undelegate(s.Ctx, undelegateMsg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	// Balance should increase by rewards (undelegated amount is locked until completion time)
	balanceAfter := s.App.BankKeeper.GetBalance(s.Ctx, daoAddr, params.BaseDenom)
	s.Require().True(balanceAfter.Amount.GTE(balanceBefore.Amount), "balance should not decrease")
}

// TestUndelegateMultipleTimes tests multiple undelegations from the same delegation
func (s *KeeperTestSuite) TestUndelegateMultipleTimes() {
	account := s.TestAccs[1] // Use a different account
	s.InitKyc(account, "did:metaearth:multi-undelegate", types.ExperienceRegionName)
	// IMPORTANT: Initialize KYC for GlobalDao too
	daoAddr := sdk.MustAccAddressFromBech32(s.Dao.GlobalDao)
	s.InitKyc(daoAddr, "did:metaearth:globaldao-multi", types.ExperienceRegionName)

	// Use GlobalDao for staking
	delegateMsg := &types.MsgStake{
		StakerAddress:    s.Dao.GlobalDao,
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(300000000)),
	}
	_, err := s.msgServer.Stake(s.Ctx, delegateMsg)
	s.Require().NoError(err)

	// Set up region with sufficient interest pool AFTER stake
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, types.ExperienceRegionName)
	if !found {
		region = types.Region{
			RegionId:           types.ExperienceRegionName,
			OperatorAddress:    s.experienceValidator.OperatorAddress,
			DelegateInterest:   sdkmath.LegacyNewDec(100000),
			DelegateAmount:     sdkmath.ZeroInt(),
			RegionShare:        sdkmath.ZeroInt(),
			RegionTreasureAddr: s.Dao.GlobalDao,
		}
	} else {
		region.DelegateInterest = sdkmath.LegacyNewDec(100000)
	}
	s.App.StakingKeeper.SetRegion(s.Ctx, region)

	// Fund region treasure account for reward payments
	regionTreasureAddr := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
	s.FundAcc(regionTreasureAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(10000000))))

	// Set validator's DelegationAmount to match the delegation
	valAddr, _ := sdk.ValAddressFromBech32(s.experienceValidator.OperatorAddress)
	val, err := s.App.StakingKeeper.GetValidator(s.Ctx, valAddr)
	s.Require().NoError(err)
	val.DelegationAmount = sdkmath.NewInt(300000000)
	s.App.StakingKeeper.SetValidator(s.Ctx, val)

	// Manually create delegation record for testing
	// daoAddr already declared above
	s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.Delegation{
		DelegatorAddress: daoAddr.String(),
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdkmath.ZeroInt(),
		UnMeidAmount:     sdkmath.NewInt(300000000), // For Experience Region
		Unmovable:        sdkmath.ZeroInt(),
		StartHeight:      s.Ctx.BlockHeight(),
	})

	// First undelegation
	msg1 := &stakingtypes.MsgUndelegate{
		DelegatorAddress: s.Dao.GlobalDao,
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100000000)),
	}
	resp1, err := s.msgServer.Undelegate(s.Ctx, msg1)
	s.Require().NoError(err)
	s.Require().NotNil(resp1)

	// Advance time
	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 100)
	s.Ctx = s.Ctx.WithBlockTime(s.Ctx.BlockTime().Add(time.Hour))

	// Second undelegation
	msg2 := &stakingtypes.MsgUndelegate{
		DelegatorAddress: s.Dao.GlobalDao,
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100000000)),
	}
	resp2, err := s.msgServer.Undelegate(s.Ctx, msg2)
	s.Require().NoError(err)
	s.Require().NotNil(resp2)

	// Both should succeed
	s.Require().True(resp2.CompletionTime.After(resp1.CompletionTime) || resp2.CompletionTime.Equal(resp1.CompletionTime))
}

// TestUndelegateEvents tests that correct events are emitted
func (s *KeeperTestSuite) TestUndelegateEvents() {
	account := s.TestAccs[2] // Use a different account
	s.InitKyc(account, "did:metaearth:event-test", types.ExperienceRegionName)
	// IMPORTANT: Initialize KYC for GlobalDao too
	daoAddr := sdk.MustAccAddressFromBech32(s.Dao.GlobalDao)
	s.InitKyc(daoAddr, "did:metaearth:globaldao-event", types.ExperienceRegionName)

	// Use GlobalDao for staking
	delegateMsg := &types.MsgStake{
		StakerAddress:    s.Dao.GlobalDao,
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100000000)),
	}
	_, err := s.msgServer.Stake(s.Ctx, delegateMsg)
	s.Require().NoError(err)

	// Set up region with sufficient interest pool AFTER stake
	region, found := s.App.StakingKeeper.GetRegion(s.Ctx, types.ExperienceRegionName)
	if !found {
		region = types.Region{
			RegionId:           types.ExperienceRegionName,
			OperatorAddress:    s.experienceValidator.OperatorAddress,
			DelegateInterest:   sdkmath.LegacyNewDec(100000),
			DelegateAmount:     sdkmath.ZeroInt(),
			RegionShare:        sdkmath.ZeroInt(),
			RegionTreasureAddr: s.Dao.GlobalDao,
		}
	} else {
		region.DelegateInterest = sdkmath.LegacyNewDec(100000)
	}
	s.App.StakingKeeper.SetRegion(s.Ctx, region)

	// Fund region treasure account for reward payments
	regionTreasureAddr := sdk.MustAccAddressFromBech32(region.RegionTreasureAddr)
	s.FundAcc(regionTreasureAddr, sdk.NewCoins(sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(10000000))))

	// Manually create delegation record for testing
	// daoAddr already declared above
	s.App.StakingKeeper.SetDelegation(s.Ctx, stakingtypes.Delegation{
		DelegatorAddress: daoAddr.String(),
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdkmath.ZeroInt(),
		UnMeidAmount:     sdkmath.NewInt(100000000), // For Experience Region
		Unmovable:        sdkmath.ZeroInt(),
		StartHeight:      s.Ctx.BlockHeight(),
	})

	// Clear events
	s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())

	// Undelegate
	undelegateMsg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: s.Dao.GlobalDao,
		ValidatorAddress: s.experienceValidator.OperatorAddress,
		Amount:           sdk.NewCoin(params.BaseDenom, sdkmath.NewInt(100000000)),
	}
	_, err = s.msgServer.Undelegate(s.Ctx, undelegateMsg)
	s.Require().NoError(err)

	// Check events
	events := s.Ctx.EventManager().Events()
	s.Require().NotEmpty(events)

	var undelegateEvent sdk.Event
	for _, event := range events {
		if event.Type == types.EventTypeUnDelegate {
			undelegateEvent = event
			break
		}
	}

	s.Require().NotNil(undelegateEvent.Type)
	s.Require().Equal(types.EventTypeUnDelegate, undelegateEvent.Type)

	// Verify event attributes
	attrs := undelegateEvent.Attributes
	s.Require().NotEmpty(attrs)

	// Check for required attributes
	hasValidator := false
	hasAmount := false
	hasDelegator := false
	hasCompletionTime := false

	for _, attr := range attrs {
		switch attr.Key {
		case stakingtypes.AttributeKeyValidator:
			hasValidator = true
		case sdk.AttributeKeyAmount:
			hasAmount = true
		case types.AttributeKeyDelegatorAddress:
			hasDelegator = true
		case stakingtypes.AttributeKeyCompletionTime:
			hasCompletionTime = true
		}
	}

	s.Require().True(hasValidator, "should have validator attribute")
	s.Require().True(hasAmount, "should have amount attribute")
	s.Require().True(hasDelegator, "should have delegator attribute")
	s.Require().True(hasCompletionTime, "should have completion time attribute")
}
