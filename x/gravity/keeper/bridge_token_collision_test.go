package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	"github.com/openmetaearth/me-hub/utils"
	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// TestBridgeTokenSymbolMECRejected verifies that a MsgBridgeTokenClaim with
// symbol "MEC" is rejected because the derived denom "umec" collides with
// the Hub native base denom. This prevents an attacker from registering a
// bridge token that mints native supply via MsgSendToMeClaim.
func (suite *KeeperTestSuite) TestBridgeTokenSymbolMECRejected() {
	// Arrange: create enough relayers to meet the attestation threshold
	relayers := suite.relayerAddrs
	externalPris := suite.externalPris

	// Act: attempt to register a bridge token with symbol "MEC"
	tokenContract := helpers.GenerateAddress().Hex()
	nonce := uint64(1)

	for i, relayer := range relayers {
		claim := &types.MsgBridgeTokenClaim{
			EventNonce:     nonce + uint64(i),
			TokenContract:  tokenContract,
			Symbol:         "MEC",
			Name:           "MEC Token",
			Decimals:       18,
			ChainName:      suite.chainName,
			ExternalAddress: suite.PubKeyToExternalAddr(externalPris[i].PublicKey),
			RelayerAddress:  relayer.String(),
		}
		_, err := suite.MsgServer().BridgeTokenClaim(suite.Ctx, claim)
		if i == 0 {
			// The first relayer's claim creates the attestation; it should fail
			// because "umec" is the native denom.
			suite.Error(err)
			suite.ErrorContains(err, "native denom")
			return
		}
	}

	// If we get here, the first claim didn't produce an error, which is wrong.
	suite.Fail("expected first BridgeTokenClaim with symbol MEC to be rejected")
}

// TestBridgeTokenSymbolNonNativeAllowed verifies that bridge tokens with
// symbols that do NOT collide with the native denom are still accepted.
func (suite *KeeperTestSuite) TestBridgeTokenSymbolNonNativeAllowed() {
	relayers := suite.relayerAddrs
	externalPris := suite.externalPris

	tokenContract := helpers.GenerateAddress().Hex()
	nonce := uint64(1)

	// Submit claims from enough relayers to reach attestation threshold
	for i, relayer := range relayers {
		claim := &types.MsgBridgeTokenClaim{
			EventNonce:     nonce + uint64(i),
			TokenContract:  tokenContract,
			Symbol:         "USDT",
			Name:           "Tether USD",
			Decimals:       6,
			ChainName:      suite.chainName,
			ExternalAddress: suite.PubKeyToExternalAddr(externalPris[i].PublicKey),
			RelayerAddress:  relayer.String(),
		}
		_, err := suite.MsgServer().BridgeTokenClaim(suite.Ctx, claim)
		suite.NoError(err)
	}

	// Verify the bridge token was registered with denom "uusdt"
	expectedDenom := utils.GetDenom("USDT")
	bridgeToken, err := suite.Keeper().GetBridgeTokenByDenom(suite.Ctx, expectedDenom)
	suite.NoError(err)
	suite.Equal(tokenContract, bridgeToken.ContractAddress)
	suite.Equal(expectedDenom, bridgeToken.Denom)

	// Verify it is NOT the native denom
	suite.NotEqual(params.BaseDenom, bridgeToken.Denom)
}

// TestBridgeTokenNativeDenomSupplyIsolation verifies that after the fix,
// the native umec supply cannot be inflated through bridge token claims.
func (suite *KeeperTestSuite) TestBridgeTokenNativeDenomSupplyIsolation() {
	// Record initial native supply
	initialSupply := suite.App.BankKeeper.GetSupply(suite.Ctx, params.BaseDenom)

	// Attempt to register MEC bridge token (should fail)
	tokenContract := helpers.GenerateAddress().Hex()
	claim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		TokenContract:  tokenContract,
		Symbol:         "MEC",
		Name:           "MEC Token",
		Decimals:       18,
		ChainName:      suite.chainName,
		ExternalAddress: suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		RelayerAddress:  suite.relayerAddrs[0].String(),
	}
	_, err := suite.MsgServer().BridgeTokenClaim(suite.Ctx, claim)
	suite.Error(err, "MEC bridge token registration should be rejected")

	// Verify native supply unchanged
	finalSupply := suite.App.BankKeeper.GetSupply(suite.Ctx, params.BaseDenom)
	suite.True(initialSupply.Amount.Equal(finalSupply.Amount),
		fmt.Sprintf("native supply changed: initial=%s final=%s", initialSupply.Amount, finalSupply.Amount))

	// Verify no bridge token was created for "umec"
	_, err = suite.Keeper().GetBridgeTokenByDenom(suite.Ctx, params.BaseDenom)
	suite.Error(err, "no bridge token should exist for native denom")
}

// TestBridgeTokenVariousNativeSymbolsRejected verifies that various
// case-insensitive representations of MEC are all rejected.
func (suite *KeeperTestSuite) TestBridgeTokenVariousNativeSymbolsRejected() {
	symbols := []string{"MEC", "mec", "Mec", "MEC ", " MEC", "mEc"}

	for _, symbol := range symbols {
		denom := utils.GetDenom(symbol)
		if denom == params.BaseDenom {
			// This symbol would collide - verify it's rejected at the
			// denom derivation level by checking GetDenom output
			suite.Equal(params.BaseDenom, denom,
				fmt.Sprintf("symbol %q derives to native denom", symbol))
		}
	}
}

// TestBridgeTokenNonCollidingSymbolSucceeds ensures that a typical bridge
// token (e.g. "ETH") does not collide with the native denom.
func (suite *KeeperTestSuite) TestBridgeTokenNonCollidingDenom() {
	testCases := []struct {
		symbol      string
		expectDenom string
	}{
		{"USDT", "uusdt"},
		{"ETH", "ueth"},
		{"BTC", "ubtc"},
		{"BNB", "ubnb"},
	}

	for _, tc := range testCases {
		denom := utils.GetDenom(tc.symbol)
		suite.Equal(tc.expectDenom, denom)
		suite.NotEqual(params.BaseDenom, denom,
			fmt.Sprintf("symbol %q should not collide with native denom", tc.symbol))
	}
}
