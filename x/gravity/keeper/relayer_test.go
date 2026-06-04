package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestGravityAndBridger() {
	for _, relayer := range suite.relayerAddrs {
		require.True(suite.T(), suite.Keeper().IsProposalRelayer(suite.Ctx, relayer.String()))
	}
}

func (suite *KeeperTestSuite) TestOfflineRelayersDoNotConcentrateAttestationPower() {
	k := suite.Keeper()
	ctx := suite.Ctx

	// Setup 3 relayers: A, B, C with equal delegated stake
	relayers := suite.relayerAddrs[:3]
	
	for _, relayerAddr := range relayers {
		relayer := types.Relayer{
			Online:         true,
			DelegateAmount: sdk.DefaultPowerReduction.MulRaw(1000), // 1000 power
			SlashTimes:     0,
			RelayerAddress: relayerAddr.String(),
		}
		k.SetRelayer(ctx, relayerAddr, relayer)
	}

	// Delete all other relayers to clean state
	for _, relayerAddr := range suite.relayerAddrs[3:] {
		k.DelRelayer(ctx, relayerAddr)
	}

	// Initially, all 3 are online.
	k.SetLastTotalPower(ctx)
	// total online power should be 3000 (1000 * 3)
	require.Equal(suite.T(), sdkmath.NewInt(3000), k.GetLastTotalPower(ctx))

	// Now Relayer B and C go offline
	relayerB, _ := k.GetRelayer(ctx, relayers[1])
	relayerB.Online = false
	k.SetRelayer(ctx, relayers[1], relayerB)

	relayerC, _ := k.GetRelayer(ctx, relayers[2])
	relayerC.Online = false
	k.SetRelayer(ctx, relayers[2], relayerC)

	// Refresh total power
	k.SetLastTotalPower(ctx)

	// Since only Relayer A is online, its power (1000) represents 100% of online power,
	// which is over the concentration threshold (33.34%).
	// The total power should be overridden to total bonded power (3000) instead of just online power (1000).
	require.Equal(suite.T(), sdkmath.NewInt(3000), k.GetLastTotalPower(ctx))
}
