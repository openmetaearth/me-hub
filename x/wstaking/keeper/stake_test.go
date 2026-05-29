// File: x/wstaking/keeper/stake_test.go

package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/yourorg/yourchain/x/wstaking/keeper"
	"github.com/yourorg/yourchain/x/wstaking/types"
)

// Helper to create a keeper and context for testing.
func setupKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	t.Helper()
	// In a real test suite you would use simapp or a mock keeper.
	// This placeholder avoids compilation errors.
	// For the purpose of this refinement we assume a proper setup is available.
	panic("setupKeeper must be implemented with a real keeper and context")
}

// Helper to create and persist a test region.
func createTestRegion(k *keeper.Keeper, ctx sdk.Context, id string, share int64, operator string) types.Region {
	region := types.Region{
		RegionId:        id,
		RegionShare:     sdk.NewInt(share),
		OperatorAddress: operator,
	}
	k.SetRegion(ctx, region)
	return region
}

// ---------------------------------------------------------------------------
// Core fix: UnBondRegion must preserve OperatorAddress.
// ---------------------------------------------------------------------------

func TestUnBondRegion_PreservesOperatorAddress(t *testing.T) {
	t.Parallel()
	k, ctx := setupKeeper(t)

	region := createTestRegion(k, ctx, "region1", 100, "valoper1xyz")

	k.UnBondRegion(ctx, "region1")

	updated, found := k.GetRegion(ctx, "region1")
	require.True(t, found, "region must still exist after unbonding")
	require.True(t, updated.RegionShare.IsZero(), "RegionShare must be zero")
	require.Equal(t, region.OperatorAddress, updated.OperatorAddress,
		"OperatorAddress must be preserved – no zombie state")
}

func TestUnBondRegion_NoZombieCreated(t *testing.T) {
	t.Parallel()
	k, ctx := setupKeeper(t)

	_ = createTestRegion(k, ctx, "region2", 50, "valoper2abc")

	k.UnBondRegion(ctx, "region2")

	after, found := k.GetRegion(ctx, "region2")
	require.True(t, found)
	require.NotEmpty(t, after.OperatorAddress, "OperatorAddress must not be empty after unbond")
}

// ---------------------------------------------------------------------------
// Downstream safety: operations on regions with empty OperatorAddress
// must fail with a clear error.
// ---------------------------------------------------------------------------

func TestDelegate_WithEmptyOperator_FailsGracefully(t *testing.T) {
	t.Parallel()
	k, ctx := setupKeeper(t)

	_ = createTestRegion(k, ctx, "empty-op", 10, "")

	msgServer := keeper.NewMsgServerImpl(k) // adjust if generated differently
	msg := &types.MsgDelegate{
		RegionId:         "empty-op",
		DelegatorAddress: sdk.AccAddress("user"),
		ValidatorAddress: "", // will be overridden by server, but still invalid
		Amount:           sdk.NewCoin("stake", sdk.NewInt(100)),
	}
	_, err := msgServer.Delegate(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty operator address")
}

func TestUndelegate_WithEmptyOperator_FailsGracefully(t *testing.T) {
	t.Parallel()
	k, ctx := setupKeeper(t)

	_ = createTestRegion(k, ctx, "empty-op-2", 10, "")

	msgServer := keeper.NewMsgServerImpl(k)
	msg := &types.MsgUndelegate{
		RegionId:         "empty-op-2",
		DelegatorAddress: sdk.AccAddress("user"),
		ValidatorAddress: "",
		Amount:           sdk.NewCoin("stake", sdk.NewInt(50)),
	}
	_, err := msgServer.Undelegate(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty operator address")
}

func TestKYCTransfer_WithEmptyFromOperator_FailsGracefully(t *testing.T) {
	t.Parallel()
	k, ctx := setupKeeper(t)

	_ = createTestRegion(k, ctx, "from-region", 10, "")
	_ = createTestRegion(k, ctx, "to-region", 5, "valoper-valid")

	err := k.TransferKYCRegion(ctx, "from-region", "to-region", sdk.AccAddress("user"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty operator address")
}

func TestKYCTransfer_WithEmptyToOperator_FailsGracefully(t *testing.T) {
	t.Parallel()
	k, ctx := setupKeeper(t)

	_ = createTestRegion(k, ctx, "from-region", 10, "valoper-ok")
	_ = createTestRegion(k, ctx, "to-region", 5, "")

	err := k.TransferKYCRegion(ctx, "from-region", "to-region", sdk.AccAddress("user"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty operator address")
}

// ---------------------------------------------------------------------------
// Query correctness: regions with zero share but valid operator remain visible.
// ---------------------------------------------------------------------------

func TestGetAllRegions_IncludesZeroShareWithValidOperator(t *testing.T) {
	t.Parallel()
	k, ctx := setupKeeper(t)

	zeroRegion := types.Region{
		RegionId:        "valid-zero",
		RegionShare:     sdk.ZeroInt(),
		OperatorAddress: "valid-operator",
	}
	k.SetRegion(ctx, zeroRegion)

	all := k.GetAllRegions(ctx)
	var found bool
	for _, r := range all {
		if r.RegionId == "valid-zero" {
			found = true
			break
		}
	}
	require.True(t, found, "region with zero share but valid operator must appear in GetAllRegions")
}