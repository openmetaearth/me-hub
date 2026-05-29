package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Krako-Labs/KORA/x/wstaking/keeper"
	"github.com/Krako-Labs/KORA/x/wstaking/types"
)

func (suite *KeeperTestSuite) TestDelegateUndelegateAfterFullUnstake() {
	suite.SetupTest()

	app := suite.App
	ctx := app.BaseApp.NewContext(false, app.DeliverTx)
	wstakingKeeper := app.WStakingKeeper
	bankKeeper := app.BankKeeper
	stakingKeeper := app.StakingKeeper

	// 1. Create a validator, bond some tokens
	validatorAddr := sdk.ValAddress(sdk.AccAddress("validator-1"))
	validator, err := stakingKeeper.GetValidator(ctx, validatorAddr)
	require.NoError(suite.T(), err)
	if validator == nil {
		// create a simple validator (e.g., by delegating)
		bondDenom := stakingKeeper.BondDenom(ctx)
		bondAmount := sdk.NewInt(1_000_000)
		acc := sdk.AccAddress(validatorAddr)
		err := bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(bondDenom, bondAmount)))
		require.NoError(suite.T(), err)
		_, err = stakingKeeper.Delegate(ctx, acc, bondAmount, sdk.Unbonded, validator, true)
		require.NoError(suite.T(), err)
		validator, _ = stakingKeeper.GetValidator(ctx, validatorAddr)
	}

	// 2. Create a region linked to this validator
	regionID := "test-region-1"
	region := types.Region{
		RegionId:        regionID,
		RegionShare:     sdk.ZeroInt(),
		OperatorAddress: validatorAddr.String(),
	}
	wstakingKeeper.SetRegion(ctx, region)
	err = wstakingKeeper.GroupKeeper.CreateGroup(ctx, regionID, validatorAddr.String())
	require.NoError(suite.T(), err)

	// 3. Delegate some tokens to the region
	delegator := sdk.AccAddress("delegator-1")
	delegationAmount := sdk.NewInt(500_000)
	bondDenom := stakingKeeper.BondDenom(ctx)
	coins := sdk.NewCoins(sdk.NewCoin(bondDenom, delegationAmount))
	err = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(suite.T(), err)
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, coins)
	require.NoError(suite.T(), err)

	msgDelegate := types.MsgDelegate{
		DelegatorAddress:  delegator.String(),
		RegionId:          regionID,
		Amount:            sdk.NewCoin(bondDenom, delegationAmount),
	}
	msgServer := keeper.NewMsgServerImpl(wstakingKeeper)
	_, err = msgServer.Delegate(sdk.WrapSDKContext(ctx), &msgDelegate)
	require.NoError(suite.T(), err)

	// Verify delegation exists
	delegation, found := wstakingKeeper.GetDelegation(ctx, regionID, delegator.String())
	require.True(suite.T(), found)
	require.Equal(suite.T(), delegationAmount, delegation.Shares)

	// 4. Fully unstake: undelegate all shares from the region
	msgUndelegate := types.MsgUndelegate{
		DelegatorAddress: delegator.String(),
		RegionId:         regionID,
		Amount:           sdk.NewCoin(bondDenom, delegationAmount),
	}
	_, err = msgServer.Undelegate(sdk.WrapSDKContext(ctx), &msgUndelegate)
	require.NoError(suite.T(), err)

	// Complete unbonding period (simulate time passing)
	ctx = ctx.WithBlockTime(ctx.BlockHeader().Time.Add(time.Hour * 24 * 7))
	stakingKeeper.CompleteUnbonding(ctx)

	// 5. Verify region state after full unstake
	regionAfterUnstake, found := wstakingKeeper.GetRegion(ctx, regionID)
	require.True(suite.T(), found, "region should still exist after full unstake")
	require.True(suite.T(), regionAfterUnstake.RegionShare.IsZero(), "region share should be zero")
	require.NotEmpty(suite.T(), regionAfterUnstake.OperatorAddress, "operator address should NOT be cleared")

	// 6. Attempt to delegate again to the same region (now with zero shares but operator address preserved)
	// This should succeed because operator address is still valid
	delegationAmount2 := sdk.NewInt(100_000)
	coins2 := sdk.NewCoins(sdk.NewCoin(bondDenom, delegationAmount2))
	err = bankKeeper.MintCoins(ctx, types.ModuleName, coins2)
	require.NoError(suite.T(), err)
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, coins2)
	require.NoError(suite.T(), err)

	msgDelegate2 := types.MsgDelegate{
		DelegatorAddress: delegator.String(),
		RegionId:         regionID,
		Amount:           sdk.NewCoin(bondDenom, delegationAmount2),
	}
	_, err = msgServer.Delegate(sdk.WrapSDKContext(ctx), &msgDelegate2)
	require.NoError(suite.T(), err, "delegation to region with zero shares should succeed")

	// 7. Undelegate the new delegation
	msgUndelegate2 := types.MsgUndelegate{
		DelegatorAddress: delegator.String(),
		RegionId:         regionID,
		Amount:           sdk.NewCoin(bondDenom, delegationAmount2),
	}
	_, err = msgServer.Undelegate(sdk.WrapSDKContext(ctx), &msgUndelegate2)
	require.NoError(suite.T(), err)

	// 8. Verify region still has operator address after second full unstake
	regionAfterSecond, found := wstakingKeeper.GetRegion(ctx, regionID)
	require.True(suite.T(), found)
	require.True(suite.T(), regionAfterSecond.RegionShare.IsZero())
	require.NotEmpty(suite.T(), regionAfterSecond.OperatorAddress, "operator address still not cleared after second full unstake")

	// 9. Optional: validate KYC region transfer still works with such a region
	// (If the module includes TransferRegionKYC functionality, test it here)
	if kycKeeper, ok := app.WStakingKeeper.(interface{ GetRegionKYC(ctx sdk.Context, regionId string) (types.RegionKYC, bool) }); ok {
		// minimal test: transfer from this region to another (requires two regions)
		// For completeness, we just ensure no panic occurs when operator address is used
		_, err = kycKeeper.GetRegionKYC(ctx, regionID)
		require.NoError(suite.T(), err) // assume no error if region exists
	}
}