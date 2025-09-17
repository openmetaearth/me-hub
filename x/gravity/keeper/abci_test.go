package keeper_test

import (
	"encoding/hex"
	"fmt"
	"github.com/st-chain/me-hub/utils"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/st-chain/me-hub/testutil/helpers"
	"github.com/st-chain/me-hub/x/gravity/types"
)

func (suite *KeeperTestSuite) TestABCIEndBlockDepositClaim() {
	normalMsg := &types.MsgBondedRelayer{
		RelayerAddress:  suite.oracleAddrs[0].String(),
		ExternalAddress: suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		DelegateAmount:  types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
		ChainName:       suite.chainName,
	}
	_, err := suite.MsgServer().BondedRelayer(sdk.WrapSDKContext(suite.ctx), normalMsg)
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
	trace, err := utils.GetIbcDenomTrace(denom, addBridgeTokenClaim.ChannelIbc)
	suite.NoError(err)
	denom = trace.IBCDenom()
	require.EqualValues(suite.T(), fmt.Sprintf("%s%s", sendToFxClaim.Amount.String(), denom), allBalances.String())
}

func (suite *KeeperTestSuite) TestRelayerUpdate() {
	if len(suite.oracleAddrs) < 10 {
		return
	}
	for i := 0; i < 10; i++ {
		msgBondedRelayer := &types.MsgBondedRelayer{
			RelayerAddress:  suite.oracleAddrs[i].String(),
			ExternalAddress: suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			DelegateAmount:  types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:       suite.chainName,
		}
		require.NoError(suite.T(), msgBondedRelayer.ValidateBasic())
		_, err := suite.MsgServer().BondedRelayer(sdk.WrapSDKContext(suite.ctx), msgBondedRelayer)

		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		oracleSets := suite.Keeper().GetRelayerSets(suite.ctx)
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

	var newRelayerList []string
	for i := 0; i < 7; i++ {
		newRelayerList = append(newRelayerList, suite.oracleAddrs[i].String())
	}
	err = proposalHandler(suite.ctx, &types.UpdateChainRelayersProposal{ // nolint:staticcheck
		Title:       "proposal 1: try update chain oracle power >= 30%, expect error",
		Description: "",
		Relayers:    newRelayerList,
		ChainName:   suite.chainName,
	})
	require.ErrorIs(suite.T(), types.ErrInvalid, err)

	expectTotalPower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(10)).Quo(sdk.DefaultPowerReduction)
	actualTotalPower := suite.Keeper().GetLastTotalPower(suite.ctx)
	require.True(suite.T(), expectTotalPower.Equal(actualTotalPower))

	expectMaxChangePower := types.AttestationProposalRelayerChangePowerThreshold.Mul(expectTotalPower).Quo(sdkmath.NewInt(100))

	expectDeletePower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(3)).Quo(sdk.DefaultPowerReduction)
	require.EqualValues(suite.T(), fmt.Sprintf("max change power, maxChangePowerThreshold: %s, deleteTotalPower: %s: %s", expectMaxChangePower.String(), expectDeletePower.String(), types.ErrInvalid), err.Error())

	var newRelayerList2 []string
	for i := 0; i < 8; i++ {
		newRelayerList2 = append(newRelayerList2, suite.oracleAddrs[i].String())
	}
	err = proposalHandler(suite.ctx, &types.UpdateChainRelayersProposal{ // nolint:staticcheck
		Title:       "proposal 2: try update chain oracle power <= 30%, expect success",
		Description: "",
		Relayers:    newRelayerList2,
		ChainName:   suite.chainName,
	})
	require.NoError(suite.T(), err)
}

func (suite *KeeperTestSuite) TestAttestationAfterRelayerUpdate() {
	if len(suite.bridgerAddrs) < 20 {
		return
	}
	for i := 0; i < 20; i++ {
		msgBondedRelayer := &types.MsgBondedRelayer{
			RelayerAddress:   suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		_, err := suite.MsgServer().BondedRelayer(sdk.WrapSDKContext(suite.ctx), msgBondedRelayer)
		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		oracleSets := suite.Keeper().GetRelayerSets(suite.ctx)
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

		var newRelayerList []string
		for i := 0; i < 15; i++ {
			newRelayerList = append(newRelayerList, suite.oracleAddrs[i].String())
		}
		_, err := suite.MsgServer().UpdateChainRelayers(suite.ctx, &types.MsgUpdateChainRelayers{
			Relayers:  newRelayerList,
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

		activeRelayers := suite.Keeper().GetAllRelayers(suite.ctx, true)
		require.NotNil(suite.T(), activeRelayers)
		require.EqualValues(suite.T(), 15, len(activeRelayers))
		for i := 0; i < 15; i++ {
			require.NotNil(suite.T(), newRelayerList[i], activeRelayers[i].RelayerAddress)
		}

		var newRelayerList2 []string
		for i := 0; i < 11; i++ {
			newRelayerList2 = append(newRelayerList2, suite.oracleAddrs[i].String())
		}
		_, err = suite.MsgServer().UpdateChainRelayers(suite.ctx, &types.MsgUpdateChainRelayers{
			Relayers:  newRelayerList2,
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

		activeRelayers = suite.Keeper().GetAllRelayers(suite.ctx, true)
		require.NotNil(suite.T(), activeRelayers)
		require.EqualValues(suite.T(), 11, len(activeRelayers))
		for i := 0; i < 11; i++ {
			require.NotNil(suite.T(), newRelayerList2[i], activeRelayers[i].RelayerAddress)
		}

		var newRelayerList3 []string
		for i := 0; i < 10; i++ {
			newRelayerList3 = append(newRelayerList3, suite.oracleAddrs[i].String())
		}
		_, err = suite.MsgServer().UpdateChainRelayers(suite.ctx, &types.MsgUpdateChainRelayers{
			Relayers:  newRelayerList3,
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

		activeRelayers = suite.Keeper().GetAllRelayers(suite.ctx, true)
		require.NotNil(suite.T(), activeRelayers)
		require.EqualValues(suite.T(), 10, len(activeRelayers))
		for i := 0; i < 10; i++ {
			require.NotNil(suite.T(), newRelayerList3[i], activeRelayers[i].RelayerAddress)
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

func (suite *KeeperTestSuite) TestRelayerDelete() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedRelayer := &types.MsgBondedRelayer{
			RelayerAddress:   suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		require.NoError(suite.T(), msgBondedRelayer.ValidateBasic())
		_, err := suite.MsgServer().BondedRelayer(sdk.WrapSDKContext(suite.ctx), msgBondedRelayer)
		require.NoError(suite.T(), err)
	}
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	allRelayers := suite.Keeper().GetAllRelayers(suite.ctx, false)
	require.NotNil(suite.T(), allRelayers)
	require.EqualValues(suite.T(), len(suite.oracleAddrs), len(allRelayers))

	oracle := suite.oracleAddrs[0]
	bridger := suite.bridgerAddrs[0]
	externalAddress := suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey)

	oracleAddr, found := suite.Keeper().GetRelayerAddressByBridgerKey(suite.ctx, bridger)
	require.True(suite.T(), found)
	require.EqualValues(suite.T(), oracle.String(), oracleAddr.String())

	oracleAddr, found = suite.Keeper().GetRelayerByExternalAddress(suite.ctx, externalAddress)
	require.True(suite.T(), found)
	require.EqualValues(suite.T(), oracle.String(), oracleAddr.String())

	oracleData, found := suite.Keeper().GetRelayer(suite.ctx, oracle)
	require.True(suite.T(), found)
	require.NotNil(suite.T(), oracleData)
	require.EqualValues(suite.T(), oracle.String(), oracleData.RelayerAddress)
	require.EqualValues(suite.T(), bridger.String(), oracleData.BridgerAddress)
	require.EqualValues(suite.T(), externalAddress, oracleData.ExternalAddress)

	require.True(suite.T(), sdkmath.NewInt(10*1e3).MulRaw(1e18).Equal(oracleData.DelegateAmount))

	newRelayerAddressList := make([]string, 0, len(suite.oracleAddrs)-1)
	for _, address := range suite.oracleAddrs[1:] {
		newRelayerAddressList = append(newRelayerAddressList, address.String())
	}

	_, err := suite.MsgServer().UpdateChainRelayers(suite.ctx, &types.MsgUpdateChainRelayers{
		Relayers:  newRelayerAddressList,
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ChainName: suite.chainName,
	})
	require.NoError(suite.T(), err)
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	oracleAddr, found = suite.Keeper().GetRelayerAddressByBridgerKey(suite.ctx, bridger)
	require.True(suite.T(), found)
	require.Equal(suite.T(), oracleAddr, oracle)

	oracleAddr, found = suite.Keeper().GetRelayerByExternalAddress(suite.ctx, externalAddress)
	require.True(suite.T(), found)
	require.Equal(suite.T(), oracleAddr, oracle)

	oracleData, found = suite.Keeper().GetRelayer(suite.ctx, oracle)
	require.True(suite.T(), found)
}

func (suite *KeeperTestSuite) TestRelayerSetSlash() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedRelayer := &types.MsgBondedRelayer{
			RelayerAddress:   suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		require.NoError(suite.T(), msgBondedRelayer.ValidateBasic())
		_, err := suite.MsgServer().BondedRelayer(sdk.WrapSDKContext(suite.ctx), msgBondedRelayer)
		require.NoError(suite.T(), err)
	}
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.Keeper().EndBlocker(suite.ctx)

	allRelayers := suite.Keeper().GetAllRelayers(suite.ctx, false)
	require.NotNil(suite.T(), allRelayers)
	require.Equal(suite.T(), len(suite.oracleAddrs), len(allRelayers))

	oracleSets := suite.Keeper().GetRelayerSets(suite.ctx)
	require.NotNil(suite.T(), oracleSets)
	require.EqualValues(suite.T(), 1, len(oracleSets))

	for i := 0; i < len(suite.oracleAddrs)-1; i++ {
		externalAddress, signature := suite.SignRelayerSetConfirm(suite.externalPris[i], oracleSets[0])
		oracleSetConfirm := &types.MsgRelayerSetConfirm{
			Nonce:           oracleSets[0].Nonce,
			BridgerAddress:  suite.bridgerAddrs[i].String(),
			ExternalAddress: externalAddress,
			Signature:       hex.EncodeToString(signature),
			ChainName:       suite.chainName,
		}
		require.NoError(suite.T(), oracleSetConfirm.ValidateBasic())
		_, err := suite.MsgServer().RelayerSetConfirm(sdk.WrapSDKContext(suite.ctx), oracleSetConfirm)
		require.NoError(suite.T(), err)
	}

	suite.Keeper().EndBlocker(suite.ctx)
	oracleSetHeight := int64(oracleSets[0].Height)
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	oracle, found := suite.Keeper().GetRelayer(suite.ctx, suite.oracleAddrs[len(suite.oracleAddrs)-1])
	require.True(suite.T(), found)
	require.True(suite.T(), oracle.Online)
	require.Equal(suite.T(), int64(0), oracle.SlashTimes)

	suite.ctx = suite.ctx.WithBlockHeight(oracleSetHeight + int64(suite.Keeper().GetParams(suite.ctx).SignedWindow) + 1)
	suite.Keeper().EndBlocker(suite.ctx)

	oracle, found = suite.Keeper().GetRelayer(suite.ctx, suite.oracleAddrs[len(suite.oracleAddrs)-1])
	require.True(suite.T(), found)
	require.False(suite.T(), oracle.Online)
	require.Equal(suite.T(), int64(1), oracle.SlashTimes)
}

func (suite *KeeperTestSuite) TestSlashRelayer() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedRelayer := &types.MsgBondedRelayer{
			RelayerAddress:   suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		require.NoError(suite.T(), msgBondedRelayer.ValidateBasic())
		_, err := suite.MsgServer().BondedRelayer(sdk.WrapSDKContext(suite.ctx), msgBondedRelayer)
		require.NoError(suite.T(), err)
	}

	params := suite.Keeper().GetParams(suite.ctx)
	err := suite.Keeper().SetParams(suite.ctx, &params)
	suite.Require().NoError(err)
	for i := 0; i < len(suite.oracleAddrs); i++ {
		oracle, found := suite.Keeper().GetRelayer(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.True(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(0), oracle.SlashTimes)

		suite.Keeper().SlashRelayer(suite.ctx, oracle.RelayerAddress)

		oracle, found = suite.Keeper().GetRelayer(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)
	}

	// repeat slash test.
	for i := 0; i < len(suite.oracleAddrs); i++ {
		oracle, found := suite.Keeper().GetRelayer(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)

		suite.Keeper().SlashRelayer(suite.ctx, oracle.RelayerAddress)

		oracle, found = suite.Keeper().GetRelayer(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)
	}
}
