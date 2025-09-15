package keeper_test

import (
	"encoding/hex"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/st-chain/me-hub/testutil/helpers"
	fxtypes "github.com/st-chain/me-hub/types"
	"github.com/st-chain/me-hub/x/gravity/types"
	trontypes "github.com/st-chain/me-hub/x/tron/types"
)

func (suite *KeeperTestSuite) TestABCIEndBlockDepositClaim() {
	normalMsg := &types.MsgBondedGravity{
		GravityAddress:   suite.oracleAddrs[0].String(),
		BridgerAddress:   suite.bridgerAddrs[0].String(),
		ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		ValidatorAddress: suite.valAddrs[0].String(),
		DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
		ChainName:        suite.chainName,
	}
	_, err := suite.MsgServer().BondedGravity(sdk.WrapSDKContext(suite.ctx), normalMsg)
	require.NoError(suite.T(), err)

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	bridgeToken := helpers.GenerateAddress().String()
	sendToFxSendAddr := helpers.GenerateAddress().String()
	if trontypes.ModuleName == suite.chainName {
		bridgeToken = trontypes.AddressFromHex(bridgeToken)
		sendToFxSendAddr = tronaddress.PubkeyToAddress(suite.externalPris[0].PublicKey).String()
	}
	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  bridgeToken,
		Name:           "Test Token",
		Symbol:         "TEST",
		Decimals:       18,
		BridgerAddress: suite.bridgerAddrs[0].String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      suite.chainName,
	}
	_, err = suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), addBridgeTokenClaim)
	require.NoError(suite.T(), err)

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	sendToFxClaim := &types.MsgSendToFxClaim{
		EventNonce:     2,
		BlockHeight:    1001,
		TokenContract:  bridgeToken,
		Amount:         sdkmath.NewInt(1234),
		Sender:         sendToFxSendAddr,
		Receiver:       sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
		TargetIbc:      hex.EncodeToString([]byte("px/transfer/channel-0")),
		BridgerAddress: suite.bridgerAddrs[0].String(),
		ChainName:      suite.chainName,
	}
	_, err = suite.MsgServer().SendToFxClaim(sdk.WrapSDKContext(suite.ctx), sendToFxClaim)
	require.NoError(suite.T(), err)

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	allBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(sendToFxClaim.Receiver))
	denom := fmt.Sprintf("%s%s", suite.chainName, bridgeToken)
	trace, err := fxtypes.GetIbcDenomTrace(denom, addBridgeTokenClaim.ChannelIbc)
	suite.NoError(err)
	denom = trace.IBCDenom()
	require.EqualValues(suite.T(), fmt.Sprintf("%s%s", sendToFxClaim.Amount.String(), denom), allBalances.String())
}

func (suite *KeeperTestSuite) TestGravityUpdate() {
	if len(suite.oracleAddrs) < 10 {
		return
	}
	for i := 0; i < 10; i++ {
		msgBondedGravity := &types.MsgBondedGravity{
			GravityAddress:   suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		require.NoError(suite.T(), msgBondedGravity.ValidateBasic())
		_, err := suite.MsgServer().BondedGravity(sdk.WrapSDKContext(suite.ctx), msgBondedGravity)

		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		oracleSets := suite.Keeper().GetGravitySets(suite.ctx)
		require.NotNil(suite.T(), oracleSets)
		require.EqualValues(suite.T(), i+1, len(oracleSets))

		power := suite.Keeper().GetLastTotalPower(suite.ctx)
		expectPower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		require.True(suite.T(), expectPower.Equal(power))
	}

	bridgeToken := helpers.GenerateAddress().String()
	if trontypes.ModuleName == suite.chainName {
		bridgeToken = trontypes.AddressFromHex(bridgeToken)
	}

	for i := 0; i < 6; i++ {
		addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  bridgeToken,
			Name:           "Test Token",
			Symbol:         "TEST",
			Decimals:       18,
			BridgerAddress: suite.bridgerAddrs[i].String(),
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}
		_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), addBridgeTokenClaim)
		require.NoError(suite.T(), err)
		endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), endBlockBeforeAttestation)
		require.False(suite.T(), endBlockBeforeAttestation.Observed)
		require.NotNil(suite.T(), endBlockBeforeAttestation.Votes)
		require.EqualValues(suite.T(), i+1, len(endBlockBeforeAttestation.Votes))

		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), endBlockAfterAttestation)
		require.False(suite.T(), endBlockAfterAttestation.Observed)
	}

	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  bridgeToken,
		Name:           "Test Token",
		Symbol:         "TEST",
		Decimals:       18,
		BridgerAddress: suite.bridgerAddrs[6].String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      suite.chainName,
	}
	_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), addBridgeTokenClaim)
	require.NoError(suite.T(), err)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	attestation := suite.Keeper().GetAttestation(suite.ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())

	require.NotNil(suite.T(), attestation)
	require.True(suite.T(), attestation.Observed)

	proposalHandler := crosschain.NewCrosschainProposalHandler(suite.app.CrosschainKeeper)

	var newGravityList []string
	for i := 0; i < 7; i++ {
		newGravityList = append(newGravityList, suite.oracleAddrs[i].String())
	}
	err = proposalHandler(suite.ctx, &types.UpdateChainGravitysProposal{ // nolint:staticcheck
		Title:       "proposal 1: try update chain oracle power >= 30%, expect error",
		Description: "",
		Gravitys:    newGravityList,
		ChainName:   suite.chainName,
	})
	require.ErrorIs(suite.T(), types.ErrInvalid, err)

	expectTotalPower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(10)).Quo(sdk.DefaultPowerReduction)
	actualTotalPower := suite.Keeper().GetLastTotalPower(suite.ctx)
	require.True(suite.T(), expectTotalPower.Equal(actualTotalPower))

	expectMaxChangePower := types.AttestationProposalGravityChangePowerThreshold.Mul(expectTotalPower).Quo(sdkmath.NewInt(100))

	expectDeletePower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(3)).Quo(sdk.DefaultPowerReduction)
	require.EqualValues(suite.T(), fmt.Sprintf("max change power, maxChangePowerThreshold: %s, deleteTotalPower: %s: %s", expectMaxChangePower.String(), expectDeletePower.String(), types.ErrInvalid), err.Error())

	var newGravityList2 []string
	for i := 0; i < 8; i++ {
		newGravityList2 = append(newGravityList2, suite.oracleAddrs[i].String())
	}
	err = proposalHandler(suite.ctx, &types.UpdateChainGravitysProposal{ // nolint:staticcheck
		Title:       "proposal 2: try update chain oracle power <= 30%, expect success",
		Description: "",
		Gravitys:    newGravityList2,
		ChainName:   suite.chainName,
	})
	require.NoError(suite.T(), err)
}

func (suite *KeeperTestSuite) TestAttestationAfterGravityUpdate() {
	if len(suite.bridgerAddrs) < 20 {
		return
	}
	for i := 0; i < 20; i++ {
		msgBondedGravity := &types.MsgBondedGravity{
			GravityAddress:   suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		_, err := suite.MsgServer().BondedGravity(sdk.WrapSDKContext(suite.ctx), msgBondedGravity)
		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		oracleSets := suite.Keeper().GetGravitySets(suite.ctx)
		require.NotNil(suite.T(), oracleSets)
		require.EqualValues(suite.T(), i+1, len(oracleSets))

		power := suite.Keeper().GetLastTotalPower(suite.ctx)
		expectPower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		require.True(suite.T(), expectPower.Equal(power))
	}

	bridgeToken := helpers.GenerateAddress().String()
	if trontypes.ModuleName == suite.chainName {
		bridgeToken = trontypes.AddressFromHex(bridgeToken)
	}

	{
		firstBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  bridgeToken,
			Name:           "Test Token",
			Symbol:         "TEST",
			Decimals:       18,
			BridgerAddress: "",
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}

		for i := 0; i < 13; i++ {
			firstBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[i].String()
			_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), firstBridgeTokenClaim)
			require.NoError(suite.T(), err)
			endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockBeforeAttestation)
			require.False(suite.T(), endBlockBeforeAttestation.Observed)
			require.NotNil(suite.T(), endBlockBeforeAttestation.Votes)
			require.EqualValues(suite.T(), i+1, len(endBlockBeforeAttestation.Votes))

			suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
			suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
			endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockAfterAttestation)
			require.False(suite.T(), endBlockAfterAttestation.Observed)
		}

		firstBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[13].String()
		_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), firstBridgeTokenClaim)
		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		attestation := suite.Keeper().GetAttestation(suite.ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())

		require.NotNil(suite.T(), attestation)
		require.True(suite.T(), attestation.Observed)
	}

	{
		secondBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     2,
			BlockHeight:    1001,
			TokenContract:  bridgeToken,
			Name:           "Test Token2",
			Symbol:         "TEST2",
			Decimals:       18,
			BridgerAddress: "",
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}

		for i := 0; i < 6; i++ {
			secondBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[i].String()
			_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), secondBridgeTokenClaim)
			require.NoError(suite.T(), err)
			endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockBeforeAttestation)
			require.False(suite.T(), endBlockBeforeAttestation.Observed)
			require.NotNil(suite.T(), endBlockBeforeAttestation.Votes)
			require.EqualValues(suite.T(), i+1, len(endBlockBeforeAttestation.Votes))

			suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
			suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
			endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockAfterAttestation)
			require.False(suite.T(), endBlockAfterAttestation.Observed)
		}

		secondClaimAttestation := suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		var newGravityList []string
		for i := 0; i < 15; i++ {
			newGravityList = append(newGravityList, suite.oracleAddrs[i].String())
		}
		_, err := suite.MsgServer().UpdateChainGravitys(suite.ctx, &types.MsgUpdateChainGravitys{
			Gravitys:  newGravityList,
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			ChainName: suite.chainName,
		})
		require.NoError(suite.T(), err)
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		activeGravitys := suite.Keeper().GetAllGravitys(suite.ctx, true)
		require.NotNil(suite.T(), activeGravitys)
		require.EqualValues(suite.T(), 15, len(activeGravitys))
		for i := 0; i < 15; i++ {
			require.NotNil(suite.T(), newGravityList[i], activeGravitys[i].GravityAddress)
		}

		var newGravityList2 []string
		for i := 0; i < 11; i++ {
			newGravityList2 = append(newGravityList2, suite.oracleAddrs[i].String())
		}
		_, err = suite.MsgServer().UpdateChainGravitys(suite.ctx, &types.MsgUpdateChainGravitys{
			Gravitys:  newGravityList2,
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			ChainName: suite.chainName,
		})
		require.NoError(suite.T(), err)
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		activeGravitys = suite.Keeper().GetAllGravitys(suite.ctx, true)
		require.NotNil(suite.T(), activeGravitys)
		require.EqualValues(suite.T(), 11, len(activeGravitys))
		for i := 0; i < 11; i++ {
			require.NotNil(suite.T(), newGravityList2[i], activeGravitys[i].GravityAddress)
		}

		var newGravityList3 []string
		for i := 0; i < 10; i++ {
			newGravityList3 = append(newGravityList3, suite.oracleAddrs[i].String())
		}
		_, err = suite.MsgServer().UpdateChainGravitys(suite.ctx, &types.MsgUpdateChainGravitys{
			Gravitys:  newGravityList3,
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			ChainName: suite.chainName,
		})
		require.NoError(suite.T(), err)
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		activeGravitys = suite.Keeper().GetAllGravitys(suite.ctx, true)
		require.NotNil(suite.T(), activeGravitys)
		require.EqualValues(suite.T(), 10, len(activeGravitys))
		for i := 0; i < 10; i++ {
			require.NotNil(suite.T(), newGravityList3[i], activeGravitys[i].GravityAddress)
		}

		secondBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[6].String()
		_, err = suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), secondBridgeTokenClaim)
		require.NoError(suite.T(), err)

		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.True(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 7, len(secondClaimAttestation.Votes))
	}
}

func (suite *KeeperTestSuite) TestGravityDelete() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedGravity := &types.MsgBondedGravity{
			GravityAddress:   suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		require.NoError(suite.T(), msgBondedGravity.ValidateBasic())
		_, err := suite.MsgServer().BondedGravity(sdk.WrapSDKContext(suite.ctx), msgBondedGravity)
		require.NoError(suite.T(), err)
	}
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	allGravitys := suite.Keeper().GetAllGravitys(suite.ctx, false)
	require.NotNil(suite.T(), allGravitys)
	require.EqualValues(suite.T(), len(suite.oracleAddrs), len(allGravitys))

	oracle := suite.oracleAddrs[0]
	bridger := suite.bridgerAddrs[0]
	externalAddress := suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey)

	oracleAddr, found := suite.Keeper().GetGravityAddressByBridgerKey(suite.ctx, bridger)
	require.True(suite.T(), found)
	require.EqualValues(suite.T(), oracle.String(), oracleAddr.String())

	oracleAddr, found = suite.Keeper().GetGravityByExternalAddress(suite.ctx, externalAddress)
	require.True(suite.T(), found)
	require.EqualValues(suite.T(), oracle.String(), oracleAddr.String())

	oracleData, found := suite.Keeper().GetGravity(suite.ctx, oracle)
	require.True(suite.T(), found)
	require.NotNil(suite.T(), oracleData)
	require.EqualValues(suite.T(), oracle.String(), oracleData.GravityAddress)
	require.EqualValues(suite.T(), bridger.String(), oracleData.BridgerAddress)
	require.EqualValues(suite.T(), externalAddress, oracleData.ExternalAddress)

	require.True(suite.T(), sdkmath.NewInt(10*1e3).MulRaw(1e18).Equal(oracleData.DelegateAmount))

	newGravityAddressList := make([]string, 0, len(suite.oracleAddrs)-1)
	for _, address := range suite.oracleAddrs[1:] {
		newGravityAddressList = append(newGravityAddressList, address.String())
	}

	_, err := suite.MsgServer().UpdateChainGravitys(suite.ctx, &types.MsgUpdateChainGravitys{
		Gravitys:  newGravityAddressList,
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ChainName: suite.chainName,
	})
	require.NoError(suite.T(), err)
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	oracleAddr, found = suite.Keeper().GetGravityAddressByBridgerKey(suite.ctx, bridger)
	require.True(suite.T(), found)
	require.Equal(suite.T(), oracleAddr, oracle)

	oracleAddr, found = suite.Keeper().GetGravityByExternalAddress(suite.ctx, externalAddress)
	require.True(suite.T(), found)
	require.Equal(suite.T(), oracleAddr, oracle)

	oracleData, found = suite.Keeper().GetGravity(suite.ctx, oracle)
	require.True(suite.T(), found)
}

func (suite *KeeperTestSuite) TestGravitySetSlash() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedGravity := &types.MsgBondedGravity{
			GravityAddress:   suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		require.NoError(suite.T(), msgBondedGravity.ValidateBasic())
		_, err := suite.MsgServer().BondedGravity(sdk.WrapSDKContext(suite.ctx), msgBondedGravity)
		require.NoError(suite.T(), err)
	}
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.Keeper().EndBlocker(suite.ctx)

	allGravitys := suite.Keeper().GetAllGravitys(suite.ctx, false)
	require.NotNil(suite.T(), allGravitys)
	require.Equal(suite.T(), len(suite.oracleAddrs), len(allGravitys))

	oracleSets := suite.Keeper().GetGravitySets(suite.ctx)
	require.NotNil(suite.T(), oracleSets)
	require.EqualValues(suite.T(), 1, len(oracleSets))

	for i := 0; i < len(suite.oracleAddrs)-1; i++ {
		externalAddress, signature := suite.SignGravitySetConfirm(suite.externalPris[i], oracleSets[0])
		oracleSetConfirm := &types.MsgGravitySetConfirm{
			Nonce:           oracleSets[0].Nonce,
			BridgerAddress:  suite.bridgerAddrs[i].String(),
			ExternalAddress: externalAddress,
			Signature:       hex.EncodeToString(signature),
			ChainName:       suite.chainName,
		}
		require.NoError(suite.T(), oracleSetConfirm.ValidateBasic())
		_, err := suite.MsgServer().GravitySetConfirm(sdk.WrapSDKContext(suite.ctx), oracleSetConfirm)
		require.NoError(suite.T(), err)
	}

	suite.Keeper().EndBlocker(suite.ctx)
	oracleSetHeight := int64(oracleSets[0].Height)
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	oracle, found := suite.Keeper().GetGravity(suite.ctx, suite.oracleAddrs[len(suite.oracleAddrs)-1])
	require.True(suite.T(), found)
	require.True(suite.T(), oracle.Online)
	require.Equal(suite.T(), int64(0), oracle.SlashTimes)

	suite.ctx = suite.ctx.WithBlockHeight(oracleSetHeight + int64(suite.Keeper().GetParams(suite.ctx).SignedWindow) + 1)
	suite.Keeper().EndBlocker(suite.ctx)

	oracle, found = suite.Keeper().GetGravity(suite.ctx, suite.oracleAddrs[len(suite.oracleAddrs)-1])
	require.True(suite.T(), found)
	require.False(suite.T(), oracle.Online)
	require.Equal(suite.T(), int64(1), oracle.SlashTimes)
}

func (suite *KeeperTestSuite) TestSlashGravity() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedGravity := &types.MsgBondedGravity{
			GravityAddress:   suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		require.NoError(suite.T(), msgBondedGravity.ValidateBasic())
		_, err := suite.MsgServer().BondedGravity(sdk.WrapSDKContext(suite.ctx), msgBondedGravity)
		require.NoError(suite.T(), err)
	}

	params := suite.Keeper().GetParams(suite.ctx)
	err := suite.Keeper().SetParams(suite.ctx, &params)
	suite.Require().NoError(err)
	for i := 0; i < len(suite.oracleAddrs); i++ {
		oracle, found := suite.Keeper().GetGravity(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.True(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(0), oracle.SlashTimes)

		suite.Keeper().SlashGravity(suite.ctx, oracle.GravityAddress)

		oracle, found = suite.Keeper().GetGravity(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)
	}

	// repeat slash test.
	for i := 0; i < len(suite.oracleAddrs); i++ {
		oracle, found := suite.Keeper().GetGravity(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)

		suite.Keeper().SlashGravity(suite.ctx, oracle.GravityAddress)

		oracle, found = suite.Keeper().GetGravity(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)
	}
}
